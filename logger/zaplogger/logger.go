package zaplogger

import (
	"fmt"

	"github.com/deadelus/go-clean-app/v2/application"
)

// SetLogger sets the logger for the Engine.
// NewZapLogger is a hook for zaplogger.NewLogger, can be replaced in tests.
var NewZapLogger = NewLogger

// SetZapLogger sets the logger for the Engine.
func SetZapLogger() application.Option {
	return func(e *application.Engine) {
		logger, closeLogger, err := NewZapLogger(
			e.Name(),
			e.Version(),
			e.Env(),
			e.Debug(),
		)

		if err != nil {
			panic(fmt.Errorf("failed to create zap logger: %w", err))
		}

		// Set the logger in the Engine
		e.SetLogger(logger)

		// Register the close function with the graceful shutdown manager
		if err := e.Gracefull().Register("zaplogger", closeLogger); err != nil {
			panic(fmt.Errorf("failed to register zap logger for graceful shutdown: %w", err))
		}
	}
}
