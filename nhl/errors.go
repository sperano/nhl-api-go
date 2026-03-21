package nhl

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Sentinel errors for well-known HTTP status codes.
// Use errors.Is(err, nhl.ErrNotFound) to check error types.
var (
	ErrBadRequest    = &APIError{StatusCode: http.StatusBadRequest}
	ErrUnauthorized  = &APIError{StatusCode: http.StatusUnauthorized}
	ErrNotFound      = &APIError{StatusCode: http.StatusNotFound}
	ErrRateLimited   = &APIError{StatusCode: http.StatusTooManyRequests}
	ErrServerError   = &APIError{StatusCode: http.StatusInternalServerError}
)

// APIError represents an NHL API error with an HTTP status code and message.
// Use errors.Is with sentinel errors (ErrNotFound, ErrRateLimited, etc.) to
// check for specific status codes. Matching is done by status code, so any
// 404 APIError will match ErrNotFound regardless of message.
type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("NHL API error (status %d): %s", e.StatusCode, e.Message)
}

// Is supports errors.Is matching by status code. For 5xx errors, any
// APIError with a 5xx status code matches ErrServerError.
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	// Match any 5xx against ErrServerError
	if t.StatusCode == http.StatusInternalServerError && e.StatusCode >= 500 {
		return true
	}
	return e.StatusCode == t.StatusCode
}

// NewAPIError creates a new APIError with the given status code and message.
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// ErrorFromStatusCode creates an APIError from an HTTP status code.
// If message is empty, the standard HTTP status text is used.
func ErrorFromStatusCode(statusCode int, message string) error {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	return NewAPIError(statusCode, message)
}

// RequestError wraps errors that occur during HTTP request execution.
type RequestError struct {
	Err error
}

// NewRequestError creates a new RequestError.
func NewRequestError(err error) *RequestError {
	return &RequestError{Err: err}
}

// Error implements the error interface.
func (e *RequestError) Error() string {
	return fmt.Sprintf("request error: %v", e.Err)
}

// Unwrap returns the wrapped error for errors.Is and errors.As.
func (e *RequestError) Unwrap() error {
	return e.Err
}

// JSONError wraps errors that occur during JSON marshaling/unmarshaling.
type JSONError struct {
	Err error
}

// NewJSONError creates a new JSONError.
func NewJSONError(err error) *JSONError {
	return &JSONError{Err: err}
}

// Error implements the error interface.
func (e *JSONError) Error() string {
	return fmt.Sprintf("JSON error: %v", e.Err)
}

// Unwrap returns the wrapped error for errors.Is and errors.As.
func (e *JSONError) Unwrap() error {
	return e.Err
}

// UnmarshalJSON implements custom JSON unmarshaling for APIError.
func (e *APIError) UnmarshalJSON(data []byte) error {
	var raw struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	e.Message = raw.Message
	e.StatusCode = raw.StatusCode
	return nil
}

// MarshalJSON implements custom JSON marshaling for APIError.
func (e *APIError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	}{
		Message:    e.Message,
		StatusCode: e.StatusCode,
	})
}
