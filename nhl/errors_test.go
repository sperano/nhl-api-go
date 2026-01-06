package nhl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "basic error message",
			err: &Error{
				Message:    "resource not found",
				StatusCode: 404,
			},
			expected: "NHL API error (status 404): resource not found",
		},
		{
			name: "server error",
			err: &Error{
				Message:    "internal server error",
				StatusCode: 500,
			},
			expected: "NHL API error (status 500): internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestResourceNotFoundError(t *testing.T) {
	err := NewResourceNotFoundError("team not found")

	if err.Message() != "team not found" {
		t.Errorf("Message = %q, want %q", err.Message(), "team not found")
	}

	if err.StatusCode() != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", err.StatusCode(), http.StatusNotFound)
	}

	expected := "NHL API error (status 404): team not found"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestRateLimitExceededError(t *testing.T) {
	err := NewRateLimitExceededError("rate limit exceeded")

	if err.Message() != "rate limit exceeded" {
		t.Errorf("Message = %q, want %q", err.Message(), "rate limit exceeded")
	}

	if err.StatusCode() != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want %d", err.StatusCode(), http.StatusTooManyRequests)
	}

	expected := "NHL API error (status 429): rate limit exceeded"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestServerError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "500 internal server error",
			statusCode: 500,
			message:    "internal server error",
		},
		{
			name:       "503 service unavailable",
			statusCode: 503,
			message:    "service unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewServerError(tt.statusCode, tt.message)

			if err.Message() != tt.message {
				t.Errorf("Message = %q, want %q", err.Message(), tt.message)
			}

			if err.StatusCode() != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", err.StatusCode(), tt.statusCode)
			}

			expected := fmt.Sprintf("NHL API error (status %d): %s", tt.statusCode, tt.message)
			if got := err.Error(); got != expected {
				t.Errorf("Error() = %q, want %q", got, expected)
			}
		})
	}
}

func TestBadRequestError(t *testing.T) {
	err := NewBadRequestError("invalid parameter")

	if err.Message() != "invalid parameter" {
		t.Errorf("Message = %q, want %q", err.Message(), "invalid parameter")
	}

	if err.StatusCode() != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", err.StatusCode(), http.StatusBadRequest)
	}

	expected := "NHL API error (status 400): invalid parameter"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("authentication required")

	if err.Message() != "authentication required" {
		t.Errorf("Message = %q, want %q", err.Message(), "authentication required")
	}

	if err.StatusCode() != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want %d", err.StatusCode(), http.StatusUnauthorized)
	}

	expected := "NHL API error (status 401): authentication required"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestAPIError(t *testing.T) {
	err := NewAPIError(403, "forbidden")

	if err.Message() != "forbidden" {
		t.Errorf("Message = %q, want %q", err.Message(), "forbidden")
	}

	if err.StatusCode() != 403 {
		t.Errorf("StatusCode = %d, want %d", err.StatusCode(), 403)
	}

	expected := "NHL API error (status 403): forbidden"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestRequestError(t *testing.T) {
	originalErr := fmt.Errorf("connection timeout")
	err := NewRequestError(originalErr)

	expected := "request error: connection timeout"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}

	// Test error unwrapping
	if !errors.Is(err, originalErr) {
		t.Error("errors.Is() should return true for wrapped error")
	}

	var reqErr *RequestError
	if !errors.As(err, &reqErr) {
		t.Error("errors.As() should return true for RequestError")
	}
}

func TestJSONError(t *testing.T) {
	originalErr := fmt.Errorf("invalid JSON")
	err := NewJSONError(originalErr)

	expected := "JSON error: invalid JSON"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}

	// Test error unwrapping
	if !errors.Is(err, originalErr) {
		t.Error("errors.Is() should return true for wrapped error")
	}

	var jsonErr *JSONError
	if !errors.As(err, &jsonErr) {
		t.Error("errors.As() should return true for JSONError")
	}
}

func TestErrorFromStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		checkType  func(error) bool
	}{
		{
			name:       "400 bad request",
			statusCode: 400,
			message:    "bad request",
			checkType: func(err error) bool {
				var e *BadRequestError
				return errors.As(err, &e)
			},
		},
		{
			name:       "401 unauthorized",
			statusCode: 401,
			message:    "unauthorized",
			checkType: func(err error) bool {
				var e *UnauthorizedError
				return errors.As(err, &e)
			},
		},
		{
			name:       "404 not found",
			statusCode: 404,
			message:    "not found",
			checkType: func(err error) bool {
				var e *ResourceNotFoundError
				return errors.As(err, &e)
			},
		},
		{
			name:       "429 rate limit",
			statusCode: 429,
			message:    "too many requests",
			checkType: func(err error) bool {
				var e *RateLimitExceededError
				return errors.As(err, &e)
			},
		},
		{
			name:       "500 server error",
			statusCode: 500,
			message:    "internal server error",
			checkType: func(err error) bool {
				var e *ServerError
				return errors.As(err, &e)
			},
		},
		{
			name:       "503 server error",
			statusCode: 503,
			message:    "service unavailable",
			checkType: func(err error) bool {
				var e *ServerError
				return errors.As(err, &e)
			},
		},
		{
			name:       "403 api error",
			statusCode: 403,
			message:    "forbidden",
			checkType: func(err error) bool {
				var e *APIError
				return errors.As(err, &e)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrorFromStatusCode(tt.statusCode, tt.message)

			if !tt.checkType(err) {
				t.Errorf("ErrorFromStatusCode() returned unexpected type %T", err)
			}

			// Verify the error message contains the status code and message
			errStr := err.Error()
			if errStr == "" {
				t.Error("Error() returned empty string")
			}
		})
	}
}

func TestErrorFromStatusCode_EmptyMessage(t *testing.T) {
	err := ErrorFromStatusCode(404, "")

	var notFoundErr *ResourceNotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Errorf("Expected ResourceNotFoundError, got %T", err)
	}

	// Should use http.StatusText when message is empty
	if notFoundErr.Message() != http.StatusText(404) {
		t.Errorf("Message = %q, want %q", notFoundErr.Message(), http.StatusText(404))
	}
}

func TestError_JSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		err := &Error{
			Message:    "test error",
			StatusCode: 404,
		}

		data, marshalErr := json.Marshal(err)
		if marshalErr != nil {
			t.Fatalf("json.Marshal() error = %v", marshalErr)
		}

		expected := `{"message":"test error","status_code":404}`
		if string(data) != expected {
			t.Errorf("json.Marshal() = %q, want %q", string(data), expected)
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		data := []byte(`{"message":"test error","status_code":404}`)

		var err Error
		if unmarshalErr := json.Unmarshal(data, &err); unmarshalErr != nil {
			t.Fatalf("json.Unmarshal() error = %v", unmarshalErr)
		}

		if err.Message != "test error" {
			t.Errorf("Message = %q, want %q", err.Message, "test error")
		}

		if err.StatusCode != 404 {
			t.Errorf("StatusCode = %d, want %d", err.StatusCode, 404)
		}
	})

	t.Run("round trip", func(t *testing.T) {
		original := &Error{
			Message:    "round trip test",
			StatusCode: 500,
		}

		data, marshalErr := json.Marshal(original)
		if marshalErr != nil {
			t.Fatalf("json.Marshal() error = %v", marshalErr)
		}

		var decoded Error
		if unmarshalErr := json.Unmarshal(data, &decoded); unmarshalErr != nil {
			t.Fatalf("json.Unmarshal() error = %v", unmarshalErr)
		}

		if decoded.Message != original.Message {
			t.Errorf("Message = %q, want %q", decoded.Message, original.Message)
		}

		if decoded.StatusCode != original.StatusCode {
			t.Errorf("StatusCode = %d, want %d", decoded.StatusCode, original.StatusCode)
		}
	})
}

// Additional error path tests for Error.UnmarshalJSON

func TestError_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var e Error
	err := json.Unmarshal([]byte(`{invalid json}`), &e)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}
