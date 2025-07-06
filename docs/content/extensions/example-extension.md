# Example Extension

Here is a complete example of a simple "Hello, World" extension.

## `my-extension/main.go`

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

## Building the Extension

```bash
go build -buildmode=plugin -o my-extension.so
```

## Installing the Extension

Copy the `my-extension.so` file to the `extensions` directory of your Evoke site.

## Running the Build

When you run `evoke build`, you will see the following output:

```
Before build from my extension!
Building...
After build from my extension!
```