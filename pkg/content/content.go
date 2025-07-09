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

type Content struct {
	Partials      *partials.Partials
	Goldmark      goldmark.Markdown
	Config        map[string]any
	Plugins       []plugins.Plugin
	LayoutCache   sync.Map
	TemplateCache sync.Map
	Pipelines     []pipelines.Pipeline
	OutputDir     string
	bufferPool    sync.Pool
}

type templateData struct {
	Global  map[string]any
	Page    map[string]any
	Content template.HTML
}

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
