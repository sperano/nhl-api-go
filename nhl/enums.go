package nhl

import (
	"encoding/json"
	"fmt"
)

// Position represents a player's position on the ice.
type Position string

const (
	// PositionCenter represents a center forward.
	PositionCenter Position = "C"
	// PositionLeftWing represents a left wing forward.
	PositionLeftWing Position = "LW"
	// PositionRightWing represents a right wing forward.
	PositionRightWing Position = "RW"
	// PositionDefense represents a defenseman.
	PositionDefense Position = "D"
	// PositionGoalie represents a goaltender.
	PositionGoalie Position = "G"
)

// Code returns the position code (e.g., "C", "LW", "RW", "D", "G").
func (p Position) Code() string {
	return string(p)
}

// Name returns the full name of the position.
func (p Position) Name() string {
	switch p {
	case PositionCenter:
		return "Center"
	case PositionLeftWing:
		return "Left Wing"
	case PositionRightWing:
		return "Right Wing"
	case PositionDefense:
		return "Defense"
	case PositionGoalie:
		return "Goalie"
	default:
		return fmt.Sprintf("Unknown(%s)", string(p))
	}
}

// String returns the full name of the position.
func (p Position) String() string {
	return p.Name()
}

// IsForward returns true if the position is a forward (C, LW, or RW).
func (p Position) IsForward() bool {
	return p == PositionCenter || p == PositionLeftWing || p == PositionRightWing
}

// IsSkater returns true if the position is a skater (not a goalie).
func (p Position) IsSkater() bool {
	return p != PositionGoalie && p.IsValid()
}

// IsValid returns true if the Position is one of the known valid positions.
func (p Position) IsValid() bool {
	switch p {
	case PositionCenter, PositionLeftWing, PositionRightWing, PositionDefense, PositionGoalie:
		return true
	default:
		return false
	}
}

// PositionFromString parses a string into a Position.
// Accepts codes ("C", "L", "LW", "R", "RW", "D", "G") and full names ("Center", "Left Wing", etc.).
// Returns an error if the string is not a valid Position.
func PositionFromString(s string) (Position, error) {
	switch s {
	case "C", "Center":
		return PositionCenter, nil
	case "L", "LW", "Left Wing", "LeftWing":
		return PositionLeftWing, nil
	case "R", "RW", "Right Wing", "RightWing":
		return PositionRightWing, nil
	case "D", "Defense", "Defenseman":
		return PositionDefense, nil
	case "G", "Goalie", "Goaltender":
		return PositionGoalie, nil
	default:
		return "", fmt.Errorf("invalid position: %q", s)
	}
}

// MustPositionFromString parses a string into a Position.
// Panics if the string is not a valid Position.
func MustPositionFromString(s string) Position {
	p, err := PositionFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}

// UnmarshalJSON implements custom JSON unmarshaling for Position.
func (p *Position) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	position, err := PositionFromString(s)
	if err != nil {
		return err
	}

	*p = position
	return nil
}

// MarshalJSON implements custom JSON marshaling for Position.
func (p Position) MarshalJSON() ([]byte, error) {
	if !p.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid position: %q", string(p))
	}
	return json.Marshal(p.Code())
}

// Handedness represents a player's shooting or catching hand.
type Handedness string

const (
	// HandednessLeft represents left-handed.
	HandednessLeft Handedness = "L"
	// HandednessRight represents right-handed.
	HandednessRight Handedness = "R"
)

// Code returns the handedness code ("L" or "R").
func (h Handedness) Code() string {
	return string(h)
}

// Name returns the full name of the handedness.
func (h Handedness) Name() string {
	switch h {
	case HandednessLeft:
		return "Left"
	case HandednessRight:
		return "Right"
	default:
		return fmt.Sprintf("Unknown(%s)", string(h))
	}
}

// String returns the full name of the handedness.
func (h Handedness) String() string {
	return h.Name()
}

// IsValid returns true if the Handedness is one of the known valid values.
func (h Handedness) IsValid() bool {
	return h == HandednessLeft || h == HandednessRight
}

