# Contributing Code

If you're a developer, you can contribute to the project by writing code. To get started, you'll need to fork the [Evoke repository](https://github.com/Bitlatte/evoke) on GitHub and clone it to your local machine.

## Development Workflow

1.  Create a new branch for your changes.
2.  Make your changes and commit them with a clear and descriptive message.
3.  Push your changes to your fork.
4.  Create a pull request to the `main` branch of the Evoke repository.

When creating a pull request, please include a detailed description of your changes and any relevant information that will help us to review them.

## Code Style

We use the standard Go code style, so please make sure your code is formatted with `gofmt` before submitting a pull request. We also use `golangci-lint` to lint our code, so please make sure your code passes all of the linter checks.

## Testing

We have a comprehensive test suite that we use to ensure the quality of our code. Before submitting a pull request, please make sure that all of the tests pass. You can run the tests with the following command:

```bash
go test ./...
