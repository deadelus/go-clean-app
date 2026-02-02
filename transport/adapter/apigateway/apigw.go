// Package apigateway provides an adapter to run an http.Handler in AWS Lambda with API Gateway.
package apigateway

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

// StartWithOptions is a hook for lambda.StartWithOptions, can be replaced in tests.
var StartWithOptions = lambda.StartWithOptions

// Adapter adapts any http.Handler for AWS Lambda + API Gateway.
type Adapter struct {
	Handler http.Handler
	Options []lambda.Option
}

// NewAdapter creates a new AWS Lambda API Gateway adapter.
func NewAdapter(h http.Handler, options ...lambda.Option) *Adapter {
	return &Adapter{
		Handler: h,
		Options: options,
	}
}

// Start starts the adapter by initializing the Lambda handler
// Start initializes the Lambda handler and blocks for API Gateway requests.
func (a *Adapter) Start() error {
	adapter := httpadapter.New(a.Handler)
	StartWithOptions(adapter.ProxyWithContext, a.Options...)
	return nil
}

// Stop satisfies the Transporter interface. For Lambda, AWS manages the lifecycle.
func (a *Adapter) Stop(ctx context.Context) error {
	return nil
}
