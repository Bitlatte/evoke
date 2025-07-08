# Introduction to Extensions

Evoke's extension system is a powerful feature that allows you to add new functionality and customize the build process. Extensions are written in Go and are loaded as plugins, giving you full access to the power of the Go language.

## What Can Extensions Do?

With extensions, you can:

- **Add new CLI commands:** Create custom commands to automate tasks, manage your project, or integrate with other tools.
- **Modify the build process:** Hook into the build process to perform actions before or after your site is built. For example, you could fetch data from an API, process images, or generate a sitemap.
- **Integrate with other tools:** Connect Evoke to other services and tools, such as content management systems, analytics platforms, or deployment pipelines.

## How Do Extensions Work?

Evoke extensions are Go plugins that are loaded at runtime. They can hook into the build process and add new commands to the CLI.

### Extension Hooks

Extensions can execute code at specific points in the build process using the following hooks:

- `BeforeBuild()`: This method is called before the build process begins.
- `AfterBuild()`: This method is called after the build process has completed.

### Custom CLI Commands

Extensions can also add new commands to the Evoke CLI. This is done by exporting a `Commands` variable from your extension, which is a slice of `*cli.Command` pointers.

## The Extension Interface

All extensions must implement the `Extension` interface:

```go
package extensions

// Extension defines the interface for an evoke extension.
type Extension interface {
	BeforeBuild() error
	AfterBuild() error
}
```

## Getting Started with Extensions

To learn how to create your own extensions, check out the following guides:

- **[Creating Extensions](./creating-extensions.md):** A step-by-step guide to creating your first extension.
- **[Example Extension](./example-extension.md):** A practical example of an extension that adds a new CLI command.
- **[Building and Installing](./building-and-installing.md):** Learn how to build and install extensions for your projects.
