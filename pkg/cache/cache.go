package cache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Cache represents the cache
type Cache struct {
	path  string
	Store map[string]string
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
	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, &c.Store)
}

// Save saves the cache to disk
func (c *Cache) Save() error {
	data, err := json.Marshal(c.Store)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.path, data, 0644)
}

// Get returns the hash for the given path
func (c *Cache) Get(path string) string {
	return c.Store[path]
}

// Set sets the hash for the given path
func (c *Cache) Set(path, hash string) {
	c.Store[path] = hash
}
