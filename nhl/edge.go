package nhl

// ===== Edge Shared Types =====

// EdgeMeasurement represents a value in both imperial and metric units.
type EdgeMeasurement struct {
	Imperial float64 `json:"imperial"`
	Metric   float64 `json:"metric"`
}

// EdgeMeasurementWithOverlay is a measurement with game context overlay.
type EdgeMeasurementWithOverlay struct {
	Imperial float64      `json:"imperial"`
	Metric   float64      `json:"metric"`
	Overlay  *EdgeOverlay `json:"overlay,omitempty"`
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

// ===== Comparison Types (used in skater/team/goalie comparison endpoints) =====

// EdgeComparisonShotSpeedDetails contains shot speed breakdown for comparisons.
type EdgeComparisonShotSpeedDetails struct {
	TopShotSpeed       *EdgeMeasurementWithOverlay `json:"topShotSpeed,omitempty"`
	AvgShotSpeed       *EdgeMeasurement            `json:"avgShotSpeed,omitempty"`
	ShotAttemptsOver100 int                        `json:"shotAttemptsOver100,omitempty"`
	ShotAttempts90To100 int                        `json:"shotAttempts90To100,omitempty"`
	ShotAttempts80To90  int                        `json:"shotAttempts80To90,omitempty"`
	ShotAttempts70To80  int                        `json:"shotAttempts70To80,omitempty"`
}

// EdgeComparisonSkatingSpeedDetails contains skating speed breakdown for comparisons.
type EdgeComparisonSkatingSpeedDetails struct {
	MaxSkatingSpeed *EdgeMeasurementWithOverlay `json:"maxSkatingSpeed,omitempty"`
	BurstsOver22    int                         `json:"burstsOver22,omitempty"`
	Bursts20To22    int                         `json:"bursts20To22,omitempty"`
	Bursts18To20    int                         `json:"bursts18To20,omitempty"`
}

// EdgeComparisonSkatingDistanceDetails contains distance breakdown for comparisons.
type EdgeComparisonSkatingDistanceDetails struct {
	DistanceTotal     *EdgeMeasurement            `json:"distanceTotal,omitempty"`
	DistancePer60     *EdgeMeasurement            `json:"distancePer60,omitempty"`
	DistanceMaxGame   *EdgeMeasurementWithOverlay `json:"distanceMaxGame,omitempty"`
	DistanceMaxPeriod *EdgeMeasurementWithOverlay `json:"distanceMaxPeriod,omitempty"`
}

// EdgeComparisonZoneTimeDetails contains zone time percentages for comparisons.
type EdgeComparisonZoneTimeDetails struct {
	OffensiveZonePctg      float64 `json:"offensiveZonePctg"`
	OffensiveZoneLeagueAvg float64 `json:"offensiveZoneLeagueAvg,omitempty"`
	NeutralZonePctg        float64 `json:"neutralZonePctg"`
	NeutralZoneLeagueAvg   float64 `json:"neutralZoneLeagueAvg,omitempty"`
	DefensiveZonePctg      float64 `json:"defensiveZonePctg"`
	DefensiveZoneLeagueAvg float64 `json:"defensiveZoneLeagueAvg,omitempty"`
}

// EdgeComparisonShotLocationDetail is a shot location breakdown by area.
type EdgeComparisonShotLocationDetail struct {
	Area         string  `json:"area"`
	SOG          int     `json:"sog"`
	Goals        int     `json:"goals"`
	ShootingPctg float64 `json:"shootingPctg"`
}

// EdgeComparisonShotLocationTotal is shot totals by location code.
type EdgeComparisonShotLocationTotal struct {
	LocationCode string  `json:"locationCode"`
	SOG          int     `json:"sog"`
	Goals        int     `json:"goals"`
	ShootingPctg float64 `json:"shootingPctg"`
}

// EdgeComparisonDistanceLast10Entry is a per-game distance entry in last 10 games.
type EdgeComparisonDistanceLast10Entry struct {
	GameCenterLink   string           `json:"gameCenterLink,omitempty"`
	GameDate         string           `json:"gameDate"`
	PlayerOnHomeTeam bool             `json:"playerOnHomeTeam,omitempty"`
	DistanceSkated   *EdgeMeasurement `json:"distanceSkated,omitempty"`
	TOI              float64          `json:"toi,omitempty"`
	HomeTeam         *EdgeOverlayTeam `json:"homeTeam,omitempty"`
	AwayTeam         *EdgeOverlayTeam `json:"awayTeam,omitempty"`
	// Team comparison uses different field names
	Distance *EdgeMeasurement `json:"distance,omitempty"`
}

// EdgeComparisonZoneStarts contains zone start percentages.
type EdgeComparisonZoneStarts struct {
	OffensiveZoneStarts float64 `json:"offensiveZoneStarts"`
	NeutralZoneStarts   float64 `json:"neutralZoneStarts"`
	DefensiveZoneStarts float64 `json:"defensiveZoneStarts"`
}

// EdgeSkaterComparison is the response from v1/edge/skater-comparison/{p}/{s}/{gt}.
// Rich composite for head-to-head display. Cached on filesystem only.
type EdgeSkaterComparison struct {
	Player                 EdgeSkaterPlayer                     `json:"player"`
	SeasonsWithEdgeStats   []EdgeSeasonAvailability             `json:"seasonsWithEdgeStats"`
	ShotSpeedDetails       *EdgeComparisonShotSpeedDetails      `json:"shotSpeedDetails,omitempty"`
	SkatingSpeedDetails    *EdgeComparisonSkatingSpeedDetails   `json:"skatingSpeedDetails,omitempty"`
	SkatingDistanceLast10  []EdgeComparisonDistanceLast10Entry  `json:"skatingDistanceLast10,omitempty"`
	SkatingDistanceDetails *EdgeComparisonSkatingDistanceDetails `json:"skatingDistanceDetails,omitempty"`
	ShotLocationDetails    []EdgeComparisonShotLocationDetail   `json:"shotLocationDetails,omitempty"`
	ShotLocationTotals     []EdgeComparisonShotLocationTotal    `json:"shotLocationTotals,omitempty"`
	ZoneTimeDetails        *EdgeComparisonZoneTimeDetails       `json:"zoneTimeDetails,omitempty"`
	ZoneStarts             *EdgeComparisonZoneStarts            `json:"zoneStarts,omitempty"`
}

// EdgeLeaderShotLocation is shot location detail in leader responses.
type EdgeLeaderShotLocation struct {
	Area              string   `json:"area"`
	SOG               *int     `json:"sog,omitempty"`               // skater: shots on goal
	SOGPercentile     *float64 `json:"sogPercentile,omitempty"`     // skater percentile
	SavePctg          *float64 `json:"savePctg,omitempty"`          // goalie: save percentage
	SavePctgPercentile *float64 `json:"savePctgPercentile,omitempty"` // goalie percentile
}

// EdgeSkaterLeader is a leader entry in the skater landing response.
type EdgeSkaterLeader struct {
	Player  EdgeSkaterPlayer    `json:"player"`
	Overlay *EdgeOverlay        `json:"overlay,omitempty"`
	// Stat fields - only one set is populated per category
	ShotSpeed           *EdgeMeasurement         `json:"shotSpeed,omitempty"`           // hardestShot
	SkatingSpeed        *EdgeMeasurement         `json:"skatingSpeed,omitempty"`        // maxSkatingSpeed
	DistanceSkated      *EdgeMeasurement         `json:"distanceSkated,omitempty"`      // totalDistanceSkated, distanceMaxGame
	SOG                 *int                     `json:"sog,omitempty"`                 // highDangerSOG
	ShotLocationDetails []EdgeLeaderShotLocation `json:"shotLocationDetails,omitempty"` // highDangerSOG
	ZoneTime            *float64                 `json:"zoneTime,omitempty"`            // offensiveZoneTime, defensiveZoneTime
}

// EdgeSkaterLanding is the response from v1/edge/skater-landing/{s}/{gt}.
// League-wide leaders in each Edge category.
type EdgeSkaterLanding struct {
	SeasonsWithEdgeStats []EdgeSeasonAvailability    `json:"seasonsWithEdgeStats"`
	Leaders              map[string]EdgeSkaterLeader `json:"leaders"`
}
