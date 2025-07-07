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
	partials      *partials.Partials
	layouts       sync.Map
	goldmark      goldmark.Markdown
	config        map[string]any
	layoutCache   sync.Map
	templateCache sync.Map
}

type templateData struct {
	Global  map[string]any
	Page    map[string]any
	Content template.HTML
}

func New(config map[string]any, partials *partials.Partials) (*Content, error) {
	return &Content{
		partials: partials,
		config:   config,
		goldmark: goldmark.New(
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		),
	}, nil
}

func (c *Content) parseFrontMatter(content []byte) (map[string]any, []byte, error) {
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

func (c *Content) getLayouts(path string) []string {
	dir := filepath.Dir(path)
	if layouts, ok := c.layoutCache.Load(dir); ok {
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
	c.layoutCache.Store(dir, layouts)
	return layouts
}

func (c *Content) ProcessHTML(path string) error {
	// Read the content of the HTML file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Find layouts
	layouts := c.getLayouts(path)

	// Execute the templates
	var processedContent bytes.Buffer
	// If there are layouts, execute them
	if len(layouts) > 0 {
		// Get the template from the cache or parse it
		t, err := c.getTemplate(layouts)
		if err != nil {
			return err
		}
		// Execute the layout
		data := templateData{
			Global:  c.config,
			Content: template.HTML(fileContent),
		}
		err = t.ExecuteTemplate(&processedContent, filepath.Base(layouts[0]), data)
		if err != nil {
			return err
		}
	} else {
		// If there are no layouts, just use the file content
		processedContent.Write(fileContent)
	}

	// Determine the output path
	outputPath := filepath.Join("dist", path[len("content"):])
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Write the processed content to the output file
	return os.WriteFile(outputPath, processedContent.Bytes(), 0644)
}

func (c *Content) ProcessMarkdown(path string) error {
	// Read the content of the Markdown file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse front matter
	frontMatter, body, err := c.parseFrontMatter(fileContent)
	if err != nil {
		return err
	}

	// Convert Markdown to HTML
	var buf bytes.Buffer
	if err := c.goldmark.Convert(body, &buf); err != nil {
		return err
	}

	// Find layouts
	layouts := c.getLayouts(path)

	// Execute the templates
	var processedContent bytes.Buffer
	// If there are layouts, execute them
	if len(layouts) > 0 {
		// Get the template from the cache or parse it
		t, err := c.getTemplate(layouts)
		if err != nil {
			return err
		}
		// Execute the layout
		data := templateData{
			Global:  c.config,
			Page:    frontMatter,
			Content: template.HTML(buf.String()),
		}
		if err := t.ExecuteTemplate(&processedContent, filepath.Base(layouts[0]), data); err != nil {
			return err
		}
	} else {
		// If there are no layouts, just use the file content
		processedContent.Write(buf.Bytes())
	}

	// Determine the output path
	outputPath := filepath.Join("dist", path[len("content"):len(path)-3]+".html")
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Write the processed content to the output file
	return os.WriteFile(outputPath, processedContent.Bytes(), 0644)
}

func (c *Content) getTemplate(layouts []string) (*template.Template, error) {
	cacheKey := strings.Join(layouts, ",")
	if t, ok := c.templateCache.Load(cacheKey); ok {
		return t.(*template.Template), nil
	}

	t, err := c.partials.Clone()
	if err != nil {
		return nil, err
	}
	if _, err = t.Template.ParseFiles(layouts...); err != nil {
		return nil, err
	}
	c.templateCache.Store(cacheKey, t.Template)
	return t.Template, nil
}
