// Package lifecycle provides a mechanism for managing application lifecycle events, particularly for graceful shutdowns.
package lifecycle

import (
	"context"
	"log"
	"sync"
)

// Lifecycle interface defines methods for managing application lifecycle events.
type Lifecycle interface {
	Done() <-chan struct{}
	Register(name string, gracefull func() error) error
}

// Gracefull represents a list of functions to be executed during graceful shutdown.
type Gracefull struct {
	functions map[string]func() error
	done      chan struct{}
}

// Done returns a channel that is closed when the graceful shutdown is complete.
func (g *Gracefull) Done() <-chan struct{} {
	return g.done
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
func (g *Gracefull) gracefullAll() {
	log.Println("Shutting down in progress...")

	wg := &sync.WaitGroup{}
	for name, gracefullFunc := range g.functions {
		wg.Add(1)
		k, v := name, gracefullFunc
		go g.gracefullOne(wg, k, v)
	}
	wg.Wait()

	log.Println("Shutdown is over.")

	g.done <- struct{}{}
}

// gracefullOne executes a single registered function and logs any errors.
func (g *Gracefull) gracefullOne(wg *sync.WaitGroup, name string, gracefullFunc func() error) {
	defer wg.Done()

	if err := gracefullFunc(); err != nil {
		log.Printf("Error during gracefull shutdown of %s: %v", name, err)

		return
	}

	log.Printf("Gracefull shutdown of %s completed successfully", name)
}
