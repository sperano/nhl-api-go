package nhl

// PlayerLanding represents comprehensive player profile data from the NHL API.
type PlayerLanding struct {
	PlayerID           PlayerID         `json:"playerId"`
	IsActive           bool             `json:"isActive"`
	CurrentTeamID      *TeamID          `json:"currentTeamId,omitempty"`
	CurrentTeamAbbrev  *string          `json:"currentTeamAbbrev,omitempty"`
	FirstName          LocalizedString  `json:"firstName"`
	LastName           LocalizedString  `json:"lastName"`
	SweaterNumber      *int             `json:"sweaterNumber,omitempty"`
	Position           Position         `json:"position"`
	Headshot           string           `json:"headshot"`
	HeroImage          *string          `json:"heroImage,omitempty"`
	HeightInInches     int              `json:"heightInInches"`
	WeightInPounds     int              `json:"weightInPounds"`
	BirthDate          string           `json:"birthDate"`
	BirthCity          *LocalizedString `json:"birthCity,omitempty"`
	BirthStateProvince *LocalizedString `json:"birthStateProvince,omitempty"`
	BirthCountry       *string          `json:"birthCountry,omitempty"`
	ShootsCatches      Handedness       `json:"shootsCatches"`
	DraftDetails       *DraftDetails    `json:"draftDetails,omitempty"`
	PlayerSlug         *string          `json:"playerSlug,omitempty"`
	FeaturedStats      *FeaturedStats   `json:"featuredStats,omitempty"`
	CareerTotals       *CareerTotals    `json:"careerTotals,omitempty"`
	SeasonTotals       []SeasonTotal    `json:"seasonTotals,omitempty"`
	Awards             []Award          `json:"awards,omitempty"`
	LastFiveGames      []GameLog        `json:"lastFiveGames,omitempty"`
}

// DraftDetails represents draft information for a player.
type DraftDetails struct {
	Year        int    `json:"year"`
	TeamAbbrev  string `json:"teamAbbrev"`
	Round       int    `json:"round"`
	PickInRound int    `json:"pickInRound"`
	OverallPick int    `json:"overallPick"`
}

// FeaturedStats represents featured statistics prominently shown on a player's page.
type FeaturedStats struct {
	Season        Season       `json:"season"`
	RegularSeason PlayerStats  `json:"regularSeason"`
	Playoffs      *PlayerStats `json:"playoffs,omitempty"`
}

// CareerTotals represents career totals for regular season and playoffs.
type CareerTotals struct {
	RegularSeason PlayerStats  `json:"regularSeason"`
	Playoffs      *PlayerStats `json:"playoffs,omitempty"`
}

// PlayerStats represents player statistics for either skaters or goalies.
type PlayerStats struct {
	// Common stats
	GamesPlayed *int `json:"gamesPlayed,omitempty"`

	// Skater stats
	Goals             *int     `json:"goals,omitempty"`
	Assists           *int     `json:"assists,omitempty"`
	Points            *int     `json:"points,omitempty"`
	PlusMinus         *int     `json:"plusMinus,omitempty"`
	PIM               *int     `json:"pim,omitempty"`
	PowerPlayGoals    *int     `json:"powerPlayGoals,omitempty"`
	PowerPlayPoints   *int     `json:"powerPlayPoints,omitempty"`
	ShortHandedGoals  *int     `json:"shortHandedGoals,omitempty"`
	ShortHandedPoints *int     `json:"shortHandedPoints,omitempty"`
	Shots             *int     `json:"shots,omitempty"`
	ShootingPctg      *float64 `json:"shootingPctg,omitempty"`
	FaceoffWinPctg    *float64 `json:"faceoffWinPctg,omitempty"`
	AvgTOI            *string  `json:"avgToi,omitempty"`

	// Goalie stats
	Wins            *int     `json:"wins,omitempty"`
	Losses          *int     `json:"losses,omitempty"`
	OTLosses        *int     `json:"otLosses,omitempty"`
	Shutouts        *int     `json:"shutouts,omitempty"`
	GoalsAgainstAvg *float64 `json:"goalsAgainstAvg,omitempty"`
	SavePctg        *float64 `json:"savePctg,omitempty"`
}

// SeasonTotal represents season-by-season statistics for a player.
type SeasonTotal struct {
	Season         Season           `json:"season"`
	GameType       GameType         `json:"gameTypeId"`
	LeagueAbbrev   string           `json:"leagueAbbrev"`
	TeamName       LocalizedString  `json:"teamName"`
	TeamCommonName *LocalizedString `json:"teamCommonName,omitempty"`
	Sequence       *int             `json:"sequence,omitempty"`
	GamesPlayed    int              `json:"gamesPlayed"`
	Goals          *int             `json:"goals,omitempty"`
	Assists        *int             `json:"assists,omitempty"`
	Points         *int             `json:"points,omitempty"`
	PlusMinus      *int             `json:"plusMinus,omitempty"`
	PIM            *int             `json:"pim,omitempty"`
}

// Award represents an award won by a player.
type Award struct {
	Trophy  LocalizedString `json:"trophy"`
	Seasons []AwardSeason   `json:"seasons"`
}

// AwardSeason represents a season when an award was won.
type AwardSeason struct {
	SeasonID Season `json:"seasonId"`
}

// GameLog represents a game log entry for a single game.
type GameLog struct {
	GameID           GameID   `json:"gameId"`
	GameDate         string   `json:"gameDate"`
	TeamAbbrev       string   `json:"teamAbbrev"`
	HomeRoadFlag     HomeRoad `json:"homeRoadFlag"`
	OpponentAbbrev   string   `json:"opponentAbbrev"`
	Goals            int      `json:"goals"`
	Assists          int      `json:"assists"`
	Points           int      `json:"points"`
	PlusMinus        int      `json:"plusMinus"`
	PowerPlayGoals   int      `json:"powerPlayGoals"`
	PowerPlayPoints  int      `json:"powerPlayPoints"`
	Shots            int      `json:"shots"`
	Shifts           int      `json:"shifts"`
	TOI              string   `json:"toi"`
	GameWinningGoals *int     `json:"gameWinningGoals,omitempty"`
	OTGoals          *int     `json:"otGoals,omitempty"`
	PIM              *int     `json:"pim,omitempty"`
}

// PlayerGameLog represents a player's game log for a season.
type PlayerGameLog struct {
	// PlayerID is not included in the API response, tracked manually
	PlayerID PlayerID  `json:"-"`
	Season   Season    `json:"seasonId"`
	GameType GameType  `json:"gameTypeId"`
	GameLog  []GameLog `json:"gameLog"`
}

// PlayerSearchResult represents a player search result from the NHL API.
type PlayerSearchResult struct {
	PlayerID           PlayerID `json:"playerId"`
	Name               string   `json:"name"`
	Position           Position `json:"positionCode"`
	TeamID             *TeamID  `json:"teamId,omitempty"`
	TeamAbbrev         *string  `json:"teamAbbrev,omitempty"`
	SweaterNumber      *int     `json:"sweaterNumber,omitempty"`
	Active             bool     `json:"active"`
	Height             *string  `json:"height,omitempty"`
	BirthCity          *string  `json:"birthCity,omitempty"`
	BirthStateProvince *string  `json:"birthStateProvince,omitempty"`
	BirthCountry       *string  `json:"birthCountry,omitempty"`
}
