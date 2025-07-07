# Creating Extensions

Evoke's extension system allows you to hook into the build process, add new commands to the CLI, or both.

## Build Process Hooks

To create an extension that hooks into the build process, you'll need to create a Go plugin that exports a variable named `EvokeExtension`. This variable must implement the `extensions.Extension` interface, which has two methods: `BeforeBuild` and `AfterBuild`.

Here's an example of a simple extension that prints a message before and after the build:

```go
package main

import "fmt"

type MyExtension struct{}

func (e *MyExtension) BeforeBuild() error {
    fmt.Println("Before build from my extension!")
    return nil
}

func (e *MyExtension) AfterBuild() error {
    fmt.Println("After build from my extension!")
    return nil
}

var EvokeExtension = &MyExtension{}

func main() {}
```

## CLI Commands

You can also create extensions that add new commands to the Evoke CLI. To do this, your extension will need to export a variable named `Commands` that is a slice of `*cli.Command` pointers.

Here's an example of an extension that adds a `serve` command:

```go
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

func main() {}
```

## Combined Extensions

An extension can provide both build hooks and CLI commands. To do this, simply export both the `EvokeExtension` and `Commands` variables from your plugin.
