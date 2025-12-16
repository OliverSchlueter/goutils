package problems

import (
	"fmt"
	"net/http"
	"time"
)

// MethodNotAllowed creates a Problem instance for HTTP 405 Method Not Allowed errors.
// It takes the HTTP method that was attempted and a list of allowed methods as parameters.
func MethodNotAllowed(method string, allowedMethods []string) *Problem {
	return &Problem{
		Type:      "MethodNotAllowed",
		Title:     "Method not allowed",
		Detail:    fmt.Sprintf("The HTTP method %s is not allowed for this resource. Allowed methods are: %v", method, allowedMethods),
		Status:    http.StatusMethodNotAllowed,
		Timestamp: time.Now(),
	}
}

// NotFound creates a Problem instance for HTTP 404 Not Found errors.
// It takes the resource type (e.g., "User", "Project") and the specific resource identifier (e.g., "12345") as parameters.
func NotFound(resourceType, resource string) *Problem {
	return &Problem{
		Type:      "NotFound",
		Title:     fmt.Sprintf("%s not found", resourceType),
		Detail:    fmt.Sprintf("The requested %s '%s' could not be found.", resourceType, resource),
		Status:    http.StatusNotFound,
		Timestamp: time.Now(),
	}
}

func AlreadyExists(resourceType, resource string) *Problem {
	return &Problem{
		Type:      "AlreadyExists",
		Title:     fmt.Sprintf("%s already exists", resourceType),
		Detail:    fmt.Sprintf("The requested %s '%s' already exists.", resourceType, resource),
		Status:    http.StatusConflict,
		Timestamp: time.Now(),
	}
}

// Unauthorized creates a Problem instance for HTTP 401 Unauthorized errors.
// It indicates that the request requires user authentication.
func Unauthorized() *Problem {
	return &Problem{
		Type:      "Unauthorized",
		Title:     "Unauthorized",
		Detail:    "Authentication is required to access this resource.",
		Status:    http.StatusUnauthorized,
		Timestamp: time.Now(),
	}
}

// Forbidden creates a Problem instance for HTTP 403 Forbidden errors.
// It indicates that the server understood the request, but refuses to authorize it.
func Forbidden() *Problem {
	return &Problem{
		Type:      "Forbidden",
		Title:     "Forbidden",
		Detail:    "You do not have permission to access this resource.",
		Status:    http.StatusForbidden,
		Timestamp: time.Now(),
	}
}

// WrongContentType creates a Problem instance for HTTP 415 Unsupported Media Type errors.
// It indicates that the request's Content-Type header does not match the expected type.
func WrongContentType(expected, actual string) *Problem {
	return &Problem{
		Type:      "WrongContentType",
		Title:     "Wrong Content-Type",
		Detail:    fmt.Sprintf("Expected Content-Type '%s', but got '%s'.", expected, actual),
		Status:    http.StatusUnsupportedMediaType,
		Timestamp: time.Now(),
	}
}

// WrongAcceptType creates a Problem instance for HTTP 406 Not Acceptable errors.
// It indicates that the request's Accept header does not match the expected type.
func WrongAcceptType(expected, actual string) *Problem {
	return &Problem{
		Type:      "WrongAcceptType",
		Title:     "Wrong Accept header",
		Detail:    fmt.Sprintf("Expected Accept header '%s', but got '%s'.", expected, actual),
		Status:    http.StatusNotAcceptable,
		Timestamp: time.Now(),
	}
}

// CouldNotDecodeBody creates a Problem instance for HTTP 400 Bad Request errors.
// It indicates that the request body could not be decoded, which is typically due to an invalid format.
func CouldNotDecodeBody() *Problem {
	return &Problem{
		Type:      "CouldNotDecodeBody",
		Title:     "Could not decode request body",
		Detail:    "The request body could not be decoded. Please check the request format.",
		Status:    http.StatusBadRequest,
		Timestamp: time.Now(),
	}
}

func ValidationError(field, reason string) *Problem {
	return &Problem{
		Type:      "ValidationError",
		Title:     "Validation error",
		Detail:    fmt.Sprintf("Validation failed for field '%s': %s", field, reason),
		Status:    http.StatusBadRequest,
		Timestamp: time.Now(),
	}
}

func TooManyRequests() *Problem {
	return &Problem{
		Type:      "TooManyRequests",
		Title:     "Too Many Requests",
		Detail:    "You have sent too many requests in a given amount of time. Please try again later.",
		Status:    http.StatusTooManyRequests,
		Timestamp: time.Now(),
	}
}

// InternalServerError creates a Problem instance for HTTP 500 Internal Server Error.
// It takes a detail message that provides more context about the error.
func InternalServerError(detail string) *Problem {
	return &Problem{
		Type:      "InternalServerError",
		Title:     "Internal Server Error",
		Detail:    detail,
		Status:    http.StatusInternalServerError,
		Timestamp: time.Now(),
	}
}

// NotImplemented creates a Problem instance for HTTP 501 Not Implemented errors.
// It indicates that the server does not support the functionality required to fulfill the request.
func NotImplemented() *Problem {
	return &Problem{
		Type:      "NotImplemented",
		Title:     "Not Implemented",
		Detail:    "This feature is not implemented yet.",
		Status:    http.StatusNotImplemented,
		Timestamp: time.Now(),
	}
}
