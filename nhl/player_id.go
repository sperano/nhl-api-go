package nhl

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// PlayerID is a wrapper type for NHL player identifiers.
// Player IDs are numeric identifiers assigned to each player (e.g., 8478402 for Connor McDavid).
type PlayerID int64

// NewPlayerID creates a new PlayerID from an int64.
func NewPlayerID(id int64) PlayerID {
	return PlayerID(id)
}

// AsInt64 returns the PlayerID as an int64.
func (p PlayerID) AsInt64() int64 {
	return int64(p)
}

// String implements the fmt.Stringer interface.
func (p PlayerID) String() string {
	return strconv.FormatInt(int64(p), 10)
}

// MarshalJSON implements json.Marshaler.
// PlayerIDs are marshaled as integers in JSON.
func (p PlayerID) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(p))
}

// UnmarshalJSON implements json.Unmarshaler.
// PlayerIDs can be unmarshaled from either integers or strings.
func (p *PlayerID) UnmarshalJSON(data []byte) error {
	// Try unmarshaling as integer first
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*p = PlayerID(i)
		return nil
	}

	// Try unmarshaling as string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("player ID must be an integer or string: %w", err)
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid player ID string: %w", err)
	}

	*p = PlayerID(i)
	return nil
}

// PlayerIDFromInt creates a PlayerID from an int.
func PlayerIDFromInt(i int) PlayerID {
	return PlayerID(i)
}

// PlayerIDFromString parses a PlayerID from a string.
func PlayerIDFromString(s string) (PlayerID, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid player ID string: %w", err)
	}

	return PlayerID(i), nil
}

// MustPlayerIDFromString parses a PlayerID from a string and panics on error.
// This should only be used in tests or when you are certain the input is valid.
func MustPlayerIDFromString(s string) PlayerID {
	id, err := PlayerIDFromString(s)
	if err != nil {
		panic(err)
	}
	return id
}
