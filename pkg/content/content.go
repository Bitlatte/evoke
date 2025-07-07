package content

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

func parseFrontMatter(content []byte) (map[string]any, []byte, error) {
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

func findLayouts(path string) ([]string, error) {
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
	return layouts, nil
}

func ProcessHTML(path string, config map[string]any, partials *partials.Partials) error {
	// Read the content of the HTML file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Find layouts
	layouts, err := findLayouts(path)
	if err != nil {
		return err
	}

	// Execute the templates
	var processedContent bytes.Buffer
	// If there are layouts, execute them
	if len(layouts) > 0 {
		// Clone the partials template
		t, err := partials.Clone()
		if err != nil {
			return err
		}
		// Parse the layout files into the template set
		_, err = t.ParseFiles(layouts...)
		if err != nil {
			return err
		}
		// Execute the layout
		config["content"] = template.HTML(fileContent)
		err = t.ExecuteTemplate(&processedContent, filepath.Base(layouts[0]), config)
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

func ProcessMarkdown(path string, config map[string]any, partials *partials.Partials) error {
	// Read the content of the Markdown file
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse front matter
	frontMatter, body, err := parseFrontMatter(fileContent)
	if err != nil {
		return err
	}

	// Create a new config map for this file to avoid modifying the global config
	fileConfig := make(map[string]any)
	for k, v := range config {
		fileConfig[k] = v
	}
	for k, v := range frontMatter {
		fileConfig[k] = v
	}

	// Convert Markdown to HTML
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	if err := md.Convert(body, &buf); err != nil {
		return err
	}

	// Find layouts
	layouts, err := findLayouts(path)
	if err != nil {
		return err
	}

	// Execute the templates
	var processedContent bytes.Buffer
	// If there are layouts, execute them
	if len(layouts) > 0 {
		// Clone the partials template
		t, err := partials.Clone()
		if err != nil {
			return err
		}
		// Parse the layout files into the template set
		if _, err := t.ParseFiles(layouts...); err != nil {
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
