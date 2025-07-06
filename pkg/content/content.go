package content

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/Bitlatte/evoke/pkg/partials"
	"github.com/yuin/goldmark"
)

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
	content, err := os.ReadFile(path)
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
		// Parse the content as a template named "content"
		_, err = t.New("content").Parse(string(content))
		if err != nil {
			return err
		}
		// Execute the layout
		err = t.ExecuteTemplate(&processedContent, filepath.Base(layouts[0]), config)
		if err != nil {
			return err
		}
	} else {
		// If there are no layouts, just use the file content
		processedContent.Write(content)
	}

	// Determine the output path
	outputPath := filepath.Join("dist", path[len("content"):])
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Write the processed content to the output file
	return os.WriteFile(outputPath, processedContent.Bytes(), 0644)
}

func ProcessMarkdown(path string, config map[string]any, partials *partials.Partials) error {
	// Read the content of the Markdown file
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Convert Markdown to HTML
	var buf bytes.Buffer
	if err := goldmark.Convert(content, &buf); err != nil {
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
		// Parse the content as a template named "content"
		_, err = t.New("content").Parse(buf.String())
		if err != nil {
			return err
		}
		// Execute the layout
		err = t.ExecuteTemplate(&processedContent, filepath.Base(layouts[0]), config)
		if err != nil {
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
