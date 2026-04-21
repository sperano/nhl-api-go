package nhl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ===== Deserialization Tests =====

func TestEdgeSkaterDetail_Deserialization(t *testing.T) {
	jsonData := `{
		"player": {
			"id": 8478402,
			"firstName": {"default": "Connor"},
			"lastName": {"default": "McDavid"},
			"birthDate": "1997-01-13",
			"shootsCatches": "L",
			"sweaterNumber": 97,
			"position": "C",
			"slug": "connor-mcdavid-8478402",
			"headshot": "https://assets.nhle.com/mugs/nhl/20242025/EDM/8478402.png",
			"goals": 30,
			"assists": 60,
			"points": 90,
			"gamesPlayed": 50,
			"team": {
				"id": 22,
				"commonName": {"default": "Oilers"},
				"placeNameWithPreposition": {"default": "Edmonton"},
				"abbrev": "EDM",
				"teamLogo": {"light": "https://assets.nhle.com/logos/nhl/svg/EDM_light.svg", "dark": "https://assets.nhle.com/logos/nhl/svg/EDM_dark.svg"},
				"slug": "edmonton-oilers",
				"conference": "Western",
				"division": "Pacific",
				"wins": 35,
				"losses": 15,
				"otLosses": 5,
				"gamesPlayed": 55,
				"points": 75
			}
		},
		"seasonsWithEdgeStats": [{"id": 20242025, "gameTypes": [2]}, {"id": 20232024, "gameTypes": [2, 3]}],
		"topShotSpeed": {
			"imperial": 102.3,
			"metric": 164.6,
			"percentile": 0.95,
			"leagueAvg": {"imperial": 85.0, "metric": 136.8},
			"overlay": {
				"player": {"firstName": {"default": "Connor"}, "lastName": {"default": "McDavid"}},
				"gameDate": "2025-01-15",
				"awayTeam": {"abbrev": "CGY", "score": 2},
				"homeTeam": {"abbrev": "EDM", "score": 5},
				"periodDescriptor": {"number": 2, "periodType": "REG", "maxRegulationPeriods": 3},
				"timeInPeriod": "14:32",
				"gameType": 2
			}
		},
		"skatingSpeed": {
			"speedMax": {
				"imperial": 23.1,
				"metric": 37.2,
				"percentile": 0.98,
				"leagueAvg": {"imperial": 21.5, "metric": 34.6},
				"overlay": {
					"player": {"firstName": {"default": "Connor"}, "lastName": {"default": "McDavid"}},
					"gameDate": "2025-02-01",
					"awayTeam": {"abbrev": "EDM", "score": 4},
					"homeTeam": {"abbrev": "VAN", "score": 1},
					"periodDescriptor": {"number": 3, "periodType": "REG", "maxRegulationPeriods": 3},
					"timeInPeriod": "08:15",
					"gameType": 2
				}
			},
			"burstsOver20": {"value": 150, "percentile": 0.92, "leagueAvg": {"value": 110.5}}
		},
		"totalDistanceSkated": {"imperial": 450.2, "metric": 724.5, "percentile": 0.88, "leagueAvg": {"imperial": 400.0, "metric": 643.7}},
		"distanceMaxGame": {
			"imperial": 12.5,
			"metric": 20.1,
			"percentile": 0.91,
			"leagueAvg": {"imperial": 10.0, "metric": 16.1}
		},
		"sogSummary": [
			{
				"locationCode": "all",
				"shots": 200,
				"shotsPercentile": 0.90,
				"shotsLeagueAvg": 150.0,
				"goals": 30,
				"goalsPercentile": 0.95,
				"goalsLeagueAvg": 20.0,
				"shootingPctg": 0.15,
				"shootingPctgPercentile": 0.88,
				"shootingPctgLeagueAvg": 0.12
			}
		],
		"sogDetails": [
			{"area": "Crease", "shots": 40, "shootingPctg": 0.25, "shotsPercentile": 0.85}
		],
		"zoneTimeDetails": {
			"offensiveZonePctg": 0.35,
			"offensiveZonePercentile": 0.80,
			"offensiveZoneLeagueAvg": 0.30,
			"offensiveZoneEvPctg": 0.33,
			"offensiveZoneEvPercentile": 0.78,
			"offensiveZoneEvLeagueAvg": 0.29,
			"neutralZonePctg": 0.35,
			"neutralZonePercentile": 0.50,
			"neutralZoneLeagueAvg": 0.36,
			"defensiveZonePctg": 0.30,
			"defensiveZonePercentile": 0.70,
			"defensiveZoneLeagueAvg": 0.34
		}
	}`

	var detail EdgeSkaterDetail
	if err := json.Unmarshal([]byte(jsonData), &detail); err != nil {
		t.Fatalf("Failed to unmarshal EdgeSkaterDetail: %v", err)
	}

	if detail.Player.ID != 8478402 {
		t.Errorf("Player.ID = %d, want 8478402", detail.Player.ID)
	}
	if detail.Player.FirstName.Default != "Connor" {
		t.Errorf("Player.FirstName = %q, want %q", detail.Player.FirstName.Default, "Connor")
	}
	if detail.Player.Team.Abbrev != "EDM" {
		t.Errorf("Player.Team.Abbrev = %q, want %q", detail.Player.Team.Abbrev, "EDM")
	}
	if len(detail.SeasonsWithEdgeStats) != 2 {
		t.Errorf("SeasonsWithEdgeStats length = %d, want 2", len(detail.SeasonsWithEdgeStats))
	}
	if detail.TopShotSpeed.Imperial != 102.3 {
		t.Errorf("TopShotSpeed.Imperial = %f, want 102.3", detail.TopShotSpeed.Imperial)
	}
	if detail.TopShotSpeed.Overlay == nil {
		t.Fatal("TopShotSpeed.Overlay is nil, want non-nil")
	}
	if detail.TopShotSpeed.Overlay.HomeTeam.Abbrev != "EDM" {
		t.Errorf("TopShotSpeed.Overlay.HomeTeam.Abbrev = %q, want %q", detail.TopShotSpeed.Overlay.HomeTeam.Abbrev, "EDM")
	}
	if detail.SkatingSpeed.BurstsOver20.Value != 150 {
		t.Errorf("SkatingSpeed.BurstsOver20.Value = %d, want 150", detail.SkatingSpeed.BurstsOver20.Value)
	}
	if detail.SkatingSpeed.BurstsOver20.LeagueAvg.Value != 110.5 {
		t.Errorf("SkatingSpeed.BurstsOver20.LeagueAvg.Value = %f, want 110.5", detail.SkatingSpeed.BurstsOver20.LeagueAvg.Value)
	}
	if len(detail.SogSummary) != 1 {
		t.Errorf("SogSummary length = %d, want 1", len(detail.SogSummary))
	}
	if detail.ZoneTimeDetails.OffensiveZonePctg != 0.35 {
		t.Errorf("ZoneTimeDetails.OffensiveZonePctg = %f, want 0.35", detail.ZoneTimeDetails.OffensiveZonePctg)
	}
}

