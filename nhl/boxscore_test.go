package nhl

import (
	"encoding/json"
	"testing"
)

func TestBoxscore_Deserialization(t *testing.T) {
	jsonData := `{
		"id": 2024020001,
		"season": 20242025,
		"gameType": 2,
		"limitedScoring": false,
		"gameDate": "2024-10-04",
		"venue": {"default": "Test Arena"},
		"venueLocation": {"default": "Test City"},
		"startTimeUTC": "2024-10-04T19:00:00Z",
		"easternUTCOffset": "-04:00",
		"venueUTCOffset": "-04:00",
		"tvBroadcasts": [],
		"gameState": "LIVE",
		"gameScheduleState": "OK",
		"periodDescriptor": {
			"number": 2,
			"periodType": "REG",
			"maxRegulationPeriods": 3
		},
		"awayTeam": {
			"id": 1,
			"commonName": {"default": "Devils"},
			"abbrev": "NJD",
			"score": 2,
			"sog": 15,
			"logo": "https://assets.nhle.com/logos/nhl/svg/NJD_light.svg",
			"darkLogo": "https://assets.nhle.com/logos/nhl/svg/NJD_dark.svg",
			"placeName": {"default": "New Jersey"},
			"placeNameWithPreposition": {"default": "New Jersey"}
		},
		"homeTeam": {
			"id": 7,
			"commonName": {"default": "Sabres"},
			"abbrev": "BUF",
			"score": 1,
			"sog": 12,
			"logo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg",
			"darkLogo": "https://assets.nhle.com/logos/nhl/svg/BUF_dark.svg",
			"placeName": {"default": "Buffalo"},
			"placeNameWithPreposition": {"default": "Buffalo"}
		},
		"clock": {
			"timeRemaining": "10:15",
			"secondsRemaining": 615,
			"running": true,
			"inIntermission": false
		},
		"playerByGameStats": {
			"awayTeam": {
				"forwards": [],
				"defense": [],
				"goalies": []
			},
			"homeTeam": {
				"forwards": [],
				"defense": [],
				"goalies": []
			}
		}
	}`

	var boxscore Boxscore
	err := json.Unmarshal([]byte(jsonData), &boxscore)
	if err != nil {
		t.Fatalf("Failed to unmarshal boxscore: %v", err)
	}

	if boxscore.ID != GameID(2024020001) {
		t.Errorf("ID = %d, want 2024020001", boxscore.ID)
	}
	if boxscore.Season != NewSeason(2024) {
		t.Errorf("Season = %v, want 2024-2025", boxscore.Season)
	}
	if boxscore.GameType != GameTypeRegularSeason {
		t.Errorf("GameType = %v, want %v", boxscore.GameType, GameTypeRegularSeason)
	}
	if boxscore.GameState != GameStateLive {
		t.Errorf("GameState = %v, want %v", boxscore.GameState, GameStateLive)
	}
	if boxscore.AwayTeam.Abbrev != "NJD" {
		t.Errorf("AwayTeam.Abbrev = %s, want NJD", boxscore.AwayTeam.Abbrev)
	}
	if boxscore.HomeTeam.Abbrev != "BUF" {
		t.Errorf("HomeTeam.Abbrev = %s, want BUF", boxscore.HomeTeam.Abbrev)
	}
	if boxscore.AwayTeam.Score != 2 {
		t.Errorf("AwayTeam.Score = %d, want 2", boxscore.AwayTeam.Score)
	}
	if boxscore.HomeTeam.Score != 1 {
		t.Errorf("HomeTeam.Score = %d, want 1", boxscore.HomeTeam.Score)
	}
	if boxscore.Clock.TimeRemaining != "10:15" {
		t.Errorf("Clock.TimeRemaining = %s, want 10:15", boxscore.Clock.TimeRemaining)
	}
	if boxscore.Clock.SecondsRemaining != 615 {
		t.Errorf("Clock.SecondsRemaining = %d, want 615", boxscore.Clock.SecondsRemaining)
	}
	if !boxscore.Clock.Running {
		t.Errorf("Clock.Running = false, want true")
	}
	if boxscore.PeriodDescriptor.Number != 2 {
		t.Errorf("PeriodDescriptor.Number = %d, want 2", boxscore.PeriodDescriptor.Number)
	}
}

