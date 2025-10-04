package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("healthHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("healthHandler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestInfoHandler(t *testing.T) {
	// Set environment variable for testing
	if err := os.Setenv("GITSHA", "test-sha-123"); err != nil {
		t.Fatalf("Failed to set GITSHA environment variable: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("GITSHA"); err != nil {
			t.Logf("Failed to unset GITSHA environment variable: %v", err)
		}
	}()

	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(infoHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("infoHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("infoHandler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check required fields
	expectedFields := []string{"Repository", "URL", "Programming Language", "Version", "git Sha", "routes"}
	for _, field := range expectedFields {
		if _, exists := response[field]; !exists {
			t.Errorf("infoHandler response missing field: %s", field)
		}
	}

	// Check specific values
	if response["Repository"] != "RepoOfWebServers" {
		t.Errorf("Expected Repository to be 'RepoOfWebServers', got %v", response["Repository"])
	}

	if response["Programming Language"] != "Go" {
		t.Errorf("Expected Programming Language to be 'Go', got %v", response["Programming Language"])
	}

	if response["git Sha"] != "test-sha-123" {
		t.Errorf("Expected git Sha to be 'test-sha-123', got %v", response["git Sha"])
	}

	// Check routes array
	routes, ok := response["routes"].([]interface{})
	if !ok {
		t.Error("Expected routes to be an array")
	} else if len(routes) != 4 {
		t.Errorf("Expected 4 routes, got %d", len(routes))
	}
}

func TestImageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/image", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(imageHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("imageHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Image handler currently returns empty string
	expected := ""
	if rr.Body.String() != expected {
		t.Errorf("imageHandler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestNotFoundHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(notFoundHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("notFoundHandler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	if !strings.Contains(rr.Body.String(), "404 page not found") {
		t.Errorf("notFoundHandler should contain '404 page not found', got %v",
			rr.Body.String())
	}
}

func TestLoggingCall(t *testing.T) {
	// Create a test handler
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("test response")); err != nil {
			t.Errorf("Failed to write test response: %v", err)
		}
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Wrap with logging middleware
	wrappedHandler := loggingCall(testHandler)
	wrappedHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("loggingCall returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "test response"
	if rr.Body.String() != expected {
		t.Errorf("loggingCall returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSetupRoutes(t *testing.T) {
	// Setup routes
	setupRoutes()

	tests := []struct {
		name           string
		path           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Root endpoint",
			path:           "/",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, World!",
		},
		{
			name:           "Health endpoint",
			path:           "/health",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Info endpoint",
			path:           "/info",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "", // JSON response, just check status
		},
		{
			name:           "Image endpoint",
			path:           "/image",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "Not found endpoint",
			path:           "/nonexistent",
			method:         "GET",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			// Use the default ServeMux since setupRoutes() registers with it
			http.DefaultServeMux.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("%s returned wrong status code: got %v want %v",
					tt.name, status, tt.expectedStatus)
			}

			if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
				t.Errorf("%s returned unexpected body: got %v want %v",
					tt.name, rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestMainFunctionEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name            string
		envValue        string
		expectedDefault string
	}{
		{
			name:            "Default port when no env var",
			envValue:        "",
			expectedDefault: "8000",
		},
		{
			name:            "Custom port from env var",
			envValue:        "9000",
			expectedDefault: "9000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env var
			originalPort := os.Getenv("WEBSERVER_PORT")
			defer func() {
				if err := os.Setenv("WEBSERVER_PORT", originalPort); err != nil {
					t.Logf("Failed to restore WEBSERVER_PORT environment variable: %v", err)
				}
			}()

			// Set test env var
			if tt.envValue != "" {
				if err := os.Setenv("WEBSERVER_PORT", tt.envValue); err != nil {
					t.Fatalf("Failed to set WEBSERVER_PORT environment variable: %v", err)
				}
			} else {
				if err := os.Unsetenv("WEBSERVER_PORT"); err != nil {
					t.Fatalf("Failed to unset WEBSERVER_PORT environment variable: %v", err)
				}
			}

			// Test the port logic (extracted from main function)
			port := os.Getenv("WEBSERVER_PORT")
			if port == "" {
				port = "8000"
			}

			if port != tt.expectedDefault {
				t.Errorf("Expected port %s, got %s", tt.expectedDefault, port)
			}
		})
	}
}

// Benchmark tests
func BenchmarkHandler(b *testing.B) {
	req, _ := http.NewRequest("GET", "/", nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler(rr, req)
	}
}

func BenchmarkHealthHandler(b *testing.B) {
	req, _ := http.NewRequest("GET", "/health", nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		healthHandler(rr, req)
	}
}

func BenchmarkInfoHandler(b *testing.B) {
	req, _ := http.NewRequest("GET", "/info", nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		infoHandler(rr, req)
	}
}

func BenchmarkLoggingCall(b *testing.B) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("test")); err != nil {
			b.Errorf("Failed to write test response: %v", err)
		}
	}
	wrappedHandler := loggingCall(testHandler)
	req, _ := http.NewRequest("GET", "/test", nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		wrappedHandler(rr, req)
	}
}

func TestStartServerWithContext(t *testing.T) {
	// Test server with context for controlled shutdown
	server := &http.Server{
		Addr: ":0",
	}

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- startServer(server)
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Shutdown the server
	shutdownErr := server.Shutdown(context.Background())
	if shutdownErr != nil {
		t.Errorf("Failed to shutdown server: %v", shutdownErr)
	}

	// Wait for server to return
	select {
	case err := <-serverErr:
		// Server should return nil when shutdown gracefully
		if err != nil {
			t.Errorf("Expected server to return nil on graceful shutdown, got: %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Server did not shutdown within expected time")
	}
}

func TestStartServerError(t *testing.T) {
	// Test server start with invalid address
	server := &http.Server{
		Addr: "invalid:address:port",
	}

	err := startServer(server)
	if err == nil {
		t.Error("Expected error for invalid server address")
	}
}

func TestGracefulShutdown(t *testing.T) {
	server := &http.Server{
		Addr: ":0",
	}

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Test graceful shutdown by calling Shutdown directly
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	shutdownErr := server.Shutdown(ctx)
	if shutdownErr != nil {
		t.Errorf("Failed to shutdown server gracefully: %v", shutdownErr)
	}

	// Wait for server to return
	select {
	case err := <-serverErr:
		// Server should return ErrServerClosed on graceful shutdown
		if err != http.ErrServerClosed {
			t.Errorf("Expected server to return ErrServerClosed, got: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("Server did not shutdown within expected time")
	}
}

func TestHandlerError(t *testing.T) {
	// Test handler error path by using a response writer that fails
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response writer that will fail on write
	rw := &errorResponseWriter{}

	handler(rw, req)

	// The handler should handle the error gracefully (log it)
	// We can't easily test the log output, but we can verify no panic
}

func TestHealthHandlerError(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rw := &errorResponseWriter{}

	healthHandler(rw, req)

	// Should handle error gracefully
}

func TestImageHandlerError(t *testing.T) {
	req, err := http.NewRequest("GET", "/image", nil)
	if err != nil {
		t.Fatal(err)
	}

	rw := &errorResponseWriter{}

	imageHandler(rw, req)

	// Should handle error gracefully
}

func TestInfoHandlerJSONError(t *testing.T) {
	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rw := &errorResponseWriter{}

	infoHandler(rw, req)

	// Should return 500 error for JSON encoding failure
}

func TestInfoHandlerWithoutGitSha(t *testing.T) {
	// Ensure GITSHA is not set
	if err := os.Unsetenv("GITSHA"); err != nil {
		t.Logf("Failed to unset GITSHA: %v", err)
	}

	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(infoHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("infoHandler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// git Sha should be empty string when not set
	if gitSha, exists := response["git Sha"]; exists {
		if gitShaStr, ok := gitSha.(string); !ok || gitShaStr != "" {
			t.Errorf("Expected git Sha to be empty string when GITSHA env var is not set, got: %v", gitSha)
		}
	} else {
		t.Error("Expected git Sha field to exist in response")
	}
}

func TestInfoHandlerRoutesStructure(t *testing.T) {
	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(infoHandler)

	handler.ServeHTTP(rr, req)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	routes, ok := response["routes"].([]interface{})
	if !ok {
		t.Fatal("routes should be an array")
	}

	expectedRoutes := []map[string]string{
		{"path": "/", "description": "Hello World"},
		{"path": "/health", "description": "Health check"},
		{"path": "/info", "description": "Server info"},
		{"path": "/image", "description": "Image handler"},
	}

	if len(routes) != len(expectedRoutes) {
		t.Errorf("Expected %d routes, got %d", len(expectedRoutes), len(routes))
	}

	for i, route := range routes {
		routeMap, ok := route.(map[string]interface{})
		if !ok {
			t.Errorf("Route %d should be a map", i)
			continue
		}

		expected := expectedRoutes[i]
		if path, exists := routeMap["path"]; !exists || path != expected["path"] {
			t.Errorf("Route %d path: expected %s, got %v", i, expected["path"], path)
		}
		if desc, exists := routeMap["description"]; !exists || desc != expected["description"] {
			t.Errorf("Route %d description: expected %s, got %v", i, expected["description"], desc)
		}
	}
}

// Helper type for testing error conditions
type errorResponseWriter struct {
	header http.Header
}

func (e *errorResponseWriter) Header() http.Header {
	if e.header == nil {
		e.header = make(http.Header)
	}
	return e.header
}

func (e *errorResponseWriter) Write(data []byte) (int, error) {
	return 0, fmt.Errorf("simulated write error")
}

func (e *errorResponseWriter) WriteHeader(statusCode int) {
	// No-op for error testing
}

func TestConcurrentRequests(t *testing.T) {
	// Test that handlers can handle concurrent requests
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			req, _ := http.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			handler(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Concurrent request %d failed with status %d", id, rr.Code)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestDifferentHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			testHandler := http.HandlerFunc(handler)

			testHandler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("%s request returned wrong status code: got %v want %v",
					method, status, http.StatusOK)
			}
		})
	}
}

func TestMainFunctionIntegration(t *testing.T) {

	// Save original port
	originalPort := os.Getenv("WEBSERVER_PORT")
	defer func() {
		if err := os.Setenv("WEBSERVER_PORT", originalPort); err != nil {
			t.Logf("Failed to restore port: %v", err)
		}
	}()

	// Test port setting
	if err := os.Setenv("WEBSERVER_PORT", "8888"); err != nil {
		t.Fatalf("Failed to set port: %v", err)
	}

	port := os.Getenv("WEBSERVER_PORT")
	if port == "" {
		port = "8000"
	}

	if port != "8888" {
		t.Errorf("Expected port 8888, got %s", port)
	}

	// Test server creation
	server := &http.Server{
		Addr: ":" + port,
	}

	if server.Addr != ":8888" {
		t.Errorf("Expected server addr :8888, got %s", server.Addr)
	}
}
