// Package init provides the functionality to initialize a new evoke project.
package init

import (
	"fmt"
	"os"

	"github.com/Bitlatte/evoke/pkg/defaults"
	"github.com/charmbracelet/huh"
)

// Run initializes a new project.
func Run() error {
	var projectName string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Value(&projectName),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	directory := projectName
	if projectName == "." {
		directory = "."
		// Check if the current directory is empty.
		files, err := os.ReadDir(".")
		if err != nil {
			return err
		}
		if len(files) > 0 {
			return fmt.Errorf("current directory is not empty")
		}
	}

	fmt.Printf("Project Name: %s, Directory: %s\n", projectName, directory)

	// Create the project directory.
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return err
	}

	// Create subdirectories.
	subdirs := []string{"content", "partials"}
	for _, subdir := range subdirs {
		if err := os.MkdirAll(fmt.Sprintf("%s/%s", directory, subdir), os.ModePerm); err != nil {
			return err
		}
	}

	// Create evoke.yaml.
	evokeYAML := []byte(fmt.Sprintf("site:\n  name: %s\n", projectName))
	if err := os.WriteFile(fmt.Sprintf("%s/evoke.yaml", directory), evokeYAML, 0644); err != nil {
		return err
	}

	// Create content/index.md.
	indexMD := []byte("---\ntitle: Home\n---\n\n# Hello, World!\n")
	if err := os.WriteFile(fmt.Sprintf("%s/content/index.md", directory), indexMD, 0644); err != nil {
		return err
	}

	// Create content/_layout.html.
	if err := os.WriteFile(fmt.Sprintf("%s/content/_layout.html", directory), []byte(defaults.Layout), 0644); err != nil {
		return err
	}

	fmt.Printf("Successfully created project in %s\n", directory)

	return nil
}