func TestSkaterStats_Deserialization(t *testing.T) {
	jsonData := `{
		"playerId": 8480002,
		"sweaterNumber": 13,
		"name": {"default": "N. Hischier"},
		"position": "C",
		"goals": 1,
		"assists": 2,
		"points": 3,
		"plusMinus": 2,
		"pim": 0,
		"hits": 3,
		"powerPlayGoals": 1,
		"sog": 4,
		"faceoffWinningPctg": 0.55,
		"toi": "18:15",
		"blockedShots": 1,
		"shifts": 27,
		"giveaways": 0,
		"takeaways": 2
	}`

	var stats SkaterStats
	err := json.Unmarshal([]byte(jsonData), &stats)
	if err != nil {
		t.Fatalf("Failed to unmarshal skater stats: %v", err)
	}

	if stats.PlayerID != PlayerID(8480002) {
		t.Errorf("PlayerID = %d, want 8480002", stats.PlayerID)
	}
	if stats.SweaterNumber != 13 {
		t.Errorf("SweaterNumber = %d, want 13", stats.SweaterNumber)
	}
	if stats.Name.Default != "N. Hischier" {
		t.Errorf("Name.Default = %s, want N. Hischier", stats.Name.Default)
	}
	if stats.Position != PositionCenter {
		t.Errorf("Position = %v, want %v", stats.Position, PositionCenter)
	}
	if stats.Goals != 1 {
		t.Errorf("Goals = %d, want 1", stats.Goals)
	}
	if stats.Assists != 2 {
		t.Errorf("Assists = %d, want 2", stats.Assists)
	}
	if stats.Points != 3 {
		t.Errorf("Points = %d, want 3", stats.Points)
	}
	if stats.PlusMinus != 2 {
		t.Errorf("PlusMinus = %d, want 2", stats.PlusMinus)
	}
	if stats.FaceoffWinningPctg != 0.55 {
		t.Errorf("FaceoffWinningPctg = %f, want 0.55", stats.FaceoffWinningPctg)
	}
}

func TestGoalieStats_Deserialization(t *testing.T) {
	jsonData := `{
		"playerId": 8474593,
		"sweaterNumber": 25,
		"name": {"default": "J. Markstrom"},
		"position": "G",
		"evenStrengthShotsAgainst": "25/26",
		"powerPlayShotsAgainst": "5/5",
		"shorthandedShotsAgainst": "0/0",
		"saveShotsAgainst": "30/31",
		"savePctg": 0.967,
		"evenStrengthGoalsAgainst": 1,
		"powerPlayGoalsAgainst": 0,
		"shorthandedGoalsAgainst": 0,
		"pim": 0,
		"goalsAgainst": 1,
		"toi": "59:38",
		"starter": true,
		"decision": "W",
		"shotsAgainst": 31,
		"saves": 30
	}`

	var stats GoalieStats
	err := json.Unmarshal([]byte(jsonData), &stats)
	if err != nil {
		t.Fatalf("Failed to unmarshal goalie stats: %v", err)
	}

	if stats.PlayerID != PlayerID(8474593) {
		t.Errorf("PlayerID = %d, want 8474593", stats.PlayerID)
	}
	if stats.SweaterNumber != 25 {
		t.Errorf("SweaterNumber = %d, want 25", stats.SweaterNumber)
	}
	if stats.Name.Default != "J. Markstrom" {
		t.Errorf("Name.Default = %s, want J. Markstrom", stats.Name.Default)
	}
	if stats.Position != PositionGoalie {
		t.Errorf("Position = %v, want %v", stats.Position, PositionGoalie)
	}
	if stats.SavePctg == nil || *stats.SavePctg != 0.967 {
		t.Errorf("SavePctg = %v, want 0.967", stats.SavePctg)
	}
	if stats.GoalsAgainst != 1 {
		t.Errorf("GoalsAgainst = %d, want 1", stats.GoalsAgainst)
	}
	if stats.Saves != 30 {
		t.Errorf("Saves = %d, want 30", stats.Saves)
	}
	if stats.ShotsAgainst != 31 {
		t.Errorf("ShotsAgainst = %d, want 31", stats.ShotsAgainst)
	}
	if stats.Starter == nil || !*stats.Starter {
		t.Errorf("Starter = %v, want true", stats.Starter)
	}
	if stats.Decision == nil || *stats.Decision != GoalieDecisionWin {
		t.Errorf("Decision = %v, want %v", stats.Decision, GoalieDecisionWin)
	}
}

