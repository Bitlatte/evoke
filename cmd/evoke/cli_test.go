package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestApp_HasExtensionCommand(t *testing.T) {
	// Arrange
	app := &cli.Command{
		Name: "evoke",
		Commands: []*cli.Command{
			{
				Name: "build",
			},
			{
				Name: "extension",
			},
		},
	}

	// Act
	found := false
	for _, cmd := range app.Commands {
		if cmd.Name == "extension" {
			found = true
			break
		}
	}

	// Assert
	assert.True(t, found, "Expected to find the 'extension' command")
}
