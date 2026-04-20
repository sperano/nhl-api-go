package nhl

// ===== Edge Team Types =====

// EdgeTeamDetail is the response from v1/edge/team-detail/{t}/{s}/{gt}.
type EdgeTeamDetail struct {
	Team                 EdgeTeamInfo            `json:"team"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	ShotSpeed            EdgeTeamShotSpeed        `json:"shotSpeed"`
	SkatingSpeed         EdgeTeamSkatingSpeed     `json:"skatingSpeed"`
	DistanceSkated       EdgeTeamDistance          `json:"distanceSkated"`
	SogSummary           []EdgeTeamSogSummary     `json:"sogSummary"`
	SogDetails           []EdgeTeamSogAreaDetail  `json:"sogDetails"`
	ZoneTimeDetails      EdgeTeamZoneTime         `json:"zoneTimeDetails"`
}

// EdgeTeamShotSpeed contains team shot speed stats.
type EdgeTeamShotSpeed struct {
	ShotAttemptsOver90 EdgeRankStat            `json:"shotAttemptsOver90"`
	TopShotSpeed       EdgeRankStatWithOverlay `json:"topShotSpeed"`
}

// EdgeTeamSkatingSpeed contains team skating speed stats.
type EdgeTeamSkatingSpeed struct {
	BurstsOver22 EdgeRankStat            `json:"burstsOver22"`
	BurstsOver20 EdgeRankStat            `json:"burstsOver20"`
	SpeedMax     EdgeRankStatWithOverlay `json:"speedMax"`
}

// EdgeTeamDistance contains team distance skated stats.
type EdgeTeamDistance struct {
	Total EdgeRankStat `json:"total"`
}

// EdgeTeamSogSummary is a team shot-on-goal summary by location code.
type EdgeTeamSogSummary struct {
	LocationCode          string  `json:"locationCode"`
	Shots                 int     `json:"shots"`
	ShotsRank             int     `json:"shotsRank"`
	ShotsLeagueAvg        float64 `json:"shotsLeagueAvg"`
	Goals                 int     `json:"goals"`
	GoalsRank             int     `json:"goalsRank"`
	GoalsLeagueAvg        float64 `json:"goalsLeagueAvg"`
	ShootingPctg          float64 `json:"shootingPctg"`
	ShootingPctgRank      int     `json:"shootingPctgRank"`
	ShootingPctgLeagueAvg float64 `json:"shootingPctgLeagueAvg"`
}

// EdgeTeamSogAreaDetail is a team shot-on-goal detail for a specific rink area.
type EdgeTeamSogAreaDetail struct {
	Area      string `json:"area"`
	Shots     int    `json:"shots"`
	ShotsRank int    `json:"shotsRank"`
}

// EdgeTeamZoneTime contains team zone time percentages and ranks.
type EdgeTeamZoneTime struct {
	OffensiveZonePctg      float64 `json:"offensiveZonePctg"`
	OffensiveZoneRank      int     `json:"offensiveZoneRank"`
	OffensiveZoneLeagueAvg float64 `json:"offensiveZoneLeagueAvg"`
	OffensiveZoneEvPctg    float64 `json:"offensiveZoneEvPctg"`
	OffensiveZoneEvRank    int     `json:"offensiveZoneEvRank"`
	NeutralZonePctg        float64 `json:"neutralZonePctg"`
	NeutralZoneRank        int     `json:"neutralZoneRank"`
	NeutralZoneLeagueAvg   float64 `json:"neutralZoneLeagueAvg"`
	DefensiveZonePctg      float64 `json:"defensiveZonePctg"`
	DefensiveZoneRank      int     `json:"defensiveZoneRank"`
	DefensiveZoneLeagueAvg float64 `json:"defensiveZoneLeagueAvg"`
}

// EdgeTeamSpeedDetail is the response from v1/edge/team-skating-speed-detail/{t}/{s}/{gt}.
type EdgeTeamSpeedDetail struct {
	Team                 EdgeTeamInfo            `json:"team"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	TopSkatingSpeeds     []EdgeTeamSpeedEntry     `json:"topSkatingSpeeds"`
}

// EdgeTeamSpeedEntry is a per-player skating speed entry within a team.
type EdgeTeamSpeedEntry struct {
	Player EdgeOverlayPlayer `json:"player"`
	Speed  EdgeMeasurement   `json:"speed"`
}

// EdgeTeamDistanceDetail is the response from v1/edge/team-skating-distance-detail/{t}/{s}/{gt}.
type EdgeTeamDistanceDetail struct {
	Team                  EdgeTeamInfo            `json:"team"`
	SeasonsWithEdgeStats  []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	SkatingDistanceLast10 []EdgeTeamDistanceEntry  `json:"skatingDistanceLast10"`
}

