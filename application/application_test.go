package application_test

import (
	"context"
	"errors"
	"syscall"
	"testing"
	"time"

	"github.com/deadelus/go-clean-app/v2/application"
	"github.com/deadelus/go-clean-app/v2/logger/zaplogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	t.Run("successful creation with option", func(t *testing.T) {
		app, err := application.New(
			application.AppName("my-test-app"),
			application.Env("dev"),
			application.Version("2.0.0"),
			application.Debug(false),
			zaplogger.SetZapLogger(),
		)
		require.NoError(t, err)
		assert.Equal(t, "my-test-app", app.Name())
		assert.Equal(t, "dev", app.Env())
		assert.Equal(t, "2.0.0", app.Version())
		assert.False(t, app.Debug())
	})

	t.Run("successful creation with default", func(t *testing.T) {
		app, err := application.New(zaplogger.SetZapLogger())
		require.NoError(t, err)
		assert.Equal(t, "application", app.Name())
		assert.Equal(t, "development", app.Env())
		assert.Equal(t, "0.1.0", app.Version())
		assert.False(t, app.Debug())
	})

	t.Run("with cli mode", func(t *testing.T) {
		app, err := application.New(
			application.AppName("TEST"),
			application.Version("1.0.0"),
			application.WithCLIMode(),
			zaplogger.SetZapLoggerForCLI(),
		)
		require.NoError(t, err)
		assert.True(t, app.CLIMode())
	})
}

func TestEngine_Methods(t *testing.T) {
	app, err := application.New(application.AppName("TEST"), application.Version("1.0.0"), zaplogger.SetZapLogger())
	require.NoError(t, err)

	assert.Equal(t, "TEST", app.Name())
	assert.Equal(t, "1.0.0", app.Version())
	assert.Equal(t, "development", app.Env())
	assert.False(t, app.Debug())
	assert.NotNil(t, app.Context())
	assert.NotNil(t, app.Gracefull())
	assert.NotNil(t, app.Logger())
	assert.Equal(t, "default-user", app.CurrentUser())
	assert.Equal(t, "default-user-agent", app.UserAgent())
	assert.False(t, app.CLIMode())
}

func TestSignalHandling(t *testing.T) {
	app, err := application.New(application.AppName("TEST"), application.Version("1.0.0"), zaplogger.SetZapLogger())
	require.NoError(t, err)

	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()

	select {
	case <-app.Context().Done():
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for context to be canceled")
	}
}

type mockLifecycle struct {
	err error
}

func TestEngine_SetContext(t *testing.T) {
	app, err := application.New(zaplogger.SetZapLogger())
	require.NoError(t, err)

	originalCtx := app.Context()
	newCtx := context.WithValue(context.Background(), "test_key", "test_value")

	app.SetContext(newCtx)

	assert.NotEqual(t, originalCtx, app.Context())
	assert.Equal(t, "test_value", app.Context().Value("test_key"))
}
func (m *mockLifecycle) Register(name string, fn func() error) error {
	return m.err
}

func (m *mockLifecycle) Shutdown() {}

func TestLoggerRegistrationErrors(t *testing.T) {
	mockErr := errors.New("mock error")

	t.Run("error in SetZapLogger should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			application.New(
				func(e *application.Engine) {
					e.SetGracefull(&mockLifecycle{err: mockErr})
				},
				zaplogger.SetZapLogger(),
			)
		})
	})

	t.Run("error in SetZapLoggerForCLI should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			application.New(
				application.WithCLIMode(),
				func(e *application.Engine) {
					e.SetGracefull(&mockLifecycle{err: mockErr})
				},
				zaplogger.SetZapLoggerForCLI(),
			)
		})
	})
}

func TestLoggerCreationErrors(t *testing.T) {
	mockErr := errors.New("mock logger creation error")

	t.Run("error in NewLogger should panic", func(t *testing.T) {
		originalNewLogger := zaplogger.NewZapLogger
		zaplogger.NewZapLogger = func(string, string, string, bool) (*zaplogger.ZapLogger, zaplogger.Gracefull, error) {
			return nil, nil, mockErr
		}
		defer func() { zaplogger.NewZapLogger = originalNewLogger }()

		assert.Panics(t, func() {
			application.New(
				zaplogger.SetZapLogger(),
			)
		})
	})

	t.Run("error in NewLoggerForCLI should panic", func(t *testing.T) {
		originalNewZapLoggerForCLI := zaplogger.NewZapLoggerForCLI
		zaplogger.NewZapLoggerForCLI = func() zap.Config {
			// Return a config that will cause Build to fail
			return zap.Config{
				Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
				Development: false,
				Sampling: &zap.SamplingConfig{
					Initial:    100,
					Thereafter: 100,
				},
				Encoding:         "json",
				EncoderConfig:    zap.NewProductionEncoderConfig(),
				OutputPaths:      []string{"/invalid/path"},
				ErrorOutputPaths: []string{"stderr"},
			}
		}
		defer func() { zaplogger.NewZapLoggerForCLI = originalNewZapLoggerForCLI }()

		assert.Panics(t, func() {
			application.New(
				application.WithCLIMode(),
				zaplogger.SetZapLoggerForCLI(),
			)
		})
	})
}
