package nhl

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestRosterUnmarshalPositionCodes(t *testing.T) {
	const rosterJSON = `{
		"forwards": [
			{"id": 8478402, "firstName": {"default": "Connor"}, "lastName": {"default": "McDavid"}, "positionCode": "C"},
			{"id": 8476454, "firstName": {"default": "Ryan"}, "lastName": {"default": "Nugent-Hopkins"}, "positionCode": "L"},
			{"id": 8477934, "firstName": {"default": "Viktor"}, "lastName": {"default": "Arvidsson"}, "positionCode": "R"}
		],
		"defensemen": [
			{"id": 8480803, "firstName": {"default": "Evan"}, "lastName": {"default": "Bouchard"}, "positionCode": "D"}
		],
		"goalies": [
			{"id": 8479973, "firstName": {"default": "Stuart"}, "lastName": {"default": "Skinner"}, "positionCode": "G"}
		]
	}`
	wantPositions := []Position{
		PositionCenter,
		PositionLeftWing,
		PositionRightWing,
		PositionDefense,
		PositionGoalie,
	}

	var roster Roster
	if err := json.Unmarshal([]byte(rosterJSON), &roster); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	players := roster.AllPlayers()
	if len(players) != len(wantPositions) {
		t.Fatalf("Roster.AllPlayers() length = %d, want %d", len(players), len(wantPositions))
	}
	for i, want := range wantPositions {
		if players[i].Position != want {
			t.Errorf("player %d Position = %q, want %q", i, players[i].Position, want)
		}
	}
	if roster.Forwards[1].Position.Code() != "LW" {
		t.Errorf("left wing Position.Code() = %q, want %q", roster.Forwards[1].Position.Code(), "LW")
	}
	if roster.Forwards[2].Position.Code() != "RW" {
		t.Errorf("right wing Position.Code() = %q, want %q", roster.Forwards[2].Position.Code(), "RW")
	}
}

func TestRosterUnmarshalRejectsUnknownPositionCode(t *testing.T) {
	const rosterJSON = `{"forwards":[{"id":8478402,"positionCode":"X"}]}`

	var roster Roster
	err := json.Unmarshal([]byte(rosterJSON), &roster)
	if err == nil {
		t.Fatal("json.Unmarshal() error = nil, want unknown position error")
	}

	var enumErr *UnknownEnumValueError
	if !errors.As(err, &enumErr) {
		t.Fatalf("json.Unmarshal() error type = %T, want *UnknownEnumValueError", err)
	}
	if enumErr.EnumType != "position" || enumErr.Value != "X" {
		t.Errorf("UnknownEnumValueError = %#v, want position X", enumErr)
	}
}