func TestEdgeGoalieDetail_Deserialization(t *testing.T) {
	jsonData := `{
		"player": {
			"id": 8479318,
			"firstName": {"default": "Igor"},
			"lastName": {"default": "Shesterkin"},
			"birthDate": "1995-12-30",
			"shootsCatches": "L",
			"sweaterNumber": 31,
			"slug": "igor-shesterkin-8479318",
			"headshot": "https://assets.nhle.com/mugs/nhl/20242025/NYR/8479318.png",
			"wins": 25,
			"losses": 10,
			"overtimeLosses": 3,
			"goalsAgainstAvg": 2.15,
			"savePctg": 0.928,
			"gamesPlayed": 38,
			"team": {
				"id": 3,
				"commonName": {"default": "Rangers"},
				"placeNameWithPreposition": {"default": "New York"},
				"abbrev": "NYR",
				"teamLogo": {"light": "https://assets.nhle.com/logos/nhl/svg/NYR_light.svg", "dark": "https://assets.nhle.com/logos/nhl/svg/NYR_dark.svg"},
				"slug": "new-york-rangers",
				"conference": "Eastern",
				"division": "Metropolitan",
				"wins": 30,
				"losses": 18,
				"otLosses": 7,
				"gamesPlayed": 55,
				"points": 67
			}
		},
		"seasonsWithEdgeStats": [{"id": 20242025, "gameTypes": [2]}],
		"stats": {
			"goalsAgainstAvg": {"value": 2.15, "percentile": 0.90, "leagueAvg": 2.85},
			"gamesAbove900": {"value": 28.0, "percentile": 0.92, "leagueAvg": 18.5},
			"goalDifferentialPer60": {"value": 1.2, "percentile": 0.88, "leagueAvg": 0.5},
			"goalSupportAvg": {"value": 3.1, "percentile": 0.65, "leagueAvg": 3.0},
			"pointPctg": {"value": 0.68, "percentile": 0.85, "leagueAvg": 0.55}
		},
		"shotLocationSummary": [
			{
				"locationCode": "all",
				"goalsAgainst": 80,
				"goalsAgainstPercentile": 0.85,
				"goalsAgainstLeagueAvg": 100.0,
				"saves": 1000,
				"savesPercentile": 0.90,
				"savesLeagueAvg": 850.0,
				"savePctg": 0.926,
				"savePctgPercentile": 0.88,
				"savePctgLeagueAvg": 0.905
			}
		],
		"shotLocationDetails": [
			{"area": "Crease", "saves": 200, "savesPercentile": 0.85, "savePctg": 0.88, "savePctgPercentile": 0.80}
		]
	}`

	var detail EdgeGoalieDetail
	if err := json.Unmarshal([]byte(jsonData), &detail); err != nil {
		t.Fatalf("Failed to unmarshal EdgeGoalieDetail: %v", err)
	}

	if detail.Player.ID != 8479318 {
		t.Errorf("Player.ID = %d, want 8479318", detail.Player.ID)
	}
	if detail.Player.SavePctg != 0.928 {
		t.Errorf("Player.SavePctg = %f, want 0.928", detail.Player.SavePctg)
	}
	if detail.Stats.GoalsAgainstAvg.Value != 2.15 {
		t.Errorf("Stats.GoalsAgainstAvg.Value = %f, want 2.15", detail.Stats.GoalsAgainstAvg.Value)
	}
	if len(detail.ShotLocationSummary) != 1 {
		t.Errorf("ShotLocationSummary length = %d, want 1", len(detail.ShotLocationSummary))
	}
	if detail.ShotLocationSummary[0].SavePctg != 0.926 {
		t.Errorf("ShotLocationSummary[0].SavePctg = %f, want 0.926", detail.ShotLocationSummary[0].SavePctg)
	}
}