func TestTVBroadcast_Deserialization(t *testing.T) {
	jsonData := `{
		"id": 123,
		"market": "NATIONAL",
		"countryCode": "US",
		"network": "ESPN",
		"sequenceNumber": 1
	}`

	var broadcast TVBroadcast
	err := json.Unmarshal([]byte(jsonData), &broadcast)
	if err != nil {
		t.Fatalf("Failed to unmarshal TV broadcast: %v", err)
	}

	if broadcast.ID != 123 {
		t.Errorf("ID = %d, want 123", broadcast.ID)
	}
	if broadcast.Market != "NATIONAL" {
		t.Errorf("Market = %s, want NATIONAL", broadcast.Market)
	}
	if broadcast.CountryCode != "US" {
		t.Errorf("CountryCode = %s, want US", broadcast.CountryCode)
	}
	if broadcast.Network != "ESPN" {
		t.Errorf("Network = %s, want ESPN", broadcast.Network)
	}
	if broadcast.SequenceNumber != 1 {
		t.Errorf("SequenceNumber = %d, want 1", broadcast.SequenceNumber)
	}
}

func TestSpecialEvent_Deserialization(t *testing.T) {
	jsonData := `{
		"parentId": 999,
		"name": {"default": "Winter Classic"},
		"lightLogoUrl": {"default": "https://example.com/logo.png"}
	}`

	var event SpecialEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal special event: %v", err)
	}

	if event.ParentID != 999 {
		t.Errorf("ParentID = %d, want 999", event.ParentID)
	}
	if event.Name.Default != "Winter Classic" {
		t.Errorf("Name.Default = %s, want Winter Classic", event.Name.Default)
	}
	if event.LightLogoURL.Default != "https://example.com/logo.png" {
		t.Errorf("LightLogoURL.Default = %s, want https://example.com/logo.png", event.LightLogoURL.Default)
	}
}

func TestPeriodDescriptor_Deserialization(t *testing.T) {
	tests := []struct {
		name       string
		jsonData   string
		wantNumber int
		wantType   PeriodType
	}{
		{
			name: "regulation period",
			jsonData: `{
				"number": 3,
				"periodType": "REG",
				"maxRegulationPeriods": 3
			}`,
			wantNumber: 3,
			wantType:   PeriodTypeRegulation,
		},
		{
			name: "overtime period",
			jsonData: `{
				"number": 4,
				"periodType": "OT",
				"maxRegulationPeriods": 3
			}`,
			wantNumber: 4,
			wantType:   PeriodTypeOvertime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var period PeriodDescriptor
			err := json.Unmarshal([]byte(tt.jsonData), &period)
			if err != nil {
				t.Fatalf("Failed to unmarshal period descriptor: %v", err)
			}

			if period.Number != tt.wantNumber {
				t.Errorf("Number = %d, want %d", period.Number, tt.wantNumber)
			}
			if period.PeriodType != tt.wantType {
				t.Errorf("PeriodType = %v, want %v", period.PeriodType, tt.wantType)
			}
			if period.MaxRegulationPeriods != 3 {
				t.Errorf("MaxRegulationPeriods = %d, want 3", period.MaxRegulationPeriods)
			}
		})
	}
}

func TestBoxscoreTeam_Deserialization(t *testing.T) {
	jsonData := `{
		"id": 8,
		"commonName": {"default": "Canadiens"},
		"abbrev": "MTL",
		"score": 3,
		"sog": 28,
		"logo": "https://assets.nhle.com/logos/nhl/svg/MTL_light.svg",
		"darkLogo": "https://assets.nhle.com/logos/nhl/svg/MTL_dark.svg",
		"placeName": {"default": "Montréal"},
		"placeNameWithPreposition": {"default": "Montréal"}
	}`

	var team BoxscoreTeam
	err := json.Unmarshal([]byte(jsonData), &team)
	if err != nil {
		t.Fatalf("Failed to unmarshal boxscore team: %v", err)
	}

	if team.ID != TeamID(8) {
		t.Errorf("ID = %d, want 8", team.ID)
	}
	if team.CommonName.Default != "Canadiens" {
		t.Errorf("CommonName.Default = %s, want Canadiens", team.CommonName.Default)
	}
	if team.Abbrev != "MTL" {
		t.Errorf("Abbrev = %s, want MTL", team.Abbrev)
	}
	if team.Score != 3 {
		t.Errorf("Score = %d, want 3", team.Score)
	}
	if team.SOG != 28 {
		t.Errorf("SOG = %d, want 28", team.SOG)
	}
}

