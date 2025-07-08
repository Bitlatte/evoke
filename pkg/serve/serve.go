package serve

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Bitlatte/evoke/pkg/build"
	"github.com/fsnotify/fsnotify"
)

// Serve starts a web server and watches for changes.
func Serve(port int) error {
	if err := build.Build(); err != nil {
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

func startServer(port int) {
	portStr := strconv.Itoa(port)
	log.Printf("Starting server on :%s", portStr)
	http.Handle("/", http.FileServer(http.Dir("dist")))
	if err := http.ListenAndServe(":"+portStr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func watchFiles(watcher *fsnotify.Watcher) {
	var lastEventTime time.Time
	var lastEventName string

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Debounce events
			if event.Name == lastEventName && time.Since(lastEventTime) < 2*time.Second {
				continue
			}
			lastEventName = event.Name
			lastEventTime = time.Now()

			log.Printf("Change detected in %s, rebuilding...", event.Name)
			if err := build.Build(); err != nil {
				log.Printf("Error rebuilding site: %v", err)
			} else {
				log.Println("Site rebuilt successfully.")
				// A simple restart is not possible, so we just rebuild.
				// A more advanced implementation would restart the server process.
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
