package nhl

import (
	"encoding/json"
	"testing"
)

func TestFixturesMarshalRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		obj  any
	}{
		{"Boxscore", FixtureBoxscore()},
		{"PlayByPlay", FixturePlayByPlay()},
		{"GameStory", FixtureGameStory()},
		{"ShiftChart", FixtureShiftChart()},
		{"SeasonSeriesMatchup", FixtureSeasonSeriesMatchup()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.obj)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if len(data) == 0 {
				t.Fatal("Marshal produced empty output")
			}
		})
	}
}
