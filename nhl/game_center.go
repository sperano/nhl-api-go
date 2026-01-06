package nhl

import (
	"fmt"
	"strconv"
)

// GameSituation represents a parsed game situation from situation code.
// The situation code is a 4-character string (e.g., "1551") that encodes:
// - Position 0: Away goalie in net (1) or not (0)
// - Position 1: Number of away skaters (1-6)
// - Position 2: Number of home skaters (1-6)
// - Position 3: Home goalie in net (1) or not (0)
type GameSituation struct {
	// AwaySkaters is the number of away team skaters on ice (1-6).
	AwaySkaters int
	// AwayGoalieIn indicates whether the away team has a goalie in net.
	AwayGoalieIn bool
	// HomeSkaters is the number of home team skaters on ice (1-6).
	HomeSkaters int
	// HomeGoalieIn indicates whether the home team has a goalie in net.
	HomeGoalieIn bool
}

// GameSituationFromCode parses a situation code string (e.g., "1551").
// Returns nil if the code is invalid (wrong length or contains non-digits).
func GameSituationFromCode(code string) *GameSituation {
	if len(code) != 4 {
		return nil
	}

	awayGoalieDigit := code[0]
	awaySkaterDigit := code[1]
	homeSkaterDigit := code[2]
	homeGoalieDigit := code[3]

	awaySkaters, err := strconv.Atoi(string(awaySkaterDigit))
	if err != nil {
		return nil
	}

	homeSkaters, err := strconv.Atoi(string(homeSkaterDigit))
	if err != nil {
		return nil
	}

	return &GameSituation{
		AwaySkaters:  awaySkaters,
		AwayGoalieIn: awayGoalieDigit == '1',
		HomeSkaters:  homeSkaters,
		HomeGoalieIn: homeGoalieDigit == '1',
	}
}

// IsEvenStrength returns true if this is even strength (5v5, 4v4, or 3v3).
func (g *GameSituation) IsEvenStrength() bool {
	return g.AwaySkaters == g.HomeSkaters
}

// IsAwayPowerPlay returns true if the away team has a power play advantage.
func (g *GameSituation) IsAwayPowerPlay() bool {
	return g.AwaySkaters > g.HomeSkaters
}

// IsHomePowerPlay returns true if the home team has a power play advantage.
func (g *GameSituation) IsHomePowerPlay() bool {
	return g.HomeSkaters > g.AwaySkaters
}

// IsEmptyNet returns true if either team has pulled their goalie.
func (g *GameSituation) IsEmptyNet() bool {
	return !g.AwayGoalieIn || !g.HomeGoalieIn
}

// StrengthDescription returns the strength description (e.g., "5v5", "5v4 PP", "6v5 EN").
func (g *GameSituation) StrengthDescription() string {
	base := fmt.Sprintf("%dv%d", g.AwaySkaters, g.HomeSkaters)

	if !g.AwayGoalieIn || !g.HomeGoalieIn {
		return base + " EN"
	}

	if g.AwaySkaters != g.HomeSkaters {
		return base + " PP"
	}

	return base
}

// String returns the strength description.
func (g *GameSituation) String() string {
	return g.StrengthDescription()
}

// PlayByPlay represents the play-by-play response with all game events.
type PlayByPlay struct {
	ID                  int64               `json:"id"`
	Season              int64               `json:"season"`
	GameType            GameType            `json:"gameType"`
	LimitedScoring      bool                `json:"limitedScoring"`
	GameDate            string              `json:"gameDate"`
	Venue               LocalizedString     `json:"venue"`
	VenueLocation       LocalizedString     `json:"venueLocation"`
	StartTimeUTC        string              `json:"startTimeUTC"`
	EasternUTCOffset    string              `json:"easternUTCOffset"`
	VenueUTCOffset      string              `json:"venueUTCOffset"`
	TVBroadcasts        []TVBroadcast       `json:"tvBroadcasts"`
	GameState           GameState           `json:"gameState"`
	GameScheduleState   GameScheduleState   `json:"gameScheduleState"`
	PeriodDescriptor    PeriodDescriptor    `json:"periodDescriptor"`
	SpecialEvent        *SpecialEvent       `json:"specialEvent,omitempty"`
	AwayTeam            BoxscoreTeam        `json:"awayTeam"`
	HomeTeam            BoxscoreTeam        `json:"homeTeam"`
	ShootoutInUse       bool                `json:"shootoutInUse"`
	OTInUse             bool                `json:"otInUse"`
	Clock               GameClock           `json:"clock"`
	DisplayPeriod       int                 `json:"displayPeriod"`
	MaxPeriods          int                 `json:"maxPeriods"`
	GameOutcome         *GameOutcome        `json:"gameOutcome,omitempty"`
	Plays               []PlayEvent         `json:"plays"`
	RosterSpots         []RosterSpot        `json:"rosterSpots"`
	RegPeriods          *int                `json:"regPeriods,omitempty"`
	Summary             *GameSummary        `json:"summary,omitempty"`
}

