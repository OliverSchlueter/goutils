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

## Installation

```bash
go get github.com/OliverSchlueter/goutils
```