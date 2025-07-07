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

func New(config map[string]any, partials *partials.Partials) *Content {
	return &Content{
		partials: partials,
		config:   config,
		goldmark: goldmark.New(
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		),
	}
}

func (c *Content) parseFrontMatter(content []byte) (map[string]any, []byte, error) {
	var frontMatter map[string]any
	body := content

	if strings.HasPrefix(string(content), "---") {
		parts := strings.SplitN(string(content), "---", 3)
		if len(parts) >= 3 {
			err := yaml.Unmarshal([]byte(parts[1]), &frontMatter)
			if err != nil {
				return nil, nil, err
			}
			body = []byte(strings.TrimSpace(parts[2]))
		}
	}

	return frontMatter, body, nil
}

func (c *Content) findLayouts(path string) ([]string, error) {
	if layouts, ok := c.layoutCache.Load(path); ok {
		return layouts.([]string), nil
	}

	var layouts []string
	dir := filepath.Dir(path)
	for {
		layoutPath := filepath.Join(dir, "_layout.html")
		if _, err := os.Stat(layoutPath); err == nil {
			layouts = append(layouts, layoutPath)
		}
		if dir == "content" || dir == "." || dir == "/" {
			break
		}
		dir = filepath.Dir(dir)
	}
	// Reverse the layouts slice so that the outermost layout is first
	for i, j := 0, len(layouts)-1; i < j; i, j = i+1, j-1 {
		layouts[i], layouts[j] = layouts[j], layouts[i]
	}

	c.layoutCache.Store(path, layouts)
	return layouts, nil
}

func (c *Content) ProcessHTML(path string) error {
	// Read the content of the HTML file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Find layouts
	layouts, err := c.findLayouts(path)
	if err != nil {
		return err
	}

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
		c.config["content"] = template.HTML(fileContent)
		err = t.ExecuteTemplate(&processedContent, filepath.Base(layouts[0]), c.config)
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

	// Create a new config map for this file to avoid modifying the global config
	fileConfig := make(map[string]any)
	for k, v := range c.config {
		fileConfig[k] = v
	}
	for k, v := range frontMatter {
		fileConfig[k] = v
	}

	// Convert Markdown to HTML
	var buf bytes.Buffer
	if err := c.goldmark.Convert(body, &buf); err != nil {
		return err
	}

	// Find layouts
	layouts, err := c.findLayouts(path)
	if err != nil {
		return err
	}

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
		fileConfig["content"] = template.HTML(buf.String())
		if err := t.ExecuteTemplate(&processedContent, filepath.Base(layouts[0]), fileConfig); err != nil {
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

	if _, err := t.Template.ParseFiles(layouts...); err != nil {
		return nil, err
	}

	c.templateCache.Store(cacheKey, t.Template)
	return t.Template, nil
}
