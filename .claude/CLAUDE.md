# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`serve-http` is an HTTP server command in the senzing-tools suite that provides:

- Senzing REST API server
- Swagger UI
- Xterm terminal interface

It connects to a Senzing database (local SQLite/PostgreSQL or remote via gRPC) to provide entity resolution services over HTTP.

## Build and Development Commands

```bash
# Install development dependencies (one-time)
make dependencies-for-development

# Install Go dependencies
make dependencies

# Build binaries (output in target/<os>-<arch>/)
make clean build

# Run without building
make run

# Run linting (golangci-lint, govulncheck, cspell)
make lint

# Run tests (requires setup first)
make clean setup test

# Run a single test
export LD_LIBRARY_PATH=/opt/senzing/er/lib
go test -v -run TestFunctionName ./path/to/package

# Run tests with coverage
make clean setup coverage

# Check coverage against thresholds
make check-coverage

# Build Docker image
make docker-build

# Apply lint fixes
make fix
```

## Running the Server

Requires Senzing C library installed at `/opt/senzing/er/lib`. Set `LD_LIBRARY_PATH`:

```bash
export LD_LIBRARY_PATH=/opt/senzing/er/lib
```

The server runs on port 8261 by default. Use `--enable-all` to enable all services.

Configuration via environment variables (prefixed `SENZING_TOOLS_`) or command-line flags:

- `SENZING_TOOLS_DATABASE_URL` - Database connection (e.g., `sqlite3://na:na@/tmp/sqlite/G2C.db`)
- `SENZING_TOOLS_GRPC_URL` - Connect to remote Senzing via gRPC instead of local database
- `SENZING_TOOLS_ENABLE_ALL` - Enable all services (API, Swagger, Xterm)
- `SENZING_TOOLS_HTTP_PORT` - Server port (default 8261)

## Architecture

### Package Structure

- `main.go` - Entry point, calls `cmd.Execute()`
- `cmd/` - Cobra command implementation with Viper configuration
  - `root.go` - Main command definition with `RootCmd`, `PreRun`, `RunE`
  - `context_*.go` - Platform-specific context variables (Xterm command differs by OS)
- `httpserver/` - HTTP server implementation
  - `httpserver_basic.go` - `BasicHTTPServer` struct and `Serve()` method
  - `static/` - Embedded static files and templates

### Key Components

**BasicHTTPServer** (`httpserver/httpserver_basic.go`): Core server struct that configures and runs the HTTP server. The `Serve()` method sets up routes via `http.ServeMux`:

- `/api/` - Senzing REST API (via go-rest-api-service)
- `/swagger/` - Swagger UI
- `/xterm/` - Terminal interface (via cloudshell)
- `/site/` - Static site templates
- `/` - Static root files

**Command Configuration** (`cmd/root.go`): Uses `go-cmdhelping` for CLI options. The `RunE` function builds a `BasicHTTPServer` from Viper config and calls `Serve()`. Supports two modes:

- Local mode: Connects directly to database via `SENZING_TOOLS_DATABASE_URL`
- gRPC mode: Connects to remote Senzing server via `SENZING_TOOLS_GRPC_URL`

### Dependencies

Key Senzing packages:

- `go-rest-api-service` - REST API implementation
- `go-cmdhelping` - CLI helpers and option management
- `go-grpcing` - gRPC URL parsing for remote Senzing connections
- `go-observing` - Observer pattern implementation

### Test Data

Test setup (`make setup`) copies `testdata/sqlite/G2C.db` to `/tmp/sqlite/G2C.db`. Tests use this SQLite database.
