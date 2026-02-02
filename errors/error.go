// Package errors provides custom error types for the application.
package errors

// ErrMissingConfig is returned when a required config is missing.
const (
	ErrMissingConfig = "missing configuration"
	ErrRuntime       = "runtime error"
)
