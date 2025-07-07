package suite

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var scenarios = []struct {
	name  string
	pages int
}{
	{"small", 100},
	{"medium", 1000},
	{"large", 10000},
}

func generateSite(b *testing.B, siteDir string, numPages int) {
	contentDir := filepath.Join(siteDir, "content")
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		b.Fatalf("Failed to create content directory: %v", err)
	}

	for i := 1; i <= numPages; i++ {
		content := fmt.Sprintf("---\ntitle: Page %d\n---\nThis is page %d.\nLorem ipsum dolor sit amet, consectetur adipiscing elit.", i, i)
		filePath := filepath.Join(contentDir, fmt.Sprintf("page-%d.md", i))
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			b.Fatalf("Failed to write page: %v", err)
		}
	}
}

func BenchmarkEvoke(b *testing.B) {
	evokePath, err := filepath.Abs("../../evoke")
	if err != nil {
		b.Fatalf("Failed to get absolute path for evoke executable: %v", err)
	}

	for _, s := range scenarios {
		b.Run(s.name, func(b *testing.B) {
			siteDir := b.TempDir()
			generateSite(b, siteDir, s.pages)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cmd := exec.Command(evokePath, "build")
				cmd.Dir = siteDir
				if err := cmd.Run(); err != nil {
					b.Fatalf("Evoke build failed: %v", err)
				}
			}
		})
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func BenchmarkHugo(b *testing.B) {
	if !commandExists("hugo") {
		b.Skip("hugo command not found, skipping benchmarks")
	}

	for _, s := range scenarios {
		b.Run(s.name, func(b *testing.B) {
			siteDir := b.TempDir()
			// Hugo requires a theme, even a blank one.
			if err := os.MkdirAll(filepath.Join(siteDir, "themes", "empty"), 0755); err != nil {
				b.Fatalf("Failed to create theme directory: %v", err)
			}
			// Create a minimal config.toml
			configContent := `baseURL = "http://example.org/"
languageCode = "en-us"
title = "My New Hugo Site"
theme = "empty"
`
			if err := os.WriteFile(filepath.Join(siteDir, "hugo.toml"), []byte(configContent), 0644); err != nil {
				b.Fatalf("Failed to write hugo.toml: %v", err)
			}

			generateSite(b, siteDir, s.pages)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cmd := exec.Command("hugo")
				cmd.Dir = siteDir
				if err := cmd.Run(); err != nil {
					b.Fatalf("Hugo build failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkJekyll(b *testing.B) {
	if !commandExists("jekyll") {
		b.Skip("jekyll command not found, skipping benchmarks")
	}

	for _, s := range scenarios {
		b.Run(s.name, func(b *testing.B) {
			siteDir := b.TempDir()
			// Create a minimal _config.yml
			configContent := `
theme: minima
`
			if err := os.WriteFile(filepath.Join(siteDir, "_config.yml"), []byte(configContent), 0644); err != nil {
				b.Fatalf("Failed to write _config.yml: %v", err)
			}

			// Jekyll uses _posts for markdown files
			postsDir := filepath.Join(siteDir, "_posts")
			generateSite(b, siteDir, s.pages)
			// Move the generated content to the _posts directory
			contentDir := filepath.Join(siteDir, "content")
			files, err := os.ReadDir(contentDir)
			if err != nil {
				b.Fatalf("Failed to read content directory: %v", err)
			}
			if err := os.MkdirAll(postsDir, 0755); err != nil {
				b.Fatalf("Failed to create _posts directory: %v", err)
			}
			for _, file := range files {
				oldPath := filepath.Join(contentDir, file.Name())
				newPath := filepath.Join(postsDir, file.Name())
				if err := os.Rename(oldPath, newPath); err != nil {
					b.Fatalf("Failed to move file to _posts: %v", err)
				}
			}
			os.Remove(contentDir)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cmd := exec.Command("jekyll", "build")
				cmd.Dir = siteDir
				if err := cmd.Run(); err != nil {
					b.Fatalf("Jekyll build failed: %v", err)
				}
			}
		})
	}
}
