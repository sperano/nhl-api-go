package nhl

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GameDate represents either the current date/time or a specific date for NHL games.
type GameDate struct {
	isNow bool
	date  time.Time
}

// Now creates a GameDate representing the current date/time.
func Now() GameDate {
	return GameDate{isNow: true}
}

// FromDate creates a GameDate from a specific time.Time.
func FromDate(t time.Time) GameDate {
	return GameDate{isNow: false, date: t}
}

// FromYMD creates a GameDate from year, month, and day components.
func FromYMD(year, month, day int) GameDate {
	return GameDate{
		isNow: false,
		date:  time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
	}
}

// Today creates a GameDate representing today's date.
func Today() GameDate {
	now := time.Now().UTC()
	return FromDate(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC))
}

// IsNow returns true if this GameDate represents the current time.
func (gd GameDate) IsNow() bool {
	return gd.isNow
}

// Date returns the underlying time.Time value.
// If IsNow is true, this returns the current time.
func (gd GameDate) Date() time.Time {
	if gd.isNow {
		return time.Now().UTC()
	}
	return gd.date
}

// ToAPIString converts the GameDate to the API format (YYYY-MM-DD).
// If IsNow is true, uses the current date.
func (gd GameDate) ToAPIString() string {
	d := gd.Date()
	return fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
}

// AddDays returns a new GameDate with the specified number of days added.
// If IsNow is true, it first resolves to the current date before adding.
func (gd GameDate) AddDays(days int) GameDate {
	d := gd.Date()
	newDate := d.AddDate(0, 0, days)
	return FromDate(newDate)
}

// String implements the fmt.Stringer interface.
func (gd GameDate) String() string {
	if gd.isNow {
		return "now"
	}
	return gd.ToAPIString()
}

// MarshalJSON implements json.Marshaler.
func (gd GameDate) MarshalJSON() ([]byte, error) {
	if gd.isNow {
		return json.Marshal("now")
	}
	return json.Marshal(gd.ToAPIString())
}

// UnmarshalJSON implements json.Unmarshaler.
func (gd *GameDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "now" {
		*gd = Now()
		return nil
	}

	// Parse YYYY-MM-DD format
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return fmt.Errorf("invalid date format: %s", s)
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid year: %s", parts[0])
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid month: %s", parts[1])
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid day: %s", parts[2])
	}

	*gd = FromYMD(year, month, day)
	return nil
}

// Season represents an NHL season.
type Season struct {
	startYear int
}

// NewSeason creates a new Season from a start year.
func NewSeason(startYear int) Season {
	return Season{startYear: startYear}
}

// FromYears creates a Season from start and end years.
// The start year is validated to match the expected end year.
func FromYears(startYear, endYear int) (Season, error) {
	expectedEnd := startYear + 1
	if endYear != expectedEnd {
		return Season{}, fmt.Errorf("invalid season years: %d-%d (expected %d-%d)",
			startYear, endYear, startYear, expectedEnd)
	}
	return NewSeason(startYear), nil
}

// StartYear returns the start year of the season.
func (s Season) StartYear() int {
	return s.startYear
}

// EndYear returns the end year of the season.
func (s Season) EndYear() int {
	return s.startYear + 1
}

// ToAPIString converts the Season to the API format (YYYYYYYY).
// For example, the 2023-2024 season is represented as "20232024".
func (s Season) ToAPIString() string {
	return fmt.Sprintf("%d%d", s.startYear, s.EndYear())
}

// String implements the fmt.Stringer interface.
// Returns the season in "YYYY-YYYY" format.
func (s Season) String() string {
	return fmt.Sprintf("%d-%d", s.startYear, s.EndYear())
}

// Parse parses a season string in either "YYYY-YYYY" or "YYYYYYYY" format.
func Parse(s string) (Season, error) {
	// Try YYYY-YYYY format first
	if strings.Contains(s, "-") {
		parts := strings.Split(s, "-")
		if len(parts) != 2 {
			return Season{}, fmt.Errorf("invalid season format: %s", s)
		}

		startYear, err := strconv.Atoi(parts[0])
		if err != nil {
			return Season{}, fmt.Errorf("invalid start year: %s", parts[0])
		}

		endYear, err := strconv.Atoi(parts[1])
		if err != nil {
			return Season{}, fmt.Errorf("invalid end year: %s", parts[1])
		}

		return FromYears(startYear, endYear)
	}

	// Try YYYYYYYY format
	if len(s) == 8 {
		startYear, err := strconv.Atoi(s[:4])
		if err != nil {
			return Season{}, fmt.Errorf("invalid start year: %s", s[:4])
		}

		endYear, err := strconv.Atoi(s[4:])
		if err != nil {
			return Season{}, fmt.Errorf("invalid end year: %s", s[4:])
		}

		return FromYears(startYear, endYear)
	}

	return Season{}, fmt.Errorf("invalid season format: %s", s)
}

// Current returns the current NHL season based on the current date.
// The NHL season typically starts in October and ends in June.
func Current() Season {
	now := time.Now()
	year := now.Year()
	month := now.Month()

	// If we're in January-June, we're in the season that started last year
	if month >= time.January && month <= time.June {
		return NewSeason(year - 1)
	}

	// Otherwise, we're in the season starting this year
	return NewSeason(year)
}

// MarshalJSON implements json.Marshaler.
func (s Season) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToAPIString())
}

// UnmarshalJSON implements json.Unmarshaler.
// Seasons can be unmarshaled from either integers (20232024) or strings ("20232024" or "2023-2024").
func (s *Season) UnmarshalJSON(data []byte) error {
	// Try unmarshaling as integer first (common API format)
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		season, err := SeasonFromInt64(i)
		if err != nil {
			return err
		}
		*s = season
		return nil
	}

	// Try unmarshaling as string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("season must be an integer or string: %w", err)
	}

	season, err := Parse(str)
	if err != nil {
		return err
	}

	*s = season
	return nil
}

// SeasonFromInt creates a Season from an integer in YYYYYYYY format.
func SeasonFromInt(i int) (Season, error) {
	if i < 10000000 || i > 99999999 {
		return Season{}, fmt.Errorf("invalid season integer: %d", i)
	}

	startYear := i / 10000
	endYear := i % 10000

	return FromYears(startYear, endYear)
}

// SeasonFromInt64 creates a Season from an int64 in YYYYYYYY format.
func SeasonFromInt64(i int64) (Season, error) {
	if i < 10000000 || i > 99999999 {
		return Season{}, fmt.Errorf("invalid season integer: %d", i)
	}

	return SeasonFromInt(int(i))
}

// ToInt converts the Season to an integer in YYYYYYYY format.
func (s Season) ToInt() int {
	i, _ := strconv.Atoi(s.ToAPIString())
	return i
}

// ToInt64 converts the Season to an int64 in YYYYYYYY format.
func (s Season) ToInt64() int64 {
	return int64(s.ToInt())
}
