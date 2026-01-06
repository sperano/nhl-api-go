package nhl

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error represents an NHL API error with an HTTP status code and message.
type Error struct {
	Message    string
	StatusCode int
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("NHL API error (status %d): %s", e.StatusCode, e.Message)
}

// ResourceNotFoundError indicates the requested resource was not found (404).
type ResourceNotFoundError struct {
	error *Error
}

// NewResourceNotFoundError creates a new ResourceNotFoundError.
func NewResourceNotFoundError(message string) *ResourceNotFoundError {
	return &ResourceNotFoundError{
		error: &Error{
			Message:    message,
			StatusCode: http.StatusNotFound,
		},
	}
}

// Error implements the error interface.
func (e *ResourceNotFoundError) Error() string {
	return e.error.Error()
}

// Message returns the error message.
func (e *ResourceNotFoundError) Message() string {
	return e.error.Message
}

// StatusCode returns the HTTP status code.
func (e *ResourceNotFoundError) StatusCode() int {
	return e.error.StatusCode
}

// RateLimitExceededError indicates rate limiting is in effect (429).
type RateLimitExceededError struct {
	error *Error
}

// NewRateLimitExceededError creates a new RateLimitExceededError.
func NewRateLimitExceededError(message string) *RateLimitExceededError {
	return &RateLimitExceededError{
		error: &Error{
			Message:    message,
			StatusCode: http.StatusTooManyRequests,
		},
	}
}

// Error implements the error interface.
func (e *RateLimitExceededError) Error() string {
	return e.error.Error()
}

// Message returns the error message.
func (e *RateLimitExceededError) Message() string {
	return e.error.Message
}

// StatusCode returns the HTTP status code.
func (e *RateLimitExceededError) StatusCode() int {
	return e.error.StatusCode
}

// ServerError indicates an internal server error (5xx).
type ServerError struct {
	error *Error
}

// NewServerError creates a new ServerError with the given status code and message.
func NewServerError(statusCode int, message string) *ServerError {
	return &ServerError{
		error: &Error{
			Message:    message,
			StatusCode: statusCode,
		},
	}
}

// Error implements the error interface.
func (e *ServerError) Error() string {
	return e.error.Error()
}

// Message returns the error message.
func (e *ServerError) Message() string {
	return e.error.Message
}

// StatusCode returns the HTTP status code.
func (e *ServerError) StatusCode() int {
	return e.error.StatusCode
}

// BadRequestError indicates a malformed request (400).
type BadRequestError struct {
	error *Error
}

// NewBadRequestError creates a new BadRequestError.
func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{
		error: &Error{
			Message:    message,
			StatusCode: http.StatusBadRequest,
		},
	}
}

// Error implements the error interface.
func (e *BadRequestError) Error() string {
	return e.error.Error()
}

// Message returns the error message.
func (e *BadRequestError) Message() string {
	return e.error.Message
}

// StatusCode returns the HTTP status code.
func (e *BadRequestError) StatusCode() int {
	return e.error.StatusCode
}

// UnauthorizedError indicates authentication is required or failed (401).
type UnauthorizedError struct {
	error *Error
}

// NewUnauthorizedError creates a new UnauthorizedError.
func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{
		error: &Error{
			Message:    message,
			StatusCode: http.StatusUnauthorized,
		},
	}
}

// Error implements the error interface.
func (e *UnauthorizedError) Error() string {
	return e.error.Error()
}

// Message returns the error message.
func (e *UnauthorizedError) Message() string {
	return e.error.Message
}

// StatusCode returns the HTTP status code.
func (e *UnauthorizedError) StatusCode() int {
	return e.error.StatusCode
}

// APIError represents a general API error with a custom status code.
type APIError struct {
	error *Error
}

// NewAPIError creates a new APIError with the given status code and message.
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		error: &Error{
			Message:    message,
			StatusCode: statusCode,
		},
	}
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return e.error.Error()
}

// Message returns the error message.
func (e *APIError) Message() string {
	return e.error.Message
}

// StatusCode returns the HTTP status code.
func (e *APIError) StatusCode() int {
	return e.error.StatusCode
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

// ErrorFromStatusCode creates an appropriate error based on HTTP status code.
func ErrorFromStatusCode(statusCode int, message string) error {
	if message == "" {
		message = http.StatusText(statusCode)
	}

	switch statusCode {
	case http.StatusBadRequest:
		return NewBadRequestError(message)
	case http.StatusUnauthorized:
		return NewUnauthorizedError(message)
	case http.StatusNotFound:
		return NewResourceNotFoundError(message)
	case http.StatusTooManyRequests:
		return NewRateLimitExceededError(message)
	default:
		if statusCode >= 500 {
			return NewServerError(statusCode, message)
		}
		return NewAPIError(statusCode, message)
	}
}

// UnmarshalJSON implements custom JSON unmarshaling for Error.
func (e *Error) UnmarshalJSON(data []byte) error {
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

// MarshalJSON implements custom JSON marshaling for Error.
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	}{
		Message:    e.Message,
		StatusCode: e.StatusCode,
	})
}
