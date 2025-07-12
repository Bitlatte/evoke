// Package serve provides the functionality to serve the site and watch for changes.
package serve

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/Bitlatte/evoke/pkg/build"

	"github.com/Bitlatte/evoke/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

type wsMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	buildMutex    sync.Mutex
	memoryFS      = make(map[string][]byte)
	memoryFSMutex sync.RWMutex
	upgrader      = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

//go:embed devtools.js
var devtoolsJS []byte

// Serve starts a web server and watches for changes.
func Serve(port int) error {
	if err := buildAndCache(); err != nil {
		return fmt.Errorf("error building site: %w", err)
	}

	go startServer(port)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}
	defer watcher.Close()

	go watchFiles(watcher)

	// Add directories and files to watch
	if err := watchRecursive(watcher, "content"); err != nil {
		return fmt.Errorf("error watching content directory: %w", err)
	}
	optionalWatch := []string{"public", "plugins", "partials", "evoke.yaml"}
	for _, item := range optionalWatch {
		if err := watcher.Add(item); err != nil {
			logger.Logger.Warn("Could not watch", "item", item, "error", err)
		}
	}

	logger.Logger.Debug("Watching for changes...")
	<-make(chan bool) // Block forever

	return nil
}

// watchRecursive recursively watches a directory for changes.
func watchRecursive(watcher *fsnotify.Watcher, path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := watcher.Add(path); err != nil {
				return fmt.Errorf("error watching path %s: %w", path, err)
			}
		}
		return nil
	})
}

// startServer starts the web server.
func startServer(port int) {
	portStr := strconv.Itoa(port)
	logger.Logger.Debug("Starting server", "port", portStr)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", memoryFileServer)
	if err := http.ListenAndServe(":"+portStr, nil); err != nil {
		logger.Logger.Fatal("Server failed", "error", err)
	}
}

// wsHandler handles websocket connections.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Logger.Error("Failed to upgrade connection", "error", err)
		return
	}
	logger.Logger.Debug("WebSocket client connected", "remoteAddr", conn.RemoteAddr())

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	defer func() {
		logger.Logger.Debug("WebSocket client disconnected", "remoteAddr", conn.RemoteAddr())
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

// broadcast sends a message to all connected websocket clients.
func broadcast(messageType string, data interface{}) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	msg := wsMessage{
		Type: messageType,
		Data: data,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Failed to marshal websocket message", "error", err)
		return
	}

	logger.Logger.Debug("Broadcasting message to clients", "clients", len(clients), "type", messageType)
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, payload); err != nil {
			logger.Logger.Error("Failed to write message to client", "error", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// memoryFileServer serves files from memory.
func memoryFileServer(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	path = strings.TrimPrefix(path, "/")

	logger.Logger.Debug("Serving file", "path", path)

	if path == "devtools.js" {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(devtoolsJS)
		return
	}

	memoryFSMutex.RLock()
	content, ok := memoryFS[path]
	memoryFSMutex.RUnlock()

	if !ok {
		logger.Logger.Warn("File not found", "path", path)
		http.NotFound(w, r)
		return
	}

	// Inject reload script
	if strings.HasSuffix(path, ".html") {
		content = bytes.Replace(content, []byte("</body>"), []byte("<script src=\"/devtools.js\"></script></body>"), 1)
	}

	w.Header().Set("Content-Type", http.DetectContentType(content))
	w.Write(content)
}

// buildAndCache builds the site and caches it in memory.
func buildAndCache() error {
	buildMutex.Lock()
	defer buildMutex.Unlock()

	logger.Logger.Debug("Building and caching site...")
	tempDir, err := os.MkdirTemp("", "evoke-build-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if err := build.Build(tempDir, false, runtime.NumCPU()); err != nil {
		return err
	}

	memoryFSMutex.Lock()
	defer memoryFSMutex.Unlock()
	memoryFS = make(map[string][]byte)

	return filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(tempDir, path)
			if err != nil {
				return err
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			memoryFS[relPath] = content
		}
		return nil
	})
}

// watchFiles watches for file changes and rebuilds the site.
func watchFiles(watcher *fsnotify.Watcher) {
	var (
		buildTicker = time.NewTicker(1 * time.Second)
		buildEvents = make(map[string]fsnotify.Event)
		mu          sync.Mutex
	)
	defer buildTicker.Stop()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			logger.Logger.Debug("Received event", "event", event)

			// Filter out events that are not relevant
			isRelevant := event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename
			if !isRelevant {
				continue
			}

			mu.Lock()
			buildEvents[event.Name] = event
			mu.Unlock()

		case <-buildTicker.C:
			mu.Lock()
			if len(buildEvents) == 0 {
				mu.Unlock()
				continue
			}
			logger.Logger.Info("Change detected, rebuilding...", "files", len(buildEvents))
			buildEvents = make(map[string]fsnotify.Event)
			mu.Unlock()

			// Check if only CSS files have changed
			onlyCSS := true
			cssFiles := []string{}
			for name := range buildEvents {
				if strings.HasSuffix(name, ".css") {
					cssFiles = append(cssFiles, name)
				} else {
					onlyCSS = false
					break
				}
			}

			if onlyCSS {
				logger.Logger.Info("CSS change detected, injecting new styles...")
				for _, cssFile := range cssFiles {
					content, err := os.ReadFile(cssFile)
					if err != nil {
						logger.Logger.Error("Failed to read CSS file", "file", cssFile, "error", err)
						continue
					}
					relPath, err := filepath.Rel("content", cssFile)
					if err != nil {
						logger.Logger.Error("Failed to get relative path", "file", cssFile, "error", err)
						continue
					}
					broadcast("css-update", map[string]string{
						"path":    "/" + relPath,
						"content": string(content),
					})
				}
			} else {
				if err := buildAndCache(); err != nil {
					logger.Logger.Error("Error rebuilding site", "error", err)
					broadcast("error", err.Error())
				} else {
					logger.Logger.Info("Site rebuilt successfully.")
					broadcast("reload", nil)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Logger.Error("Watcher error", "error", err)
		}
	}
}
