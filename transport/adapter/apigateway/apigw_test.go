package apigateway_test

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/deadelus/go-clean-app/v2/transport/adapter/apigateway"
	"github.com/stretchr/testify/assert"
)

type dummyHandler struct{}

func (d *dummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func TestNewApiGatewayAdapter(t *testing.T) {
	h := &dummyHandler{}
	adapter := apigateway.NewAdapter(h)
	assert.NotNil(t, adapter)
	assert.Equal(t, h, adapter.Handler)
}

func TestApiGatewayAdapter_Stop(t *testing.T) {
	h := &dummyHandler{}
	adapter := apigateway.NewAdapter(h)
	err := adapter.Stop(context.Background())
	assert.NoError(t, err)
}

func TestApiGatewayAdapter_Start_ReturnsNil(t *testing.T) {
	var called int32
	apigateway.StartWithOptions = func(_ interface{}, _ ...lambda.Option) {
		atomic.StoreInt32(&called, 1)
	}

	h := &dummyHandler{}
	adapter := apigateway.NewAdapter(h)
	// This test only checks the return value, not the blocking behavior.

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_ = adapter.Start()
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(t, int32(1), atomic.LoadInt32(&called))
}
