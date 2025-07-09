package util_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Bitlatte/evoke/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "util-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src.txt")
	destFile := filepath.Join(tmpDir, "dest.txt")
	os.WriteFile(srcFile, []byte("hello"), 0644)

	// Act
	err = util.CopyFile(srcFile, destFile)

	// Assert
	assert.NoError(t, err)
	assert.FileExists(t, destFile)
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(content))
}

func TestCopyDirectory(t *testing.T) {
	// Arrange
	tmpDir, err := os.MkdirTemp("", "util-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	destDir := filepath.Join(tmpDir, "dest")
	os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("root"), 0644)
	os.WriteFile(filepath.Join(srcDir, "subdir", "file.txt"), []byte("nested"), 0644)

	// Act
	err = util.CopyDirectory(srcDir, destDir)

	// Assert
	assert.NoError(t, err)
	assert.FileExists(t, filepath.Join(destDir, "file.txt"))
	assert.FileExists(t, filepath.Join(destDir, "subdir", "file.txt"))

	rootContent, err := os.ReadFile(filepath.Join(destDir, "file.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "root", string(rootContent))

	nestedContent, err := os.ReadFile(filepath.Join(destDir, "subdir", "file.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "nested", string(nestedContent))
}

func BenchmarkCopyFile(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "util-benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src.txt")
	destFile := filepath.Join(tmpDir, "dest.txt")
	// Create a 1MB file to make the benchmark meaningful
	data := make([]byte, 1024*1024)
	if err := os.WriteFile(srcFile, data, 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if err := util.CopyFile(srcFile, destFile); err != nil {
			b.Fatal(err)
		}
	}
}

func generateBenchmarkDirectory(b *testing.B, path string, numFiles, numDirs int) {
	if err := os.MkdirAll(path, 0755); err != nil {
		b.Fatal(err)
	}

	// Create files in the current directory
	for i := 0; i < numFiles; i++ {
		filePath := filepath.Join(path, fmt.Sprintf("file-%d.txt", i))
		data := make([]byte, 1024) // 1KB files
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			b.Fatal(err)
		}
	}

	// Create subdirectories
	if numDirs > 0 {
		for i := 0; i < 2; i++ { // Create a couple of subdirs
			subDirPath := filepath.Join(path, fmt.Sprintf("subdir-%d", i))
			generateBenchmarkDirectory(b, subDirPath, numFiles, numDirs-1)
		}
	}
}

func BenchmarkCopyDirectory(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "util-benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	destDir := filepath.Join(tmpDir, "dest")

	// Generate a directory with 10 files and 2 levels of subdirectories
	generateBenchmarkDirectory(b, srcDir, 10, 2)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Clean the destination directory for each iteration
		os.RemoveAll(destDir)
		if err := util.CopyDirectory(srcDir, destDir); err != nil {
			b.Fatal(err)
		}
	}
}
