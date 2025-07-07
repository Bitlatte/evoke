package detailed

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"math/rand"

	"github.com/Bitlatte/evoke/pkg/build"
	"github.com/Bitlatte/evoke/pkg/content"
)

func generateSite(b *testing.B, siteDir string, numPages int, avgContentSize int) {
	contentDir := filepath.Join(siteDir, "content")
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		b.Fatalf("Failed to create content directory: %v", err)
	}

	for i := 1; i <= numPages; i++ {
		filePath := filepath.Join(contentDir, fmt.Sprintf("page-%d.md", i))
		file, err := os.Create(filePath)
		if err != nil {
			b.Fatalf("Failed to create page: %v", err)
		}
		defer file.Close()

		fmt.Fprintf(file, "---\ntitle: Page %d\n---\n", i)
		contentSize := avgContentSize/2 + rand.Intn(avgContentSize)
		for j := 0; j < contentSize; j++ {
			fmt.Fprintln(file, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
		}
	}
}

func BenchmarkBuildProcess(b *testing.B) {
	benchmarks := []struct {
		name           string
		numPages       int
		avgContentSize int
	}{
		{"Tiny", 1, 10},
		{"Small", 100, 50},
		{"Medium", 1000, 100},
		{"Large", 1000, 500},
		{"Huge", 10000, 500},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			siteDir := b.TempDir()
			generateSite(b, siteDir, bm.numPages, bm.avgContentSize)

			// Store the original working directory
			originalDir, err := os.Getwd()
			if err != nil {
				b.Fatalf("Failed to get current directory: %v", err)
			}
			defer os.Chdir(originalDir)

			// Change to the temporary site directory
			if err := os.Chdir(siteDir); err != nil {
				b.Fatalf("Failed to change directory: %v", err)
			}

			config, err := build.LoadConfiguration()
			if err != nil {
				b.Fatalf("Failed to load configuration: %v", err)
			}

			partials, err := build.LoadPartials()
			if err != nil {
				b.Fatalf("Failed to load partials: %v", err)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				contentProcessor, err := content.New(config, partials)
				if err != nil {
					b.Fatalf("Failed to create content processor: %v", err)
				}
				build.ProcessContentWithProcessor(contentProcessor)
			}
		})
	}
}
