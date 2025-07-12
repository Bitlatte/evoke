// Package build provides the functionality to build the site.
package build

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"html/template"

	"github.com/Bitlatte/evoke/pkg/cache"
	"github.com/Bitlatte/evoke/pkg/config"
	"github.com/Bitlatte/evoke/pkg/content"
	"github.com/Bitlatte/evoke/pkg/dag"
	"github.com/Bitlatte/evoke/pkg/defaults"
	"github.com/Bitlatte/evoke/pkg/diff"
	"github.com/Bitlatte/evoke/pkg/hash"
	"github.com/Bitlatte/evoke/pkg/logger"
	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/Bitlatte/evoke/pkg/pipelines"
	"github.com/Bitlatte/evoke/pkg/plugins"
	"github.com/Bitlatte/evoke/pkg/util"
	"github.com/yuin/goldmark"

	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// LoadPlugins loads the build plugins.
func LoadPlugins() ([]plugins.Plugin, error) {
	logger.Logger.Debug("Loading plugins...")
	if _, err := os.Stat("plugins"); os.IsNotExist(err) {
		logger.Logger.Debug("No plugins directory found, skipping plugin loading.")
		return nil, nil
	}
	p, err := plugins.LoadPlugins()
	if err != nil {
		return nil, err
	}
	logger.Logger.Debug("Plugins loaded.", "count", len(p))
	return p, nil
}

// RunOnPreBuildHooks runs the OnPreBuild hooks for the given plugins.
func RunOnPreBuildHooks(loadedPlugins []plugins.Plugin) error {
	logger.Logger.Debug("Running OnPreBuild hooks...")
	for _, p := range loadedPlugins {
		logger.Logger.Debug("Running OnPreBuild hook", "plugin", p.Name())
		if err := p.OnPreBuild(); err != nil {
			return fmt.Errorf("error running OnPreBuild hook: %w", err)
		}
	}
	return nil
}

