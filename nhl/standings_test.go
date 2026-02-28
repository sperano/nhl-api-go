package nhl

import (
	"encoding/json"
	"testing"
)

func TestStandingsResponseDeserialization(t *testing.T) {
	jsonData := `{
		"standings": [
			{
				"conferenceAbbrev": "E",
				"conferenceName": "Eastern",
				"divisionAbbrev": "ATL",
				"divisionName": "Atlantic",
				"teamName": {"default": "Buffalo Sabres"},
				"teamCommonName": {"default": "Sabres"},
				"teamAbbrev": {"default": "BUF"},
				"teamLogo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg",
				"wins": 10,
				"losses": 5,
				"otLosses": 2,
				"points": 22
			}
		]
	}`

	var response StandingsResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal StandingsResponse: %v", err)
	}

	if len(response.Standings) != 1 {
		t.Fatalf("expected 1 standing, got %d", len(response.Standings))
	}

	standing := response.Standings[0]
	if standing.TeamAbbrev.Default != "BUF" {
		t.Errorf("expected TeamAbbrev.Default = BUF, got %s", standing.TeamAbbrev.Default)
	}
	if standing.Wins != 10 {
		t.Errorf("expected Wins = 10, got %d", standing.Wins)
	}
	if standing.Points != 22 {
		t.Errorf("expected Points = 22, got %d", standing.Points)
	}
}

func TestStandingToTeamConversion(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("W"),
		ConferenceName:   stringPtr("Western"),
		DivisionAbbrev:   "PAC",
		DivisionName:     "Pacific",
		TeamName:         LocalizedString{Default: "Vegas Golden Knights"},
		TeamCommonName:   LocalizedString{Default: "Golden Knights"},
		TeamAbbrev:       LocalizedString{Default: "VGK"},
		TeamLogo:         "https://assets.nhle.com/logos/nhl/svg/VGK_light.svg",
		Wins:             12,
		Losses:           3,
		OTLosses:         1,
		Points:           25,
	}

	team := standing.ToTeam()

	if team.FullName != "Vegas Golden Knights" {
		t.Errorf("expected FullName = Vegas Golden Knights, got %s", team.FullName)
	}
	if team.TeamCommonName.Default != "Golden Knights" {
		t.Errorf("expected TeamCommonName.Default = Golden Knights, got %s", team.TeamCommonName.Default)
	}
	if team.Tricode != "VGK" {
		t.Errorf("expected Tricode = VGK, got %s", team.Tricode)
	}
	if team.TeamLogo != "https://assets.nhle.com/logos/nhl/svg/VGK_light.svg" {
		t.Errorf("expected TeamLogo = https://assets.nhle.com/logos/nhl/svg/VGK_light.svg, got %s", team.TeamLogo)
	}
	if team.Conference.Abbrev != "W" {
		t.Errorf("expected Conference.Abbrev = W, got %s", team.Conference.Abbrev)
	}
	if team.Conference.Name != "Western" {
		t.Errorf("expected Conference.Name = Western, got %s", team.Conference.Name)
	}
	if team.Division.Abbrev != "PAC" {
		t.Errorf("expected Division.Abbrev = PAC, got %s", team.Division.Abbrev)
	}
	if team.Division.Name != "Pacific" {
		t.Errorf("expected Division.Name = Pacific, got %s", team.Division.Name)
	}
}

func TestStandingDisplay(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("E"),
		ConferenceName:   stringPtr("Eastern"),
		DivisionAbbrev:   "ATL",
		DivisionName:     "Atlantic",
		TeamName:         LocalizedString{Default: "Boston Bruins"},
		TeamCommonName:   LocalizedString{Default: "Bruins"},
		TeamAbbrev:       LocalizedString{Default: "BOS"},
		TeamLogo:         "https://assets.nhle.com/logos/nhl/svg/BOS_light.svg",
		Wins:             15,
		Losses:           2,
		OTLosses:         1,
		Points:           31,
	}

	expected := "BOS: 31 pts (15-2-1)"
	if standing.String() != expected {
		t.Errorf("expected %q, got %q", expected, standing.String())
	}
}

