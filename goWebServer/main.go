package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a new router
	mux := http.NewServeMux()

	// Add a simple route
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Processing request for /hello")
		fmt.Fprint(w, "Hello, World!")
	})

	// Add a route that reads query parameters
	mux.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "Guest"
		}
		log.Printf("Greeting %s", name)
		fmt.Fprintf(w, "Hello, %s!", name)
	})

	// Add a very slow endpoint
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		clientId := r.URL.Query().Get("id")
		if clientId == "" {
			clientId = "unknown"
		}
		log.Printf("Processing slow request from client %s", clientId)

		// Simulate work taking 10 seconds
		time.Sleep(10 * time.Second)

		fmt.Fprintf(w, "Slow response for client %s!", clientId)
		log.Printf("Completed slow request from client %s", clientId)
	})

	// Create the server and assign our router
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)
	// Notify the channel when we receive SIGINT (Ctrl+C) or SIGTERM (kill command)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Wait until we receive a signal
	sig := <-quit
	log.Printf("Received signal: %v. Starting shutdown...", sig)

	// Create a context with a timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt to shut down gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server has stopped gracefully")
}
