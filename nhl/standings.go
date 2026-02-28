package nhl

import "fmt"

// Standing represents a team's standing entry with complete statistics.
// Contains conference, division, team identification, and win/loss records.
type Standing struct {
	ConferenceAbbrev *string         `json:"conferenceAbbrev,omitempty"`
	ConferenceName   *string         `json:"conferenceName,omitempty"`
	DivisionAbbrev   string          `json:"divisionAbbrev"`
	DivisionName     string          `json:"divisionName"`
	TeamName         LocalizedString `json:"teamName"`
	TeamCommonName   LocalizedString `json:"teamCommonName"`
	TeamAbbrev       LocalizedString `json:"teamAbbrev"`
	TeamLogo         string          `json:"teamLogo"`
	Wins             int             `json:"wins"`
	Losses           int             `json:"losses"`
	OTLosses         int             `json:"otLosses"`
	Points           int             `json:"points"`
}

const (
	unknownConferenceAbbrev = "UNK"
	unknownConferenceName   = "Unknown"
)

// conferenceAbbrev returns the conference abbreviation, or a default if not set.
// Used internally to handle historical data before conferences existed.
func (s *Standing) conferenceAbbrev() string {
	if s.ConferenceAbbrev == nil {
		return unknownConferenceAbbrev
	}
	return *s.ConferenceAbbrev
}

// conferenceName returns the conference name, or a default if not set.
// Used internally to handle historical data before conferences existed.
func (s *Standing) conferenceName() string {
	if s.ConferenceName == nil {
		return unknownConferenceName
	}
	return *s.ConferenceName
}

// ToTeam converts a Standing entry into a Team struct.
// This is useful for extracting team metadata from standings data.
func (s *Standing) ToTeam() Team {
	return Team{
		FullName:       s.TeamName.Default,
		TeamCommonName: LocalizedString{Default: s.TeamCommonName.Default},
		TeamPlaceName:  LocalizedString{Default: s.TeamName.Default},
		Tricode:        s.TeamAbbrev.Default,
		TeamLogo:       s.TeamLogo,
		Conference: Conference{
			Abbrev: s.conferenceAbbrev(),
			Name:   s.conferenceName(),
		},
		Division: Division{
			Abbrev: s.DivisionAbbrev,
			Name:   s.DivisionName,
		},
	}
}

// GamesPlayed calculates the total number of games played.
// Returns the sum of wins, losses, and overtime losses.
func (s *Standing) GamesPlayed() int {
	return s.Wins + s.Losses + s.OTLosses
}

// String implements fmt.Stringer for Standing.
// Returns a formatted string like "BOS: 31 pts (15-2-1)".
func (s Standing) String() string {
	return fmt.Sprintf("%s: %d pts (%d-%d-%d)",
		s.TeamAbbrev.Default,
		s.Points,
		s.Wins,
		s.Losses,
		s.OTLosses,
	)
}

// StandingsResponse represents the API response for standings queries.
// Contains a list of team standings.
type StandingsResponse struct {
	Standings []Standing `json:"standings"`
}

// SeasonInfo represents season metadata including standings date range.
// Used in the seasons manifest to identify valid season periods.
type SeasonInfo struct {
	ID             Season   `json:"id"`
	StandingsStart GameDate `json:"standingsStart"`
	StandingsEnd   GameDate `json:"standingsEnd"`
}

func (s SeasonInfo) Label() string {
	startYear := s.ID.StartYear()
	return fmt.Sprintf("%d-%02d", startYear, startYear+1)
}

// SeasonsResponse represents the API response for the seasons manifest.
// Contains a list of all available seasons.
type SeasonsResponse struct {
	Seasons []SeasonInfo `json:"seasons"`
}
