package extensions

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"

	"github.com/urfave/cli/v3"
)

// Extension defines the interface for an evoke extension.
type Extension interface {
	BeforeBuild() error
	AfterBuild() error
}

func LoadBuildExtensions() ([]Extension, error) {
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
				// This is not a build extension, so we can safely ignore it.
				return nil
			}

			// The looked-up symbol is a pointer to the variable.
			// We need to dereference it to get the actual value.
			val := reflect.ValueOf(sym)
			if val.Kind() != reflect.Ptr {
				return fmt.Errorf("EvokeExtension symbol is not a pointer, but %s", val.Type())
			}
			// Dereference the pointer to get the actual extension variable.
			extVal := val.Elem()

			ext, ok := extVal.Interface().(Extension)
			if !ok {
				return fmt.Errorf("unexpected type from module symbol: %T. Expected a type that implements the extensions.Extension interface", extVal.Interface())
			}

			extensions = append(extensions, ext)
		}

		return nil
	})

	return extensions, err
}

// LoadCliCommands loads all the commands from the extensions.
func LoadCliCommands() ([]*cli.Command, error) {
	var commands []*cli.Command

	// Check if the extensions directory exists
	if _, err := os.Stat("extensions"); os.IsNotExist(err) {
		return commands, nil
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

			sym, err := p.Lookup("Commands")
			if err != nil {
				// This is not a CLI extension, so we can safely ignore it.
				return nil
			}

			// The looked-up symbol is a pointer to the variable.
			// We need to dereference it to get the actual value.
			val := reflect.ValueOf(sym)
			if val.Kind() != reflect.Ptr {
				return fmt.Errorf("Commands symbol in %s is not a pointer, but %s", path, val.Type())
			}
			// Dereference the pointer to get the actual extension variable.
			cmdVal := val.Elem()

			cmds, ok := cmdVal.Interface().([]*cli.Command)
			if !ok {
				return fmt.Errorf("unexpected type from module symbol in %s: %T. Expected a slice of cli.Command pointers", path, cmdVal.Interface())
			}

			commands = append(commands, cmds...)
		}

		return nil
	})

	return commands, err
}
