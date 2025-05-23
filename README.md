# serve-http

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: serve-http is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

`serve-http` is a command in the [senzing-tools] suite of tools.
This command is an
HTTP server application that supports requests to HTTP applications via network access.

[![Go Reference Badge]][Package reference]
[![Go Report Card Badge]][Go Report Card]
[![License Badge]][License]
[![go-test-linux.yaml Badge]][go-test-linux.yaml]
[![go-test-darwin.yaml Badge]][go-test-darwin.yaml]
[![go-test-windows.yaml Badge]][go-test-windows.yaml]

[![golangci-lint.yaml Badge]][golangci-lint.yaml]

## Overview

`serve-http` supports the

Senzing SDKs for accessing the gRPC server:

1. Go: [sz-sdk-go-grpc]
1. Python: [sz-sdk-python-grpc]

A simple demonstration using `senzing-tools` and a SQLite database.

```console
export LD_LIBRARY_PATH=/opt/senzing/er/lib/
export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db
senzing-tools init-database
senzing-tools serve-http --enable-all

```

Then visit [localhost:8261]

## Install

1. The `serve-http` command is installed with the [senzing-tools] suite of tools.
   See senzing-tools [install].

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/er/lib/
senzing-tools serve-http [flags]
```

1. For options and flags:
    1. [Online documentation]
    1. Runtime documentation:

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/er/lib/
        senzing-tools serve-http --help
        ```

1. In addition to the following simple usage examples, there are additional [Examples].

### Using command line options

1. :pencil2: Specify database using command line option.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    senzing-tools serve-http \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --enable-all

    ```

1. Visit [localhost:8261]
1. Run `senzing-tools serve-http --help` or see [Parameters] for additional parameters.

### Using environment variables

1. :pencil2: Specify database using environment variable.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export SENZING_TOOLS_ENABLE_ALL=true
    senzing-tools serve-http
    ```

1. Visit [localhost:8261]
1. Run `senzing-tools serve-http --help` or see [Parameters] for additional parameters.

### Using Docker

This usage shows how to initialize a database with a Docker container.

1. This usage specifies a URL of an external database.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        --env SENZING_TOOLS_ENABLE_ALL=true \
        --interactive \
        --publish 8261:8261 \
        --rm \
        --tty \
        senzing/senzing-tools serve-http

    ```

1. Visit [localhost:8261]
1. See [Parameters] for additional parameters.

### Parameters

- **[SENZING_TOOLS_DATABASE_URL]**
- **[SENZING_TOOLS_ENABLE_ALL]**
- **[SENZING_TOOLS_ENABLE_SZCONFIG]**
- **[SENZING_TOOLS_ENABLE_SZCONFIGMGR]**
- **[SENZING_TOOLS_ENABLE_SZDIAGNOSTIC]**
- **[SENZING_TOOLS_ENABLE_SZENGINE]**
- **[SENZING_TOOLS_ENABLE_SZPRODUCT]**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL]**
- **[SENZING_TOOLS_ENGINE_MODULE_NAME]**
- **[SENZING_TOOLS_GRPC_PORT]**
- **[SENZING_TOOLS_LOG_LEVEL]**

## References

1. [Command reference]
1. [Development]
1. [Errors]
1. [Examples]

[Command reference]: https://hub.senzing.com/senzing-tools/senzing-tools_serve-http.html
[Development]: docs/development.md
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[Go Reference Badge]: https://pkg.go.dev/badge/github.com/senzing-garage/serve-http.svg
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/senzing-garage/serve-http
[Go Report Card]: https://goreportcard.com/report/github.com/senzing-garage/serve-http
[go-test-darwin.yaml Badge]: https://github.com/senzing-garage/serve-http/actions/workflows/go-test-darwin.yaml/badge.svg
[go-test-darwin.yaml]: https://github.com/senzing-garage/serve-http/actions/workflows/go-test-darwin.yaml
[go-test-linux.yaml Badge]: https://github.com/senzing-garage/serve-http/actions/workflows/go-test-linux.yaml/badge.svg
[go-test-linux.yaml]: https://github.com/senzing-garage/serve-http/actions/workflows/go-test-linux.yaml
[go-test-windows.yaml Badge]: https://github.com/senzing-garage/serve-http/actions/workflows/go-test-windows.yaml/badge.svg
[go-test-windows.yaml]: https://github.com/senzing-garage/serve-http/actions/workflows/go-test-windows.yaml
[golangci-lint.yaml Badge]: https://github.com/senzing-garage/serve-http/actions/workflows/golangci-lint.yaml/badge.svg
[golangci-lint.yaml]: https://github.com/senzing-garage/serve-http/actions/workflows/golangci-lint.yaml
[install]: https://github.com/senzing-garage/senzing-tools#install
[License Badge]: https://img.shields.io/badge/License-Apache2-brightgreen.svg
[License]: https://github.com/senzing-garage/serve-http/blob/main/LICENSE
[localhost:8261]: http://localhost:8261
[Online documentation]: https://hub.senzing.com/senzing-tools/senzing-tools_serve-http.html
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/serve-http
[Parameters]: #parameters
[Senzing Garage]: https://github.com/senzing-garage
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[SENZING_TOOLS_DATABASE_URL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url
[SENZING_TOOLS_ENABLE_ALL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_all
[SENZING_TOOLS_ENABLE_SZCONFIG]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szconfig
[SENZING_TOOLS_ENABLE_SZCONFIGMGR]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szconfigmgr
[SENZING_TOOLS_ENABLE_SZDIAGNOSTIC]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szdiagnostic
[SENZING_TOOLS_ENABLE_SZENGINE]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szengine
[SENZING_TOOLS_ENABLE_SzPRODUCT]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_szproduct
[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json
[SENZING_TOOLS_ENGINE_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level
[SENZING_TOOLS_ENGINE_MODULE_NAME]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_module_name
[SENZING_TOOLS_GRPC_PORT]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_grpc_port
[SENZING_TOOLS_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level
[senzing-tools]: https://github.com/senzing-garage/senzing-tools
[Senzing]: https://senzing.com/
[sz-sdk-go-grpc]: https://github.com/senzing/sz-sdk-go-grpc
[sz-sdk-python-grpc]: https://github.com/senzing-garage/sz-sdk-python-grpc
