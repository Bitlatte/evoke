# Building and Installing Plugins

This guide explains how to build your plugins into executables that Evoke can use and how to install them in your project.

## Building Your Plugin

Evoke plugins are compiled as standalone executables. To build your plugin, navigate to your project's root directory and use the `go build` command.

### Example

Let's say you have a plugin located in `plugins/my-plugin/main.go`. You would build it with the following command:

```bash
go build -o plugins/my-plugin/my-plugin plugins/my-plugin/main.go
```

This command does the following:

- `go build`: The standard Go command to compile packages and dependencies.
- `-o plugins/my-plugin/my-plugin`: This specifies the output file name and location. By convention, place the compiled plugin directly in the plugin's directory.
- `plugins/my-plugin/main.go`: This is the path to your plugin's source code.

## Installing Your Plugin

Once you have built your plugin, it is already "installed" and ready to be used by Evoke. Evoke automatically discovers and loads any executable files found in the `plugins` directory.

There are no further steps required. The next time you run an `evoke` command, your plugin's hooks will be active.

## Cross-Compilation

If you are developing a plugin that you want to distribute to others, you will need to compile it for different operating systems and architectures. You can do this by setting the `GOOS` and `GOARCH` environment variables before running the `go build` command.

For example, to build your plugin for Windows, you would run the following command:

```bash
GOOS=windows GOARCH=amd64 go build -o plugins/my-plugin/my-plugin.exe plugins/my-plugin/main.go
```

To build your plugin for Linux, you would run the following command:

```bash
GOOS=linux GOARCH=amd64 go build -o plugins/my-plugin/my-plugin plugins/my-plugin/main.go
```

## Distributing Your Plugin

Once you have built your plugin for different operating systems and architectures, you can distribute it to others. The easiest way to do this is to create a zip file containing the compiled plugin and any other assets that it needs.

You can then share this zip file with others, and they can install it by unzipping it into their project's `plugins` directory.