// RecentPlays returns the most recent N plays (most recent first).
func (p *PlayByPlay) RecentPlays(count int) []*PlayEvent {
	if count <= 0 {
		return nil
	}

	playCount := len(p.Plays)
	if playCount == 0 {
		return nil
	}

	if count > playCount {
		count = playCount
	}

	recent := make([]*PlayEvent, count)
	for i := 0; i < count; i++ {
		recent[i] = &p.Plays[playCount-1-i]
	}

	return recent
}

// Goals returns all goal events in the game.
func (p *PlayByPlay) Goals() []*PlayEvent {
	goals := make([]*PlayEvent, 0)
	for i := range p.Plays {
		if p.Plays[i].TypeDescKey == PlayEventTypeGoal {
			goals = append(goals, &p.Plays[i])
		}
	}
	return goals
}

// Penalties returns all penalty events in the game.
func (p *PlayByPlay) Penalties() []*PlayEvent {
	penalties := make([]*PlayEvent, 0)
	for i := range p.Plays {
		if p.Plays[i].TypeDescKey == PlayEventTypePenalty {
			penalties = append(penalties, &p.Plays[i])
		}
	}
	return penalties
}

// Shots returns all shot events (on goal, missed, blocked) in the game.
func (p *PlayByPlay) Shots() []*PlayEvent {
	shots := make([]*PlayEvent, 0)
	for i := range p.Plays {
		if p.Plays[i].TypeDescKey.IsScoringChance() {
			shots = append(shots, &p.Plays[i])
		}
	}
	return shots
}

// PlaysInPeriod returns all plays for a specific period number.
func (p *PlayByPlay) PlaysInPeriod(period int) []*PlayEvent {
	plays := make([]*PlayEvent, 0)
	for i := range p.Plays {
		if p.Plays[i].PeriodDescriptor.Number == period {
			plays = append(plays, &p.Plays[i])
		}
	}
	return plays
}

// GetPlayer returns a player from the roster by player ID.
// Returns nil if the player is not found.
func (p *PlayByPlay) GetPlayer(playerID int64) *RosterSpot {
	for i := range p.RosterSpots {
		if p.RosterSpots[i].PlayerID == playerID {
			return &p.RosterSpots[i]
		}
	}
	return nil
}

// TeamRoster returns all players for a specific team ID.
func (p *PlayByPlay) TeamRoster(teamID int64) []*RosterSpot {
	roster := make([]*RosterSpot, 0)
	for i := range p.RosterSpots {
		if p.RosterSpots[i].TeamID == teamID {
			roster = append(roster, &p.RosterSpots[i])
		}
	}
	return roster
}

// CurrentSituation returns the current game situation from the latest play.
// Returns nil if there are no plays or the situation code is invalid.
func (p *PlayByPlay) CurrentSituation() *GameSituation {
	if len(p.Plays) == 0 {
		return nil
	}
	return p.Plays[len(p.Plays)-1].Situation()
}

// GameOutcome represents game outcome information.
type GameOutcome struct {
	LastPeriodType PeriodType `json:"lastPeriodType"`
}

