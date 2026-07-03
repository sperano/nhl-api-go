package nhl

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Date represents a date without time components (YYYY-MM-DD).
// Used for API responses that contain concrete date values.
type Date struct {
	time.Time
}

// DateLayout is the format used for parsing and formatting dates.
const DateLayout = "2006-01-02"

// NewDate creates a Date from year, month, and day components.
func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// NewDateYMD creates a Date from year, month (as int), and day components.
func NewDateYMD(year, month, day int) Date {
	return NewDate(year, time.Month(month), day)
}

// DateFromTime creates a Date from a time.Time, truncating time components.
func DateFromTime(t time.Time) Date {
	return NewDate(t.Year(), t.Month(), t.Day())
}

// ParseDate parses a date string in YYYY-MM-DD format.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(DateLayout, s)
	if err != nil {
		return Date{}, fmt.Errorf("invalid date format %q: %w", s, err)
	}
	return Date{t}, nil
}

// MustParseDate parses a date string in YYYY-MM-DD format, panicking on error.
// Intended for test data setup.
func MustParseDate(s string) Date {
	d, err := ParseDate(s)
	if err != nil {
		panic(err)
	}
	return d
}

// String returns the date in YYYY-MM-DD format.
func (d Date) String() string {
	return d.Time.Format(DateLayout)
}

// Equal returns true if two dates represent the same calendar day. It compares
// the year/month/day components rather than the underlying instants, so two
// Date values for the same day compare equal even if their time-of-day or
// location differ (the package constructors normalize to midnight UTC, but a
// Date populated directly need not).
func (d Date) Equal(other Date) bool {
	y1, m1, d1 := d.Time.Date()
	y2, m2, d2 := other.Time.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// MarshalJSON implements json.Marshaler.
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Date) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// GobEncode implements gob.GobEncoder.
func (d Date) GobEncode() ([]byte, error) {
	return d.Time.GobEncode()
}

// GobDecode implements gob.GobDecoder.
func (d *Date) GobDecode(data []byte) error {
	return d.Time.GobDecode(data)
}

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

// APIString converts the GameDate to the API format (YYYY-MM-DD).
// If IsNow is true, uses the current date.
func (gd GameDate) APIString() string {
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
	return gd.APIString()
}

// MarshalJSON implements json.Marshaler.
func (gd GameDate) MarshalJSON() ([]byte, error) {
	if gd.isNow {
		return json.Marshal("now")
	}
	return json.Marshal(gd.APIString())
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

	// Delegate to ParseDate (time.Parse with DateLayout), which enforces the
	// exact YYYY-MM-DD shape and validates ranges. The previous hand-rolled
	// split+Atoi accepted overflow values like "2024-13-40" because FromYMD
	// relies on time.Date, which silently normalizes them.
	d, err := ParseDate(s)
	if err != nil {
		return err
	}
	*gd = FromDate(d.Time)
	return nil
}

// GobEncode implements gob.GobEncoder for GameDate.
func (gd GameDate) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(gd.isNow); err != nil {
		return nil, err
	}
	if err := enc.Encode(gd.date); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode implements gob.GobDecoder for GameDate.
func (gd *GameDate) GobDecode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&gd.isNow); err != nil {
		return err
	}
	return dec.Decode(&gd.date)
}

// Season represents an NHL season. It stores both the start and end years
// explicitly rather than assuming end == start+1: while NHL season IDs are
// conventionally cross-year (e.g. 20232024), seasons played within a single
// calendar year exist (the COVID-shortened 2020-21 season ran Jan–May 2021),
// so the end year is recorded as parsed rather than derived.
type Season struct {
	startYear int
	endYear   int
}

// NewSeason creates a new Season from a start year, using the conventional
// cross-year end (startYear + 1).
func NewSeason(startYear int) Season {
	return Season{startYear: startYear, endYear: startYear + 1}
}

// FromYears creates a Season from start and end years.
// The end year must equal startYear (a single-calendar-year season) or
// startYear + 1 (the typical cross-year season like 2023-2024); both are
// stored as given. Returns an error for any other range.
func FromYears(startYear, endYear int) (Season, error) {
	if endYear == startYear || endYear == startYear+1 {
		return Season{startYear: startYear, endYear: endYear}, nil
	}
	return Season{}, fmt.Errorf("invalid season years: %d-%d (expected %d-%d or %d-%d)",
		startYear, endYear, startYear, startYear, startYear, startYear+1)
}

// StartYear returns the start year of the season.
func (s Season) StartYear() int {
	return s.startYear
}

// EndYear returns the end year of the season.
func (s Season) EndYear() int {
	return s.endYear
}

// APIString converts the Season to the API format (YYYYYYYY).
// For example, the 2023-2024 season is represented as "20232024".
func (s Season) APIString() string {
	return fmt.Sprintf("%04d%04d", s.startYear, s.endYear)
}

// String implements the fmt.Stringer interface.
// Returns the season in "YYYY-YYYY" format.
func (s Season) String() string {
	return fmt.Sprintf("%d-%d", s.startYear, s.endYear)
}

// GobEncode implements gob.GobEncoder for Season.
func (s Season) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s.startYear); err != nil {
		return nil, err
	}
	if err := enc.Encode(s.endYear); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode implements gob.GobDecoder for Season.
// It tolerates legacy data that encoded only the start year (before the end
// year was stored): in that case the end year defaults to startYear + 1, the
// conventional cross-year value, preserving backward compatibility.
func (s *Season) GobDecode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&s.startYear); err != nil {
		return err
	}
	if err := dec.Decode(&s.endYear); err != nil {
		if err == io.EOF {
			s.endYear = s.startYear + 1
			return nil
		}
		return err
	}
	return nil
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

// Current returns the current NHL season based on the current date (UTC).
// The NHL season typically starts in October and ends in June. UTC is used
// for consistency with Today() and GameDate.Date(); the season-rollover months
// (June/July) fall in the offseason, so the exact time zone is not significant.
func Current() Season {
	now := time.Now().UTC()
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
	return json.Marshal(s.APIString())
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

// ID returns the Season as an integer in YYYYYYYY format.
func (s Season) ID() int {
	return s.startYear*10000 + s.endYear
}

// Int64 returns the Season as an int64 in YYYYYYYY format.
func (s Season) Int64() int64 {
	return int64(s.ID())
}
