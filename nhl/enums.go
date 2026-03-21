package nhl

// Custom domain methods for enum types. The type declarations, constants,
// and boilerplate methods (String, IsValid, FromString, JSON) are generated
// in enums_generated.go — see internal/enumgen for the generator.

// IsForward returns true if the position is a forward (C, LW, RW, or F).
func (v Position) IsForward() bool {
	return v == PositionCenter || v == PositionLeftWing || v == PositionRightWing || v == PositionForward
}

// IsSkater returns true if the position is a skater (not a goalie).
func (v Position) IsSkater() bool {
	return v != PositionGoalie && v.IsValid()
}

// IsOvertime returns true if the period is overtime or shootout.
func (v PeriodType) IsOvertime() bool {
	return v == PeriodTypeOvertime || v == PeriodTypeShootout
}

// IsScoringChance returns true if the event is a scoring chance (shot, goal, etc.).
func (v PlayEventType) IsScoringChance() bool {
	switch v {
	case PlayEventTypeShotOnGoal, PlayEventTypeMissedShot, PlayEventTypeBlockedShot, PlayEventTypeGoal:
		return true
	default:
		return false
	}
}

// IsGoal returns true if the event is a goal.
func (v PlayEventType) IsGoal() bool {
	return v == PlayEventTypeGoal
}

// IsPeriodBoundary returns true if the event marks the start or end of a period.
func (v PlayEventType) IsPeriodBoundary() bool {
	switch v {
	case PlayEventTypeGameStart, PlayEventTypePeriodStart, PlayEventTypePeriodEnd, PlayEventTypeGameEnd:
		return true
	default:
		return false
	}
}