func TestGameClock_Deserialization(t *testing.T) {
	tests := []struct {
		name              string
		jsonData          string
		wantTimeRemaining string
		wantSecondsRemain int
		wantRunning       bool
		wantIntermission  bool
	}{
		{
			name: "running clock",
			jsonData: `{
				"timeRemaining": "05:30",
				"secondsRemaining": 330,
				"running": true,
				"inIntermission": false
			}`,
			wantTimeRemaining: "05:30",
			wantSecondsRemain: 330,
			wantRunning:       true,
			wantIntermission:  false,
		},
		{
			name: "end of period",
			jsonData: `{
				"timeRemaining": "00:00",
				"secondsRemaining": 0,
				"running": false,
				"inIntermission": true
			}`,
			wantTimeRemaining: "00:00",
			wantSecondsRemain: 0,
			wantRunning:       false,
			wantIntermission:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var clock GameClock
			err := json.Unmarshal([]byte(tt.jsonData), &clock)
			if err != nil {
				t.Fatalf("Failed to unmarshal game clock: %v", err)
			}

			if clock.TimeRemaining != tt.wantTimeRemaining {
				t.Errorf("TimeRemaining = %s, want %s", clock.TimeRemaining, tt.wantTimeRemaining)
			}
			if clock.SecondsRemaining != tt.wantSecondsRemain {
				t.Errorf("SecondsRemaining = %d, want %d", clock.SecondsRemaining, tt.wantSecondsRemain)
			}
			if clock.Running != tt.wantRunning {
				t.Errorf("Running = %v, want %v", clock.Running, tt.wantRunning)
			}
			if clock.InIntermission != tt.wantIntermission {
				t.Errorf("InIntermission = %v, want %v", clock.InIntermission, tt.wantIntermission)
			}
		})
	}
}

func TestBoxscore_WithSpecialEvent(t *testing.T) {
	jsonData := `{
		"id": 2024020001,
		"season": 20242025,
		"gameType": 2,
		"limitedScoring": false,
		"gameDate": "2024-10-04",
		"venue": {"default": "Test Arena"},
		"venueLocation": {"default": "Test City"},
		"startTimeUTC": "2024-10-04T19:00:00Z",
		"easternUTCOffset": "-04:00",
		"venueUTCOffset": "-04:00",
		"tvBroadcasts": [],
		"gameState": "LIVE",
		"gameScheduleState": "OK",
		"periodDescriptor": {
			"number": 1,
			"periodType": "REG",
			"maxRegulationPeriods": 3
		},
		"specialEvent": {
			"parentId": 1000,
			"name": {"default": "Stadium Series"},
			"lightLogoUrl": {"default": "https://example.com/stadium.png"}
		},
		"awayTeam": {
			"id": 1,
			"commonName": {"default": "Devils"},
			"abbrev": "NJD",
			"score": 0,
			"sog": 0,
			"logo": "https://assets.nhle.com/logos/nhl/svg/NJD_light.svg",
			"darkLogo": "https://assets.nhle.com/logos/nhl/svg/NJD_dark.svg",
			"placeName": {"default": "New Jersey"},
			"placeNameWithPreposition": {"default": "New Jersey"}
		},
		"homeTeam": {
			"id": 7,
			"commonName": {"default": "Sabres"},
			"abbrev": "BUF",
			"score": 0,
			"sog": 0,
			"logo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg",
			"darkLogo": "https://assets.nhle.com/logos/nhl/svg/BUF_dark.svg",
			"placeName": {"default": "Buffalo"},
			"placeNameWithPreposition": {"default": "Buffalo"}
		},
		"clock": {
			"timeRemaining": "20:00",
			"secondsRemaining": 1200,
			"running": false,
			"inIntermission": false
		},
		"playerByGameStats": {
			"awayTeam": {
				"forwards": [],
				"defense": [],
				"goalies": []
			},
			"homeTeam": {
				"forwards": [],
				"defense": [],
				"goalies": []
			}
		}
	}`

	var boxscore Boxscore
	err := json.Unmarshal([]byte(jsonData), &boxscore)
	if err != nil {
		t.Fatalf("Failed to unmarshal boxscore: %v", err)
	}

	if boxscore.SpecialEvent == nil {
		t.Fatal("SpecialEvent is nil, want non-nil")
	}
	if boxscore.SpecialEvent.Name.Default != "Stadium Series" {
		t.Errorf("SpecialEvent.Name.Default = %s, want Stadium Series", boxscore.SpecialEvent.Name.Default)
	}
}

