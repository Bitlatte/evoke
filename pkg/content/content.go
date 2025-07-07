package content

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

type Content struct {
	Partials      *partials.Partials
	Layouts       sync.Map
	Goldmark      goldmark.Markdown
	Config        map[string]any
	LayoutCache   sync.Map
	TemplateCache sync.Map
	bufferPool    sync.Pool
}

type templateData struct {
	Global  map[string]any
	Page    map[string]any
	Content template.HTML
}

func New(config map[string]any, partials *partials.Partials) (*Content, error) {
	return &Content{
		Partials: partials,
		Config:   config,
		Goldmark: goldmark.New(
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		),
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
	// Reverse the layouts slice so that the outermost layout is first
	for i, j := 0, len(layouts)-1; i < j; i, j = i+1, j-1 {
		layouts[i], layouts[j] = layouts[j], layouts[i]
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

	// Find layouts
	layouts := c.GetLayouts(path)

	// Get a buffer from the pool
	processedContent := c.bufferPool.Get().(*bytes.Buffer)
	processedContent.Reset()
	defer c.bufferPool.Put(processedContent)

	// If there are layouts, execute them
	if len(layouts) > 0 {
		// Get the template from the cache or parse it
		t, err := c.GetTemplate(layouts)
		if err != nil {
			return err
		}

		// Read the file content into a buffer to pass to the template
		fileContent := c.bufferPool.Get().(*bytes.Buffer)
		fileContent.Reset()
		defer c.bufferPool.Put(fileContent)
		_, err = fileContent.ReadFrom(file)
		if err != nil {
			return err
		}

		// Execute the layout
		data := templateData{
			Global:  c.Config,
			Content: template.HTML(fileContent.String()),
		}
		err = t.ExecuteTemplate(processedContent, filepath.Base(layouts[0]), data)
		if err != nil {
			return err
		}
	} else {
		// If there are no layouts, just copy the file content
		_, err := processedContent.ReadFrom(file)
		if err != nil {
			return err
		}
	}

	// Determine the output path
	outputPath := filepath.Join("dist", path[len("content"):])
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Write the processed content to the output file
	return os.WriteFile(outputPath, processedContent.Bytes(), 0644)
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

	// Parse front matter
	frontMatter, body, err := c.ParseFrontMatter(bodyContent.Bytes())
	if err != nil {
		return err
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
	processedContent := c.bufferPool.Get().(*bytes.Buffer)
	processedContent.Reset()
	defer c.bufferPool.Put(processedContent)

	// If there are layouts, execute them
	if len(layouts) > 0 {
		// Get the template from the cache or parse it
		t, err := c.GetTemplate(layouts)
		if err != nil {
			return err
		}
		// Execute the layout
		data := templateData{
			Global:  c.Config,
			Page:    frontMatter,
			Content: template.HTML(mdOutput.String()),
		}
		if err := t.ExecuteTemplate(processedContent, filepath.Base(layouts[0]), data); err != nil {
			return err
		}
	} else {
		// If there are no layouts, just use the file content
		processedContent.Write(mdOutput.Bytes())
	}

	// Determine the output path
	outputPath := filepath.Join("dist", path[len("content"):len(path)-3]+".html")
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Write the processed content to the output file
	return os.WriteFile(outputPath, processedContent.Bytes(), 0644)
}

func (c *Content) GetTemplate(layouts []string) (*template.Template, error) {
	cacheKey := strings.Join(layouts, ",")
	if t, ok := c.TemplateCache.Load(cacheKey); ok {
		return t.(*template.Template), nil
	}

	t, err := c.Partials.Clone()
	if err != nil {
		return nil, err
	}
	if _, err = t.Template.ParseFiles(layouts...); err != nil {
		return nil, err
	}
	c.TemplateCache.Store(cacheKey, t.Template)
	return t.Template, nil
}
