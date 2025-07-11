// The main package for the evoke command.
package main

import (
	"context"
	"os"

	"github.com/Bitlatte/evoke/pkg/build"
	init_pkg "github.com/Bitlatte/evoke/pkg/init"
	"github.com/Bitlatte/evoke/pkg/logger"
	"github.com/Bitlatte/evoke/pkg/serve"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

var version = "dev"

// main is the entry point for the evoke command. It sets up the CLI commands
// and executes them.
func main() {
	cmd := &cli.Command{
		Name:    "evoke",
		Usage:   "a powerful little static site generator",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "builds your content into static HTML",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "verbose",
						Usage: "Enable verbose logging",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Bool("verbose") {
						logger.Logger.SetLevel(log.DebugLevel)
					}
					logger.Logger.Info("Starting build...")
					err := build.Build("dist")
					if err != nil {
						logger.Logger.Error("Build failed", "error", err)
						return err
					}
					logger.Logger.Info("Build complete!")
					return nil
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
					if cmd.Bool("verbose") {
						logger.Logger.SetLevel(log.DebugLevel)
					}
					logger.Logger.Info("Starting server...", "port", port)
					return serve.Serve(port)
				},
			},
			{
				Name:  "init",
				Usage: "Initialize a new project",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return init_pkg.Run()
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Logger.Fatal(err)
	}
}
