// Package context provides the application context and lifecycle management.
package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/deadelus/go-clean-app/v2/lifecycle"
	"github.com/deadelus/go-clean-app/v2/logger"
)

const (
	// AppNameEnvName is the environment variable name for the application name.
	AppNameEnvName = "APP_NAME"
	// AppVersionEnvName is the default environment variable name for the application version.
	AppVersionEnvName = "APP_VERSION"
	// LoggerModeEnvName is the environment variable name for the logger mode.s
	LoggerModeEnvName = "APP_ENV"
	// AppDebugEnvName is the environment variable name for the debug mode.
	AppDebugEnvName = "APP_DEBUG"
)

// Application interface defines the methods for the application context.
type Application interface {
	Name() string
	Version() string
	Env() string
	Debug() bool
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
	appName, appVersion, appEnv string
	appDebug                    bool
	ctx                         context.Context
	gracefull                   lifecycle.Lifecycle
	logger                      logger.Logger
}

// Force interface compliance
// Ensure that Engine implements the Application interface.
var _ Application = &Engine{}

// New creates a new Engine instance with the specified application name and version.
// It initializes the context, logger, and graceful shutdown manager.
// It also sets up signal handling for graceful shutdown.
func New(options ...Option) (*Engine, error) {
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

	engine := &Engine{
		ctx:       ctx,
		gracefull: lifecycle.NewGracefullShutdown(ctx),
	}

	for _, option := range options {
		option(engine)
	}

	if engine.appName == "" {
		engine.appName = "application"
	}

	if engine.appVersion == "" {
		engine.appVersion = "0.1.0"
	}

	if engine.appEnv == "" {
		engine.appEnv = "development"
	}

	return engine, nil
}

// Name returns the name of the application.
func (e *Engine) Name() string {
	return e.appName
}

// Version returns the version of the application.
func (e *Engine) Version() string {
	return e.appVersion
}

// Env returns the environment of the application.
func (e *Engine) Env() string {
	return e.appEnv
}

// Debug returns whether the application is in debug mode.
func (e *Engine) Debug() bool {
	return e.appDebug
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

// SetLogger sets the logger for the application engine.
func (e *Engine) SetLogger(l logger.Logger) {
	e.logger = l
}

// SetGracefull sets the lifecycle manager for the application engine.
func (e *Engine) SetGracefull(l lifecycle.Lifecycle) {
	e.gracefull = l
}

// SetContext sets the context for the application engine.
func (e *Engine) SetContext(ctx context.Context) {
	e.ctx = ctx
}
