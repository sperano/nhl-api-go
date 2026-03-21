package main

// EnumDef defines a string-based enum type for code generation.
type EnumDef struct {
	TypeName              string
	Doc                   string // type doc comment
	ErrorLabel            string // label for error messages (e.g., "position")
	HasCode               bool   // generate Code() method
	HasName               bool   // generate Name() method, String() returns Name()
	HasDisplayName        bool   // String() uses switch with display names (when HasName is false)
	AllowEmpty            bool   // allow empty string in UnmarshalJSON/MarshalJSON
	SkipMarshalValidation bool   // don't validate in MarshalJSON (e.g., Handedness)
	Values                []ValueDef
}

// ValueDef defines a single enum value.
type ValueDef struct {
	Name        string   // constant name (e.g., "PositionCenter")
	Value       string   // underlying string value (e.g., "C")
	DisplayName string   // display name for Name()/String() (e.g., "Center")
	Aliases     []string // all accepted strings in FromString (empty = use cast-and-validate)
	Doc         string   // doc comment
}

// HasAliases returns true if any value has explicit aliases defined.
func (e EnumDef) HasAliases() bool {
	for _, v := range e.Values {
		if len(v.Aliases) > 0 {
			return true
		}
	}
	return false
}

var enums = []EnumDef{
	{
		TypeName:   "Position",
		Doc:        "Position represents a player's position on the ice.",
		ErrorLabel: "position",
		HasCode:    true,
		HasName:    true,
		Values: []ValueDef{
			{Name: "PositionCenter", Value: "C", DisplayName: "Center", Aliases: []string{"C", "Center"}, Doc: "PositionCenter represents a center forward."},
			{Name: "PositionLeftWing", Value: "LW", DisplayName: "Left Wing", Aliases: []string{"L", "LW", "Left Wing", "LeftWing"}, Doc: "PositionLeftWing represents a left wing forward."},
			{Name: "PositionRightWing", Value: "RW", DisplayName: "Right Wing", Aliases: []string{"R", "RW", "Right Wing", "RightWing"}, Doc: "PositionRightWing represents a right wing forward."},
			{Name: "PositionForward", Value: "F", DisplayName: "Forward", Aliases: []string{"F", "Forward"}, Doc: "PositionForward represents a generic forward (used in historical data)."},
			{Name: "PositionDefense", Value: "D", DisplayName: "Defense", Aliases: []string{"D", "Defense", "Defenseman"}, Doc: "PositionDefense represents a defenseman."},
			{Name: "PositionGoalie", Value: "G", DisplayName: "Goalie", Aliases: []string{"G", "Goalie", "Goaltender"}, Doc: "PositionGoalie represents a goaltender."},
		},
	},
	{
		TypeName:              "Handedness",
		Doc:                   "Handedness represents a player's shooting or catching hand.",
		ErrorLabel:            "handedness",
		HasCode:               true,
		HasName:               true,
		AllowEmpty:            true,
		SkipMarshalValidation: true,
		Values: []ValueDef{
			{Name: "HandednessLeft", Value: "L", DisplayName: "Left", Aliases: []string{"L", "Left"}, Doc: "HandednessLeft represents left-handed."},
			{Name: "HandednessRight", Value: "R", DisplayName: "Right", Aliases: []string{"R", "Right"}, Doc: "HandednessRight represents right-handed."},
		},
	},
	{
		TypeName:       "GoalieDecision",
		Doc:            "GoalieDecision represents the decision (result) for a goalie in a game.",
		ErrorLabel:     "goalie decision",
		HasDisplayName: true,
		Values: []ValueDef{
			{Name: "GoalieDecisionWin", Value: "W", DisplayName: "Win", Aliases: []string{"W", "Win"}, Doc: "GoalieDecisionWin represents a win."},
			{Name: "GoalieDecisionLoss", Value: "L", DisplayName: "Loss", Aliases: []string{"L", "Loss"}, Doc: "GoalieDecisionLoss represents a loss."},
			{Name: "GoalieDecisionTie", Value: "T", DisplayName: "Tie", Aliases: []string{"T", "Tie"}, Doc: "GoalieDecisionTie represents a tie (used in older games)."},
			{Name: "GoalieDecisionOvertimeLoss", Value: "OTL", DisplayName: "Overtime Loss", Aliases: []string{"O", "OTL", "Overtime Loss", "OvertimeLoss"}, Doc: "GoalieDecisionOvertimeLoss represents an overtime/shootout loss."},
		},
	},
	{
		TypeName:   "PeriodType",
		Doc:        "PeriodType represents the type of period in a hockey game.",
		ErrorLabel: "period type",
		HasCode:    true,
		HasName:    true,
		Values: []ValueDef{
			{Name: "PeriodTypeRegulation", Value: "REG", DisplayName: "Regulation", Aliases: []string{"REG", "Regulation"}, Doc: "PeriodTypeRegulation represents a regulation period."},
			{Name: "PeriodTypeOvertime", Value: "OT", DisplayName: "Overtime", Aliases: []string{"OT", "Overtime"}, Doc: "PeriodTypeOvertime represents an overtime period."},
			{Name: "PeriodTypeShootout", Value: "SO", DisplayName: "Shootout", Aliases: []string{"SO", "Shootout"}, Doc: "PeriodTypeShootout represents a shootout."},
		},
	},
	{
		TypeName:   "HomeRoad",
		Doc:        "HomeRoad represents whether a team is home or away (road).",
		ErrorLabel: "home/road",
		HasCode:    true,
		HasName:    true,
		Values: []ValueDef{
			{Name: "HomeRoadHome", Value: "H", DisplayName: "Home", Aliases: []string{"H", "Home"}, Doc: "HomeRoadHome represents the home team."},
			{Name: "HomeRoadRoad", Value: "R", DisplayName: "Road", Aliases: []string{"R", "Road", "Away"}, Doc: "HomeRoadRoad represents the away/road team."},
		},
	},
	{
		TypeName:       "ZoneCode",
		Doc:            "ZoneCode represents a zone on the ice.",
		ErrorLabel:     "zone code",
		HasCode:        true,
		HasDisplayName: true,
		Values: []ValueDef{
			{Name: "ZoneCodeOffensive", Value: "O", DisplayName: "Offensive", Aliases: []string{"O", "Offensive"}, Doc: "ZoneCodeOffensive represents the offensive zone."},
			{Name: "ZoneCodeDefensive", Value: "D", DisplayName: "Defensive", Aliases: []string{"D", "Defensive"}, Doc: "ZoneCodeDefensive represents the defensive zone."},
			{Name: "ZoneCodeNeutral", Value: "N", DisplayName: "Neutral", Aliases: []string{"N", "Neutral"}, Doc: "ZoneCodeNeutral represents the neutral zone."},
		},
	},
	{
		TypeName:   "DefendingSide",
		Doc:        "DefendingSide represents which side of the ice a team is defending.",
		ErrorLabel: "defending side",
		AllowEmpty: true,
		Values: []ValueDef{
			{Name: "DefendingSideLeft", Value: "left", Aliases: []string{"left"}, Doc: "DefendingSideLeft represents defending the left side."},
			{Name: "DefendingSideRight", Value: "right", Aliases: []string{"right"}, Doc: "DefendingSideRight represents defending the right side."},
		},
	},
	{
		TypeName:   "GameScheduleState",
		Doc:        "GameScheduleState represents the state of a scheduled game.",
		ErrorLabel: "game schedule state",
		Values: []ValueDef{
			{Name: "GameScheduleStateOK", Value: "OK", Doc: "GameScheduleStateOK represents a normal scheduled game."},
			{Name: "GameScheduleStateDontPlay", Value: "DONT_PLAY", Doc: "GameScheduleStateDontPlay represents a game that won't be played."},
			{Name: "GameScheduleStatePostponed", Value: "PPD", Doc: "GameScheduleStatePostponed represents a postponed game."},
			{Name: "GameScheduleStateSuspended", Value: "SUSP", Doc: "GameScheduleStateSuspended represents a suspended game."},
			{Name: "GameScheduleStateTBD", Value: "TBD", Doc: "GameScheduleStateTBD represents a game with time to be determined."},
			{Name: "GameScheduleStateCompleted", Value: "COMPLETED", Doc: "GameScheduleStateCompleted represents a completed game."},
			{Name: "GameScheduleStateCancelled", Value: "CNCL", Doc: "GameScheduleStateCancelled represents a cancelled game."},
		},
	},
	{
		TypeName:   "PlayEventType",
		Doc:        "PlayEventType represents a type of play event during a game.",
		ErrorLabel: "play event type",
		Values: []ValueDef{
			{Name: "PlayEventTypeGameStart", Value: "game-start", Doc: "PlayEventTypeGameStart represents the start of a game."},
			{Name: "PlayEventTypePeriodStart", Value: "period-start", Doc: "PlayEventTypePeriodStart represents the start of a period."},
			{Name: "PlayEventTypePeriodEnd", Value: "period-end", Doc: "PlayEventTypePeriodEnd represents the end of a period."},
			{Name: "PlayEventTypeGameEnd", Value: "game-end", Doc: "PlayEventTypeGameEnd represents the end of a game."},
			{Name: "PlayEventTypeFaceoff", Value: "faceoff", Doc: "PlayEventTypeFaceoff represents a faceoff."},
			{Name: "PlayEventTypeHit", Value: "hit", Doc: "PlayEventTypeHit represents a hit."},
			{Name: "PlayEventTypeGiveaway", Value: "giveaway", Doc: "PlayEventTypeGiveaway represents a giveaway."},
			{Name: "PlayEventTypeTakeaway", Value: "takeaway", Doc: "PlayEventTypeTakeaway represents a takeaway."},
			{Name: "PlayEventTypeShotOnGoal", Value: "shot-on-goal", Doc: "PlayEventTypeShotOnGoal represents a shot on goal."},
			{Name: "PlayEventTypeMissedShot", Value: "missed-shot", Doc: "PlayEventTypeMissedShot represents a missed shot."},
			{Name: "PlayEventTypeBlockedShot", Value: "blocked-shot", Doc: "PlayEventTypeBlockedShot represents a blocked shot."},
			{Name: "PlayEventTypeGoal", Value: "goal", Doc: "PlayEventTypeGoal represents a goal."},
			{Name: "PlayEventTypePenalty", Value: "penalty", Doc: "PlayEventTypePenalty represents a penalty."},
			{Name: "PlayEventTypeStoppage", Value: "stoppage", Doc: "PlayEventTypeStoppage represents a stoppage in play."},
			{Name: "PlayEventTypeDelayedPenalty", Value: "delayed-penalty", Doc: "PlayEventTypeDelayedPenalty represents a delayed penalty."},
			{Name: "PlayEventTypeFailedShotAttempt", Value: "failed-shot-attempt", Doc: "PlayEventTypeFailedShotAttempt represents a failed shot attempt."},
			{Name: "PlayEventTypeShootoutComplete", Value: "shootout-complete", Doc: "PlayEventTypeShootoutComplete represents the completion of a shootout."},
			{Name: "PlayEventTypeUnknown", Value: "unknown", Doc: "PlayEventTypeUnknown represents an unknown event type."},
		},
	},
	{
		TypeName:   "GameState",
		Doc:        "GameState represents the current state of an NHL game.",
		ErrorLabel: "game state",
		Values: []ValueDef{
			{Name: "GameStateFuture", Value: "FUT", Doc: "GameStateFuture represents a game that has not yet started."},
			{Name: "GameStatePreGame", Value: "PRE", Doc: "GameStatePreGame represents a game in pre-game state."},
			{Name: "GameStateLive", Value: "LIVE", Doc: "GameStateLive represents a game currently in progress."},
			{Name: "GameStateFinal", Value: "FINAL", Doc: "GameStateFinal represents a completed game."},
			{Name: "GameStateOff", Value: "OFF", Doc: "GameStateOff represents a game that is off/completed (alternative to FINAL)."},
			{Name: "GameStatePostponed", Value: "PPD", Doc: "GameStatePostponed represents a postponed game."},
			{Name: "GameStateSuspended", Value: "SUSP", Doc: "GameStateSuspended represents a suspended game."},
			{Name: "GameStateCritical", Value: "CRIT", Doc: "GameStateCritical represents a game in critical state."},
		},
	},
}
