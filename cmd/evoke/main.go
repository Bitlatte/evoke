package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Bitlatte/evoke/pkg/build"
	"github.com/Bitlatte/evoke/pkg/serve"
	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	cmd := &cli.Command{
		Name:    "evoke",
		Usage:   "a powerful little static site generator",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "builds your content into static HTML",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return build.Build()
				},
			},
			{
				Name:  "serve",
				Usage: "Build and serve the site, watching for changes",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "port",
						Value: 8990,
						Usage: "port to serve the site on",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					port := cmd.Value("port").(int)
					fmt.Printf("Serving on port %d\n", port)
					return serve.Serve(port)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
