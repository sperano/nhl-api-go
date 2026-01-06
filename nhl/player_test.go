package nhl

import (
	"encoding/json"
	"testing"
)

func TestPlayerStatsDeserialization(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		validate func(*testing.T, PlayerStats)
	}{
		{
			name: "complete skater stats",
			jsonData: `{
				"gamesPlayed": 82,
				"goals": 41,
				"assists": 52,
				"points": 93,
				"plusMinus": 21,
				"pim": 40,
				"powerPlayGoals": 15,
				"powerPlayPoints": 40,
				"shortHandedGoals": 2,
				"shortHandedPoints": 3,
				"shots": 305,
				"shootingPctg": 0.134,
				"faceoffWinPctg": 0.489,
				"avgToi": "21:30"
			}`,
			validate: func(t *testing.T, stats PlayerStats) {
				if stats.GamesPlayed == nil || *stats.GamesPlayed != 82 {
					t.Errorf("expected GamesPlayed=82, got %v", stats.GamesPlayed)
				}
				if stats.Goals == nil || *stats.Goals != 41 {
					t.Errorf("expected Goals=41, got %v", stats.Goals)
				}
				if stats.Assists == nil || *stats.Assists != 52 {
					t.Errorf("expected Assists=52, got %v", stats.Assists)
				}
				if stats.Points == nil || *stats.Points != 93 {
					t.Errorf("expected Points=93, got %v", stats.Points)
				}
				if stats.PlusMinus == nil || *stats.PlusMinus != 21 {
					t.Errorf("expected PlusMinus=21, got %v", stats.PlusMinus)
				}
				if stats.PIM == nil || *stats.PIM != 40 {
					t.Errorf("expected PIM=40, got %v", stats.PIM)
				}
				if stats.PowerPlayGoals == nil || *stats.PowerPlayGoals != 15 {
					t.Errorf("expected PowerPlayGoals=15, got %v", stats.PowerPlayGoals)
				}
				if stats.PowerPlayPoints == nil || *stats.PowerPlayPoints != 40 {
					t.Errorf("expected PowerPlayPoints=40, got %v", stats.PowerPlayPoints)
				}
				if stats.ShortHandedGoals == nil || *stats.ShortHandedGoals != 2 {
					t.Errorf("expected ShortHandedGoals=2, got %v", stats.ShortHandedGoals)
				}
				if stats.ShortHandedPoints == nil || *stats.ShortHandedPoints != 3 {
					t.Errorf("expected ShortHandedPoints=3, got %v", stats.ShortHandedPoints)
				}
				if stats.Shots == nil || *stats.Shots != 305 {
					t.Errorf("expected Shots=305, got %v", stats.Shots)
				}
				if stats.ShootingPctg == nil || *stats.ShootingPctg != 0.134 {
					t.Errorf("expected ShootingPctg=0.134, got %v", stats.ShootingPctg)
				}
				if stats.FaceoffWinPctg == nil || *stats.FaceoffWinPctg != 0.489 {
					t.Errorf("expected FaceoffWinPctg=0.489, got %v", stats.FaceoffWinPctg)
				}
				if stats.AvgTOI == nil || *stats.AvgTOI != "21:30" {
					t.Errorf("expected AvgTOI=21:30, got %v", stats.AvgTOI)
				}
			},
		},
		{
			name: "complete goalie stats",
			jsonData: `{
				"gamesPlayed": 64,
				"wins": 40,
				"losses": 18,
				"otLosses": 6,
				"shutouts": 5,
				"goalsAgainstAvg": 2.24,
				"savePctg": 0.922
			}`,
			validate: func(t *testing.T, stats PlayerStats) {
				if stats.GamesPlayed == nil || *stats.GamesPlayed != 64 {
					t.Errorf("expected GamesPlayed=64, got %v", stats.GamesPlayed)
				}
				if stats.Wins == nil || *stats.Wins != 40 {
					t.Errorf("expected Wins=40, got %v", stats.Wins)
				}
				if stats.Losses == nil || *stats.Losses != 18 {
					t.Errorf("expected Losses=18, got %v", stats.Losses)
				}
				if stats.OTLosses == nil || *stats.OTLosses != 6 {
					t.Errorf("expected OTLosses=6, got %v", stats.OTLosses)
				}
				if stats.Shutouts == nil || *stats.Shutouts != 5 {
					t.Errorf("expected Shutouts=5, got %v", stats.Shutouts)
				}
				if stats.GoalsAgainstAvg == nil || *stats.GoalsAgainstAvg != 2.24 {
					t.Errorf("expected GoalsAgainstAvg=2.24, got %v", stats.GoalsAgainstAvg)
				}
				if stats.SavePctg == nil || *stats.SavePctg != 0.922 {
					t.Errorf("expected SavePctg=0.922, got %v", stats.SavePctg)
				}
			},
		},
		{
			name: "partial stats",
			jsonData: `{
				"gamesPlayed": 10,
				"goals": 3,
				"assists": 5
			}`,
			validate: func(t *testing.T, stats PlayerStats) {
				if stats.GamesPlayed == nil || *stats.GamesPlayed != 10 {
					t.Errorf("expected GamesPlayed=10, got %v", stats.GamesPlayed)
				}
				if stats.Goals == nil || *stats.Goals != 3 {
					t.Errorf("expected Goals=3, got %v", stats.Goals)
				}
				if stats.Assists == nil || *stats.Assists != 5 {
					t.Errorf("expected Assists=5, got %v", stats.Assists)
				}
				if stats.Points != nil {
					t.Errorf("expected Points=nil, got %v", stats.Points)
				}
				if stats.Wins != nil {
					t.Errorf("expected Wins=nil, got %v", stats.Wins)
				}
			},
		},
		{
			name:     "empty stats",
			jsonData: `{}`,
			validate: func(t *testing.T, stats PlayerStats) {
				if stats.GamesPlayed != nil {
					t.Errorf("expected GamesPlayed=nil, got %v", stats.GamesPlayed)
				}
				if stats.Goals != nil {
					t.Errorf("expected Goals=nil, got %v", stats.Goals)
				}
				if stats.Wins != nil {
					t.Errorf("expected Wins=nil, got %v", stats.Wins)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stats PlayerStats
			if err := json.Unmarshal([]byte(tt.jsonData), &stats); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			tt.validate(t, stats)
		})
	}
}

