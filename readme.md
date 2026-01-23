# Go Clean App

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Test Coverage](https://img.shields.io/badge/coverage-100.0%25-brightgreen)](https://github.com/deadelus/go-clean-app)
[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/deadelus/go-clean-app)

A lightweight Go library providing a robust application skeleton with lifecycle management, structured logging, and graceful shutdown capabilities.

## 📋 Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [API Reference](#api-reference)
- [License](#license)

## ✨ Features

- **Standard Application Engine**: A ready-to-use `Engine` implementing the `Application` interface.
- **Graceful Shutdown**: Built-in lifecycle management to handle OS signals (SIGTERM, SIGINT) and cleanup tasks.
- **Structured Logging**: Decoupled logger interface with a production-ready Zap implementation.
- **Functional Options**: Clean and extensible configuration via the options pattern.
- **Environment Management**: Categorize your app lifecycle (Development, Staging, Production, Testing).
- **No Side-Effects**: Explicit initialization without magic global states or automatic file loading.

## 🚀 Installation

```bash
go get github.com/deadelus/go-clean-app
```

## 🏃 Quick Start

### Basic Server Application

```go
package main

import (
	"fmt"
	"github.com/deadelus/go-clean-app/v2/application"
	"github.com/deadelus/go-clean-app/v2/lifecycle"
	"github.com/deadelus/go-clean-app/v2/logger/zaplogger"
)

func main() {
	// Initialize the application engine with options
	app := application.New(
		application.AppName("my-service"),
		application.Version("1.0.0"),
		application.Env(application.EnvDevelopment),
		application.Debug(true),
		zaplogger.SetZapLogger(), // Configure Zap as the logger
	)

	// Access common properties
	app.Logger().Info(fmt.Sprintf("Starting %s (%s)", app.AppName(), app.Version()))

	// Register components for graceful shutdown
	app.Gracefull().Register("database", func() error {
		fmt.Println("Closing database connections...")
		return nil
	})

	// Run your logic using the application context
	go func() {
		app.Logger().Info("Application is running...")
		// Use app.Context() for cancellation propagation
	}()

	// The engine automatically listens for SIGINT/SIGTERM
	// Block until shutdown happens
	<-app.Context().Done()

	// Wait for graceful shutdown to complete (recommended)
	<-app.Gracefull().Done()
	
	app.Logger().Info("Shutdown complete")
}
```

### CLI Application

For CLI tools, you can use a optimized logger configuration:

```go
app := application.New(
    application.AppName("my-cli"),
    zaplogger.SetZapLoggerForCLI(),
)
```

## ⚙️ Configuration

Configuration is managed through functional options passed to `application.New()`:

| Option | Description |
|--------|-------------|
| `application.AppName(string)` | Sets the application name. |
| `application.Version(string)` | Sets the application version. |
| `application.Env(string)` | Sets the environment (Development, Production, etc.). |
| `application.Debug(bool)` | Enables/disables debug mode. |
| `zaplogger.SetZapLogger()` | Attaches a Zap-based structured logger. |
| `zaplogger.SetZapLoggerForCLI()` | Attaches a Zap logger optimized for CLI output. |

## 🏗 Architecture

The library follows clean architecture principles by decoupling the core engine from specific implementations:

- **`application`**: Defines the `Application` interface and provides the default `Engine`.
- **`logger`**: Defines the `Logger` interface to keep the application logic agnostic of the logging library.
- **`lifecycle`**: Manages the application state and shutdown hooks.
- **`errors`**: Centralized error constants for the library.

## 📚 API Reference

### Application Interface

```go
type Application interface {
	Gracefull() lifecycle.Lifecycle
	Logger() logger.Logger
	Name() string
	Version() string
	Env() string
	Debug() bool
	Context() context.Context
}
```

### Engine Methods

- `Gracefull()`: Returns the `Lifecycle` manager to register shutdown hooks and wait for shutdown completion (with `Done()`).
- `Context()`: Returns the application context that is canceled when the app shuts down.
- `Logger()`: Returns the configured logger instance.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Copyright © 2026 Geoffrey Trambolho (deadelus)
