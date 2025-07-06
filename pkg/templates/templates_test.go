package templates_test

import (
	"os"
	"testing"

	"github.com/Bitlatte/evoke/pkg/templates"
	"github.com/stretchr/testify/assert"
)

func TestLoadTemplates_LoadsAllTemplates(t *testing.T) {
	// Arrange
	templatesDir := "templates"
	os.MkdirAll(templatesDir, 0755)
	os.WriteFile("templates/base.html", []byte("{{.Title}}"), 0644)
	os.WriteFile("templates/post.html", []byte("{{.Content}}"), 0644)

	// Act
	loadedTemplates, err := templates.LoadTemplates()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, loadedTemplates.Lookup("base.html"))
	assert.NotNil(t, loadedTemplates.Lookup("post.html"))

	// Clean up
	os.RemoveAll(templatesDir)
}
