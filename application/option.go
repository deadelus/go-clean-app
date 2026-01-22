// Package application provides the application context and lifecycle management
package application

import (
	"context"
)

// Option is a function that configures the Engine.
type Option func(*Engine)

// WithCLIMode is an option to set the application to run in CLI mode.
func WithCLIMode() Option {
	return func(e *Engine) {
		e.ctx = context.WithValue(e.ctx, "cli_mode", true)
	}
}

// Version is an Option that sets the application version in the Engine.
// It allows the application version to be configured at runtime.
// This option can be used to set the version of the application when creating a new Engine instance.
// It is useful for applications that need to report their version or for logging purposes.
func Version(version string) Option {
	return func(e *Engine) {
		e.appVersion = version
	}
}

// AppName is an Option that sets the application name in the Engine.
// It allows the application name to be configured at runtime.
// This option can be used to set the name of the application when creating a new Engine instance.
// It is useful for applications that need to report their name or for logging purposes.
func AppName(name string) Option {
	return func(e *Engine) {
		e.appName = name
	}
}

// Debug is an Option that sets the debug mode of the application in the Engine.
// It allows the debug mode to be configured at runtime.
// This option can be used to enable or disable debug mode when creating a new Engine instance.
// It is useful for applications that need to run in different modes for development and production.
func Debug(debug bool) Option {
	return func(e *Engine) {
		e.appDebug = debug
	}
}

// Env is an Option that sets the application environment in the Engine.
// It allows the application environment to be configured at runtime.
// This option can be used to set the environment (e.g., development, production) when creating a new Engine instance.
// It is useful for applications that need to behave differently based on the environment they are running in.
func Env(env string) Option {
	return func(e *Engine) {
		e.appEnv = env
	}
}
