package nhl

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestClubSkaterStatsDeserialization(t *testing.T) {
	jsonData := `{
		"playerId": 8475233,
		"headshot": "https://assets.nhle.com/mugs/nhl/20242025/MTL/8475233.png",
		"firstName": {
			"default": "David"
		},
		"lastName": {
			"default": "Savard"
		},
		"positionCode": "D",
		"gamesPlayed": 75,
		"goals": 1,
		"assists": 14,
		"points": 15,
		"plusMinus": -8,
		"penaltyMinutes": 36,
		"powerPlayGoals": 0,
		"shorthandedGoals": 0,
		"gameWinningGoals": 0,
		"overtimeGoals": 0,
		"shots": 48,
		"shootingPctg": 0.020833,
		"avgTimeOnIcePerGame": 995.36,
		"avgShiftsPerGame": 19.84,
		"faceoffWinPctg": 0.0
	}`

	var stats ClubSkaterStats
	if err := json.Unmarshal([]byte(jsonData), &stats); err != nil {
		t.Fatalf("Failed to unmarshal ClubSkaterStats: %v", err)
	}

	if stats.PlayerID != PlayerID(8475233) {
		t.Errorf("Expected PlayerID 8475233, got %d", stats.PlayerID)
	}
	if stats.FirstName.Default != "David" {
		t.Errorf("Expected FirstName 'David', got '%s'", stats.FirstName.Default)
	}
	if stats.LastName.Default != "Savard" {
		t.Errorf("Expected LastName 'Savard', got '%s'", stats.LastName.Default)
	}
	if stats.Position != PositionDefense {
		t.Errorf("Expected Position Defense, got %v", stats.Position)
	}
	if stats.GamesPlayed != 75 {
		t.Errorf("Expected GamesPlayed 75, got %d", stats.GamesPlayed)
	}
	if stats.Goals != 1 {
		t.Errorf("Expected Goals 1, got %d", stats.Goals)
	}
	if stats.Assists != 14 {
		t.Errorf("Expected Assists 14, got %d", stats.Assists)
	}
	if stats.Points != 15 {
		t.Errorf("Expected Points 15, got %d", stats.Points)
	}
	if stats.PlusMinus != -8 {
		t.Errorf("Expected PlusMinus -8, got %d", stats.PlusMinus)
	}
	if stats.Shots != 48 {
		t.Errorf("Expected Shots 48, got %d", stats.Shots)
	}
}

func TestClubGoalieStatsDeserialization(t *testing.T) {
	jsonData := `{
		"playerId": 8478470,
		"headshot": "https://assets.nhle.com/mugs/nhl/20242025/MTL/8478470.png",
		"firstName": {
			"default": "Sam"
		},
		"lastName": {
			"default": "Montembeault"
		},
		"gamesPlayed": 62,
		"gamesStarted": 60,
		"wins": 31,
		"losses": 24,
		"overtimeLosses": 7,
		"goalsAgainstAverage": 2.818349,
		"savePercentage": 0.901669,
		"shotsAgainst": 1678,
		"saves": 1513,
		"goalsAgainst": 166,
		"shutouts": 4,
		"goals": 0,
		"assists": 1,
		"points": 1,
		"penaltyMinutes": 0,
		"timeOnIce": 212039
	}`

	var stats ClubGoalieStats
	if err := json.Unmarshal([]byte(jsonData), &stats); err != nil {
		t.Fatalf("Failed to unmarshal ClubGoalieStats: %v", err)
	}

	if stats.PlayerID != PlayerID(8478470) {
		t.Errorf("Expected PlayerID 8478470, got %d", stats.PlayerID)
	}
	if stats.FirstName.Default != "Sam" {
		t.Errorf("Expected FirstName 'Sam', got '%s'", stats.FirstName.Default)
	}
	if stats.LastName.Default != "Montembeault" {
		t.Errorf("Expected LastName 'Montembeault', got '%s'", stats.LastName.Default)
	}
	if stats.GamesPlayed != 62 {
		t.Errorf("Expected GamesPlayed 62, got %d", stats.GamesPlayed)
	}
	if stats.Wins != 31 {
		t.Errorf("Expected Wins 31, got %d", stats.Wins)
	}
	if stats.Losses != 24 {
		t.Errorf("Expected Losses 24, got %d", stats.Losses)
	}
	if stats.OvertimeLosses != 7 {
		t.Errorf("Expected OvertimeLosses 7, got %d", stats.OvertimeLosses)
	}
	if stats.Shutouts != 4 {
		t.Errorf("Expected Shutouts 4, got %d", stats.Shutouts)
	}
}

