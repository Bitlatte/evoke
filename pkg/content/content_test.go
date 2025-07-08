package content_test

import (
	"os"
	"testing"

	"github.com/Bitlatte/evoke/pkg/content"
	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
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
	gm := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	contentProcessor, err := content.New(config, loadedPartials, gm)
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
	gm := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	contentProcessor, err := content.New(config, loadedPartials, gm)
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
	benchmarks := []struct {
		name    string
		content string
		layout  string
	}{
		{"Small", "<h1>Hello</h1>", "<html><body>{{.Content}}</body></html>"},
		{"Medium", "<h1>Hello</h1><p>" + "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " + "</p>", "<html><body><header></header>{{.Content}}<footer></footer></body></html>"},
		{"Large", "<h1>Hello</h1><p>" + "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " + "</p>" + "<div>" + "Donec a diam lectus. Sed sit amet ipsum mauris. " + "</div>", "<html><body><header></header><main>{{.Content}}</main><footer></footer></body></html>"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			distDir := "dist"
			contentDir := "content"
			partialsDir := "partials"
			os.RemoveAll(distDir)
			os.MkdirAll(contentDir, 0755)
			os.MkdirAll(partialsDir, 0755)
			defer os.RemoveAll(distDir)
			defer os.RemoveAll(contentDir)
			defer os.RemoveAll(partialsDir)

			os.WriteFile("content/index.html", []byte(bm.content), 0644)
			os.WriteFile("content/_layout.html", []byte(bm.layout), 0644)

			config := map[string]interface{}{}
			loadedPartials, _ := partials.LoadPartials()
			gm := goldmark.New(
				goldmark.WithRendererOptions(
					html.WithUnsafe(),
				),
			)
			contentProcessor, err := content.New(config, loadedPartials, gm)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for b.Loop() {
				err := contentProcessor.ProcessHTML("content/index.html")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkProcessMarkdown(b *testing.B) {
	benchmarks := []struct {
		name    string
		content string
		layout  string
	}{
		{"Small", "# My Post", "<html><body>{{.Content}}</body></html>"},
		{"Medium", "# My Post\n\n" + "Lorem ipsum dolor sit amet, consectetur adipiscing elit. ", "<html><body><header></header>{{.Content}}<footer></footer></body></html>"},
		{"Large", "# My Post\n\n" + "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " + "\n\n" + "Donec a diam lectus. Sed sit amet ipsum mauris. ", "<html><body><header></header><main>{{.Content}}</main><footer></footer></body></html>"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			distDir := "dist"
			contentDir := "content"
			partialsDir := "partials"
			os.RemoveAll(distDir)
			os.MkdirAll(contentDir, 0755)
			os.MkdirAll(partialsDir, 0755)
			defer os.RemoveAll(distDir)
			defer os.RemoveAll(contentDir)
			defer os.RemoveAll(partialsDir)

			os.WriteFile("content/post.md", []byte(bm.content), 0644)
			os.WriteFile("content/_layout.html", []byte(bm.layout), 0644)

			config := map[string]interface{}{}
			loadedPartials, _ := partials.LoadPartials()
			gm := goldmark.New(
				goldmark.WithRendererOptions(
					html.WithUnsafe(),
				),
			)
			contentProcessor, err := content.New(config, loadedPartials, gm)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for b.Loop() {
				err := contentProcessor.ProcessMarkdown("content/post.md")
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
