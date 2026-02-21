package nhl

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// unmarshalNumericID handles JSON unmarshaling for numeric ID types.
// Accepts both integer and string representations for flexibility with different API responses.
func unmarshalNumericID(data []byte, typeName string) (int64, error) {
	// Try unmarshaling as integer first (most common case)
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		return i, nil
	}

	// Try unmarshaling as string (some APIs return IDs as strings)
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return 0, fmt.Errorf("%s must be an integer or string: %w", typeName, err)
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s string: %w", typeName, err)
	}
	return i, nil
}

// parseNumericID parses a string into an int64 for ID types.
func parseNumericID(s string, typeName string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s string: %w", typeName, err)
	}
	return i, nil
}