func TestClubStatsDeserialization(t *testing.T) {
	jsonData := `{
		"season": "20242025",
		"gameType": 2,
		"skaters": [
			{
				"playerId": 8475233,
				"headshot": "https://assets.nhle.com/mugs/nhl/20242025/MTL/8475233.png",
				"firstName": {"default": "David"},
				"lastName": {"default": "Savard"},
				"positionCode": "D",
				"gamesPlayed": 75,
				"goals": 1,
				"assists": 14,
				"points": 15,
				"plusMinus": -8,
				"penaltyMinutes": 36,
				"powerPlayGoals": 0,
				"shorthandedGoals": 0,
				"gameWinningGoals": 0,
				"overtimeGoals": 0,
				"shots": 48,
				"shootingPctg": 0.020833,
				"avgTimeOnIcePerGame": 995.36,
				"avgShiftsPerGame": 19.84,
				"faceoffWinPctg": 0.0
			}
		],
		"goalies": [
			{
				"playerId": 8478470,
				"headshot": "https://assets.nhle.com/mugs/nhl/20242025/MTL/8478470.png",
				"firstName": {"default": "Sam"},
				"lastName": {"default": "Montembeault"},
				"gamesPlayed": 62,
				"gamesStarted": 60,
				"wins": 31,
				"losses": 24,
				"overtimeLosses": 7,
				"goalsAgainstAverage": 2.818349,
				"savePercentage": 0.901669,
				"shotsAgainst": 1678,
				"saves": 1513,
				"goalsAgainst": 166,
				"shutouts": 4,
				"goals": 0,
				"assists": 1,
				"points": 1,
				"penaltyMinutes": 0,
				"timeOnIce": 212039
			}
		]
	}`

	var stats ClubStats
	if err := json.Unmarshal([]byte(jsonData), &stats); err != nil {
		t.Fatalf("Failed to unmarshal ClubStats: %v", err)
	}

	if stats.Season != "20242025" {
		t.Errorf("Expected Season '20242025', got '%s'", stats.Season)
	}
	if stats.GameType != GameTypeRegularSeason {
		t.Errorf("Expected GameType RegularSeason, got %v", stats.GameType)
	}
	if len(stats.Skaters) != 1 {
		t.Errorf("Expected 1 skater, got %d", len(stats.Skaters))
	}
	if len(stats.Goalies) != 1 {
		t.Errorf("Expected 1 goalie, got %d", len(stats.Goalies))
	}
}

func TestSeasonGameTypesDeserialization(t *testing.T) {
	jsonData := `{
		"season": 20242025,
		"gameTypes": [2, 3]
	}`

	var season SeasonGameTypes
	if err := json.Unmarshal([]byte(jsonData), &season); err != nil {
		t.Fatalf("Failed to unmarshal SeasonGameTypes: %v", err)
	}

	if season.Season != NewSeason(2024) {
		t.Errorf("Expected Season 20242025, got %s", season.Season)
	}
	if len(season.GameTypes) != 2 {
		t.Fatalf("Expected 2 game types, got %d", len(season.GameTypes))
	}
	if season.GameTypes[0] != GameTypeRegularSeason {
		t.Errorf("Expected first GameType RegularSeason, got %v", season.GameTypes[0])
	}
	if season.GameTypes[1] != GameTypePlayoffs {
		t.Errorf("Expected second GameType Playoffs, got %v", season.GameTypes[1])
	}
}

