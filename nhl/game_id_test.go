package nhl

import (
	"encoding/json"
	"testing"
)

func TestNewGameID(t *testing.T) {
	id := NewGameID(2023020001)

	if id.AsInt64() != 2023020001 {
		t.Errorf("AsInt64() = %d, want %d", id.AsInt64(), 2023020001)
	}
}

func TestGameID_String(t *testing.T) {
	tests := []struct {
		name     string
		gameID   GameID
		expected string
	}{
		{
			name:     "regular season game",
			gameID:   GameID(2023020001),
			expected: "2023020001",
		},
		{
			name:     "playoff game",
			gameID:   GameID(2023030001),
			expected: "2023030001",
		},
		{
			name:     "preseason game",
			gameID:   GameID(2023010001),
			expected: "2023010001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.gameID.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGameID_JSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		id := GameID(2023020001)
		data, err := json.Marshal(id)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		expected := "2023020001"
		if string(data) != expected {
			t.Errorf("json.Marshal() = %q, want %q", string(data), expected)
		}
	})

	t.Run("unmarshal from integer", func(t *testing.T) {
		data := []byte("2023020001")
		var id GameID
		if err := json.Unmarshal(data, &id); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		if id.AsInt64() != 2023020001 {
			t.Errorf("AsInt64() = %d, want %d", id.AsInt64(), 2023020001)
		}
	})

	t.Run("unmarshal from string", func(t *testing.T) {
		data := []byte(`"2023020001"`)
		var id GameID
		if err := json.Unmarshal(data, &id); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		if id.AsInt64() != 2023020001 {
			t.Errorf("AsInt64() = %d, want %d", id.AsInt64(), 2023020001)
		}
	})

	t.Run("unmarshal invalid string", func(t *testing.T) {
		data := []byte(`"invalid"`)
		var id GameID
		if err := json.Unmarshal(data, &id); err == nil {
			t.Error("json.Unmarshal() should return error for invalid string")
		}
	})

	t.Run("unmarshal invalid type", func(t *testing.T) {
		data := []byte(`true`)
		var id GameID
		if err := json.Unmarshal(data, &id); err == nil {
			t.Error("json.Unmarshal() should return error for boolean")
		}
	})

	t.Run("round trip", func(t *testing.T) {
		original := GameID(2023020001)
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		var decoded GameID
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		if decoded != original {
			t.Errorf("Round trip failed: got %d, want %d", decoded, original)
		}
	})
}

func TestGameID_Season(t *testing.T) {
	tests := []struct {
		name      string
		gameID    GameID
		wantStart int
		wantErr   bool
	}{
		{
			name:      "2023-2024 regular season",
			gameID:    GameID(2023020001),
			wantStart: 2023,
			wantErr:   false,
		},
		{
			name:      "2024-2025 playoffs",
			gameID:    GameID(2024030001),
			wantStart: 2024,
			wantErr:   false,
		},
		{
			name:    "invalid - too short",
			gameID:  GameID(202302001),
			wantErr: true,
		},
		{
			name:    "invalid - too long",
			gameID:  GameID(20230200011),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			season, err := tt.gameID.Season()

			if tt.wantErr {
				if err == nil {
					t.Error("Season() should return error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Season() error = %v", err)
			}

			if season.StartYear() != tt.wantStart {
				t.Errorf("Season().StartYear() = %d, want %d", season.StartYear(), tt.wantStart)
			}
		})
	}
}

func TestGameID_GameType(t *testing.T) {
	tests := []struct {
		name     string
		gameID   GameID
		wantType int
		wantErr  bool
	}{
		{
			name:     "preseason",
			gameID:   GameID(2023010001),
			wantType: 1,
			wantErr:  false,
		},
		{
			name:     "regular season",
			gameID:   GameID(2023020001),
			wantType: 2,
			wantErr:  false,
		},
		{
			name:     "playoffs",
			gameID:   GameID(2023030001),
			wantType: 3,
			wantErr:  false,
		},
		{
			name:     "all-star",
			gameID:   GameID(2023040001),
			wantType: 4,
			wantErr:  false,
		},
		{
			name:    "invalid game ID",
			gameID:  GameID(12345),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameType, err := tt.gameID.GameType()

			if tt.wantErr {
				if err == nil {
					t.Error("GameType() should return error")
				}
				return
			}

			if err != nil {
				t.Fatalf("GameType() error = %v", err)
			}

			if gameType != tt.wantType {
				t.Errorf("GameType() = %d, want %d", gameType, tt.wantType)
			}
		})
	}
}

func TestGameID_GameNumber(t *testing.T) {
	tests := []struct {
		name       string
		gameID     GameID
		wantNumber int
		wantErr    bool
	}{
		{
			name:       "game 1",
			gameID:     GameID(2023020001),
			wantNumber: 1,
			wantErr:    false,
		},
		{
			name:       "game 1230",
			gameID:     GameID(2023021230),
			wantNumber: 1230,
			wantErr:    false,
		},
		{
			name:       "game 0417",
			gameID:     GameID(2023030417),
			wantNumber: 417,
			wantErr:    false,
		},
		{
			name:    "invalid game ID",
			gameID:  GameID(12345),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameNumber, err := tt.gameID.GameNumber()

			if tt.wantErr {
				if err == nil {
					t.Error("GameNumber() should return error")
				}
				return
			}

			if err != nil {
				t.Fatalf("GameNumber() error = %v", err)
			}

			if gameNumber != tt.wantNumber {
				t.Errorf("GameNumber() = %d, want %d", gameNumber, tt.wantNumber)
			}
		})
	}
}

