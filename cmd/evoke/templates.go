package main

import (
	"html/template"
	"os"
	"path/filepath"
)

func loadTemplates() (*template.Template, error) {
	templates := template.New("")
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			_, err := templates.ParseFiles(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return templates, nil
}