func TestSeasonGameTypesDisplay(t *testing.T) {
	tests := []struct {
		name     string
		season   SeasonGameTypes
		expected string
	}{
		{
			name: "Regular season and playoffs",
			season: SeasonGameTypes{
				Season:    NewSeason(2024),
				GameTypes: []GameType{GameTypeRegularSeason, GameTypePlayoffs},
			},
			expected: "2024-2025: Regular Season, Playoffs",
		},
		{
			name: "Regular season only",
			season: SeasonGameTypes{
				Season:    NewSeason(2023),
				GameTypes: []GameType{GameTypeRegularSeason},
			},
			expected: "2023-2024: Regular Season",
		},
		{
			name: "All-star only",
			season: SeasonGameTypes{
				Season:    NewSeason(2023),
				GameTypes: []GameType{GameTypeAllStar},
			},
			expected: "2023-2024: All-Star",
		},
		{
			name: "Preseason only",
			season: SeasonGameTypes{
				Season:    NewSeason(2024),
				GameTypes: []GameType{GameTypePreseason},
			},
			expected: "2024-2025: Preseason",
		},
		{
			name: "All types mixed order",
			season: SeasonGameTypes{
				Season:    NewSeason(2024),
				GameTypes: []GameType{GameTypePreseason, GameTypeRegularSeason, GameTypeAllStar, GameTypePlayoffs},
			},
			expected: "2024-2025: Preseason, Regular Season, All-Star, Playoffs",
		},
		{
			name: "Empty game types",
			season: SeasonGameTypes{
				Season:    NewSeason(2024),
				GameTypes: []GameType{},
			},
			expected: "2024-2025: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.season.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSkaterStatsDisplay(t *testing.T) {
	stats := ClubSkaterStats{
		PlayerID:            PlayerID(8475233),
		Headshot:            "test.png",
		FirstName:           LocalizedString{Default: "David"},
		LastName:            LocalizedString{Default: "Savard"},
		Position:            PositionDefense,
		GamesPlayed:         75,
		Goals:               1,
		Assists:             14,
		Points:              15,
		PlusMinus:           -8,
		PenaltyMinutes:      36,
		PowerPlayGoals:      0,
		ShorthandedGoals:    0,
		GameWinningGoals:    0,
		OvertimeGoals:       0,
		Shots:               48,
		ShootingPctg:        0.020833,
		AvgTimeOnIcePerGame: 995.36,
		AvgShiftsPerGame:    19.84,
		FaceoffWinPctg:      0.0,
	}

	expected := "David Savard - 75 GP, 1 G, 14 A, 15 PTS"
	result := stats.String()
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestGoalieStatsDisplay(t *testing.T) {
	stats := ClubGoalieStats{
		PlayerID:            PlayerID(8478470),
		Headshot:            "test.png",
		FirstName:           LocalizedString{Default: "Sam"},
		LastName:            LocalizedString{Default: "Montembeault"},
		GamesPlayed:         62,
		GamesStarted:        60,
		Wins:                31,
		Losses:              24,
		OvertimeLosses:      7,
		GoalsAgainstAverage: 2.818349,
		SavePercentage:      0.901669,
		ShotsAgainst:        1678,
		Saves:               1513,
		GoalsAgainst:        166,
		Shutouts:            4,
		Goals:               0,
		Assists:             1,
		Points:              1,
		PenaltyMinutes:      0,
		TimeOnIce:           212039,
	}

	expected := "Sam Montembeault - 62 GP, 31-24-7, 2.818 GAA, 0.902 SV%"
	result := stats.String()
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestClubStatsEmptyLists(t *testing.T) {
	jsonData := `{
		"season": "20242025",
		"gameType": 2,
		"skaters": [],
		"goalies": []
	}`

	var stats ClubStats
	if err := json.Unmarshal([]byte(jsonData), &stats); err != nil {
		t.Fatalf("Failed to unmarshal ClubStats: %v", err)
	}

	if len(stats.Skaters) != 0 {
		t.Errorf("Expected 0 skaters, got %d", len(stats.Skaters))
	}
	if len(stats.Goalies) != 0 {
		t.Errorf("Expected 0 goalies, got %d", len(stats.Goalies))
	}
}

func TestClubStatsWithAllGameTypes(t *testing.T) {
	tests := []struct {
		name     string
		gameType int
		expected GameType
	}{
		{"Preseason", 1, GameTypePreseason},
		{"Regular Season", 2, GameTypeRegularSeason},
		{"Playoffs", 3, GameTypePlayoffs},
		{"All-Star", 4, GameTypeAllStar},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData := fmt.Sprintf(`{
				"season": "20242025",
				"gameType": %d,
				"skaters": [],
				"goalies": []
			}`, tt.gameType)

			var stats ClubStats
			if err := json.Unmarshal([]byte(jsonData), &stats); err != nil {
				t.Fatalf("Failed to unmarshal ClubStats: %v", err)
			}

			if stats.GameType != tt.expected {
				t.Errorf("Expected GameType %v, got %v", tt.expected, stats.GameType)
			}
		})
	}
}

func TestSeasonGameTypesUnknownGameType(t *testing.T) {
	jsonData := `{
		"season": 20242025,
		"gameTypes": [2, 99]
	}`

	var season SeasonGameTypes
	err := json.Unmarshal([]byte(jsonData), &season)
	if err == nil {
		t.Error("Expected error for unknown game type, got nil")
	}
}

func TestClubSkaterStatsSerializationRoundTrip(t *testing.T) {
	original := ClubSkaterStats{
		PlayerID:            PlayerID(8475233),
		Headshot:            "test.png",
		FirstName:           LocalizedString{Default: "David"},
		LastName:            LocalizedString{Default: "Savard"},
		Position:            PositionDefense,
		GamesPlayed:         75,
		Goals:               1,
		Assists:             14,
		Points:              15,
		PlusMinus:           -8,
		PenaltyMinutes:      36,
		PowerPlayGoals:      0,
		ShorthandedGoals:    0,
		GameWinningGoals:    0,
		OvertimeGoals:       0,
		Shots:               48,
		ShootingPctg:        0.020833,
		AvgTimeOnIcePerGame: 995.36,
		AvgShiftsPerGame:    19.84,
		FaceoffWinPctg:      0.0,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded ClubSkaterStats
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.PlayerID != original.PlayerID {
		t.Errorf("PlayerID mismatch: expected %d, got %d", original.PlayerID, decoded.PlayerID)
	}
	if decoded.FirstName.Default != original.FirstName.Default {
		t.Errorf("FirstName mismatch: expected %s, got %s", original.FirstName.Default, decoded.FirstName.Default)
	}
	if decoded.Goals != original.Goals {
		t.Errorf("Goals mismatch: expected %d, got %d", original.Goals, decoded.Goals)
	}
}

func TestClubGoalieStatsSerializationRoundTrip(t *testing.T) {
	original := ClubGoalieStats{
		PlayerID:            PlayerID(8478470),
		Headshot:            "test.png",
		FirstName:           LocalizedString{Default: "Sam"},
		LastName:            LocalizedString{Default: "Montembeault"},
		GamesPlayed:         62,
		GamesStarted:        60,
		Wins:                31,
		Losses:              24,
		OvertimeLosses:      7,
		GoalsAgainstAverage: 2.818349,
		SavePercentage:      0.901669,
		ShotsAgainst:        1678,
		Saves:               1513,
		GoalsAgainst:        166,
		Shutouts:            4,
		Goals:               0,
		Assists:             1,
		Points:              1,
		PenaltyMinutes:      0,
		TimeOnIce:           212039,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded ClubGoalieStats
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.PlayerID != original.PlayerID {
		t.Errorf("PlayerID mismatch: expected %d, got %d", original.PlayerID, decoded.PlayerID)
	}
	if decoded.Wins != original.Wins {
		t.Errorf("Wins mismatch: expected %d, got %d", original.Wins, decoded.Wins)
	}
	if decoded.Shutouts != original.Shutouts {
		t.Errorf("Shutouts mismatch: expected %d, got %d", original.Shutouts, decoded.Shutouts)
	}
}

func TestClubStatsSerializationRoundTrip(t *testing.T) {
	original := ClubStats{
		Season:   "20242025",
		GameType: GameTypeRegularSeason,
		Skaters:  []ClubSkaterStats{},
		Goalies:  []ClubGoalieStats{},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded ClubStats
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Season != original.Season {
		t.Errorf("Season mismatch: expected %s, got %s", original.Season, decoded.Season)
	}
	if decoded.GameType != original.GameType {
		t.Errorf("GameType mismatch: expected %v, got %v", original.GameType, decoded.GameType)
	}
}

func TestSeasonGameTypesSerializationRoundTrip(t *testing.T) {
	original := SeasonGameTypes{
		Season:    NewSeason(2024),
		GameTypes: []GameType{GameTypeRegularSeason, GameTypePlayoffs},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded SeasonGameTypes
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Season != original.Season {
		t.Errorf("Season mismatch: expected %s, got %s", original.Season, decoded.Season)
	}
	if len(decoded.GameTypes) != len(original.GameTypes) {
		t.Fatalf("GameTypes length mismatch: expected %d, got %d", len(original.GameTypes), len(decoded.GameTypes))
	}
	for i := range original.GameTypes {
		if decoded.GameTypes[i] != original.GameTypes[i] {
			t.Errorf("GameType[%d] mismatch: expected %v, got %v", i, original.GameTypes[i], decoded.GameTypes[i])
		}
	}
}

func TestSeasonGameTypesMarshalAsIntegers(t *testing.T) {
	season := SeasonGameTypes{
		Season:    NewSeason(2024),
		GameTypes: []GameType{GameTypeRegularSeason, GameTypePlayoffs},
	}

	data, err := json.Marshal(season)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify that gameTypes are serialized as integers
	expected := `{"season":20242025,"gameTypes":[2,3]}`
	if string(data) != expected {
		t.Errorf("Expected JSON '%s', got '%s'", expected, string(data))
	}
}

func TestSeasonGameTypesMarshalInvalidGameType(t *testing.T) {
	season := SeasonGameTypes{
		Season:    NewSeason(2024),
		GameTypes: []GameType{GameType(99)}, // Invalid game type
	}

	_, err := json.Marshal(season)
	if err == nil {
		t.Error("Expected error when marshaling invalid game type, got nil")
	}
}

func TestClubSkaterStatsEquality(t *testing.T) {
	stats1 := ClubSkaterStats{
		PlayerID:    PlayerID(8475233),
		Headshot:    "test.png",
		FirstName:   LocalizedString{Default: "David"},
		LastName:    LocalizedString{Default: "Savard"},
		Position:    PositionDefense,
		GamesPlayed: 75,
		Goals:       1,
		Assists:     14,
		Points:      15,
		PlusMinus:   -8,
	}

	stats2 := stats1
	stats3 := stats1
	stats3.Goals = 10

	if stats1 != stats2 {
		t.Error("stats1 and stats2 should be equal")
	}
	if stats1 == stats3 {
		t.Error("stats1 and stats3 should not be equal")
	}
}

func TestClubGoalieStatsEquality(t *testing.T) {
	stats1 := ClubGoalieStats{
		PlayerID:    PlayerID(8478470),
		Headshot:    "test.png",
		FirstName:   LocalizedString{Default: "Sam"},
		LastName:    LocalizedString{Default: "Montembeault"},
		GamesPlayed: 62,
		Wins:        31,
		Losses:      24,
	}

	stats2 := stats1
	stats3 := stats1
	stats3.Wins = 40

	if stats1 != stats2 {
		t.Error("stats1 and stats2 should be equal")
	}
	if stats1 == stats3 {
		t.Error("stats1 and stats3 should not be equal")
	}
}

func TestClubStatsEquality(t *testing.T) {
	stats1 := ClubStats{
		Season:   "20242025",
		GameType: GameTypeRegularSeason,
		Skaters:  []ClubSkaterStats{},
		Goalies:  []ClubGoalieStats{},
	}

	stats2 := stats1
	stats3 := stats1
	stats3.GameType = GameTypePlayoffs

	if stats1.Season != stats2.Season || stats1.GameType != stats2.GameType {
		t.Error("stats1 and stats2 should be equal")
	}
	if stats1.GameType == stats3.GameType {
		t.Error("stats1 and stats3 should not be equal")
	}
}

func TestSeasonGameTypesEquality(t *testing.T) {
	season1 := SeasonGameTypes{
		Season:    NewSeason(2024),
		GameTypes: []GameType{GameTypeRegularSeason},
	}

	season2 := SeasonGameTypes{
		Season:    NewSeason(2024),
		GameTypes: []GameType{GameTypeRegularSeason},
	}

	season3 := SeasonGameTypes{
		Season:    NewSeason(2023),
		GameTypes: []GameType{GameTypeRegularSeason},
	}

	if season1.Season != season2.Season {
		t.Error("season1 and season2 should be equal")
	}
	if season1.Season == season3.Season {
		t.Error("season1 and season3 should not be equal")
	}
}

func TestClubSkaterStatsWithNegativeStats(t *testing.T) {
	jsonData := `{
		"playerId": 123,
		"headshot": "test.png",
		"firstName": {"default": "Test"},
		"lastName": {"default": "Player"},
		"positionCode": "C",
		"gamesPlayed": 10,
		"goals": 0,
		"assists": 0,
		"points": 0,
		"plusMinus": -15,
		"penaltyMinutes": 0,
		"powerPlayGoals": 0,
		"shorthandedGoals": 0,
		"gameWinningGoals": 0,
		"overtimeGoals": 0,
		"shots": 0,
		"shootingPctg": 0.0,
		"avgTimeOnIcePerGame": 0.0,
		"avgShiftsPerGame": 0.0,
		"faceoffWinPctg": 0.0
	}`

	var stats ClubSkaterStats
	if err := json.Unmarshal([]byte(jsonData), &stats); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if stats.PlusMinus != -15 {
		t.Errorf("Expected PlusMinus -15, got %d", stats.PlusMinus)
	}
}

func TestClubGoalieStatsWithZeroValues(t *testing.T) {
	jsonData := `{
		"playerId": 123,
		"headshot": "test.png",
		"firstName": {"default": "Test"},
		"lastName": {"default": "Goalie"},
		"gamesPlayed": 0,
		"gamesStarted": 0,
		"wins": 0,
		"losses": 0,
		"overtimeLosses": 0,
		"goalsAgainstAverage": 0.0,
		"savePercentage": 0.0,
		"shotsAgainst": 0,
		"saves": 0,
		"goalsAgainst": 0,
		"shutouts": 0,
		"goals": 0,
		"assists": 0,
		"points": 0,
		"penaltyMinutes": 0,
		"timeOnIce": 0
	}`

	var stats ClubGoalieStats
	if err := json.Unmarshal([]byte(jsonData), &stats); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if stats.Wins != 0 {
		t.Errorf("Expected Wins 0, got %d", stats.Wins)
	}
	if stats.Shutouts != 0 {
		t.Errorf("Expected Shutouts 0, got %d", stats.Shutouts)
	}
}

// Additional error path tests for club_stats.go

func TestSeasonGameTypes_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var s SeasonGameTypes
	err := json.Unmarshal([]byte(`{invalid json}`), &s)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestSeasonGameTypes_UnmarshalJSON_InvalidGameType(t *testing.T) {
	var s SeasonGameTypes
	err := json.Unmarshal([]byte(`{"season":20232024,"gameTypes":[999]}`), &s)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid game type")
	}
}
