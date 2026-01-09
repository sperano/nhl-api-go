package nhl

// Boxscore represents the boxscore response with detailed game and player statistics.
type Boxscore struct {
	ID                GameID            `json:"id"`
	Season            Season            `json:"season"`
	GameType          GameType          `json:"gameType"`
	LimitedScoring    bool              `json:"limitedScoring"`
	GameDate          string            `json:"gameDate"`
	Venue             LocalizedString   `json:"venue"`
	VenueLocation     LocalizedString   `json:"venueLocation"`
	StartTimeUTC      string            `json:"startTimeUTC"`
	EasternUTCOffset  string            `json:"easternUTCOffset"`
	VenueUTCOffset    string            `json:"venueUTCOffset"`
	TVBroadcasts      []TVBroadcast     `json:"tvBroadcasts"`
	GameState         GameState         `json:"gameState"`
	GameScheduleState GameScheduleState `json:"gameScheduleState"`
	PeriodDescriptor  PeriodDescriptor  `json:"periodDescriptor"`
	SpecialEvent      *SpecialEvent     `json:"specialEvent,omitempty"`
	AwayTeam          BoxscoreTeam      `json:"awayTeam"`
	HomeTeam          BoxscoreTeam      `json:"homeTeam"`
	Clock             GameClock         `json:"clock"`
	PlayerByGameStats PlayerByGameStats `json:"playerByGameStats"`
}

// TVBroadcast represents TV broadcast information for a game.
type TVBroadcast struct {
	ID             int64  `json:"id"`
	Market         string `json:"market"`
	CountryCode    string `json:"countryCode"`
	Network        string `json:"network"`
	SequenceNumber int    `json:"sequenceNumber"`
}

// SpecialEvent represents special event information for a game.
type SpecialEvent struct {
	ParentID     int64           `json:"parentId"`
	Name         LocalizedString `json:"name"`
	LightLogoURL LocalizedString `json:"lightLogoUrl"`
}

// PeriodDescriptor represents period descriptor with game period information.
type PeriodDescriptor struct {
	Number               int        `json:"number"`
	PeriodType           PeriodType `json:"periodType"`
	MaxRegulationPeriods int        `json:"maxRegulationPeriods"`
}

// BoxscoreTeam represents team information in the boxscore.
type BoxscoreTeam struct {
	ID                       TeamID          `json:"id"`
	CommonName               LocalizedString `json:"commonName"`
	Abbrev                   string          `json:"abbrev"`
	Score                    int             `json:"score"`
	SOG                      int             `json:"sog"`
	Logo                     string          `json:"logo"`
	DarkLogo                 string          `json:"darkLogo"`
	PlaceName                LocalizedString `json:"placeName"`
	PlaceNameWithPreposition LocalizedString `json:"placeNameWithPreposition"`
}

// GameClock represents game clock information.
type GameClock struct {
	TimeRemaining    string `json:"timeRemaining"`
	SecondsRemaining int    `json:"secondsRemaining"`
	Running          bool   `json:"running"`
	InIntermission   bool   `json:"inIntermission"`
}

// PlayerByGameStats represents player statistics organized by team.
type PlayerByGameStats struct {
	AwayTeam TeamPlayerStats `json:"awayTeam"`
	HomeTeam TeamPlayerStats `json:"homeTeam"`
}

// TeamPlayerStats represents a team's player statistics grouped by position.
type TeamPlayerStats struct {
	Forwards []SkaterStats `json:"forwards"`
	Defense  []SkaterStats `json:"defense"`
	Goalies  []GoalieStats `json:"goalies"`
}

// TeamGameStats represents aggregated team statistics for game comparison.
type TeamGameStats struct {
	ShotsOnGoal            int
	FaceoffWins            int
	FaceoffTotal           int
	PowerPlayGoals         int
	PowerPlayOpportunities int
	PenaltyMinutes         int
	Hits                   int
	BlockedShots           int
	Giveaways              int
	Takeaways              int
}

// FromTeamPlayerStats calculates aggregated team statistics from individual player stats.
func FromTeamPlayerStats(stats *TeamPlayerStats) TeamGameStats {
	teamStats := TeamGameStats{}

	aggregateSkaterStats(&teamStats, stats)
	aggregateGoalieStats(&teamStats, stats)

	return teamStats
}

// aggregateSkaterStats aggregates statistics from forwards and defensemen.
func aggregateSkaterStats(teamStats *TeamGameStats, stats *TeamPlayerStats) {
	allSkaters := make([]SkaterStats, 0, len(stats.Forwards)+len(stats.Defense))
	allSkaters = append(allSkaters, stats.Forwards...)
	allSkaters = append(allSkaters, stats.Defense...)

	for i := range allSkaters {
		skater := &allSkaters[i]
		teamStats.ShotsOnGoal += skater.SOG
		teamStats.PowerPlayGoals += skater.PowerPlayGoals
		teamStats.PenaltyMinutes += skater.PIM
		teamStats.Hits += skater.Hits
		teamStats.BlockedShots += skater.BlockedShots
		teamStats.Giveaways += skater.Giveaways
		teamStats.Takeaways += skater.Takeaways

		addFaceoffStats(teamStats, skater)
	}
}

