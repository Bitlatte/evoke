// Package serve provides the functionality to serve the site and watch for changes.
package serve

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Bitlatte/evoke/pkg/build"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

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
			log.Printf("Warning: could not watch %s: %v", item, err)
		}
	}

	log.Printf("Watching for changes...")
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
	log.Printf("Starting server on :%s", portStr)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", memoryFileServer)
	if err := http.ListenAndServe(":"+portStr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// wsHandler handles websocket connections.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	defer func() {
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
func broadcast(message []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
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

	memoryFSMutex.RLock()
	content, ok := memoryFS[path]
	memoryFSMutex.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	// Inject reload script
	if strings.HasSuffix(path, ".html") {
		content = bytes.Replace(content, []byte("</body>"), []byte("<script>var conn=new WebSocket(\"ws://\"+location.host+\"/ws\");conn.onmessage=function(e){if(e.data===\"reload\"){location.reload()}};</script></body>"), 1)
	}

	w.Header().Set("Content-Type", http.DetectContentType(content))
	w.Write(content)
}

// buildAndCache builds the site and caches it in memory.
func buildAndCache() error {
	buildMutex.Lock()
	defer buildMutex.Unlock()

	tempDir, err := os.MkdirTemp("", "evoke-build-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if err := build.Build(tempDir); err != nil {
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
	var timer *time.Timer
	var lastBuildTime time.Time
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			isRelevantEvent := event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename

			if !isRelevantEvent {
				continue
			}

			if time.Since(lastBuildTime) < 2*time.Second {
				continue
			}
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(1*time.Second, func() {
				log.Printf("Change detected in %s, rebuilding...", event.Name)
				if err := buildAndCache(); err != nil {
					log.Printf("Error rebuilding site: %v", err)
				} else {
					log.Println("Site rebuilt successfully.")
					lastBuildTime = time.Now()
					broadcast([]byte("reload"))
				}
			})
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
