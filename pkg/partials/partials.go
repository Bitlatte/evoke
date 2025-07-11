// Package partials provides functionality for loading and managing HTML partials.
package partials

import (
	"html/template"
	"os"
	"path/filepath"
)

// Partials holds the parsed partial templates.
type Partials struct {
	*template.Template
}

// LoadPartials walks the "partials" directory and parses all the files as templates.
func LoadPartials() (*Partials, error) {
	t := template.New("")
	err := filepath.Walk("partials", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Read the content of the partial file
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Create a new template with the base name of the file
			_, err = t.New(filepath.Base(path)).Parse(string(content))
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Partials{t}, nil
}

// Clone creates a deep copy of the Partials templates.
func (p *Partials) Clone() (*Partials, error) {
	if p.Template == nil {
		return &Partials{template.New("")}, nil
	}
	cloned, err := p.Template.Clone()
	if err != nil {
		return nil, err
	}
	return &Partials{cloned}, nil
}
