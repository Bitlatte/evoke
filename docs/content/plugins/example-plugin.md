# Example Plugin: Modify Content

In this example, we'll create a plugin that modifies the content of a page before it's rendered. This plugin will find all instances of the word "Hello" and replace them with "Hello from our plugin!".

## 1. Create the Plugin Directory

First, create a directory for the plugin:

```bash
mkdir -p plugins/modifier
```

## 2. Create the `main.go` File

Next, create a `main.go` file inside `plugins/modifier`:

```go
package main

import (
	"bytes"
	"context"

	"github.com/Bitlatte/evoke/pkg/plugins"
	"github.com/Bitlatte/evoke/proto"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of the plugin.
type ModifierPlugin struct{}

func (p *ModifierPlugin) OnPreBuild(ctx context.Context, req *proto.OnPreBuildRequest) (*proto.OnPreBuildResponse, error) {
	return &proto.OnPreBuildResponse{}, nil
}

func (p *ModifierPlugin) OnConfigLoaded(ctx context.Context, req *proto.OnConfigLoadedRequest) (*proto.OnConfigLoadedResponse, error) {
	return &proto.OnConfigLoadedResponse{Config: req.Config}, nil
}

func (p *ModifierPlugin) OnPublicAssetsCopied(ctx context.Context, req *proto.OnPublicAssetsCopiedRequest) (*proto.OnPublicAssetsCopiedResponse, error) {
	return &proto.OnPublicAssetsCopiedResponse{}, nil
}

func (p *ModifierPlugin) OnContentLoaded(ctx context.Context, req *proto.OnContentLoadedRequest) (*proto.OnContentLoadedResponse, error) {
	newContent := bytes.ReplaceAll(req.Content, []byte("Hello"), []byte("Hello from our plugin!"))
	return &proto.OnContentLoadedResponse{Content: newContent}, nil
}

func (p *ModifierPlugin) OnContentRender(ctx context.Context, req *proto.OnContentRenderRequest) (*proto.OnContentRenderResponse, error) {
	return &proto.OnContentRenderResponse{Content: req.Content}, nil
}

func (p *ModifierPlugin) OnHTMLRendered(ctx context.Context, req *proto.OnHTMLRenderedRequest) (*proto.OnHTMLRenderedResponse, error) {
	return &proto.OnHTMLRenderedResponse{Content: req.Content}, nil
}

func (p *ModifierPlugin) OnPostBuild(ctx context.Context, req *proto.OnPostBuildRequest) (*proto.OnPostBuildResponse, error) {
	return &proto.OnPostBuildResponse{}, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.Handshake,
		Plugins: map[string]plugin.Plugin{
			"evoke": &plugins.EvokePlugin{Impl: &ModifierPlugin{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
```

This code defines a new plugin that implements the `OnContentLoaded` hook. This hook will be called for each content file after it's read from disk, and it will replace all instances of "Hello" with "Hello from our plugin!".

## 3. Build the Plugin

Now, build the plugin as a Go executable:

```bash
go build -o plugins/modifier/modifier plugins/modifier/main.go
```

This will create the `modifier` executable file in your `plugins/modifier` directory.

## 4. Use the New Plugin

First, create a content file with the word "Hello" in it. For example, create `content/index.html` with the following content:

```html
<h1>Hello, World!</h1>
```

Then, run a build:

```bash
evoke build
```

Now, if you open the `dist/index.html` file, you should see that the content has been modified:

```html
<h1>Hello from our plugin!, World!</h1>
