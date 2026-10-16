# Guardian

Gateway service for the Slotty platform.

## Stack

Go · Go modules

## Development requirements

- Go 1.27.2 or later (as declared in `go.mod`)
- Make for the shortcut commands below

## Quick start

Run from the repository root:

```sh
make init  # Download dependencies declared in go.mod
make run   # Run the current executable
```

The module is already initialized; you do not need to run `go mod init` again.
There are currently no external dependencies to download.

Without Make:

```sh
go mod download
go run ./cmd/guardian
```

There is no HTTP address or endpoint to call yet.

## Development commands

```sh
make build  # Build bin/guardian
make fmt    # Format Go source files
make test   # Run all package tests
make vet    # Run Go's static checks
make check  # Run tests and static checks
```

Tests have not been added yet. A successful `make test` currently confirms that
packages compile, but does not verify application behavior.

## Dependencies and code

```sh
go get <module-path>  # Add or update a dependency
go mod tidy          # Reconcile dependencies with source imports
```

Commit `go.mod` and, when generated, `go.sum` with dependency changes.

```text
cmd/guardian/main.go  Executable entry point and startup
internal/auth/       Authentication package placeholder
internal/proxy/      Proxy package placeholder
```

The current module path is `hqdumbcoders/guardian`. Imports of packages in this
module use that prefix. `internal/` keeps these packages private to the project.

Agent development instructions are in `AGENTS.md`.
