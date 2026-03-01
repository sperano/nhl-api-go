package nhl

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	gd := Now()

	if !gd.IsNow() {
		t.Error("Now() should create a GameDate with IsNow() == true")
	}

	// Date() should return current time
	now := time.Now().UTC()
	date := gd.Date()

	// Allow small time difference (within 1 second)
	diff := date.Sub(now)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("Date() returned %v, expected close to %v", date, now)
	}
}

func TestFromDate(t *testing.T) {
	testTime := time.Date(2023, 10, 15, 12, 30, 0, 0, time.UTC)
	gd := FromDate(testTime)

	if gd.IsNow() {
		t.Error("FromDate() should create a GameDate with IsNow() == false")
	}

	if !gd.Date().Equal(testTime) {
		t.Errorf("Date() = %v, want %v", gd.Date(), testTime)
	}
}

func TestFromYMD(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month int
		day   int
	}{
		{
			name:  "2023-10-15",
			year:  2023,
			month: 10,
			day:   15,
		},
		{
			name:  "2024-01-01",
			year:  2024,
			month: 1,
			day:   1,
		},
		{
			name:  "2024-12-31",
			year:  2024,
			month: 12,
			day:   31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gd := FromYMD(tt.year, tt.month, tt.day)

			if gd.IsNow() {
				t.Error("FromYMD() should create a GameDate with IsNow() == false")
			}

			date := gd.Date()
			if date.Year() != tt.year {
				t.Errorf("Year = %d, want %d", date.Year(), tt.year)
			}

			if int(date.Month()) != tt.month {
				t.Errorf("Month = %d, want %d", date.Month(), tt.month)
			}

			if date.Day() != tt.day {
				t.Errorf("Day = %d, want %d", date.Day(), tt.day)
			}
		})
	}
}

func TestToday(t *testing.T) {
	gd := Today()

	if gd.IsNow() {
		t.Error("Today() should create a GameDate with IsNow() == false")
	}

	now := time.Now().UTC()
	date := gd.Date()

	if date.Year() != now.Year() || date.Month() != now.Month() || date.Day() != now.Day() {
		t.Errorf("Today() = %v, expected today's date", date)
	}

	// Time should be midnight UTC
	if date.Hour() != 0 || date.Minute() != 0 || date.Second() != 0 {
		t.Errorf("Today() should set time to midnight, got %02d:%02d:%02d",
			date.Hour(), date.Minute(), date.Second())
	}
}

func TestGameDate_ToAPIString(t *testing.T) {
	tests := []struct {
		name     string
		gameDate GameDate
		expected string
	}{
		{
			name:     "specific date",
			gameDate: FromYMD(2023, 10, 15),
			expected: "2023-10-15",
		},
		{
			name:     "date with single-digit month",
			gameDate: FromYMD(2023, 1, 5),
			expected: "2023-01-05",
		},
		{
			name:     "date with single-digit day",
			gameDate: FromYMD(2023, 12, 5),
			expected: "2023-12-05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.gameDate.ToAPIString()
			if result != tt.expected {
				t.Errorf("ToAPIString() = %q, want %q", result, tt.expected)
			}
		})
	}

	t.Run("now", func(t *testing.T) {
		gd := Now()
		result := gd.ToAPIString()

		// Should match current date in YYYY-MM-DD format
		now := time.Now().UTC()
		expected := now.Format("2006-01-02")

		if result != expected {
			t.Errorf("ToAPIString() for Now() = %q, want %q", result, expected)
		}
	})
}