func TestEdgeTeamDetail_Deserialization(t *testing.T) {
	jsonData := `{
		"team": {
			"id": 22,
			"commonName": {"default": "Oilers"},
			"placeNameWithPreposition": {"default": "Edmonton"},
			"abbrev": "EDM",
			"teamLogo": {"light": "https://assets.nhle.com/logos/nhl/svg/EDM_light.svg", "dark": "https://assets.nhle.com/logos/nhl/svg/EDM_dark.svg"},
			"slug": "edmonton-oilers",
			"conference": "Western",
			"division": "Pacific",
			"wins": 35,
			"losses": 15,
			"otLosses": 5,
			"gamesPlayed": 55,
			"points": 75
		},
		"seasonsWithEdgeStats": [{"id": 20242025, "gameTypes": [2]}],
		"shotSpeed": {
			"shotAttemptsOver90": {"value": 120, "rank": 3, "leagueAvg": {"value": 95.5}},
			"topShotSpeed": {
				"imperial": 105.2,
				"metric": 169.3,
				"rank": 5,
				"leagueAvg": {"imperial": 98.0, "metric": 157.7}
			}
		},
		"skatingSpeed": {
			"burstsOver22": {"value": 80, "rank": 2},
			"burstsOver20": {"value": 500, "rank": 8, "leagueAvg": {"value": 420.0}},
			"speedMax": {
				"imperial": 24.5,
				"metric": 39.4,
				"rank": 1,
				"leagueAvg": {"imperial": 22.8, "metric": 36.7}
			}
		},
		"distanceSkated": {
			"total": {"value": 5000, "rank": 12, "leagueAvg": {"value": 4800.0}}
		},
		"sogSummary": [
			{
				"locationCode": "all",
				"shots": 1800,
				"shotsRank": 5,
				"shotsLeagueAvg": 1600.0,
				"goals": 200,
				"goalsRank": 3,
				"goalsLeagueAvg": 170.0,
				"shootingPctg": 0.111,
				"shootingPctgRank": 8,
				"shootingPctgLeagueAvg": 0.106
			}
		],
		"sogDetails": [
			{"area": "High Slot", "shots": 450, "shotsRank": 4}
		],
		"zoneTimeDetails": {
			"offensiveZonePctg": 0.34,
			"offensiveZoneRank": 5,
			"offensiveZoneLeagueAvg": 0.31,
			"offensiveZoneEvPctg": 0.32,
			"offensiveZoneEvRank": 6,
			"neutralZonePctg": 0.34,
			"neutralZoneRank": 15,
			"neutralZoneLeagueAvg": 0.35,
			"defensiveZonePctg": 0.32,
			"defensiveZoneRank": 20,
			"defensiveZoneLeagueAvg": 0.34
		}
	}`

	var detail EdgeTeamDetail
	if err := json.Unmarshal([]byte(jsonData), &detail); err != nil {
		t.Fatalf("Failed to unmarshal EdgeTeamDetail: %v", err)
	}

	if detail.Team.ID != 22 {
		t.Errorf("Team.ID = %d, want 22", detail.Team.ID)
	}
	if detail.Team.Abbrev != "EDM" {
		t.Errorf("Team.Abbrev = %q, want %q", detail.Team.Abbrev, "EDM")
	}
	if detail.ShotSpeed.ShotAttemptsOver90.Rank != 3 {
		t.Errorf("ShotSpeed.ShotAttemptsOver90.Rank = %d, want 3", detail.ShotSpeed.ShotAttemptsOver90.Rank)
	}
	if detail.ShotSpeed.ShotAttemptsOver90.LeagueAvg == nil {
		t.Fatal("ShotSpeed.ShotAttemptsOver90.LeagueAvg is nil, want non-nil")
	}
	if detail.SkatingSpeed.BurstsOver22.LeagueAvg != nil {
		t.Errorf("SkatingSpeed.BurstsOver22.LeagueAvg = %v, want nil", detail.SkatingSpeed.BurstsOver22.LeagueAvg)
	}
	if detail.SkatingSpeed.SpeedMax.Rank != 1 {
		t.Errorf("SkatingSpeed.SpeedMax.Rank = %d, want 1", detail.SkatingSpeed.SpeedMax.Rank)
	}
	if detail.DistanceSkated.Total.Value != 5000 {
		t.Errorf("DistanceSkated.Total.Value = %d, want 5000", detail.DistanceSkated.Total.Value)
	}
	if detail.ZoneTimeDetails.OffensiveZoneRank != 5 {
		t.Errorf("ZoneTimeDetails.OffensiveZoneRank = %d, want 5", detail.ZoneTimeDetails.OffensiveZoneRank)
	}
}

