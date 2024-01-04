# serve-http

If you are beginning your journey with
[Senzing](https://senzing.com/),
please start with
[Senzing Quick Start guides](https://docs.senzing.com/quickstart/).

You are in the
[Senzing Garage](https://github.com/senzing-garage)
where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: serve-http is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

`serve-http` is a command in the
[senzing-tools](https://github.com/senzing-garage/senzing-tools)
suite of tools.
This command is an
HTTP server application that supports requests to HTTP applications via network access.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing-garage/serve-http.svg)](https://pkg.go.dev/github.com/senzing-garage/serve-http)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing-garage/serve-http)](https://goreportcard.com/report/github.com/senzing-garage/serve-http)
[![License](https://img.shields.io/badge/License-Apache2-brightgreen.svg)](https://github.com/senzing-garage/serve-http/blob/main/LICENSE)

[![gosec.yaml](https://github.com/senzing-garage/serve-http/actions/workflows/gosec.yaml/badge.svg)](https://github.com/senzing-garage/serve-http/actions/workflows/gosec.yaml)
[![go-test-linux.yaml](https://github.com/senzing-garage/serve-http/actions/workflows/go-test-linux.yaml/badge.svg)](https://github.com/senzing-garage/serve-http/actions/workflows/go-test-linux.yaml)
[![go-test-darwin.yaml](https://github.com/senzing-garage/serve-http/actions/workflows/go-test-darwin.yaml/badge.svg)](https://github.com/senzing-garage/serve-http/actions/workflows/go-test-darwin.yaml)
[![go-test-windows.yaml](https://github.com/senzing-garage/serve-http/actions/workflows/go-test-windows.yaml/badge.svg)](https://github.com/senzing-garage/serve-http/actions/workflows/go-test-windows.yaml)

## Overview

`serve-http` supports the

Senzing SDKs for accessing the gRPC server:

1. Go: [g2-sdk-go-grpc](https://github.com/senzing/g2-sdk-go-grpc)
1. Python: [g2-sdk-python-grpc](https://github.com/senzing-garage/g2-sdk-python-grpc)

A simple demonstration using `senzing-tools` and a SQLite database.

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db
senzing-tools init-database
senzing-tools serve-http --enable-all

```

Then visit [localhost:8261](http://localhost:8261)

## Install

1. The `serve-http` command is installed with the
   [senzing-tools](https://github.com/senzing-garage/senzing-tools)
   suite of tools.
   See senzing-tools [install](https://github.com/senzing-garage/senzing-tools#install).

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
senzing-tools serve-http [flags]
```

1. For options and flags:
    1. [Online documentation](https://hub.senzing.com/senzing-tools/senzing-tools_serve-http.html)
    1. Runtime documentation:

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
        senzing-tools serve-http --help
        ```

1. In addition to the following simple usage examples, there are additional [Examples](docs/examples.md).

### Using command line options

1. :pencil2: Specify database using command line option.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools serve-http \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --enable-all

    ```

1. Visit [localhost:8261](http://localhost:8261)
1. Run `senzing-tools serve-http --help` or see [Parameters](#parameters) for additional parameters.

### Using environment variables

1. :pencil2: Specify database using environment variable.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export SENZING_TOOLS_ENABLE_ALL=true
    senzing-tools serve-http
    ```

1. Visit [localhost:8261](http://localhost:8261)
1. Run `senzing-tools serve-http --help` or see [Parameters](#parameters) for additional parameters.

### Using Docker

This usage shows how to initialze a database with a Docker container.

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

1. Visit [localhost:8261](http://localhost:8261)
1. See [Parameters](#parameters) for additional parameters.

### Parameters

- **[SENZING_TOOLS_DATABASE_URL](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url)**
- **[SENZING_TOOLS_ENABLE_ALL](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_all)**
- **[SENZING_TOOLS_ENABLE_G2CONFIG](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2config)**
- **[SENZING_TOOLS_ENABLE_G2CONFIGMGR](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2configmgr)**
- **[SENZING_TOOLS_ENABLE_G2DIAGNOSTIC](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2diagnostic)**
- **[SENZING_TOOLS_ENABLE_G2ENGINE](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2engine)**
- **[SENZING_TOOLS_ENABLE_G2PRODUCT](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_enable_g2product)**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json)**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level)**
- **[SENZING_TOOLS_ENGINE_MODULE_NAME](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_module_name)**
- **[SENZING_TOOLS_GRPC_PORT](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_grpc_port)**
- **[SENZING_TOOLS_LOG_LEVEL](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level)**

## References

1. [Command reference](https://hub.senzing.com/senzing-tools/senzing-tools_serve-http.html)
1. [Development](docs/development.md)
1. [Errors](docs/errors.md)
1. [Examples](docs/examples.md)
