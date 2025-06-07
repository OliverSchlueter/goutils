package middleware

import (
	"bytes"
	"github.com/nats-io/nats.go"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestNatsLogging(t *testing.T) {
	// Setup a buffer to capture log output
	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuffer, nil))
	slog.SetDefault(logger)

	// Create a mock NATS message
	msg := &nats.Msg{
		Subject: "test.subject",
		Reply:   "test.reply",
		Data:    []byte("test data"),
	}

	// Create a handler function that will be wrapped
	handlerCalled := false
	handler := func(msg *nats.Msg) {
		handlerCalled = true
		// Simulate some work
		time.Sleep(10 * time.Millisecond)
	}

	// Wrap the handler with our logging middleware
	wrappedHandler := NatsLogging(handler)

	// Call the wrapped handler
	wrappedHandler(msg)

	// Verify the original handler was called
	if !handlerCalled {
		t.Error("Original handler was not called")
	}

	// Verify log contains expected fields
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "NATS message received") {
		t.Errorf("Expected log to contain 'NATS message received', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "subject=test.subject") {
		t.Errorf("Expected log to contain subject, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "reply=test.reply") {
		t.Errorf("Expected log to contain reply, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "elapsed_time=") {
		t.Errorf("Expected log to contain elapsed time, got: %s", logOutput)
	}
}
