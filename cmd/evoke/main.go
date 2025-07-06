package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "evoke",
		Usage: "a powerful little static site generator",
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "builds your content into static HTML",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return build()
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func build() error {
	fmt.Println("Building...")
	// Create the output directory
	if err := os.MkdirAll("dist", 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}
	return nil
}