func TestEdgeTeamZoneTimeDetails_Deserialization(t *testing.T) {
	// Structure matches real API response from /v1/edge/team-zone-time-details/{t}/{s}/{gt}
	jsonData := `{
		"zoneTimeDetails": [
			{"strengthCode": "all", "offensiveZonePctg": 0.43, "offensiveZoneRank": 3, "offensiveZoneLeagueAvg": 0.41, "neutralZonePctg": 0.17, "neutralZoneRank": 30, "neutralZoneLeagueAvg": 0.18, "defensiveZonePctg": 0.40, "defensiveZoneRank": 5, "defensiveZoneLeagueAvg": 0.41},
			{"strengthCode": "es", "offensiveZonePctg": 0.42, "offensiveZoneRank": 4, "offensiveZoneLeagueAvg": 0.41, "neutralZonePctg": 0.18, "neutralZoneRank": 30, "neutralZoneLeagueAvg": 0.19, "defensiveZonePctg": 0.40, "defensiveZoneRank": 6, "defensiveZoneLeagueAvg": 0.41},
			{"strengthCode": "pp", "offensiveZonePctg": 0.62, "offensiveZoneRank": 4, "offensiveZoneLeagueAvg": 0.59, "neutralZonePctg": 0.14, "neutralZoneRank": 24, "neutralZoneLeagueAvg": 0.14, "defensiveZonePctg": 0.25, "defensiveZoneRank": 4, "defensiveZoneLeagueAvg": 0.27},
			{"strengthCode": "pk", "offensiveZonePctg": 0.29, "offensiveZoneRank": 3, "offensiveZoneLeagueAvg": 0.27, "neutralZonePctg": 0.14, "neutralZoneRank": 13, "neutralZoneLeagueAvg": 0.14, "defensiveZonePctg": 0.57, "defensiveZoneRank": 6, "defensiveZoneLeagueAvg": 0.59}
		],
		"shotDifferential": {
			"shotAttemptDifferential": 5.01,
			"shotAttemptDifferentialRank": 3,
			"sogDifferential": 0.12,
			"sogDifferentialRank": 2
		}
	}`

	var detail EdgeTeamZoneTimeDetails
	if err := json.Unmarshal([]byte(jsonData), &detail); err != nil {
		t.Fatalf("Failed to unmarshal EdgeTeamZoneTimeDetails: %v", err)
	}

	if len(detail.ZoneTimeDetails) != 4 {
		t.Fatalf("ZoneTimeDetails length = %d, want 4", len(detail.ZoneTimeDetails))
	}
	if detail.ZoneTimeDetails[0].StrengthCode != "all" {
		t.Errorf("ZoneTimeDetails[0].StrengthCode = %q, want %q", detail.ZoneTimeDetails[0].StrengthCode, "all")
	}
	if detail.ZoneTimeDetails[2].OffensiveZonePctg != 0.62 {
		t.Errorf("ZoneTimeDetails[2].OffensiveZonePctg = %f, want 0.62", detail.ZoneTimeDetails[2].OffensiveZonePctg)
	}
	if detail.ShotDifferential == nil {
		t.Fatal("ShotDifferential is nil, want non-nil")
	}
	if detail.ShotDifferential.ShotAttemptDifferentialRank != 3 {
		t.Errorf("ShotDifferential.ShotAttemptDifferentialRank = %d, want 3", detail.ShotDifferential.ShotAttemptDifferentialRank)
	}
}

