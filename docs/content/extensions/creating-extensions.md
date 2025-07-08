# Creating Extensions

This guide will walk you through the process of creating your first Evoke extension. We'll create a simple extension that adds a new CLI command to your project.

## Prerequisites

Before you begin, make sure you have Go installed on your system.

## Step 1: Create a New Directory

First, create a new directory for your extension inside your project's `extensions` directory. For this example, we'll create a `hello` extension.

```bash
mkdir -p extensions/hello
```

## Step 2: Create the Extension File

Inside the `extensions/hello` directory, create a new file named `main.go`. This file will contain the code for your extension.

```go
package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

// Commands is a slice of *cli.Command pointers that will be added to the Evoke CLI.
var Commands = []*cli.Command{
	{
		Name:  "hello",
		Usage: "prints a friendly greeting",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("Hello from my first extension!")
			return nil
		},
	},
}

func main() {}
```

In this example, we're creating a new command called `hello` that simply prints a message to the console.

## Step 3: Build the Extension

To build the extension, you'll need to compile it as a Go plugin. Run the following command from your project's root directory:

```bash
go build -buildmode=plugin -o extensions/hello.so extensions/hello/main.go
```

This will create a new file named `hello.so` in your `extensions` directory. This is the compiled plugin that Evoke will load.

## Step 4: Run the New Command

Now that your extension is built, you can run the new `hello` command:

```bash
evoke hello
```

You should see the following output:

```
Hello from my first extension!
```

## Build Process Hooks

You can also create extensions that hook into the build process. To do this, your extension will need to export a variable named `EvokeExtension` that implements the `extensions.Extension` interface.

Here's an example of an extension that prints a message before and after the build:

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

## Combined Extensions

An extension can provide both build hooks and CLI commands. To do this, simply export both the `EvokeExtension` and `Commands` variables from your plugin.
