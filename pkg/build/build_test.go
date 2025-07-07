package build_test

import (
	"os"
	"testing"

	"github.com/Bitlatte/evoke/pkg/build"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "evoke-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change to the temporary directory
	originalWd, err := os.Getwd()
	assert.NoError(t, err)
	err = os.Chdir(tmpDir)
	assert.NoError(t, err)
	defer os.Chdir(originalWd)

	// Create the necessary directories
	os.Mkdir("content", 0755)
	os.Mkdir("partials", 0755)
	os.Mkdir("public", 0755)
	os.Mkdir("extensions", 0755)

	// Create dummy files
	os.WriteFile("evoke.yaml", []byte("title: My Site"), 0644)
	os.WriteFile("content/index.md", []byte("# Hello World"), 0644)
	os.WriteFile("content/_layout.html", []byte("<html><body>{{.Content}}</body></html>"), 0644)
	os.WriteFile("partials/header.html", []byte("<header>My Header</header>"), 0644)
	os.WriteFile("public/style.css", []byte("body { color: red; }"), 0644)

	// Run the build
	err = build.Build()
	assert.NoError(t, err)

	// Assert the results
	assert.FileExists(t, "dist/index.html")
	assert.FileExists(t, "dist/style.css")

	// Verify the content of the created file
	content, err := os.ReadFile("dist/index.html")
	assert.NoError(t, err)
	assert.Contains(t, string(content), "<html><body><h1>Hello World</h1>\n</body></html>")

	css, err := os.ReadFile("dist/style.css")
	assert.NoError(t, err)
	assert.Equal(t, "body { color: red; }", string(css))
}

func BenchmarkBuild(b *testing.B) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "evoke-benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to the temporary directory
	originalWd, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}
	err = os.Chdir(tmpDir)
	if err != nil {
		b.Fatal(err)
	}
	defer os.Chdir(originalWd)

	// Create the necessary directories
	os.Mkdir("content", 0755)
	os.Mkdir("partials", 0755)
	os.Mkdir("public", 0755)
	os.Mkdir("extensions", 0755)

	// Create dummy files
	os.WriteFile("evoke.yaml", []byte("title: My Site"), 0644)
	os.WriteFile("content/index.md", []byte("# Hello World"), 0644)
	os.WriteFile("content/_layout.html", []byte("<html><body>{{.Content}}</body></html>"), 0644)
	os.WriteFile("partials/header.html", []byte("<header>My Header</header>"), 0644)
	os.WriteFile("public/style.css", []byte("body { color: red; }"), 0644)

	b.ResetTimer()
	b.ReportAllocs()

	// Run the build
	for i := 0; i < b.N; i++ {
		err = build.Build()
		if err != nil {
			b.Fatal(err)
		}
	}
}