func TestGameDate_AddDays(t *testing.T) {
	tests := []struct {
		name     string
		base     GameDate
		days     int
		expected string
	}{
		{
			name:     "add positive days",
			base:     FromYMD(2023, 10, 15),
			days:     5,
			expected: "2023-10-20",
		},
		{
			name:     "add negative days",
			base:     FromYMD(2023, 10, 15),
			days:     -5,
			expected: "2023-10-10",
		},
		{
			name:     "add zero days",
			base:     FromYMD(2023, 10, 15),
			days:     0,
			expected: "2023-10-15",
		},
		{
			name:     "cross month boundary",
			base:     FromYMD(2023, 10, 30),
			days:     5,
			expected: "2023-11-04",
		},
		{
			name:     "cross year boundary",
			base:     FromYMD(2023, 12, 30),
			days:     5,
			expected: "2024-01-04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.AddDays(tt.days)

			if result.IsNow() {
				t.Error("AddDays() should return a concrete date, not Now")
			}

			if result.ToAPIString() != tt.expected {
				t.Errorf("AddDays(%d) = %q, want %q", tt.days, result.ToAPIString(), tt.expected)
			}
		})
	}

	t.Run("add to now", func(t *testing.T) {
		gd := Now()
		result := gd.AddDays(1)

		if result.IsNow() {
			t.Error("AddDays() on Now should return a concrete date")
		}

		// Should be tomorrow's date
		tomorrow := time.Now().UTC().AddDate(0, 0, 1)
		expected := tomorrow.Format("2006-01-02")

		if result.ToAPIString() != expected {
			t.Errorf("AddDays(1) on Now() = %q, want %q", result.ToAPIString(), expected)
		}
	})
}

