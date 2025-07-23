package application_test

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/deadelus/go-clean-app/src/application"
	"github.com/deadelus/go-clean-app/src/logger/zaplogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	t.Run("successful creation with default app name", func(t *testing.T) {
		os.Unsetenv("TEST_APP_NAME")
		app, err := application.New("TEST_APP_NAME", application.SetOptionVersion("1.0.0"), application.SetZapLogger())
		require.NoError(t, err)
		assert.NotNil(t, app)
		assert.Equal(t, "application", app.Name())
		assert.Equal(t, "1.0.0", app.Version())
		assert.NotNil(t, app.Context())
		assert.NotNil(t, app.Gracefull())
		assert.NotNil(t, app.Logger())
	})

	t.Run("successful creation with app name from env", func(t *testing.T) {
		os.Setenv("TEST_APP_NAME", "my-test-app")
		defer os.Unsetenv("TEST_APP_NAME")
		app, err := application.New("TEST_APP_NAME", application.SetOptionVersion("1.0.0"), application.SetZapLogger())
		require.NoError(t, err)
		assert.Equal(t, "my-test-app", app.Name())
	})

	t.Run("successful creation with version from env", func(t *testing.T) {
		os.Setenv("TEST_APP_VERSION", "2.0.0")
		defer os.Unsetenv("TEST_APP_VERSION")
		app, err := application.New("TEST_APP_NAME", application.SetVersionFromSpecifiedEnv("TEST_APP_VERSION"), application.SetZapLogger())
		require.NoError(t, err)
		assert.Equal(t, "2.0.0", app.Version())
	})

	t.Run("successful creation with version from default env", func(t *testing.T) {
		os.Setenv(application.AppVersionEnvName, "3.0.0")
		defer os.Unsetenv(application.AppVersionEnvName)
		app, err := application.New("TEST_APP_NAME", application.SetVersionFromEnv(), application.SetZapLogger())
		require.NoError(t, err)
		assert.Equal(t, "3.0.0", app.Version())
	})

	t.Run("error when version env var is not set", func(t *testing.T) {
		os.Unsetenv("MISSING_VERSION_ENV")
		_, err := application.New("TEST_APP_NAME", application.SetVersionFromSpecifiedEnv("MISSING_VERSION_ENV"))
		assert.Error(t, err)
	})

	t.Run("with cli mode", func(t *testing.T) {
		app, err := application.New(
			"TEST_APP_NAME",
			application.SetOptionVersion("1.0.0"),
			application.WithCLIMode(),
			application.SetZapLoggerForCLI(),
		)
		require.NoError(t, err)
		assert.True(t, app.CLIMode())
	})
}

func TestEngine_Methods(t *testing.T) {
	app, err := application.New("TEST_APP_NAME", application.SetOptionVersion("1.0.0"), application.SetZapLogger())
	require.NoError(t, err)

	assert.Equal(t, "application", app.Name())
	assert.Equal(t, "1.0.0", app.Version())
	assert.NotNil(t, app.Context())
	assert.NotNil(t, app.Gracefull())
	assert.NotNil(t, app.Logger())
	assert.Equal(t, "default-user", app.CurrentUser())
	assert.Equal(t, "default-user-agent", app.UserAgent())
	assert.False(t, app.CLIMode())
}

func TestSignalHandling(t *testing.T) {
	os.Setenv("TEST_APP_NAME", "signal-test-app")
	defer os.Unsetenv("TEST_APP_NAME")

	app, err := application.New("TEST_APP_NAME", application.SetOptionVersion("1.0.0"), application.SetZapLogger())
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
	app, err := application.New("TEST_APP_NAME", application.SetOptionVersion("1.0.0"), application.SetZapLogger())
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

	t.Run("error in SetZapLogger", func(t *testing.T) {
		_, err := application.New(
			"TEST_APP_NAME",
			application.SetOptionVersion("1.0.0"),
			func(e *application.Engine) error {
				e.SetGracefull(&mockLifecycle{err: mockErr})
				return nil
			},
			application.SetZapLogger(),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register zap logger for graceful shutdown")
	})

	t.Run("error in SetZapLoggerForCLI", func(t *testing.T) {
		_, err := application.New(
			"TEST_APP_NAME",
			application.SetOptionVersion("1.0.0"),
			application.WithCLIMode(),
			func(e *application.Engine) error {
				e.SetGracefull(&mockLifecycle{err: mockErr})
				return nil
			},
			application.SetZapLoggerForCLI(),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register zap logger for CLI graceful shutdown")
	})
}

func TestLoggerCreationErrors(t *testing.T) {
	mockErr := errors.New("mock logger creation error")

	t.Run("error in NewLogger", func(t *testing.T) {
		originalNewLogger := application.NewZapLogger
		application.NewZapLogger = func(string, string, string) (*zaplogger.ZapLogger, zaplogger.Gracefull, error) {
			return nil, nil, mockErr
		}
		defer func() { application.NewZapLogger = originalNewLogger }()

		_, err := application.New(
			"TEST_APP_NAME",
			application.SetOptionVersion("1.0.0"),
			application.SetZapLogger(),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create zap logger")
	})

	t.Run("error in NewLoggerForCLI", func(t *testing.T) {
		originalNewZapLoggerForCLI := application.NewZapLoggerForCLI
		application.NewZapLoggerForCLI = func() zap.Config {
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
		defer func() { application.NewZapLoggerForCLI = originalNewZapLoggerForCLI }()

		_, err := application.New(
			"TEST_APP_NAME",
			application.SetOptionVersion("1.0.0"),
			application.WithCLIMode(),
			application.SetZapLoggerForCLI(),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create zap logger for CLI")
	})
}
