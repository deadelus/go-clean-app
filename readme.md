# Go Clean App

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/deadelus/go-clean-app)](https://goreportcard.com/report/github.com/deadelus/go-clean-app)
[![Version](https://img.shields.io/badge/version-2.1.0-blue.svg)](https://github.com/deadelus/go-clean-app)

A lightweight Go library providing a robust application skeleton with lifecycle management, structured logging, and graceful shutdown capabilities. Built for modern Go services, it helps you focus on your business logic while handling the "plumbing" of a production-ready application.

## 📋 Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Transport Adapters](#transport-adapters)
- [Project Structure](#project-structure)
- [License](#license)

## ✨ Features

- **Standard Application Engine**: A ready-to-use `Engine` implementing the `Application` interface.
- **Graceful Shutdown**: Built-in lifecycle management to handle OS signals (SIGTERM, SIGINT) and cleanup tasks.
- **Structured Logging**: Decoupled logger interface with a production-ready Zap implementation.
- **Functional Options**: Clean and extensible configuration via the options pattern.
- **Environment Management**: Categorize your app lifecycle (Development, Staging, Production, Testing).
- **Transport Abstraction**: Easily run your application as a local HTTP server or an AWS Lambda function.
- **Explicit over Magic**: No side-effects, no global states, and no automatic file loading.

## 🚀 Installation

```bash
go get github.com/deadelus/go-clean-app/v2
```

## 🏃 Quick Start

### Basic Server Application

See [main.go](main.go) for a complete example. You can run it with:

```bash
go run .
```

### CLI Application

For CLI tools, use the optimized CLI logger which provides a cleaner output:

```go
app, _ := application.New(
    application.AppName("my-cli"),
    zaplogger.SetZapLoggerForCLI(),
    application.WithCLIMode(),
)
```

## ⚙️ Configuration

Configuration is managed through functional options passed to `application.New()`:

| Option | Description | Default |
|--------|-------------|---------|
| `AppName(string)` | Sets the application name. | `"application"` |
| `Version(string)` | Sets the application version. | `"2.1.0"` |
| `Env(string)` | Sets the environment (use `application.EnvDevelopment`, etc.). | `application.EnvDevelopment` |
| `Debug(bool)` | Enables/disables debug mode. | `false` |
| `WithCLIMode()` | Opt-in for CLI-specific behavior. | `false` |
| `zaplogger.SetZapLogger()` | Attaches a Zap-based structured logger. | - |
| `zaplogger.SetZapLoggerForCLI()` | Attaches a Zap logger optimized for CLI. | - |

## 🌐 Transport Adapters

`go-clean-app` provides adapters to run the same business logic in different environments:

### Local HTTP Server
Ideal for local development or Docker-based deployments.
```go
server := local.NewAdapter(handler, 8080)
server.Start()
```

### AWS Lambda (API Gateway)
Perfect for serverless deployments on AWS. It uses `aws-lambda-go-api-proxy` to wrap your standard `http.Handler`.
```go
import "github.com/deadelus/go-clean-app/v2/transport/adapter/apigateway"

lambdaAdapter := apigateway.NewAdapter(handler)
lambdaAdapter.Start()
```

## 🏗 Project Structure

- `application/`: Core engine and configuration options.
- `lifecycle/`: Graceful shutdown management.
- `logger/`: Generic logging interface and Zap implementation.
- `transport/`: Abstractions for HTTP, Lambda, and more.
- `errors/`: Custom error handling utilities.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Copyright © 2026 Geoffrey Trambolho (deadelus)
