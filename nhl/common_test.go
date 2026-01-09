package nhl

import (
	"encoding/json"
	"testing"
	"time"
)

// TestLocalizedString_UnmarshalJSON tests unmarshaling LocalizedString from JSON.
func TestLocalizedString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "standard format",
			input:    `{"default": "Boston"}`,
			expected: "Boston",
			wantErr:  false,
		},
		{
			name:     "plain string fallback",
			input:    `"Boston"`,
			expected: "Boston",
			wantErr:  false,
		},
		{
			name:     "empty default",
			input:    `{"default": ""}`,
			expected: "",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    `""`,
			expected: "",
			wantErr:  false,
		},
		{
			name:     "invalid JSON",
			input:    `{invalid}`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "number",
			input:    `123`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ls LocalizedString
			err := json.Unmarshal([]byte(tt.input), &ls)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if ls.Default != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, ls.Default)
			}
		})
	}
}

// TestLocalizedString_MarshalJSON tests marshaling LocalizedString to JSON.
func TestLocalizedString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    LocalizedString
		expected string
	}{
		{
			name:     "standard value",
			input:    LocalizedString{Default: "Boston"},
			expected: `{"default":"Boston"}`,
		},
		{
			name:     "empty value",
			input:    LocalizedString{Default: ""},
			expected: `{"default":""}`,
		},
		{
			name:     "special characters",
			input:    LocalizedString{Default: "Montréal"},
			expected: `{"default":"Montréal"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

// TestLocalizedString_String tests the String method.
func TestLocalizedString_String(t *testing.T) {
	tests := []struct {
		name     string
		input    LocalizedString
		expected string
	}{
		{
			name:     "non-empty",
			input:    LocalizedString{Default: "Bruins"},
			expected: "Bruins",
		},
		{
			name:     "empty",
			input:    LocalizedString{Default: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestLocalizedString_RoundTrip tests marshaling and unmarshaling.
func TestLocalizedString_RoundTrip(t *testing.T) {
	original := LocalizedString{Default: "Toronto"}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded LocalizedString
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Default != original.Default {
		t.Errorf("expected %q, got %q", original.Default, decoded.Default)
	}
}

// TestConference_UnmarshalJSON tests unmarshaling Conference from JSON.
func TestConference_UnmarshalJSON(t *testing.T) {
	input := `{"abbrev": "E", "name": "Eastern"}`
	var conf Conference

	if err := json.Unmarshal([]byte(input), &conf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Abbrev != "E" {
		t.Errorf("expected abbrev %q, got %q", "E", conf.Abbrev)
	}

	if conf.Name != "Eastern" {
		t.Errorf("expected name %q, got %q", "Eastern", conf.Name)
	}
}

// TestDivision_UnmarshalJSON tests unmarshaling Division from JSON.
func TestDivision_UnmarshalJSON(t *testing.T) {
	input := `{"abbrev": "ATL", "name": "Atlantic"}`
	var div Division

	if err := json.Unmarshal([]byte(input), &div); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if div.Abbrev != "ATL" {
		t.Errorf("expected abbrev %q, got %q", "ATL", div.Abbrev)
	}

	if div.Name != "Atlantic" {
		t.Errorf("expected name %q, got %q", "Atlantic", div.Name)
	}
}

// TestFranchise_UnmarshalJSON tests unmarshaling Franchise from JSON.
func TestFranchise_UnmarshalJSON(t *testing.T) {
	input := `{
		"id": 6,
		"fullName": "Boston Bruins",
		"teamCommonName": "Bruins",
		"teamPlaceName": "Boston"
	}`
	var franchise Franchise

	if err := json.Unmarshal([]byte(input), &franchise); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if franchise.ID != 6 {
		t.Errorf("expected id 6, got %d", franchise.ID)
	}

	if franchise.FullName != "Boston Bruins" {
		t.Errorf("expected fullName %q, got %q", "Boston Bruins", franchise.FullName)
	}

	if franchise.TeamCommonName != "Bruins" {
		t.Errorf("expected teamCommonName %q, got %q", "Bruins", franchise.TeamCommonName)
	}

	if franchise.TeamPlaceName != "Boston" {
		t.Errorf("expected teamPlaceName %q, got %q", "Boston", franchise.TeamPlaceName)
	}
}

// TestTeam_UnmarshalJSON tests unmarshaling Team from JSON.
func TestTeam_UnmarshalJSON(t *testing.T) {
	input := `{
		"id": 6,
		"franchiseId": 6,
		"fullName": "Boston Bruins",
		"leagueAbbrev": "NHL",
		"rawTricode": "BOS",
		"tricode": "BOS",
		"teamPlaceName": {"default": "Boston"},
		"teamCommonName": {"default": "Bruins"},
		"teamLogo": "https://assets.nhle.com/logos/nhl/svg/BOS_light.svg",
		"conference": {"abbrev": "E", "name": "Eastern"},
		"division": {"abbrev": "ATL", "name": "Atlantic"}
	}`
	var team Team

	if err := json.Unmarshal([]byte(input), &team); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if team.ID != 6 {
		t.Errorf("expected id 6, got %d", team.ID)
	}

	if team.FranchiseID != 6 {
		t.Errorf("expected franchiseId 6, got %d", team.FranchiseID)
	}

	if team.FullName != "Boston Bruins" {
		t.Errorf("expected fullName %q, got %q", "Boston Bruins", team.FullName)
	}

	if team.RawTricode != "BOS" {
		t.Errorf("expected rawTricode %q, got %q", "BOS", team.RawTricode)
	}

	if team.Tricode != "BOS" {
		t.Errorf("expected tricode %q, got %q", "BOS", team.Tricode)
	}

	if team.TeamPlaceName.Default != "Boston" {
		t.Errorf("expected teamPlaceName %q, got %q", "Boston", team.TeamPlaceName.Default)
	}

	if team.TeamCommonName.Default != "Bruins" {
		t.Errorf("expected teamCommonName %q, got %q", "Bruins", team.TeamCommonName.Default)
	}

	if team.Conference.Abbrev != "E" {
		t.Errorf("expected conference abbrev %q, got %q", "E", team.Conference.Abbrev)
	}

	if team.Division.Name != "Atlantic" {
		t.Errorf("expected division name %q, got %q", "Atlantic", team.Division.Name)
	}
}

// TestRosterPlayer_UnmarshalJSON tests unmarshaling RosterPlayer from JSON.
func TestRosterPlayer_UnmarshalJSON(t *testing.T) {
	input := `{
		"id": 8478402,
		"headshot": "https://assets.nhle.com/mugs/nhl/20232024/BOS/8478402.png",
		"firstName": {"default": "David"},
		"lastName": {"default": "Pastrnak"},
		"sweaterNumber": 88,
		"position": "RW",
		"shootsCatches": "R",
		"heightInInches": 72,
		"weightInPounds": 194,
		"birthDate": "1996-05-25",
		"birthCity": {"default": "Havirov"},
		"birthStateProvince": {"default": "Moravskoslezsky"},
		"birthCountry": "CZE"
	}`
	var player RosterPlayer

	if err := json.Unmarshal([]byte(input), &player); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if player.ID != 8478402 {
		t.Errorf("expected id 8478402, got %d", player.ID)
	}

	if player.FirstName.Default != "David" {
		t.Errorf("expected firstName %q, got %q", "David", player.FirstName.Default)
	}

	if player.LastName.Default != "Pastrnak" {
		t.Errorf("expected lastName %q, got %q", "Pastrnak", player.LastName.Default)
	}

	if player.SweaterNumber != 88 {
		t.Errorf("expected sweaterNumber 88, got %d", player.SweaterNumber)
	}

	if player.Position != PositionRightWing {
		t.Errorf("expected position RW, got %s", player.Position)
	}

	if player.ShootsCatches != HandednessRight {
		t.Errorf("expected shootsCatches R, got %s", player.ShootsCatches)
	}

	if player.HeightInInches != 72 {
		t.Errorf("expected heightInInches 72, got %d", player.HeightInInches)
	}

	if player.WeightInPounds != 194 {
		t.Errorf("expected weightInPounds 194, got %d", player.WeightInPounds)
	}

	if player.BirthDate != "1996-05-25" {
		t.Errorf("expected birthDate %q, got %q", "1996-05-25", player.BirthDate)
	}

	if player.BirthCity == nil || player.BirthCity.Default != "Havirov" {
		t.Errorf("expected birthCity Havirov, got %v", player.BirthCity)
	}

	if player.BirthCountry != "CZE" {
		t.Errorf("expected birthCountry CZE, got %q", player.BirthCountry)
	}
}

// TestRosterPlayer_UnmarshalJSON_OptionalFields tests unmarshaling with missing optional fields.
func TestRosterPlayer_UnmarshalJSON_OptionalFields(t *testing.T) {
	input := `{
		"id": 8478402,
		"headshot": "",
		"firstName": {"default": "David"},
		"lastName": {"default": "Pastrnak"},
		"sweaterNumber": 88,
		"position": "RW",
		"shootsCatches": "R",
		"heightInInches": 72,
		"weightInPounds": 194,
		"birthDate": "1996-05-25",
		"birthCountry": "CZE"
	}`
	var player RosterPlayer

	if err := json.Unmarshal([]byte(input), &player); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if player.BirthCity != nil {
		t.Errorf("expected birthCity to be nil, got %v", player.BirthCity)
	}

	if player.BirthStateProvince != nil {
		t.Errorf("expected birthStateProvince to be nil, got %v", player.BirthStateProvince)
	}
}

// TestRosterPlayer_FullName tests the FullName method.
func TestRosterPlayer_FullName(t *testing.T) {
	tests := []struct {
		name     string
		player   RosterPlayer
		expected string
	}{
		{
			name: "standard name",
			player: RosterPlayer{
				FirstName: LocalizedString{Default: "David"},
				LastName:  LocalizedString{Default: "Pastrnak"},
			},
			expected: "David Pastrnak",
		},
		{
			name: "empty first name",
			player: RosterPlayer{
				FirstName: LocalizedString{Default: ""},
				LastName:  LocalizedString{Default: "Pastrnak"},
			},
			expected: " Pastrnak",
		},
		{
			name: "empty last name",
			player: RosterPlayer{
				FirstName: LocalizedString{Default: "David"},
				LastName:  LocalizedString{Default: ""},
			},
			expected: "David ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.player.FullName()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestRosterPlayer_Age tests the Age method.
func TestRosterPlayer_Age(t *testing.T) {
	tests := []struct {
		name      string
		birthDate string
		wantAge   int
	}{
		{
			name:      "valid birth date",
			birthDate: "1996-05-25",
			wantAge:   calculateExpectedAge("1996-05-25"),
		},
		{
			name:      "recent birth date",
			birthDate: "2000-01-01",
			wantAge:   calculateExpectedAge("2000-01-01"),
		},
		{
			name:      "invalid birth date format",
			birthDate: "invalid",
			wantAge:   -1,
		},
		{
			name:      "empty birth date",
			birthDate: "",
			wantAge:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := RosterPlayer{BirthDate: tt.birthDate}
			age := player.Age()
			if age != tt.wantAge {
				t.Errorf("expected age %d, got %d", tt.wantAge, age)
			}
		})
	}
}

// calculateExpectedAge is a helper function to calculate the expected age.
func calculateExpectedAge(birthDateStr string) int {
	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return -1
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

// TestRosterPlayer_HeightFeetInches tests the HeightFeetInches method.
func TestRosterPlayer_HeightFeetInches(t *testing.T) {
	tests := []struct {
		name           string
		heightInInches int
		expected       string
	}{
		{
			name:           "6 feet tall",
			heightInInches: 72,
			expected:       "6'0\"",
		},
		{
			name:           "6 feet 1 inch",
			heightInInches: 73,
			expected:       "6'1\"",
		},
		{
			name:           "5 feet 9 inches",
			heightInInches: 69,
			expected:       "5'9\"",
		},
		{
			name:           "7 feet",
			heightInInches: 84,
			expected:       "7'0\"",
		},
		{
			name:           "zero height",
			heightInInches: 0,
			expected:       "0'0\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := RosterPlayer{HeightInInches: tt.heightInInches}
			result := player.HeightFeetInches()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestRosterPlayer_BirthPlace tests the BirthPlace method.
func TestRosterPlayer_BirthPlace(t *testing.T) {
	havirov := LocalizedString{Default: "Havirov"}
	moravskoslezsky := LocalizedString{Default: "Moravskoslezsky"}
	boston := LocalizedString{Default: "Boston"}
	massachusetts := LocalizedString{Default: "Massachusetts"}

	tests := []struct {
		name     string
		player   RosterPlayer
		expected string
	}{
		{
			name: "all fields present",
			player: RosterPlayer{
				BirthCity:          &havirov,
				BirthStateProvince: &moravskoslezsky,
				BirthCountry:       "CZE",
			},
			expected: "Havirov, Moravskoslezsky, CZE",
		},
		{
			name: "city and country only",
			player: RosterPlayer{
				BirthCity:    &boston,
				BirthCountry: "USA",
			},
			expected: "Boston, USA",
		},
		{
			name: "country only",
			player: RosterPlayer{
				BirthCountry: "CAN",
			},
			expected: "CAN",
		},
		{
			name: "state and country",
			player: RosterPlayer{
				BirthStateProvince: &massachusetts,
				BirthCountry:       "USA",
			},
			expected: "Massachusetts, USA",
		},
		{
			name:     "no fields",
			player:   RosterPlayer{},
			expected: "",
		},
		{
			name: "empty string fields",
			player: RosterPlayer{
				BirthCity:    &LocalizedString{Default: ""},
				BirthCountry: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.player.BirthPlace()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestRoster_UnmarshalJSON tests unmarshaling Roster from JSON.
func TestRoster_UnmarshalJSON(t *testing.T) {
	input := `{
		"forwards": [
			{
				"id": 1,
				"headshot": "",
				"firstName": {"default": "Center"},
				"lastName": {"default": "One"},
				"sweaterNumber": 11,
				"position": "C",
				"shootsCatches": "L",
				"heightInInches": 72,
				"weightInPounds": 200,
				"birthDate": "1990-01-01",
				"birthCountry": "CAN"
			}
		],
		"defensemen": [
			{
				"id": 2,
				"headshot": "",
				"firstName": {"default": "Defense"},
				"lastName": {"default": "One"},
				"sweaterNumber": 2,
				"position": "D",
				"shootsCatches": "R",
				"heightInInches": 74,
				"weightInPounds": 210,
				"birthDate": "1992-01-01",
				"birthCountry": "USA"
			}
		],
		"goalies": [
			{
				"id": 3,
				"headshot": "",
				"firstName": {"default": "Goalie"},
				"lastName": {"default": "One"},
				"sweaterNumber": 30,
				"position": "G",
				"shootsCatches": "L",
				"heightInInches": 76,
				"weightInPounds": 220,
				"birthDate": "1988-01-01",
				"birthCountry": "FIN"
			}
		]
	}`

	var roster Roster
	if err := json.Unmarshal([]byte(input), &roster); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(roster.Forwards) != 1 {
		t.Errorf("expected 1 forward, got %d", len(roster.Forwards))
	}

	if len(roster.Defensemen) != 1 {
		t.Errorf("expected 1 defenseman, got %d", len(roster.Defensemen))
	}

	if len(roster.Goalies) != 1 {
		t.Errorf("expected 1 goalie, got %d", len(roster.Goalies))
	}

	if roster.Forwards[0].Position != PositionCenter {
		t.Errorf("expected forward position C, got %s", roster.Forwards[0].Position)
	}

	if roster.Defensemen[0].Position != PositionDefense {
		t.Errorf("expected defenseman position D, got %s", roster.Defensemen[0].Position)
	}

	if roster.Goalies[0].Position != PositionGoalie {
		t.Errorf("expected goalie position G, got %s", roster.Goalies[0].Position)
	}
}

// TestRoster_AllPlayers tests the AllPlayers method.
func TestRoster_AllPlayers(t *testing.T) {
	roster := Roster{
		Forwards: []RosterPlayer{
			{ID: 1, Position: PositionCenter},
			{ID: 2, Position: PositionLeftWing},
		},
		Defensemen: []RosterPlayer{
			{ID: 3, Position: PositionDefense},
		},
		Goalies: []RosterPlayer{
			{ID: 4, Position: PositionGoalie},
		},
	}

	all := roster.AllPlayers()

	expectedCount := 4
	if len(all) != expectedCount {
		t.Errorf("expected %d players, got %d", expectedCount, len(all))
	}

	// Verify order: forwards, defensemen, goalies
	if all[0].ID != 1 || all[1].ID != 2 {
		t.Error("forwards should come first")
	}

	if all[2].ID != 3 {
		t.Error("defensemen should come second")
	}

	if all[3].ID != 4 {
		t.Error("goalies should come last")
	}
}

// TestRoster_PlayerCount tests the PlayerCount method.
func TestRoster_PlayerCount(t *testing.T) {
	tests := []struct {
		name     string
		roster   Roster
		expected int
	}{
		{
			name: "full roster",
			roster: Roster{
				Forwards:   make([]RosterPlayer, 12),
				Defensemen: make([]RosterPlayer, 6),
				Goalies:    make([]RosterPlayer, 2),
			},
			expected: 20,
		},
		{
			name: "empty roster",
			roster: Roster{
				Forwards:   []RosterPlayer{},
				Defensemen: []RosterPlayer{},
				Goalies:    []RosterPlayer{},
			},
			expected: 0,
		},
		{
			name: "partial roster",
			roster: Roster{
				Forwards:   make([]RosterPlayer, 5),
				Defensemen: []RosterPlayer{},
				Goalies:    make([]RosterPlayer, 1),
			},
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := tt.roster.PlayerCount()
			if count != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, count)
			}
		})
	}
}

// TestRoster_MarshalUnmarshalRoundTrip tests marshaling and unmarshaling a complete roster.
func TestRoster_MarshalUnmarshalRoundTrip(t *testing.T) {
	original := Roster{
		Forwards: []RosterPlayer{
			{
				ID:             8478402,
				FirstName:      LocalizedString{Default: "David"},
				LastName:       LocalizedString{Default: "Pastrnak"},
				SweaterNumber:  88,
				Position:       PositionRightWing,
				ShootsCatches:  HandednessRight,
				HeightInInches: 72,
				WeightInPounds: 194,
				BirthDate:      "1996-05-25",
				BirthCountry:   "CZE",
			},
		},
		Defensemen: []RosterPlayer{
			{
				ID:             8470794,
				FirstName:      LocalizedString{Default: "Zdeno"},
				LastName:       LocalizedString{Default: "Chara"},
				SweaterNumber:  33,
				Position:       PositionDefense,
				ShootsCatches:  HandednessLeft,
				HeightInInches: 81,
				WeightInPounds: 250,
				BirthDate:      "1977-03-18",
				BirthCountry:   "SVK",
			},
		},
		Goalies: []RosterPlayer{
			{
				ID:             8471695,
				FirstName:      LocalizedString{Default: "Tuukka"},
				LastName:       LocalizedString{Default: "Rask"},
				SweaterNumber:  40,
				Position:       PositionGoalie,
				ShootsCatches:  HandednessLeft,
				HeightInInches: 75,
				WeightInPounds: 176,
				BirthDate:      "1987-03-10",
				BirthCountry:   "FIN",
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded Roster
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if len(decoded.Forwards) != len(original.Forwards) {
		t.Errorf("expected %d forwards, got %d", len(original.Forwards), len(decoded.Forwards))
	}

	if len(decoded.Defensemen) != len(original.Defensemen) {
		t.Errorf("expected %d defensemen, got %d", len(original.Defensemen), len(decoded.Defensemen))
	}

	if len(decoded.Goalies) != len(original.Goalies) {
		t.Errorf("expected %d goalies, got %d", len(original.Goalies), len(decoded.Goalies))
	}

	// Verify a sample player
	if decoded.Forwards[0].FullName() != "David Pastrnak" {
		t.Errorf("expected forward name 'David Pastrnak', got %q", decoded.Forwards[0].FullName())
	}
}