func TestEdgeSkaterSpeedDetail_Deserialization(t *testing.T) {
	jsonData := `{
		"player": {
			"id": 8478402,
			"firstName": {"default": "Connor"}, "lastName": {"default": "McDavid"},
			"birthDate": "1997-01-13", "shootsCatches": "L", "sweaterNumber": 97,
			"position": "C", "slug": "connor-mcdavid-8478402", "headshot": "h",
			"goals": 30, "assists": 60, "points": 90, "gamesPlayed": 50,
			"team": {"id": 22, "commonName": {"default": "Oilers"}, "placeNameWithPreposition": {"default": "Edmonton"}, "abbrev": "EDM", "teamLogo": {"light": "l", "dark": "d"}, "slug": "s", "conference": "W", "division": "P", "wins": 35, "losses": 15, "otLosses": 5, "gamesPlayed": 55, "points": 75}
		},
		"seasonsWithEdgeStats": [{"id": 20242025, "gameTypes": [2]}],
		"topSkatingSpeeds": [
			{"gameDate": "2025-01-15", "awayTeam": {"abbrev": "CGY", "score": 2}, "homeTeam": {"abbrev": "EDM", "score": 5}, "speed": {"imperial": 23.1, "metric": 37.2}},
			{"gameDate": "2025-01-10", "awayTeam": {"abbrev": "EDM", "score": 3}, "homeTeam": {"abbrev": "VAN", "score": 1}, "speed": {"imperial": 22.8, "metric": 36.7}}
		]
	}`

	var detail EdgeSkaterSpeedDetail
	if err := json.Unmarshal([]byte(jsonData), &detail); err != nil {
		t.Fatalf("Failed to unmarshal EdgeSkaterSpeedDetail: %v", err)
	}

	if len(detail.TopSkatingSpeeds) != 2 {
		t.Fatalf("TopSkatingSpeeds length = %d, want 2", len(detail.TopSkatingSpeeds))
	}
	if detail.TopSkatingSpeeds[0].Speed.Imperial != 23.1 {
		t.Errorf("TopSkatingSpeeds[0].Speed.Imperial = %f, want 23.1", detail.TopSkatingSpeeds[0].Speed.Imperial)
	}
}

