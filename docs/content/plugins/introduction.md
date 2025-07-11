# Introduction to Plugins

Evoke's plugin system is a powerful feature that allows you to add new functionality and customize the build process. Plugins are built on a gRPC-based architecture, which means you can write them in any language that supports gRPC, including Go, Python, JavaScript, and more.

## What Can Plugins Do?

With plugins, you can:

- **Modify the build process:** Hook into the build process to perform a wide range of actions, such as fetching data from an API, processing images, or generating a sitemap.
- **Customize content rendering:** Intercept and modify content at various stages of the rendering pipeline. For example, you could add a custom Markdown renderer, or post-process the HTML before it's written to disk.
- **Integrate with other tools:** Connect Evoke to other services and tools, such as content management systems, analytics platforms, or deployment pipelines.

## How Do Plugins Work?

Evoke plugins are standalone executables that communicate with the Evoke application over gRPC. When you run a build, Evoke discovers and launches any plugins in your `plugins` directory, and then communicates with them at various points in the build process.

### Plugin Hooks

Plugins can execute code at specific points in the build process using the following hooks:

- `OnPreBuild()`: This method is called before the build process begins.
- `OnConfigLoaded()`: This method is called after the configuration is loaded, but before it is used.
- `OnPublicAssetsCopied()`: This method is called after the public assets have been copied to the output directory.
- `OnContentLoaded()`: This method is called after a content file has been read from disk, but before it is processed.
- `OnContentRender()`: This method is called before a content file is rendered.
- `OnHTMLRendered()`: This method is called after a content file has been rendered to HTML.
- `OnPostBuild()`: This method is called after the build process has completed.

## The Plugin Interface

All plugins must implement the `Plugin` service, which is defined in the `plugin.proto` file. You can find the full definition of the service and its messages in the [Plugin Service Definition](./plugin-service-definition.html) documentation.

## Getting Started with Plugins

To learn how to create your own plugins, check out the following guides:

- **[Creating Plugins](./creating-plugins.html):** A step-by-step guide to creating your first plugin.
- **[Example Plugin](./example-plugin.html):** A practical example of a plugin that adds a new CLI command.
- **[Building and Installing](./building-and-installing.html):** Learn how to build and install plugins for your projects.