// PlayEvent represents an individual play event in the game.
type PlayEvent struct {
	EventID               int64              `json:"eventId"`
	PeriodDescriptor      PeriodDescriptor   `json:"periodDescriptor"`
	TimeInPeriod          string             `json:"timeInPeriod"`
	TimeRemaining         string             `json:"timeRemaining"`
	SituationCode         string             `json:"situationCode"`
	HomeTeamDefendingSide DefendingSide      `json:"homeTeamDefendingSide"`
	TypeCode              int                `json:"typeCode"`
	TypeDescKey           PlayEventType      `json:"typeDescKey"`
	SortOrder             int                `json:"sortOrder"`
	Details               *PlayEventDetails  `json:"details,omitempty"`
	PPTReplayURL          *string            `json:"pptReplayUrl,omitempty"`
}

// Situation parses the situation code into a GameSituation.
// Returns nil if the situation code is invalid.
func (p *PlayEvent) Situation() *GameSituation {
	return GameSituationFromCode(p.SituationCode)
}

// PlayEventDetails represents details for a play event (varies by event type).
type PlayEventDetails struct {
	// Coordinate details
	XCoord           *int       `json:"xCoord,omitempty"`
	YCoord           *int       `json:"yCoord,omitempty"`
	ZoneCode         *ZoneCode  `json:"zoneCode,omitempty"`
	EventOwnerTeamID *int64     `json:"eventOwnerTeamId,omitempty"`

	// Shot details
	ShotType         *string `json:"shotType,omitempty"`
	ShootingPlayerID *int64  `json:"shootingPlayerId,omitempty"`
	GoalieInNetID    *int64  `json:"goalieInNetId,omitempty"`

	// Blocked shot details
	BlockingPlayerID *int64 `json:"blockingPlayerId,omitempty"`

	// Goal details
	ScoringPlayerID        *int64  `json:"scoringPlayerId,omitempty"`
	ScoringPlayerTotal     *int    `json:"scoringPlayerTotal,omitempty"`
	Assist1PlayerID        *int64  `json:"assist1PlayerId,omitempty"`
	Assist1PlayerTotal     *int    `json:"assist1PlayerTotal,omitempty"`
	Assist2PlayerID        *int64  `json:"assist2PlayerId,omitempty"`
	Assist2PlayerTotal     *int    `json:"assist2PlayerTotal,omitempty"`
	AwayScore              *int    `json:"awayScore,omitempty"`
	HomeScore              *int    `json:"homeScore,omitempty"`
	HighlightClip          *int64  `json:"highlightClip,omitempty"`
	HighlightClipSharingURL *string `json:"highlightClipSharingUrl,omitempty"`
	DiscreteClip           *int64  `json:"discreteClip,omitempty"`

	// Penalty details
	TypeCode             *string `json:"typeCode,omitempty"`
	DescKey              *string `json:"descKey,omitempty"`
	Duration             *int    `json:"duration,omitempty"`
	CommittedByPlayerID  *int64  `json:"committedByPlayerId,omitempty"`
	DrawnByPlayerID      *int64  `json:"drawnByPlayerId,omitempty"`

	// Hit details
	HittingPlayerID *int64 `json:"hittingPlayerId,omitempty"`
	HitteePlayerID  *int64 `json:"hitteePlayerId,omitempty"`

	// Faceoff details
	WinningPlayerID *int64 `json:"winningPlayerId,omitempty"`
	LosingPlayerID  *int64 `json:"losingPlayerId,omitempty"`

	// General details
	PlayerID *int64  `json:"playerId,omitempty"`
	Reason   *string `json:"reason,omitempty"`
	AwaySOG  *int    `json:"awaySOG,omitempty"`
	HomeSOG  *int    `json:"homeSOG,omitempty"`
}

// RosterSpot represents player information in the roster.
type RosterSpot struct {
	TeamID        int64           `json:"teamId"`
	PlayerID      int64           `json:"playerId"`
	FirstName     LocalizedString `json:"firstName"`
	LastName      LocalizedString `json:"lastName"`
	SweaterNumber int             `json:"sweaterNumber"`
	Position      Position        `json:"positionCode"`
	Headshot      string          `json:"headshot"`
}

