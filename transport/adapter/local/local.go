// Local adapter uses the standard net/http package to serve HTTP requests locally.
package local

import (
	"context"
	"fmt"
	"net/http"
)

// L'adapter local devient générique : il accepte n'importe quel http.Handler
type LocalAdapter struct {
	Server *http.Server
}

func NewAdapter(h http.Handler, port int) *LocalAdapter {
	return &LocalAdapter{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: h,
		},
	}
}

func (a *LocalAdapter) Start() error                   { return a.Server.ListenAndServe() }
func (a *LocalAdapter) Stop(ctx context.Context) error { return a.Server.Shutdown(ctx) }
