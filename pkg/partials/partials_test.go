package partials_test

import (
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

func BenchmarkLoadPartials(b *testing.B) {
	// Arrange
	partialsDir := "partials"
	os.MkdirAll(partialsDir, 0755)
	os.WriteFile("partials/base.html", []byte("{{.Title}}"), 0644)
	os.WriteFile("partials/post.html", []byte("{{.Content}}"), 0644)

	b.ReportAllocs()

	// Act
	for b.Loop() {
		_, err := partials.LoadPartials()
		if err != nil {
			b.Fatal(err)
		}
	}

	// Clean up
	os.RemoveAll(partialsDir)
}
