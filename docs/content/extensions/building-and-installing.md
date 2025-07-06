# Building and Installing Extensions

To build your extension, run the following command:

```bash
go build -buildmode=plugin -o my-extension.so
```

This will create a file named `my-extension.so` in the root of your project.

To install your extension, simply copy the `.so` file to the `extensions` directory of your Evoke site.