// EdgeTeamDistanceEntry is a per-game team distance entry.
type EdgeTeamDistanceEntry struct {
	GameDate string          `json:"gameDate"`
	AwayTeam EdgeOverlayTeam `json:"awayTeam"`
	HomeTeam EdgeOverlayTeam `json:"homeTeam"`
	Distance EdgeMeasurement `json:"distance"`
}

// EdgeTeamShotSpeedDetail is the response from v1/edge/team-shot-speed-detail/{t}/{s}/{gt}.
type EdgeTeamShotSpeedDetail struct {
	Team                 EdgeTeamInfo            `json:"team"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	HardestShots         []EdgeTeamShotSpeedEntry `json:"hardestShots"`
}

// EdgeTeamShotSpeedEntry is a per-player shot speed entry within a team.
type EdgeTeamShotSpeedEntry struct {
	Player EdgeOverlayPlayer `json:"player"`
	Speed  EdgeMeasurement   `json:"speed"`
}

// EdgeTeamShotLocationDetail is the response from v1/edge/team-shot-location-detail/{t}/{s}/{gt}.
type EdgeTeamShotLocationDetail struct {
	Team                 EdgeTeamInfo                `json:"team"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability     `json:"seasonsWithEdgeStats"`
	ShotLocationDetails  []EdgeTeamShotLocationEntry  `json:"shotLocationDetails"`
}

// EdgeTeamShotLocationEntry is a shot location breakdown for a specific area.
type EdgeTeamShotLocationEntry struct {
	Area  string `json:"area"`
	Shots int    `json:"shots"`
	Goals int    `json:"goals"`
}

// EdgeTeamZoneTimeDetails is the response from v1/edge/team-zone-time-details/{t}/{s}/{gt}.
// This is DISTINCT from the zone time data embedded in team-detail — it breaks down
// zone time by strength code (all/es/pp/pk) and includes shot differential.
// Imported to DB in edge_team_zone_time_by_strength table.
type EdgeTeamZoneTimeDetails struct {
	Team                 EdgeTeamInfo                    `json:"team"`
	SeasonsWithEdgeStats []EdgeSeasonAvailability         `json:"seasonsWithEdgeStats"`
	ZoneTimeDetails      []EdgeTeamZoneTimeByStrength     `json:"zoneTimeDetails"`
	ShotDifferential     []EdgeTeamShotDifferentialEntry  `json:"shotDifferential"`
}

// EdgeTeamZoneTimeByStrength is zone time broken down by strength code.
type EdgeTeamZoneTimeByStrength struct {
	StrengthCode      string  `json:"strengthCode"`
	OffensiveZonePctg float64 `json:"offensiveZonePctg"`
	OffensiveZoneRank int     `json:"offensiveZoneRank"`
	NeutralZonePctg   float64 `json:"neutralZonePctg"`
	NeutralZoneRank   int     `json:"neutralZoneRank"`
	DefensiveZonePctg float64 `json:"defensiveZonePctg"`
	DefensiveZoneRank int     `json:"defensiveZoneRank"`
}

// EdgeTeamShotDifferentialEntry is shot differential by strength code.
type EdgeTeamShotDifferentialEntry struct {
	StrengthCode            string  `json:"strengthCode"`
	ForPerGame              float64 `json:"forPerGame"`
	ForPerGameRank          int     `json:"forPerGameRank"`
	AgainstPerGame          float64 `json:"againstPerGame"`
	AgainstPerGameRank      int     `json:"againstPerGameRank"`
	DifferentialPerGame     float64 `json:"differentialPerGame"`
	DifferentialPerGameRank int     `json:"differentialPerGameRank"`
}

// EdgeTeamComparison is the response from v1/edge/team-comparison/{t}/{s}/{gt}.
// Rich composite for head-to-head display. Cached on filesystem only.
type EdgeTeamComparison struct {
	Team                   EdgeTeamInfo            `json:"team"`
	SeasonsWithEdgeStats   []EdgeSeasonAvailability `json:"seasonsWithEdgeStats"`
	ShotSpeedDetails       interface{}              `json:"shotSpeedDetails"`
	SkatingSpeedDetails    interface{}              `json:"skatingSpeedDetails"`
	SkatingDistanceLast10  interface{}              `json:"skatingDistanceLast10"`
	SkatingDistanceDetails interface{}              `json:"skatingDistanceDetails"`
	ShotLocationDetails    interface{}              `json:"shotLocationDetails"`
	ShotLocationTotals     interface{}              `json:"shotLocationTotals"`
	ZoneTimeDetails        interface{}              `json:"zoneTimeDetails"`
	ShotDifferential       interface{}              `json:"shotDifferential"`
}

// EdgeTeamLanding is the response from v1/edge/team-landing/{s}/{gt}.
// League-wide leaders in each Edge category. Cached on filesystem only.
type EdgeTeamLanding struct {
	Leaders map[string]interface{} `json:"leaders"`
}
