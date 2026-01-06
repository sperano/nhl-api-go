package nhl

import (
	"encoding/json"
	"fmt"
	"time"
)

// LocalizedString represents a localized string from the NHL API.
// The NHL API returns localized strings in the format: {"default": "value"}
type LocalizedString struct {
	Default string `json:"default"`
}

// String returns the default localized string value.
func (l LocalizedString) String() string {
	return l.Default
}

// UnmarshalJSON implements custom JSON unmarshaling for LocalizedString.
// It handles both the standard {"default": "value"} format and plain string values.
func (l *LocalizedString) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as an object first
	var obj struct {
		Default string `json:"default"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		l.Default = obj.Default
		return nil
	}

	// Fall back to plain string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal LocalizedString: %w", err)
	}
	l.Default = s
	return nil
}

// MarshalJSON implements custom JSON marshaling for LocalizedString.
func (l LocalizedString) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Default string `json:"default"`
	}{
		Default: l.Default,
	})
}

// Conference represents an NHL conference.
type Conference struct {
	Abbrev string `json:"abbrev"`
	Name   string `json:"name"`
}

// Division represents an NHL division.
type Division struct {
	Abbrev string `json:"abbrev"`
	Name   string `json:"name"`
}

// Franchise represents an NHL franchise.
type Franchise struct {
	ID              int64  `json:"id"`
	FullName        string `json:"fullName"`
	TeamCommonName  string `json:"teamCommonName"`
	TeamPlaceName   string `json:"teamPlaceName"`
}

// Team represents an NHL team with all its metadata.
type Team struct {
	ID              int64           `json:"id"`
	FranchiseID     int64           `json:"franchiseId"`
	FullName        string          `json:"fullName"`
	LeagueAbbrev    string          `json:"leagueAbbrev"`
	RawTricode      string          `json:"rawTricode"`
	Tricode         string          `json:"tricode"`
	TeamPlaceName   LocalizedString `json:"teamPlaceName"`
	TeamCommonName  LocalizedString `json:"teamCommonName"`
	TeamLogo        string          `json:"teamLogo"`
	Conference      Conference      `json:"conference"`
	Division        Division        `json:"division"`
}

// Roster represents a team's roster organized by position.
type Roster struct {
	Forwards   []RosterPlayer `json:"forwards"`
	Defensemen []RosterPlayer `json:"defensemen"`
	Goalies    []RosterPlayer `json:"goalies"`
}

// AllPlayers returns all players on the roster in a single slice.
func (r *Roster) AllPlayers() []RosterPlayer {
	all := make([]RosterPlayer, 0, len(r.Forwards)+len(r.Defensemen)+len(r.Goalies))
	all = append(all, r.Forwards...)
	all = append(all, r.Defensemen...)
	all = append(all, r.Goalies...)
	return all
}

// PlayerCount returns the total number of players on the roster.
func (r *Roster) PlayerCount() int {
	return len(r.Forwards) + len(r.Defensemen) + len(r.Goalies)
}

// RosterPlayer represents a player on a team's roster.
type RosterPlayer struct {
	ID                  int64            `json:"id"`
	Headshot            string           `json:"headshot"`
	FirstName           LocalizedString  `json:"firstName"`
	LastName            LocalizedString  `json:"lastName"`
	SweaterNumber       int              `json:"sweaterNumber"`
	Position            Position         `json:"position"`
	ShootsCatches       Handedness       `json:"shootsCatches"`
	HeightInInches      int              `json:"heightInInches"`
	WeightInPounds      int              `json:"weightInPounds"`
	BirthDate           string           `json:"birthDate"`
	BirthCity           *LocalizedString `json:"birthCity,omitempty"`
	BirthStateProvince  *LocalizedString `json:"birthStateProvince,omitempty"`
	BirthCountry        string           `json:"birthCountry"`
}

// FullName returns the player's full name (first name + last name).
func (p *RosterPlayer) FullName() string {
	return p.FirstName.Default + " " + p.LastName.Default
}

// Age calculates the player's age based on their birth date.
// Returns the age in years, or -1 if the birth date cannot be parsed.
func (p *RosterPlayer) Age() int {
	birthDate, err := time.Parse("2006-01-02", p.BirthDate)
	if err != nil {
		return -1
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// Adjust if birthday hasn't occurred yet this year
	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

// HeightFeetInches returns the player's height formatted as feet and inches.
// For example, a player who is 73 inches tall returns "6'1\"".
func (p *RosterPlayer) HeightFeetInches() string {
	feet := p.HeightInInches / 12
	inches := p.HeightInInches % 12
	return fmt.Sprintf("%d'%d\"", feet, inches)
}

// BirthPlace returns a formatted string of the player's birth place.
// Format depends on which fields are available:
// - "City, State/Province, Country" (all fields present)
// - "City, Country" (no state/province)
// - "Country" (only country)
func (p *RosterPlayer) BirthPlace() string {
	parts := make([]string, 0, 3)

	if p.BirthCity != nil && p.BirthCity.Default != "" {
		parts = append(parts, p.BirthCity.Default)
	}

	if p.BirthStateProvince != nil && p.BirthStateProvince.Default != "" {
		parts = append(parts, p.BirthStateProvince.Default)
	}

	if p.BirthCountry != "" {
		parts = append(parts, p.BirthCountry)
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += ", "
		}
		result += part
	}

	return result
}
