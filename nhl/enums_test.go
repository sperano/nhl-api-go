package nhl

import (
	"encoding/json"
	"testing"
)

// Position tests

func TestPosition_Code(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		want     string
	}{
		{"center", PositionCenter, "C"},
		{"left wing", PositionLeftWing, "LW"},
		{"right wing", PositionRightWing, "RW"},
		{"defense", PositionDefense, "D"},
		{"goalie", PositionGoalie, "G"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.Code(); got != tt.want {
				t.Errorf("Position.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_Name(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		want     string
	}{
		{"center", PositionCenter, "Center"},
		{"left wing", PositionLeftWing, "Left Wing"},
		{"right wing", PositionRightWing, "Right Wing"},
		{"defense", PositionDefense, "Defense"},
		{"goalie", PositionGoalie, "Goalie"},
		{"unknown", Position("X"), "Unknown(X)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.Name(); got != tt.want {
				t.Errorf("Position.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_IsForward(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		want     bool
	}{
		{"center is forward", PositionCenter, true},
		{"left wing is forward", PositionLeftWing, true},
		{"right wing is forward", PositionRightWing, true},
		{"defense not forward", PositionDefense, false},
		{"goalie not forward", PositionGoalie, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.IsForward(); got != tt.want {
				t.Errorf("Position.IsForward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_IsSkater(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		want     bool
	}{
		{"center is skater", PositionCenter, true},
		{"left wing is skater", PositionLeftWing, true},
		{"right wing is skater", PositionRightWing, true},
		{"defense is skater", PositionDefense, true},
		{"goalie not skater", PositionGoalie, false},
		{"invalid not skater", Position("X"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.IsSkater(); got != tt.want {
				t.Errorf("Position.IsSkater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		want     bool
	}{
		{"center valid", PositionCenter, true},
		{"left wing valid", PositionLeftWing, true},
		{"right wing valid", PositionRightWing, true},
		{"defense valid", PositionDefense, true},
		{"goalie valid", PositionGoalie, true},
		{"unknown invalid", Position("X"), false},
		{"empty invalid", Position(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.IsValid(); got != tt.want {
				t.Errorf("Position.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPositionFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Position
		wantErr bool
	}{
		{"code center", "C", PositionCenter, false},
		{"name center", "Center", PositionCenter, false},
		{"code left wing short", "L", PositionLeftWing, false},
		{"code left wing", "LW", PositionLeftWing, false},
		{"name left wing", "Left Wing", PositionLeftWing, false},
		{"name left wing no space", "LeftWing", PositionLeftWing, false},
		{"code right wing short", "R", PositionRightWing, false},
		{"code right wing", "RW", PositionRightWing, false},
		{"name right wing", "Right Wing", PositionRightWing, false},
		{"name right wing no space", "RightWing", PositionRightWing, false},
		{"code defense", "D", PositionDefense, false},
		{"name defense", "Defense", PositionDefense, false},
		{"name defenseman", "Defenseman", PositionDefense, false},
		{"code goalie", "G", PositionGoalie, false},
		{"name goalie", "Goalie", PositionGoalie, false},
		{"name goaltender", "Goaltender", PositionGoalie, false},
		{"invalid", "X", Position(""), true},
		{"empty", "", Position(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PositionFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PositionFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("PositionFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustPositionFromString(t *testing.T) {
	t.Run("valid position", func(t *testing.T) {
		got := MustPositionFromString("C")
		if got != PositionCenter {
			t.Errorf("MustPositionFromString() = %v, want %v", got, PositionCenter)
		}
	})

	t.Run("invalid position panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustPositionFromString() did not panic")
			}
		}()
		MustPositionFromString("INVALID")
	})
}

func TestPosition_JSON(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		json     string
		wantErr  bool
	}{
		{"center", PositionCenter, `"C"`, false},
		{"left wing", PositionLeftWing, `"LW"`, false},
		{"right wing", PositionRightWing, `"RW"`, false},
		{"defense", PositionDefense, `"D"`, false},
		{"goalie", PositionGoalie, `"G"`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name+" marshal", func(t *testing.T) {
			got, err := json.Marshal(tt.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.json {
				t.Errorf("json.Marshal() = %v, want %v", string(got), tt.json)
			}
		})

		t.Run(tt.name+" unmarshal", func(t *testing.T) {
			var got Position
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.position {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.position)
			}
		})
	}
}

// Handedness tests

func TestHandedness_Code(t *testing.T) {
	tests := []struct {
		name       string
		handedness Handedness
		want       string
	}{
		{"left", HandednessLeft, "L"},
		{"right", HandednessRight, "R"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.handedness.Code(); got != tt.want {
				t.Errorf("Handedness.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandedness_Name(t *testing.T) {
	tests := []struct {
		name       string
		handedness Handedness
		want       string
	}{
		{"left", HandednessLeft, "Left"},
		{"right", HandednessRight, "Right"},
		{"unknown", Handedness("X"), "Unknown(X)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.handedness.Name(); got != tt.want {
				t.Errorf("Handedness.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandedness_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		handedness Handedness
		want       bool
	}{
		{"left valid", HandednessLeft, true},
		{"right valid", HandednessRight, true},
		{"invalid", Handedness("X"), false},
		{"empty invalid", Handedness(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.handedness.IsValid(); got != tt.want {
				t.Errorf("Handedness.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandednessFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Handedness
		wantErr bool
	}{
		{"code left", "L", HandednessLeft, false},
		{"name left", "Left", HandednessLeft, false},
		{"code right", "R", HandednessRight, false},
		{"name right", "Right", HandednessRight, false},
		{"invalid", "X", Handedness(""), true},
		{"empty", "", Handedness(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandednessFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandednessFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("HandednessFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandedness_JSON(t *testing.T) {
	tests := []struct {
		name       string
		handedness Handedness
		json       string
	}{
		{"left", HandednessLeft, `"L"`},
		{"right", HandednessRight, `"R"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.handedness)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(data) != tt.json {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.json)
			}

			var got Handedness
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got != tt.handedness {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.handedness)
			}
		})
	}
}

// GoalieDecision tests

func TestGoalieDecision_String(t *testing.T) {
	tests := []struct {
		name     string
		decision GoalieDecision
		want     string
	}{
		{"win", GoalieDecisionWin, "Win"},
		{"loss", GoalieDecisionLoss, "Loss"},
		{"tie", GoalieDecisionTie, "Tie"},
		{"overtime loss", GoalieDecisionOvertimeLoss, "Overtime Loss"},
		{"unknown", GoalieDecision("X"), "Unknown(X)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.decision.String(); got != tt.want {
				t.Errorf("GoalieDecision.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoalieDecision_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		decision GoalieDecision
		want     bool
	}{
		{"win valid", GoalieDecisionWin, true},
		{"loss valid", GoalieDecisionLoss, true},
		{"tie valid", GoalieDecisionTie, true},
		{"overtime loss valid", GoalieDecisionOvertimeLoss, true},
		{"invalid", GoalieDecision("X"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.decision.IsValid(); got != tt.want {
				t.Errorf("GoalieDecision.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoalieDecisionFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    GoalieDecision
		wantErr bool
	}{
		{"code win", "W", GoalieDecisionWin, false},
		{"name win", "Win", GoalieDecisionWin, false},
		{"code loss", "L", GoalieDecisionLoss, false},
		{"name loss", "Loss", GoalieDecisionLoss, false},
		{"code tie", "T", GoalieDecisionTie, false},
		{"name tie", "Tie", GoalieDecisionTie, false},
		{"code otl short", "O", GoalieDecisionOvertimeLoss, false},
		{"code otl", "OTL", GoalieDecisionOvertimeLoss, false},
		{"name overtime loss", "Overtime Loss", GoalieDecisionOvertimeLoss, false},
		{"name overtime loss no space", "OvertimeLoss", GoalieDecisionOvertimeLoss, false},
		{"invalid", "X", GoalieDecision(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoalieDecisionFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoalieDecisionFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GoalieDecisionFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoalieDecision_JSON(t *testing.T) {
	tests := []struct {
		name     string
		decision GoalieDecision
		json     string
	}{
		{"win", GoalieDecisionWin, `"W"`},
		{"loss", GoalieDecisionLoss, `"L"`},
		{"tie", GoalieDecisionTie, `"T"`},
		{"overtime loss", GoalieDecisionOvertimeLoss, `"OTL"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.decision)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(data) != tt.json {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.json)
			}

			var got GoalieDecision
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got != tt.decision {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.decision)
			}
		})
	}
}

// PeriodType tests

func TestPeriodType_Code(t *testing.T) {
	tests := []struct {
		name       string
		periodType PeriodType
		want       string
	}{
		{"regulation", PeriodTypeRegulation, "REG"},
		{"overtime", PeriodTypeOvertime, "OT"},
		{"shootout", PeriodTypeShootout, "SO"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.periodType.Code(); got != tt.want {
				t.Errorf("PeriodType.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodType_Name(t *testing.T) {
	tests := []struct {
		name       string
		periodType PeriodType
		want       string
	}{
		{"regulation", PeriodTypeRegulation, "Regulation"},
		{"overtime", PeriodTypeOvertime, "Overtime"},
		{"shootout", PeriodTypeShootout, "Shootout"},
		{"unknown", PeriodType("X"), "Unknown(X)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.periodType.Name(); got != tt.want {
				t.Errorf("PeriodType.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodType_IsOvertime(t *testing.T) {
	tests := []struct {
		name       string
		periodType PeriodType
		want       bool
	}{
		{"regulation not overtime", PeriodTypeRegulation, false},
		{"overtime is overtime", PeriodTypeOvertime, true},
		{"shootout is overtime", PeriodTypeShootout, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.periodType.IsOvertime(); got != tt.want {
				t.Errorf("PeriodType.IsOvertime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodTypeFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    PeriodType
		wantErr bool
	}{
		{"code regulation", "REG", PeriodTypeRegulation, false},
		{"name regulation", "Regulation", PeriodTypeRegulation, false},
		{"code overtime", "OT", PeriodTypeOvertime, false},
		{"name overtime", "Overtime", PeriodTypeOvertime, false},
		{"code shootout", "SO", PeriodTypeShootout, false},
		{"name shootout", "Shootout", PeriodTypeShootout, false},
		{"invalid", "X", PeriodType(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PeriodTypeFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeriodTypeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("PeriodTypeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodType_JSON(t *testing.T) {
	tests := []struct {
		name       string
		periodType PeriodType
		json       string
	}{
		{"regulation", PeriodTypeRegulation, `"REG"`},
		{"overtime", PeriodTypeOvertime, `"OT"`},
		{"shootout", PeriodTypeShootout, `"SO"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.periodType)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(data) != tt.json {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.json)
			}

			var got PeriodType
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got != tt.periodType {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.periodType)
			}
		})
	}
}

// HomeRoad tests

func TestHomeRoad_Code(t *testing.T) {
	tests := []struct {
		name     string
		homeRoad HomeRoad
		want     string
	}{
		{"home", HomeRoadHome, "H"},
		{"road", HomeRoadRoad, "R"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.homeRoad.Code(); got != tt.want {
				t.Errorf("HomeRoad.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHomeRoad_Name(t *testing.T) {
	tests := []struct {
		name     string
		homeRoad HomeRoad
		want     string
	}{
		{"home", HomeRoadHome, "Home"},
		{"road", HomeRoadRoad, "Road"},
		{"unknown", HomeRoad("X"), "Unknown(X)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.homeRoad.Name(); got != tt.want {
				t.Errorf("HomeRoad.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHomeRoadFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    HomeRoad
		wantErr bool
	}{
		{"code home", "H", HomeRoadHome, false},
		{"name home", "Home", HomeRoadHome, false},
		{"code road", "R", HomeRoadRoad, false},
		{"name road", "Road", HomeRoadRoad, false},
		{"name away", "Away", HomeRoadRoad, false},
		{"invalid", "X", HomeRoad(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HomeRoadFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HomeRoadFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("HomeRoadFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHomeRoad_JSON(t *testing.T) {
	tests := []struct {
		name     string
		homeRoad HomeRoad
		json     string
	}{
		{"home", HomeRoadHome, `"H"`},
		{"road", HomeRoadRoad, `"R"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.homeRoad)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(data) != tt.json {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.json)
			}

			var got HomeRoad
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got != tt.homeRoad {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.homeRoad)
			}
		})
	}
}

// ZoneCode tests

func TestZoneCode_Code(t *testing.T) {
	tests := []struct {
		name     string
		zoneCode ZoneCode
		want     string
	}{
		{"offensive", ZoneCodeOffensive, "O"},
		{"defensive", ZoneCodeDefensive, "D"},
		{"neutral", ZoneCodeNeutral, "N"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zoneCode.Code(); got != tt.want {
				t.Errorf("ZoneCode.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZoneCode_String(t *testing.T) {
	tests := []struct {
		name     string
		zoneCode ZoneCode
		want     string
	}{
		{"offensive", ZoneCodeOffensive, "Offensive"},
		{"defensive", ZoneCodeDefensive, "Defensive"},
		{"neutral", ZoneCodeNeutral, "Neutral"},
		{"unknown", ZoneCode("X"), "Unknown(X)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zoneCode.String(); got != tt.want {
				t.Errorf("ZoneCode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZoneCodeFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ZoneCode
		wantErr bool
	}{
		{"code offensive", "O", ZoneCodeOffensive, false},
		{"name offensive", "Offensive", ZoneCodeOffensive, false},
		{"code defensive", "D", ZoneCodeDefensive, false},
		{"name defensive", "Defensive", ZoneCodeDefensive, false},
		{"code neutral", "N", ZoneCodeNeutral, false},
		{"name neutral", "Neutral", ZoneCodeNeutral, false},
		{"invalid", "X", ZoneCode(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ZoneCodeFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZoneCodeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ZoneCodeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZoneCode_JSON(t *testing.T) {
	tests := []struct {
		name     string
		zoneCode ZoneCode
		json     string
	}{
		{"offensive", ZoneCodeOffensive, `"O"`},
		{"defensive", ZoneCodeDefensive, `"D"`},
		{"neutral", ZoneCodeNeutral, `"N"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.zoneCode)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(data) != tt.json {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.json)
			}

			var got ZoneCode
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got != tt.zoneCode {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.zoneCode)
			}
		})
	}
}

// DefendingSide tests

func TestDefendingSide_String(t *testing.T) {
	tests := []struct {
		name string
		side DefendingSide
		want string
	}{
		{"left", DefendingSideLeft, "left"},
		{"right", DefendingSideRight, "right"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.side.String(); got != tt.want {
				t.Errorf("DefendingSide.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefendingSide_IsValid(t *testing.T) {
	tests := []struct {
		name string
		side DefendingSide
		want bool
	}{
		{"left valid", DefendingSideLeft, true},
		{"right valid", DefendingSideRight, true},
		{"invalid", DefendingSide("center"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.side.IsValid(); got != tt.want {
				t.Errorf("DefendingSide.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefendingSideFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    DefendingSide
		wantErr bool
	}{
		{"left", "left", DefendingSideLeft, false},
		{"right", "right", DefendingSideRight, false},
		{"invalid", "center", DefendingSide(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DefendingSideFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefendingSideFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("DefendingSideFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefendingSide_JSON(t *testing.T) {
	tests := []struct {
		name string
		side DefendingSide
		json string
	}{
		{"left", DefendingSideLeft, `"left"`},
		{"right", DefendingSideRight, `"right"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.side)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(data) != tt.json {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.json)
			}

			var got DefendingSide
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got != tt.side {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.side)
			}
		})
	}
}

// GameScheduleState tests

func TestGameScheduleState_String(t *testing.T) {
	tests := []struct {
		name  string
		state GameScheduleState
		want  string
	}{
		{"ok", GameScheduleStateOK, "OK"},
		{"dont play", GameScheduleStateDontPlay, "DONT_PLAY"},
		{"postponed", GameScheduleStatePostponed, "PPD"},
		{"suspended", GameScheduleStateSuspended, "SUSP"},
		{"tbd", GameScheduleStateTBD, "TBD"},
		{"completed", GameScheduleStateCompleted, "COMPLETED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("GameScheduleState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameScheduleState_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		state GameScheduleState
		want  bool
	}{
		{"ok valid", GameScheduleStateOK, true},
		{"dont play valid", GameScheduleStateDontPlay, true},
		{"postponed valid", GameScheduleStatePostponed, true},
		{"suspended valid", GameScheduleStateSuspended, true},
		{"tbd valid", GameScheduleStateTBD, true},
		{"completed valid", GameScheduleStateCompleted, true},
		{"invalid", GameScheduleState("INVALID"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsValid(); got != tt.want {
				t.Errorf("GameScheduleState.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameScheduleStateFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    GameScheduleState
		wantErr bool
	}{
		{"ok", "OK", GameScheduleStateOK, false},
		{"dont play", "DONT_PLAY", GameScheduleStateDontPlay, false},
		{"postponed", "PPD", GameScheduleStatePostponed, false},
		{"suspended", "SUSP", GameScheduleStateSuspended, false},
		{"tbd", "TBD", GameScheduleStateTBD, false},
		{"completed", "COMPLETED", GameScheduleStateCompleted, false},
		{"invalid", "INVALID", GameScheduleState(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GameScheduleStateFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameScheduleStateFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GameScheduleStateFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameScheduleState_JSON(t *testing.T) {
	tests := []struct {
		name  string
		state GameScheduleState
		json  string
	}{
		{"ok", GameScheduleStateOK, `"OK"`},
		{"dont play", GameScheduleStateDontPlay, `"DONT_PLAY"`},
		{"postponed", GameScheduleStatePostponed, `"PPD"`},
		{"suspended", GameScheduleStateSuspended, `"SUSP"`},
		{"tbd", GameScheduleStateTBD, `"TBD"`},
		{"completed", GameScheduleStateCompleted, `"COMPLETED"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.state)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(data) != tt.json {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.json)
			}

			var got GameScheduleState
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got != tt.state {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.state)
			}
		})
	}
}

// PlayEventType tests

func TestPlayEventType_String(t *testing.T) {
	tests := []struct {
		name      string
		eventType PlayEventType
		want      string
	}{
		{"game start", PlayEventTypeGameStart, "game-start"},
		{"goal", PlayEventTypeGoal, "goal"},
		{"shot on goal", PlayEventTypeShotOnGoal, "shot-on-goal"},
		{"penalty", PlayEventTypePenalty, "penalty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eventType.String(); got != tt.want {
				t.Errorf("PlayEventType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayEventType_IsScoringChance(t *testing.T) {
	tests := []struct {
		name      string
		eventType PlayEventType
		want      bool
	}{
		{"shot on goal is scoring chance", PlayEventTypeShotOnGoal, true},
		{"missed shot is scoring chance", PlayEventTypeMissedShot, true},
		{"goal is scoring chance", PlayEventTypeGoal, true},
		{"blocked shot is scoring chance", PlayEventTypeBlockedShot, true},
		{"penalty not scoring chance", PlayEventTypePenalty, false},
		{"faceoff not scoring chance", PlayEventTypeFaceoff, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eventType.IsScoringChance(); got != tt.want {
				t.Errorf("PlayEventType.IsScoringChance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayEventType_IsGoal(t *testing.T) {
	tests := []struct {
		name      string
		eventType PlayEventType
		want      bool
	}{
		{"goal is goal", PlayEventTypeGoal, true},
		{"shot on goal not goal", PlayEventTypeShotOnGoal, false},
		{"penalty not goal", PlayEventTypePenalty, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eventType.IsGoal(); got != tt.want {
				t.Errorf("PlayEventType.IsGoal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayEventType_IsPeriodBoundary(t *testing.T) {
	tests := []struct {
		name      string
		eventType PlayEventType
		want      bool
	}{
		{"game start is boundary", PlayEventTypeGameStart, true},
		{"period start is boundary", PlayEventTypePeriodStart, true},
		{"period end is boundary", PlayEventTypePeriodEnd, true},
		{"game end is boundary", PlayEventTypeGameEnd, true},
		{"goal not boundary", PlayEventTypeGoal, false},
		{"faceoff not boundary", PlayEventTypeFaceoff, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eventType.IsPeriodBoundary(); got != tt.want {
				t.Errorf("PlayEventType.IsPeriodBoundary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayEventType_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		eventType PlayEventType
		want      bool
	}{
		{"game start valid", PlayEventTypeGameStart, true},
		{"goal valid", PlayEventTypeGoal, true},
		{"unknown valid", PlayEventTypeUnknown, true},
		{"invalid", PlayEventType("invalid-event"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.eventType.IsValid(); got != tt.want {
				t.Errorf("PlayEventType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayEventTypeFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    PlayEventType
		wantErr bool
	}{
		{"game start", "game-start", PlayEventTypeGameStart, false},
		{"goal", "goal", PlayEventTypeGoal, false},
		{"shot on goal", "shot-on-goal", PlayEventTypeShotOnGoal, false},
		{"penalty", "penalty", PlayEventTypePenalty, false},
		{"unknown", "unknown", PlayEventTypeUnknown, false},
		{"invalid", "invalid-event", PlayEventType(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PlayEventTypeFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlayEventTypeFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("PlayEventTypeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayEventType_JSON(t *testing.T) {
	tests := []struct {
		name      string
		eventType PlayEventType
		json      string
	}{
		{"game start", PlayEventTypeGameStart, `"game-start"`},
		{"goal", PlayEventTypeGoal, `"goal"`},
		{"shot on goal", PlayEventTypeShotOnGoal, `"shot-on-goal"`},
		{"penalty", PlayEventTypePenalty, `"penalty"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.eventType)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(data) != tt.json {
				t.Errorf("json.Marshal() = %v, want %v", string(data), tt.json)
			}

			var got PlayEventType
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got != tt.eventType {
				t.Errorf("json.Unmarshal() = %v, want %v", got, tt.eventType)
			}
		})
	}
}

// Additional tests for uncovered String() methods

func TestPosition_String(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		want     string
	}{
		{"center", PositionCenter, "Center"},
		{"defense", PositionDefense, "Defense"},
		{"goalie", PositionGoalie, "Goalie"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.String(); got != tt.want {
				t.Errorf("Position.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandedness_String(t *testing.T) {
	tests := []struct {
		name       string
		handedness Handedness
		want       string
	}{
		{"left", HandednessLeft, "Left"},
		{"right", HandednessRight, "Right"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.handedness.String(); got != tt.want {
				t.Errorf("Handedness.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeriodType_String(t *testing.T) {
	tests := []struct {
		name       string
		periodType PeriodType
		want       string
	}{
		{"regulation", PeriodTypeRegulation, "Regulation"},
		{"overtime", PeriodTypeOvertime, "Overtime"},
		{"shootout", PeriodTypeShootout, "Shootout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.periodType.String(); got != tt.want {
				t.Errorf("PeriodType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHomeRoad_String(t *testing.T) {
	tests := []struct {
		name     string
		homeRoad HomeRoad
		want     string
	}{
		{"home", HomeRoadHome, "Home"},
		{"road", HomeRoadRoad, "Road"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.homeRoad.String(); got != tt.want {
				t.Errorf("HomeRoad.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Tests for Must*FromString functions

func TestMustHandednessFromString(t *testing.T) {
	t.Run("valid handedness", func(t *testing.T) {
		got := MustHandednessFromString("L")
		if got != HandednessLeft {
			t.Errorf("MustHandednessFromString() = %v, want %v", got, HandednessLeft)
		}
	})

	t.Run("invalid handedness panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustHandednessFromString() did not panic")
			}
		}()
		MustHandednessFromString("INVALID")
	})
}

func TestMustGoalieDecisionFromString(t *testing.T) {
	t.Run("valid decision", func(t *testing.T) {
		got := MustGoalieDecisionFromString("W")
		if got != GoalieDecisionWin {
			t.Errorf("MustGoalieDecisionFromString() = %v, want %v", got, GoalieDecisionWin)
		}
	})

	t.Run("invalid decision panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustGoalieDecisionFromString() did not panic")
			}
		}()
		MustGoalieDecisionFromString("INVALID")
	})
}

func TestMustPeriodTypeFromString(t *testing.T) {
	t.Run("valid period type", func(t *testing.T) {
		got := MustPeriodTypeFromString("REG")
		if got != PeriodTypeRegulation {
			t.Errorf("MustPeriodTypeFromString() = %v, want %v", got, PeriodTypeRegulation)
		}
	})

	t.Run("invalid period type panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustPeriodTypeFromString() did not panic")
			}
		}()
		MustPeriodTypeFromString("INVALID")
	})
}

func TestMustHomeRoadFromString(t *testing.T) {
	t.Run("valid home/road", func(t *testing.T) {
		got := MustHomeRoadFromString("H")
		if got != HomeRoadHome {
			t.Errorf("MustHomeRoadFromString() = %v, want %v", got, HomeRoadHome)
		}
	})

	t.Run("invalid home/road panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustHomeRoadFromString() did not panic")
			}
		}()
		MustHomeRoadFromString("INVALID")
	})
}

func TestMustZoneCodeFromString(t *testing.T) {
	t.Run("valid zone code", func(t *testing.T) {
		got := MustZoneCodeFromString("O")
		if got != ZoneCodeOffensive {
			t.Errorf("MustZoneCodeFromString() = %v, want %v", got, ZoneCodeOffensive)
		}
	})

	t.Run("invalid zone code panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustZoneCodeFromString() did not panic")
			}
		}()
		MustZoneCodeFromString("INVALID")
	})
}

func TestMustDefendingSideFromString(t *testing.T) {
	t.Run("valid defending side", func(t *testing.T) {
		got := MustDefendingSideFromString("left")
		if got != DefendingSideLeft {
			t.Errorf("MustDefendingSideFromString() = %v, want %v", got, DefendingSideLeft)
		}
	})

	t.Run("invalid defending side panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustDefendingSideFromString() did not panic")
			}
		}()
		MustDefendingSideFromString("INVALID")
	})
}

func TestMustGameScheduleStateFromString(t *testing.T) {
	t.Run("valid schedule state", func(t *testing.T) {
		got := MustGameScheduleStateFromString("OK")
		if got != GameScheduleStateOK {
			t.Errorf("MustGameScheduleStateFromString() = %v, want %v", got, GameScheduleStateOK)
		}
	})

	t.Run("invalid schedule state panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustGameScheduleStateFromString() did not panic")
			}
		}()
		MustGameScheduleStateFromString("INVALID")
	})
}

func TestMustPlayEventTypeFromString(t *testing.T) {
	t.Run("valid event type", func(t *testing.T) {
		got := MustPlayEventTypeFromString("goal")
		if got != PlayEventTypeGoal {
			t.Errorf("MustPlayEventTypeFromString() = %v, want %v", got, PlayEventTypeGoal)
		}
	})

	t.Run("invalid event type panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustPlayEventTypeFromString() did not panic")
			}
		}()
		MustPlayEventTypeFromString("INVALID")
	})
}

// Tests for MarshalJSON error paths

func TestPosition_MarshalJSON_Invalid(t *testing.T) {
	invalid := Position("INVALID")
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Error("MarshalJSON() should error for invalid position")
	}
}

func TestHandedness_MarshalJSON_Empty(t *testing.T) {
	// Empty handedness should marshal successfully to support players
	// with missing data from the NHL API (e.g., player 8449312)
	empty := Handedness("")
	data, err := json.Marshal(empty)
	if err != nil {
		t.Errorf("MarshalJSON() unexpected error: %v", err)
	}
	if string(data) != `""` {
		t.Errorf("MarshalJSON() = %s, want \"\"", string(data))
	}

	// Non-empty invalid values also marshal as-is
	invalid := Handedness("INVALID")
	data, err = json.Marshal(invalid)
	if err != nil {
		t.Errorf("MarshalJSON() unexpected error: %v", err)
	}
	if string(data) != `"INVALID"` {
		t.Errorf("MarshalJSON() = %s, want \"INVALID\"", string(data))
	}
}

func TestGoalieDecision_MarshalJSON_Invalid(t *testing.T) {
	invalid := GoalieDecision("INVALID")
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Error("MarshalJSON() should error for invalid goalie decision")
	}
}

func TestPeriodType_MarshalJSON_Invalid(t *testing.T) {
	invalid := PeriodType("INVALID")
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Error("MarshalJSON() should error for invalid period type")
	}
}

func TestHomeRoad_MarshalJSON_Invalid(t *testing.T) {
	invalid := HomeRoad("INVALID")
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Error("MarshalJSON() should error for invalid home/road")
	}
}

func TestZoneCode_MarshalJSON_Invalid(t *testing.T) {
	invalid := ZoneCode("INVALID")
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Error("MarshalJSON() should error for invalid zone code")
	}
}

func TestDefendingSide_MarshalJSON_Invalid(t *testing.T) {
	invalid := DefendingSide("INVALID")
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Error("MarshalJSON() should error for invalid defending side")
	}
}

func TestGameScheduleState_MarshalJSON_Invalid(t *testing.T) {
	invalid := GameScheduleState("INVALID")
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Error("MarshalJSON() should error for invalid game schedule state")
	}
}

func TestPlayEventType_MarshalJSON_Invalid(t *testing.T) {
	invalid := PlayEventType("INVALID")
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Error("MarshalJSON() should error for invalid play event type")
	}
}

func TestPeriodType_IsValid(t *testing.T) {
	tests := []struct {
		name       string
		periodType PeriodType
		want       bool
	}{
		{"regulation valid", PeriodTypeRegulation, true},
		{"overtime valid", PeriodTypeOvertime, true},
		{"shootout valid", PeriodTypeShootout, true},
		{"invalid", PeriodType("INVALID"), false},
		{"empty invalid", PeriodType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.periodType.IsValid(); got != tt.want {
				t.Errorf("PeriodType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZoneCode_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		zoneCode ZoneCode
		want     bool
	}{
		{"offensive valid", ZoneCodeOffensive, true},
		{"defensive valid", ZoneCodeDefensive, true},
		{"neutral valid", ZoneCodeNeutral, true},
		{"invalid", ZoneCode("INVALID"), false},
		{"empty invalid", ZoneCode(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.zoneCode.IsValid(); got != tt.want {
				t.Errorf("ZoneCode.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Additional UnmarshalJSON error path tests

func TestPosition_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var p Position
	err := json.Unmarshal([]byte(`{invalid json}`), &p)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestHandedness_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var h Handedness
	err := json.Unmarshal([]byte(`{invalid json}`), &h)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestGoalieDecision_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var g GoalieDecision
	err := json.Unmarshal([]byte(`{invalid json}`), &g)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestPeriodType_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var p PeriodType
	err := json.Unmarshal([]byte(`{invalid json}`), &p)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestHomeRoad_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var h HomeRoad
	err := json.Unmarshal([]byte(`{invalid json}`), &h)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestZoneCode_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var z ZoneCode
	err := json.Unmarshal([]byte(`{invalid json}`), &z)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestDefendingSide_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var d DefendingSide
	err := json.Unmarshal([]byte(`{invalid json}`), &d)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestGameScheduleState_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var g GameScheduleState
	err := json.Unmarshal([]byte(`{invalid json}`), &g)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestPlayEventType_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var p PlayEventType
	err := json.Unmarshal([]byte(`{invalid json}`), &p)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

// Additional error path tests for UnmarshalJSON - invalid value errors

func TestPosition_UnmarshalJSON_InvalidValue(t *testing.T) {
	var p Position
	err := json.Unmarshal([]byte(`"INVALID"`), &p)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid position value")
	}
}

func TestHandedness_UnmarshalJSON_InvalidValue(t *testing.T) {
	var h Handedness
	err := json.Unmarshal([]byte(`"INVALID"`), &h)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid handedness value")
	}
}

func TestGoalieDecision_UnmarshalJSON_InvalidValue(t *testing.T) {
	var g GoalieDecision
	err := json.Unmarshal([]byte(`"INVALID"`), &g)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid goalie decision value")
	}
}

func TestPeriodType_UnmarshalJSON_InvalidValue(t *testing.T) {
	var p PeriodType
	err := json.Unmarshal([]byte(`"INVALID"`), &p)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid period type value")
	}
}

func TestHomeRoad_UnmarshalJSON_InvalidValue(t *testing.T) {
	var h HomeRoad
	err := json.Unmarshal([]byte(`"INVALID"`), &h)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid home/road value")
	}
}

func TestZoneCode_UnmarshalJSON_InvalidValue(t *testing.T) {
	var z ZoneCode
	err := json.Unmarshal([]byte(`"INVALID"`), &z)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid zone code value")
	}
}

func TestDefendingSide_UnmarshalJSON_InvalidValue(t *testing.T) {
	var d DefendingSide
	err := json.Unmarshal([]byte(`"INVALID"`), &d)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid defending side value")
	}
}

func TestGameScheduleState_UnmarshalJSON_InvalidValue(t *testing.T) {
	var g GameScheduleState
	err := json.Unmarshal([]byte(`"INVALID"`), &g)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid game schedule state value")
	}
}

func TestPlayEventType_UnmarshalJSON_InvalidValue(t *testing.T) {
	var p PlayEventType
	err := json.Unmarshal([]byte(`"INVALID"`), &p)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid play event type value")
	}
}
