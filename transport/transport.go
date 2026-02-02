// Package transport provides interfaces and implementations for various transport mechanisms.
package transport

import (
	"context"
)

// Transporter defines the interface for transport mechanisms
type Transporter interface {
	Start() error
	Stop(ctx context.Context) error
}
