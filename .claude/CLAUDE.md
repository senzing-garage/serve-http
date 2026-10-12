# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`serve-http` is an HTTP server command in the senzing-tools suite that provides:

- Xterm terminal interface
- Static console site

The Senzing REST API and Swagger UI were removed in 0.4.0, along with the archived `go-rest-api-service` dependency.
The Senzing settings (from `SENZING_TOOLS_DATABASE_URL` or `SENZING_TOOLS_CORE_SETTINGS`) are passed to the Xterm shell as `SENZING_ENGINE_CONFIGURATION_JSON`.

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

The binary does not link the Senzing C library. Senzing tools run from the Xterm shell need it at `/opt/senzing/er/lib`, so set `LD_LIBRARY_PATH`:

```bash
export LD_LIBRARY_PATH=/opt/senzing/er/lib
```

The server runs on port 8261 by default. Use `--enable-all` to enable all services.

Configuration via environment variables (prefixed `SENZING_TOOLS_`) or command-line flags:

- `SENZING_TOOLS_DATABASE_URL` - Database connection (e.g., `sqlite3://na:na@/tmp/sqlite/G2C.db`)
- `SENZING_TOOLS_ENABLE_ALL` - Enable all services (Xterm)
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

- `/xterm/` - Terminal interface (via cloudshell)
- `/site/` - Static site templates
- `/` - Static root files

**Command Configuration** (`cmd/root.go`): Uses `go-cmdhelping` for CLI options. The `RunE` function builds a `BasicHTTPServer` from Viper config and calls `Serve()`. Senzing settings are built from `SENZING_TOOLS_DATABASE_URL` or `SENZING_TOOLS_CORE_SETTINGS`.

### Dependencies

Key Senzing packages:

- `go-cmdhelping` - CLI helpers and option management
- `go-helpers` - Settings building and error wrapping

### Test Data

Test setup (`make setup`) copies `testdata/sqlite/G2C.db` to `/tmp/sqlite/G2C.db`. Tests use this SQLite database.
