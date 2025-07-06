package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "evoke",
		Usage: "a powerful little static site generator",
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "builds your content into static HTML",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return build()
				},
			},
			{
				Name:  "extension",
				Usage: "manage extensions",
				Commands: []*cli.Command{
					{
						Name:  "get",
						Usage: "get a new extension from a url",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							fmt.Println("Getting extension...")
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func build() error {
	fmt.Println("Building...")
	// Create the output directory
	if err := os.MkdirAll("dist", 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Copy the public directory
	if err := copyDirectory("public", "dist"); err != nil {
		return fmt.Errorf("error copying public directory: %w", err)
	}

	// Process content
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	templates, err := loadTemplates()
	if err != nil {
		return fmt.Errorf("error loading templates: %w", err)
	}

	err = filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".html" {
				return processHTML(path, config, templates)
			} else if ext == ".md" {
				return processMarkdown(path, config, templates)
			}
		}

		return nil
	})

	return err
}

func copyDirectory(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a parallel structure in the destination
		destPath := filepath.Join(dest, path[len(src):])

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy the file
		return copyFile(path, destPath)
	})
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
