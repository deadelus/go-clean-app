package main

import (
	"log"

	"github.com/deadelus/go-clean-app/v2/application"
)

func main() {
	// Initialise the application
	engine, err := application.New()
	if err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	// Here you start your servers, workers, etc. (e.g., go startServer())

	// Wait for the stop signal via the context
	<-engine.Context().Done()

	// Here you can perform actions before shutdown if necessary

	// Wait for the graceful shutdown to complete
	<-engine.Gracefull().Done()

	log.Println("Shutdown is over.")
}
