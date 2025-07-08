# Building and Installing Extensions

This guide explains how to build your Go-based extensions into plugins that Evoke can use and how to install them in your project.

## Building Your Extension

Evoke extensions are compiled as Go plugins. To build your extension, navigate to your project's root directory and use the `go build` command with the `-buildmode=plugin` flag.

### Example

Let's say you have an extension located in `extensions/my-extension/main.go`. You would build it with the following command:

```bash
go build -buildmode=plugin -o extensions/my-extension.so extensions/my-extension/main.go
```

This command does the following:

- `go build`: The standard Go command to compile packages and dependencies.
- `-buildmode=plugin`: This flag tells the Go compiler to create a shared object (`.so`) file that can be loaded by a Go program at runtime.
- `-o extensions/my-extension.so`: This specifies the output file name and location. By convention, place the compiled plugin directly in the `extensions` directory.
- `extensions/my-extension/main.go`: This is the path to your extension's source code.

## Installing Your Extension

Once you have built your extension, it is already "installed" and ready to be used by Evoke. Evoke automatically discovers and loads any `.so` files found in the `extensions` directory.

There are no further steps required. The next time you run an `evoke` command, your extension's hooks and CLI commands will be available.

## The `extension get` Command (Future Feature)

You may have noticed the `evoke extension get` command. This command is intended to simplify the process of downloading and installing extensions from remote repositories (like GitHub).

**Please note that this feature is not yet fully implemented.**

In the future, you will be able to install extensions with a single command, like this:

```bash
evoke extension get github.com/user/my-evoke-extension
```

For now, please follow the manual build and installation process described above.
