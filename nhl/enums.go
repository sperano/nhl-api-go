package nhl

import "fmt"

// Custom domain methods for enum types. The type declarations, constants,
// and boilerplate methods (String, IsValid, FromString, JSON) are generated
// in enums_generated.go — see internal/enumgen for the generator.

// UnknownEnumValueError indicates that a value received from the API did not
// match any known value of a strongly-typed enum. It is returned by the
// generated FromString/UnmarshalJSON methods and by GameType's equivalents,
// and is preserved through the library's JSONError wrapping. Because enum
// UnmarshalJSON failures propagate up and fail the whole response decode,
// callers can use errors.As to recover exactly which enum type and raw value
// were unrecognized — for example, to log a value the NHL API newly
// introduced before the library was updated to recognize it.
//
// Note: Go's encoding/json does not attach the JSON field path to errors
// returned by a custom UnmarshalJSON, so this error identifies the enum type
// and value but not its position within the response. The library's request
// methods wrap decode failures with the request URL for additional context.
type UnknownEnumValueError struct {
	// EnumType is a human-readable label for the enum, e.g. "game state".
	EnumType string
	// Value is the unrecognized raw value as received from the API.
	Value string
}

// Error implements the error interface. The message is intentionally identical
// to the pre-existing format ("invalid <type>: <value>") for backward
// compatibility with callers that matched on the string.
func (e *UnknownEnumValueError) Error() string {
	return fmt.Sprintf("invalid %s: %q", e.EnumType, e.Value)
}

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