// RunOnConfigLoadedHooks runs the OnConfigLoaded hooks for the given plugins.
func RunOnConfigLoadedHooks(loadedPlugins []plugins.Plugin, config []byte) ([]byte, error) {
	logger.Logger.Debug("Running OnConfigLoaded hooks...")
	for _, p := range loadedPlugins {
		logger.Logger.Debug("Running OnConfigLoaded hook", "plugin", p.Name())
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
	logger.Logger.Debug("Running OnPublicAssetsCopied hooks...")
	for _, p := range loadedPlugins {
		logger.Logger.Debug("Running OnPublicAssetsCopied hook", "plugin", p.Name())
		if err := p.OnPublicAssetsCopied(); err != nil {
			return fmt.Errorf("error running OnPublicAssetsCopied hook: %w", err)
		}
	}
	return nil
}

// CreateOutputDirectory creates the output directory.
func CreateOutputDirectory(outputDir string) error {
	logger.Logger.Debug("Creating output directory...", "path", outputDir)
	return os.MkdirAll(outputDir, 0755)
}

// CopyPublicDirectory copies the public directory to the output directory.
func CopyPublicDirectory(outputDir string) error {
	logger.Logger.Debug("Copying public directory...")
	if _, err := os.Stat("public"); !os.IsNotExist(err) {
		if err := util.CopyDirectory("public", outputDir); err != nil {
			return fmt.Errorf("error copying public directory: %w", err)
		}
	}
	logger.Logger.Debug("Public directory copied.")
	return nil
}

// LoadConfiguration loads the configuration.
func LoadConfiguration() (map[string]interface{}, error) {
	logger.Logger.Debug("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	logger.Logger.Debug("Configuration loaded.")
	return cfg, nil
}

// LoadPartials loads the partials.
func LoadPartials() (*partials.Partials, error) {
	logger.Logger.Debug("Loading partials...")
	if _, err := os.Stat("partials"); !os.IsNotExist(err) {
		p, err := partials.LoadPartials()
		if err != nil {
			return nil, err
		}
		logger.Logger.Debug("Partials loaded.")
		return p, nil
	}
	logger.Logger.Debug("No partials directory found, skipping partial loading.")
	return &partials.Partials{}, nil
}

// ProcessContent processes the content.
func ProcessContent(outputDir string, loadedConfig map[string]interface{}, t *partials.Partials, loadedPlugins []plugins.Plugin) error {
	logger.Logger.Debug("Processing content...")
	gm := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var p []pipelines.Pipeline
	p = append(p, pipelines.NewMarkdownPipeline(gm))
	p = append(p, pipelines.NewHTMLPipeline())
	p = append(p, pipelines.NewCopyPipeline())

	// Register plugin pipelines
	for _, plugin := range loadedPlugins {
		pluginPipelines, err := plugin.RegisterPipelines()
		if err != nil {
			return fmt.Errorf("error registering pipelines: %w", err)
		}
		for _, pipeline := range pluginPipelines {
			p = append(p, pipelines.NewGRPCPipeline(plugin, pipeline.Name))
		}
	}

	contentProcessor, err := content.New(outputDir, loadedConfig, t, gm, loadedPlugins, p)
	if err != nil {
		return fmt.Errorf("error creating content processor: %w", err)
	}

	// Create a new cache
	c, err := cache.New(filepath.Join(outputDir, ".cache"))
	if err != nil {
		return fmt.Errorf("error creating cache: %w", err)
	}

	// Build the dependency graph
	d, err := dag.BuildGraph("content", "partials")
	if err != nil {
		return fmt.Errorf("error building dependency graph: %w", err)
	}

	// Get the files to rebuild
	toRebuild, err := getFilesToRebuild(c, d)
	if err != nil {
		return fmt.Errorf("error getting files to rebuild: %w", err)
	}

	if err := ProcessContentWithProcessor(contentProcessor, loadedConfig, toRebuild); err != nil {
		return err
	}

	// Save the cache
	if err := c.Save(); err != nil {
		return fmt.Errorf("error saving cache: %w", err)
	}

	return nil
}

// ProcessContentWithProcessor processes the content with a given processor.
func ProcessContentWithProcessor(contentProcessor *content.Content, loadedConfig map[string]interface{}, toRebuild map[string]bool) error {
	if _, statErr := os.Stat("content"); os.IsNotExist(statErr) {
		return nil // No content directory, nothing to do.
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context is always cancelled

	var wg sync.WaitGroup
	jobs := make(chan pipelines.Asset)
	errs := make(chan error, 1)

	// Function to handle error and cancel context
	handleError := func(err error) {
		select {
		case errs <- err:
		default:
		}
		cancel()
	}

	// Start worker goroutines
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case asset, ok := <-jobs:
					if !ok {
						return
					}

					var processedAsset *pipelines.Asset
					var err error
					var pipeline pipelines.Pipeline
					ext := filepath.Ext(asset.Path)
					switch ext {
					case ".md":
						pipeline = contentProcessor.Pipelines[0]
					case ".html":
						pipeline = contentProcessor.Pipelines[1]
					default:
						// Check for a plugin pipeline
						for _, p := range contentProcessor.Pipelines[3:] {
							grpcPipeline, ok := p.(*pipelines.GRPC)
							if !ok {
								continue
							}
							// This is a bit of a hack, but we need to get the extensions
							// for the pipeline. We'll do this by calling the RegisterPipelines
							// method on the plugin.
							pluginPipelines, err := grpcPipeline.Plugin.RegisterPipelines()
							if err != nil {
								handleError(fmt.Errorf("error getting plugin pipelines: %w", err))
								return
							}

							for _, pp := range pluginPipelines {
								if pp.Name == grpcPipeline.Name() {
									for _, e := range pp.Extensions {
										if e == ext {
											pipeline = p
											break
										}
									}
								}
								if pipeline != nil {
									break
								}
							}
							if pipeline != nil {
								break
							}
						}
					}

					// If no pipeline was found, use the copy pipeline
					if pipeline == nil {
						pipeline = contentProcessor.Pipelines[2]
					}

					logger.Logger.Debug("Processing asset", "path", asset.Path, "pipeline", pipeline.Name())
					processedAsset, err = pipeline.Process(&asset)
					if err != nil {
						handleError(fmt.Errorf("pipeline error for %s: %w", asset.Path, err))
						return
					}

					if filepath.Ext(processedAsset.Path) == ".html" {
						layouts := getLayouts(processedAsset.Path, contentProcessor.Partials)
						buf := new(bytes.Buffer)
						if _, err := buf.ReadFrom(processedAsset.Content); err != nil {
							handleError(fmt.Errorf("buffer read error for %s: %w", asset.Path, err))
							return
						}
						processedContent, err := processLayouts(layouts, buf.Bytes(), processedAsset.Metadata, contentProcessor.Partials, loadedConfig)
						if err != nil {
							handleError(fmt.Errorf("layout error for %s: %w", asset.Path, err))
							return
						}

						outputPath := filepath.Join(contentProcessor.OutputDir, processedAsset.Path[len("content"):])
						if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
							handleError(err)
							return
						}

						// Check if the file exists
						if _, err := os.Stat(outputPath); os.IsNotExist(err) {
							// If it doesn't exist, write the whole file
							if err := os.WriteFile(outputPath, processedContent, 0644); err != nil {
								handleError(err)
								return
							}
						} else {
							// If it does exist, apply a patch
							existingContent, err := os.ReadFile(outputPath)
							if err != nil {
								handleError(err)
								return
							}

							newContent, err := diff.Merge(existingContent, processedContent)
							if err != nil {
								handleError(err)
								return
							}

							if err := os.WriteFile(outputPath, newContent, 0644); err != nil {
								handleError(err)
								return
							}
						}
					}
				}
			}
		}()
	}

	// Start file walker in a separate goroutine
	go func() {
		defer close(jobs)
		err := filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if ctx.Err() != nil {
				return ctx.Err() // Stop walking if context is cancelled
			}
			if !info.IsDir() && info.Name()[0] != '_' {
				if _, ok := toRebuild[path]; ok {
					content, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					select {
					case jobs <- pipelines.Asset{Path: path, Content: bytes.NewReader(content)}:
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			}
			return nil
		})
		if err != nil {
			handleError(err)
		}
	}()

	wg.Wait()
	close(errs)

	// Return the first error that occurred
	return <-errs
}

