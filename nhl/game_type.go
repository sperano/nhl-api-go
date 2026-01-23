package nhl

import (
	"encoding/json"
	"fmt"
)

// GameType represents the type of NHL game.
type GameType int

const (
	// GameTypePreseason represents a preseason game.
	GameTypePreseason GameType = 1
	// GameTypeRegularSeason represents a regular season game.
	GameTypeRegularSeason GameType = 2
	// GameTypePlayoffs represents a playoff game.
	GameTypePlayoffs GameType = 3
	// GameTypeAllStar represents an all-star game.
	GameTypeAllStar GameType = 4
	// GameTypeOlympics represents an Olympic tournament game.
	GameTypeOlympics GameType = 9
	// GameTypeYoungStars represents a YoungStars game (rookies/sophomores) during All-Star Weekend.
	GameTypeYoungStars GameType = 10
	// GameTypePWHLShowcase represents a PWHL 3-on-3 showcase game during All-Star Weekend.
	GameTypePWHLShowcase GameType = 12
	// GameTypeWomensAllStar represents a Women's All-Star game.
	GameTypeWomensAllStar GameType = 19
	// GameType4Nations represents a 4 Nations Face-Off game.
	GameType4Nations GameType = 20
)

// ToInt returns the integer representation of the GameType.
func (g GameType) ToInt() int {
	return int(g)
}

// String returns the string representation of the GameType.
func (g GameType) String() string {
	switch g {
	case GameTypePreseason:
		return "Preseason"
	case GameTypeRegularSeason:
		return "Regular Season"
	case GameTypePlayoffs:
		return "Playoffs"
	case GameTypeAllStar:
		return "All-Star"
	case GameTypeOlympics:
		return "Olympics"
	case GameTypeYoungStars:
		return "YoungStars"
	case GameTypePWHLShowcase:
		return "PWHL Showcase"
	case GameTypeWomensAllStar:
		return "Women's All-Star"
	case GameType4Nations:
		return "4 Nations Face-Off"
	default:
		return fmt.Sprintf("Unknown(%d)", g)
	}
}

// IsValid returns true if the GameType is one of the known valid types.
func (g GameType) IsValid() bool {
	switch g {
	case GameTypePreseason, GameTypeRegularSeason, GameTypePlayoffs, GameTypeAllStar, GameTypeOlympics, GameTypeYoungStars, GameTypePWHLShowcase, GameTypeWomensAllStar, GameType4Nations:
		return true
	default:
		return false
	}
}

// GameTypeFromInt parses an integer into a GameType.
// Returns an error if the integer is not a valid GameType.
func GameTypeFromInt(i int) (GameType, error) {
	g := GameType(i)
	if !g.IsValid() {
		return 0, fmt.Errorf("invalid game type: %d", i)
	}
	return g, nil
}

// MustGameTypeFromInt parses an integer into a GameType.
// Panics if the integer is not a valid GameType.
func MustGameTypeFromInt(i int) GameType {
	g, err := GameTypeFromInt(i)
	if err != nil {
		panic(err)
	}
	return g
}

// GameTypeFromString parses a string into a GameType.
// Accepts both numeric strings ("1", "2", etc.) and descriptive strings
// ("Preseason", "Regular Season", etc.).
// Returns an error if the string is not a valid GameType.
func GameTypeFromString(s string) (GameType, error) {
	switch s {
	case "1", "Preseason":
		return GameTypePreseason, nil
	case "2", "Regular Season", "RegularSeason":
		return GameTypeRegularSeason, nil
	case "3", "Playoffs":
		return GameTypePlayoffs, nil
	case "4", "All-Star", "AllStar":
		return GameTypeAllStar, nil
	case "9", "Olympics":
		return GameTypeOlympics, nil
	case "10", "YoungStars", "Young Stars":
		return GameTypeYoungStars, nil
	case "12", "PWHL Showcase", "PWHLShowcase":
		return GameTypePWHLShowcase, nil
	case "19", "Women's All-Star", "WomensAllStar":
		return GameTypeWomensAllStar, nil
	case "20", "4 Nations Face-Off", "4NationsFaceOff":
		return GameType4Nations, nil
	default:
		return 0, fmt.Errorf("invalid game type: %q", s)
	}
}

// MustGameTypeFromString parses a string into a GameType.
// Panics if the string is not a valid GameType.
func MustGameTypeFromString(s string) GameType {
	g, err := GameTypeFromString(s)
	if err != nil {
		panic(err)
	}
	return g
}

// UnmarshalJSON implements custom JSON unmarshaling for GameType.
// Accepts both integer and string representations.
func (g *GameType) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as int first
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		gameType, err := GameTypeFromInt(i)
		if err != nil {
			return err
		}
		*g = gameType
		return nil
	}

	// Try to unmarshal as string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("game type must be an integer or string: %w", err)
	}

	gameType, err := GameTypeFromString(s)
	if err != nil {
		return err
	}

	*g = gameType
	return nil
}

// MarshalJSON implements custom JSON marshaling for GameType.
// Always marshals as an integer.
func (g GameType) MarshalJSON() ([]byte, error) {
	if !g.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid game type: %d", g)
	}
	return json.Marshal(g.ToInt())
}
