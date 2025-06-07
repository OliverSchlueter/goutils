package problems_test

import (
	"encoding/json"
	"github.com/OliverSchlueter/goutils/broker"
	"github.com/OliverSchlueter/goutils/problems"
	"github.com/nats-io/nats.go"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProblem_WriteToHTTP(t *testing.T) {
	problem := &problems.Problem{
		Type:      "TestError",
		Title:     "Test Error",
		Detail:    "This is a test error",
		Status:    http.StatusBadRequest,
		Timestamp: time.Now(),
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Send the problem
	problem.WriteToHTTP(rr)

	// Check status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	assert.Equal(t, "application/problem+json", contentType)

	// Verify the response body
	var responseProblem problems.Problem
	err := json.Unmarshal(rr.Body.Bytes(), &responseProblem)
	require.NoError(t, err)

	assert.Equal(t, problem.Type, responseProblem.Type)
	assert.Equal(t, problem.Title, responseProblem.Title)
	assert.Equal(t, problem.Detail, responseProblem.Detail)
	assert.Equal(t, problem.Status, responseProblem.Status)
	assert.WithinDuration(t, problem.Timestamp, responseProblem.Timestamp, time.Second)
}

func TestProblem_WriteToBroker(t *testing.T) {
	// Create a problem
	problem := &problems.Problem{
		Type:      "TestError",
		Title:     "Test Error",
		Detail:    "This is a test error",
		Status:    http.StatusBadRequest,
		Timestamp: time.Now(),
	}

	// Create a fake broker
	fakeBroker := broker.NewFakeBroker()

	// Create a channel to receive the published message
	receivedMsg := make(chan []byte, 1)

	// Subscribe to messages
	err := fakeBroker.Subscribe("test.subject", func(msg *nats.Msg) {
		receivedMsg <- msg.Data
	})
	require.NoError(t, err)

	// Write problem to broker
	problem.WriteToBroker(fakeBroker, "test.subject")

	// Wait for the message to be received
	var data []byte
	select {
	case data = <-receivedMsg:
		// Got data
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for message")
	}

	// Verify the published data
	var receivedProblem problems.Problem
	err = json.Unmarshal(data, &receivedProblem)
	require.NoError(t, err)

	assert.Equal(t, problem.Type, receivedProblem.Type)
	assert.Equal(t, problem.Title, receivedProblem.Title)
	assert.Equal(t, problem.Detail, receivedProblem.Detail)
	assert.Equal(t, problem.Status, receivedProblem.Status)
	assert.WithinDuration(t, problem.Timestamp, receivedProblem.Timestamp, time.Second)
}

func TestMethodNotAllowed(t *testing.T) {
	method := "POST"
	allowedMethods := []string{"GET", "PUT"}

	problem := problems.MethodNotAllowed(method, allowedMethods)

	assert.Equal(t, "MethodNotAllowed", problem.Type)
	assert.Equal(t, "Method not allowed", problem.Title)
	assert.Contains(t, problem.Detail, method)
	assert.Contains(t, problem.Detail, "GET")
	assert.Contains(t, problem.Detail, "PUT")
	assert.Equal(t, http.StatusMethodNotAllowed, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestNotFound(t *testing.T) {
	resourceType := "User"
	resource := "12345"

	problem := problems.NotFound(resourceType, resource)

	assert.Equal(t, "NotFound", problem.Type)
	assert.Equal(t, "User not found", problem.Title)
	assert.Contains(t, problem.Detail, resourceType)
	assert.Contains(t, problem.Detail, resource)
	assert.Equal(t, http.StatusNotFound, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestAlreadyExists(t *testing.T) {
	resourceType := "User"
	resource := "12345"

	problem := problems.AlreadyExists(resourceType, resource)

	assert.Equal(t, "AlreadyExists", problem.Type)
	assert.Equal(t, "User already exists", problem.Title)
	assert.Contains(t, problem.Detail, resourceType)
	assert.Contains(t, problem.Detail, resource)
	assert.Equal(t, http.StatusConflict, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestUnauthorized(t *testing.T) {
	problem := problems.Unauthorized()

	assert.Equal(t, "Unauthorized", problem.Type)
	assert.Equal(t, "Unauthorized", problem.Title)
	assert.Contains(t, problem.Detail, "Authentication")
	assert.Equal(t, http.StatusUnauthorized, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestForbidden(t *testing.T) {
	problem := problems.Forbidden()

	assert.Equal(t, "Forbidden", problem.Type)
	assert.Equal(t, "Forbidden", problem.Title)
	assert.Contains(t, problem.Detail, "permission")
	assert.Equal(t, http.StatusForbidden, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestWrongContentType(t *testing.T) {
	expected := "application/json"
	actual := "text/plain"

	problem := problems.WrongContentType(expected, actual)

	assert.Equal(t, "WrongContentType", problem.Type)
	assert.Equal(t, "Wrong Content-Type", problem.Title)
	assert.Contains(t, problem.Detail, expected)
	assert.Contains(t, problem.Detail, actual)
	assert.Equal(t, http.StatusUnsupportedMediaType, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestWrongAcceptType(t *testing.T) {
	expected := "application/json"
	actual := "text/plain"

	problem := problems.WrongAcceptType(expected, actual)

	assert.Equal(t, "WrongAcceptType", problem.Type)
	assert.Equal(t, "Wrong Accept header", problem.Title)
	assert.Contains(t, problem.Detail, expected)
	assert.Contains(t, problem.Detail, actual)
	assert.Equal(t, http.StatusNotAcceptable, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestCouldNotDecodeBody(t *testing.T) {
	problem := problems.CouldNotDecodeBody()

	assert.Equal(t, "CouldNotDecodeBody", problem.Type)
	assert.Equal(t, "Could not decode request body", problem.Title)
	assert.Contains(t, problem.Detail, "request body")
	assert.Equal(t, http.StatusBadRequest, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestValidationError(t *testing.T) {
	field := "email"
	reason := "invalid format"

	problem := problems.ValidationError(field, reason)

	assert.Equal(t, "ValidationError", problem.Type)
	assert.Equal(t, "Validation error", problem.Title)
	assert.Contains(t, problem.Detail, field)
	assert.Contains(t, problem.Detail, reason)
	assert.Equal(t, http.StatusBadRequest, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestInternalServerError(t *testing.T) {
	detail := "Database connection failed"

	problem := problems.InternalServerError(detail)

	assert.Equal(t, "InternalServerError", problem.Type)
	assert.Equal(t, "Internal Server Error", problem.Title)
	assert.Equal(t, detail, problem.Detail)
	assert.Equal(t, http.StatusInternalServerError, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}

func TestNotImplemented(t *testing.T) {
	problem := problems.NotImplemented()

	assert.Equal(t, "NotImplemented", problem.Type)
	assert.Equal(t, "Not Implemented", problem.Title)
	assert.Contains(t, problem.Detail, "not implemented")
	assert.Equal(t, http.StatusNotImplemented, problem.Status)
	assert.WithinDuration(t, time.Now(), problem.Timestamp, time.Second)
}
