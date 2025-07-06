package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

func processHTML(path string, config map[string]interface{}, templates *template.Template) error {
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