func TestGameID_Validate(t *testing.T) {
	tests := []struct {
		name    string
		gameID  GameID
		wantErr bool
	}{
		{
			name:    "valid preseason game",
			gameID:  GameID(2023010001),
			wantErr: false,
		},
		{
			name:    "valid regular season game",
			gameID:  GameID(2023020001),
			wantErr: false,
		},
		{
			name:    "valid playoff game",
			gameID:  GameID(2023030001),
			wantErr: false,
		},
		{
			name:    "valid all-star game",
			gameID:  GameID(2023040001),
			wantErr: false,
		},
		{
			name:    "invalid - too short",
			gameID:  GameID(202302001),
			wantErr: true,
		},
		{
			name:    "invalid - too long",
			gameID:  GameID(20230200011),
			wantErr: true,
		},
		{
			name:    "invalid - game type 00",
			gameID:  GameID(2023000001),
			wantErr: true,
		},
		{
			name:    "invalid - game type 05",
			gameID:  GameID(2023050001),
			wantErr: true,
		},
		{
			name:    "invalid - game type 99",
			gameID:  GameID(2023990001),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gameID.Validate()

			if tt.wantErr {
				if err == nil {
					t.Error("Validate() should return error")
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			}
		})
	}
}

func TestGameIDFromInt(t *testing.T) {
	id := GameIDFromInt(2023020001)

	if id.AsInt64() != 2023020001 {
		t.Errorf("AsInt64() = %d, want %d", id.AsInt64(), 2023020001)
	}
}

func TestGameIDFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{
			name:     "valid game ID",
			input:    "2023020001",
			expected: 2023020001,
			wantErr:  false,
		},
		{
			name:     "valid game ID with leading zeros",
			input:    "2023020001",
			expected: 2023020001,
			wantErr:  false,
		},
		{
			name:    "invalid - not a number",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "invalid - empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid - float",
			input:   "2023020001.5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GameIDFromString(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("GameIDFromString() should return error")
				}
				return
			}

			if err != nil {
				t.Fatalf("GameIDFromString() error = %v", err)
			}

			if id.AsInt64() != tt.expected {
				t.Errorf("AsInt64() = %d, want %d", id.AsInt64(), tt.expected)
			}
		})
	}
}

func TestMustGameIDFromString(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		id := MustGameIDFromString("2023020001")

		if id.AsInt64() != 2023020001 {
			t.Errorf("AsInt64() = %d, want %d", id.AsInt64(), 2023020001)
		}
	})

	t.Run("invalid input panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustGameIDFromString() should panic on invalid input")
			}
		}()

		MustGameIDFromString("invalid")
	})
}

func TestGameID_CompleteExample(t *testing.T) {
	// Create a game ID for game 1230 of the 2023-2024 regular season
	id := GameID(2023021230)

	// Validate it
	if err := id.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	// Extract components
	season, err := id.Season()
	if err != nil {
		t.Fatalf("Season() error = %v", err)
	}

	if season.StartYear() != 2023 {
		t.Errorf("Season start year = %d, want %d", season.StartYear(), 2023)
	}

	gameType, err := id.GameType()
	if err != nil {
		t.Fatalf("GameType() error = %v", err)
	}

	if gameType != 2 {
		t.Errorf("Game type = %d, want %d (regular season)", gameType, 2)
	}

	gameNumber, err := id.GameNumber()
	if err != nil {
		t.Fatalf("GameNumber() error = %v", err)
	}

	if gameNumber != 1230 {
		t.Errorf("Game number = %d, want %d", gameNumber, 1230)
	}

	// Test string conversion
	if id.String() != "2023021230" {
		t.Errorf("String() = %q, want %q", id.String(), "2023021230")
	}
}

func TestGameID_EdgeCases(t *testing.T) {
	t.Run("minimum valid game ID", func(t *testing.T) {
		id := GameID(1000010000)
		if err := id.Validate(); err != nil {
			t.Errorf("Validate() error = %v, should be valid", err)
		}
	})

	t.Run("maximum valid game ID", func(t *testing.T) {
		id := GameID(9999049999)
		if err := id.Validate(); err != nil {
			t.Errorf("Validate() error = %v, should be valid", err)
		}
	})

	t.Run("zero game ID", func(t *testing.T) {
		id := GameID(0)
		if err := id.Validate(); err == nil {
			t.Error("Validate() should return error for zero")
		}
	})

	t.Run("negative game ID", func(t *testing.T) {
		id := GameID(-2023020001)
		if err := id.Validate(); err == nil {
			t.Error("Validate() should return error for negative ID")
		}
	})
}

func TestGameID_TypeConversions(t *testing.T) {
	t.Run("int64 to GameID to int64", func(t *testing.T) {
		var original int64 = 2023020001
		id := NewGameID(original)
		result := id.AsInt64()

		if result != original {
			t.Errorf("Round trip failed: got %d, want %d", result, original)
		}
	})

	t.Run("string to GameID to string", func(t *testing.T) {
		original := "2023020001"
		id, err := GameIDFromString(original)
		if err != nil {
			t.Fatalf("GameIDFromString() error = %v", err)
		}

		result := id.String()
		if result != original {
			t.Errorf("Round trip failed: got %q, want %q", result, original)
		}
	})
}

// Additional error path tests for GameID.Validate

func TestGameID_Validate_InvalidGameType(t *testing.T) {
	// Create a game ID with an invalid game type (05, which is > 4)
	gameID := GameID(2023050001)
	err := gameID.Validate()
	if err == nil {
		t.Error("Validate() should error on invalid game type")
	}
}