func TestStandingsResponseWithExtraFields(t *testing.T) {
	jsonData := `{
		"wildCardIndicator": true,
		"standingsDateTimeUtc": "2024-01-15T12:00:00Z",
		"standings": [
			{
				"conferenceAbbrev": "E",
				"conferenceName": "Eastern",
				"divisionAbbrev": "ATL",
				"divisionName": "Atlantic",
				"teamName": {"default": "Buffalo Sabres"},
				"teamCommonName": {"default": "Sabres"},
				"teamAbbrev": {"default": "BUF"},
				"teamLogo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg",
				"wins": 10,
				"losses": 5,
				"otLosses": 2,
				"points": 22
			}
		]
	}`

	var response StandingsResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal StandingsResponse: %v", err)
	}

	if len(response.Standings) != 1 {
		t.Fatalf("expected 1 standing, got %d", len(response.Standings))
	}
	if response.Standings[0].TeamAbbrev.Default != "BUF" {
		t.Errorf("expected TeamAbbrev.Default = BUF, got %s", response.Standings[0].TeamAbbrev.Default)
	}
}

func TestStandingsResponseEmpty(t *testing.T) {
	jsonData := `{
		"wildCardIndicator": true,
		"standings": []
	}`

	var response StandingsResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal StandingsResponse: %v", err)
	}

	if len(response.Standings) != 0 {
		t.Errorf("expected 0 standings, got %d", len(response.Standings))
	}
}

func TestStandingsWithoutConferenceFields(t *testing.T) {
	jsonData := `{
		"standings": [
			{
				"divisionAbbrev": "EAST",
				"divisionName": "East",
				"teamName": {"default": "Boston Bruins"},
				"teamCommonName": {"default": "Bruins"},
				"teamAbbrev": {"default": "BOS"},
				"teamLogo": "https://assets.nhle.com/logos/nhl/svg/BOS_light.svg",
				"wins": 20,
				"losses": 10,
				"otLosses": 5,
				"points": 45
			}
		]
	}`

	var response StandingsResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal StandingsResponse: %v", err)
	}

	if len(response.Standings) != 1 {
		t.Fatalf("expected 1 standing, got %d", len(response.Standings))
	}

	standing := response.Standings[0]
	if standing.ConferenceAbbrev != nil {
		t.Errorf("expected ConferenceAbbrev = nil, got %v", standing.ConferenceAbbrev)
	}
	if standing.ConferenceName != nil {
		t.Errorf("expected ConferenceName = nil, got %v", standing.ConferenceName)
	}
	if standing.DivisionAbbrev != "EAST" {
		t.Errorf("expected DivisionAbbrev = EAST, got %s", standing.DivisionAbbrev)
	}
	if standing.TeamAbbrev.Default != "BOS" {
		t.Errorf("expected TeamAbbrev.Default = BOS, got %s", standing.TeamAbbrev.Default)
	}
	if standing.Wins != 20 {
		t.Errorf("expected Wins = 20, got %d", standing.Wins)
	}
}

func TestStandingToTeamWithoutConference(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: nil,
		ConferenceName:   nil,
		DivisionAbbrev:   "EAST",
		DivisionName:     "East",
		TeamName:         LocalizedString{Default: "Montreal Canadiens"},
		TeamCommonName:   LocalizedString{Default: "Canadiens"},
		TeamAbbrev:       LocalizedString{Default: "MTL"},
		TeamLogo:         "https://assets.nhle.com/logos/nhl/svg/MTL_light.svg",
		Wins:             25,
		Losses:           8,
		OTLosses:         3,
		Points:           53,
	}

	team := standing.ToTeam()

	if team.FullName != "Montreal Canadiens" {
		t.Errorf("expected FullName = Montreal Canadiens, got %s", team.FullName)
	}
	if team.TeamCommonName.Default != "Canadiens" {
		t.Errorf("expected TeamCommonName.Default = Canadiens, got %s", team.TeamCommonName.Default)
	}
	if team.Tricode != "MTL" {
		t.Errorf("expected Tricode = MTL, got %s", team.Tricode)
	}
	if team.Conference.Abbrev != unknownConferenceAbbrev {
		t.Errorf("expected Conference.Abbrev = %s, got %s", unknownConferenceAbbrev, team.Conference.Abbrev)
	}
	if team.Conference.Name != unknownConferenceName {
		t.Errorf("expected Conference.Name = %s, got %s", unknownConferenceName, team.Conference.Name)
	}
	if team.Division.Abbrev != "EAST" {
		t.Errorf("expected Division.Abbrev = EAST, got %s", team.Division.Abbrev)
	}
	if team.Division.Name != "East" {
		t.Errorf("expected Division.Name = East, got %s", team.Division.Name)
	}
}

