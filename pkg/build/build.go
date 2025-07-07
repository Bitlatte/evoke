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

// LoadExtensions loads the build extensions.
func LoadExtensions() ([]extensions.Extension, error) {
	return extensions.LoadBuildExtensions()
}

// RunBeforeBuildHooks runs the BeforeBuild hooks for the given extensions.
func RunBeforeBuildHooks(loadedExtensions []extensions.Extension) error {
	for _, ext := range loadedExtensions {
		if err := ext.BeforeBuild(); err != nil {
			return fmt.Errorf("error running BeforeBuild hook: %w", err)
		}
	}
	return nil
}

// CreateOutputDirectory creates the output directory.
func CreateOutputDirectory() error {
	return os.MkdirAll("dist", 0755)
}

// CopyPublicDirectory copies the public directory to the output directory.
func CopyPublicDirectory() error {
	if _, err := os.Stat("public"); !os.IsNotExist(err) {
		if err := util.CopyDirectory("public", "dist"); err != nil {
			return fmt.Errorf("error copying public directory: %w", err)
		}
	}
	return nil
}

// LoadConfiguration loads the configuration.
func LoadConfiguration() (map[string]interface{}, error) {
	return config.LoadConfig()
}

// LoadPartials loads the partials.
func LoadPartials() (*partials.Partials, error) {
	if _, err := os.Stat("partials"); !os.IsNotExist(err) {
		return partials.LoadPartials()
	}
	return nil, nil
}

// ProcessContent processes the content.
func ProcessContent(loadedConfig map[string]interface{}, t *partials.Partials) error {
	contentProcessor, err := content.New(loadedConfig, t)
	if err != nil {
		return fmt.Errorf("error creating content processor: %w", err)
	}
	return ProcessContentWithProcessor(contentProcessor)
}

// ProcessContentWithProcessor processes the content with a given processor.
func ProcessContentWithProcessor(contentProcessor *content.Content) error {
	if _, statErr := os.Stat("content"); !os.IsNotExist(statErr) {
		if statErr != nil {
			return fmt.Errorf("error checking content directory: %w", statErr)
		}
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
						err = contentProcessor.ProcessHTML(path)
					case ".md":
						err = contentProcessor.ProcessMarkdown(path)
					}
					if err != nil {
						errs <- err
					}
				}
			}()
		}

		err := filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
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
	return nil
}

// RunAfterBuildHooks runs the AfterBuild hooks for the given extensions.
func RunAfterBuildHooks(loadedExtensions []extensions.Extension) error {
	for _, ext := range loadedExtensions {
		if err := ext.AfterBuild(); err != nil {
			return fmt.Errorf("error running AfterBuild hook: %w", err)
		}
	}
	return nil
}

func Build() error {
	// Load extensions
	loadedExtensions, err := LoadExtensions()
	if err != nil {
		return fmt.Errorf("error loading extensions: %w", err)
	}

	// Run BeforeBuild hooks
	if err := RunBeforeBuildHooks(loadedExtensions); err != nil {
		return err
	}

	// Create the output directory
	if err := CreateOutputDirectory(); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Copy the public directory
	if err := CopyPublicDirectory(); err != nil {
		return err
	}

	// Process content
	loadedConfig, err := LoadConfiguration()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Load partials
	t, err := LoadPartials()
	if err != nil {
		return fmt.Errorf("error loading partials: %w", err)
	}

	if err := ProcessContent(loadedConfig, t); err != nil {
		return err
	}

	// Run AfterBuild hooks
	if err := RunAfterBuildHooks(loadedExtensions); err != nil {
		return err
	}

	return nil
}
