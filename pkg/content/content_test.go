package content_test

import (
	"os"
	"testing"

	"github.com/Bitlatte/evoke/pkg/content"
	"github.com/Bitlatte/evoke/pkg/templates"
	"github.com/stretchr/testify/assert"
)

func TestProcessHTML_CreatesFileInDist(t *testing.T) {
	// Arrange
	distDir := "dist"
	contentDir := "content"
	templatesDir := "templates"
	os.RemoveAll(distDir)
	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(templatesDir, 0755)

	// Create a dummy content file
	os.WriteFile("content/index.html", []byte("<h1>{{.Title}}</h1>"), 0644)
	// Create a dummy template file
	os.WriteFile("templates/base.html", []byte("<html><body>{{.Content}}</body></html>"), 0644)

	config := map[string]interface{}{"Title": "My Test Site"}
	loadedTemplates, _ := templates.LoadTemplates()

	// Act
	err := content.ProcessHTML("content/index.html", config, loadedTemplates)

	// Assert
	assert.NoError(t, err)
	assert.FileExists(t, "dist/index.html")

	// Verify the content of the created file
	content, err := os.ReadFile("dist/index.html")
	assert.NoError(t, err)
	assert.Contains(t, string(content), "<h1>My Test Site</h1>")

	// Clean up
	os.RemoveAll(distDir)
	os.RemoveAll(contentDir)
	os.RemoveAll(templatesDir)
}

func TestProcessMarkdown_CreatesFileInDist(t *testing.T) {
	// Arrange
	distDir := "dist"
	contentDir := "content"
	templatesDir := "templates"
	os.RemoveAll(distDir)
	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(templatesDir, 0755)

	// Create a dummy content file
	os.WriteFile("content/post.md", []byte("# My Post"), 0644)
	// Create a dummy template file
	os.WriteFile("templates/post.html", []byte("<html><body>{{.Content}}</body></html>"), 0644)

	config := map[string]interface{}{}
	loadedTemplates, _ := templates.LoadTemplates()

	// Act
	err := content.ProcessMarkdown("content/post.md", config, loadedTemplates)

	// Assert
	assert.NoError(t, err)
	assert.FileExists(t, "dist/post.html")

	// Verify the content of the created file
	content, err := os.ReadFile("dist/post.html")
	assert.NoError(t, err)
	assert.Contains(t, string(content), "<h1>My Post</h1>")

	// Clean up
	os.RemoveAll(distDir)
	os.RemoveAll(contentDir)
	os.RemoveAll(templatesDir)
}
