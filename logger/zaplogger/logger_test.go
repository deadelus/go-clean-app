package zaplogger_test

import (
	"testing"

	"github.com/deadelus/go-clean-app/v2/application"
	"github.com/deadelus/go-clean-app/v2/logger/zaplogger"
	"github.com/deadelus/go-clean-app/v2/logger/zaplogger/mock"
	"github.com/stretchr/testify/assert"
)

func TestSetZapLogger(t *testing.T) {
	// Use default NewZapLogger
	engine, err := application.New()
	assert.NoError(t, err)

	opt := zaplogger.SetZapLogger()
	assert.NotNil(t, opt)

	// Should not panic and should set logger
	assert.NotPanics(t, func() {
		opt(engine)
	})

	assert.NotNil(t, engine.Logger())
}

func TestSetZapLogger_PanicOnLoggerError(t *testing.T) {
	// Patch NewZapLogger to return error
	zaplogger.NewZapLogger = func(name, version, env string, debug bool) (*zaplogger.ZapLogger, zaplogger.Gracefull, error) {
		return nil, nil, assert.AnError
	}

	engine, err := application.New()
	assert.NoError(t, err)

	opt := zaplogger.SetZapLogger()
	assert.NotNil(t, opt)

	// Should panic due to logger creation error
	assert.Panics(t, func() {
		opt(engine)
	})
}

func TestSetZapLogger_PanicOnGracefulRegisterError(t *testing.T) {
	// Patch NewZapLogger to return valid logger and close function
	zaplogger.NewZapLogger = func(name, version, env string, debug bool) (*zaplogger.ZapLogger, zaplogger.Gracefull, error) {
		return &zaplogger.ZapLogger{}, func() error { return nil }, nil
	}

	engine, err := application.New()
	assert.NoError(t, err)

	engine.SetGracefull(&mock.Lifecycle{})

	opt := zaplogger.SetZapLogger()
	assert.NotNil(t, opt)

	// Should panic due to Register error
	assert.Panics(t, func() {
		opt(engine)
	})
}
