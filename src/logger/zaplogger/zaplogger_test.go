package zaplogger_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/deadelus/go-clean-app/src/logger/zaplogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	t.Run("development mode", func(t *testing.T) {
		os.Setenv("LOGGER_MODE", "development")
		defer os.Unsetenv("LOGGER_MODE")

		logger, graceful, err := zaplogger.NewLogger("test-app", "v1.0.0", "LOGGER_MODE")
		require.NoError(t, err)
		require.NotNil(t, logger)
		require.NotNil(t, graceful)

		assert.NotNil(t, logger)
		graceful()
	})

	t.Run("production mode", func(t *testing.T) {
		os.Setenv("LOGGER_MODE", "production")
		defer os.Unsetenv("LOGGER_MODE")

		logger, graceful, err := zaplogger.NewLogger("test-app", "v1.0.0", "LOGGER_MODE")
		require.NoError(t, err)
		require.NotNil(t, logger)
		require.NotNil(t, graceful)

		assert.NotNil(t, logger)
		graceful()
	})

	t.Run("default mode", func(t *testing.T) {
		os.Unsetenv("LOGGER_MODE")

		logger, graceful, err := zaplogger.NewLogger("test-app", "v1.0.0", "LOGGER_MODE")
		require.NoError(t, err)
		require.NotNil(t, logger)
		require.NotNil(t, graceful)

		assert.NotNil(t, logger)
		graceful()
	})
}

func TestGetFromExternalLogger(t *testing.T) {
	zapLogger, _ := zap.NewProduction()
	logger, graceful, err := zaplogger.GetFromExternalLogger(zapLogger)
	require.NoError(t, err)
	require.NotNil(t, logger)
	require.NotNil(t, graceful)

	assert.Equal(t, zapLogger, logger.Logger)
	graceful()
}

func TestZapLogger_LoggingMethods(t *testing.T) {
	var buffer bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writer := zapcore.AddSync(&buffer)
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	zapLogger := zap.New(core)

	logger, _, _ := zaplogger.GetFromExternalLogger(zapLogger)

	tests := []struct {
		name    string
		logFunc func(msg string, fields ...any)
		level   zapcore.Level
	}{
		{"Info", logger.Info, zapcore.InfoLevel},
		{"Error", logger.Error, zapcore.ErrorLevel},
		{"Debug", logger.Debug, zapcore.DebugLevel},
		{"Warn", logger.Warn, zapcore.WarnLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer.Reset()
			tt.logFunc("test message", map[string]interface{}{"key": "value"}, zap.String("zkey", "zvalue"))

			var logOutput map[string]interface{}
			err := json.Unmarshal(buffer.Bytes(), &logOutput)
			require.NoError(t, err)

			assert.Equal(t, tt.level.String(), logOutput["level"])
			assert.Equal(t, "test message", logOutput["msg"])
			assert.Equal(t, "value", logOutput["key"])
			assert.Equal(t, "zvalue", logOutput["zkey"])
		})
	}
}

func TestConvertToZapFields(t *testing.T) {
	t.Run("with zap.Field", func(t *testing.T) {
		fields := []any{zap.String("key", "value")}
		zapFields := zaplogger.ConvertToZapFields(fields...)
		assert.Len(t, zapFields, 1)
		assert.Equal(t, zap.String("key", "value"), zapFields[0])
	})

	t.Run("with map[string]interface{}", func(t *testing.T) {
		fields := []any{map[string]interface{}{"key": "value", "num": 123}}
		zapFields := zaplogger.ConvertToZapFields(fields...)
		assert.Len(t, zapFields, 2)
		// Order is not guaranteed in maps
		assert.Contains(t, zapFields, zap.String("key", "value"))
		assert.Contains(t, zapFields, zap.Int("num", 123))
	})

	t.Run("with map[string]any", func(t *testing.T) {
		fields := []any{map[string]any{"key": "value", "bool": true}}
		zapFields := zaplogger.ConvertToZapFields(fields...)
		assert.Len(t, zapFields, 2)
		assert.Contains(t, zapFields, zap.String("key", "value"))
		assert.Contains(t, zapFields, zap.Bool("bool", true))
	})

	t.Run("with other type", func(t *testing.T) {
		fields := []any{"just a string"}
		zapFields := zaplogger.ConvertToZapFields(fields...)
		assert.Len(t, zapFields, 1)
		assert.Equal(t, zap.Any("field", "just a string"), zapFields[0])
	})

	t.Run("with mixed types", func(t *testing.T) {
		fields := []any{
			zap.String("zkey", "zvalue"),
			map[string]interface{}{"mkey": "mvalue"},
			"a string",
		}
		zapFields := zaplogger.ConvertToZapFields(fields...)
		assert.Len(t, zapFields, 3)
		assert.Contains(t, zapFields, zap.String("zkey", "zvalue"))
		assert.Contains(t, zapFields, zap.String("mkey", "mvalue"))
		assert.Contains(t, zapFields, zap.Any("field", "a string"))
	})
}

func TestConvertMapToZapFields(t *testing.T) {
	m := map[string]interface{}{
		"error":   assert.AnError,
		"string":  "hello",
		"int":     123,
		"int64":   int64(456),
		"float64": 78.9,
		"bool":    true,
		"any":     []string{"a", "b"},
	}

	fields := zaplogger.ConvertMapToZapFields(m)
	assert.Len(t, fields, 7)

	expectedFields := []zap.Field{
		zap.Error(assert.AnError),
		zap.String("string", "hello"),
		zap.Int("int", 123),
		zap.Int64("int64", 456),
		zap.Float64("float64", 78.9),
		zap.Bool("bool", true),
		zap.Any("any", []string{"a", "b"}),
	}

	for _, expected := range expectedFields {
		assert.Contains(t, fields, expected)
	}
}