func TestEdgeGoalie5v5Detail_Deserialization(t *testing.T) {
	jsonData := `{
		"player": {
			"id": 8479318,
			"firstName": {"default": "Igor"}, "lastName": {"default": "Shesterkin"},
			"birthDate": "1995-12-30", "shootsCatches": "L", "sweaterNumber": 31,
			"slug": "igor-shesterkin-8479318", "headshot": "h",
			"wins": 25, "losses": 10, "overtimeLosses": 3, "goalsAgainstAvg": 2.15,
			"savePctg": 0.928, "gamesPlayed": 38,
			"team": {"id": 3, "commonName": {"default": "Rangers"}, "placeNameWithPreposition": {"default": "New York"}, "abbrev": "NYR", "teamLogo": {"light": "l", "dark": "d"}, "slug": "s", "conference": "E", "division": "M", "wins": 30, "losses": 18, "otLosses": 7, "gamesPlayed": 55, "points": 67}
		},
		"seasonsWithEdgeStats": [{"id": 20242025, "gameTypes": [2]}],
		"savePctg5v5Last10": [
			{"gameDate": "2025-02-01", "awayTeam": {"abbrev": "NYR", "score": 3}, "homeTeam": {"abbrev": "BOS", "score": 1}, "savePctg": 0.950},
			{"gameDate": "2025-01-28", "awayTeam": {"abbrev": "PIT", "score": 2}, "homeTeam": {"abbrev": "NYR", "score": 4}, "savePctg": 0.935}
		]
	}`

	var detail EdgeGoalie5v5Detail
	if err := json.Unmarshal([]byte(jsonData), &detail); err != nil {
		t.Fatalf("Failed to unmarshal EdgeGoalie5v5Detail: %v", err)
	}

	if len(detail.SavePctg5v5Last10) != 2 {
		t.Fatalf("SavePctg5v5Last10 length = %d, want 2", len(detail.SavePctg5v5Last10))
	}
	if detail.SavePctg5v5Last10[0].SavePctg != 0.950 {
		t.Errorf("SavePctg5v5Last10[0].SavePctg = %f, want 0.950", detail.SavePctg5v5Last10[0].SavePctg)
	}
}

