package nhl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name:     "not found",
			err:      NewAPIError(404, "resource not found"),
			expected: "NHL API error (status 404): resource not found",
		},
		{
			name:     "server error",
			err:      NewAPIError(500, "internal server error"),
			expected: "NHL API error (status 500): internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestAPIError_Is(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		wantIs   bool
	}{
		{
			name:   "404 matches ErrNotFound",
			err:    NewAPIError(404, "team not found"),
			target: ErrNotFound,
			wantIs: true,
		},
		{
			name:   "429 matches ErrRateLimited",
			err:    NewAPIError(429, "slow down"),
			target: ErrRateLimited,
			wantIs: true,
		},
		{
			name:   "400 matches ErrBadRequest",
			err:    NewAPIError(400, "invalid"),
			target: ErrBadRequest,
			wantIs: true,
		},
		{
			name:   "401 matches ErrUnauthorized",
			err:    NewAPIError(401, "denied"),
			target: ErrUnauthorized,
			wantIs: true,
		},
		{
			name:   "500 matches ErrServerError",
			err:    NewAPIError(500, "boom"),
			target: ErrServerError,
			wantIs: true,
		},
		{
			name:   "503 matches ErrServerError (any 5xx)",
			err:    NewAPIError(503, "unavailable"),
			target: ErrServerError,
			wantIs: true,
		},
		{
			name:   "502 matches ErrServerError (any 5xx)",
			err:    NewAPIError(502, "bad gateway"),
			target: ErrServerError,
			wantIs: true,
		},
		{
			name:   "404 does not match ErrBadRequest",
			err:    NewAPIError(404, "not found"),
			target: ErrBadRequest,
			wantIs: false,
		},
		{
			name:   "403 does not match any sentinel",
			err:    NewAPIError(403, "forbidden"),
			target: ErrNotFound,
			wantIs: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.wantIs {
				t.Errorf("errors.Is() = %v, want %v", got, tt.wantIs)
			}
		})
	}
}

func TestAPIError_As(t *testing.T) {
	err := ErrorFromStatusCode(404, "not found")

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("errors.As() should return true for *APIError")
	}

	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
	if apiErr.Message != "not found" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "not found")
	}
}

func TestAPIError_Fields(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{"bad request", 400, "invalid parameter"},
		{"unauthorized", 401, "authentication required"},
		{"not found", 404, "team not found"},
		{"rate limited", 429, "rate limit exceeded"},
		{"server error", 500, "internal server error"},
		{"service unavailable", 503, "service unavailable"},
		{"generic", 403, "forbidden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.statusCode, tt.message)

			if err.Message != tt.message {
				t.Errorf("Message = %q, want %q", err.Message, tt.message)
			}
			if err.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", err.StatusCode, tt.statusCode)
			}

			expected := fmt.Sprintf("NHL API error (status %d): %s", tt.statusCode, tt.message)
			if got := err.Error(); got != expected {
				t.Errorf("Error() = %q, want %q", got, expected)
			}
		})
	}
}

func TestRequestError(t *testing.T) {
	originalErr := fmt.Errorf("connection timeout")
	err := NewRequestError(originalErr)

	expected := "request error: connection timeout"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}

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
		sentinel   error
	}{
		{"400 bad request", 400, "bad request", ErrBadRequest},
		{"401 unauthorized", 401, "unauthorized", ErrUnauthorized},
		{"404 not found", 404, "not found", ErrNotFound},
		{"429 rate limit", 429, "too many requests", ErrRateLimited},
		{"500 server error", 500, "internal server error", ErrServerError},
		{"503 server error", 503, "service unavailable", ErrServerError},
		{"403 generic", 403, "forbidden", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ErrorFromStatusCode(tt.statusCode, tt.message)

			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatal("ErrorFromStatusCode should return *APIError")
			}

			if apiErr.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.statusCode)
			}

			if tt.sentinel != nil && !errors.Is(err, tt.sentinel) {
				t.Errorf("errors.Is(err, sentinel) = false, want true")
			}
		})
	}
}

func TestErrorFromStatusCode_EmptyMessage(t *testing.T) {
	err := ErrorFromStatusCode(404, "")

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("Expected *APIError")
	}

	if apiErr.Message != http.StatusText(404) {
		t.Errorf("Message = %q, want %q", apiErr.Message, http.StatusText(404))
	}

	if !errors.Is(err, ErrNotFound) {
		t.Error("errors.Is(err, ErrNotFound) should be true")
	}
}

func TestAPIError_JSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		err := NewAPIError(404, "test error")

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

		var err APIError
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
		original := NewAPIError(500, "round trip test")

		data, marshalErr := json.Marshal(original)
		if marshalErr != nil {
			t.Fatalf("json.Marshal() error = %v", marshalErr)
		}

		var decoded APIError
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

func TestAPIError_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var e APIError
	err := json.Unmarshal([]byte(`{invalid json}`), &e)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}