func TestPlayerStatsSerialization(t *testing.T) {
	gamesPlayed := 82
	goals := 41
	assists := 52
	points := 93
	shootingPctg := 0.134

	stats := PlayerStats{
		GamesPlayed:  &gamesPlayed,
		Goals:        &goals,
		Assists:      &assists,
		Points:       &points,
		ShootingPctg: &shootingPctg,
	}

	data, err := json.Marshal(stats)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled PlayerStats
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.GamesPlayed == nil || *unmarshaled.GamesPlayed != gamesPlayed {
		t.Errorf("expected GamesPlayed=%d, got %v", gamesPlayed, unmarshaled.GamesPlayed)
	}
	if unmarshaled.Goals == nil || *unmarshaled.Goals != goals {
		t.Errorf("expected Goals=%d, got %v", goals, unmarshaled.Goals)
	}
	if unmarshaled.Assists == nil || *unmarshaled.Assists != assists {
		t.Errorf("expected Assists=%d, got %v", assists, unmarshaled.Assists)
	}
	if unmarshaled.Points == nil || *unmarshaled.Points != points {
		t.Errorf("expected Points=%d, got %v", points, unmarshaled.Points)
	}
	if unmarshaled.ShootingPctg == nil || *unmarshaled.ShootingPctg != shootingPctg {
		t.Errorf("expected ShootingPctg=%f, got %v", shootingPctg, unmarshaled.ShootingPctg)
	}
}

