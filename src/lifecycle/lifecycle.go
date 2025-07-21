// Package lifecycle provides a mechanism for managing application lifecycle events, particularly for graceful shutdowns.
package lifecycle

import (
	"context"
	"log"
	"sync"
)

// Lifecycle interface defines methods for managing application lifecycle events.
type Lifecycle interface {
	Register(name string, gracefull func() error) error
}

// Gracefull represents a list of functions to be executed during graceful shutdown.
type Gracefull struct {
	functions map[string]func() error
	done      chan struct{}
}

// NewGracefullShutdown is the constructor of the shutdown ochestrator.
func NewGracefullShutdown(ctx context.Context) *Gracefull {
	life := &Gracefull{
		functions: make(map[string]func() error),
		done:      make(chan struct{}),
	}

	go func() {
		<-ctx.Done()
		life.gracefullAll()
	}()

	return life
}

// Register adds a function to the list of functions to be executed during graceful shutdown.
func (g *Gracefull) Register(name string, gracefull func() error) error {
	if _, exists := g.functions[name]; exists {
		return nil // Already registered
	}
	g.functions[name] = gracefull
	return nil
}

// gracefullAll executes all registered functions in the order they were added.
func (life *Gracefull) gracefullAll() {
	log.Println("Shutting down in progress...")

	wg := &sync.WaitGroup{}
	for name, gracefullFunc := range life.functions {
		wg.Add(1)
		k, v := name, gracefullFunc
		go life.gracefullOne(wg, k, v)
	}
	wg.Wait()

	log.Println("Shutdown is over.")

	life.done <- struct{}{}
}

// gracefullOne executes a single registered function and logs any errors.
func (life *Gracefull) gracefullOne(wg *sync.WaitGroup, name string, gracefullFunc func() error) {
	defer wg.Done()

	if err := gracefullFunc(); err != nil {
		log.Printf("Error during gracefull shutdown of %s: %v", name, err)

		return
	}

	log.Printf("Gracefull shutdown of %s completed successfully", name)
}