// RunOnPostBuildHooks runs the OnPostBuild hooks for the given plugins.
func RunOnPostBuildHooks(loadedPlugins []plugins.Plugin) error {
	logger.Logger.Debug("Running OnPostBuild hooks...")
	for _, p := range loadedPlugins {
		logger.Logger.Debug("Running OnPostBuild hook", "plugin", p.Name())
		if err := p.OnPostBuild(); err != nil {
			return fmt.Errorf("error running OnPostBuild hook: %w", err)
		}
	}
	return nil
}

// getLayouts returns the layouts for a given path.
func getLayouts(path string, p *partials.Partials) []string {
	// This is a simplified version of the original getLayouts function.
	// A more robust implementation would cache the layouts.
	var layouts []string
	currentDir := filepath.Dir(path)
	for {
		layoutPath := filepath.Join(currentDir, "_layout.html")
		if _, err := os.Stat(layoutPath); err == nil {
			layouts = append(layouts, layoutPath)
		}
		if currentDir == "content" || currentDir == "." || currentDir == "/" {
			break
		}
		currentDir = filepath.Dir(currentDir)
	}
	if len(layouts) == 0 {
		return []string{"default"}
	}
	return layouts
}

// processLayouts processes the layouts for a given content file.
func processLayouts(layouts []string, content []byte, frontMatter map[string]any, p *partials.Partials, config map[string]any) ([]byte, error) {
	processedContent := content

	for _, layoutPath := range layouts {
		layoutContent := new(bytes.Buffer)

		t, err := p.Clone()
		if err != nil {
			return nil, err
		}

		if layoutPath == "default" {
			if _, err = t.Template.Parse(defaults.Layout); err != nil {
				return nil, err
			}
		} else {
			if _, err = t.Template.ParseFiles(layoutPath); err != nil {
				return nil, err
			}
		}

		data := struct {
			Site    map[string]any
			Page    map[string]any
			Content template.HTML
		}{
			Site:    config,
			Page:    frontMatter,
			Content: template.HTML(processedContent),
		}

		if layoutPath == "default" {
			if err := t.Template.Execute(layoutContent, data); err != nil {
				return nil, err
			}
		} else {
			if err := t.Template.ExecuteTemplate(layoutContent, filepath.Base(layoutPath), data); err != nil {
				return nil, err
			}
		}
		processedContent = layoutContent.Bytes()
	}

	return processedContent, nil
}

// getFilesToRebuild returns a map of files to rebuild.
func getFilesToRebuild(c *cache.Cache, d *dag.Graph) (map[string]bool, error) {
	toRebuild := make(map[string]bool)

	for path, node := range d.Nodes {
		h, err := hash.New(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		if c.Get(path) != h {
			toRebuild[path] = true
			// If a partial is modified, we need to rebuild all of its dependents
			if filepath.Base(filepath.Dir(path)) == "partials" {
				for _, dependent := range d.GetDependents(node) {
					toRebuild[dependent.Path] = true
				}
			}
		}
		c.Set(path, h)
	}

	return toRebuild, nil
}

// Build builds the site.
func Build(outputDir string, clean bool) error {
	// If clean is true, remove the cache file
	if clean {
		if err := os.Remove(filepath.Join(outputDir, ".cache")); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("error removing cache: %w", err)
			}
		}
	}

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
