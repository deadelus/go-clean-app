# Go Clean App

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Test Coverage](https://img.shields.io/badge/coverage-100.0%25-brightgreen)](https://github.com/deadelus/go-clean-app)
[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/deadelus/go-clean-app)

A Go library that provides a ready-to-use application architecture with lifecycle management, structured logging, and graceful shutdown.

## 📋 Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Examples](#examples)
- [Architecture](#architecture)
- [API Reference](#api-reference)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- **Lifecycle Management**: Graceful application startup and shutdown
- **Structured Logging**: Integration with Zap for high-performance logging
- **Signal Handling**: Automatic capture of system signals (SIGTERM, SIGINT)
- **Flexible Configuration**: Environment variables support with godotenv
- **CLI and Server Modes**: Suitable for different types of applications
- **Clean Interface**: Simple and extensible API

## 🚀 Installation

```bash
go get github.com/deadelus/go-clean-app
```

## 🏃 Quick Start

### Basic Application

```go
package main

import (
    "log"
    "github.com/deadelus/go-clean-app/src/application"
)

func main() {
    // Create a new application
    app, err := application.New(
        application.AppNameEnvName,
        application.SetOptionVersion("1.0.0"),
        application.SetZapLogger(),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Use the logger
    app.Logger().Info("Application started", 
        map[string]interface{}{
            "name": app.Name(),
            "version": app.Version(),
        },
    )

    // Register a cleanup function
    app.Gracefull().Register("cleanup", func() error {
        app.Logger().Info("Cleanup in progress...")
        return nil
    })

    // Wait for graceful shutdown
    <-app.Context().Done()
    app.Logger().Info("Application stopped")
}
```

### With Custom Options

```go
package main

import (
    "log"
    "github.com/deadelus/go-clean-app/src/application"
)

func main() {
    app, err := application.New(
        "MY_APP_NAME", // Environment variable for the name
        application.SetVersionFromEnv(), // Version from APP_VERSION
        application.SetZapLogger(),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Your application logic here...
    
    <-app.Context().Done()
}
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file at the root of your project:

```env
# Application name
APP_NAME=my-application

# Application version
APP_VERSION=1.2.3

# Logging mode (development/production)
APP_ENV=development
```

### Supported Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | "application" |
| `APP_VERSION` | Application version | - |
| `APP_ENV` | Application mode (dev/prod) | "development" |

## 📚 Examples

### CLI Application with Subcommands

```go
package main

import (
    "context"
    "log"
    "github.com/deadelus/go-clean-app/src/application"
)

func main() {
    app, err := application.New(
        "CLI_APP_NAME",
        application.SetOptionVersion("1.0.0"),
        application.WithCLIMode(),
        application.SetZapLoggerForCLI(),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Simulate a CLI task
    go func() {
        app.Logger().Info("Executing CLI task...")
        // Your logic here
        app.Logger().Info("Task completed")
    }()

    <-app.Context().Done()
}
```

### HTTP Server Application

```go
package main

import (
    "context"
    "net/http"
    "time"
    "log"
    "github.com/deadelus/go-clean-app/src/application"
)

func main() {
    app, err := application.New(
        "HTTP_SERVER_NAME",
        application.SetOptionVersion("1.0.0"),
        application.SetZapLogger(),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Configure HTTP server
    srv := &http.Server{
        Addr:    ":8080",
        Handler: http.DefaultServeMux,
    }

    // Register server shutdown
    app.Gracefull().Register("http-server", func() error {
        app.Logger().Info("Shutting down HTTP server...")
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        return srv.Shutdown(ctx)
    })

    // Start the server
    go func() {
        app.Logger().Info("HTTP server started", map[string]interface{}{
            "addr": srv.Addr,
        })
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            app.Logger().Error("HTTP server error", map[string]interface{}{
                "error": err.Error(),
            })
        }
    }()

    <-app.Context().Done()
}
```

## 🏗️ Architecture

```
src/
├── application/          # Main application module
│   ├── application.go   # Main interface and engine
│   ├── version.go       # Version management
│   ├── logger.go        # Production logger configuration
│   └── cli-logger.go    # CLI logger configuration
├── lifecycle/           # Lifecycle management
│   └── lifecycle.go     # Graceful shutdown
├── logger/              # Logging interfaces
│   ├── logger.go        # Main interface
│   └── zaplogger/       # Zap implementation
│       └── zaplogger.go
└── cerr/               # Error handling
    └── error.go
```

### Main Components

- **Application Engine**: Main entry point managing the lifecycle
- **Lifecycle Manager**: Graceful shutdown orchestration
- **Logger**: Abstraction for structured logging
- **Context Management**: Context and signal management

## 📖 API Reference

### Application Interface

```go
type Application interface {
    Name() string                    // Application name
    Version() string                 // Application version  
    Context() context.Context        // Main context
    Gracefull() lifecycle.Lifecycle  // Graceful shutdown manager
    Logger() logger.Logger           // Structured logger
    CurrentUser() string             // Current user
    UserAgent() string               // User agent
}
```

### Version Options

```go
// Static version
application.SetOptionVersion("1.0.0")

// Version from environment variable
application.SetVersionFromEnv() // Uses APP_VERSION

// Version from custom variable
application.SetVersionFromSpecifiedEnv("MY_VERSION_VAR")
```

### Logger Methods

```go
logger.Info(msg string, fields ...interface{})    // Info level log
logger.Error(msg string, fields ...interface{})   // Error level log
logger.Debug(msg string, fields ...interface{})   // Debug level log
logger.Warn(msg string, fields ...interface{})    // Warning level log
logger.Close()                                     // Close logger
```

### Lifecycle Manager

```go
// Register a graceful shutdown function
app.Gracefull().Register("service-name", func() error {
    // Cleanup logic
    return nil
})
```

## 🔧 Development

### Prerequisites

- Go 1.24.5+
- Go modules enabled

### Installation for Development

```bash
git clone https://github.com/deadelus/go-clean-app.git
cd go-clean-app
go mod download
```

### Testing

```bash
go test -cover ./...
```

### Test Coverage

We aim for high test coverage to ensure the quality and reliability of the library. You can generate a coverage report using the following command:

```bash
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```


### Vendor

To create a vendor directory with all dependencies:

```bash
go mod vendor
```

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Zap](https://github.com/uber-go/zap) - High-performance structured logger
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading
- [multierr](https://github.com/uber-go/multierr) - Multiple error handling

---

**Note**: This library aims to provide a solid foundation for developing Go applications with clean architecture and integrated best practices.