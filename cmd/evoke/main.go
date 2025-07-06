package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Bitlatte/evoke/pkg/build"
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
					return build.Build()
				},
			},
			{
				Name:  "extension",
				Usage: "manage extensions",
				Commands: []*cli.Command{
					{
						Name:  "get",
						Usage: "get a new extension from a url",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							fmt.Println("Getting extension...")
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