func TestGamesPlayedTypicalSeason(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("E"),
		ConferenceName:   stringPtr("Eastern"),
		DivisionAbbrev:   "ATL",
		DivisionName:     "Atlantic",
		TeamName:         LocalizedString{Default: "Toronto Maple Leafs"},
		TeamCommonName:   LocalizedString{Default: "Maple Leafs"},
		TeamAbbrev:       LocalizedString{Default: "TOR"},
		TeamLogo:         "https://assets.nhle.com/logos/nhl/svg/TOR_light.svg",
		Wins:             15,
		Losses:           10,
		OTLosses:         2,
		Points:           32,
	}

	expected := 27
	if standing.GamesPlayed() != expected {
		t.Errorf("expected GamesPlayed = %d, got %d", expected, standing.GamesPlayed())
	}
}

func TestGamesPlayedZeroGames(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("W"),
		ConferenceName:   stringPtr("Western"),
		DivisionAbbrev:   "CEN",
		DivisionName:     "Central",
		TeamName:         LocalizedString{Default: "Test Team"},
		TeamCommonName:   LocalizedString{Default: "Test"},
		TeamAbbrev:       LocalizedString{Default: "TST"},
		TeamLogo:         "https://example.com/logo.svg",
		Wins:             0,
		Losses:           0,
		OTLosses:         0,
		Points:           0,
	}

	expected := 0
	if standing.GamesPlayed() != expected {
		t.Errorf("expected GamesPlayed = %d, got %d", expected, standing.GamesPlayed())
	}
}

func TestGamesPlayedOnlyWins(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("E"),
		ConferenceName:   stringPtr("Eastern"),
		DivisionAbbrev:   "ATL",
		DivisionName:     "Atlantic",
		TeamName:         LocalizedString{Default: "Undefeated Team"},
		TeamCommonName:   LocalizedString{Default: "Undefeated"},
		TeamAbbrev:       LocalizedString{Default: "UND"},
		TeamLogo:         "https://example.com/logo.svg",
		Wins:             10,
		Losses:           0,
		OTLosses:         0,
		Points:           20,
	}

	expected := 10
	if standing.GamesPlayed() != expected {
		t.Errorf("expected GamesPlayed = %d, got %d", expected, standing.GamesPlayed())
	}
}

func TestGamesPlayedOnlyLosses(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("W"),
		ConferenceName:   stringPtr("Western"),
		DivisionAbbrev:   "PAC",
		DivisionName:     "Pacific",
		TeamName:         LocalizedString{Default: "Winless Team"},
		TeamCommonName:   LocalizedString{Default: "Winless"},
		TeamAbbrev:       LocalizedString{Default: "WLS"},
		TeamLogo:         "https://example.com/logo.svg",
		Wins:             0,
		Losses:           15,
		OTLosses:         0,
		Points:           0,
	}

	expected := 15
	if standing.GamesPlayed() != expected {
		t.Errorf("expected GamesPlayed = %d, got %d", expected, standing.GamesPlayed())
	}
}

func TestGamesPlayedOnlyOTLosses(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("E"),
		ConferenceName:   stringPtr("Eastern"),
		DivisionAbbrev:   "MET",
		DivisionName:     "Metropolitan",
		TeamName:         LocalizedString{Default: "OT Loss Team"},
		TeamCommonName:   LocalizedString{Default: "OT Loss"},
		TeamAbbrev:       LocalizedString{Default: "OTL"},
		TeamLogo:         "https://example.com/logo.svg",
		Wins:             0,
		Losses:           0,
		OTLosses:         5,
		Points:           5,
	}

	expected := 5
	if standing.GamesPlayed() != expected {
		t.Errorf("expected GamesPlayed = %d, got %d", expected, standing.GamesPlayed())
	}
}