// HandednessFromString parses a string into a Handedness.
// Accepts both codes ("L", "R") and full names ("Left", "Right").
// Returns an error if the string is not a valid Handedness.
func HandednessFromString(s string) (Handedness, error) {
	switch s {
	case "L", "Left":
		return HandednessLeft, nil
	case "R", "Right":
		return HandednessRight, nil
	default:
		return "", fmt.Errorf("invalid handedness: %q", s)
	}
}

// MustHandednessFromString parses a string into a Handedness.
// Panics if the string is not a valid Handedness.
func MustHandednessFromString(s string) Handedness {
	h, err := HandednessFromString(s)
	if err != nil {
		panic(err)
	}
	return h
}

// UnmarshalJSON implements custom JSON unmarshaling for Handedness.
// Empty strings are accepted to support players with missing data from the NHL API.
func (h *Handedness) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Allow empty strings for players with missing handedness data
	if s == "" {
		*h = Handedness("")
		return nil
	}

	handedness, err := HandednessFromString(s)
	if err != nil {
		return err
	}

	*h = handedness
	return nil
}

// MarshalJSON implements custom JSON marshaling for Handedness.
// Empty handedness marshals as an empty string to support players with
// missing data from the NHL API.
func (h Handedness) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(h))
}

// GoalieDecision represents the decision (result) for a goalie in a game.
type GoalieDecision string

const (
	// GoalieDecisionWin represents a win.
	GoalieDecisionWin GoalieDecision = "W"
	// GoalieDecisionLoss represents a loss.
	GoalieDecisionLoss GoalieDecision = "L"
	// GoalieDecisionTie represents a tie (used in older games).
	GoalieDecisionTie GoalieDecision = "T"
	// GoalieDecisionOvertimeLoss represents an overtime/shootout loss.
	GoalieDecisionOvertimeLoss GoalieDecision = "OTL"
)

// String returns the string representation of the GoalieDecision.
func (g GoalieDecision) String() string {
	switch g {
	case GoalieDecisionWin:
		return "Win"
	case GoalieDecisionLoss:
		return "Loss"
	case GoalieDecisionTie:
		return "Tie"
	case GoalieDecisionOvertimeLoss:
		return "Overtime Loss"
	default:
		return fmt.Sprintf("Unknown(%s)", string(g))
	}
}

// IsValid returns true if the GoalieDecision is one of the known valid values.
func (g GoalieDecision) IsValid() bool {
	switch g {
	case GoalieDecisionWin, GoalieDecisionLoss, GoalieDecisionTie, GoalieDecisionOvertimeLoss:
		return true
	default:
		return false
	}
}

// GoalieDecisionFromString parses a string into a GoalieDecision.
// Returns an error if the string is not a valid GoalieDecision.
func GoalieDecisionFromString(s string) (GoalieDecision, error) {
	switch s {
	case "W", "Win":
		return GoalieDecisionWin, nil
	case "L", "Loss":
		return GoalieDecisionLoss, nil
	case "T", "Tie":
		return GoalieDecisionTie, nil
	case "O", "OTL", "Overtime Loss", "OvertimeLoss":
		return GoalieDecisionOvertimeLoss, nil
	default:
		return "", fmt.Errorf("invalid goalie decision: %q", s)
	}
}

// MustGoalieDecisionFromString parses a string into a GoalieDecision.
// Panics if the string is not a valid GoalieDecision.
func MustGoalieDecisionFromString(s string) GoalieDecision {
	g, err := GoalieDecisionFromString(s)
	if err != nil {
		panic(err)
	}
	return g
}

// UnmarshalJSON implements custom JSON unmarshaling for GoalieDecision.
func (g *GoalieDecision) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	decision, err := GoalieDecisionFromString(s)
	if err != nil {
		return err
	}

	*g = decision
	return nil
}

// MarshalJSON implements custom JSON marshaling for GoalieDecision.
func (g GoalieDecision) MarshalJSON() ([]byte, error) {
	if !g.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid goalie decision: %q", string(g))
	}
	return json.Marshal(string(g))
}

// PeriodType represents the type of period in a hockey game.
type PeriodType string

const (
	// PeriodTypeRegulation represents a regulation period.
	PeriodTypeRegulation PeriodType = "REG"
	// PeriodTypeOvertime represents an overtime period.
	PeriodTypeOvertime PeriodType = "OT"
	// PeriodTypeShootout represents a shootout.
	PeriodTypeShootout PeriodType = "SO"
)

