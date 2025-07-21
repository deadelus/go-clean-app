package application

import (
	"fmt"

	"github.com/deadelus/go-clean-app/src/logger/zaplogger"
	"go.uber.org/zap"
)

// SetZapLoggerForCLI sets the logger for the Engine specifically for CLI applications.
func setZapLoggerForCLI() Option {
	return func(e *Engine) error {
		config := zap.NewDevelopmentConfig()
		l, err := config.Build(
			zap.AddStacktrace(zap.PanicLevel),
			zap.WithCaller(false),
		)

		if err != nil {
			return fmt.Errorf("failed to create zap logger for CLI: %w", err)
		}

		logger, closeLogger, _ := zaplogger.GetFromExternalLogger(l)

		// Set the logger in the Engine
		e.logger = logger

		// Register the close function with the graceful shutdown manager
		if err := e.gracefull.Register("zaplogger-cli", closeLogger); err != nil {
			return fmt.Errorf("failed to register zap logger for CLI graceful shutdown: %w", err)
		}

		return nil
	}
}
