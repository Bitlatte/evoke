package build

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/Bitlatte/evoke/pkg/config"
	"github.com/Bitlatte/evoke/pkg/content"
	"github.com/Bitlatte/evoke/pkg/extensions"
	"github.com/Bitlatte/evoke/pkg/templates"
	"github.com/Bitlatte/evoke/pkg/util"
)

func Build() error {
	// Load extensions
	loadedExtensions, err := extensions.LoadExtensions()
	if err != nil {
		return fmt.Errorf("error loading extensions: %w", err)
	}

	// Run BeforeBuild hooks
	for _, ext := range loadedExtensions {
		if err := ext.BeforeBuild(); err != nil {
			return fmt.Errorf("error running BeforeBuild hook: %w", err)
		}
	}

	fmt.Println("Building...")
	// Create the output directory
	if err := os.MkdirAll("dist", 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Copy the public directory
	if _, err := os.Stat("public"); !os.IsNotExist(err) {
		if err := util.CopyDirectory("public", "dist"); err != nil {
			return fmt.Errorf("error copying public directory: %w", err)
		}
	}

	// Process content
	loadedConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	var loadedTemplates *template.Template
	if _, err := os.Stat("templates"); !os.IsNotExist(err) {
		loadedTemplates, err = templates.LoadTemplates()
		if err != nil {
			return fmt.Errorf("error loading templates: %w", err)
		}
	}

	if _, statErr := os.Stat("content"); !os.IsNotExist(statErr) {
		if statErr != nil {
			err = fmt.Errorf("error checking content directory: %w", statErr)
		} else {
			err = filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					ext := filepath.Ext(path)
					switch ext {
					case ".html":
						return content.ProcessHTML(path, loadedConfig, loadedTemplates)
					case ".md":
						return content.ProcessMarkdown(path, loadedConfig, loadedTemplates)
					}
				}

				return nil
			})
		}
	}

	if err != nil {
		return err
	}

	// Run AfterBuild hooks
	for _, ext := range loadedExtensions {
		if err := ext.AfterBuild(); err != nil {
			return fmt.Errorf("error running AfterBuild hook: %w", err)
		}
	}

	return nil
}
