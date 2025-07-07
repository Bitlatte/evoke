package content_test

import (
	"os"
	"testing"

	"github.com/Bitlatte/evoke/pkg/content"
	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/stretchr/testify/assert"
)

func TestProcessHTML_WithLayout(t *testing.T) {
	// Arrange
	distDir := "dist"
	contentDir := "content"
	partialsDir := "partials"
	os.RemoveAll(distDir)
	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(partialsDir, 0755)

	// Create a dummy content file
	os.WriteFile("content/index.html", []byte("<h1>Hello</h1>"), 0644)
	// Create a dummy layout file
	os.WriteFile("content/_layout.html", []byte("<html><body>{{.Content}}</body></html>"), 0644)

	config := map[string]interface{}{}
	loadedPartials, _ := partials.LoadPartials()
	contentProcessor, err := content.New(config, loadedPartials)
	assert.NoError(t, err)

	// Act
	err = contentProcessor.ProcessHTML("content/index.html")

	// Assert
	assert.NoError(t, err)
	assert.FileExists(t, "dist/index.html")

	// Verify the content of the created file
	content, err := os.ReadFile("dist/index.html")
	assert.NoError(t, err)
	assert.Contains(t, string(content), "<html><body><h1>Hello</h1></body></html>")

	// Clean up
	os.RemoveAll(distDir)
	os.RemoveAll(contentDir)
	os.RemoveAll(partialsDir)
}

func TestProcessMarkdown_WithLayout(t *testing.T) {
	// Arrange
	distDir := "dist"
	contentDir := "content"
	partialsDir := "partials"
	os.RemoveAll(distDir)
	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(partialsDir, 0755)

	// Create a dummy content file
	os.WriteFile("content/post.md", []byte("# My Post"), 0644)
	// Create a dummy layout file
	os.WriteFile("content/_layout.html", []byte("<html><body>{{.Content}}</body></html>"), 0644)

	config := map[string]interface{}{}
	loadedPartials, _ := partials.LoadPartials()
	contentProcessor, err := content.New(config, loadedPartials)
	assert.NoError(t, err)

	// Act
	err = contentProcessor.ProcessMarkdown("content/post.md")

	// Assert
	assert.NoError(t, err)
	assert.FileExists(t, "dist/post.html")

	// Verify the content of the created file
	content, err := os.ReadFile("dist/post.html")
	assert.NoError(t, err)
	assert.Contains(t, string(content), "<html><body><h1>My Post</h1>\n</body></html>")

	// Clean up
	os.RemoveAll(distDir)
	os.RemoveAll(contentDir)
	os.RemoveAll(partialsDir)
}

func BenchmarkProcessHTML(b *testing.B) {
	// Arrange
	distDir := "dist"
	contentDir := "content"
	partialsDir := "partials"
	os.RemoveAll(distDir)
	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(partialsDir, 0755)

	// Create a dummy content file
	os.WriteFile("content/index.html", []byte("<h1>Hello</h1>"), 0644)
	// Create a dummy layout file
	os.WriteFile("content/_layout.html", []byte("<html><body>{{.Content}}</body></html>"), 0644)

	config := map[string]interface{}{}
	loadedPartials, _ := partials.LoadPartials()
	contentProcessor, err := content.New(config, loadedPartials)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	// Act
	for i := 0; i < b.N; i++ {
		err := contentProcessor.ProcessHTML("content/index.html")
		if err != nil {
			b.Fatal(err)
		}
	}

	// Clean up
	os.RemoveAll(distDir)
	os.RemoveAll(contentDir)
	os.RemoveAll(partialsDir)
}

func BenchmarkProcessMarkdown(b *testing.B) {
	// Arrange
	distDir := "dist"
	contentDir := "content"
	partialsDir := "partials"
	os.RemoveAll(distDir)
	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(partialsDir, 0755)

	// Create a dummy content file
	os.WriteFile("content/post.md", []byte("# My Post"), 0644)
	// Create a dummy layout file
	os.WriteFile("content/_layout.html", []byte("<html><body>{{.Content}}</body></html>"), 0644)

	config := map[string]interface{}{}
	loadedPartials, _ := partials.LoadPartials()
	contentProcessor, err := content.New(config, loadedPartials)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	// Act
	for i := 0; i < b.N; i++ {
		err := contentProcessor.ProcessMarkdown("content/post.md")
		if err != nil {
			b.Fatal(err)
		}
	}

	// Clean up
	os.RemoveAll(distDir)
	os.RemoveAll(contentDir)
	os.RemoveAll(partialsDir)
}
