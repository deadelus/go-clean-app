package zaplogger_test

import (
	"testing"

	"github.com/deadelus/go-clean-app/v2/application"
	"github.com/deadelus/go-clean-app/v2/logger/zaplogger"
	"github.com/deadelus/go-clean-app/v2/logger/zaplogger/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetZapLoggerForCLI(t *testing.T) {
	// Replace NewZapLoggerForCLI with a custom config for testing
	zaplogger.NewZapLoggerForCLI = func() zap.Config {
		return zap.NewDevelopmentConfig()
	}

	engine, err := application.New()
	assert.NoError(t, err)

	opt := zaplogger.SetZapLoggerForCLI()
	assert.NotNil(t, opt)

	// Should not panic and should set logger
	assert.NotPanics(t, func() {
		opt(engine)
	})

	assert.NotNil(t, engine.Logger())
}

func TestSetZapLoggerForCLI_PanicOnBuildError(t *testing.T) {
	// Simulate zap.Config.Build error by returning a config with invalid output path
	zaplogger.NewZapLoggerForCLI = func() zap.Config {
		cfg := zap.NewDevelopmentConfig()
		// Set an invalid output path to force Build to fail
		cfg.OutputPaths = []string{"/dev/null/doesnotexist"}
		return cfg
	}

	engine, err := application.New()
	assert.NoError(t, err)

	opt := zaplogger.SetZapLoggerForCLI()
	assert.NotNil(t, opt)

	// Should panic due to logger build error
	assert.Panics(t, func() {
		opt(engine)
	})
}

func TestSetZapLoggerForCLI_PanicOnGracefulRegisterError(t *testing.T) {
	// Use a valid zap config
	zaplogger.NewZapLoggerForCLI = func() zap.Config {
		return zap.NewDevelopmentConfig()
	}

	engine, err := application.New()
	assert.NoError(t, err)

	// Patch Gracefull().Register to return error
	engine.SetGracefull(&mock.Lifecycle{})

	opt := zaplogger.SetZapLoggerForCLI()
	assert.NotNil(t, opt)

	// Should panic due to Register error
	assert.Panics(t, func() {
		opt(engine)
	})
}