// GameMatchup represents the game matchup/landing response.
type GameMatchup struct {
	ID                  int64               `json:"id"`
	Season              int64               `json:"season"`
	GameType            GameType            `json:"gameType"`
	LimitedScoring      bool                `json:"limitedScoring"`
	GameDate            string              `json:"gameDate"`
	Venue               LocalizedString     `json:"venue"`
	VenueLocation       LocalizedString     `json:"venueLocation"`
	StartTimeUTC        string              `json:"startTimeUTC"`
	EasternUTCOffset    string              `json:"easternUTCOffset"`
	VenueUTCOffset      string              `json:"venueUTCOffset"`
	VenueTimezone       string              `json:"venueTimezone"`
	PeriodDescriptor    PeriodDescriptor    `json:"periodDescriptor"`
	TVBroadcasts        []TVBroadcast       `json:"tvBroadcasts"`
	GameState           GameState           `json:"gameState"`
	GameScheduleState   GameScheduleState   `json:"gameScheduleState"`
	SpecialEvent        *SpecialEvent       `json:"specialEvent,omitempty"`
	AwayTeam            MatchupTeam         `json:"awayTeam"`
	HomeTeam            MatchupTeam         `json:"homeTeam"`
	ShootoutInUse       bool                `json:"shootoutInUse"`
	MaxPeriods          int                 `json:"maxPeriods"`
	RegPeriods          int                 `json:"regPeriods"`
	OTInUse             bool                `json:"otInUse"`
	TiesInUse           bool                `json:"tiesInUse"`
	Summary             *GameSummary        `json:"summary,omitempty"`
	Clock               *GameClock          `json:"clock,omitempty"`
}

// MatchupTeam represents team information in game matchup.
type MatchupTeam struct {
	ID                        int64           `json:"id"`
	CommonName                LocalizedString `json:"commonName"`
	Abbrev                    string          `json:"abbrev"`
	PlaceName                 LocalizedString `json:"placeName"`
	PlaceNameWithPreposition  LocalizedString `json:"placeNameWithPreposition"`
	Score                     int             `json:"score"`
	SOG                       int             `json:"sog"`
	Logo                      string          `json:"logo"`
	DarkLogo                  string          `json:"darkLogo"`
}

// GameSummary represents game summary with scoring and penalties.
type GameSummary struct {
	Scoring    []PeriodScoring     `json:"scoring"`
	Shootout   *[]ShootoutAttempt  `json:"shootout,omitempty"`
	ThreeStars *[]ThreeStar        `json:"threeStars,omitempty"`
	Penalties  []PeriodPenalties   `json:"penalties"`
}

// PeriodScoring represents scoring summary for a period.
type PeriodScoring struct {
	PeriodDescriptor PeriodDescriptor `json:"periodDescriptor"`
	Goals            []GoalSummary    `json:"goals"`
}

// GoalSummary represents goal summary information.
type GoalSummary struct {
	SituationCode           string           `json:"situationCode"`
	EventID                 int64            `json:"eventId"`
	Strength                string           `json:"strength"`
	PlayerID                int64            `json:"playerId"`
	FirstName               LocalizedString  `json:"firstName"`
	LastName                LocalizedString  `json:"lastName"`
	Name                    LocalizedString  `json:"name"`
	TeamAbbrev              LocalizedString  `json:"teamAbbrev"`
	Headshot                string           `json:"headshot"`
	HighlightClipSharingURL *string          `json:"highlightClipSharingUrl,omitempty"`
	HighlightClip           *int64           `json:"highlightClip,omitempty"`
	DiscreteClip            *int64           `json:"discreteClip,omitempty"`
	GoalsToDate             *int             `json:"goalsToDate,omitempty"`
	AwayScore               int              `json:"awayScore"`
	HomeScore               int              `json:"homeScore"`
	LeadingTeamAbbrev       *LocalizedString `json:"leadingTeamAbbrev,omitempty"`
	TimeInPeriod            string           `json:"timeInPeriod"`
	ShotType                string           `json:"shotType"`
	GoalModifier            string           `json:"goalModifier"`
	Assists                 []AssistSummary  `json:"assists"`
	HomeTeamDefendingSide   DefendingSide    `json:"homeTeamDefendingSide"`
	IsHome                  bool             `json:"isHome"`
}

