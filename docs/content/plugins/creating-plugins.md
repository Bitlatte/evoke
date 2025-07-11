# Creating Plugins

This guide will walk you through the process of creating your first Evoke plugin. We'll create a simple plugin that prints a message during the build process.

## Prerequisites

Before you begin, make sure you have the following installed:

- Go
- Protobuf Compiler (`protoc`)
- Go gRPC plugins

You can install the Protobuf Compiler and the Go gRPC plugins with the following commands:

```bash
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

## Step 1: Create a New Directory

First, create a new directory for your plugin inside your project's `plugins` directory. For this example, we'll create a `hello` plugin.

```bash
mkdir -p plugins/hello
```

## Step 2: Create the Plugin File

Inside the `plugins/hello` directory, create a new file named `main.go`. This file will contain the code for your plugin.

```go
package main

import (
	"context"
	"fmt"

	"github.com/Bitlatte/evoke/pkg/plugins"
	"github.com/Bitlatte/evoke/proto"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of the plugin.
type HelloPlugin struct{}

func (p *HelloPlugin) OnPreBuild(ctx context.Context, req *proto.OnPreBuildRequest) (*proto.OnPreBuildResponse, error) {
	fmt.Println("Hello from the OnPreBuild hook!")
	return &proto.OnPreBuildResponse{}, nil
}

func (p *HelloPlugin) OnConfigLoaded(ctx context.Context, req *proto.OnConfigLoadedRequest) (*proto.OnConfigLoadedResponse, error) {
	fmt.Println("Hello from the OnConfigLoaded hook!")
	return &proto.OnConfigLoadedResponse{Config: req.Config}, nil
}

func (p *HelloPlugin) OnPublicAssetsCopied(ctx context.Context, req *proto.OnPublicAssetsCopiedRequest) (*proto.OnPublicAssetsCopiedResponse, error) {
	fmt.Println("Hello from the OnPublicAssetsCopied hook!")
	return &proto.OnPublicAssetsCopiedResponse{}, nil
}

func (p *HelloPlugin) OnContentLoaded(ctx context.Context, req *proto.OnContentLoadedRequest) (*proto.OnContentLoadedResponse, error) {
	fmt.Printf("Hello from the OnContentLoaded hook for %s!\n", req.Path)
	return &proto.OnContentLoadedResponse{Content: req.Content}, nil
}

func (p *HelloPlugin) OnContentRender(ctx context.Context, req *proto.OnContentRenderRequest) (*proto.OnContentRenderResponse, error) {
	fmt.Printf("Hello from the OnContentRender hook for %s!\n", req.Path)
	return &proto.OnContentRenderResponse{Content: req.Content}, nil
}

func (p *HelloPlugin) OnHTMLRendered(ctx context.Context, req *proto.OnHTMLRenderedRequest) (*proto.OnHTMLRenderedResponse, error) {
	fmt.Printf("Hello from the OnHTMLRendered hook for %s!\n", req.Path)
	return &proto.OnHTMLRenderedResponse{Content: req.Content}, nil
}

func (p *HelloPlugin) OnPostBuild(ctx context.Context, req *proto.OnPostBuildRequest) (*proto.OnPostBuildResponse, error) {
	fmt.Println("Hello from the OnPostBuild hook!")
	return &proto.OnPostBuildResponse{}, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		Plugins: map[string]plugin.Plugin{
			"evoke": &plugins.EvokePlugin{Impl: &HelloPlugin{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
```

## Step 3: Build the Plugin

To build the plugin, run the following command from your project's root directory:

```bash
go build -o plugins/hello/hello plugins/hello/main.go
```

This will create a new executable file named `hello` in your `plugins/hello` directory. This is the compiled plugin that Evoke will load.

## Step 4: Run a Build

Now that your plugin is built, run the `evoke build` command:

```bash
evoke build
```

You should see the messages from your plugin printed to the console during the build process.