func TestGameDate_String(t *testing.T) {
	tests := []struct {
		name     string
		gameDate GameDate
		expected string
	}{
		{
			name:     "now",
			gameDate: Now(),
			expected: "now",
		},
		{
			name:     "specific date",
			gameDate: FromYMD(2023, 10, 15),
			expected: "2023-10-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.gameDate.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGameDate_JSON(t *testing.T) {
	t.Run("marshal now", func(t *testing.T) {
		gd := Now()
		data, err := json.Marshal(gd)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		expected := `"now"`
		if string(data) != expected {
			t.Errorf("json.Marshal() = %q, want %q", string(data), expected)
		}
	})

	t.Run("marshal date", func(t *testing.T) {
		gd := FromYMD(2023, 10, 15)
		data, err := json.Marshal(gd)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		expected := `"2023-10-15"`
		if string(data) != expected {
			t.Errorf("json.Marshal() = %q, want %q", string(data), expected)
		}
	})

	t.Run("unmarshal now", func(t *testing.T) {
		data := []byte(`"now"`)
		var gd GameDate
		if err := json.Unmarshal(data, &gd); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		if !gd.IsNow() {
			t.Error("Unmarshaled 'now' should have IsNow() == true")
		}
	})

	t.Run("unmarshal date", func(t *testing.T) {
		data := []byte(`"2023-10-15"`)
		var gd GameDate
		if err := json.Unmarshal(data, &gd); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		if gd.IsNow() {
			t.Error("Unmarshaled date should have IsNow() == false")
		}

		if gd.ToAPIString() != "2023-10-15" {
			t.Errorf("Unmarshaled date = %q, want %q", gd.ToAPIString(), "2023-10-15")
		}
	})

	t.Run("unmarshal invalid format", func(t *testing.T) {
		data := []byte(`"invalid"`)
		var gd GameDate
		if err := json.Unmarshal(data, &gd); err == nil {
			t.Error("json.Unmarshal() should return error for invalid format")
		}
	})
}

func TestNewSeason(t *testing.T) {
	season := NewSeason(2023)

	if season.StartYear() != 2023 {
		t.Errorf("StartYear() = %d, want %d", season.StartYear(), 2023)
	}

	if season.EndYear() != 2024 {
		t.Errorf("EndYear() = %d, want %d", season.EndYear(), 2024)
	}
}

func TestFromYears(t *testing.T) {
	t.Run("valid years", func(t *testing.T) {
		season, err := FromYears(2023, 2024)
		if err != nil {
			t.Fatalf("FromYears() error = %v", err)
		}

		if season.StartYear() != 2023 {
			t.Errorf("StartYear() = %d, want %d", season.StartYear(), 2023)
		}

		if season.EndYear() != 2024 {
			t.Errorf("EndYear() = %d, want %d", season.EndYear(), 2024)
		}
	})

	t.Run("invalid years", func(t *testing.T) {
		_, err := FromYears(2023, 2025)
		if err == nil {
			t.Error("FromYears() should return error for invalid year range")
		}
	})

	t.Run("reversed years", func(t *testing.T) {
		_, err := FromYears(2024, 2023)
		if err == nil {
			t.Error("FromYears() should return error for reversed years")
		}
	})
}

func TestSeason_ToAPIString(t *testing.T) {
	tests := []struct {
		name     string
		season   Season
		expected string
	}{
		{
			name:     "2023-2024 season",
			season:   NewSeason(2023),
			expected: "20232024",
		},
		{
			name:     "2024-2025 season",
			season:   NewSeason(2024),
			expected: "20242025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.season.ToAPIString()
			if result != tt.expected {
				t.Errorf("ToAPIString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSeason_String(t *testing.T) {
	tests := []struct {
		name     string
		season   Season
		expected string
	}{
		{
			name:     "2023-2024 season",
			season:   NewSeason(2023),
			expected: "2023-2024",
		},
		{
			name:     "2024-2025 season",
			season:   NewSeason(2024),
			expected: "2024-2025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.season.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{
			name:      "dash format",
			input:     "2023-2024",
			wantStart: 2023,
			wantEnd:   2024,
			wantErr:   false,
		},
		{
			name:      "concatenated format",
			input:     "20232024",
			wantStart: 2023,
			wantEnd:   2024,
			wantErr:   false,
		},
		{
			name:    "invalid dash format",
			input:   "2023-2025",
			wantErr: true,
		},
		{
			name:    "invalid concatenated format",
			input:   "20232025",
			wantErr: true,
		},
		{
			name:    "invalid format too short",
			input:   "2023",
			wantErr: true,
		},
		{
			name:    "invalid format empty",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			season, err := Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("Parse() should return error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if season.StartYear() != tt.wantStart {
				t.Errorf("StartYear() = %d, want %d", season.StartYear(), tt.wantStart)
			}

			if season.EndYear() != tt.wantEnd {
				t.Errorf("EndYear() = %d, want %d", season.EndYear(), tt.wantEnd)
			}
		})
	}
}

func TestCurrent(t *testing.T) {
	season := Current()
	now := time.Now()
	year := now.Year()
	month := now.Month()

	var expectedStart int
	if month >= time.January && month <= time.June {
		expectedStart = year - 1
	} else {
		expectedStart = year
	}

	if season.StartYear() != expectedStart {
		t.Errorf("Current() StartYear() = %d, want %d", season.StartYear(), expectedStart)
	}
}

func TestSeason_JSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		season := NewSeason(2023)
		data, err := json.Marshal(season)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		expected := `"20232024"`
		if string(data) != expected {
			t.Errorf("json.Marshal() = %q, want %q", string(data), expected)
		}
	})

	t.Run("unmarshal concatenated", func(t *testing.T) {
		data := []byte(`"20232024"`)
		var season Season
		if err := json.Unmarshal(data, &season); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		if season.StartYear() != 2023 {
			t.Errorf("StartYear() = %d, want %d", season.StartYear(), 2023)
		}

		if season.EndYear() != 2024 {
			t.Errorf("EndYear() = %d, want %d", season.EndYear(), 2024)
		}
	})

	t.Run("unmarshal dash format", func(t *testing.T) {
		data := []byte(`"2023-2024"`)
		var season Season
		if err := json.Unmarshal(data, &season); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}

		if season.StartYear() != 2023 {
			t.Errorf("StartYear() = %d, want %d", season.StartYear(), 2023)
		}
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		data := []byte(`"invalid"`)
		var season Season
		if err := json.Unmarshal(data, &season); err == nil {
			t.Error("json.Unmarshal() should return error for invalid format")
		}
	})
}

func TestSeasonFromInt(t *testing.T) {
	tests := []struct {
		name      string
		input     int
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{
			name:      "valid season",
			input:     20232024,
			wantStart: 2023,
			wantEnd:   2024,
			wantErr:   false,
		},
		{
			name:    "too small",
			input:   2023,
			wantErr: true,
		},
		{
			name:    "too large",
			input:   202320241,
			wantErr: true,
		},
		{
			name:    "invalid year range",
			input:   20232025,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			season, err := SeasonFromInt(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("SeasonFromInt() should return error")
				}
				return
			}

			if err != nil {
				t.Fatalf("SeasonFromInt() error = %v", err)
			}

			if season.StartYear() != tt.wantStart {
				t.Errorf("StartYear() = %d, want %d", season.StartYear(), tt.wantStart)
			}

			if season.EndYear() != tt.wantEnd {
				t.Errorf("EndYear() = %d, want %d", season.EndYear(), tt.wantEnd)
			}
		})
	}
}