func TestBoxscore_WithTVBroadcasts(t *testing.T) {
	jsonData := `{
		"id": 2024020001,
		"season": 20242025,
		"gameType": 2,
		"limitedScoring": false,
		"gameDate": "2024-10-04",
		"venue": {"default": "Test Arena"},
		"venueLocation": {"default": "Test City"},
		"startTimeUTC": "2024-10-04T19:00:00Z",
		"easternUTCOffset": "-04:00",
		"venueUTCOffset": "-04:00",
		"tvBroadcasts": [
			{
				"id": 1,
				"market": "NATIONAL",
				"countryCode": "US",
				"network": "ESPN",
				"sequenceNumber": 1
			},
			{
				"id": 2,
				"market": "AWAY",
				"countryCode": "US",
				"network": "MSG",
				"sequenceNumber": 2
			}
		],
		"gameState": "LIVE",
		"gameScheduleState": "OK",
		"periodDescriptor": {
			"number": 1,
			"periodType": "REG",
			"maxRegulationPeriods": 3
		},
		"awayTeam": {
			"id": 1,
			"commonName": {"default": "Devils"},
			"abbrev": "NJD",
			"score": 0,
			"sog": 0,
			"logo": "https://assets.nhle.com/logos/nhl/svg/NJD_light.svg",
			"darkLogo": "https://assets.nhle.com/logos/nhl/svg/NJD_dark.svg",
			"placeName": {"default": "New Jersey"},
			"placeNameWithPreposition": {"default": "New Jersey"}
		},
		"homeTeam": {
			"id": 7,
			"commonName": {"default": "Sabres"},
			"abbrev": "BUF",
			"score": 0,
			"sog": 0,
			"logo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg",
			"darkLogo": "https://assets.nhle.com/logos/nhl/svg/BUF_dark.svg",
			"placeName": {"default": "Buffalo"},
			"placeNameWithPreposition": {"default": "Buffalo"}
		},
		"clock": {
			"timeRemaining": "20:00",
			"secondsRemaining": 1200,
			"running": false,
			"inIntermission": false
		},
		"playerByGameStats": {
			"awayTeam": {
				"forwards": [],
				"defense": [],
				"goalies": []
			},
			"homeTeam": {
				"forwards": [],
				"defense": [],
				"goalies": []
			}
		}
	}`

	var boxscore Boxscore
	err := json.Unmarshal([]byte(jsonData), &boxscore)
	if err != nil {
		t.Fatalf("Failed to unmarshal boxscore: %v", err)
	}

	if len(boxscore.TVBroadcasts) != 2 {
		t.Fatalf("len(TVBroadcasts) = %d, want 2", len(boxscore.TVBroadcasts))
	}
	if boxscore.TVBroadcasts[0].Network != "ESPN" {
		t.Errorf("TVBroadcasts[0].Network = %s, want ESPN", boxscore.TVBroadcasts[0].Network)
	}
	if boxscore.TVBroadcasts[1].Network != "MSG" {
		t.Errorf("TVBroadcasts[1].Network = %s, want MSG", boxscore.TVBroadcasts[1].Network)
	}
}

func TestGoalieStats_MissingOptionalFields(t *testing.T) {
	jsonData := `{
		"playerId": 8475123,
		"sweaterNumber": 30,
		"name": {"default": "J. Doe"},
		"position": "G",
		"evenStrengthShotsAgainst": "0/0",
		"powerPlayShotsAgainst": "0/0",
		"shorthandedShotsAgainst": "0/0",
		"saveShotsAgainst": "0/0",
		"evenStrengthGoalsAgainst": 0,
		"powerPlayGoalsAgainst": 0,
		"shorthandedGoalsAgainst": 0,
		"goalsAgainst": 0,
		"toi": "00:00",
		"shotsAgainst": 0,
		"saves": 0
	}`

	var stats GoalieStats
	err := json.Unmarshal([]byte(jsonData), &stats)
	if err != nil {
		t.Fatalf("Failed to unmarshal goalie stats: %v", err)
	}

	if stats.PlayerID != PlayerID(8475123) {
		t.Errorf("PlayerID = %d, want 8475123", stats.PlayerID)
	}
	if stats.SavePctg != nil {
		t.Errorf("SavePctg = %v, want nil", stats.SavePctg)
	}
	if stats.PIM != nil {
		t.Errorf("PIM = %v, want nil", stats.PIM)
	}
	if stats.Starter != nil {
		t.Errorf("Starter = %v, want nil", stats.Starter)
	}
	if stats.Decision != nil {
		t.Errorf("Decision = %v, want nil", stats.Decision)
	}
}

