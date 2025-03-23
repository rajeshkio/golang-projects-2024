package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

type customLogger struct {
	http.ResponseWriter
	statusCode int
}

// captures the status code before writing it

func (c *customLogger) WriteHeader(code int) {
	c.statusCode = code
	c.ResponseWriter.WriteHeader(code)
}
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		customLogs := &customLogger{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Record start time
		start := time.Now()
		log.Printf("Request Received: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Printf(
			"REQUEST: method=%s path=%s ip=%s status=%d duration=%s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			customLogs.statusCode,
			duration,
		)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				http.Error(w, "Internal Server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func RespondWithError(w http.ResponseWriter, status int, message string, err error) {
	response := ErrorResponse{
		Status:  status,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to encode response", err)
	}
}
func greetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = " World"
	}
	response := map[string]string{"message": "Greeting, " + name + "!"}
	RespondWithJSON(w, http.StatusOK, response)
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("This is a deliberate panic!")
}
func main() {
	// Create a new router
	mux := http.NewServeMux()

	// Add a simple route
	// mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("Processing request for /hello")
	// 	fmt.Fprint(w, "Hello, World!")
	// })

	// Add a route that reads query parameters
	mux.Handle("/greet", Chain(
		http.HandlerFunc(greetHandler),
		LoggingMiddleware,
		RecoveryMiddleware,
		CORSMiddleware,
	))

	// Add a very slow endpoint
	// mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
	// 	clientId := r.URL.Query().Get("id")
	// 	if clientId == "" {
	// 		clientId = "unknown"
	// 	}
	// 	log.Printf("Processing slow request from client %s", clientId)

	// 	// Simulate work taking 10 seconds
	// 	time.Sleep(10 * time.Second)

	// 	fmt.Fprintf(w, "Slow response for client %s!", clientId)
	// 	log.Printf("Completed slow request from client %s", clientId)
	// })

	mux.Handle("/panic", Chain(
		http.HandlerFunc(panicHandler),
		LoggingMiddleware,
		RecoveryMiddleware,
		CORSMiddleware,
	))

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
