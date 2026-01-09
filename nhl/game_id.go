package nhl

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// GameID is a wrapper type for NHL game identifiers.
// Game IDs are 10-digit integers in the format: SSSGTNNNN where:
// - SSS is the first 4 digits of the season (e.g., 2023 for 2023-2024)
// - GT is the game type (01=preseason, 02=regular, 03=playoffs, 04=all-star, 12=PWHL showcase)
// - NNNN is the game number
type GameID int64

// NewGameID creates a new GameID from an int64.
func NewGameID(id int64) GameID {
	return GameID(id)
}

// AsInt64 returns the GameID as an int64.
func (g GameID) AsInt64() int64 {
	return int64(g)
}

// String implements the fmt.Stringer interface.
func (g GameID) String() string {
	return strconv.FormatInt(int64(g), 10)
}

// MarshalJSON implements json.Marshaler.
// GameIDs are marshaled as integers in JSON.
func (g GameID) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(g))
}

// UnmarshalJSON implements json.Unmarshaler.
// GameIDs can be unmarshaled from either integers or strings.
func (g *GameID) UnmarshalJSON(data []byte) error {
	// Try unmarshaling as integer first
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*g = GameID(i)
		return nil
	}

	// Try unmarshaling as string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("game ID must be an integer or string: %w", err)
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid game ID string: %w", err)
	}

	*g = GameID(i)
	return nil
}

// Season extracts the season from the game ID.
// Returns the season in YYYYYYYY format (e.g., 20232024).
func (g GameID) Season() (Season, error) {
	id := int64(g)
	if id < 1000000000 || id > 9999999999 {
		return Season{}, fmt.Errorf("invalid game ID: %d", id)
	}

	// Extract first 4 digits for the start year
	startYear := int(id / 1000000)

	return NewSeason(startYear), nil
}

// GameType extracts the game type code from the game ID.
// Returns: 01 (preseason), 02 (regular season), 03 (playoffs), 04 (all-star)
func (g GameID) GameType() (int, error) {
	id := int64(g)
	if id < 1000000000 || id > 9999999999 {
		return 0, fmt.Errorf("invalid game ID: %d", id)
	}

	// Extract digits 5-6 (0-indexed positions 4-5)
	gameType := int((id / 10000) % 100)

	return gameType, nil
}

// GameNumber extracts the game number from the game ID.
func (g GameID) GameNumber() (int, error) {
	id := int64(g)
	if id < 1000000000 || id > 9999999999 {
		return 0, fmt.Errorf("invalid game ID: %d", id)
	}

	// Extract last 4 digits
	gameNumber := int(id % 10000)

	return gameNumber, nil
}

// Validate checks if the GameID is in a valid format.
func (g GameID) Validate() error {
	id := int64(g)

	// Check if it's a 10-digit number
	if id < 1000000000 || id > 9999999999 {
		return fmt.Errorf("game ID must be 10 digits, got: %d", id)
	}

	// Validate game type
	gameTypeInt, err := g.GameType()
	if err != nil {
		return err
	}

	if !GameType(gameTypeInt).IsValid() {
		return fmt.Errorf("invalid game type: %02d", gameTypeInt)
	}

	return nil
}

// GameIDFromInt creates a GameID from an int.
func GameIDFromInt(i int) GameID {
	return GameID(i)
}

// GameIDFromString parses a GameID from a string.
func GameIDFromString(s string) (GameID, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid game ID string: %w", err)
	}

	return GameID(i), nil
}

// MustGameIDFromString parses a GameID from a string and panics on error.
// This should only be used in tests or when you are certain the input is valid.
func MustGameIDFromString(s string) GameID {
	id, err := GameIDFromString(s)
	if err != nil {
		panic(err)
	}
	return id
}
