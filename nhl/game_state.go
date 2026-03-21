package nhl

// Custom domain methods for GameState. The type declaration, constants,
// and boilerplate methods are generated in enums_generated.go.

// HasStarted returns true if the game has started (is live, final, or off).
func (v GameState) HasStarted() bool {
	switch v {
	case GameStateLive, GameStateFinal, GameStateOff, GameStateCritical:
		return true
	default:
		return false
	}
}

// IsFinal returns true if the game is completed (final or off).
func (v GameState) IsFinal() bool {
	return v == GameStateFinal || v == GameStateOff
}

// IsLive returns true if the game is currently in progress.
func (v GameState) IsLive() bool {
	return v == GameStateLive || v == GameStateCritical
}

// IsScheduled returns true if the game is scheduled (future or pre-game).
func (v GameState) IsScheduled() bool {
	return v == GameStateFuture || v == GameStatePreGame
}
