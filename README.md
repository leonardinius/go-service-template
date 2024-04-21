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
| Github Actions        | Github Actions |
| Logging               | slog  |
| E2E Testing           | E2E testing skeleton |

## Usage

This repository is intended to serve as a template for new services. To use it, follow these steps:

- Fork this repository.
- Find & replace `go-service-template` with your service name.

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