func TestDraftDetailsDeserialization(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		validate func(*testing.T, DraftDetails)
	}{
		{
			name: "first overall pick",
			jsonData: `{
				"year": 2015,
				"teamAbbrev": "EDM",
				"round": 1,
				"pickInRound": 1,
				"overallPick": 1
			}`,
			validate: func(t *testing.T, draft DraftDetails) {
				if draft.Year != 2015 {
					t.Errorf("expected Year=2015, got %d", draft.Year)
				}
				if draft.TeamAbbrev != "EDM" {
					t.Errorf("expected TeamAbbrev=EDM, got %s", draft.TeamAbbrev)
				}
				if draft.Round != 1 {
					t.Errorf("expected Round=1, got %d", draft.Round)
				}
				if draft.PickInRound != 1 {
					t.Errorf("expected PickInRound=1, got %d", draft.PickInRound)
				}
				if draft.OverallPick != 1 {
					t.Errorf("expected OverallPick=1, got %d", draft.OverallPick)
				}
			},
		},
		{
			name: "later round pick",
			jsonData: `{
				"year": 2010,
				"teamAbbrev": "BOS",
				"round": 5,
				"pickInRound": 15,
				"overallPick": 145
			}`,
			validate: func(t *testing.T, draft DraftDetails) {
				if draft.Year != 2010 {
					t.Errorf("expected Year=2010, got %d", draft.Year)
				}
				if draft.TeamAbbrev != "BOS" {
					t.Errorf("expected TeamAbbrev=BOS, got %s", draft.TeamAbbrev)
				}
				if draft.Round != 5 {
					t.Errorf("expected Round=5, got %d", draft.Round)
				}
				if draft.PickInRound != 15 {
					t.Errorf("expected PickInRound=15, got %d", draft.PickInRound)
				}
				if draft.OverallPick != 145 {
					t.Errorf("expected OverallPick=145, got %d", draft.OverallPick)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var draft DraftDetails
			if err := json.Unmarshal([]byte(tt.jsonData), &draft); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			tt.validate(t, draft)
		})
	}
}

func TestDraftDetailsSerialization(t *testing.T) {
	draft := DraftDetails{
		Year:        2015,
		TeamAbbrev:  "EDM",
		Round:       1,
		PickInRound: 1,
		OverallPick: 1,
	}

	data, err := json.Marshal(draft)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled DraftDetails
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled != draft {
		t.Errorf("expected %+v, got %+v", draft, unmarshaled)
	}
}

func TestFeaturedStatsDeserialization(t *testing.T) {
	jsonData := `{
		"season": 20232024,
		"regularSeason": {
			"gamesPlayed": 76,
			"goals": 32,
			"assists": 68,
			"points": 100
		},
		"playoffs": {
			"gamesPlayed": 25,
			"goals": 8,
			"assists": 34,
			"points": 42
		}
	}`

	var stats FeaturedStats
	if err := json.Unmarshal([]byte(jsonData), &stats); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if stats.Season != 20232024 {
		t.Errorf("expected Season=20232024, got %d", stats.Season)
	}
	if stats.RegularSeason.GamesPlayed == nil || *stats.RegularSeason.GamesPlayed != 76 {
		t.Errorf("expected RegularSeason.GamesPlayed=76, got %v", stats.RegularSeason.GamesPlayed)
	}
	if stats.RegularSeason.Points == nil || *stats.RegularSeason.Points != 100 {
		t.Errorf("expected RegularSeason.Points=100, got %v", stats.RegularSeason.Points)
	}
	if stats.Playoffs == nil {
		t.Fatal("expected Playoffs to be non-nil")
	}
	if stats.Playoffs.GamesPlayed == nil || *stats.Playoffs.GamesPlayed != 25 {
		t.Errorf("expected Playoffs.GamesPlayed=25, got %v", stats.Playoffs.GamesPlayed)
	}
}

