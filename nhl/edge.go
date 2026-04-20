package nhl

// ===== Edge Shared Types =====

// EdgeMeasurement represents a value in both imperial and metric units.
type EdgeMeasurement struct {
	Imperial float64 `json:"imperial"`
	Metric   float64 `json:"metric"`
}

// EdgePercentileStat is a measurement with league-relative percentile.
type EdgePercentileStat struct {
	Imperial   float64         `json:"imperial"`
	Metric     float64         `json:"metric"`
	Percentile float64         `json:"percentile"`
	LeagueAvg  EdgeMeasurement `json:"leagueAvg"`
}

// EdgePercentileStatWithOverlay adds game context to a percentile stat.
type EdgePercentileStatWithOverlay struct {
	Imperial   float64         `json:"imperial"`
	Metric     float64         `json:"metric"`
	Percentile float64         `json:"percentile"`
	LeagueAvg  EdgeMeasurement `json:"leagueAvg"`
	Overlay    *EdgeOverlay    `json:"overlay,omitempty"`
}

// EdgeCountLeagueAvg is the league average for a count-based stat.
type EdgeCountLeagueAvg struct {
	Value float64 `json:"value"`
}

// EdgeCountPercentileStat is a count-based stat with percentile and league average.
type EdgeCountPercentileStat struct {
	Value      int                `json:"value"`
	Percentile float64            `json:"percentile"`
	LeagueAvg  EdgeCountLeagueAvg `json:"leagueAvg"`
}

// EdgeRankStat is a count-based stat with league rank (1-32) instead of percentile.
type EdgeRankStat struct {
	Value     int                 `json:"value"`
	Rank      int                 `json:"rank"`
	LeagueAvg *EdgeCountLeagueAvg `json:"leagueAvg,omitempty"`
}

// EdgeRankStatWithOverlay adds game context to a rank-based stat.
type EdgeRankStatWithOverlay struct {
	Imperial  float64         `json:"imperial"`
	Metric    float64         `json:"metric"`
	Rank      int             `json:"rank"`
	LeagueAvg EdgeMeasurement `json:"leagueAvg"`
	Overlay   *EdgeOverlay    `json:"overlay,omitempty"`
}

// EdgeOverlay provides game context for a "best-of" stat.
type EdgeOverlay struct {
	Player           EdgeOverlayPlayer `json:"player"`
	GameDate         string            `json:"gameDate"`
	AwayTeam         EdgeOverlayTeam   `json:"awayTeam"`
	HomeTeam         EdgeOverlayTeam   `json:"homeTeam"`
	GameOutcome      *GameOutcome      `json:"gameOutcome,omitempty"`
	PeriodDescriptor PeriodDescriptor  `json:"periodDescriptor"`
	TimeInPeriod     string            `json:"timeInPeriod"`
	GameType         int               `json:"gameType"`
}

// EdgeOverlayPlayer identifies a player within an overlay.
type EdgeOverlayPlayer struct {
	FirstName LocalizedString `json:"firstName"`
	LastName  LocalizedString `json:"lastName"`
}

// EdgeOverlayTeam identifies a team within an overlay.
type EdgeOverlayTeam struct {
	Abbrev string `json:"abbrev"`
	Score  int    `json:"score"`
}

// EdgeSeasonAvailability indicates which seasons/game types have Edge data.
type EdgeSeasonAvailability struct {
	ID        int   `json:"id"`
	GameTypes []int `json:"gameTypes"`
}

// EdgeTeamLogo contains light/dark logo URLs.
type EdgeTeamLogo struct {
	Light string `json:"light"`
	Dark  string `json:"dark"`
}

// EdgeTeamInfo is the team metadata embedded in Edge responses.
type EdgeTeamInfo struct {
	ID                       int             `json:"id"`
	CommonName               LocalizedString `json:"commonName"`
	PlaceNameWithPreposition LocalizedString `json:"placeNameWithPreposition"`
	Abbrev                   string          `json:"abbrev"`
	TeamLogo                 EdgeTeamLogo    `json:"teamLogo"`
	Slug                     string          `json:"slug"`
	Conference               string          `json:"conference"`
	Division                 string          `json:"division"`
	Wins                     int             `json:"wins"`
	Losses                   int             `json:"losses"`
	OTLosses                 int             `json:"otLosses"`
	GamesPlayed              int             `json:"gamesPlayed"`
	Points                   int             `json:"points"`
}

// ===== Edge Skater Types =====