// Code returns the period type code (e.g., "REG", "OT", "SO").
func (p PeriodType) Code() string {
	return string(p)
}

// Name returns the full name of the period type.
func (p PeriodType) Name() string {
	switch p {
	case PeriodTypeRegulation:
		return "Regulation"
	case PeriodTypeOvertime:
		return "Overtime"
	case PeriodTypeShootout:
		return "Shootout"
	default:
		return fmt.Sprintf("Unknown(%s)", string(p))
	}
}

// String returns the full name of the period type.
func (p PeriodType) String() string {
	return p.Name()
}

// IsOvertime returns true if the period is overtime or shootout.
func (p PeriodType) IsOvertime() bool {
	return p == PeriodTypeOvertime || p == PeriodTypeShootout
}

// IsValid returns true if the PeriodType is one of the known valid types.
func (p PeriodType) IsValid() bool {
	switch p {
	case PeriodTypeRegulation, PeriodTypeOvertime, PeriodTypeShootout:
		return true
	default:
		return false
	}
}

// PeriodTypeFromString parses a string into a PeriodType.
// Accepts both codes ("REG", "OT", "SO") and full names.
// Returns an error if the string is not a valid PeriodType.
func PeriodTypeFromString(s string) (PeriodType, error) {
	switch s {
	case "REG", "Regulation":
		return PeriodTypeRegulation, nil
	case "OT", "Overtime":
		return PeriodTypeOvertime, nil
	case "SO", "Shootout":
		return PeriodTypeShootout, nil
	default:
		return "", fmt.Errorf("invalid period type: %q", s)
	}
}

// MustPeriodTypeFromString parses a string into a PeriodType.
// Panics if the string is not a valid PeriodType.
func MustPeriodTypeFromString(s string) PeriodType {
	p, err := PeriodTypeFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}

// UnmarshalJSON implements custom JSON unmarshaling for PeriodType.
func (p *PeriodType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	periodType, err := PeriodTypeFromString(s)
	if err != nil {
		return err
	}

	*p = periodType
	return nil
}

// MarshalJSON implements custom JSON marshaling for PeriodType.
func (p PeriodType) MarshalJSON() ([]byte, error) {
	if !p.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid period type: %q", string(p))
	}
	return json.Marshal(p.Code())
}

// HomeRoad represents whether a team is home or away (road).
type HomeRoad string

const (
	// HomeRoadHome represents the home team.
	HomeRoadHome HomeRoad = "H"
	// HomeRoadRoad represents the away/road team.
	HomeRoadRoad HomeRoad = "R"
)

// Code returns the home/road code ("H" or "R").
func (h HomeRoad) Code() string {
	return string(h)
}

// Name returns the full name (Home or Road).
func (h HomeRoad) Name() string {
	switch h {
	case HomeRoadHome:
		return "Home"
	case HomeRoadRoad:
		return "Road"
	default:
		return fmt.Sprintf("Unknown(%s)", string(h))
	}
}

// String returns the full name (Home or Road).
func (h HomeRoad) String() string {
	return h.Name()
}

// IsValid returns true if the HomeRoad is one of the known valid values.
func (h HomeRoad) IsValid() bool {
	return h == HomeRoadHome || h == HomeRoadRoad
}

// HomeRoadFromString parses a string into a HomeRoad.
// Accepts both codes ("H", "R") and full names ("Home", "Road", "Away").
// Returns an error if the string is not a valid HomeRoad.
func HomeRoadFromString(s string) (HomeRoad, error) {
	switch s {
	case "H", "Home":
		return HomeRoadHome, nil
	case "R", "Road", "Away":
		return HomeRoadRoad, nil
	default:
		return "", fmt.Errorf("invalid home/road: %q", s)
	}
}

// MustHomeRoadFromString parses a string into a HomeRoad.
// Panics if the string is not a valid HomeRoad.
func MustHomeRoadFromString(s string) HomeRoad {
	h, err := HomeRoadFromString(s)
	if err != nil {
		panic(err)
	}
	return h
}

// UnmarshalJSON implements custom JSON unmarshaling for HomeRoad.
func (h *HomeRoad) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	homeRoad, err := HomeRoadFromString(s)
	if err != nil {
		return err
	}

	*h = homeRoad
	return nil
}

