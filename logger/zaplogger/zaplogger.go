// Package zaplogger provides a logger implementation using the zap logging library.
package zaplogger

import (
	"fmt"
	"runtime"

	"go.uber.org/zap"
)

// ZapLogger is a logger implementation using the zap logging library.
// It implements the Logger interface defined in pkg/logger/logger.go.
type ZapLogger struct {
	Logger *zap.Logger
}

type Gracefull func() error

// BuildConfig is a helper to allow testing config.Build() errors
var BuildConfig = func(appDebug bool) zap.Config {
	if appDebug {
		return zap.NewDevelopmentConfig()
	}
	return zap.NewProductionConfig()
}

// NewZapLogger creates a new ZapLogger instance.
// It initializes the zap logger and returns a ZapLogger instance.
// If there is an error during initialization, it returns the error.
func NewLogger(
	appName string,
	appVersion string,
	appEnv string,
	appDebug bool,
) (*ZapLogger, Gracefull, error) {

	switch appEnv {
	case "development", "dev":
		fmt.Println("Logger mode set to development")
	case "production", "prod":
		fmt.Println("Logger mode set to production")
	default:
		fmt.Println("Logger mode not set or invalid, defaulting to development")
		appEnv = "development"
	}

	config := BuildConfig(appDebug)
	var zapOptions []zap.Option

	zapOptions = append(zapOptions, zap.AddStacktrace(zap.PanicLevel))

	if appDebug {
		zapOptions = append(zapOptions, zap.WithCaller(false))
	}

	logger, err := config.Build(zapOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create zap Logger: %w", err)
	}

	logger = logger.Named(appName).With(
		zap.String("app_version", appVersion),
		zap.String("app_env", appEnv),
		zap.Bool("app_debug", appDebug),
		zap.String("go_version", runtime.Version()))

	zl := &ZapLogger{Logger: logger}

	gracefull := func() error {
		zl.Close()
		return nil
	}

	return zl, gracefull, nil
}

// GetFromExternalZapLogger sets the zap logger for the ZapLogger instance.
func GetFromExternalLogger(logger *zap.Logger) (*ZapLogger, Gracefull, error) {
	zl := &ZapLogger{Logger: logger}

	gracefull := func() error {
		zl.Close()
		return nil
	}

	return zl, gracefull, nil
}

// Info logs an info message with the provided fields.
func (z *ZapLogger) Info(msg string, fields ...any) {
	z.Logger.Info(msg, ConvertToZapFields(fields...)...)
}

// Error logs an error message with the provided fields.
func (z *ZapLogger) Error(msg string, fields ...any) {
	z.Logger.Error(msg, ConvertToZapFields(fields...)...)
}

// Debug logs a debug message with the provided fields.
func (z *ZapLogger) Debug(msg string, fields ...any) {
	z.Logger.Debug(msg, ConvertToZapFields(fields...)...)
}

// Warn logs a warning message with the provided fields.
func (z *ZapLogger) Warn(msg string, fields ...any) {
	z.Logger.Warn(msg, ConvertToZapFields(fields...)...)
}

// Close flushes the logger and releases any resources.
// It ensures that all buffered log entries are written out.
// If there is an error during flushing, it logs the error using the zap logger.
// This method should be called when the application is shutting down to ensure proper cleanup.
func (z *ZapLogger) Close() {
	z.Logger.Sync()
}

// ConvertToZapFields converts various field types to zap.Field
func ConvertToZapFields(fields ...any) []zap.Field {
	var zapFields []zap.Field

	for _, field := range fields {
		// Si c'est déjà un zap.Field, on l'ajoute directement
		if f, ok := field.(zap.Field); ok {
			zapFields = append(zapFields, f)
			continue
		}

		// Si c'est une map[string]any
		if m, ok := field.(map[string]any); ok {
			zapFields = append(zapFields, ConvertMapToZapFields(m)...)
			continue
		}

		// Par défaut, on ajoute comme zap.Any
		zapFields = append(zapFields, zap.Any("field", field))
	}

	return zapFields
}

// ConvertMapToZapFields convertit une map en slice de zap.Field
func ConvertMapToZapFields(m map[string]interface{}) []zap.Field {
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
