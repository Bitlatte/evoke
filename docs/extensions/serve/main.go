package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/urfave/cli/v3"
)

var Commands = []*cli.Command{
	{
		Name:  "serve",
		Usage: "serve the dist directory",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fs := http.FileServer(http.Dir("./dist"))
			http.Handle("/", fs)

			fmt.Println("Serving on http://localhost:8080")
			log.Fatal(http.ListenAndServe(":8080", nil))
			return nil
		},
	},
}
