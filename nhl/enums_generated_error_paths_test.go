package nhl

import (
	"encoding/json"
	"testing"
)

// These tests target generator-template surfaces in enums_generated.go that
// the per-enum tests in enums_test.go don't otherwise exercise: the inner
// json.Unmarshal error branch, the empty-string branches (where the
// template emits them), and a handful of constants whose individual cases
// are otherwise unreached. Grouped here rather than scattered into the
// per-enum sections to make the intent — "lock down the generated
// template's behavior" — visible at a glance.

// nonStringJSONInput is syntactically valid JSON of the wrong shape for
// any enum: the outer json package hands it to our UnmarshalJSON method
// as-is, and the inner json.Unmarshal(data, &s) then fails because a
// number can't unmarshal into a string. Malformed JSON like "{not json}"
// would short-circuit at the outer parser and never reach our method, so
// it cannot exercise this branch.
const nonStringJSONInput = "123"

func TestEnumsGenerated_NonStringJSONRejected(t *testing.T) {
	cases := []struct {
		name   string
		decode func() error
	}{
		{"Position", func() error { var v Position; return json.Unmarshal([]byte(nonStringJSONInput), &v) }},
		{"Handedness", func() error { var v Handedness; return json.Unmarshal([]byte(nonStringJSONInput), &v) }},
		{"GoalieDecision", func() error { var v GoalieDecision; return json.Unmarshal([]byte(nonStringJSONInput), &v) }},
		{"PeriodType", func() error { var v PeriodType; return json.Unmarshal([]byte(nonStringJSONInput), &v) }},
		{"HomeRoad", func() error { var v HomeRoad; return json.Unmarshal([]byte(nonStringJSONInput), &v) }},
		{"ZoneCode", func() error { var v ZoneCode; return json.Unmarshal([]byte(nonStringJSONInput), &v) }},
		{"DefendingSide", func() error { var v DefendingSide; return json.Unmarshal([]byte(nonStringJSONInput), &v) }},
		{"GameScheduleState", func() error { var v GameScheduleState; return json.Unmarshal([]byte(nonStringJSONInput), &v) }},
		{"PlayEventType", func() error { var v PlayEventType; return json.Unmarshal([]byte(nonStringJSONInput), &v) }},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.decode(); err == nil {
				t.Errorf("%s.UnmarshalJSON(%s) should reject non-string JSON", c.name, nonStringJSONInput)
			}
		})
	}
}

// TestEnumsGenerated_EmptyStringAccepted verifies the empty-string branch
// in UnmarshalJSON for the three enums whose template emits it. Empty is
// treated as a valid zero value (e.g. an absent player position field in
// the API), distinct from an invalid string which must still return an
// error. The other six enums in enums_generated.go don't have this
// branch and would reject empty-string input.
func TestEnumsGenerated_EmptyStringAccepted(t *testing.T) {
	t.Run("Position", func(t *testing.T) {
		var v Position
		if err := json.Unmarshal([]byte(`""`), &v); err != nil {
			t.Fatalf("UnmarshalJSON(\"\") error = %v", err)
		}
		if v != "" {
			t.Errorf("UnmarshalJSON(\"\") = %q, want \"\"", v)
		}
	})

	t.Run("Handedness", func(t *testing.T) {
		var v Handedness
		if err := json.Unmarshal([]byte(`""`), &v); err != nil {
			t.Fatalf("UnmarshalJSON(\"\") error = %v", err)
		}
		if v != "" {
			t.Errorf("UnmarshalJSON(\"\") = %q, want \"\"", v)
		}
	})

	t.Run("PeriodType", func(t *testing.T) {
		var v PeriodType
		if err := json.Unmarshal([]byte(`""`), &v); err != nil {
			t.Fatalf("UnmarshalJSON(\"\") error = %v", err)
		}
		if v != "" {
			t.Errorf("UnmarshalJSON(\"\") = %q, want \"\"", v)
		}
	})
}

// TestPosition_MarshalJSON_Empty pins the "" -> json.Marshal("") branch in
// Position.MarshalJSON. The other valid-position cases are already
// covered by TestPosition_JSON in enums_test.go.
func TestPosition_MarshalJSON_Empty(t *testing.T) {
	got, err := json.Marshal(Position(""))
	if err != nil {
		t.Fatalf("Marshal(empty Position) error = %v", err)
	}
	if string(got) != `""` {
		t.Errorf("Marshal(empty Position) = %s, want \"\"", string(got))
	}
}

// TestPosition_Forward_NameAndFromString covers PositionForward, the
// only constant in the file whose Name() case and FromString() case had
// no test. PositionForward exists for historical data ("F" was the
// generic forward position before C/LW/RW were broken out).
func TestPosition_Forward_NameAndFromString(t *testing.T) {
	t.Run("Name", func(t *testing.T) {
		if got := PositionForward.Name(); got != "Forward" {
			t.Errorf("PositionForward.Name() = %q, want %q", got, "Forward")
		}
	})

	t.Run("FromString F", func(t *testing.T) {
		got, err := PositionFromString("F")
		if err != nil {
			t.Fatalf("PositionFromString(\"F\") error = %v", err)
		}
		if got != PositionForward {
			t.Errorf("PositionFromString(\"F\") = %v, want PositionForward", got)
		}
	})

	t.Run("FromString Forward", func(t *testing.T) {
		got, err := PositionFromString("Forward")
		if err != nil {
			t.Fatalf("PositionFromString(\"Forward\") error = %v", err)
		}
		if got != PositionForward {
			t.Errorf("PositionFromString(\"Forward\") = %v, want PositionForward", got)
		}
	})
}