// AssistSummary represents assist summary information.
type AssistSummary struct {
	PlayerID      int64           `json:"playerId"`
	FirstName     LocalizedString `json:"firstName"`
	LastName      LocalizedString `json:"lastName"`
	Name          LocalizedString `json:"name"`
	AssistsToDate int             `json:"assistsToDate"`
	SweaterNumber int             `json:"sweaterNumber"`
}

// ShootoutAttempt represents shootout attempt information.
type ShootoutAttempt struct {
	Sequence   int             `json:"sequence"`
	PlayerID   int64           `json:"playerId"`
	TeamAbbrev LocalizedString `json:"teamAbbrev"`
	FirstName  LocalizedString `json:"firstName"`
	LastName   LocalizedString `json:"lastName"`
	ShotType   string          `json:"shotType"`
	Result     string          `json:"result"`
	Headshot   string          `json:"headshot"`
	GameWinner bool            `json:"gameWinner"`
}

// ThreeStar represents three stars selection.
type ThreeStar struct {
	Star                 int      `json:"star"`
	PlayerID             int64    `json:"playerId"`
	TeamAbbrev           string   `json:"teamAbbrev"`
	Headshot             string   `json:"headshot"`
	Name                 LocalizedString `json:"name"`
	SweaterNo            int      `json:"sweaterNo"`
	Position             Position `json:"position"`
	// Skater stats
	Goals                *int     `json:"goals,omitempty"`
	Assists              *int     `json:"assists,omitempty"`
	Points               *int     `json:"points,omitempty"`
	// Goalie stats
	GoalsAgainstAverage  *float64 `json:"goalsAgainstAverage,omitempty"`
	SavePctg             *float64 `json:"savePctg,omitempty"`
}

// PeriodPenalties represents penalty summary for a period.
type PeriodPenalties struct {
	PeriodDescriptor PeriodDescriptor  `json:"periodDescriptor"`
	Penalties        []PenaltySummary  `json:"penalties"`
}

// PenaltySummary represents penalty summary information.
type PenaltySummary struct {
	TimeInPeriod      string           `json:"timeInPeriod"`
	PenaltyType       string           `json:"type"`
	Duration          int              `json:"duration"`
	CommittedByPlayer *PenaltyPlayer   `json:"committedByPlayer,omitempty"`
	TeamAbbrev        LocalizedString  `json:"teamAbbrev"`
	DrawnBy           *PenaltyPlayer   `json:"drawnBy,omitempty"`
	DescKey           string           `json:"descKey"`
	ServedBy          *LocalizedString `json:"servedBy,omitempty"`
	EventID           *int64           `json:"eventId,omitempty"`
}

// PenaltyPlayer represents player information in penalty summary.
type PenaltyPlayer struct {
	FirstName     LocalizedString `json:"firstName"`
	LastName      LocalizedString `json:"lastName"`
	SweaterNumber int             `json:"sweaterNumber"`
}

// ShiftChart represents shift chart data.
type ShiftChart struct {
	Data []ShiftEntry `json:"data"`
}

// ShiftEntry represents an individual shift entry for a player.
type ShiftEntry struct {
	ID               int64   `json:"id"`
	DetailCode       int     `json:"detailCode"`
	Duration         string  `json:"duration"`
	EndTime          string  `json:"endTime"`
	EventDescription *string `json:"eventDescription,omitempty"`
	EventNumber      int64   `json:"eventNumber"`
	FirstName        string  `json:"firstName"`
	GameID           int64   `json:"gameId"`
	HexValue         string  `json:"hexValue"`
	LastName         string  `json:"lastName"`
	Period           int     `json:"period"`
	PlayerID         int64   `json:"playerId"`
	ShiftNumber      int     `json:"shiftNumber"`
	StartTime        string  `json:"startTime"`
	TeamAbbrev       string  `json:"teamAbbrev"`
	TeamID           int64   `json:"teamId"`
	TeamName         string  `json:"teamName"`
	TypeCode         int     `json:"typeCode"`
}

