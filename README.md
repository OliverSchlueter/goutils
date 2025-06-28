# GoUtils

Just a collection of useful Go utilities for various tasks.

## Features

- **sloki**: structured logger that supports multiple output formats
- **broker**: an abstraction layer for message brokers (e.g. for Nats)
- **middleware**: a collection of commonly used middlewares
- **featureflags**: a simple feature flag implementation
- **container**: start and stop docker containers for testing purposes
- **cloudevents**: utilities for working with CloudEvents ([1.0.2](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md))
- **problem**: a structured error handling package ([RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) compliant)
- **healthcheck**: a health check handler for HTTP servers

## Installation

```bash
go get github.com/OliverSchlueter/goutils
```

## sloki

```go
lokiService := sloki.NewService(sloki.Configuration{
    URL:          "http://localhost:3100/loki/api/v1/push",
    Service:      "my-service",
    ConsoleLevel: slog.LevelDebug,
    LokiLevel:    slog.LevelInfo,
    EnableLoki:   true,
})
slog.SetDefault(slog.New(lokiService))

slog.Info("Hello, world!", "key", "value")
```

The field `limits_config.allow_structured_metadata` in the loki configuration must be set to `true` to allow structured metadata.