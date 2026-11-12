# Guardian

Gateway service for the Slotty platform.

## Stack

Go · Go modules · Docker · Docker Compose · NGINX

## Development requirements

- Docker Engine
- Docker Compose plugin

For local Go development

- Go 1.27.2 or later (as declared in `go.mod`)
- Make for the shortcut commands below

## Quick start

Clone the repository and enter its directory:

```sh
git clone https://github.com/hqdumbcoders/guardian.git
cd guardian
```

Create a local environment file from the example:

```sh
cp example.env .env
```

Start the application and NGINX:

```sh
docker compose up --build -d
```

Guardian will be available through NGINX on port `80`.

Check application health:

```sh
curl -i -X GET http://localhost/health
```

A successful response returns `HTTP 200 OK`.

Check container status:

```sh
docker compose ps
```

Guardian should have the `healthy` status.

Stop the application:

```sh
docker compose down
```

To view application logs:

```sh
docker compose logs -f
```

## Configuration

The application requires the following environment variables:

Variable=Description 

ERDTREE_URL=URL of the Erdtree Core service    
GLADOS_URL=URL of the Glados Identity service 

The `.env` file is used for local configuration 
and must not be committed to the repository. 
Use `example.env` as a template

## Development 

When running Guardian locally, load the 
required environment variables from .env 
before executing make run. Docker Compose 
loads the variables from .env for container 
configuration.

Run from the repository root:

```sh
make init  # Download dependencies declared in go.mod
make run   # Run the Go application locally
```

The module is already initialized; you do not need to run `go mod init` again.
There are currently no external dependencies to download.

Without Make:

```sh
go mod download
go run ./cmd/guardian
```

The application requires ERDTREE_URL and GLADOS_URL to be set when running locally

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

## Project structure 

```text
cmd/guardian/main.go       Executable entry point and startup
internal/auth/             Authentication package placeholder
internal/config/           Application configuration       
internal/proxy/            Proxy package placeholder
nginx/nginx.conf           NGINX configuration
Dockerfile                 Multi-stage Docker build
docker-compose.yml         Application and NGINX services
example.env                Environment variables template
```

The current module path is `hqdumbcoders/guardian`. Imports of packages in this
module use that prefix. `internal/` keeps these packages private to the project.

Agent development instructions are in `AGENTS.md`.
