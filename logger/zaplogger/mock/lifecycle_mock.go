// Package mock provides mock implementations for zaplogger lifecycle.
package mock

import "github.com/stretchr/testify/assert"

// Lifecycle is a mock implementation of lifecycle.Lifecycle for tests.
type Lifecycle struct{}

// Register is a mock implementation of the Register method.
func (f *Lifecycle) Register(name string, fn func() error) error { return assert.AnError }

// Done is a mock implementation of the Done method.
func (f *Lifecycle) Done() <-chan struct{} { return make(chan struct{}) }
