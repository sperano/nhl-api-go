package nhl

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// TeamID is a wrapper type for NHL team identifiers.
// Team IDs are numeric identifiers assigned to each team (e.g., 10 for Toronto Maple Leafs).
type TeamID int64

// NewTeamID creates a new TeamID from an int64.
func NewTeamID(id int64) TeamID {
	return TeamID(id)
}

// AsInt64 returns the TeamID as an int64.
func (t TeamID) AsInt64() int64 {
	return int64(t)
}

// String implements the fmt.Stringer interface.
func (t TeamID) String() string {
	return strconv.FormatInt(int64(t), 10)
}

// MarshalJSON implements json.Marshaler.
// TeamIDs are marshaled as integers in JSON.
func (t TeamID) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(t))
}

// UnmarshalJSON implements json.Unmarshaler.
// TeamIDs can be unmarshaled from either integers or strings.
func (t *TeamID) UnmarshalJSON(data []byte) error {
	// Try unmarshaling as integer first
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*t = TeamID(i)
		return nil
	}

	// Try unmarshaling as string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("team ID must be an integer or string: %w", err)
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid team ID string: %w", err)
	}

	*t = TeamID(i)
	return nil
}

// TeamIDFromInt creates a TeamID from an int.
func TeamIDFromInt(i int) TeamID {
	return TeamID(i)
}

// TeamIDFromString parses a TeamID from a string.
func TeamIDFromString(s string) (TeamID, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid team ID string: %w", err)
	}

	return TeamID(i), nil
}

// MustTeamIDFromString parses a TeamID from a string and panics on error.
// This should only be used in tests or when you are certain the input is valid.
func MustTeamIDFromString(s string) TeamID {
	id, err := TeamIDFromString(s)
	if err != nil {
		panic(err)
	}
	return id
}