func TestTeamPlayerStats_Deserialization(t *testing.T) {
	jsonData := `{
		"forwards": [
			{
				"playerId": 8480002,
				"sweaterNumber": 13,
				"name": {"default": "N. Hischier"},
				"position": "C",
				"goals": 1,
				"assists": 2,
				"points": 3,
				"plusMinus": 2,
				"pim": 0,
				"hits": 3,
				"powerPlayGoals": 1,
				"sog": 4,
				"faceoffWinningPctg": 0.55,
				"toi": "18:15",
				"blockedShots": 1,
				"shifts": 27,
				"giveaways": 0,
				"takeaways": 2
			}
		],
		"defense": [],
		"goalies": []
	}`

	var stats TeamPlayerStats
	err := json.Unmarshal([]byte(jsonData), &stats)
	if err != nil {
		t.Fatalf("Failed to unmarshal team player stats: %v", err)
	}

	if len(stats.Forwards) != 1 {
		t.Fatalf("len(Forwards) = %d, want 1", len(stats.Forwards))
	}
	if len(stats.Defense) != 0 {
		t.Errorf("len(Defense) = %d, want 0", len(stats.Defense))
	}
	if len(stats.Goalies) != 0 {
		t.Errorf("len(Goalies) = %d, want 0", len(stats.Goalies))
	}
	if stats.Forwards[0].PlayerID != PlayerID(8480002) {
		t.Errorf("Forwards[0].PlayerID = %d, want 8480002", stats.Forwards[0].PlayerID)
	}
}

func TestTeamPlayerStats_EmptyArrays(t *testing.T) {
	jsonData := `{
		"forwards": [],
		"defense": [],
		"goalies": []
	}`

	var stats TeamPlayerStats
	err := json.Unmarshal([]byte(jsonData), &stats)
	if err != nil {
		t.Fatalf("Failed to unmarshal team player stats: %v", err)
	}

	if len(stats.Forwards) != 0 {
		t.Errorf("len(Forwards) = %d, want 0", len(stats.Forwards))
	}
	if len(stats.Defense) != 0 {
		t.Errorf("len(Defense) = %d, want 0", len(stats.Defense))
	}
	if len(stats.Goalies) != 0 {
		t.Errorf("len(Goalies) = %d, want 0", len(stats.Goalies))
	}
}

func TestTeamGameStats_FromEmptyTeam(t *testing.T) {
	teamStats := TeamPlayerStats{
		Forwards: []SkaterStats{},
		Defense:  []SkaterStats{},
		Goalies:  []GoalieStats{},
	}

	gameStats := FromTeamPlayerStats(&teamStats)

	if gameStats.ShotsOnGoal != 0 {
		t.Errorf("ShotsOnGoal = %d, want 0", gameStats.ShotsOnGoal)
	}
	if gameStats.Hits != 0 {
		t.Errorf("Hits = %d, want 0", gameStats.Hits)
	}
	if gameStats.PenaltyMinutes != 0 {
		t.Errorf("PenaltyMinutes = %d, want 0", gameStats.PenaltyMinutes)
	}
}

func TestTeamGameStats_FromSkaters(t *testing.T) {
	teamStats := TeamPlayerStats{
		Forwards: []SkaterStats{
			{
				PlayerID:           PlayerID(1),
				SweaterNumber:      13,
				Name:               LocalizedString{Default: "Player 1"},
				Position:           PositionCenter,
				Goals:              1,
				Assists:            2,
				Points:             3,
				PlusMinus:          1,
				PIM:                2,
				Hits:               5,
				PowerPlayGoals:     1,
				SOG:                4,
				FaceoffWinningPctg: 0.6,
				TOI:                "18:00",
				BlockedShots:       2,
				Shifts:             25,
				Giveaways:          1,
				Takeaways:          3,
			},
		},
		Defense: []SkaterStats{
			{
				PlayerID:           PlayerID(2),
				SweaterNumber:      44,
				Name:               LocalizedString{Default: "Player 2"},
				Position:           PositionDefense,
				Goals:              0,
				Assists:            1,
				Points:             1,
				PlusMinus:          0,
				PIM:                4,
				Hits:               8,
				PowerPlayGoals:     0,
				SOG:                3,
				FaceoffWinningPctg: 0.0,
				TOI:                "22:00",
				BlockedShots:       5,
				Shifts:             30,
				Giveaways:          2,
				Takeaways:          1,
			},
		},
		Goalies: []GoalieStats{},
	}

	gameStats := FromTeamPlayerStats(&teamStats)

	if gameStats.ShotsOnGoal != 7 {
		t.Errorf("ShotsOnGoal = %d, want 7", gameStats.ShotsOnGoal)
	}
	if gameStats.Hits != 13 {
		t.Errorf("Hits = %d, want 13", gameStats.Hits)
	}
	if gameStats.PenaltyMinutes != 6 {
		t.Errorf("PenaltyMinutes = %d, want 6", gameStats.PenaltyMinutes)
	}
	if gameStats.PowerPlayGoals != 1 {
		t.Errorf("PowerPlayGoals = %d, want 1", gameStats.PowerPlayGoals)
	}
	if gameStats.BlockedShots != 7 {
		t.Errorf("BlockedShots = %d, want 7", gameStats.BlockedShots)
	}
	if gameStats.Giveaways != 3 {
		t.Errorf("Giveaways = %d, want 3", gameStats.Giveaways)
	}
	if gameStats.Takeaways != 4 {
		t.Errorf("Takeaways = %d, want 4", gameStats.Takeaways)
	}
}

