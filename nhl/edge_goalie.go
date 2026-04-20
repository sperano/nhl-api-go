package nhl

// ===== Edge Goalie Types =====

// EdgeGoalieDetail is the response from v1/edge/goalie-detail/{g}/{s}/{gt}.
type EdgeGoalieDetail struct {
	Player               EdgeGoaliePlayer                `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability         `json:"seasonsWithEdgeStats"`
	Stats                EdgeGoalieStatsSummary           `json:"stats"`
	ShotLocationSummary  []EdgeGoalieShotLocationSummary  `json:"shotLocationSummary"`
	ShotLocationDetails  []EdgeGoalieShotLocationArea     `json:"shotLocationDetails"`
}

// EdgeGoaliePlayer is the player metadata embedded in goalie Edge responses.
type EdgeGoaliePlayer struct {
	ID              int             `json:"id"`
	FirstName       LocalizedString `json:"firstName"`
	LastName        LocalizedString `json:"lastName"`
	BirthDate       string          `json:"birthDate"`
	ShootsCatches   string          `json:"shootsCatches"`
	SweaterNumber   int             `json:"sweaterNumber"`
	Slug            string          `json:"slug"`
	Headshot        string          `json:"headshot"`
	Wins            int             `json:"wins"`
	Losses          int             `json:"losses"`
	OvertimeLosses  int             `json:"overtimeLosses"`
	GoalsAgainstAvg float64         `json:"goalsAgainstAvg"`
	SavePctg        float64         `json:"savePctg"`
	GamesPlayed     int             `json:"gamesPlayed"`
	Team            EdgeTeamInfo    `json:"team"`
}

// EdgeGoalieStatsSummary contains the top-level goalie stat entries.
type EdgeGoalieStatsSummary struct {
	GoalsAgainstAvg       EdgeGoalieStatEntry `json:"goalsAgainstAvg"`
	GamesAbove900         EdgeGoalieStatEntry `json:"gamesAbove900"`
	GoalDifferentialPer60 EdgeGoalieStatEntry `json:"goalDifferentialPer60"`
	GoalSupportAvg        EdgeGoalieStatEntry `json:"goalSupportAvg"`
	PointPctg             EdgeGoalieStatEntry `json:"pointPctg"`
}

// EdgeGoalieStatEntry is a single goalie stat with percentile and league average.
type EdgeGoalieStatEntry struct {
	Value      float64 `json:"value"`
	Percentile float64 `json:"percentile"`
	LeagueAvg  float64 `json:"leagueAvg"`
}

// EdgeGoalieShotLocationSummary is a shot location summary by location code.
type EdgeGoalieShotLocationSummary struct {
	LocationCode           string  `json:"locationCode"`
	GoalsAgainst           int     `json:"goalsAgainst"`
	GoalsAgainstPercentile float64 `json:"goalsAgainstPercentile"`
	GoalsAgainstLeagueAvg  float64 `json:"goalsAgainstLeagueAvg"`
	Saves                  int     `json:"saves"`
	SavesPercentile        float64 `json:"savesPercentile"`
	SavesLeagueAvg         float64 `json:"savesLeagueAvg"`
	SavePctg               float64 `json:"savePctg"`
	SavePctgPercentile     float64 `json:"savePctgPercentile"`
	SavePctgLeagueAvg      float64 `json:"savePctgLeagueAvg"`
}

// EdgeGoalieShotLocationArea is a shot location detail for a specific rink area.
type EdgeGoalieShotLocationArea struct {
	Area               string  `json:"area"`
	Saves              int     `json:"saves"`
	SavesPercentile    float64 `json:"savesPercentile"`
	SavePctg           float64 `json:"savePctg"`
	SavePctgPercentile float64 `json:"savePctgPercentile"`
}

// EdgeGoalie5v5Detail is the response from v1/edge/goalie-5v5-detail/{g}/{s}/{gt}.
type EdgeGoalie5v5Detail struct {
	Player               EdgeGoaliePlayer        `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	SavePctg5v5Last10    []EdgeGoalie5v5Entry     `json:"savePctg5v5Last10"`
}

// EdgeGoalie5v5Entry is a per-game 5v5 save percentage entry.
type EdgeGoalie5v5Entry struct {
	GameDate string          `json:"gameDate"`
	AwayTeam EdgeOverlayTeam `json:"awayTeam"`
	HomeTeam EdgeOverlayTeam `json:"homeTeam"`
	SavePctg float64         `json:"savePctg"`
}

// EdgeGoalieShotLocationDetail is the response from v1/edge/goalie-shot-location-detail/{g}/{s}/{gt}.
type EdgeGoalieShotLocationDetail struct {
	Player               EdgeGoaliePlayer             `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability      `json:"seasonsWithEdgeStats"`
	ShotLocationDetails  []EdgeGoalieShotLocationEntry `json:"shotLocationDetails"`
}

// EdgeGoalieShotLocationEntry is a per-area shot location detail for a goalie.
type EdgeGoalieShotLocationEntry struct {
	Area     string  `json:"area"`
	Saves    int     `json:"saves"`
	SavePctg float64 `json:"savePctg"`
}

// EdgeGoalieSavePctgDetail is the response from v1/edge/goalie-save-percentage-detail/{g}/{s}/{gt}.
type EdgeGoalieSavePctgDetail struct {
	Player               EdgeGoaliePlayer          `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability   `json:"seasonsWithEdgeStats"`
	SavePctgLast10       []EdgeGoalieSavePctgEntry  `json:"savePctgLast10"`
	SavePctgDetails      []EdgeGoalieSavePctgEntry  `json:"savePctgDetails"`
}

// EdgeGoalieSavePctgEntry is a per-game save percentage entry.
type EdgeGoalieSavePctgEntry struct {
	GameDate string          `json:"gameDate"`
	AwayTeam EdgeOverlayTeam `json:"awayTeam"`
	HomeTeam EdgeOverlayTeam `json:"homeTeam"`
	SavePctg float64         `json:"savePctg"`
}

// EdgeGoalieComparison is the response from v1/edge/goalie-comparison/{g}/{s}/{gt}.
// Rich composite for head-to-head display. Cached on filesystem only.
type EdgeGoalieComparison struct {
	Player               EdgeGoaliePlayer        `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	ShotLocationSummary  interface{}              `json:"shotLocationSummary"`
	ShotLocationDetails  interface{}              `json:"shotLocationDetails"`
	SavePctg5v5Last10    interface{}              `json:"savePctg5v5Last10"`
	SavePctg5v5Details   interface{}              `json:"savePctg5v5Details"`
	SavePctgLast10       interface{}              `json:"savePctgLast10"`
	SavePctgDetails      interface{}              `json:"savePctgDetails"`
}

// EdgeGoalieLanding is the response from v1/edge/goalie-landing/{s}/{gt}.
// League-wide leaders in each Edge category. Cached on filesystem only.
type EdgeGoalieLanding struct {
	Leaders map[string]interface{} `json:"leaders"`
}
