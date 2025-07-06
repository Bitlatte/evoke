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
