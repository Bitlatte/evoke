package build

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/Bitlatte/evoke/pkg/config"
	"github.com/Bitlatte/evoke/pkg/content"
	"github.com/Bitlatte/evoke/pkg/extensions"
	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/Bitlatte/evoke/pkg/util"
)

func Build() error {
	// Load extensions
	loadedExtensions, err := extensions.LoadBuildExtensions()
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

	// Load partials
	var t *partials.Partials
	if _, err := os.Stat("partials"); !os.IsNotExist(err) {
		t, err = partials.LoadPartials()
		if err != nil {
			return fmt.Errorf("error loading partials: %w", err)
		}
	}

	if _, statErr := os.Stat("content"); !os.IsNotExist(statErr) {
		if statErr != nil {
			err = fmt.Errorf("error checking content directory: %w", statErr)
		} else {
			var wg sync.WaitGroup
			jobs := make(chan string)
			errs := make(chan error, runtime.NumCPU())

			for i := 0; i < runtime.NumCPU(); i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for path := range jobs {
						ext := filepath.Ext(path)
						var err error
						switch ext {
						case ".html":
							err = content.ProcessHTML(path, loadedConfig, t)
						case ".md":
							err = content.ProcessMarkdown(path, loadedConfig, t)
						}
						if err != nil {
							errs <- err
						}
					}
				}()
			}

			err = filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() && info.Name()[0] != '_' {
					jobs <- path
				}

				return nil
			})

			close(jobs)
			wg.Wait()
			close(errs)

			if err != nil {
				return err
			}

			for e := range errs {
				// For now, we'll just return the first error we see.
				// A more robust solution might collect all errors.
				return e
			}
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