// SeasonSeriesMatchup represents season series matchup.
type SeasonSeriesMatchup struct {
	SeasonSeries     []SeriesGame    `json:"seasonSeries"`
	SeasonSeriesWins SeriesWins      `json:"seasonSeriesWins"`
	GameInfo         SeriesGameInfo  `json:"gameInfo"`
}

// SeriesGame represents an individual game in the season series.
type SeriesGame struct {
	ID                  int64               `json:"id"`
	Season              int64               `json:"season"`
	GameType            GameType            `json:"gameType"`
	GameDate            string              `json:"gameDate"`
	StartTimeUTC        string              `json:"startTimeUTC"`
	EasternUTCOffset    string              `json:"easternUTCOffset"`
	VenueUTCOffset      string              `json:"venueUTCOffset"`
	GameState           GameState           `json:"gameState"`
	GameScheduleState   GameScheduleState   `json:"gameScheduleState"`
	AwayTeam            SeriesTeam          `json:"awayTeam"`
	HomeTeam            SeriesTeam          `json:"homeTeam"`
	PeriodDescriptor    PeriodDescriptor    `json:"periodDescriptor"`
	GameCenterLink      string              `json:"gameCenterLink"`
	GameOutcome         GameOutcome         `json:"gameOutcome"`
}

// SeriesTeam represents team information in season series.
type SeriesTeam struct {
	ID     int64  `json:"id"`
	Abbrev string `json:"abbrev"`
	Logo   string `json:"logo"`
	Score  int    `json:"score"`
}

// SeriesWins represents season series win counts.
type SeriesWins struct {
	AwayTeamWins int `json:"awayTeamWins"`
	HomeTeamWins int `json:"homeTeamWins"`
}

// SeriesGameInfo represents game information including officials and scratches.
type SeriesGameInfo struct {
	Referees []LocalizedString `json:"referees"`
	Linesmen []LocalizedString `json:"linesmen"`
	AwayTeam TeamGameInfo      `json:"awayTeam"`
	HomeTeam TeamGameInfo      `json:"homeTeam"`
}

// TeamGameInfo represents team-specific game information.
type TeamGameInfo struct {
	HeadCoach LocalizedString   `json:"headCoach"`
	Scratches []ScratchedPlayer `json:"scratches"`
}

// ScratchedPlayer represents scratched player information.
type ScratchedPlayer struct {
	ID        int64           `json:"id"`
	FirstName LocalizedString `json:"firstName"`
	LastName  LocalizedString `json:"lastName"`
}

// GameStory represents game story.
type GameStory struct {
	ID                  int64               `json:"id"`
	Season              int64               `json:"season"`
	GameType            GameType            `json:"gameType"`
	LimitedScoring      bool                `json:"limitedScoring"`
	GameDate            string              `json:"gameDate"`
	Venue               LocalizedString     `json:"venue"`
	VenueLocation       LocalizedString     `json:"venueLocation"`
	StartTimeUTC        string              `json:"startTimeUTC"`
	EasternUTCOffset    string              `json:"easternUTCOffset"`
	VenueUTCOffset      string              `json:"venueUTCOffset"`
	VenueTimezone       string              `json:"venueTimezone"`
	TVBroadcasts        []TVBroadcast       `json:"tvBroadcasts"`
	GameState           GameState           `json:"gameState"`
	GameScheduleState   GameScheduleState   `json:"gameScheduleState"`
	AwayTeam            StoryTeam           `json:"awayTeam"`
	HomeTeam            StoryTeam           `json:"homeTeam"`
	ShootoutInUse       bool                `json:"shootoutInUse"`
	MaxPeriods          int                 `json:"maxPeriods"`
	RegPeriods          int                 `json:"regPeriods"`
	OTInUse             bool                `json:"otInUse"`
	TiesInUse           bool                `json:"tiesInUse"`
	Summary             *GameSummary        `json:"summary,omitempty"`
}

// StoryTeam represents team information in game story.
type StoryTeam struct {
	ID        int64           `json:"id"`
	Name      LocalizedString `json:"name"`
	Abbrev    string          `json:"abbrev"`
	PlaceName LocalizedString `json:"placeName"`
	Score     int             `json:"score"`
	SOG       int             `json:"sog"`
	Logo      string          `json:"logo"`
}
