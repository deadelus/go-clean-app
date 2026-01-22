package zaplogger

import (
	"fmt"

	"github.com/deadelus/go-clean-app/v2/application"
	"go.uber.org/zap"
)

// NewZapLoggerForCLI is a hook for zap.NewDevelopmentConfig, can be replaced in tests.
var NewZapLoggerForCLI = zap.NewDevelopmentConfig

// SetZapLoggerForCLI sets the logger for the Engine specifically for CLI applications.
func SetZapLoggerForCLI() application.Option {
	return func(e *application.Engine) {
		config := NewZapLoggerForCLI()
		l, err := config.Build(
			zap.AddStacktrace(zap.PanicLevel),
			zap.WithCaller(false),
		)

		if err != nil {
			panic(fmt.Errorf("failed to create zap logger for CLI: %w", err))
		}

		logger, closeLogger, _ := GetFromExternalLogger(l)

		// Set the logger in the Engine
		e.SetLogger(logger)

		// Register the close function with the graceful shutdown manager
		if err := e.Gracefull().Register("zaplogger-cli", closeLogger); err != nil {
			panic(fmt.Errorf("failed to register zap logger for CLI graceful shutdown: %w", err))
		}
	}
}