// MarshalJSON implements custom JSON marshaling for HomeRoad.
func (h HomeRoad) MarshalJSON() ([]byte, error) {
	if !h.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid home/road: %q", string(h))
	}
	return json.Marshal(h.Code())
}

// ZoneCode represents a zone on the ice.
type ZoneCode string

const (
	// ZoneCodeOffensive represents the offensive zone.
	ZoneCodeOffensive ZoneCode = "O"
	// ZoneCodeDefensive represents the defensive zone.
	ZoneCodeDefensive ZoneCode = "D"
	// ZoneCodeNeutral represents the neutral zone.
	ZoneCodeNeutral ZoneCode = "N"
)

// Code returns the zone code ("O", "D", or "N").
func (z ZoneCode) Code() string {
	return string(z)
}

// String returns the full name of the zone.
func (z ZoneCode) String() string {
	switch z {
	case ZoneCodeOffensive:
		return "Offensive"
	case ZoneCodeDefensive:
		return "Defensive"
	case ZoneCodeNeutral:
		return "Neutral"
	default:
		return fmt.Sprintf("Unknown(%s)", string(z))
	}
}

// IsValid returns true if the ZoneCode is one of the known valid codes.
func (z ZoneCode) IsValid() bool {
	switch z {
	case ZoneCodeOffensive, ZoneCodeDefensive, ZoneCodeNeutral:
		return true
	default:
		return false
	}
}

// ZoneCodeFromString parses a string into a ZoneCode.
// Accepts both codes ("O", "D", "N") and full names.
// Returns an error if the string is not a valid ZoneCode.
func ZoneCodeFromString(s string) (ZoneCode, error) {
	switch s {
	case "O", "Offensive":
		return ZoneCodeOffensive, nil
	case "D", "Defensive":
		return ZoneCodeDefensive, nil
	case "N", "Neutral":
		return ZoneCodeNeutral, nil
	default:
		return "", fmt.Errorf("invalid zone code: %q", s)
	}
}

// MustZoneCodeFromString parses a string into a ZoneCode.
// Panics if the string is not a valid ZoneCode.
func MustZoneCodeFromString(s string) ZoneCode {
	z, err := ZoneCodeFromString(s)
	if err != nil {
		panic(err)
	}
	return z
}

// UnmarshalJSON implements custom JSON unmarshaling for ZoneCode.
func (z *ZoneCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	zoneCode, err := ZoneCodeFromString(s)
	if err != nil {
		return err
	}

	*z = zoneCode
	return nil
}

// MarshalJSON implements custom JSON marshaling for ZoneCode.
func (z ZoneCode) MarshalJSON() ([]byte, error) {
	if !z.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid zone code: %q", string(z))
	}
	return json.Marshal(z.Code())
}

// DefendingSide represents which side of the ice a team is defending.
type DefendingSide string

const (
	// DefendingSideLeft represents defending the left side.
	DefendingSideLeft DefendingSide = "left"
	// DefendingSideRight represents defending the right side.
	DefendingSideRight DefendingSide = "right"
)

// String returns the string representation of the DefendingSide.
func (d DefendingSide) String() string {
	return string(d)
}

// IsValid returns true if the DefendingSide is one of the known valid values.
func (d DefendingSide) IsValid() bool {
	return d == DefendingSideLeft || d == DefendingSideRight
}

// DefendingSideFromString parses a string into a DefendingSide.
// Returns an error if the string is not a valid DefendingSide.
func DefendingSideFromString(s string) (DefendingSide, error) {
	switch s {
	case "left":
		return DefendingSideLeft, nil
	case "right":
		return DefendingSideRight, nil
	default:
		return "", fmt.Errorf("invalid defending side: %q", s)
	}
}

// MustDefendingSideFromString parses a string into a DefendingSide.
// Panics if the string is not a valid DefendingSide.
func MustDefendingSideFromString(s string) DefendingSide {
	d, err := DefendingSideFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

// UnmarshalJSON implements custom JSON unmarshaling for DefendingSide.
func (d *DefendingSide) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	side, err := DefendingSideFromString(s)
	if err != nil {
		return err
	}

	*d = side
	return nil
}

