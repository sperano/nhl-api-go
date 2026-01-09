package nhl

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ClubSkaterStats represents skater season statistics for a team.
type ClubSkaterStats struct {
	PlayerID            PlayerID        `json:"playerId"`
	Headshot            string          `json:"headshot"`
	FirstName           LocalizedString `json:"firstName"`
	LastName            LocalizedString `json:"lastName"`
	Position            Position        `json:"positionCode"`
	GamesPlayed         int             `json:"gamesPlayed"`
	Goals               int             `json:"goals"`
	Assists             int             `json:"assists"`
	Points              int             `json:"points"`
	PlusMinus           int             `json:"plusMinus"`
	PenaltyMinutes      int             `json:"penaltyMinutes"`
	PowerPlayGoals      int             `json:"powerPlayGoals"`
	ShorthandedGoals    int             `json:"shorthandedGoals"`
	GameWinningGoals    int             `json:"gameWinningGoals"`
	OvertimeGoals       int             `json:"overtimeGoals"`
	Shots               int             `json:"shots"`
	ShootingPctg        float64         `json:"shootingPctg"`
	AvgTimeOnIcePerGame float64         `json:"avgTimeOnIcePerGame"`
	AvgShiftsPerGame    float64         `json:"avgShiftsPerGame"`
	FaceoffWinPctg      float64         `json:"faceoffWinPctg"`
}

// String returns a formatted string representation of the skater stats.
func (s ClubSkaterStats) String() string {
	return fmt.Sprintf("%s %s - %d GP, %d G, %d A, %d PTS",
		s.FirstName.Default,
		s.LastName.Default,
		s.GamesPlayed,
		s.Goals,
		s.Assists,
		s.Points,
	)
}

// ClubGoalieStats represents goalie season statistics for a team.
type ClubGoalieStats struct {
	PlayerID            PlayerID        `json:"playerId"`
	Headshot            string          `json:"headshot"`
	FirstName           LocalizedString `json:"firstName"`
	LastName            LocalizedString `json:"lastName"`
	GamesPlayed         int             `json:"gamesPlayed"`
	GamesStarted        int             `json:"gamesStarted"`
	Wins                int             `json:"wins"`
	Losses              int             `json:"losses"`
	OvertimeLosses      int             `json:"overtimeLosses"`
	GoalsAgainstAverage float64         `json:"goalsAgainstAverage"`
	SavePercentage      float64         `json:"savePercentage"`
	ShotsAgainst        int             `json:"shotsAgainst"`
	Saves               int             `json:"saves"`
	GoalsAgainst        int             `json:"goalsAgainst"`
	Shutouts            int             `json:"shutouts"`
	Goals               int             `json:"goals"`
	Assists             int             `json:"assists"`
	Points              int             `json:"points"`
	PenaltyMinutes      int             `json:"penaltyMinutes"`
	TimeOnIce           int64           `json:"timeOnIce"`
}

// String returns a formatted string representation of the goalie stats.
func (g ClubGoalieStats) String() string {
	return fmt.Sprintf("%s %s - %d GP, %d-%d-%d, %.3f GAA, %.3f SV%%",
		g.FirstName.Default,
		g.LastName.Default,
		g.GamesPlayed,
		g.Wins,
		g.Losses,
		g.OvertimeLosses,
		g.GoalsAgainstAverage,
		g.SavePercentage,
	)
}

// ClubStats represents club statistics response containing skater and goalie stats.
type ClubStats struct {
	Season   string            `json:"season"`
	GameType GameType          `json:"gameType"`
	Skaters  []ClubSkaterStats `json:"skaters"`
	Goalies  []ClubGoalieStats `json:"goalies"`
}

// SeasonGameTypes represents season game type availability for a team.
type SeasonGameTypes struct {
	Season    Season     `json:"season"`
	GameTypes []GameType `json:"gameTypes"`
}

// UnmarshalJSON implements custom JSON unmarshaling for SeasonGameTypes.
// The NHL API returns gameTypes as an array of integers, which we convert to GameType enums.
func (s *SeasonGameTypes) UnmarshalJSON(data []byte) error {
	var raw struct {
		Season    int   `json:"season"`
		GameTypes []int `json:"gameTypes"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	season, err := SeasonFromInt(raw.Season)
	if err != nil {
		return fmt.Errorf("unmarshaling season: %w", err)
	}
	s.Season = season
	s.GameTypes = make([]GameType, 0, len(raw.GameTypes))

	for _, gt := range raw.GameTypes {
		gameType, err := GameTypeFromInt(gt)
		if err != nil {
			return fmt.Errorf("unmarshaling game types: %w", err)
		}
		s.GameTypes = append(s.GameTypes, gameType)
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for SeasonGameTypes.
// Serializes GameTypes as an array of integers to match the NHL API format.
func (s SeasonGameTypes) MarshalJSON() ([]byte, error) {
	gameTypeInts := make([]int, len(s.GameTypes))
	for i, gt := range s.GameTypes {
		if !gt.IsValid() {
			return nil, fmt.Errorf("invalid game type at index %d: %d", i, gt)
		}
		gameTypeInts[i] = gt.ToInt()
	}

	raw := struct {
		Season    int   `json:"season"`
		GameTypes []int `json:"gameTypes"`
	}{
		Season:    s.Season.ToInt(),
		GameTypes: gameTypeInts,
	}

	return json.Marshal(raw)
}

// String returns a formatted string representation of the season game types.
func (s SeasonGameTypes) String() string {
	gameTypeStrs := make([]string, len(s.GameTypes))
	for i, gt := range s.GameTypes {
		gameTypeStrs[i] = gt.String()
	}
	return fmt.Sprintf("%s: %s", s.Season.String(), strings.Join(gameTypeStrs, ", "))
}
