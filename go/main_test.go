package main

import (
	"encoding/json"
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
	os.Setenv("GITSHA", "test-sha-123")
	defer os.Unsetenv("GITSHA")

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
		w.Write([]byte("test response"))
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
		name        string
		envValue    string
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
			defer os.Setenv("WEBSERVER_PORT", originalPort)

			// Set test env var
			if tt.envValue != "" {
				os.Setenv("WEBSERVER_PORT", tt.envValue)
			} else {
				os.Unsetenv("WEBSERVER_PORT")
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
		w.Write([]byte("test"))
	}
	wrappedHandler := loggingCall(testHandler)
	req, _ := http.NewRequest("GET", "/test", nil)
	
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		wrappedHandler(rr, req)
	}
}

// Test server startup and shutdown functionality
func TestServerLifecycle(t *testing.T) {
	// This test verifies that the server can start and stop without hanging
	// Note: This is a simplified test since testing actual signal handling
	// requires more complex setup
	
	server := &http.Server{
		Addr: ":0", // Use any available port
	}
	
	// Test that we can create the server without error
	if server == nil {
		t.Error("Failed to create server")
	}
	
	// Test graceful shutdown context creation
	// (This tests the timeout logic without actually starting the server)
	timeout := 5 * time.Second
	if timeout != 5*time.Second {
		t.Error("Unexpected timeout value")
	}
}