func TestTeamGameStats_WithGoalies(t *testing.T) {
	pim := 2
	teamStats := TeamPlayerStats{
		Forwards: []SkaterStats{},
		Defense:  []SkaterStats{},
		Goalies: []GoalieStats{
			{
				PlayerID:                 PlayerID(1),
				SweaterNumber:            35,
				Name:                     LocalizedString{Default: "Goalie 1"},
				Position:                 PositionGoalie,
				EvenStrengthShotsAgainst: "20/22",
				PowerPlayShotsAgainst:    "3/5",
				ShorthandedShotsAgainst:  "0/0",
				SaveShotsAgainst:         "23/27",
				SavePctg:                 floatPtr(0.852),
				EvenStrengthGoalsAgainst: 2,
				PowerPlayGoalsAgainst:    2,
				ShorthandedGoalsAgainst:  0,
				PIM:                      &pim,
				GoalsAgainst:             4,
				TOI:                      "60:00",
				Starter:                  boolPtr(true),
				Decision:                 goalieDecisionPtr(GoalieDecisionLoss),
				ShotsAgainst:             27,
				Saves:                    23,
			},
		},
	}

	gameStats := FromTeamPlayerStats(&teamStats)

	if gameStats.PenaltyMinutes != 2 {
		t.Errorf("PenaltyMinutes = %d, want 2", gameStats.PenaltyMinutes)
	}
	if gameStats.PowerPlayOpportunities != 2 {
		t.Errorf("PowerPlayOpportunities = %d, want 2", gameStats.PowerPlayOpportunities)
	}
}

func TestTeamGameStats_FaceoffPercentage_ZeroFaceoffs(t *testing.T) {
	gameStats := TeamGameStats{
		ShotsOnGoal:            30,
		FaceoffWins:            0,
		FaceoffTotal:           0,
		PowerPlayGoals:         1,
		PowerPlayOpportunities: 4,
		PenaltyMinutes:         8,
		Hits:                   25,
		BlockedShots:           15,
		Giveaways:              5,
		Takeaways:              7,
	}

	got := gameStats.FaceoffPercentage()
	if got != 0.0 {
		t.Errorf("FaceoffPercentage() = %f, want 0.0", got)
	}
}

func TestTeamGameStats_FaceoffPercentage(t *testing.T) {
	gameStats := TeamGameStats{
		ShotsOnGoal:            30,
		FaceoffWins:            30,
		FaceoffTotal:           60,
		PowerPlayGoals:         1,
		PowerPlayOpportunities: 4,
		PenaltyMinutes:         8,
		Hits:                   25,
		BlockedShots:           15,
		Giveaways:              5,
		Takeaways:              7,
	}

	got := gameStats.FaceoffPercentage()
	want := 50.0
	if got != want {
		t.Errorf("FaceoffPercentage() = %f, want %f", got, want)
	}
}

func TestTeamGameStats_PowerPlayPercentage_ZeroOpportunities(t *testing.T) {
	gameStats := TeamGameStats{
		ShotsOnGoal:            30,
		FaceoffWins:            30,
		FaceoffTotal:           60,
		PowerPlayGoals:         0,
		PowerPlayOpportunities: 0,
		PenaltyMinutes:         8,
		Hits:                   25,
		BlockedShots:           15,
		Giveaways:              5,
		Takeaways:              7,
	}

	got := gameStats.PowerPlayPercentage()
	if got != 0.0 {
		t.Errorf("PowerPlayPercentage() = %f, want 0.0", got)
	}
}

func TestTeamGameStats_PowerPlayPercentage(t *testing.T) {
	gameStats := TeamGameStats{
		ShotsOnGoal:            30,
		FaceoffWins:            30,
		FaceoffTotal:           60,
		PowerPlayGoals:         2,
		PowerPlayOpportunities: 5,
		PenaltyMinutes:         8,
		Hits:                   25,
		BlockedShots:           15,
		Giveaways:              5,
		Takeaways:              7,
	}

	got := gameStats.PowerPlayPercentage()
	want := 40.0
	if got != want {
		t.Errorf("PowerPlayPercentage() = %f, want %f", got, want)
	}
}