// MarshalJSON implements custom JSON marshaling for DefendingSide.
// Empty strings are allowed for historical games that lack this data.
func (d DefendingSide) MarshalJSON() ([]byte, error) {
	// Allow empty strings for old games that don't have defending side data
	if d == "" {
		return json.Marshal("")
	}
	if !d.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid defending side: %q", string(d))
	}
	return json.Marshal(string(d))
}

// GameScheduleState represents the state of a scheduled game.
type GameScheduleState string

const (
	// GameScheduleStateOK represents a normal scheduled game.
	GameScheduleStateOK GameScheduleState = "OK"
	// GameScheduleStateDontPlay represents a game that won't be played.
	GameScheduleStateDontPlay GameScheduleState = "DONT_PLAY"
	// GameScheduleStatePostponed represents a postponed game.
	GameScheduleStatePostponed GameScheduleState = "PPD"
	// GameScheduleStateSuspended represents a suspended game.
	GameScheduleStateSuspended GameScheduleState = "SUSP"
	// GameScheduleStateTBD represents a game with time to be determined.
	GameScheduleStateTBD GameScheduleState = "TBD"
	// GameScheduleStateCompleted represents a completed game.
	GameScheduleStateCompleted GameScheduleState = "COMPLETED"
	// GameScheduleStateCancelled represents a cancelled game.
	GameScheduleStateCancelled GameScheduleState = "CNCL"
)

// String returns the string representation of the GameScheduleState.
func (g GameScheduleState) String() string {
	return string(g)
}

// IsValid returns true if the GameScheduleState is one of the known valid states.
func (g GameScheduleState) IsValid() bool {
	switch g {
	case GameScheduleStateOK, GameScheduleStateDontPlay, GameScheduleStatePostponed,
		GameScheduleStateSuspended, GameScheduleStateTBD, GameScheduleStateCompleted,
		GameScheduleStateCancelled:
		return true
	default:
		return false
	}
}

// GameScheduleStateFromString parses a string into a GameScheduleState.
// Returns an error if the string is not a valid GameScheduleState.
func GameScheduleStateFromString(s string) (GameScheduleState, error) {
	g := GameScheduleState(s)
	if !g.IsValid() {
		return "", fmt.Errorf("invalid game schedule state: %q", s)
	}
	return g, nil
}

// MustGameScheduleStateFromString parses a string into a GameScheduleState.
// Panics if the string is not a valid GameScheduleState.
func MustGameScheduleStateFromString(s string) GameScheduleState {
	g, err := GameScheduleStateFromString(s)
	if err != nil {
		panic(err)
	}
	return g
}

// UnmarshalJSON implements custom JSON unmarshaling for GameScheduleState.
func (g *GameScheduleState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	state, err := GameScheduleStateFromString(s)
	if err != nil {
		return err
	}

	*g = state
	return nil
}

// MarshalJSON implements custom JSON marshaling for GameScheduleState.
func (g GameScheduleState) MarshalJSON() ([]byte, error) {
	if !g.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid game schedule state: %q", string(g))
	}
	return json.Marshal(string(g))
}

// PlayEventType represents a type of play event during a game.
type PlayEventType string

