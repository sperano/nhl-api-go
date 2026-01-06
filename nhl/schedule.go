package nhl

import "fmt"

// ScheduleGame represents a game in the NHL schedule with comprehensive game information.
type ScheduleGame struct {
	ID           int64        `json:"id"`
	GameType     GameType     `json:"gameType"`
	GameDate     *string      `json:"gameDate,omitempty"`
	StartTimeUTC string       `json:"startTimeUTC"`
	AwayTeam     ScheduleTeam `json:"awayTeam"`
	HomeTeam     ScheduleTeam `json:"homeTeam"`
	GameState    GameState    `json:"gameState"`
}

// String implements fmt.Stringer for ScheduleGame.
// Returns a formatted string like "BUF @ TOR on 2023-10-10 [FUT]" or "BUF @ TOR [FUT]" if no date.
func (s ScheduleGame) String() string {
	if s.GameDate != nil {
		return fmt.Sprintf("%s @ %s on %s [%s]", s.AwayTeam.Abbrev, s.HomeTeam.Abbrev, *s.GameDate, s.GameState)
	}
	return fmt.Sprintf("%s @ %s [%s]", s.AwayTeam.Abbrev, s.HomeTeam.Abbrev, s.GameState)
}

// ScheduleTeam represents team information within a schedule context.
// Contains basic team identification and optional score information.
type ScheduleTeam struct {
	ID        int64            `json:"id"`
	Abbrev    string           `json:"abbrev"`
	PlaceName *LocalizedString `json:"placeName,omitempty"`
	Logo      string           `json:"logo"`
	Score     *int             `json:"score,omitempty"`
}

// DailySchedule represents the schedule for a single day.
// Contains navigation dates for previous and next days with games.
type DailySchedule struct {
	NextStartDate     *string        `json:"nextStartDate,omitempty"`
	PreviousStartDate *string        `json:"previousStartDate,omitempty"`
	Date              string         `json:"date"`
	Games             []ScheduleGame `json:"games"`
	NumberOfGames     int            `json:"numberOfGames"`
}

// WeeklyScheduleResponse represents a week's worth of games organized by day.
// Used for retrieving a week-long schedule from the API.
type WeeklyScheduleResponse struct {
	NextStartDate     string    `json:"nextStartDate"`
	PreviousStartDate string    `json:"previousStartDate"`
	GameWeek          []GameDay `json:"gameWeek"`
}

// GameDay represents all games scheduled for a specific day.
type GameDay struct {
	Date  string         `json:"date"`
	Games []ScheduleGame `json:"games"`
}

// TeamScheduleResponse represents a team-specific schedule response.
// Used for monthly or weekly team schedules.
type TeamScheduleResponse struct {
	Games []ScheduleGame `json:"games"`
}

// DailyScores represents game scores for a specific day.
// Includes navigation to previous and next days.
type DailyScores struct {
	PrevDate    string      `json:"prevDate"`
	CurrentDate string      `json:"currentDate"`
	NextDate    string      `json:"nextDate"`
	Games       []GameScore `json:"games"`
}

// GameScore represents a single game's score information.
// Similar to ScheduleGame but focused on score display.
type GameScore struct {
	ID        int64        `json:"id"`
	GameType  GameType     `json:"gameType"`
	GameState GameState    `json:"gameState"`
	AwayTeam  ScheduleTeam `json:"awayTeam"`
	HomeTeam  ScheduleTeam `json:"homeTeam"`
}

// String implements fmt.Stringer for GameScore.
// Returns a formatted string like "BUF 3 @ TOR 2 [FINAL]" or "BUF - @ TOR - [FUT]" for no scores.
func (g GameScore) String() string {
	formatScore := func(score *int) string {
		if score == nil {
			return "-"
		}
		return fmt.Sprintf("%d", *score)
	}

	return fmt.Sprintf("%s %s @ %s %s [%s]",
		g.AwayTeam.Abbrev,
		formatScore(g.AwayTeam.Score),
		g.HomeTeam.Abbrev,
		formatScore(g.HomeTeam.Score),
		g.GameState,
	)
}