func TestCareerTotalsDeserialization(t *testing.T) {
	jsonData := `{
		"regularSeason": {
			"gamesPlayed": 900,
			"goals": 400,
			"assists": 500,
			"points": 900
		},
		"playoffs": {
			"gamesPlayed": 100,
			"goals": 50,
			"assists": 60,
			"points": 110
		}
	}`

	var totals CareerTotals
	if err := json.Unmarshal([]byte(jsonData), &totals); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if totals.RegularSeason.GamesPlayed == nil || *totals.RegularSeason.GamesPlayed != 900 {
		t.Errorf("expected RegularSeason.GamesPlayed=900, got %v", totals.RegularSeason.GamesPlayed)
	}
	if totals.Playoffs == nil {
		t.Fatal("expected Playoffs to be non-nil")
	}
	if totals.Playoffs.GamesPlayed == nil || *totals.Playoffs.GamesPlayed != 100 {
		t.Errorf("expected Playoffs.GamesPlayed=100, got %v", totals.Playoffs.GamesPlayed)
	}
}

func TestSeasonTotalDeserialization(t *testing.T) {
	jsonData := `{
		"season": 20232024,
		"gameTypeId": 2,
		"leagueAbbrev": "NHL",
		"teamName": {"default": "Edmonton Oilers"},
		"teamCommonName": {"default": "Oilers"},
		"sequence": 1,
		"gamesPlayed": 82,
		"goals": 41,
		"assists": 52,
		"points": 93,
		"plusMinus": 21,
		"pim": 40
	}`

	var seasonTotal SeasonTotal
	if err := json.Unmarshal([]byte(jsonData), &seasonTotal); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if seasonTotal.Season != 20232024 {
		t.Errorf("expected Season=20232024, got %d", seasonTotal.Season)
	}
	if seasonTotal.GameType != GameTypeRegularSeason {
		t.Errorf("expected GameType=RegularSeason, got %v", seasonTotal.GameType)
	}
	if seasonTotal.LeagueAbbrev != "NHL" {
		t.Errorf("expected LeagueAbbrev=NHL, got %s", seasonTotal.LeagueAbbrev)
	}
	if seasonTotal.TeamName.Default != "Edmonton Oilers" {
		t.Errorf("expected TeamName=Edmonton Oilers, got %s", seasonTotal.TeamName.Default)
	}
	if seasonTotal.TeamCommonName == nil || seasonTotal.TeamCommonName.Default != "Oilers" {
		t.Errorf("expected TeamCommonName=Oilers, got %v", seasonTotal.TeamCommonName)
	}
	if seasonTotal.Sequence == nil || *seasonTotal.Sequence != 1 {
		t.Errorf("expected Sequence=1, got %v", seasonTotal.Sequence)
	}
	if seasonTotal.GamesPlayed != 82 {
		t.Errorf("expected GamesPlayed=82, got %d", seasonTotal.GamesPlayed)
	}
	if seasonTotal.Goals == nil || *seasonTotal.Goals != 41 {
		t.Errorf("expected Goals=41, got %v", seasonTotal.Goals)
	}
}

func TestAwardDeserialization(t *testing.T) {
	jsonData := `{
		"trophy": {"default": "Hart Memorial Trophy"},
		"seasons": [
			{"seasonId": 20142015},
			{"seasonId": 20162017}
		]
	}`

	var award Award
	if err := json.Unmarshal([]byte(jsonData), &award); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if award.Trophy.Default != "Hart Memorial Trophy" {
		t.Errorf("expected Trophy=Hart Memorial Trophy, got %s", award.Trophy.Default)
	}
	if len(award.Seasons) != 2 {
		t.Fatalf("expected 2 seasons, got %d", len(award.Seasons))
	}
	if award.Seasons[0].SeasonID != 20142015 {
		t.Errorf("expected Seasons[0].SeasonID=20142015, got %d", award.Seasons[0].SeasonID)
	}
	if award.Seasons[1].SeasonID != 20162017 {
		t.Errorf("expected Seasons[1].SeasonID=20162017, got %d", award.Seasons[1].SeasonID)
	}
}

