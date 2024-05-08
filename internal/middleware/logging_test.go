package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock handler that will be wrapped by the logging middleware
func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from Mock Handler"))
}

// TestLoggingMiddleware tests the logging middleware
func TestLoggingMiddleware(t *testing.T) {
	// Create a request to pass to our middleware
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Wrap the mock handler in the logging middleware
	handler := Logging(http.HandlerFunc(mockHandler))

	// Call the handler with the request and response recorder
	handler.ServeHTTP(rr, req)

	// Verify the status code is what we expect
	if rr.Code != http.StatusOK {
		t.Errorf("Unexpected status code: got %v, want %v", rr.Code, http.StatusOK)
	}

	// Verify the response body is what we expect
	expected := "Hello from Mock Handler"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Unexpected response body: got %v, want %v", rr.Body.String(), expected)
	}
}
