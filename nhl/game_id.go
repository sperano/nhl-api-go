package nhl

import "fmt"

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