func TestGamesPlayedFullSeason(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("W"),
		ConferenceName:   stringPtr("Western"),
		DivisionAbbrev:   "CEN",
		DivisionName:     "Central",
		TeamName:         LocalizedString{Default: "Colorado Avalanche"},
		TeamCommonName:   LocalizedString{Default: "Avalanche"},
		TeamAbbrev:       LocalizedString{Default: "COL"},
		TeamLogo:         "https://assets.nhle.com/logos/nhl/svg/COL_light.svg",
		Wins:             50,
		Losses:           20,
		OTLosses:         12,
		Points:           112,
	}

	expected := 82
	if standing.GamesPlayed() != expected {
		t.Errorf("expected GamesPlayed = %d (full season), got %d", expected, standing.GamesPlayed())
	}
}

func TestGamesPlayedWithExistingStandings(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("E"),
		ConferenceName:   stringPtr("Eastern"),
		DivisionAbbrev:   "ATL",
		DivisionName:     "Atlantic",
		TeamName:         LocalizedString{Default: "Buffalo Sabres"},
		TeamCommonName:   LocalizedString{Default: "Sabres"},
		TeamAbbrev:       LocalizedString{Default: "BUF"},
		TeamLogo:         "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg",
		Wins:             10,
		Losses:           5,
		OTLosses:         2,
		Points:           22,
	}

	expected := 17
	if standing.GamesPlayed() != expected {
		t.Errorf("expected GamesPlayed = %d, got %d", expected, standing.GamesPlayed())
	}
}

func TestSeasonInfoDeserialization(t *testing.T) {
	jsonData := `{
		"id": 20232024,
		"standingsStart": "2023-10-10",
		"standingsEnd": "2024-04-18"
	}`

	var season SeasonInfo
	if err := json.Unmarshal([]byte(jsonData), &season); err != nil {
		t.Fatalf("failed to unmarshal SeasonInfo: %v", err)
	}

	if season.ID != NewSeason(2023) {
		t.Errorf("expected ID = 20232024, got %s", season.ID)
	}
	if season.StandingsStart != FromYMD(2023, 10, 10) {
		t.Errorf("expected StandingsStart = 2023-10-10, got %s", season.StandingsStart)
	}
	if season.StandingsEnd != FromYMD(2024, 4, 18) {
		t.Errorf("expected StandingsEnd = 2024-04-18, got %s", season.StandingsEnd)
	}
}

func TestSeasonsResponseDeserialization(t *testing.T) {
	jsonData := `{
		"seasons": [
			{
				"id": 20222023,
				"standingsStart": "2022-10-07",
				"standingsEnd": "2023-04-13"
			},
			{
				"id": 20232024,
				"standingsStart": "2023-10-10",
				"standingsEnd": "2024-04-18"
			}
		]
	}`

	var response SeasonsResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal SeasonsResponse: %v", err)
	}

	if len(response.Seasons) != 2 {
		t.Fatalf("expected 2 seasons, got %d", len(response.Seasons))
	}

	if response.Seasons[0].ID != NewSeason(2022) {
		t.Errorf("expected first season ID = 20222023, got %s", response.Seasons[0].ID)
	}
	if response.Seasons[1].ID != NewSeason(2023) {
		t.Errorf("expected second season ID = 20232024, got %s", response.Seasons[1].ID)
	}
}