func TestSeasonFromInt64(t *testing.T) {
	season, err := SeasonFromInt64(20232024)
	if err != nil {
		t.Fatalf("SeasonFromInt64() error = %v", err)
	}

	if season.StartYear() != 2023 {
		t.Errorf("StartYear() = %d, want %d", season.StartYear(), 2023)
	}

	if season.EndYear() != 2024 {
		t.Errorf("EndYear() = %d, want %d", season.EndYear(), 2024)
	}
}

func TestSeason_ToInt(t *testing.T) {
	tests := []struct {
		startYear int
		expected  int
	}{
		{2024, 20242025},
		{2023, 20232024},
		{2000, 20002001},
		{1999, 19992000}, // Y2K boundary
		{1917, 19171918}, // First NHL season
	}

	for _, tc := range tests {
		season := NewSeason(tc.startYear)
		result := season.ToInt()
		if result != tc.expected {
			t.Errorf("NewSeason(%d).ToInt() = %d, want %d", tc.startYear, result, tc.expected)
		}
	}
}

func TestSeason_ID(t *testing.T) {
	tests := []struct {
		startYear int
		expected  int
	}{
		{2024, 20242025},
		{2023, 20232024},
		{2000, 20002001},
		{1999, 19992000}, // Y2K boundary
		{1917, 19171918}, // First NHL season
	}

	for _, tc := range tests {
		season := NewSeason(tc.startYear)
		result := season.ID()
		if result != tc.expected {
			t.Errorf("NewSeason(%d).ID() = %d, want %d", tc.startYear, result, tc.expected)
		}
	}
}

func TestSeason_ToInt64(t *testing.T) {
	season := NewSeason(2023)
	result := season.ToInt64()

	var expected int64 = 20232024
	if result != expected {
		t.Errorf("ToInt64() = %d, want %d", result, expected)
	}
}

func TestGameDate_EdgeCases(t *testing.T) {
	t.Run("leap year date", func(t *testing.T) {
		gd := FromYMD(2024, 2, 29)
		if gd.ToAPIString() != "2024-02-29" {
			t.Errorf("Leap year date failed: %s", gd.ToAPIString())
		}
	})

	t.Run("add days across leap year", func(t *testing.T) {
		gd := FromYMD(2024, 2, 28)
		next := gd.AddDays(1)
		if next.ToAPIString() != "2024-02-29" {
			t.Errorf("AddDays across leap year failed: %s", next.ToAPIString())
		}
	})
}

func TestSeason_EdgeCases(t *testing.T) {
	t.Run("century year", func(t *testing.T) {
		season := NewSeason(2000)
		if season.String() != "2000-2001" {
			t.Errorf("Century year season failed: %s", season.String())
		}
	})

	t.Run("parse with whitespace", func(t *testing.T) {
		input := strings.TrimSpace("  2023-2024  ")
		_, err := Parse(input)
		if err != nil {
			t.Errorf("Parse should handle trimmed whitespace, got error: %v", err)
		}
	})
}

// Additional error path tests for date.go

