package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "Hello, World!"); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "OK"); err != nil {
		log.Printf("Error writing health response: %v", err)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Replace Version with git tag / git release
	w.Header().Set("Content-Type", "application/json")
	fv := map[string]interface{}{
		"Repository":           "RepoOfWebServers",
		"URL":                  "https://github.com/Mattible/RepoOfWebServers",
		"Programming Language": "Go",
		"Version":              "0.1.0",
		"git Sha":              os.Getenv("GITSHA"),
		// "git_tag":              os.Getenv("GITTAG"),

		"routes": []map[string]string{
			{"path": "/", "description": "Hello World"},
			{"path": "/health", "description": "Health check"},
			{"path": "/info", "description": "Server info"},
			{"path": "/image", "description": "Image handler"},
		},
	}
	if err := json.NewEncoder(w).Encode(fv); err != nil {
		http.Error(w, "Failed to encode info", http.StatusInternalServerError)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO Add image path to a cloud provided CDN Cache or local Storage
	if _, err := fmt.Fprintf(w, ""); err != nil {
		log.Printf("Error writing image response: %v", err)
	}
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			loggingCall(handler)(w, r)
		case "/health":
			loggingCall(healthHandler)(w, r)
		case "/info":
			loggingCall(infoHandler)(w, r)
		case "/image":
			loggingCall(imageHandler)(w, r)
		default:
			loggingCall(notFoundHandler)(w, r)
		}
	})
}

func startServer(server *http.Server) error {
	log.Printf("Server is listening on port %s", server.Addr)
	log.Printf("Press Ctrl+C to shutdown server...")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %v", err)
	}
	return nil
}

func gracefulShutdown(server *http.Server) {
	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown gracefully")
	}
}

func loggingCall(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: Method=%s, URL=%s, RemoteAddr=%s", r.Method, r.URL, r.RemoteAddr)
		next(w, r)
	}
}

func main() {
	port := os.Getenv("WEBSERVER_PORT")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Addr: ":" + port,
	}

	// Setup routes
	setupRoutes()

	// Start server in a goroutine
	go func() {
		if err := startServer(server); err != nil {
			log.Fatal(err)
		}
	}()

	// Handle graceful shutdown in main goroutine
	gracefulShutdown(server)
}
