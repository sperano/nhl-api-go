package nhl

import (
	"encoding/json"
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
	i, err := unmarshalNumericID(data, "player ID")
	if err != nil {
		return err
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
	i, err := parseNumericID(s, "player ID")
	if err != nil {
		return 0, err
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
