package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	// Setup a buffer to capture log output
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, nil))
	slog.SetDefault(logger)

	// Create a test handler that does nothing
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do nothing
	})

	// Wrap the test handler with our logging middleware
	loggingHandler := Logging(testHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Call the middleware
	loggingHandler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify log contains expected fields
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "Request received") {
		t.Errorf("Expected log to contain 'Request received', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "status=200") {
		t.Errorf("Expected log to contain status code, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "elapsed_time=") {
		t.Errorf("Expected log to contain elapsed time, got: %s", logOutput)
	}
}