func TestEdgeOptionalOverlay_MarshalRoundTrip(t *testing.T) {
	// Test that optional overlay fields round-trip correctly
	original := EdgePercentileStatWithOverlay{
		Imperial:   100.0,
		Metric:     160.9,
		Percentile: 0.85,
		LeagueAvg:  EdgeMeasurement{Imperial: 90.0, Metric: 144.8},
		Overlay:    nil, // no overlay
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded EdgePercentileStatWithOverlay
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Imperial != original.Imperial {
		t.Errorf("Imperial = %f, want %f", decoded.Imperial, original.Imperial)
	}
	if decoded.Overlay != nil {
		t.Errorf("Overlay = %v, want nil", decoded.Overlay)
	}
}

func TestEdgeRankStat_OptionalLeagueAvg(t *testing.T) {
	// LeagueAvg is optional on EdgeRankStat
	jsonWithAvg := `{"value": 80, "rank": 2, "leagueAvg": {"value": 60.5}}`
	jsonWithoutAvg := `{"value": 80, "rank": 2}`

	var with EdgeRankStat
	if err := json.Unmarshal([]byte(jsonWithAvg), &with); err != nil {
		t.Fatalf("Failed to unmarshal with leagueAvg: %v", err)
	}
	if with.LeagueAvg == nil {
		t.Fatal("LeagueAvg is nil, want non-nil")
	}
	if with.LeagueAvg.Value != 60.5 {
		t.Errorf("LeagueAvg.Value = %f, want 60.5", with.LeagueAvg.Value)
	}

	var without EdgeRankStat
	if err := json.Unmarshal([]byte(jsonWithoutAvg), &without); err != nil {
		t.Fatalf("Failed to unmarshal without leagueAvg: %v", err)
	}
	if without.LeagueAvg != nil {
		t.Errorf("LeagueAvg = %v, want nil", without.LeagueAvg)
	}
}

// ===== Client Method Tests =====

func TestEdgeSkaterDetail_Client(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/edge/skater-detail/8478402/20242025/2"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q, want %q", r.URL.Path, expectedPath)
		}
		fmt.Fprint(w, `{
			"player": {"id": 8478402, "firstName": {"default": "Connor"}, "lastName": {"default": "McDavid"}, "birthDate": "1997-01-13", "shootsCatches": "L", "sweaterNumber": 97, "position": "C", "slug": "s", "headshot": "h", "goals": 30, "assists": 60, "points": 90, "gamesPlayed": 50, "team": {"id": 22, "commonName": {"default": "Oilers"}, "placeNameWithPreposition": {"default": "Edmonton"}, "abbrev": "EDM", "teamLogo": {"light": "l", "dark": "d"}, "slug": "s", "conference": "W", "division": "P", "wins": 35, "losses": 15, "otLosses": 5, "gamesPlayed": 55, "points": 75}},
			"seasonsWithEdgeStats": [],
			"topShotSpeed": {"imperial": 100.0, "metric": 160.9, "percentile": 0.85, "leagueAvg": {"imperial": 90.0, "metric": 144.8}},
			"skatingSpeed": {"speedMax": {"imperial": 22.0, "metric": 35.4, "percentile": 0.90, "leagueAvg": {"imperial": 21.0, "metric": 33.8}}, "burstsOver20": {"value": 100, "percentile": 0.80, "leagueAvg": {"value": 90.0}}},
			"totalDistanceSkated": {"imperial": 400.0, "metric": 643.7, "percentile": 0.80, "leagueAvg": {"imperial": 380.0, "metric": 611.5}},
			"distanceMaxGame": {"imperial": 10.0, "metric": 16.1, "percentile": 0.75, "leagueAvg": {"imperial": 9.5, "metric": 15.3}},
			"sogSummary": [],
			"sogDetails": [],
			"zoneTimeDetails": {"offensiveZonePctg": 0.30, "offensiveZonePercentile": 0.70, "offensiveZoneLeagueAvg": 0.29, "offensiveZoneEvPctg": 0.28, "offensiveZoneEvPercentile": 0.65, "offensiveZoneEvLeagueAvg": 0.27, "neutralZonePctg": 0.35, "neutralZonePercentile": 0.50, "neutralZoneLeagueAvg": 0.36, "defensiveZonePctg": 0.35, "defensiveZonePercentile": 0.60, "defensiveZoneLeagueAvg": 0.35}
		}`)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	result, err := client.EdgeSkaterDetail(context.Background(), PlayerID(8478402), NewSeason(2024), GameTypeRegularSeason)
	if err != nil {
		t.Fatalf("EdgeSkaterDetail returned error: %v", err)
	}
	if result.Player.ID != 8478402 {
		t.Errorf("Player.ID = %d, want 8478402", result.Player.ID)
	}
}