func TestStandingSerialization(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("E"),
		ConferenceName:   stringPtr("Eastern"),
		DivisionAbbrev:   "ATL",
		DivisionName:     "Atlantic",
		TeamName:         LocalizedString{Default: "Buffalo Sabres"},
		TeamCommonName:   LocalizedString{Default: "Sabres"},
		TeamAbbrev:       LocalizedString{Default: "BUF"},
		TeamLogo:         "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg",
		Wins:             10,
		Losses:           5,
		OTLosses:         2,
		Points:           22,
	}

	data, err := json.Marshal(standing)
	if err != nil {
		t.Fatalf("failed to marshal Standing: %v", err)
	}

	var unmarshaled Standing
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal Standing: %v", err)
	}

	if unmarshaled.TeamAbbrev.Default != standing.TeamAbbrev.Default {
		t.Errorf("expected TeamAbbrev.Default = %s, got %s", standing.TeamAbbrev.Default, unmarshaled.TeamAbbrev.Default)
	}
	if unmarshaled.Wins != standing.Wins {
		t.Errorf("expected Wins = %d, got %d", standing.Wins, unmarshaled.Wins)
	}
}

func TestSeasonInfoSerialization(t *testing.T) {
	season := SeasonInfo{
		ID:             NewSeason(2023),
		StandingsStart: FromYMD(2023, 10, 10),
		StandingsEnd:   FromYMD(2024, 4, 18),
	}

	data, err := json.Marshal(season)
	if err != nil {
		t.Fatalf("failed to marshal SeasonInfo: %v", err)
	}

	var unmarshaled SeasonInfo
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal SeasonInfo: %v", err)
	}

	if unmarshaled.ID != season.ID {
		t.Errorf("expected ID = %s, got %s", season.ID, unmarshaled.ID)
	}
	if unmarshaled.StandingsStart != season.StandingsStart {
		t.Errorf("expected StandingsStart = %s, got %s", season.StandingsStart, unmarshaled.StandingsStart)
	}
}

func TestStandingsResponseSerialization(t *testing.T) {
	response := StandingsResponse{
		Standings: []Standing{
			{
				ConferenceAbbrev: stringPtr("E"),
				ConferenceName:   stringPtr("Eastern"),
				DivisionAbbrev:   "ATL",
				DivisionName:     "Atlantic",
				TeamName:         LocalizedString{Default: "Buffalo Sabres"},
				TeamCommonName:   LocalizedString{Default: "Sabres"},
				TeamAbbrev:       LocalizedString{Default: "BUF"},
				TeamLogo:         "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg",
				Wins:             10,
				Losses:           5,
				OTLosses:         2,
				Points:           22,
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal StandingsResponse: %v", err)
	}

	var unmarshaled StandingsResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal StandingsResponse: %v", err)
	}

	if len(unmarshaled.Standings) != len(response.Standings) {
		t.Errorf("expected %d standings, got %d", len(response.Standings), len(unmarshaled.Standings))
	}
}

func TestSeasonsResponseSerialization(t *testing.T) {
	response := SeasonsResponse{
		Seasons: []SeasonInfo{
			{
				ID:             NewSeason(2023),
				StandingsStart: FromYMD(2023, 10, 10),
				StandingsEnd:   FromYMD(2024, 4, 18),
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal SeasonsResponse: %v", err)
	}

	var unmarshaled SeasonsResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal SeasonsResponse: %v", err)
	}

	if len(unmarshaled.Seasons) != len(response.Seasons) {
		t.Errorf("expected %d seasons, got %d", len(response.Seasons), len(unmarshaled.Seasons))
	}
}

func TestStandingWithAllZeroStats(t *testing.T) {
	standing := Standing{
		ConferenceAbbrev: stringPtr("E"),
		ConferenceName:   stringPtr("Eastern"),
		DivisionAbbrev:   "ATL",
		DivisionName:     "Atlantic",
		TeamName:         LocalizedString{Default: "New Team"},
		TeamCommonName:   LocalizedString{Default: "New"},
		TeamAbbrev:       LocalizedString{Default: "NEW"},
		TeamLogo:         "https://example.com/logo.svg",
		Wins:             0,
		Losses:           0,
		OTLosses:         0,
		Points:           0,
	}

	expected := "NEW: 0 pts (0-0-0)"
	if standing.String() != expected {
		t.Errorf("expected %q, got %q", expected, standing.String())
	}

	if standing.GamesPlayed() != 0 {
		t.Errorf("expected GamesPlayed = 0, got %d", standing.GamesPlayed())
	}
}