// EdgeSkaterDetail is the response from v1/edge/skater-detail/{p}/{s}/{gt}.
type EdgeSkaterDetail struct {
	Player               EdgeSkaterPlayer              `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability       `json:"seasonsWithEdgeStats"`
	TopShotSpeed         EdgePercentileStatWithOverlay  `json:"topShotSpeed"`
	SkatingSpeed         EdgeSkaterSpeed                `json:"skatingSpeed"`
	TotalDistanceSkated  EdgePercentileStat             `json:"totalDistanceSkated"`
	DistanceMaxGame      EdgePercentileStatWithOverlay  `json:"distanceMaxGame"`
	SogSummary           []EdgeSkaterSogSummary         `json:"sogSummary"`
	SogDetails           []EdgeSogAreaDetail            `json:"sogDetails"`
	ZoneTimeDetails      EdgeSkaterZoneTimeSummary      `json:"zoneTimeDetails"`
}

// EdgeSkaterPlayer is the player metadata embedded in skater Edge responses.
type EdgeSkaterPlayer struct {
	ID            int             `json:"id"`
	FirstName     LocalizedString `json:"firstName"`
	LastName      LocalizedString `json:"lastName"`
	BirthDate     string          `json:"birthDate"`
	ShootsCatches string          `json:"shootsCatches"`
	SweaterNumber int             `json:"sweaterNumber"`
	Position      string          `json:"position"`
	Slug          string          `json:"slug"`
	Headshot      string          `json:"headshot"`
	Goals         int             `json:"goals"`
	Assists       int             `json:"assists"`
	Points        int             `json:"points"`
	GamesPlayed   int             `json:"gamesPlayed"`
	Team          EdgeTeamInfo    `json:"team"`
}

// EdgeSkaterSpeed contains skating speed stats with burst counts.
type EdgeSkaterSpeed struct {
	SpeedMax     EdgePercentileStatWithOverlay `json:"speedMax"`
	BurstsOver20 EdgeCountPercentileStat       `json:"burstsOver20"`
}

// EdgeSkaterSogSummary is a shot-on-goal summary by location code.
type EdgeSkaterSogSummary struct {
	LocationCode           string  `json:"locationCode"`
	Shots                  int     `json:"shots"`
	ShotsPercentile        float64 `json:"shotsPercentile"`
	ShotsLeagueAvg         float64 `json:"shotsLeagueAvg"`
	Goals                  int     `json:"goals"`
	GoalsPercentile        float64 `json:"goalsPercentile"`
	GoalsLeagueAvg         float64 `json:"goalsLeagueAvg"`
	ShootingPctg           float64 `json:"shootingPctg"`
	ShootingPctgPercentile float64 `json:"shootingPctgPercentile"`
	ShootingPctgLeagueAvg  float64 `json:"shootingPctgLeagueAvg"`
}

// EdgeSogAreaDetail is a shot-on-goal detail for a specific rink area.
type EdgeSogAreaDetail struct {
	Area            string  `json:"area"`
	Shots           int     `json:"shots,omitempty"`
	ShootingPctg    float64 `json:"shootingPctg,omitempty"`
	ShotsPercentile float64 `json:"shotsPercentile,omitempty"`
}

// EdgeSkaterZoneTimeSummary contains zone time percentages and percentiles for a skater.
type EdgeSkaterZoneTimeSummary struct {
	OffensiveZonePctg         float64 `json:"offensiveZonePctg"`
	OffensiveZonePercentile   float64 `json:"offensiveZonePercentile"`
	OffensiveZoneLeagueAvg    float64 `json:"offensiveZoneLeagueAvg"`
	OffensiveZoneEvPctg       float64 `json:"offensiveZoneEvPctg"`
	OffensiveZoneEvPercentile float64 `json:"offensiveZoneEvPercentile"`
	OffensiveZoneEvLeagueAvg  float64 `json:"offensiveZoneEvLeagueAvg"`
	NeutralZonePctg           float64 `json:"neutralZonePctg"`
	NeutralZonePercentile     float64 `json:"neutralZonePercentile"`
	NeutralZoneLeagueAvg      float64 `json:"neutralZoneLeagueAvg"`
	DefensiveZonePctg         float64 `json:"defensiveZonePctg"`
	DefensiveZonePercentile   float64 `json:"defensiveZonePercentile"`
	DefensiveZoneLeagueAvg    float64 `json:"defensiveZoneLeagueAvg"`
}

// EdgeSkaterSpeedDetail is the response from v1/edge/skater-skating-speed-detail/{p}/{s}/{gt}.
type EdgeSkaterSpeedDetail struct {
	Player               EdgeSkaterPlayer        `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	TopSkatingSpeeds     []EdgeSpeedEntry         `json:"topSkatingSpeeds"`
}

