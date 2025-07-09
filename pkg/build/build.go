package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/Bitlatte/evoke/pkg/config"
	"github.com/Bitlatte/evoke/pkg/content"
	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/Bitlatte/evoke/pkg/plugins"
	"github.com/Bitlatte/evoke/pkg/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// LoadPlugins loads the build plugins.
func LoadPlugins() ([]plugins.Plugin, error) {
	if _, err := os.Stat("plugins"); os.IsNotExist(err) {
		return nil, nil
	}
	return plugins.LoadPlugins()
}

// RunOnPreBuildHooks runs the OnPreBuild hooks for the given plugins.
func RunOnPreBuildHooks(loadedPlugins []plugins.Plugin) error {
	for _, p := range loadedPlugins {
		if err := p.OnPreBuild(); err != nil {
			return fmt.Errorf("error running OnPreBuild hook: %w", err)
		}
	}
	return nil
}

// RunOnConfigLoadedHooks runs the OnConfigLoaded hooks for the given plugins.
func RunOnConfigLoadedHooks(loadedPlugins []plugins.Plugin, config []byte) ([]byte, error) {
	for _, p := range loadedPlugins {
		var err error
		config, err = p.OnConfigLoaded(config)
		if err != nil {
			return nil, fmt.Errorf("error running OnConfigLoaded hook: %w", err)
		}
	}
	return config, nil
}

// RunOnPublicAssetsCopiedHooks runs the OnPublicAssetsCopied hooks for the given plugins.
func RunOnPublicAssetsCopiedHooks(loadedPlugins []plugins.Plugin) error {
	for _, p := range loadedPlugins {
		if err := p.OnPublicAssetsCopied(); err != nil {
			return fmt.Errorf("error running OnPublicAssetsCopied hook: %w", err)
		}
	}
	return nil
}

// CreateOutputDirectory creates the output directory.
func CreateOutputDirectory(outputDir string) error {
	return os.MkdirAll(outputDir, 0755)
}

// CopyPublicDirectory copies the public directory to the output directory.
func CopyPublicDirectory(outputDir string) error {
	if _, err := os.Stat("public"); !os.IsNotExist(err) {
		if err := util.CopyDirectory("public", outputDir); err != nil {
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
	return &partials.Partials{}, nil
}

// ProcessContent processes the content.
func ProcessContent(outputDir string, loadedConfig map[string]interface{}, t *partials.Partials, loadedPlugins []plugins.Plugin) error {
	gm := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	contentProcessor, err := content.New(outputDir, loadedConfig, t, gm, loadedPlugins)
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

// RunOnPostBuildHooks runs the OnPostBuild hooks for the given plugins.
func RunOnPostBuildHooks(loadedPlugins []plugins.Plugin) error {
	for _, p := range loadedPlugins {
		if err := p.OnPostBuild(); err != nil {
			return fmt.Errorf("error running OnPostBuild hook: %w", err)
		}
	}
	return nil
}

func Build(outputDir string) error {
	// Load plugins
	loadedPlugins, err := LoadPlugins()
	if err != nil {
		return fmt.Errorf("error loading plugins: %w", err)
	}

	// Run OnPreBuild hooks
	if err := RunOnPreBuildHooks(loadedPlugins); err != nil {
		return err
	}

	// Create the output directory
	if err := CreateOutputDirectory(outputDir); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Copy the public directory
	if err := CopyPublicDirectory(outputDir); err != nil {
		return err
	}

	// Run OnPublicAssetsCopied hooks
	if err := RunOnPublicAssetsCopiedHooks(loadedPlugins); err != nil {
		return err
	}

	// Process content
	loadedConfig, err := LoadConfiguration()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Run OnConfigLoaded hooks
	configBytes, err := json.Marshal(loadedConfig)
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}
	configBytes, err = RunOnConfigLoadedHooks(loadedPlugins, configBytes)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(configBytes, &loadedConfig); err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Load partials
	t, err := LoadPartials()
	if err != nil {
		return fmt.Errorf("error loading partials: %w", err)
	}

	if err := ProcessContent(outputDir, loadedConfig, t, loadedPlugins); err != nil {
		return err
	}

	// Run OnPostBuild hooks
	if err := RunOnPostBuildHooks(loadedPlugins); err != nil {
		return err
	}

	return nil
}
