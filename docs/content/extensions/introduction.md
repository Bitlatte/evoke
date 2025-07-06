# Introduction to Extensions

Evoke extensions allow you to extend the functionality of the static site generator. Extensions are written in Go and are loaded as plugins.

## Extension Hooks

Extensions can hook into the build process using the following methods:

*   `BeforeBuild()`: This method is called before the build process begins.
*   `AfterBuild()`: This method is called after the build process has completed.

## Extension Interface

All extensions must implement the `Extension` interface:

```go
package extensions

// Extension defines the interface for an evoke extension.
type Extension interface {
	BeforeBuild() error
	AfterBuild() error
}
```