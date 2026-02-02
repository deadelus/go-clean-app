package local_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/deadelus/go-clean-app/v2/transport/adapter/local"
	"github.com/stretchr/testify/assert"
)

type dummyHandler struct{}

func (d *dummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func TestNewAdapter(t *testing.T) {
	h := &dummyHandler{}
	adapter := local.NewAdapter(h, 12345)
	assert.NotNil(t, adapter)
	assert.Equal(t, ":12345", adapter.Server.Addr)
	assert.Equal(t, h, adapter.Server.Handler)
}

func TestAdapter_Start_And_Stop(t *testing.T) {
	h := &dummyHandler{}
	adapter := local.NewAdapter(h, 12346)

	// Start server in a goroutine
	go func() {
		_ = adapter.Start()
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Stop the server
	err := adapter.Stop(context.Background())
	assert.NoError(t, err)
}
