# Creating Extensions

To create an extension, you'll need to create a new Go project.

1.  **Create a new directory for your extension:**

    ```bash
    mkdir my-extension
    cd my-extension
    ```

2.  **Initialize a new Go module:**

    ```bash
    go mod init my-extension
    ```

3.  **Create a new Go file:**

    Create a file named `main.go` with the following content:

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

    **Note:** The `main` function is required, but it doesn't need to do anything.