const (
	// PlayEventTypeGameStart represents the start of a game.
	PlayEventTypeGameStart PlayEventType = "game-start"
	// PlayEventTypePeriodStart represents the start of a period.
	PlayEventTypePeriodStart PlayEventType = "period-start"
	// PlayEventTypePeriodEnd represents the end of a period.
	PlayEventTypePeriodEnd PlayEventType = "period-end"
	// PlayEventTypeGameEnd represents the end of a game.
	PlayEventTypeGameEnd PlayEventType = "game-end"
	// PlayEventTypeFaceoff represents a faceoff.
	PlayEventTypeFaceoff PlayEventType = "faceoff"
	// PlayEventTypeHit represents a hit.
	PlayEventTypeHit PlayEventType = "hit"
	// PlayEventTypeGiveaway represents a giveaway.
	PlayEventTypeGiveaway PlayEventType = "giveaway"
	// PlayEventTypeTakeaway represents a takeaway.
	PlayEventTypeTakeaway PlayEventType = "takeaway"
	// PlayEventTypeShotOnGoal represents a shot on goal.
	PlayEventTypeShotOnGoal PlayEventType = "shot-on-goal"
	// PlayEventTypeMissedShot represents a missed shot.
	PlayEventTypeMissedShot PlayEventType = "missed-shot"
	// PlayEventTypeBlockedShot represents a blocked shot.
	PlayEventTypeBlockedShot PlayEventType = "blocked-shot"
	// PlayEventTypeGoal represents a goal.
	PlayEventTypeGoal PlayEventType = "goal"
	// PlayEventTypePenalty represents a penalty.
	PlayEventTypePenalty PlayEventType = "penalty"
	// PlayEventTypeStoppage represents a stoppage in play.
	PlayEventTypeStoppage PlayEventType = "stoppage"
	// PlayEventTypeDelayedPenalty represents a delayed penalty.
	PlayEventTypeDelayedPenalty PlayEventType = "delayed-penalty"
	// PlayEventTypeFailedShotAttempt represents a failed shot attempt.
	PlayEventTypeFailedShotAttempt PlayEventType = "failed-shot-attempt"
	// PlayEventTypeShootoutComplete represents the completion of a shootout.
	PlayEventTypeShootoutComplete PlayEventType = "shootout-complete"
	// PlayEventTypeUnknown represents an unknown event type.
	PlayEventTypeUnknown PlayEventType = "unknown"
)

// String returns the string representation of the PlayEventType.
func (p PlayEventType) String() string {
	return string(p)
}

// IsScoringChance returns true if the event is a scoring chance (shot, goal, etc.).
func (p PlayEventType) IsScoringChance() bool {
	switch p {
	case PlayEventTypeShotOnGoal, PlayEventTypeMissedShot, PlayEventTypeBlockedShot, PlayEventTypeGoal:
		return true
	default:
		return false
	}
}

// IsGoal returns true if the event is a goal.
func (p PlayEventType) IsGoal() bool {
	return p == PlayEventTypeGoal
}

// IsPeriodBoundary returns true if the event marks the start or end of a period.
func (p PlayEventType) IsPeriodBoundary() bool {
	switch p {
	case PlayEventTypeGameStart, PlayEventTypePeriodStart, PlayEventTypePeriodEnd, PlayEventTypeGameEnd:
		return true
	default:
		return false
	}
}

// IsValid returns true if the PlayEventType is one of the known valid types.
func (p PlayEventType) IsValid() bool {
	switch p {
	case PlayEventTypeGameStart, PlayEventTypePeriodStart, PlayEventTypePeriodEnd,
		PlayEventTypeGameEnd, PlayEventTypeFaceoff, PlayEventTypeHit, PlayEventTypeGiveaway,
		PlayEventTypeTakeaway, PlayEventTypeShotOnGoal, PlayEventTypeMissedShot,
		PlayEventTypeBlockedShot, PlayEventTypeGoal, PlayEventTypePenalty,
		PlayEventTypeStoppage, PlayEventTypeDelayedPenalty, PlayEventTypeFailedShotAttempt,
		PlayEventTypeShootoutComplete, PlayEventTypeUnknown:
		return true
	default:
		return false
	}
}

// PlayEventTypeFromString parses a string into a PlayEventType.
// Returns an error if the string is not a valid PlayEventType.
func PlayEventTypeFromString(s string) (PlayEventType, error) {
	p := PlayEventType(s)
	if !p.IsValid() {
		return "", fmt.Errorf("invalid play event type: %q", s)
	}
	return p, nil
}

// MustPlayEventTypeFromString parses a string into a PlayEventType.
// Panics if the string is not a valid PlayEventType.
func MustPlayEventTypeFromString(s string) PlayEventType {
	p, err := PlayEventTypeFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}

// UnmarshalJSON implements custom JSON unmarshaling for PlayEventType.
func (p *PlayEventType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	eventType, err := PlayEventTypeFromString(s)
	if err != nil {
		return err
	}

	*p = eventType
	return nil
}

// MarshalJSON implements custom JSON marshaling for PlayEventType.
func (p PlayEventType) MarshalJSON() ([]byte, error) {
	if !p.IsValid() {
		return nil, fmt.Errorf("cannot marshal invalid play event type: %q", string(p))
	}
	return json.Marshal(string(p))
}
