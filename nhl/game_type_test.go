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
			if got := tt.gameType.Int(); got != tt.want {
				t.Errorf("GameType.Int() = %v, want %v", got, tt.want)
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
		{"world cup", GameTypeWorldCup, "World Cup"},
		{"world cup 2004", GameTypeWorldCup2004, "World Cup 2004"},
		{"world cup pre-tournament", GameTypeWorldCupPreTournament, "World Cup Pre-Tournament"},
		{"olympics", GameTypeOlympics, "Olympics"},
		{"young stars", GameTypeYoungStars, "YoungStars"},
		{"pwhl showcase", GameTypePWHLShowcase, "PWHL Showcase"},
		{"lockout lost", GameTypeLockoutLost, "Lockout Lost"},
		{"canada cup", GameTypeCanadaCup, "Canada Cup"},
		{"exhibition overseas", GameTypeExhibitionOverseas, "Exhibition Overseas"},
		{"womens all-star", GameTypeWomensAllStar, "Women's All-Star"},
		{"4 nations", GameType4Nations, "4 Nations Face-Off"},
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
		{"world cup valid", GameTypeWorldCup, true},
		{"world cup 2004 valid", GameTypeWorldCup2004, true},
		{"world cup pre-tournament valid", GameTypeWorldCupPreTournament, true},
		{"olympics valid", GameTypeOlympics, true},
		{"young stars valid", GameTypeYoungStars, true},
		{"pwhl showcase valid", GameTypePWHLShowcase, true},
		{"lockout lost valid", GameTypeLockoutLost, true},
		{"canada cup valid", GameTypeCanadaCup, true},
		{"exhibition overseas valid", GameTypeExhibitionOverseas, true},
		{"womens all-star valid", GameTypeWomensAllStar, true},
		{"4 nations valid", GameType4Nations, true},
		{"zero invalid", GameType(0), false},
		{"negative invalid", GameType(-1), false},
		{"too high invalid", GameType(5), false},
		{"unknown invalid", GameType(99), false},
		// Pin the iota gaps: these integer values fall between defined
		// constants and must be rejected.
		{"between 4 and 6 invalid", GameType(11), false},
		{"between 14 and 18 invalid", GameType(15), false},
		{"above 20 invalid", GameType(21), false},
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
		{"world cup", 6, GameTypeWorldCup, false},
		{"world cup 2004", 7, GameTypeWorldCup2004, false},
		{"world cup pre-tournament", 8, GameTypeWorldCupPreTournament, false},
		{"olympics", 9, GameTypeOlympics, false},
		{"young stars", 10, GameTypeYoungStars, false},
		{"pwhl showcase", 12, GameTypePWHLShowcase, false},
		{"lockout lost", 13, GameTypeLockoutLost, false},
		{"canada cup", 14, GameTypeCanadaCup, false},
		{"exhibition overseas", 18, GameTypeExhibitionOverseas, false},
		{"womens all-star", 19, GameTypeWomensAllStar, false},
		{"4 nations", 20, GameType4Nations, false},
		{"zero error", 0, GameType(0), true},
		{"negative error", -1, GameType(0), true},
		{"too high error", 5, GameType(0), true},
		{"iota gap 11 error", 11, GameType(0), true},
		{"iota gap 15 error", 15, GameType(0), true},
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
		{"numeric world cup pre-tournament", "8", GameTypeWorldCupPreTournament, false},
		{"display world cup pre-tournament", "World Cup Pre-Tournament", GameTypeWorldCupPreTournament, false},
		{"camel world cup pre-tournament", "WorldCupPreTournament", GameTypeWorldCupPreTournament, false},
		{"snake world cup pre-tournament", "world_cup_pre_tournament", GameTypeWorldCupPreTournament, false},
		{"numeric lockout lost", "13", GameTypeLockoutLost, false},
		{"display lockout lost", "Lockout Lost", GameTypeLockoutLost, false},
		{"camel lockout lost", "LockoutLost", GameTypeLockoutLost, false},
		{"snake lockout lost", "lockout_lost", GameTypeLockoutLost, false},
		{"numeric canada cup", "14", GameTypeCanadaCup, false},
		{"display canada cup", "Canada Cup", GameTypeCanadaCup, false},
		{"camel canada cup", "CanadaCup", GameTypeCanadaCup, false},
		{"snake canada cup", "canada_cup", GameTypeCanadaCup, false},
		{"numeric exhibition overseas", "18", GameTypeExhibitionOverseas, false},
		{"display exhibition overseas", "Exhibition Overseas", GameTypeExhibitionOverseas, false},
		{"camel exhibition overseas", "ExhibitionOverseas", GameTypeExhibitionOverseas, false},
		{"snake exhibition overseas", "exhibition_overseas", GameTypeExhibitionOverseas, false},
		{"numeric world cup", "6", GameTypeWorldCup, false},
		{"display world cup", "World Cup", GameTypeWorldCup, false},
		{"camel world cup", "WorldCup", GameTypeWorldCup, false},
		{"snake world cup", "world_cup", GameTypeWorldCup, false},
		{"numeric world cup 2004", "7", GameTypeWorldCup2004, false},
		{"display world cup 2004", "World Cup 2004", GameTypeWorldCup2004, false},
		{"camel world cup 2004", "WorldCup2004", GameTypeWorldCup2004, false},
		{"snake world cup 2004", "world_cup_2004", GameTypeWorldCup2004, false},
		{"numeric olympics", "9", GameTypeOlympics, false},
		{"display olympics", "Olympics", GameTypeOlympics, false},
		{"snake olympics", "olympics", GameTypeOlympics, false},
		{"numeric young stars", "10", GameTypeYoungStars, false},
		{"display young stars no space", "YoungStars", GameTypeYoungStars, false},
		{"display young stars spaced", "Young Stars", GameTypeYoungStars, false},
		{"snake young stars", "young_stars", GameTypeYoungStars, false},
		{"numeric pwhl showcase", "12", GameTypePWHLShowcase, false},
		{"display pwhl showcase", "PWHL Showcase", GameTypePWHLShowcase, false},
		{"camel pwhl showcase", "PWHLShowcase", GameTypePWHLShowcase, false},
		{"snake pwhl showcase", "pwhl_showcase", GameTypePWHLShowcase, false},
		{"numeric womens all-star", "19", GameTypeWomensAllStar, false},
		{"display womens all-star", "Women's All-Star", GameTypeWomensAllStar, false},
		{"camel womens all-star", "WomensAllStar", GameTypeWomensAllStar, false},
		{"snake womens all-star", "womens_all_star", GameTypeWomensAllStar, false},
		{"numeric 4 nations", "20", GameType4Nations, false},
		{"display 4 nations", "4 Nations Face-Off", GameType4Nations, false},
		{"camel 4 nations", "4NationsFaceOff", GameType4Nations, false},
		{"snake 4 nations", "four_nations", GameType4Nations, false},
		{"empty error", "", GameType(0), true},
		{"unknown numeric error", "99", GameType(0), true},
		{"unknown text error", "Unknown", GameType(0), true},
		{"snake_case preseason", "preseason", GameTypePreseason, false},
		{"snake_case regular_season", "regular_season", GameTypeRegularSeason, false},
		{"snake_case playoffs", "playoffs", GameTypePlayoffs, false},
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
		{"int world cup pre-tournament", `8`, GameTypeWorldCupPreTournament, false},
		{"int lockout lost", `13`, GameTypeLockoutLost, false},
		{"int canada cup", `14`, GameTypeCanadaCup, false},
		{"int exhibition overseas", `18`, GameTypeExhibitionOverseas, false},
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
		{"world cup pre-tournament", GameTypeWorldCupPreTournament, `8`, false},
		{"lockout lost", GameTypeLockoutLost, `13`, false},
		{"canada cup", GameTypeCanadaCup, `14`, false},
		{"exhibition overseas", GameTypeExhibitionOverseas, `18`, false},
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
		GameTypeWorldCupPreTournament,
		GameTypeLockoutLost,
		GameTypeCanadaCup,
		GameTypeExhibitionOverseas,
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
