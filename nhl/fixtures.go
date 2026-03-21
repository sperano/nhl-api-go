package nhl

// Test fixtures for constructing valid NHL objects that survive json.Marshal round-trips.
// NHL enum types (GameType, GameState, PeriodType, GameScheduleState) reject zero values
// during marshaling. These constructors provide the minimum valid fields.

// FixtureBoxscore returns a Boxscore with the minimum fields needed for json.Marshal.
func FixtureBoxscore() *Boxscore {
	return &Boxscore{
		GameType:          GameTypeRegularSeason,
		GameState:         GameStateFinal,
		GameScheduleState: GameScheduleStateOK,
		PeriodDescriptor:  PeriodDescriptor{PeriodType: PeriodTypeRegulation},
	}
}

// FixturePlayByPlay returns a PlayByPlay with the minimum fields needed for json.Marshal.
func FixturePlayByPlay() *PlayByPlay {
	return &PlayByPlay{
		GameType:          GameTypeRegularSeason,
		GameState:         GameStateFinal,
		GameScheduleState: GameScheduleStateOK,
		PeriodDescriptor:  PeriodDescriptor{PeriodType: PeriodTypeRegulation},
	}
}

// FixtureGameStory returns a GameStory with the minimum fields needed for json.Marshal.
func FixtureGameStory() *GameStory {
	return &GameStory{
		GameType:          GameTypeRegularSeason,
		GameState:         GameStateFinal,
		GameScheduleState: GameScheduleStateOK,
	}
}

// FixtureShiftChart returns a ShiftChart. ShiftChart has no strict enum fields,
// so a zero value is valid, but this is provided for consistency.
func FixtureShiftChart() *ShiftChart {
	return &ShiftChart{}
}

// FixtureSeasonSeriesMatchup returns a SeasonSeriesMatchup. The struct has no
// top-level strict enum fields (enums are only in nested SeriesGame slices).
func FixtureSeasonSeriesMatchup() *SeasonSeriesMatchup {
	return &SeasonSeriesMatchup{}
}
