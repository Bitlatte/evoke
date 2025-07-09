package pipelines_test

import (
	"strings"
	"testing"

	"github.com/Bitlatte/evoke/pkg/pipelines"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestMarkdownPipeline(t *testing.T) {
	// Arrange
	gm := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	pipeline := pipelines.NewMarkdownPipeline(gm)
	asset := &pipelines.Asset{
		Path:    "content/post.md",
		Content: strings.NewReader("<h1>My Post</h1>"),
	}

	// Act
	processedAsset, err := pipeline.Process(asset)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "content/post.html", processedAsset.Path)
	assert.NotNil(t, processedAsset.Content)
}

func TestHTMLPipeline(t *testing.T) {
	// Arrange
	pipeline := pipelines.NewHTMLPipeline()
	asset := &pipelines.Asset{
		Path:    "content/index.html",
		Content: strings.NewReader("<h1>Hello</h1>"),
	}

	// Act
	processedAsset, err := pipeline.Process(asset)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "content/index.html", processedAsset.Path)
	assert.NotNil(t, processedAsset.Content)
}

func TestCopyPipeline(t *testing.T) {
	// Arrange
	pipeline := pipelines.NewCopyPipeline()
	asset := &pipelines.Asset{
		Path:    "content/image.jpg",
		Content: strings.NewReader(""),
	}

	// Act
	processedAsset, err := pipeline.Process(asset)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "content/image.jpg", processedAsset.Path)
}

const loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

func generateMarkdownContent(paragraphs int) string {
	var b strings.Builder
	b.WriteString("# Test Post\n\n")
	for i := 0; i < paragraphs; i++ {
		b.WriteString(loremIpsum)
		b.WriteString("\n\n")
	}
	return b.String()
}

func generateHTMLContent(paragraphs int) string {
	var b strings.Builder
	b.WriteString("<html><head><title>Test</title></head><body>\n")
	b.WriteString("<h1>Test Post</h1>\n")
	for i := 0; i < paragraphs; i++ {
		b.WriteString("<p>")
		b.WriteString(loremIpsum)
		b.WriteString("</p>\n")
	}
	b.WriteString("</body></html>")
	return b.String()
}

func BenchmarkMarkdownPipeline(b *testing.B) {
	// Arrange
	gm := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	pipeline := pipelines.NewMarkdownPipeline(gm)
	content := generateMarkdownContent(100) // 100 paragraphs of text

	b.ResetTimer()
	b.ReportAllocs()

	// Act
	for b.Loop() {
		// This asset needs to be created inside the loop because the pipeline
		// consumes the Content reader.
		asset := &pipelines.Asset{
			Path:    "content/post.md",
			Content: strings.NewReader(content),
		}
		_, err := pipeline.Process(asset)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHTMLPipeline(b *testing.B) {
	// Arrange
	pipeline := pipelines.NewHTMLPipeline()
	content := generateHTMLContent(100) // 100 paragraphs of text

	b.ResetTimer()
	b.ReportAllocs()

	// Act
	for b.Loop() {
		// This asset needs to be created inside the loop because the pipeline
		// consumes the Content reader.
		asset := &pipelines.Asset{
			Path:    "content/index.html",
			Content: strings.NewReader(content),
		}
		_, err := pipeline.Process(asset)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCopyPipeline(b *testing.B) {
	// Arrange
	pipeline := pipelines.NewCopyPipeline()
	content := []byte(strings.Repeat("a", 1024*1024)) // 1MB of data

	b.ResetTimer()
	b.ReportAllocs()

	// Act
	for b.Loop() {
		asset := &pipelines.Asset{
			Path:    "content/image.jpg",
			Content: strings.NewReader(string(content)),
		}
		_, err := pipeline.Process(asset)
		if err != nil {
			b.Fatal(err)
		}
	}
}
