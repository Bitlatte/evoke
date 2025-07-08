# Example Extension: A Simple Web Server

In this example, we'll create a practical extension that adds a `serve` command to Evoke. This command will start a local web server to serve your built site from the `dist` directory, allowing you to preview your changes locally.

## 1. Create the Extension Directory

First, create a directory for the extension:

```bash
mkdir -p extensions/serve
```

## 2. Create the `main.go` File

Next, create a `main.go` file inside `extensions/serve`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/urfave/cli/v3"
)

// Commands exports the `serve` command.
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

This code defines a new `serve` command that uses Go's built-in HTTP server to serve files from the `./dist` directory on port `8080`.

## 3. Build the Extension

Now, build the extension as a Go plugin:

```bash
go build -buildmode=plugin -o extensions/serve.so extensions/serve/main.go
```

This will create the `serve.so` plugin file in your `extensions` directory.

## 4. Use the New Command

First, make sure you have some content in your `dist` directory by running a build:

```bash
evoke build
```

Then, you can use your new `serve` command to preview your site:

```bash
evoke serve
```

You'll see the following output:

```
Serving on http://localhost:8080
```

Now you can open your web browser and navigate to `http://localhost:8080` to see your site in action.
