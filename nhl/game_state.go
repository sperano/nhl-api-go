package nhl

import (
	"encoding/json"
	"fmt"
)

// GameState represents the current state of an NHL game.
type GameState string

const (
	// GameStateFuture represents a game that has not yet started.
	GameStateFuture GameState = "FUT"
	// GameStatePreGame represents a game in pre-game state.
	GameStatePreGame GameState = "PRE"
	// GameStateLive represents a game currently in progress.
	GameStateLive GameState = "LIVE"
	// GameStateFinal represents a completed game.
	GameStateFinal GameState = "FINAL"
	// GameStateOff represents a game that is off/completed (alternative to FINAL).
	GameStateOff GameState = "OFF"
	// GameStatePostponed represents a postponed game.
	GameStatePostponed GameState = "PPD"
	// GameStateSuspended represents a suspended game.
	GameStateSuspended GameState = "SUSP"
	// GameStateCritical represents a game in critical state.
	GameStateCritical GameState = "CRIT"
)

// String returns the string representation of the GameState.
func (g GameState) String() string {
	return string(g)
}

// HasStarted returns true if the game has started (is live, final, or off).
func (g GameState) HasStarted() bool {
	switch g {
	case GameStateLive, GameStateFinal, GameStateOff, GameStateCritical:
		return true
	default:
		return false
	}
}

// IsFinal returns true if the game is completed (final or off).
func (g GameState) IsFinal() bool {
	return g == GameStateFinal || g == GameStateOff
}

// IsLive returns true if the game is currently in progress.
func (g GameState) IsLive() bool {
	return g == GameStateLive || g == GameStateCritical
}

// IsScheduled returns true if the game is scheduled (future or pre-game).
func (g GameState) IsScheduled() bool {
	return g == GameStateFuture || g == GameStatePreGame
}

// IsValid returns true if the GameState is one of the known valid states.
func (g GameState) IsValid() bool {
	switch g {
	case GameStateFuture, GameStatePreGame, GameStateLive, GameStateFinal,
		GameStateOff, GameStatePostponed, GameStateSuspended, GameStateCritical:
		return true
	default:
		return false
	}
}

// GameStateFromString parses a string into a GameState.
// Returns an error if the string is not a valid GameState.
func GameStateFromString(s string) (GameState, error) {
	g := GameState(s)
	if !g.IsValid() {
		return "", fmt.Errorf("invalid game state: %q", s)
	}
	return g, nil
}

// MustGameStateFromString parses a string into a GameState.
// Panics if the string is not a valid GameState.
func MustGameStateFromString(s string) GameState {
	g, err := GameStateFromString(s)
	if err != nil {
		panic(err)
	}
	return g
}

// UnmarshalJSON implements custom JSON unmarshaling for GameState.
func (g *GameState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	state, err := GameStateFromString(s)
	if err != nil {
		return err
	}

	*g = state
	return nil
}

// MarshalJSON implements custom JSON marshaling for GameState.
func (g GameState) MarshalJSON() ([]byte, error) {
	if !g.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid game state: %q", string(g))
	}
	return json.Marshal(string(g))
}
