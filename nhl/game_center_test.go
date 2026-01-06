package nhl

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGameSituationFromCode(t *testing.T) {
	tests := []struct {
		name             string
		code             string
		wantNil          bool
		wantAwaySkaters  int
		wantAwayGoalieIn bool
		wantHomeSkaters  int
		wantHomeGoalieIn bool
	}{
		{
			name:             "even strength 5v5",
			code:             "1551",
			wantAwaySkaters:  5,
			wantAwayGoalieIn: true,
			wantHomeSkaters:  5,
			wantHomeGoalieIn: true,
		},
		{
			name:             "away power play 5v4",
			code:             "1541",
			wantAwaySkaters:  5,
			wantAwayGoalieIn: true,
			wantHomeSkaters:  4,
			wantHomeGoalieIn: true,
		},
		{
			name:             "home power play 4v5",
			code:             "1451",
			wantAwaySkaters:  4,
			wantAwayGoalieIn: true,
			wantHomeSkaters:  5,
			wantHomeGoalieIn: true,
		},
		{
			name:             "away empty net 6v5",
			code:             "0651",
			wantAwaySkaters:  6,
			wantAwayGoalieIn: false,
			wantHomeSkaters:  5,
			wantHomeGoalieIn: true,
		},
		{
			name:             "home empty net 5v6",
			code:             "1560",
			wantAwaySkaters:  5,
			wantAwayGoalieIn: true,
			wantHomeSkaters:  6,
			wantHomeGoalieIn: false,
		},
		{
			name:    "invalid code - too short",
			code:    "155",
			wantNil: true,
		},
		{
			name:    "invalid code - too long",
			code:    "15512",
			wantNil: true,
		},
		{
			name:    "invalid code - contains non-digits",
			code:    "abcd",
			wantNil: true,
		},
		{
			name:    "empty string",
			code:    "",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GameSituationFromCode(tt.code)

			if tt.wantNil {
				if got != nil {
					t.Errorf("GameSituationFromCode() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("GameSituationFromCode() = nil, want non-nil")
			}

			if got.AwaySkaters != tt.wantAwaySkaters {
				t.Errorf("AwaySkaters = %d, want %d", got.AwaySkaters, tt.wantAwaySkaters)
			}
			if got.AwayGoalieIn != tt.wantAwayGoalieIn {
				t.Errorf("AwayGoalieIn = %v, want %v", got.AwayGoalieIn, tt.wantAwayGoalieIn)
			}
			if got.HomeSkaters != tt.wantHomeSkaters {
				t.Errorf("HomeSkaters = %d, want %d", got.HomeSkaters, tt.wantHomeSkaters)
			}
			if got.HomeGoalieIn != tt.wantHomeGoalieIn {
				t.Errorf("HomeGoalieIn = %v, want %v", got.HomeGoalieIn, tt.wantHomeGoalieIn)
			}
		})
	}
}

func TestGameSituation_IsEvenStrength(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"5v5", "1551", true},
		{"4v4", "1441", true},
		{"3v3", "1331", true},
		{"5v4 power play", "1541", false},
		{"4v5 power play", "1451", false},
		{"6v5 empty net", "0651", false},
		{"5v6 empty net", "1560", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			situation := GameSituationFromCode(tt.code)
			if situation == nil {
				t.Fatalf("GameSituationFromCode() = nil")
			}

			if got := situation.IsEvenStrength(); got != tt.want {
				t.Errorf("IsEvenStrength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameSituation_IsAwayPowerPlay(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"5v5 even", "1551", false},
		{"5v4 away PP", "1541", true},
		{"4v5 home PP", "1451", false},
		{"6v5 away EN", "0651", true},
		{"5v3 away 2-man advantage", "1531", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			situation := GameSituationFromCode(tt.code)
			if situation == nil {
				t.Fatalf("GameSituationFromCode() = nil")
			}

			if got := situation.IsAwayPowerPlay(); got != tt.want {
				t.Errorf("IsAwayPowerPlay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameSituation_IsHomePowerPlay(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"5v5 even", "1551", false},
		{"5v4 away PP", "1541", false},
		{"4v5 home PP", "1451", true},
		{"5v6 home EN", "1560", true},
		{"3v5 home 2-man advantage", "1351", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			situation := GameSituationFromCode(tt.code)
			if situation == nil {
				t.Fatalf("GameSituationFromCode() = nil")
			}

			if got := situation.IsHomePowerPlay(); got != tt.want {
				t.Errorf("IsHomePowerPlay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameSituation_IsEmptyNet(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"5v5 both goalies in", "1551", false},
		{"6v5 away empty net", "0651", true},
		{"5v6 home empty net", "1560", true},
		{"6v6 both empty net", "0660", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			situation := GameSituationFromCode(tt.code)
			if situation == nil {
				t.Fatalf("GameSituationFromCode() = nil")
			}

			if got := situation.IsEmptyNet(); got != tt.want {
				t.Errorf("IsEmptyNet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameSituation_StrengthDescription(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{"5v5 even", "1551", "5v5"},
		{"4v4 even", "1441", "4v4"},
		{"3v3 even", "1331", "3v3"},
		{"5v4 PP", "1541", "5v4 PP"},
		{"4v5 PP", "1451", "4v5 PP"},
		{"6v5 EN", "0651", "6v5 EN"},
		{"5v6 EN", "1560", "5v6 EN"},
		{"6v4 PP and EN", "0641", "6v4 EN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			situation := GameSituationFromCode(tt.code)
			if situation == nil {
				t.Fatalf("GameSituationFromCode() = nil")
			}

			if got := situation.StrengthDescription(); got != tt.want {
				t.Errorf("StrengthDescription() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGameSituation_String(t *testing.T) {
	situation := GameSituationFromCode("1551")
	if situation == nil {
		t.Fatal("GameSituationFromCode() = nil")
	}

	want := "5v5"
	if got := situation.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestPlayEvent_Deserialization_Goal(t *testing.T) {
	jsonData := `{
		"eventId": 274,
		"periodDescriptor": {
			"number": 1,
			"periodType": "REG",
			"maxRegulationPeriods": 3
		},
		"timeInPeriod": "08:39",
		"timeRemaining": "11:21",
		"situationCode": "1551",
		"homeTeamDefendingSide": "right",
		"typeCode": 505,
		"typeDescKey": "goal",
		"sortOrder": 146,
		"details": {
			"xCoord": 71,
			"yCoord": -12,
			"zoneCode": "O",
			"shotType": "snap",
			"scoringPlayerId": 8476474,
			"scoringPlayerTotal": 1,
			"assist1PlayerId": 8480192,
			"assist1PlayerTotal": 1,
			"eventOwnerTeamId": 1,
			"goalieInNetId": 8480045,
			"awayScore": 1,
			"homeScore": 0,
			"highlightClip": 6362848229112,
			"discreteClip": 6362846260112
		}
	}`

	var event PlayEvent
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if event.EventID != 274 {
		t.Errorf("EventID = %d, want 274", event.EventID)
	}
	if event.TypeDescKey != PlayEventTypeGoal {
		t.Errorf("TypeDescKey = %q, want %q", event.TypeDescKey, PlayEventTypeGoal)
	}
	if event.TimeInPeriod != "08:39" {
		t.Errorf("TimeInPeriod = %q, want %q", event.TimeInPeriod, "08:39")
	}

	if event.Details == nil {
		t.Fatal("Details = nil, want non-nil")
	}

	details := event.Details
	if details.ScoringPlayerID == nil || *details.ScoringPlayerID != 8476474 {
		t.Errorf("ScoringPlayerID = %v, want 8476474", details.ScoringPlayerID)
	}
	if details.ScoringPlayerTotal == nil || *details.ScoringPlayerTotal != 1 {
		t.Errorf("ScoringPlayerTotal = %v, want 1", details.ScoringPlayerTotal)
	}
	if details.Assist1PlayerID == nil || *details.Assist1PlayerID != 8480192 {
		t.Errorf("Assist1PlayerID = %v, want 8480192", details.Assist1PlayerID)
	}
	if details.AwayScore == nil || *details.AwayScore != 1 {
		t.Errorf("AwayScore = %v, want 1", details.AwayScore)
	}
	if details.HomeScore == nil || *details.HomeScore != 0 {
		t.Errorf("HomeScore = %v, want 0", details.HomeScore)
	}
	if details.ShotType == nil || *details.ShotType != "snap" {
		t.Errorf("ShotType = %v, want %q", details.ShotType, "snap")
	}
}

func TestPlayEvent_Deserialization_Penalty(t *testing.T) {
	jsonData := `{
		"eventId": 135,
		"periodDescriptor": {
			"number": 1,
			"periodType": "REG",
			"maxRegulationPeriods": 3
		},
		"timeInPeriod": "01:37",
		"timeRemaining": "18:23",
		"situationCode": "1560",
		"homeTeamDefendingSide": "right",
		"typeCode": 509,
		"typeDescKey": "penalty",
		"sortOrder": 45,
		"details": {
			"xCoord": 1,
			"yCoord": -37,
			"zoneCode": "N",
			"typeCode": "MIN",
			"descKey": "slashing",
			"duration": 2,
			"committedByPlayerId": 8475287,
			"drawnByPlayerId": 8479420,
			"eventOwnerTeamId": 1
		}
	}`

	var event PlayEvent
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if event.EventID != 135 {
		t.Errorf("EventID = %d, want 135", event.EventID)
	}
	if event.TypeDescKey != PlayEventTypePenalty {
		t.Errorf("TypeDescKey = %q, want %q", event.TypeDescKey, PlayEventTypePenalty)
	}

	if event.Details == nil {
		t.Fatal("Details = nil, want non-nil")
	}

	details := event.Details
	if details.TypeCode == nil || *details.TypeCode != "MIN" {
		t.Errorf("TypeCode = %v, want %q", details.TypeCode, "MIN")
	}
	if details.DescKey == nil || *details.DescKey != "slashing" {
		t.Errorf("DescKey = %v, want %q", details.DescKey, "slashing")
	}
	if details.Duration == nil || *details.Duration != 2 {
		t.Errorf("Duration = %v, want 2", details.Duration)
	}
	if details.CommittedByPlayerID == nil || *details.CommittedByPlayerID != 8475287 {
		t.Errorf("CommittedByPlayerID = %v, want 8475287", details.CommittedByPlayerID)
	}
	if details.DrawnByPlayerID == nil || *details.DrawnByPlayerID != 8479420 {
		t.Errorf("DrawnByPlayerID = %v, want 8479420", details.DrawnByPlayerID)
	}
}

func TestPlayEvent_Deserialization_ShotOnGoal(t *testing.T) {
	jsonData := `{
		"eventId": 103,
		"periodDescriptor": {
			"number": 1,
			"periodType": "REG",
			"maxRegulationPeriods": 3
		},
		"timeInPeriod": "00:08",
		"timeRemaining": "19:52",
		"situationCode": "1551",
		"homeTeamDefendingSide": "right",
		"typeCode": 506,
		"typeDescKey": "shot-on-goal",
		"sortOrder": 13,
		"details": {
			"xCoord": 56,
			"yCoord": -39,
			"zoneCode": "O",
			"shotType": "wrist",
			"shootingPlayerId": 8483495,
			"goalieInNetId": 8480045,
			"eventOwnerTeamId": 1,
			"awaySOG": 1,
			"homeSOG": 0
		}
	}`

	var event PlayEvent
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if event.EventID != 103 {
		t.Errorf("EventID = %d, want 103", event.EventID)
	}
	if event.TypeDescKey != PlayEventTypeShotOnGoal {
		t.Errorf("TypeDescKey = %q, want %q", event.TypeDescKey, PlayEventTypeShotOnGoal)
	}

	if event.Details == nil {
		t.Fatal("Details = nil, want non-nil")
	}

	details := event.Details
	if details.ShotType == nil || *details.ShotType != "wrist" {
		t.Errorf("ShotType = %v, want %q", details.ShotType, "wrist")
	}
	if details.ShootingPlayerID == nil || *details.ShootingPlayerID != 8483495 {
		t.Errorf("ShootingPlayerID = %v, want 8483495", details.ShootingPlayerID)
	}
	if details.GoalieInNetID == nil || *details.GoalieInNetID != 8480045 {
		t.Errorf("GoalieInNetID = %v, want 8480045", details.GoalieInNetID)
	}
	if details.AwaySOG == nil || *details.AwaySOG != 1 {
		t.Errorf("AwaySOG = %v, want 1", details.AwaySOG)
	}
	if details.HomeSOG == nil || *details.HomeSOG != 0 {
		t.Errorf("HomeSOG = %v, want 0", details.HomeSOG)
	}
}

func TestPlayEvent_Deserialization_Faceoff(t *testing.T) {
	jsonData := `{
		"eventId": 151,
		"periodDescriptor": {
			"number": 1,
			"periodType": "REG",
			"maxRegulationPeriods": 3
		},
		"timeInPeriod": "00:00",
		"timeRemaining": "20:00",
		"situationCode": "1551",
		"homeTeamDefendingSide": "right",
		"typeCode": 502,
		"typeDescKey": "faceoff",
		"sortOrder": 11,
		"details": {
			"eventOwnerTeamId": 1,
			"losingPlayerId": 8478043,
			"winningPlayerId": 8480002,
			"xCoord": 0,
			"yCoord": 0,
			"zoneCode": "N"
		}
	}`

	var event PlayEvent
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if event.EventID != 151 {
		t.Errorf("EventID = %d, want 151", event.EventID)
	}
	if event.TypeDescKey != PlayEventTypeFaceoff {
		t.Errorf("TypeDescKey = %q, want %q", event.TypeDescKey, PlayEventTypeFaceoff)
	}

	if event.Details == nil {
		t.Fatal("Details = nil, want non-nil")
	}

	details := event.Details
	if details.WinningPlayerID == nil || *details.WinningPlayerID != 8480002 {
		t.Errorf("WinningPlayerID = %v, want 8480002", details.WinningPlayerID)
	}
	if details.LosingPlayerID == nil || *details.LosingPlayerID != 8478043 {
		t.Errorf("LosingPlayerID = %v, want 8478043", details.LosingPlayerID)
	}
	if details.ZoneCode == nil || *details.ZoneCode != ZoneCodeNeutral {
		t.Errorf("ZoneCode = %v, want %q", details.ZoneCode, ZoneCodeNeutral)
	}
}

func TestPlayEvent_Deserialization_BlockedShot(t *testing.T) {
	jsonData := `{
		"eventId": 63,
		"periodDescriptor": {
			"number": 1,
			"periodType": "REG",
			"maxRegulationPeriods": 3
		},
		"timeInPeriod": "01:15",
		"timeRemaining": "18:45",
		"situationCode": "1551",
		"homeTeamDefendingSide": "right",
		"typeCode": 508,
		"typeDescKey": "blocked-shot",
		"sortOrder": 24,
		"details": {
			"xCoord": -73,
			"yCoord": 28,
			"zoneCode": "D",
			"blockingPlayerId": 8481568,
			"shootingPlayerId": 8479323,
			"eventOwnerTeamId": 3,
			"reason": "blocked"
		}
	}`

	var event PlayEvent
	if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if event.TypeDescKey != PlayEventTypeBlockedShot {
		t.Errorf("TypeDescKey = %q, want %q", event.TypeDescKey, PlayEventTypeBlockedShot)
	}

	if event.Details == nil {
		t.Fatal("Details = nil, want non-nil")
	}

	details := event.Details
	if details.BlockingPlayerID == nil || *details.BlockingPlayerID != 8481568 {
		t.Errorf("BlockingPlayerID = %v, want 8481568", details.BlockingPlayerID)
	}
	if details.ShootingPlayerID == nil || *details.ShootingPlayerID != 8479323 {
		t.Errorf("ShootingPlayerID = %v, want 8479323", details.ShootingPlayerID)
	}
}

func TestPlayEvent_Situation(t *testing.T) {
	event := PlayEvent{
		SituationCode: "1551",
	}

	situation := event.Situation()
	if situation == nil {
		t.Fatal("Situation() = nil, want non-nil")
	}

	if !situation.IsEvenStrength() {
		t.Error("expected even strength situation")
	}
}

func TestRosterSpot_Deserialization(t *testing.T) {
	jsonData := `{
		"teamId": 1,
		"playerId": 8474593,
		"firstName": {"default": "Jacob"},
		"lastName": {"default": "Markstrom"},
		"sweaterNumber": 25,
		"positionCode": "G",
		"headshot": "https://assets.nhle.com/mugs/nhl/20242025/NJD/8474593.png"
	}`

	var roster RosterSpot
	if err := json.Unmarshal([]byte(jsonData), &roster); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if roster.TeamID != 1 {
		t.Errorf("TeamID = %d, want 1", roster.TeamID)
	}
	if roster.PlayerID != 8474593 {
		t.Errorf("PlayerID = %d, want 8474593", roster.PlayerID)
	}
	if roster.FirstName.Default != "Jacob" {
		t.Errorf("FirstName.Default = %q, want %q", roster.FirstName.Default, "Jacob")
	}
	if roster.LastName.Default != "Markstrom" {
		t.Errorf("LastName.Default = %q, want %q", roster.LastName.Default, "Markstrom")
	}
	if roster.SweaterNumber != 25 {
		t.Errorf("SweaterNumber = %d, want 25", roster.SweaterNumber)
	}
	if roster.Position != PositionGoalie {
		t.Errorf("Position = %q, want %q", roster.Position, PositionGoalie)
	}
}

func TestGameOutcome_Deserialization(t *testing.T) {
	jsonData := `{"lastPeriodType": "REG"}`

	var outcome GameOutcome
	if err := json.Unmarshal([]byte(jsonData), &outcome); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if outcome.LastPeriodType != PeriodTypeRegulation {
		t.Errorf("LastPeriodType = %q, want %q", outcome.LastPeriodType, PeriodTypeRegulation)
	}
}

func TestShiftEntry_Deserialization(t *testing.T) {
	jsonData := `{
		"id": 14376602,
		"detailCode": 0,
		"duration": "17:15",
		"endTime": "17:15",
		"eventDescription": null,
		"eventNumber": 101,
		"firstName": "Jacob",
		"gameId": 2024020001,
		"hexValue": "#C8102E",
		"lastName": "Markstrom",
		"period": 1,
		"playerId": 8474593,
		"shiftNumber": 1,
		"startTime": "00:00",
		"teamAbbrev": "NJD",
		"teamId": 1,
		"teamName": "New Jersey Devils",
		"typeCode": 517
	}`

	var shift ShiftEntry
	if err := json.Unmarshal([]byte(jsonData), &shift); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if shift.ID != 14376602 {
		t.Errorf("ID = %d, want 14376602", shift.ID)
	}
	if shift.DetailCode != 0 {
		t.Errorf("DetailCode = %d, want 0", shift.DetailCode)
	}
	if shift.Duration != "17:15" {
		t.Errorf("Duration = %q, want %q", shift.Duration, "17:15")
	}
	if shift.EndTime != "17:15" {
		t.Errorf("EndTime = %q, want %q", shift.EndTime, "17:15")
	}
	if shift.EventDescription != nil {
		t.Errorf("EventDescription = %v, want nil", shift.EventDescription)
	}
	if shift.EventNumber != 101 {
		t.Errorf("EventNumber = %d, want 101", shift.EventNumber)
	}
	if shift.FirstName != "Jacob" {
		t.Errorf("FirstName = %q, want %q", shift.FirstName, "Jacob")
	}
	if shift.GameID != 2024020001 {
		t.Errorf("GameID = %d, want 2024020001", shift.GameID)
	}
	if shift.HexValue != "#C8102E" {
		t.Errorf("HexValue = %q, want %q", shift.HexValue, "#C8102E")
	}
	if shift.LastName != "Markstrom" {
		t.Errorf("LastName = %q, want %q", shift.LastName, "Markstrom")
	}
	if shift.Period != 1 {
		t.Errorf("Period = %d, want 1", shift.Period)
	}
	if shift.PlayerID != 8474593 {
		t.Errorf("PlayerID = %d, want 8474593", shift.PlayerID)
	}
	if shift.ShiftNumber != 1 {
		t.Errorf("ShiftNumber = %d, want 1", shift.ShiftNumber)
	}
	if shift.StartTime != "00:00" {
		t.Errorf("StartTime = %q, want %q", shift.StartTime, "00:00")
	}
	if shift.TeamAbbrev != "NJD" {
		t.Errorf("TeamAbbrev = %q, want %q", shift.TeamAbbrev, "NJD")
	}
	if shift.TeamID != 1 {
		t.Errorf("TeamID = %d, want 1", shift.TeamID)
	}
	if shift.TeamName != "New Jersey Devils" {
		t.Errorf("TeamName = %q, want %q", shift.TeamName, "New Jersey Devils")
	}
	if shift.TypeCode != 517 {
		t.Errorf("TypeCode = %d, want 517", shift.TypeCode)
	}
}

func TestShiftChart_Deserialization(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"id": 14376602,
				"detailCode": 0,
				"duration": "17:15",
				"endTime": "17:15",
				"eventDescription": null,
				"eventNumber": 101,
				"firstName": "Jacob",
				"gameId": 2024020001,
				"hexValue": "#C8102E",
				"lastName": "Markstrom",
				"period": 1,
				"playerId": 8474593,
				"shiftNumber": 1,
				"startTime": "00:00",
				"teamAbbrev": "NJD",
				"teamId": 1,
				"teamName": "New Jersey Devils",
				"typeCode": 517
			}
		]
	}`

	var chart ShiftChart
	if err := json.Unmarshal([]byte(jsonData), &chart); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(chart.Data) != 1 {
		t.Fatalf("len(Data) = %d, want 1", len(chart.Data))
	}

	shift := chart.Data[0]
	if shift.PlayerID != 8474593 {
		t.Errorf("PlayerID = %d, want 8474593", shift.PlayerID)
	}
	if shift.FirstName != "Jacob" {
		t.Errorf("FirstName = %q, want %q", shift.FirstName, "Jacob")
	}
	if shift.LastName != "Markstrom" {
		t.Errorf("LastName = %q, want %q", shift.LastName, "Markstrom")
	}
}

func TestPlayByPlay_RecentPlays(t *testing.T) {
	pbp := &PlayByPlay{
		Plays: []PlayEvent{
			{EventID: 1},
			{EventID: 2},
			{EventID: 3},
			{EventID: 4},
			{EventID: 5},
		},
	}

	tests := []struct {
		name      string
		count     int
		wantCount int
		wantFirst int64
	}{
		{"get 3 recent", 3, 3, 5},
		{"get all", 5, 5, 5},
		{"get more than available", 10, 5, 5},
		{"get zero", 0, 0, 0},
		{"get negative", -1, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recent := pbp.RecentPlays(tt.count)

			if len(recent) != tt.wantCount {
				t.Errorf("len(RecentPlays(%d)) = %d, want %d", tt.count, len(recent), tt.wantCount)
			}

			if tt.wantCount > 0 && recent[0].EventID != tt.wantFirst {
				t.Errorf("first EventID = %d, want %d", recent[0].EventID, tt.wantFirst)
			}
		})
	}
}

func TestPlayByPlay_Goals(t *testing.T) {
	pbp := &PlayByPlay{
		Plays: []PlayEvent{
			{EventID: 1, TypeDescKey: PlayEventTypeFaceoff},
			{EventID: 2, TypeDescKey: PlayEventTypeGoal},
			{EventID: 3, TypeDescKey: PlayEventTypeShotOnGoal},
			{EventID: 4, TypeDescKey: PlayEventTypeGoal},
			{EventID: 5, TypeDescKey: PlayEventTypePenalty},
		},
	}

	goals := pbp.Goals()
	if len(goals) != 2 {
		t.Fatalf("len(Goals()) = %d, want 2", len(goals))
	}

	if goals[0].EventID != 2 {
		t.Errorf("goals[0].EventID = %d, want 2", goals[0].EventID)
	}
	if goals[1].EventID != 4 {
		t.Errorf("goals[1].EventID = %d, want 4", goals[1].EventID)
	}
}

func TestPlayByPlay_Penalties(t *testing.T) {
	pbp := &PlayByPlay{
		Plays: []PlayEvent{
			{EventID: 1, TypeDescKey: PlayEventTypeFaceoff},
			{EventID: 2, TypeDescKey: PlayEventTypePenalty},
			{EventID: 3, TypeDescKey: PlayEventTypeGoal},
			{EventID: 4, TypeDescKey: PlayEventTypePenalty},
		},
	}

	penalties := pbp.Penalties()
	if len(penalties) != 2 {
		t.Fatalf("len(Penalties()) = %d, want 2", len(penalties))
	}

	if penalties[0].EventID != 2 {
		t.Errorf("penalties[0].EventID = %d, want 2", penalties[0].EventID)
	}
	if penalties[1].EventID != 4 {
		t.Errorf("penalties[1].EventID = %d, want 4", penalties[1].EventID)
	}
}

func TestPlayByPlay_Shots(t *testing.T) {
	pbp := &PlayByPlay{
		Plays: []PlayEvent{
			{EventID: 1, TypeDescKey: PlayEventTypeFaceoff},
			{EventID: 2, TypeDescKey: PlayEventTypeShotOnGoal},
			{EventID: 3, TypeDescKey: PlayEventTypeMissedShot},
			{EventID: 4, TypeDescKey: PlayEventTypeBlockedShot},
			{EventID: 5, TypeDescKey: PlayEventTypeGoal},
			{EventID: 6, TypeDescKey: PlayEventTypePenalty},
		},
	}

	shots := pbp.Shots()
	if len(shots) != 4 {
		t.Fatalf("len(Shots()) = %d, want 4", len(shots))
	}

	expectedIDs := []int64{2, 3, 4, 5}
	for i, shot := range shots {
		if shot.EventID != expectedIDs[i] {
			t.Errorf("shots[%d].EventID = %d, want %d", i, shot.EventID, expectedIDs[i])
		}
	}
}

func TestPlayByPlay_PlaysInPeriod(t *testing.T) {
	pbp := &PlayByPlay{
		Plays: []PlayEvent{
			{EventID: 1, PeriodDescriptor: PeriodDescriptor{Number: 1}},
			{EventID: 2, PeriodDescriptor: PeriodDescriptor{Number: 1}},
			{EventID: 3, PeriodDescriptor: PeriodDescriptor{Number: 2}},
			{EventID: 4, PeriodDescriptor: PeriodDescriptor{Number: 2}},
			{EventID: 5, PeriodDescriptor: PeriodDescriptor{Number: 3}},
		},
	}

	tests := []struct {
		period    int
		wantCount int
	}{
		{1, 2},
		{2, 2},
		{3, 1},
		{4, 0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("period %d", tt.period), func(t *testing.T) {
			plays := pbp.PlaysInPeriod(tt.period)
			if len(plays) != tt.wantCount {
				t.Errorf("len(PlaysInPeriod(%d)) = %d, want %d", tt.period, len(plays), tt.wantCount)
			}
		})
	}
}

func TestPlayByPlay_GetPlayer(t *testing.T) {
	pbp := &PlayByPlay{
		RosterSpots: []RosterSpot{
			{PlayerID: 8474593, FirstName: LocalizedString{Default: "Jacob"}},
			{PlayerID: 8476474, FirstName: LocalizedString{Default: "John"}},
		},
	}

	tests := []struct {
		name     string
		playerID int64
		wantNil  bool
		wantName string
	}{
		{"existing player", 8474593, false, "Jacob"},
		{"another existing player", 8476474, false, "John"},
		{"non-existing player", 9999999, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := pbp.GetPlayer(tt.playerID)

			if tt.wantNil {
				if player != nil {
					t.Errorf("GetPlayer(%d) = %v, want nil", tt.playerID, player)
				}
				return
			}

			if player == nil {
				t.Fatalf("GetPlayer(%d) = nil, want non-nil", tt.playerID)
			}

			if player.FirstName.Default != tt.wantName {
				t.Errorf("FirstName.Default = %q, want %q", player.FirstName.Default, tt.wantName)
			}
		})
	}
}

func TestPlayByPlay_TeamRoster(t *testing.T) {
	pbp := &PlayByPlay{
		RosterSpots: []RosterSpot{
			{PlayerID: 1, TeamID: 1},
			{PlayerID: 2, TeamID: 1},
			{PlayerID: 3, TeamID: 2},
			{PlayerID: 4, TeamID: 2},
			{PlayerID: 5, TeamID: 2},
		},
	}

	tests := []struct {
		teamID    int64
		wantCount int
	}{
		{1, 2},
		{2, 3},
		{3, 0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("team %d", tt.teamID), func(t *testing.T) {
			roster := pbp.TeamRoster(tt.teamID)
			if len(roster) != tt.wantCount {
				t.Errorf("len(TeamRoster(%d)) = %d, want %d", tt.teamID, len(roster), tt.wantCount)
			}
		})
	}
}

func TestPlayByPlay_CurrentSituation(t *testing.T) {
	tests := []struct {
		name    string
		plays   []PlayEvent
		wantNil bool
		want    string
	}{
		{
			name:    "no plays",
			plays:   []PlayEvent{},
			wantNil: true,
		},
		{
			name: "valid situation",
			plays: []PlayEvent{
				{SituationCode: "1551"},
			},
			wantNil: false,
			want:    "5v5",
		},
		{
			name: "power play situation",
			plays: []PlayEvent{
				{SituationCode: "1541"},
			},
			wantNil: false,
			want:    "5v4 PP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pbp := &PlayByPlay{Plays: tt.plays}
			situation := pbp.CurrentSituation()

			if tt.wantNil {
				if situation != nil {
					t.Errorf("CurrentSituation() = %v, want nil", situation)
				}
				return
			}

			if situation == nil {
				t.Fatal("CurrentSituation() = nil, want non-nil")
			}

			if got := situation.String(); got != tt.want {
				t.Errorf("CurrentSituation().String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Additional error path tests for GameSituationFromCode

func TestGameSituationFromCode_InvalidFormat(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{"too short", "123"},
		{"too long", "12345"},
		{"invalid away skaters", "1X51"},
		{"invalid home skaters", "15X1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GameSituationFromCode(tt.code)
			if result != nil {
				t.Error("GameSituationFromCode() should return nil for invalid format")
			}
		})
	}
}