func TestSeason_Parse_ErrorPaths(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"invalid format - too many parts", "2023-2024-2025", true},
		{"invalid start year in dash format", "ABCD-2024", true},
		{"invalid end year in dash format", "2023-ABCD", true},
		{"invalid start year in 8-digit format", "ABCD2024", true},
		{"invalid end year in 8-digit format", "2023ABCD", true},
		{"invalid format - too short", "202324", true},
		{"invalid format - random string", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSeasonFromInt64_ErrorPath(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		wantErr bool
	}{
		{"too small", 9999999, true},
		{"too large", 100000000, true},
		{"negative", -20232024, true},
		{"valid", 20232024, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SeasonFromInt64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SeasonFromInt64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGameDate_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var gd GameDate
	err := json.Unmarshal([]byte(`{invalid json}`), &gd)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestSeason_UnmarshalJSON_InvalidJSON(t *testing.T) {
	var s Season
	err := json.Unmarshal([]byte(`{invalid json}`), &s)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid JSON")
	}
}

func TestSeason_Current_Coverage(t *testing.T) {
	// This test just ensures Current() is called and doesn't panic
	season := Current()
	if season.StartYear() < 1900 || season.StartYear() > 2100 {
		t.Errorf("Current() returned unreasonable year: %d", season.StartYear())
	}
}

// Additional error path tests for GameDate.UnmarshalJSON

func TestGameDate_UnmarshalJSON_InvalidDateFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"too few parts", `"2024-01"`},
		{"too many parts", `"2024-01-15-extra"`},
		{"invalid year", `"ABCD-01-15"`},
		{"invalid month", `"2024-AB-15"`},
		{"invalid day", `"2024-01-AB"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gd GameDate
			err := json.Unmarshal([]byte(tt.input), &gd)
			if err == nil {
				t.Error("UnmarshalJSON() should error on invalid date format")
			}
		})
	}
}

func TestGameDate_UnmarshalJSON_Now(t *testing.T) {
	var gd GameDate
	err := json.Unmarshal([]byte(`"now"`), &gd)
	if err != nil {
		t.Errorf("UnmarshalJSON() error = %v", err)
	}
	if !gd.IsNow() {
		t.Error("UnmarshalJSON() with 'now' should create a GameDate with IsNow() = true")
	}
}

func TestSeason_UnmarshalJSON_InvalidFormat(t *testing.T) {
	var s Season
	err := json.Unmarshal([]byte(`"invalid-format"`), &s)
	if err == nil {
		t.Error("UnmarshalJSON() should error on invalid season format")
	}
}

func TestSeason_Gob(t *testing.T) {
	original := NewSeason(2023)

	// Encode
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(original); err != nil {
		t.Fatalf("GobEncode() error = %v", err)
	}

	// Decode
	var decoded Season
	dec := gob.NewDecoder(&buf)
	if err := dec.Decode(&decoded); err != nil {
		t.Fatalf("GobDecode() error = %v", err)
	}

	if decoded.StartYear() != original.StartYear() {
		t.Errorf("GobDecode() StartYear = %d, want %d", decoded.StartYear(), original.StartYear())
	}
}

func TestSeason_GobInStruct(t *testing.T) {
	// Test that Season can be gob-encoded when embedded in another struct
	type Container struct {
		Name   string
		Season Season
		Value  int
	}

	original := Container{
		Name:   "test",
		Season: NewSeason(2024),
		Value:  42,
	}

	// Encode
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(original); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	// Decode
	var decoded Container
	dec := gob.NewDecoder(&buf)
	if err := dec.Decode(&decoded); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if decoded.Name != original.Name {
		t.Errorf("Name = %s, want %s", decoded.Name, original.Name)
	}
	if decoded.Season.StartYear() != original.Season.StartYear() {
		t.Errorf("Season.StartYear() = %d, want %d", decoded.Season.StartYear(), original.Season.StartYear())
	}
	if decoded.Value != original.Value {
		t.Errorf("Value = %d, want %d", decoded.Value, original.Value)
	}
}

// Benchmarks comparing ToInt() (string-based) vs ID() (arithmetic-based)

func BenchmarkSeason_ToInt(b *testing.B) {
	season := NewSeason(2023)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = season.ToInt()
	}
}

func BenchmarkSeason_ID(b *testing.B) {
	season := NewSeason(2023)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = season.ID()
	}
}