// addFaceoffStats adds faceoff statistics from a skater to the team totals.
// Note: This logic currently only counts centers for faceoffs, which may not be entirely accurate
// as wings can also take faceoffs in certain situations.
func addFaceoffStats(teamStats *TeamGameStats, skater *SkaterStats) {
	if skater.Position == PositionCenter && skater.FaceoffWinningPctg > 0.0 {
		// Estimate total faceoffs using shifts as a proxy for faceoff participation
		estimatedFaceoffs := skater.Shifts
		teamStats.FaceoffTotal += estimatedFaceoffs
		teamStats.FaceoffWins += int(float64(estimatedFaceoffs)*skater.FaceoffWinningPctg + 0.5)
	}
}

// aggregateGoalieStats aggregates statistics from goalies.
func aggregateGoalieStats(teamStats *TeamGameStats, stats *TeamPlayerStats) {
	for i := range stats.Goalies {
		goalie := &stats.Goalies[i]
		if goalie.PIM != nil {
			teamStats.PenaltyMinutes += *goalie.PIM
		}
		// Count power play opportunities from goals against
		teamStats.PowerPlayOpportunities += goalie.PowerPlayGoalsAgainst
	}
}

// FaceoffPercentage calculates the faceoff winning percentage for the team.
// Returns 0.0 if the team has not taken any faceoffs.
func (t *TeamGameStats) FaceoffPercentage() float64 {
	if t.FaceoffTotal > 0 {
		return (float64(t.FaceoffWins) / float64(t.FaceoffTotal)) * 100.0
	}
	return 0.0
}

// PowerPlayPercentage calculates the power play percentage for the team.
// Returns 0.0 if the team has not had any power play opportunities.
func (t *TeamGameStats) PowerPlayPercentage() float64 {
	if t.PowerPlayOpportunities > 0 {
		return (float64(t.PowerPlayGoals) / float64(t.PowerPlayOpportunities)) * 100.0
	}
	return 0.0
}

// SkaterStats represents skater (forward/defense) statistics.
type SkaterStats struct {
	PlayerID           PlayerID        `json:"playerId"`
	SweaterNumber      int             `json:"sweaterNumber"`
	Name               LocalizedString `json:"name"`
	Position           Position        `json:"position"`
	Goals              int             `json:"goals"`
	Assists            int             `json:"assists"`
	Points             int             `json:"points"`
	PlusMinus          int             `json:"plusMinus"`
	PIM                int             `json:"pim"`
	Hits               int             `json:"hits"`
	PowerPlayGoals     int             `json:"powerPlayGoals"`
	SOG                int             `json:"sog"`
	FaceoffWinningPctg float64         `json:"faceoffWinningPctg"`
	TOI                string          `json:"toi"`
	BlockedShots       int             `json:"blockedShots"`
	Shifts             int             `json:"shifts"`
	Giveaways          int             `json:"giveaways"`
	Takeaways          int             `json:"takeaways"`
}

// GoalieStats represents goalie statistics.
type GoalieStats struct {
	PlayerID                 PlayerID        `json:"playerId"`
	SweaterNumber            int             `json:"sweaterNumber"`
	Name                     LocalizedString `json:"name"`
	Position                 Position        `json:"position"`
	EvenStrengthShotsAgainst string          `json:"evenStrengthShotsAgainst"`
	PowerPlayShotsAgainst    string          `json:"powerPlayShotsAgainst"`
	ShorthandedShotsAgainst  string          `json:"shorthandedShotsAgainst"`
	SaveShotsAgainst         string          `json:"saveShotsAgainst"`
	SavePctg                 *float64        `json:"savePctg,omitempty"`
	EvenStrengthGoalsAgainst int             `json:"evenStrengthGoalsAgainst"`
	PowerPlayGoalsAgainst    int             `json:"powerPlayGoalsAgainst"`
	ShorthandedGoalsAgainst  int             `json:"shorthandedGoalsAgainst"`
	PIM                      *int            `json:"pim,omitempty"`
	GoalsAgainst             int             `json:"goalsAgainst"`
	TOI                      string          `json:"toi"`
	Starter                  *bool           `json:"starter,omitempty"`
	Decision                 *GoalieDecision `json:"decision,omitempty"`
	ShotsAgainst             int             `json:"shotsAgainst"`
	Saves                    int             `json:"saves"`
}
