package lifecycle_test

import (
	"context"
	"errors"
	"testing"

	"github.com/deadelus/go-clean-app/v2/lifecycle"
	"github.com/stretchr/testify/assert"
)

func TestGracefull_Register(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g := lifecycle.NewGracefullShutdown(ctx)

	fn1 := func() error { return nil }
	err := g.Register("test1", fn1)
	assert.NoError(t, err)

	// Test already registered
	err = g.Register("test1", fn1)
	assert.NoError(t, err)
}

func TestGracefull_Shutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	g := lifecycle.NewGracefullShutdown(ctx)

	called := false
	g.Register("test", func() error {
		called = true
		return nil
	})

	cancel()

	<-g.Done() // wait for shutdown to complete
	assert.True(t, called)
}

func TestGracefull_Shutdown_Error(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	g := lifecycle.NewGracefullShutdown(ctx)

	errMock := errors.New("mock error")
	g.Register("test-error", func() error {
		return errMock
	})

	cancel()

	<-g.Done() // wait for shutdown to complete
}
