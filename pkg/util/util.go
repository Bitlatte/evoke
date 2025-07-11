// Package util provides utility functions for the evoke static site generator.
package util

import (
	"io"
	"os"
	"path/filepath"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 32*1024)
	},
}

// CopyDirectory copies a directory from src to dest.
func CopyDirectory(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a parallel structure in the destination
		destPath := filepath.Join(dest, path[len(src):])

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy the file
		return CopyFile(path, destPath)
	})
}

// CopyFile copies a file from src to dest.
func CopyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	_, err = io.CopyBuffer(destFile, sourceFile, buf)
	return err
}
