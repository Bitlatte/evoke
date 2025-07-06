package partials

import (
	"html/template"
	"os"
	"path/filepath"
)

type Partials struct {
	*template.Template
}

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

func (p *Partials) Clone() (*Partials, error) {
	cloned, err := p.Template.Clone()
	if err != nil {
		return nil, err
	}
	return &Partials{cloned}, nil
}
