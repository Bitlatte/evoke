package util

import (
	"io"
	"os"
	"path/filepath"
)

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

	_, err = io.Copy(destFile, sourceFile)
	return err
}
