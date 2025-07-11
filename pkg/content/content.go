package content

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"sync"

	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/Bitlatte/evoke/pkg/pipelines"
	"github.com/Bitlatte/evoke/pkg/plugins"
	"github.com/yuin/goldmark"
)

// Content is the main struct for handling content processing.
type Content struct {
	// Partials are the HTML partials that can be used in layouts.
	Partials *partials.Partials
	// Goldmark is the Goldmark instance used for rendering Markdown.
	Goldmark goldmark.Markdown
	// Config is the site configuration.
	Config map[string]any
	// Plugins are the plugins that are currently loaded.
	Plugins []plugins.Plugin
	// LayoutCache is a cache of layouts that have been found for a given directory.
	LayoutCache sync.Map
	// TemplateCache is a cache of templates that have been parsed.
	TemplateCache sync.Map
	// Pipelines are the content pipelines that are currently loaded.
	Pipelines []pipelines.Pipeline
	// OutputDir is the directory where the site will be built.
	OutputDir  string
	bufferPool sync.Pool
}

// templateData is the data that is passed to the layout templates.
type templateData struct {
	Global  map[string]any
	Page    map[string]any
	Content template.HTML
}

// New creates a new Content struct.
func New(outputDir string, config map[string]any, partials *partials.Partials, gm goldmark.Markdown, plugins []plugins.Plugin, pipelines []pipelines.Pipeline) (*Content, error) {
	return &Content{
		Partials:  partials,
		Config:    config,
		Goldmark:  gm,
		Plugins:   plugins,
		Pipelines: pipelines,
		OutputDir: outputDir,
		bufferPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}, nil
}

// GetLayouts returns the layouts for a given path.
func (c *Content) GetLayouts(path string) []string {
	dir := filepath.Dir(path)
	if layouts, ok := c.LayoutCache.Load(dir); ok {
		return layouts.([]string)
	}

	var layouts []string
	currentDir := dir
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
	c.LayoutCache.Store(dir, layouts)
	return layouts
}

// ProcessLayouts processes the layouts for a given content file.
func (c *Content) ProcessLayouts(layouts []string, content []byte, frontMatter map[string]any) ([]byte, error) {
	processedContent := content

	for _, layoutPath := range layouts {
		// Get a buffer from the pool
		layoutContent := c.bufferPool.Get().(*bytes.Buffer)
		layoutContent.Reset()
		defer c.bufferPool.Put(layoutContent)

		// Get the template from the cache or parse it
		t, err := c.GetTemplate(layoutPath)
		if err != nil {
			return nil, err
		}

		// Execute the layout
		data := templateData{
			Global:  c.Config,
			Page:    frontMatter,
			Content: template.HTML(processedContent),
		}

		if err := t.ExecuteTemplate(layoutContent, filepath.Base(layoutPath), data); err != nil {
			return nil, err
		}
		processedContent = layoutContent.Bytes()
	}

	return processedContent, nil
}

// GetTemplate returns a template from the cache or parses it if it's not in the cache.
func (c *Content) GetTemplate(layout string) (*template.Template, error) {
	if t, ok := c.TemplateCache.Load(layout); ok {
		return t.(*template.Template), nil
	}

	t, err := c.Partials.Clone()
	if err != nil {
		return nil, err
	}
	if _, err = t.Template.ParseFiles(layout); err != nil {
		return nil, err
	}
	c.TemplateCache.Store(layout, t.Template)
	return t.Template, nil
}
