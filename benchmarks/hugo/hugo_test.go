package hugo

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func generateHugoSite(b *testing.B, siteDir string, numPages int, avgContentSize int) {
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

func BenchmarkHugo(b *testing.B) {
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
			generateHugoSite(b, siteDir, bm.numPages, bm.avgContentSize)

			// Create a minimal config.toml file
			configContent := `baseURL = "http://example.org/"
languageCode = "en-us"
title = "My New Hugo Site"
`
			if err := os.WriteFile(filepath.Join(siteDir, "config.toml"), []byte(configContent), 0644); err != nil {
				b.Fatalf("Failed to write config.toml: %v", err)
			}

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

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				cmd := exec.Command("hugo")
				if err := cmd.Run(); err != nil {
					b.Fatalf("Failed to run hugo: %v", err)
				}
			}
		})
	}
}
