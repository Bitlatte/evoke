package cache

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"sync"
)

// Cache represents the cache
type Cache struct {
	path  string
	Store map[string]string
	mu    sync.RWMutex
}

// New creates a new cache
func New(path string) (*Cache, error) {
	c := &Cache{
		path:  path,
		Store: make(map[string]string),
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err
		}
		if _, err := os.Create(path); err != nil {
			return nil, err
		}
	}

	if err := c.Load(); err != nil {
		return nil, err
	}

	return c, nil
}

// Load loads the cache from disk
func (c *Cache) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Not an error if the file doesn't exist yet
		}
		return err
	}
	defer file.Close()

	// Check if the file is empty
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.Size() == 0 {
		return nil
	}

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&c.Store)
}

// Save saves the cache to disk
func (c *Cache) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	file, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(c.Store)
}

// Get returns the hash for the given path
func (c *Cache) Get(path string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Store[path]
}

// Set sets the hash for the given path
func (c *Cache) Set(path, hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Store[path] = hash
}