func TestSkaterStats_ZeroValues(t *testing.T) {
	jsonData := `{
		"playerId": 8480000,
		"sweaterNumber": 99,
		"name": {"default": "Test Player"},
		"position": "LW",
		"goals": 0,
		"assists": 0,
		"points": 0,
		"plusMinus": 0,
		"pim": 0,
		"hits": 0,
		"powerPlayGoals": 0,
		"sog": 0,
		"faceoffWinningPctg": 0.0,
		"toi": "00:00",
		"blockedShots": 0,
		"shifts": 0,
		"giveaways": 0,
		"takeaways": 0
	}`

	var stats SkaterStats
	err := json.Unmarshal([]byte(jsonData), &stats)
	if err != nil {
		t.Fatalf("Failed to unmarshal skater stats: %v", err)
	}

	if stats.Goals != 0 {
		t.Errorf("Goals = %d, want 0", stats.Goals)
	}
	if stats.Assists != 0 {
		t.Errorf("Assists = %d, want 0", stats.Assists)
	}
	if stats.SOG != 0 {
		t.Errorf("SOG = %d, want 0", stats.SOG)
	}
	if stats.FaceoffWinningPctg != 0.0 {
		t.Errorf("FaceoffWinningPctg = %f, want 0.0", stats.FaceoffWinningPctg)
	}
}

func TestSkaterStats_NegativePlusMinus(t *testing.T) {
	jsonData := `{
		"playerId": 8480000,
		"sweaterNumber": 99,
		"name": {"default": "Test Player"},
		"position": "RW",
		"goals": 1,
		"assists": 0,
		"points": 1,
		"plusMinus": -3,
		"pim": 2,
		"hits": 5,
		"powerPlayGoals": 0,
		"sog": 3,
		"faceoffWinningPctg": 0.0,
		"toi": "12:30",
		"blockedShots": 0,
		"shifts": 18,
		"giveaways": 2,
		"takeaways": 0
	}`

	var stats SkaterStats
	err := json.Unmarshal([]byte(jsonData), &stats)
	if err != nil {
		t.Fatalf("Failed to unmarshal skater stats: %v", err)
	}

	if stats.PlusMinus != -3 {
		t.Errorf("PlusMinus = %d, want -3", stats.PlusMinus)
	}
}

func TestBoxscore_RoundTripJSON(t *testing.T) {
	original := Boxscore{
		ID:                GameID(2024020001),
		Season:            NewSeason(2024),
		GameType:          GameTypeRegularSeason,
		LimitedScoring:    false,
		GameDate:          "2024-10-04",
		Venue:             LocalizedString{Default: "Test Arena"},
		VenueLocation:     LocalizedString{Default: "Test City"},
		StartTimeUTC:      "2024-10-04T19:00:00Z",
		EasternUTCOffset:  "-04:00",
		VenueUTCOffset:    "-04:00",
		TVBroadcasts:      []TVBroadcast{},
		GameState:         GameStateLive,
		GameScheduleState: GameScheduleStateOK,
		PeriodDescriptor: PeriodDescriptor{
			Number:               2,
			PeriodType:           PeriodTypeRegulation,
			MaxRegulationPeriods: 3,
		},
		AwayTeam: BoxscoreTeam{
			ID:         TeamID(1),
			CommonName: LocalizedString{Default: "Devils"},
			Abbrev:     "NJD",
			Score:      2,
			SOG:        15,
		},
		HomeTeam: BoxscoreTeam{
			ID:         TeamID(7),
			CommonName: LocalizedString{Default: "Sabres"},
			Abbrev:     "BUF",
			Score:      1,
			SOG:        12,
		},
		Clock: GameClock{
			TimeRemaining:    "10:15",
			SecondsRemaining: 615,
			Running:          true,
			InIntermission:   false,
		},
		PlayerByGameStats: PlayerByGameStats{
			AwayTeam: TeamPlayerStats{},
			HomeTeam: TeamPlayerStats{},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal boxscore: %v", err)
	}

	var decoded Boxscore
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal boxscore: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID = %d, want %d", decoded.ID, original.ID)
	}
	if decoded.GameType != original.GameType {
		t.Errorf("GameType = %v, want %v", decoded.GameType, original.GameType)
	}
	if decoded.GameState != original.GameState {
		t.Errorf("GameState = %v, want %v", decoded.GameState, original.GameState)
	}
}

// Helper functions for creating pointers to values
func floatPtr(f float64) *float64 {
	return &f
}

func boolPtr(b bool) *bool {
	return &b
}

func goalieDecisionPtr(g GoalieDecision) *GoalieDecision {
	return &g
}
