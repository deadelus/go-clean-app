package errors_test

import (
	"testing"

	"github.com/deadelus/go-clean-app/v2/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert.Equal(t, "missing configuration", errors.ErrMissingConfig)
	assert.Equal(t, "runtime error", errors.ErrRuntime)
}
