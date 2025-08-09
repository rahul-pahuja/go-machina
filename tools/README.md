# GoMachina Tools

This directory contains tools used for code generation in the GoMachina project.

## Mock Generation

The project uses `go generate` to automatically generate mock functions for testing. To regenerate the mocks, run:

```bash
go generate ./machina
```

This will generate the `mocks_test.go` file in the `machina` package with all the mock functions needed for testing.

## Adding New Mocks

To add new mock functions, edit the `generate_mocks.go` file and add a new entry to the `mockFunctions` slice with the following fields:

- `Name`: The name of the mock function
- `Description`: A comment describing the mock function (optional)
- `Signature`: The function signature
- `Body`: The function body

After adding new mocks, regenerate the mocks file by running `go generate ./machina`.