// EdgeSpeedEntry is a per-game skating speed entry.
type EdgeSpeedEntry struct {
	GameDate string          `json:"gameDate"`
	AwayTeam EdgeOverlayTeam `json:"awayTeam"`
	HomeTeam EdgeOverlayTeam `json:"homeTeam"`
	Speed    EdgeMeasurement `json:"speed"`
}

// EdgeSkaterDistanceDetail is the response from v1/edge/skater-skating-distance-detail/{p}/{s}/{gt}.
type EdgeSkaterDistanceDetail struct {
	Player                EdgeSkaterPlayer        `json:"player"`
	SeasonsWithEdgeStats  []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	SkatingDistanceLast10 []EdgeDistanceEntry      `json:"skatingDistanceLast10"`
}

// EdgeDistanceEntry is a per-game distance skated entry.
type EdgeDistanceEntry struct {
	GameDate string          `json:"gameDate"`
	AwayTeam EdgeOverlayTeam `json:"awayTeam"`
	HomeTeam EdgeOverlayTeam `json:"homeTeam"`
	Distance EdgeMeasurement `json:"distance"`
}

// EdgeSkaterShotSpeedDetail is the response from v1/edge/skater-shot-speed-detail/{p}/{s}/{gt}.
type EdgeSkaterShotSpeedDetail struct {
	Player               EdgeSkaterPlayer        `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	HardestShots         []EdgeShotSpeedEntry     `json:"hardestShots"`
}

// EdgeShotSpeedEntry is a per-game shot speed entry.
type EdgeShotSpeedEntry struct {
	GameDate string          `json:"gameDate"`
	AwayTeam EdgeOverlayTeam `json:"awayTeam"`
	HomeTeam EdgeOverlayTeam `json:"homeTeam"`
	Speed    EdgeMeasurement `json:"speed"`
}

// EdgeSkaterShotLocationDetail is the response from v1/edge/skater-shot-location-detail/{p}/{s}/{gt}.
type EdgeSkaterShotLocationDetail struct {
	Player               EdgeSkaterPlayer        `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	ShotLocationDetails  []EdgeShotLocationEntry  `json:"shotLocationDetails"`
}

// EdgeShotLocationEntry is a shot location breakdown for a specific area.
type EdgeShotLocationEntry struct {
	Area         string  `json:"area"`
	Shots        int     `json:"shots"`
	Goals        int     `json:"goals"`
	ShootingPctg float64 `json:"shootingPctg"`
}

// EdgeSkaterZoneTimeDetail is the response from v1/edge/skater-zone-time/{p}/{s}/{gt}.
type EdgeSkaterZoneTimeDetail struct {
	Player               EdgeSkaterPlayer        `json:"player"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	ZoneTimeDetails      []EdgeZoneTimeEntry      `json:"zoneTimeDetails"`
}

// EdgeZoneTimeEntry is a zone time breakdown by strength code.
type EdgeZoneTimeEntry struct {
	StrengthCode      string  `json:"strengthCode"`
	OffensiveZonePctg float64 `json:"offensiveZonePctg"`
	NeutralZonePctg   float64 `json:"neutralZonePctg"`
	DefensiveZonePctg float64 `json:"defensiveZonePctg"`
}

// EdgeSkaterComparison is the response from v1/edge/skater-comparison/{p}/{s}/{gt}.
// Rich composite for head-to-head display. Cached on filesystem only.
type EdgeSkaterComparison struct {
	Player                 EdgeSkaterPlayer        `json:"player"`
	SeasonsWithEdgeStats   []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	ShotSpeedDetails       interface{}              `json:"shotSpeedDetails"`
	SkatingSpeedDetails    interface{}              `json:"skatingSpeedDetails"`
	SkatingDistanceLast10  interface{}              `json:"skatingDistanceLast10"`
	SkatingDistanceDetails interface{}              `json:"skatingDistanceDetails"`
	ShotLocationDetails    interface{}              `json:"shotLocationDetails"`
	ShotLocationTotals     interface{}              `json:"shotLocationTotals"`
	ZoneTimeDetails        interface{}              `json:"zoneTimeDetails"`
	ZoneStarts             interface{}              `json:"zoneStarts"`
}

// EdgeSkaterLanding is the response from v1/edge/skater-landing/{s}/{gt}.
// League-wide leaders in each Edge category. Cached on filesystem only.
type EdgeSkaterLanding struct {
	Leaders map[string]interface{} `json:"leaders"`
}
