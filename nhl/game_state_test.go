package nhl

import (
	"encoding/json"
	"testing"
)

func TestGameState_String(t *testing.T) {
	tests := []struct {
		name  string
		state GameState
		want  string
	}{
		{"future", GameStateFuture, "FUT"},
		{"pre-game", GameStatePreGame, "PRE"},
		{"live", GameStateLive, "LIVE"},
		{"final", GameStateFinal, "FINAL"},
		{"off", GameStateOff, "OFF"},
		{"postponed", GameStatePostponed, "PPD"},
		{"suspended", GameStateSuspended, "SUSP"},
		{"critical", GameStateCritical, "CRIT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("GameState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameState_HasStarted(t *testing.T) {
	tests := []struct {
		name  string
		state GameState
		want  bool
	}{
		{"future not started", GameStateFuture, false},
		{"pre-game not started", GameStatePreGame, false},
		{"live has started", GameStateLive, true},
		{"final has started", GameStateFinal, true},
		{"off has started", GameStateOff, true},
		{"postponed not started", GameStatePostponed, false},
		{"suspended not started", GameStateSuspended, false},
		{"critical has started", GameStateCritical, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.HasStarted(); got != tt.want {
				t.Errorf("GameState.HasStarted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameState_IsFinal(t *testing.T) {
	tests := []struct {
		name  string
		state GameState
		want  bool
	}{
		{"future not final", GameStateFuture, false},
		{"pre-game not final", GameStatePreGame, false},
		{"live not final", GameStateLive, false},
		{"final is final", GameStateFinal, true},
		{"off is final", GameStateOff, true},
		{"postponed not final", GameStatePostponed, false},
		{"suspended not final", GameStateSuspended, false},
		{"critical not final", GameStateCritical, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsFinal(); got != tt.want {
				t.Errorf("GameState.IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameState_IsLive(t *testing.T) {
	tests := []struct {
		name  string
		state GameState
		want  bool
	}{
		{"future not live", GameStateFuture, false},
		{"pre-game not live", GameStatePreGame, false},
		{"live is live", GameStateLive, true},
		{"final not live", GameStateFinal, false},
		{"off not live", GameStateOff, false},
		{"postponed not live", GameStatePostponed, false},
		{"suspended not live", GameStateSuspended, false},
		{"critical is live", GameStateCritical, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsLive(); got != tt.want {
				t.Errorf("GameState.IsLive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameState_IsScheduled(t *testing.T) {
	tests := []struct {
		name  string
		state GameState
		want  bool
	}{
		{"future is scheduled", GameStateFuture, true},
		{"pre-game is scheduled", GameStatePreGame, true},
		{"live not scheduled", GameStateLive, false},
		{"final not scheduled", GameStateFinal, false},
		{"off not scheduled", GameStateOff, false},
		{"postponed not scheduled", GameStatePostponed, false},
		{"suspended not scheduled", GameStateSuspended, false},
		{"critical not scheduled", GameStateCritical, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsScheduled(); got != tt.want {
				t.Errorf("GameState.IsScheduled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameState_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		state GameState
		want  bool
	}{
		{"future valid", GameStateFuture, true},
		{"pre-game valid", GameStatePreGame, true},
		{"live valid", GameStateLive, true},
		{"final valid", GameStateFinal, true},
		{"off valid", GameStateOff, true},
		{"postponed valid", GameStatePostponed, true},
		{"suspended valid", GameStateSuspended, true},
		{"critical valid", GameStateCritical, true},
		{"empty invalid", GameState(""), false},
		{"unknown invalid", GameState("UNKNOWN"), false},
		{"lowercase invalid", GameState("live"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsValid(); got != tt.want {
				t.Errorf("GameState.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameStateFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    GameState
		wantErr bool
	}{
		{"future", "FUT", GameStateFuture, false},
		{"pre-game", "PRE", GameStatePreGame, false},
		{"live", "LIVE", GameStateLive, false},
		{"final", "FINAL", GameStateFinal, false},
		{"off", "OFF", GameStateOff, false},
		{"postponed", "PPD", GameStatePostponed, false},
		{"suspended", "SUSP", GameStateSuspended, false},
		{"critical", "CRIT", GameStateCritical, false},
		{"empty error", "", GameState(""), true},
		{"unknown error", "UNKNOWN", GameState(""), true},
		{"lowercase error", "live", GameState(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GameStateFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameStateFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GameStateFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustGameStateFromString(t *testing.T) {
	t.Run("valid state", func(t *testing.T) {
		got := MustGameStateFromString("LIVE")
		if got != GameStateLive {
			t.Errorf("MustGameStateFromString() = %v, want %v", got, GameStateLive)
		}
	})

	t.Run("invalid state panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustGameStateFromString() did not panic")
			}
		}()
		MustGameStateFromString("INVALID")
	})
}

func TestGameState_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    GameState
		wantErr bool
	}{
		{"future", `"FUT"`, GameStateFuture, false},
		{"pre-game", `"PRE"`, GameStatePreGame, false},
		{"live", `"LIVE"`, GameStateLive, false},
		{"final", `"FINAL"`, GameStateFinal, false},
		{"off", `"OFF"`, GameStateOff, false},
		{"postponed", `"PPD"`, GameStatePostponed, false},
		{"suspended", `"SUSP"`, GameStateSuspended, false},
		{"critical", `"CRIT"`, GameStateCritical, false},
		{"invalid state", `"INVALID"`, GameState(""), true},
		{"invalid json", `invalid`, GameState(""), true},
		{"null", `null`, GameState(""), true},
		{"number", `123`, GameState(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got GameState
			err := json.Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameState.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GameState.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameState_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		state   GameState
		want    string
		wantErr bool
	}{
		{"future", GameStateFuture, `"FUT"`, false},
		{"pre-game", GameStatePreGame, `"PRE"`, false},
		{"live", GameStateLive, `"LIVE"`, false},
		{"final", GameStateFinal, `"FINAL"`, false},
		{"off", GameStateOff, `"OFF"`, false},
		{"postponed", GameStatePostponed, `"PPD"`, false},
		{"suspended", GameStateSuspended, `"SUSP"`, false},
		{"critical", GameStateCritical, `"CRIT"`, false},
		{"invalid state", GameState("INVALID"), "", true},
		{"empty state", GameState(""), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameState.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("GameState.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestGameState_JSONRoundTrip(t *testing.T) {
	states := []GameState{
		GameStateFuture,
		GameStatePreGame,
		GameStateLive,
		GameStateFinal,
		GameStateOff,
		GameStatePostponed,
		GameStateSuspended,
		GameStateCritical,
	}

	for _, state := range states {
		t.Run(state.String(), func(t *testing.T) {
			data, err := json.Marshal(state)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			var got GameState
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if got != state {
				t.Errorf("round trip failed: got %v, want %v", got, state)
			}
		})
	}
}
