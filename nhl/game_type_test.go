package nhl

import (
	"encoding/json"
	"testing"
)

func TestGameType_ToInt(t *testing.T) {
	tests := []struct {
		name     string
		gameType GameType
		want     int
	}{
		{"preseason", GameTypePreseason, 1},
		{"regular season", GameTypeRegularSeason, 2},
		{"playoffs", GameTypePlayoffs, 3},
		{"all-star", GameTypeAllStar, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gameType.ToInt(); got != tt.want {
				t.Errorf("GameType.ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameType_String(t *testing.T) {
	tests := []struct {
		name     string
		gameType GameType
		want     string
	}{
		{"preseason", GameTypePreseason, "Preseason"},
		{"regular season", GameTypeRegularSeason, "Regular Season"},
		{"playoffs", GameTypePlayoffs, "Playoffs"},
		{"all-star", GameTypeAllStar, "All-Star"},
		{"unknown", GameType(99), "Unknown(99)"},
		{"zero", GameType(0), "Unknown(0)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gameType.String(); got != tt.want {
				t.Errorf("GameType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		gameType GameType
		want     bool
	}{
		{"preseason valid", GameTypePreseason, true},
		{"regular season valid", GameTypeRegularSeason, true},
		{"playoffs valid", GameTypePlayoffs, true},
		{"all-star valid", GameTypeAllStar, true},
		{"zero invalid", GameType(0), false},
		{"negative invalid", GameType(-1), false},
		{"too high invalid", GameType(5), false},
		{"unknown invalid", GameType(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gameType.IsValid(); got != tt.want {
				t.Errorf("GameType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameTypeFromInt(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    GameType
		wantErr bool
	}{
		{"preseason", 1, GameTypePreseason, false},
		{"regular season", 2, GameTypeRegularSeason, false},
		{"playoffs", 3, GameTypePlayoffs, false},
		{"all-star", 4, GameTypeAllStar, false},
		{"zero error", 0, GameType(0), true},
		{"negative error", -1, GameType(0), true},
		{"too high error", 5, GameType(0), true},
		{"unknown error", 99, GameType(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GameTypeFromInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameTypeFromInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GameTypeFromInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustGameTypeFromInt(t *testing.T) {
	t.Run("valid type", func(t *testing.T) {
		got := MustGameTypeFromInt(2)
		if got != GameTypeRegularSeason {
			t.Errorf("MustGameTypeFromInt() = %v, want %v", got, GameTypeRegularSeason)
		}
	})

	t.Run("invalid type panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustGameTypeFromInt() did not panic")
			}
		}()
		MustGameTypeFromInt(99)
	})
}

func TestGameTypeFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    GameType
		wantErr bool
	}{
		{"numeric preseason", "1", GameTypePreseason, false},
		{"text preseason", "Preseason", GameTypePreseason, false},
		{"numeric regular season", "2", GameTypeRegularSeason, false},
		{"text regular season", "Regular Season", GameTypeRegularSeason, false},
		{"text regular season no space", "RegularSeason", GameTypeRegularSeason, false},
		{"numeric playoffs", "3", GameTypePlayoffs, false},
		{"text playoffs", "Playoffs", GameTypePlayoffs, false},
		{"numeric all-star", "4", GameTypeAllStar, false},
		{"text all-star", "All-Star", GameTypeAllStar, false},
		{"text all-star no hyphen", "AllStar", GameTypeAllStar, false},
		{"empty error", "", GameType(0), true},
		{"unknown numeric error", "99", GameType(0), true},
		{"unknown text error", "Unknown", GameType(0), true},
		{"lowercase error", "preseason", GameType(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GameTypeFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameTypeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GameTypeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustGameTypeFromString(t *testing.T) {
	t.Run("valid type", func(t *testing.T) {
		got := MustGameTypeFromString("Playoffs")
		if got != GameTypePlayoffs {
			t.Errorf("MustGameTypeFromString() = %v, want %v", got, GameTypePlayoffs)
		}
	})

	t.Run("invalid type panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustGameTypeFromString() did not panic")
			}
		}()
		MustGameTypeFromString("Invalid")
	})
}

func TestGameType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    GameType
		wantErr bool
	}{
		{"int preseason", `1`, GameTypePreseason, false},
		{"int regular season", `2`, GameTypeRegularSeason, false},
		{"int playoffs", `3`, GameTypePlayoffs, false},
		{"int all-star", `4`, GameTypeAllStar, false},
		{"string preseason", `"Preseason"`, GameTypePreseason, false},
		{"string regular season", `"Regular Season"`, GameTypeRegularSeason, false},
		{"string playoffs", `"Playoffs"`, GameTypePlayoffs, false},
		{"string all-star", `"All-Star"`, GameTypeAllStar, false},
		{"numeric string preseason", `"1"`, GameTypePreseason, false},
		{"numeric string regular season", `"2"`, GameTypeRegularSeason, false},
		{"invalid int", `99`, GameType(0), true},
		{"invalid string", `"Unknown"`, GameType(0), true},
		{"null", `null`, GameType(0), true},
		{"boolean", `true`, GameType(0), true},
		{"array", `[1]`, GameType(0), true},
		{"object", `{"type":1}`, GameType(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got GameType
			err := json.Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameType.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GameType.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		gameType GameType
		want     string
		wantErr  bool
	}{
		{"preseason", GameTypePreseason, `1`, false},
		{"regular season", GameTypeRegularSeason, `2`, false},
		{"playoffs", GameTypePlayoffs, `3`, false},
		{"all-star", GameTypeAllStar, `4`, false},
		{"invalid type", GameType(99), "", true},
		{"zero type", GameType(0), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.gameType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.want {
				t.Errorf("GameType.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestGameType_JSONRoundTrip(t *testing.T) {
	gameTypes := []GameType{
		GameTypePreseason,
		GameTypeRegularSeason,
		GameTypePlayoffs,
		GameTypeAllStar,
	}

	for _, gameType := range gameTypes {
		t.Run(gameType.String(), func(t *testing.T) {
			data, err := json.Marshal(gameType)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			var got GameType
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if got != gameType {
				t.Errorf("round trip failed: got %v, want %v", got, gameType)
			}
		})
	}
}

func TestGameType_UnmarshalJSON_IntAndStringEquivalence(t *testing.T) {
	tests := []struct {
		name     string
		intJSON  string
		strJSON  string
		expected GameType
	}{
		{"preseason", `1`, `"Preseason"`, GameTypePreseason},
		{"regular season", `2`, `"Regular Season"`, GameTypeRegularSeason},
		{"playoffs", `3`, `"Playoffs"`, GameTypePlayoffs},
		{"all-star", `4`, `"All-Star"`, GameTypeAllStar},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fromInt GameType
			if err := json.Unmarshal([]byte(tt.intJSON), &fromInt); err != nil {
				t.Fatalf("unmarshal int error = %v", err)
			}

			var fromStr GameType
			if err := json.Unmarshal([]byte(tt.strJSON), &fromStr); err != nil {
				t.Fatalf("unmarshal string error = %v", err)
			}

			if fromInt != tt.expected {
				t.Errorf("int unmarshal: got %v, want %v", fromInt, tt.expected)
			}

			if fromStr != tt.expected {
				t.Errorf("string unmarshal: got %v, want %v", fromStr, tt.expected)
			}

			if fromInt != fromStr {
				t.Errorf("int and string unmarshaled to different values: %v vs %v", fromInt, fromStr)
			}
		})
	}
}
