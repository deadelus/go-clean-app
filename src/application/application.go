// Package context provides the application context and lifecycle management.
package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/deadelus/go-clean-app/src/lifecycle"
	"github.com/deadelus/go-clean-app/src/logger"
	"github.com/joho/godotenv"
)

const (
	// AppNameEnvName is the environment variable name for the application name.
	AppNameEnvName = "APP_NAME"
	// AppVersionEnvName is the default environment variable name for the application version.
	AppVersionEnvName = "APP_VERSION"
	// LoggerModeEnvName is the environment variable name for the logger mode.s
	LoggerModeEnvName = "APP_ENV"
)

// Application interface defines the methods for the application context.
type Application interface {
	Name() string
	Version() string
	Context() context.Context
	Gracefull() lifecycle.Lifecycle
	Logger() logger.Logger
	CurrentUser() string
	UserAgent() string
}

// Engine is the main application structure that implements the Application interface.
// It manages the application lifecycle, logging, and context.
// It also handles graceful shutdown and signal handling.
// The Engine can be extended with additional options for configuration.
type Engine struct {
	appName, appVersion string
	ctx                 context.Context
	gracefull           lifecycle.Lifecycle
	logger              logger.Logger
}

// Option is a function that configures the Engine.
type Option func(*Engine) error

// Force interface compliance
// Ensure that Engine implements the Application interface.
var _ Application = &Engine{}

// New creates a new Engine instance with the specified application name and version.
// It initializes the context, logger, and graceful shutdown manager.
// It also sets up signal handling for graceful shutdown.
func New(appNameEnvName string, version VersionOption, options ...Option) (*Engine, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-c:
			signal.Stop(c)
			cancel()
		case <-ctx.Done():
			// Context cancelled, do nothing
		}
	}()

	var appName string
	if appNameEnv := os.Getenv(appNameEnvName); appNameEnv != "" {
		appName = appNameEnv
	} else {
		appName = "application"
	}

	engine := &Engine{
		appName:   appName,
		ctx:       ctx,
		gracefull: lifecycle.NewGracefullShutdown(ctx),
	}

	if err := version(engine); err != nil {
		return nil, fmt.Errorf("failed to set application version: %w", err)
	}

	if engine.CLIMode() {
		if err := setZapLoggerForCLI()(engine); err != nil {
			return nil, fmt.Errorf("failed to set zap logger for cli: %w", err)
		}
	} else {
		if err := setZapLogger()(engine); err != nil {
			return nil, fmt.Errorf("failed to set zap logger: %w", err)
		}
	}

	if err := setZapLoggerForCLI()(engine); err != nil {
		return nil, fmt.Errorf("failed to set zap logger for cli: %w", err)
	}

	for _, option := range options {
		if err := option(engine); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return engine, nil
}

// Name returns the name of the application.
func (e *Engine) Name() string {
	return e.appName
}

// Context returns the context of the application.
func (e *Engine) Context() context.Context {
	return e.ctx
}

// Gracefull returns the lifecycle manager for graceful shutdown.
func (e *Engine) Gracefull() lifecycle.Lifecycle {
	return e.gracefull
}

// CurrentUser returns the current user of the application.
func (e *Engine) CurrentUser() string {
	// Implement logic to retrieve the current user
	return "default-user"
}

// UserAgent returns the user agent of the application.
func (e *Engine) UserAgent() string {
	// Implement logic to retrieve the user agent
	return "default-user-agent"
}

// CLIMode checks if the application is running in CLI mode.
func (e *Engine) CLIMode() bool {
	return e.ctx.Value("cli_mode") == true
}

// Logger returns the logger instance for the application.
func (e *Engine) Logger() logger.Logger {
	return e.logger
}