func TestEdgeTeamDetail_Client(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/edge/team-detail/22/20242025/2"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q, want %q", r.URL.Path, expectedPath)
		}
		fmt.Fprint(w, `{
			"team": {"id": 22, "commonName": {"default": "Oilers"}, "placeNameWithPreposition": {"default": "Edmonton"}, "abbrev": "EDM", "teamLogo": {"light": "l", "dark": "d"}, "slug": "s", "conference": "W", "division": "P", "wins": 35, "losses": 15, "otLosses": 5, "gamesPlayed": 55, "points": 75},
			"seasonsWithEdgeStats": [],
			"shotSpeed": {"shotAttemptsOver90": {"value": 100, "rank": 5}, "topShotSpeed": {"imperial": 100.0, "metric": 160.9, "rank": 3, "leagueAvg": {"imperial": 95.0, "metric": 152.9}}},
			"skatingSpeed": {"burstsOver22": {"value": 70, "rank": 4}, "burstsOver20": {"value": 400, "rank": 10}, "speedMax": {"imperial": 23.0, "metric": 37.0, "rank": 2, "leagueAvg": {"imperial": 22.0, "metric": 35.4}}},
			"distanceSkated": {"total": {"value": 4500, "rank": 15}},
			"sogSummary": [],
			"sogDetails": [],
			"zoneTimeDetails": {"offensiveZonePctg": 0.30, "offensiveZoneRank": 10, "offensiveZoneLeagueAvg": 0.31, "offensiveZoneEvPctg": 0.29, "offensiveZoneEvRank": 12, "neutralZonePctg": 0.35, "neutralZoneRank": 16, "neutralZoneLeagueAvg": 0.35, "defensiveZonePctg": 0.35, "defensiveZoneRank": 18, "defensiveZoneLeagueAvg": 0.34}
		}`)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	result, err := client.EdgeTeamDetail(context.Background(), TeamID(22), NewSeason(2024), GameTypeRegularSeason)
	if err != nil {
		t.Fatalf("EdgeTeamDetail returned error: %v", err)
	}
	if result.Team.Abbrev != "EDM" {
		t.Errorf("Team.Abbrev = %q, want %q", result.Team.Abbrev, "EDM")
	}
}

func TestEdgeLanding_Client(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/edge/skater-landing/20242025/2"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %q, want %q", r.URL.Path, expectedPath)
		}
		// Real API returns objects per category, not arrays
		fmt.Fprint(w, `{
			"seasonsWithEdgeStats": [{"id": 20242025, "gameTypes": [2]}],
			"leaders": {
				"hardestShot": {
					"player": {"id": 8478402, "firstName": {"default": "Connor"}, "lastName": {"default": "McDavid"}},
					"shotSpeed": {"imperial": 100.0, "metric": 160.9}
				}
			}
		}`)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	result, err := client.EdgeSkaterLanding(context.Background(), NewSeason(2024), GameTypeRegularSeason)
	if err != nil {
		t.Fatalf("EdgeSkaterLanding returned error: %v", err)
	}
	if result.Leaders == nil {
		t.Fatal("Leaders is nil, want non-nil")
	}
	if leader, ok := result.Leaders["hardestShot"]; !ok {
		t.Error("Leaders missing 'hardestShot' key")
	} else if leader.ShotSpeed == nil {
		t.Error("hardestShot.ShotSpeed is nil")
	}
}

func TestEdge_404Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	_, err := client.EdgeSkaterDetail(context.Background(), PlayerID(9999999), NewSeason(2020), GameTypeRegularSeason)
	if err == nil {
		t.Fatal("Expected error for 404 response, got nil")
	}

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %T: %v", err, err)
	}
}
