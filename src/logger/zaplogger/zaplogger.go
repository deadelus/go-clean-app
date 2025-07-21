// Package zaplogger provides a logger implementation using the zap logging library.
package zaplogger

import (
	"fmt"
	"os"
	"runtime"

	"go.uber.org/zap"
)

// ZapLogger is a logger implementation using the zap logging library.
// It implements the Logger interface defined in pkg/logger/logger.go.
type ZapLogger struct {
	logger *zap.Logger
}

type Gracefull func() error

// NewZapLogger creates a new ZapLogger instance.
// It initializes the zap logger and returns a ZapLogger instance.
// If there is an error during initialization, it returns the error.
func NewLogger(
	appName string,
	appVersion string,
	loggerModeEnvName string,
) (*ZapLogger, Gracefull, error) {

	modeStr := os.Getenv(loggerModeEnvName)
	var debugMode bool

	switch modeStr {
	case "development", "dev":
		fmt.Println("Logger mode set to development")
		debugMode = true
	case "production", "prod":
		fmt.Println("Logger mode set to production")
		debugMode = false
	default:
		fmt.Println("Logger mode not set or invalid, defaulting to development")
		debugMode = true
	}

	var config zap.Config
	var zapOptions []zap.Option

	zapOptions = append(zapOptions, zap.AddStacktrace(zap.PanicLevel))

	if debugMode {
		config = zap.NewDevelopmentConfig()
		zapOptions = append(zapOptions, zap.WithCaller(false))
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build(zapOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	logger = logger.Named(appName).With(
		zap.String("app_version", appVersion),
		zap.String("go_version", runtime.Version()))

	zl := &ZapLogger{logger: logger}

	gracefull := func() error {
		zl.Close()
		return nil
	}

	return zl, gracefull, nil
}

// GetFromExternalZapLogger sets the zap logger for the ZapLogger instance.
func GetFromExternalLogger(logger *zap.Logger) (*ZapLogger, Gracefull, error) {
	zl := &ZapLogger{logger: logger}

	gracefull := func() error {
		zl.Close()
		return nil
	}

	return zl, gracefull, nil
}

// Info logs an info message with the provided fields.
func (z *ZapLogger) Info(msg string, fields ...any) {
	z.logger.Info(msg, convertToZapFields(fields...)...)
}

// Error logs an error message with the provided fields.
func (z *ZapLogger) Error(msg string, fields ...any) {
	z.logger.Error(msg, convertToZapFields(fields...)...)
}

// Debug logs a debug message with the provided fields.
func (z *ZapLogger) Debug(msg string, fields ...any) {
	z.logger.Debug(msg, convertToZapFields(fields...)...)
}

// Warn logs a warning message with the provided fields.
func (z *ZapLogger) Warn(msg string, fields ...any) {
	z.logger.Warn(msg, convertToZapFields(fields...)...)
}

// Close flushes the logger and releases any resources.
// It ensures that all buffered log entries are written out.
// If there is an error during flushing, it logs the error using the zap logger.
// This method should be called when the application is shutting down to ensure proper cleanup.
func (z *ZapLogger) Close() {
	if err := z.logger.Sync(); err != nil {
		z.logger.Error("Failed to stop logger", zap.Error(err))
	}
}

// convertToZapFields converts various field types to zap.Field
func convertToZapFields(fields ...any) []zap.Field {
	var zapFields []zap.Field

	for _, field := range fields {
		// Si c'est déjà un zap.Field, on l'ajoute directement
		if f, ok := field.(zap.Field); ok {
			zapFields = append(zapFields, f)
			continue
		}

		// Si c'est une map[string]interface{} ou map[string]any
		if m, ok := field.(map[string]interface{}); ok {
			zapFields = append(zapFields, convertMapToZapFields(m)...)
			continue
		}

		if m, ok := field.(map[string]any); ok {
			converted := make(map[string]interface{})
			for k, v := range m {
				converted[k] = v
			}
			zapFields = append(zapFields, convertMapToZapFields(converted)...)
			continue
		}

		// Par défaut, on ajoute comme zap.Any
		zapFields = append(zapFields, zap.Any("field", field))
	}

	return zapFields
}

// convertMapToZapFields convertit une map en slice de zap.Field
func convertMapToZapFields(m map[string]interface{}) []zap.Field {
	var fields []zap.Field

	for key, value := range m {
		switch v := value.(type) {
		case error:
			fields = append(fields, zap.Error(v))
		case string:
			fields = append(fields, zap.String(key, v))
		case int:
			fields = append(fields, zap.Int(key, v))
		case int64:
			fields = append(fields, zap.Int64(key, v))
		case float64:
			fields = append(fields, zap.Float64(key, v))
		case bool:
			fields = append(fields, zap.Bool(key, v))
		default:
			fields = append(fields, zap.Any(key, v))
		}
	}

	return fields
}
