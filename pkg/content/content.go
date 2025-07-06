package content

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
)

func ProcessHTML(path string, config map[string]interface{}, templates *template.Template) error {
	// Read the content of the HTML file
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Create a new template, and parse the content of the file
	tmpl, err := template.New("").Parse(string(content))
	if err != nil {
		return err
	}

	// Execute the template with the config
	var processedContent bytes.Buffer
	err = tmpl.Execute(&processedContent, config)
	if err != nil {
		return err
	}

	// Determine the output path
	outputPath := filepath.Join("dist", path[len("content"):])
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Write the processed content to the output file
	return os.WriteFile(outputPath, processedContent.Bytes(), 0644)
}

func ProcessMarkdown(path string, config map[string]interface{}, templates *template.Template) error {
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

	// Determine the template name
	templateName := "post.html" // This will need to be more dynamic later

	// Execute the template with the config and content
	var processedContent bytes.Buffer
	err = templates.ExecuteTemplate(&processedContent, templateName, map[string]interface{}{
		"Content": template.HTML(buf.String()),
		"Evoke":   config,
	})
	if err != nil {
		return err
	}

	// Determine the output path
	outputPath := filepath.Join("dist", path[len("content"):len(path)-3]+".html")
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	// Write the processed content to the output file
	return os.WriteFile(outputPath, processedContent.Bytes(), 0644)
}
