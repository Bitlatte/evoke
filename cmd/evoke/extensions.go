package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

// Extension defines the interface for an evoke extension.
type Extension interface {
	BeforeBuild() error
	AfterBuild() error
}

func loadExtensions() ([]Extension, error) {
	var extensions []Extension

	// Check if the extensions directory exists
	if _, err := os.Stat("extensions"); os.IsNotExist(err) {
		return extensions, nil
	}

	err := filepath.Walk("extensions", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".so" {
			p, err := plugin.Open(path)
			if err != nil {
				return fmt.Errorf("could not open plugin at %s: %w", path, err)
			}

			sym, err := p.Lookup("EvokeExtension")
			if err != nil {
				return fmt.Errorf("could not find EvokeExtension symbol in %s: %w", path, err)
			}

			var ext Extension
			switch v := sym.(type) {
			case Extension:
				ext = v
			case *Extension:
				ext = *v
			default:
				return fmt.Errorf("unexpected type from module symbol: %T", sym)
			}

			extensions = append(extensions, ext)
		}

		return nil
	})

	return extensions, err
}
