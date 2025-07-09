package content

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"sync"

	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/Bitlatte/evoke/pkg/plugins"
	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type Content struct {
	Partials      *partials.Partials
	Goldmark      goldmark.Markdown
	Config        map[string]any
	Plugins       []plugins.Plugin
	LayoutCache   sync.Map
	TemplateCache sync.Map
	OutputDir     string
	bufferPool    sync.Pool
}

type templateData struct {
	Global  map[string]any
	Page    map[string]any
	Content template.HTML
}

func New(outputDir string, config map[string]any, partials *partials.Partials, gm goldmark.Markdown, plugins []plugins.Plugin) (*Content, error) {
	return &Content{
		Partials:  partials,
		Config:    config,
		Goldmark:  gm,
		Plugins:   plugins,
		OutputDir: outputDir,
		bufferPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}, nil
}

func (c *Content) ParseFrontMatter(content []byte) (map[string]any, []byte, error) {
	var frontMatter map[string]any
	body := content

	if bytes.HasPrefix(content, []byte("---")) {
		end := bytes.Index(content[3:], []byte("---"))
		if end != -1 {
			err := yaml.Unmarshal(content[3:end+3], &frontMatter)
			if err != nil {
				return nil, nil, err
			}
			body = bytes.TrimSpace(content[end+6:])
		}
	}

	return frontMatter, body, nil
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

func (c *Content) ProcessHTML(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get a buffer from the pool
	fileContent := c.bufferPool.Get().(*bytes.Buffer)
	fileContent.Reset()
	defer c.bufferPool.Put(fileContent)

	// Read the file into the buffer
	_, err = fileContent.ReadFrom(file)
	if err != nil {
		return err
	}

	// Run OnContentLoaded hooks
	for _, p := range c.Plugins {
		fileContentBytes, err := p.OnContentLoaded(path, fileContent.Bytes())
		if err != nil {
			return err
		}
		fileContent.Reset()
		fileContent.Write(fileContentBytes)
	}

	// Find layouts
	layouts := c.GetLayouts(path)

	// Get a buffer from the pool for the final processed content
	processedContent, err := c.ProcessLayouts(layouts, fileContent.Bytes(), nil)
	if err != nil {
		return err
	}

	// Run OnHTMLRendered hooks
	for _, p := range c.Plugins {
		processedContent, err = p.OnHTMLRendered(path, processedContent)
		if err != nil {
			return err
		}
	}

	// Determine the output path
	outputPath := filepath.Join(c.OutputDir, path[len("content"):])
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Write the processed content to the output file
	return os.WriteFile(outputPath, processedContent, 0644)
}

func (c *Content) ProcessMarkdown(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get buffers from the pool
	bodyContent := c.bufferPool.Get().(*bytes.Buffer)
	bodyContent.Reset()
	defer c.bufferPool.Put(bodyContent)

	// Read the file into the buffer
	_, err = bodyContent.ReadFrom(file)
	if err != nil {
		return err
	}

	// Run OnContentLoaded hooks
	for _, p := range c.Plugins {
		bodyContentBytes, err := p.OnContentLoaded(path, bodyContent.Bytes())
		if err != nil {
			return err
		}
		bodyContent.Reset()
		bodyContent.Write(bodyContentBytes)
	}

	// Parse front matter
	frontMatter, body, err := c.ParseFrontMatter(bodyContent.Bytes())
	if err != nil {
		return err
	}

	// Run OnContentRender hooks
	for _, p := range c.Plugins {
		body, err = p.OnContentRender(path, body)
		if err != nil {
			return err
		}
	}

	// Convert Markdown to HTML
	mdOutput := c.bufferPool.Get().(*bytes.Buffer)
	mdOutput.Reset()
	defer c.bufferPool.Put(mdOutput)
	if err := c.Goldmark.Convert(body, mdOutput); err != nil {
		return err
	}

	// Find layouts
	layouts := c.GetLayouts(path)

	// Get a buffer from the pool for the final processed content
	processedContent, err := c.ProcessLayouts(layouts, mdOutput.Bytes(), frontMatter)
	if err != nil {
		return err
	}

	// Run OnHTMLRendered hooks
	for _, p := range c.Plugins {
		processedContent, err = p.OnHTMLRendered(path, processedContent)
		if err != nil {
			return err
		}
	}

	// Determine the output path
	outputPath := filepath.Join(c.OutputDir, path[len("content"):len(path)-3]+".html")
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Write the processed content to the output file
	return os.WriteFile(outputPath, processedContent, 0644)
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