func TestGameLogDeserialization(t *testing.T) {
	jsonData := `{
		"gameId": 2023020001,
		"gameDate": "2023-10-10",
		"teamAbbrev": "EDM",
		"homeRoadFlag": "H",
		"opponentAbbrev": "VAN",
		"goals": 2,
		"assists": 3,
		"points": 5,
		"plusMinus": 2,
		"powerPlayGoals": 1,
		"powerPlayPoints": 2,
		"shots": 6,
		"shifts": 25,
		"toi": "21:30",
		"gameWinningGoals": 1,
		"otGoals": 0,
		"pim": 2
	}`

	var gameLog GameLog
	if err := json.Unmarshal([]byte(jsonData), &gameLog); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if gameLog.GameID != 2023020001 {
		t.Errorf("expected GameID=2023020001, got %d", gameLog.GameID)
	}
	if gameLog.GameDate != "2023-10-10" {
		t.Errorf("expected GameDate=2023-10-10, got %s", gameLog.GameDate)
	}
	if gameLog.TeamAbbrev != "EDM" {
		t.Errorf("expected TeamAbbrev=EDM, got %s", gameLog.TeamAbbrev)
	}
	if gameLog.HomeRoadFlag != HomeRoadHome {
		t.Errorf("expected HomeRoadFlag=H, got %v", gameLog.HomeRoadFlag)
	}
	if gameLog.OpponentAbbrev != "VAN" {
		t.Errorf("expected OpponentAbbrev=VAN, got %s", gameLog.OpponentAbbrev)
	}
	if gameLog.Goals != 2 {
		t.Errorf("expected Goals=2, got %d", gameLog.Goals)
	}
	if gameLog.Assists != 3 {
		t.Errorf("expected Assists=3, got %d", gameLog.Assists)
	}
	if gameLog.Points != 5 {
		t.Errorf("expected Points=5, got %d", gameLog.Points)
	}
	if gameLog.GameWinningGoals == nil || *gameLog.GameWinningGoals != 1 {
		t.Errorf("expected GameWinningGoals=1, got %v", gameLog.GameWinningGoals)
	}
	if gameLog.OTGoals == nil || *gameLog.OTGoals != 0 {
		t.Errorf("expected OTGoals=0, got %v", gameLog.OTGoals)
	}
	if gameLog.PIM == nil || *gameLog.PIM != 2 {
		t.Errorf("expected PIM=2, got %v", gameLog.PIM)
	}
}

func TestPlayerGameLogDeserialization(t *testing.T) {
	jsonData := `{
		"seasonId": 20232024,
		"gameTypeId": 2,
		"gameLog": [
			{
				"gameId": 2023020001,
				"gameDate": "2023-10-10",
				"teamAbbrev": "EDM",
				"homeRoadFlag": "H",
				"opponentAbbrev": "VAN",
				"goals": 2,
				"assists": 3,
				"points": 5,
				"plusMinus": 2,
				"powerPlayGoals": 1,
				"powerPlayPoints": 2,
				"shots": 6,
				"shifts": 25,
				"toi": "21:30"
			}
		]
	}`

	var playerGameLog PlayerGameLog
	if err := json.Unmarshal([]byte(jsonData), &playerGameLog); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if playerGameLog.Season != 20232024 {
		t.Errorf("expected Season=20232024, got %d", playerGameLog.Season)
	}
	if playerGameLog.GameType != GameTypeRegularSeason {
		t.Errorf("expected GameType=RegularSeason, got %v", playerGameLog.GameType)
	}
	if len(playerGameLog.GameLog) != 1 {
		t.Fatalf("expected 1 game log, got %d", len(playerGameLog.GameLog))
	}
	if playerGameLog.GameLog[0].GameID != 2023020001 {
		t.Errorf("expected GameLog[0].GameID=2023020001, got %d", playerGameLog.GameLog[0].GameID)
	}
}

