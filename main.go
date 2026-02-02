// Package main implements a simple application using go-clean-app
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/deadelus/go-clean-app/v2/application"
	"github.com/deadelus/go-clean-app/v2/transport/adapter/local"
)

func main() {
	// Initialise the application
	engine, err := application.New()
	if err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	// Example HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from local server!"))
	})

	localServer := local.NewAdapter(handler, 8080)

	// Register local server shutdown with the graceful shutdown manager
	engine.Gracefull().Register("local-server", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return localServer.Stop(ctx)
	})

	// Start local server on port 8080
	go func() {
		if err := localServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("local server error: %v", err)
		}
	}()
	log.Println("Local server started on :8080")

	// Wait for the stop signal via the context (e.g. SIGINT/SIGTERM)
	<-engine.Context().Done()
	// here we could perform other cleanup tasks if needed
	// before shutting down the application

	// Begin graceful shutdown
	// triggererd by context cancellation (e.g. SIGINT/SIGTERM)
	<-engine.Gracefull().Done()
}
