# go-service-template

Go service template

| Feature               | Notes |
|-----------------------|-------|
| Server                | HTTP and gRPC servers, or request-reply worker |
|                       | - gRPC, HTTP: <http://buf.build>, <https://connectrpc.com> |
|                       | - Worker: <https://nats.io> |
| Insights              | Opentelemetry tracing support (HTTP, gRPC) |
|                       | Prometheus metrics |
| Build                 | Makefile |
|                       | Github Actions |
| Logging               | slog  |
| E2E Testing           | E2E testing skeleton |

## Usage

### Create a new service

This repository is intended to serve as a template for new services. To use it, follow these steps:

- Fork this repository.
- Find & replace `github.com/leonardinius/go-service-template` with you package name and `go-service-template` with your service name.
- Create initial tags `v0.0.1` and set `Read and write permissions` under `Settings` / `Actions` / `General` / `Workflow permissions`.

### Development

Run make to see the available commands:

```bash
$ make
Usage: make <target>
 Default
        help                  Display this help
        all                   Formats, builds and tests the codebase
 Build/Run
        clean                 Format all go files
        gen                   Runs all codegen and docs tasks
        test                  Runs all tests (excluding e2e)
        e2e                   Runs e2e tests
        lint                  Runs all linters
        build                 Builds all artifacts
        run                   Runs service locally. Use ARGS="" make run to pass arguments
        watch                 Runs in watch mode. Example: `make watch ARGS="http"`
```

## Structure

```raw
.
├── api                 <- API definitions
│   ├── docs              - generated API documentation
│   └── proto             - protobuf files (gRPC, HTTP)
├── app                 <- Application
│   ├── cmd
│   │   └── ...           - CLI commands
│   └── main.go           - main entry point
├── internal            <- packages (internal; as app does not expose any packages)
│   ├── apigen            - gRPC generated code (buf.dev)
│   ├── apiserv           - API server (gRPC, HTTP)
│   ├── insights          - Opentelemetry tracing, Prometheus metrics
│   ├── log               - slog logging
│   ├── natsworker        - NATS.io worker
│   └── services          - API implementation (Business logic)
└── teste2e             <- E2E testing
    ├── internal          - internal packages
    ├── natsworkere2e     - NATS.io worker E2E testing
    └── servehttpe2e      - HTTP server E2E testing
```