func TestPlayerGameLogPlayerIDNotSerialized(t *testing.T) {
	playerGameLog := PlayerGameLog{
		PlayerID: 8478402,
		Season:   20232024,
		GameType: GameTypeRegularSeason,
		GameLog:  []GameLog{},
	}

	data, err := json.Marshal(playerGameLog)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if _, ok := result["playerId"]; ok {
		t.Error("playerId should not be serialized")
	}
	if _, ok := result["PlayerID"]; ok {
		t.Error("PlayerID should not be serialized")
	}
}

func TestPlayerSearchResultDeserialization(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		validate func(*testing.T, PlayerSearchResult)
	}{
		{
			name: "complete player search result",
			jsonData: `{
				"playerId": "8478402",
				"name": "Connor McDavid",
				"positionCode": "C",
				"teamId": "22",
				"teamAbbrev": "EDM",
				"sweaterNumber": 97,
				"active": true,
				"height": "6'1\"",
				"birthCity": "Richmond Hill",
				"birthStateProvince": "ON",
				"birthCountry": "CAN"
			}`,
			validate: func(t *testing.T, result PlayerSearchResult) {
				if result.PlayerID != "8478402" {
					t.Errorf("expected PlayerID=8478402, got %s", result.PlayerID)
				}
				if result.Name != "Connor McDavid" {
					t.Errorf("expected Name=Connor McDavid, got %s", result.Name)
				}
				if result.Position != PositionCenter {
					t.Errorf("expected Position=C, got %v", result.Position)
				}
				if result.TeamID == nil || *result.TeamID != "22" {
					t.Errorf("expected TeamID=22, got %v", result.TeamID)
				}
				if result.TeamAbbrev == nil || *result.TeamAbbrev != "EDM" {
					t.Errorf("expected TeamAbbrev=EDM, got %v", result.TeamAbbrev)
				}
				if result.SweaterNumber == nil || *result.SweaterNumber != 97 {
					t.Errorf("expected SweaterNumber=97, got %v", result.SweaterNumber)
				}
				if !result.Active {
					t.Error("expected Active=true")
				}
				if result.Height == nil || *result.Height != "6'1\"" {
					t.Errorf("expected Height=6'1\", got %v", result.Height)
				}
				if result.BirthCity == nil || *result.BirthCity != "Richmond Hill" {
					t.Errorf("expected BirthCity=Richmond Hill, got %v", result.BirthCity)
				}
				if result.BirthStateProvince == nil || *result.BirthStateProvince != "ON" {
					t.Errorf("expected BirthStateProvince=ON, got %v", result.BirthStateProvince)
				}
				if result.BirthCountry == nil || *result.BirthCountry != "CAN" {
					t.Errorf("expected BirthCountry=CAN, got %v", result.BirthCountry)
				}
			},
		},
		{
			name: "minimal player search result",
			jsonData: `{
				"playerId": "8475790",
				"name": "Sidney Crosby",
				"positionCode": "C",
				"active": true
			}`,
			validate: func(t *testing.T, result PlayerSearchResult) {
				if result.PlayerID != "8475790" {
					t.Errorf("expected PlayerID=8475790, got %s", result.PlayerID)
				}
				if result.Name != "Sidney Crosby" {
					t.Errorf("expected Name=Sidney Crosby, got %s", result.Name)
				}
				if result.Position != PositionCenter {
					t.Errorf("expected Position=C, got %v", result.Position)
				}
				if !result.Active {
					t.Error("expected Active=true")
				}
				if result.TeamID != nil {
					t.Errorf("expected TeamID=nil, got %v", result.TeamID)
				}
				if result.SweaterNumber != nil {
					t.Errorf("expected SweaterNumber=nil, got %v", result.SweaterNumber)
				}
			},
		},
		{
			name: "inactive player",
			jsonData: `{
				"playerId": "8471675",
				"name": "Jaromir Jagr",
				"positionCode": "RW",
				"active": false
			}`,
			validate: func(t *testing.T, result PlayerSearchResult) {
				if result.PlayerID != "8471675" {
					t.Errorf("expected PlayerID=8471675, got %s", result.PlayerID)
				}
				if result.Name != "Jaromir Jagr" {
					t.Errorf("expected Name=Jaromir Jagr, got %s", result.Name)
				}
				if result.Position != PositionRightWing {
					t.Errorf("expected Position=RW, got %v", result.Position)
				}
				if result.Active {
					t.Error("expected Active=false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result PlayerSearchResult
			if err := json.Unmarshal([]byte(tt.jsonData), &result); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			tt.validate(t, result)
		})
	}
}

func TestPlayerLandingDeserialization(t *testing.T) {
	jsonData := `{
		"playerId": 8478402,
		"isActive": true,
		"currentTeamId": 22,
		"currentTeamAbbrev": "EDM",
		"firstName": {"default": "Connor"},
		"lastName": {"default": "McDavid"},
		"sweaterNumber": 97,
		"position": "C",
		"headshot": "https://assets.nhle.com/mugs/nhl/20232024/EDM/8478402.png",
		"heroImage": "https://assets.nhle.com/mugs/actionshots/1296x729/8478402.jpg",
		"heightInInches": 73,
		"weightInPounds": 193,
		"birthDate": "1997-01-13",
		"birthCity": {"default": "Richmond Hill"},
		"birthStateProvince": {"default": "ON"},
		"birthCountry": "CAN",
		"shootsCatches": "L",
		"draftDetails": {
			"year": 2015,
			"teamAbbrev": "EDM",
			"round": 1,
			"pickInRound": 1,
			"overallPick": 1
		},
		"playerSlug": "connor-mcdavid-8478402",
		"featuredStats": {
			"season": 20232024,
			"regularSeason": {
				"gamesPlayed": 76,
				"goals": 32,
				"assists": 68,
				"points": 100
			}
		},
		"careerTotals": {
			"regularSeason": {
				"gamesPlayed": 645,
				"goals": 335,
				"assists": 638,
				"points": 973
			}
		},
		"seasonTotals": [
			{
				"season": 20232024,
				"gameTypeId": 2,
				"leagueAbbrev": "NHL",
				"teamName": {"default": "Edmonton Oilers"},
				"gamesPlayed": 76,
				"goals": 32,
				"assists": 68,
				"points": 100
			}
		],
		"awards": [
			{
				"trophy": {"default": "Hart Memorial Trophy"},
				"seasons": [
					{"seasonId": 20162017}
				]
			}
		],
		"lastFiveGames": [
			{
				"gameId": 2023020001,
				"gameDate": "2023-10-10",
				"teamAbbrev": "EDM",
				"homeRoadFlag": "H",
				"opponentAbbrev": "VAN",
				"goals": 2,
				"assists": 3,
				"points": 5,
				"plusMinus": 2,
				"powerPlayGoals": 1,
				"powerPlayPoints": 2,
				"shots": 6,
				"shifts": 25,
				"toi": "21:30"
			}
		]
	}`

	var player PlayerLanding
	if err := json.Unmarshal([]byte(jsonData), &player); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if player.PlayerID != 8478402 {
		t.Errorf("expected PlayerID=8478402, got %d", player.PlayerID)
	}
	if !player.IsActive {
		t.Error("expected IsActive=true")
	}
	if player.CurrentTeamID == nil || *player.CurrentTeamID != 22 {
		t.Errorf("expected CurrentTeamID=22, got %v", player.CurrentTeamID)
	}
	if player.CurrentTeamAbbrev == nil || *player.CurrentTeamAbbrev != "EDM" {
		t.Errorf("expected CurrentTeamAbbrev=EDM, got %v", player.CurrentTeamAbbrev)
	}
	if player.FirstName.Default != "Connor" {
		t.Errorf("expected FirstName=Connor, got %s", player.FirstName.Default)
	}
	if player.LastName.Default != "McDavid" {
		t.Errorf("expected LastName=McDavid, got %s", player.LastName.Default)
	}
	if player.SweaterNumber == nil || *player.SweaterNumber != 97 {
		t.Errorf("expected SweaterNumber=97, got %v", player.SweaterNumber)
	}
	if player.Position != PositionCenter {
		t.Errorf("expected Position=C, got %v", player.Position)
	}
	if player.HeightInInches != 73 {
		t.Errorf("expected HeightInInches=73, got %d", player.HeightInInches)
	}
	if player.WeightInPounds != 193 {
		t.Errorf("expected WeightInPounds=193, got %d", player.WeightInPounds)
	}
	if player.BirthDate != "1997-01-13" {
		t.Errorf("expected BirthDate=1997-01-13, got %s", player.BirthDate)
	}
	if player.BirthCity == nil || player.BirthCity.Default != "Richmond Hill" {
		t.Errorf("expected BirthCity=Richmond Hill, got %v", player.BirthCity)
	}
	if player.ShootsCatches != HandednessLeft {
		t.Errorf("expected ShootsCatches=L, got %v", player.ShootsCatches)
	}
	if player.DraftDetails == nil {
		t.Fatal("expected DraftDetails to be non-nil")
	}
	if player.DraftDetails.Year != 2015 {
		t.Errorf("expected DraftDetails.Year=2015, got %d", player.DraftDetails.Year)
	}
	if player.FeaturedStats == nil {
		t.Fatal("expected FeaturedStats to be non-nil")
	}
	if player.CareerTotals == nil {
		t.Fatal("expected CareerTotals to be non-nil")
	}
	if len(player.SeasonTotals) != 1 {
		t.Errorf("expected 1 season total, got %d", len(player.SeasonTotals))
	}
	if len(player.Awards) != 1 {
		t.Errorf("expected 1 award, got %d", len(player.Awards))
	}
	if len(player.LastFiveGames) != 1 {
		t.Errorf("expected 1 game in LastFiveGames, got %d", len(player.LastFiveGames))
	}
}

func TestPlayerLandingMinimalDeserialization(t *testing.T) {
	jsonData := `{
		"playerId": 8478402,
		"isActive": false,
		"firstName": {"default": "Connor"},
		"lastName": {"default": "McDavid"},
		"position": "C",
		"headshot": "https://assets.nhle.com/mugs/nhl/20232024/EDM/8478402.png",
		"heightInInches": 73,
		"weightInPounds": 193,
		"birthDate": "1997-01-13",
		"shootsCatches": "L"
	}`

	var player PlayerLanding
	if err := json.Unmarshal([]byte(jsonData), &player); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if player.PlayerID != 8478402 {
		t.Errorf("expected PlayerID=8478402, got %d", player.PlayerID)
	}
	if player.IsActive {
		t.Error("expected IsActive=false")
	}
	if player.CurrentTeamID != nil {
		t.Errorf("expected CurrentTeamID=nil, got %v", player.CurrentTeamID)
	}
	if player.DraftDetails != nil {
		t.Errorf("expected DraftDetails=nil, got %v", player.DraftDetails)
	}
	if player.FeaturedStats != nil {
		t.Errorf("expected FeaturedStats=nil, got %v", player.FeaturedStats)
	}
	if len(player.SeasonTotals) != 0 {
		t.Errorf("expected 0 season totals, got %d", len(player.SeasonTotals))
	}
}
