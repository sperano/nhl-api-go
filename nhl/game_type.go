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
	// GameTypeWorldCup represents a World Cup of Hockey game.
	GameTypeWorldCup GameType = 6
	// GameTypeWorldCup2004 represents World Cup 2004 (NHL API uses both 6 and 7 for World Cup).
	GameTypeWorldCup2004 GameType = 7
	// GameTypeWorldCupPreTournament represents a World Cup pre-tournament game.
	GameTypeWorldCupPreTournament GameType = 8
	// GameTypeOlympics represents an Olympic tournament game.
	GameTypeOlympics GameType = 9
	// GameTypeYoungStars represents a YoungStars game (rookies/sophomores) during All-Star Weekend.
	GameTypeYoungStars GameType = 10
	// GameTypePWHLShowcase represents a PWHL 3-on-3 showcase game during All-Star Weekend.
	GameTypePWHLShowcase GameType = 12
	// GameTypeLockoutLost represents a game lost due to a lockout.
	GameTypeLockoutLost GameType = 13
	// GameTypeCanadaCup represents a Canada Cup game.
	GameTypeCanadaCup GameType = 14
	// GameTypeExhibitionOverseas represents an exhibition game played overseas.
	GameTypeExhibitionOverseas GameType = 18
	// GameTypeWomensAllStar represents a Women's All-Star game.
	GameTypeWomensAllStar GameType = 19
	// GameType4Nations represents a 4 Nations Face-Off game.
	GameType4Nations GameType = 20
)

// Int returns the integer representation of the GameType.
func (g GameType) Int() int {
	return int(g)
}

// Label returns the snake_case label for the GameType, suitable for use as a
// PostgreSQL enum value or a normalized identifier.
func (g GameType) Label() string {
	switch g {
	case GameTypePreseason:
		return "preseason"
	case GameTypeRegularSeason:
		return "regular_season"
	case GameTypePlayoffs:
		return "playoffs"
	case GameTypeAllStar:
		return "all_star"
	case GameTypeWorldCup:
		return "world_cup"
	case GameTypeWorldCup2004:
		return "world_cup_2004"
	case GameTypeWorldCupPreTournament:
		return "world_cup_pre_tournament"
	case GameTypeOlympics:
		return "olympics"
	case GameTypeYoungStars:
		return "young_stars"
	case GameTypePWHLShowcase:
		return "pwhl_showcase"
	case GameTypeLockoutLost:
		return "lockout_lost"
	case GameTypeCanadaCup:
		return "canada_cup"
	case GameTypeExhibitionOverseas:
		return "exhibition_overseas"
	case GameTypeWomensAllStar:
		return "womens_all_star"
	case GameType4Nations:
		return "four_nations"
	default:
		return fmt.Sprintf("unknown_%d", g)
	}
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
	case GameTypeWorldCup:
		return "World Cup"
	case GameTypeWorldCup2004:
		return "World Cup 2004"
	case GameTypeWorldCupPreTournament:
		return "World Cup Pre-Tournament"
	case GameTypeOlympics:
		return "Olympics"
	case GameTypeYoungStars:
		return "YoungStars"
	case GameTypePWHLShowcase:
		return "PWHL Showcase"
	case GameTypeLockoutLost:
		return "Lockout Lost"
	case GameTypeCanadaCup:
		return "Canada Cup"
	case GameTypeExhibitionOverseas:
		return "Exhibition Overseas"
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
	case GameTypePreseason, GameTypeRegularSeason, GameTypePlayoffs, GameTypeAllStar, GameTypeWorldCup, GameTypeWorldCup2004, GameTypeWorldCupPreTournament, GameTypeOlympics, GameTypeYoungStars, GameTypePWHLShowcase, GameTypeLockoutLost, GameTypeCanadaCup, GameTypeExhibitionOverseas, GameTypeWomensAllStar, GameType4Nations:
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
	case "1", "Preseason", "preseason":
		return GameTypePreseason, nil
	case "2", "Regular Season", "RegularSeason", "regular_season":
		return GameTypeRegularSeason, nil
	case "3", "Playoffs", "playoffs":
		return GameTypePlayoffs, nil
	case "4", "All-Star", "AllStar", "all_star":
		return GameTypeAllStar, nil
	case "6", "World Cup", "WorldCup", "world_cup":
		return GameTypeWorldCup, nil
	case "7", "World Cup 2004", "WorldCup2004", "world_cup_2004":
		return GameTypeWorldCup2004, nil
	case "8", "World Cup Pre-Tournament", "WorldCupPreTournament", "world_cup_pre_tournament":
		return GameTypeWorldCupPreTournament, nil
	case "9", "Olympics", "olympics":
		return GameTypeOlympics, nil
	case "10", "YoungStars", "Young Stars", "young_stars":
		return GameTypeYoungStars, nil
	case "12", "PWHL Showcase", "PWHLShowcase", "pwhl_showcase":
		return GameTypePWHLShowcase, nil
	case "13", "Lockout Lost", "LockoutLost", "lockout_lost":
		return GameTypeLockoutLost, nil
	case "14", "Canada Cup", "CanadaCup", "canada_cup":
		return GameTypeCanadaCup, nil
	case "18", "Exhibition Overseas", "ExhibitionOverseas", "exhibition_overseas":
		return GameTypeExhibitionOverseas, nil
	case "19", "Women's All-Star", "WomensAllStar", "womens_all_star":
		return GameTypeWomensAllStar, nil
	case "20", "4 Nations Face-Off", "4NationsFaceOff", "four_nations":
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
	return json.Marshal(g.Int())
}
