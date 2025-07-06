package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTemplates_LoadsAllTemplates(t *testing.T) {
	// Arrange
	templatesDir := "templates"
	os.MkdirAll(templatesDir, 0755)
	os.WriteFile("templates/base.html", []byte("{{.Title}}"), 0644)
	os.WriteFile("templates/post.html", []byte("{{.Content}}"), 0644)

	// Act
	templates, err := loadTemplates()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, templates.Lookup("base.html"))
	assert.NotNil(t, templates.Lookup("post.html"))

	// Clean up
	os.RemoveAll(templatesDir)
}
