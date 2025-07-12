package util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ToOutputPath converts a content path to an output path, removing any route
// groups.
func ToOutputPath(path string) string {
	parts := strings.Split(path, string(filepath.Separator))
	var newParts []string
	for _, part := range parts {
		if len(part) > 2 && part[0] == '(' && part[len(part)-1] == ')' {
			continue
		}
		newParts = append(newParts, part)
	}
	return filepath.Join(newParts...)[len("content/"):]
}

// CopyDirectory copies a directory from src to dest.
func CopyDirectory(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a new path in the destination directory
		newPath := filepath.Join(dest, path[len(src):])

		if info.IsDir() {
			os.MkdirAll(newPath, info.Mode())
		} else {
			if err := CopyFile(path, newPath); err != nil {
				return err
			}
		}

		return nil
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

	_, err = io.Copy(destFile, sourceFile)
	return err
}
