package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildCommand_CreatesDistDirectory(t *testing.T) {
	// Arrange
	distDir := "dist"
	os.RemoveAll(distDir) // Clean up before the test

	// Act
	build()

	// Assert
	assert.DirExists(t, distDir)

	// Clean up after the test
	os.RemoveAll(distDir)
}

func TestBuildCommand_CopiesPublicDirectory(t *testing.T) {
	// Arrange
	distDir := "dist"
	publicDir := "public"
	os.RemoveAll(distDir) // Clean up before the test

	// Create a dummy public directory with a file
	os.MkdirAll(publicDir, 0755)
	os.WriteFile("public/style.css", []byte("body {}"), 0644)

	// Act
	build()

	// Assert
	assert.FileExists(t, "dist/style.css")

	// Clean up after the test
	os.RemoveAll(distDir)
	os.RemoveAll(publicDir)
}
