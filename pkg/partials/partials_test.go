package partials_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/stretchr/testify/assert"
)

func TestLoadPartials_LoadsAllPartials(t *testing.T) {
	// Arrange
	partialsDir := "partials"
	os.MkdirAll(partialsDir, 0755)
	os.WriteFile("partials/base.html", []byte("{{.Title}}"), 0644)
	os.WriteFile("partials/post.html", []byte("{{.Content}}"), 0644)

	// Act
	loadedPartials, err := partials.LoadPartials()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, loadedPartials.Lookup("base.html"))
	assert.NotNil(t, loadedPartials.Lookup("post.html"))

	// Clean up
	os.RemoveAll(partialsDir)
}

func generateBenchmarkPartials(b *testing.B, numPartials int) {
	partialsDir := "partials"
	os.MkdirAll(partialsDir, 0755)

	for i := 0; i < numPartials; i++ {
		filename := fmt.Sprintf("partials/partial-%d.html", i)
		content := fmt.Sprintf(`
			<div class="partial-%d">
				<h2>Partial Number %d</h2>
				<p>This is some content for the partial.</p>
				{{ block "content" . }}{{ end }}
			</div>
		`, i, i)
		os.WriteFile(filename, []byte(content), 0644)
	}
}

func BenchmarkLoadPartials(b *testing.B) {
	// Arrange
	generateBenchmarkPartials(b, 50) // Generate 50 partial files
	defer os.RemoveAll("partials")

	b.ResetTimer()
	b.ReportAllocs()

	// Act
	for i := 0; i < b.N; i++ {
		_, err := partials.LoadPartials()
		if err != nil {
			b.Fatal(err)
		}
	}
}
