package util_test

import (
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
	if err := os.WriteFile(srcFile, []byte("hello"), 0644); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()

	for b.Loop() {
		if err := util.CopyFile(srcFile, destFile); err != nil {
			b.Fatal(err)
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
	if err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755); err != nil {
		b.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("root"), 0644); err != nil {
		b.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "subdir", "file.txt"), []byte("nested"), 0644); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()

	for b.Loop() {
		if err := util.CopyDirectory(srcDir, destDir); err != nil {
			b.Fatal(err)
		}
	}
}
