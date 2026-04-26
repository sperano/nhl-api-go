package nhl

import (
	"encoding/json"
	"testing"
)

// These tests cover PlayerID and TeamID, the two ID wrappers produced from
// the same template as GameID (see internal/idgen). GameID has its own
// dedicated suite in game_id_test.go; this file only proves the generator
// surface for PlayerID and TeamID, since the shared helpers in id.go are
// already exercised through GameID.

const (
	samplePlayerID = int64(8478402) // Connor McDavid
	sampleTeamID   = int64(10)      // Toronto Maple Leafs
)

func TestPlayerID_GeneratorSurface(t *testing.T) {
	t.Run("New and Int64", func(t *testing.T) {
		id := NewPlayerID(samplePlayerID)
		if id.Int64() != samplePlayerID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), samplePlayerID)
		}
	})

	t.Run("String", func(t *testing.T) {
		got := PlayerID(samplePlayerID).String()
		if got != "8478402" {
			t.Errorf("String() = %q, want %q", got, "8478402")
		}
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		data, err := json.Marshal(PlayerID(samplePlayerID))
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if string(data) != "8478402" {
			t.Errorf("Marshal() = %q, want %q", string(data), "8478402")
		}
	})

	t.Run("UnmarshalJSON int", func(t *testing.T) {
		var id PlayerID
		if err := json.Unmarshal([]byte("8478402"), &id); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if id.Int64() != samplePlayerID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), samplePlayerID)
		}
	})

	t.Run("UnmarshalJSON string", func(t *testing.T) {
		var id PlayerID
		if err := json.Unmarshal([]byte(`"8478402"`), &id); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if id.Int64() != samplePlayerID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), samplePlayerID)
		}
	})

	t.Run("UnmarshalJSON invalid", func(t *testing.T) {
		var id PlayerID
		if err := json.Unmarshal([]byte("true"), &id); err == nil {
			t.Error("Unmarshal() should reject boolean JSON")
		}
	})

	t.Run("FromInt", func(t *testing.T) {
		id := PlayerIDFromInt(int(samplePlayerID))
		if id.Int64() != samplePlayerID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), samplePlayerID)
		}
	})

	t.Run("FromString valid", func(t *testing.T) {
		id, err := PlayerIDFromString("8478402")
		if err != nil {
			t.Fatalf("FromString() error = %v", err)
		}
		if id.Int64() != samplePlayerID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), samplePlayerID)
		}
	})

	t.Run("FromString invalid", func(t *testing.T) {
		if _, err := PlayerIDFromString("notanumber"); err == nil {
			t.Error("FromString() should reject non-numeric input")
		}
	})

	t.Run("MustFromString valid", func(t *testing.T) {
		id := MustPlayerIDFromString("8478402")
		if id.Int64() != samplePlayerID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), samplePlayerID)
		}
	})

	t.Run("MustFromString panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustPlayerIDFromString() did not panic on invalid input")
			}
		}()
		MustPlayerIDFromString("notanumber")
	})
}

func TestTeamID_GeneratorSurface(t *testing.T) {
	t.Run("New and Int64", func(t *testing.T) {
		id := NewTeamID(sampleTeamID)
		if id.Int64() != sampleTeamID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), sampleTeamID)
		}
	})

	t.Run("String", func(t *testing.T) {
		got := TeamID(sampleTeamID).String()
		if got != "10" {
			t.Errorf("String() = %q, want %q", got, "10")
		}
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		data, err := json.Marshal(TeamID(sampleTeamID))
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if string(data) != "10" {
			t.Errorf("Marshal() = %q, want %q", string(data), "10")
		}
	})

	t.Run("UnmarshalJSON int", func(t *testing.T) {
		var id TeamID
		if err := json.Unmarshal([]byte("10"), &id); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if id.Int64() != sampleTeamID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), sampleTeamID)
		}
	})

	t.Run("UnmarshalJSON string", func(t *testing.T) {
		var id TeamID
		if err := json.Unmarshal([]byte(`"10"`), &id); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if id.Int64() != sampleTeamID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), sampleTeamID)
		}
	})

	t.Run("UnmarshalJSON invalid", func(t *testing.T) {
		var id TeamID
		if err := json.Unmarshal([]byte("true"), &id); err == nil {
			t.Error("Unmarshal() should reject boolean JSON")
		}
	})

	t.Run("FromInt", func(t *testing.T) {
		id := TeamIDFromInt(int(sampleTeamID))
		if id.Int64() != sampleTeamID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), sampleTeamID)
		}
	})

	t.Run("FromString valid", func(t *testing.T) {
		id, err := TeamIDFromString("10")
		if err != nil {
			t.Fatalf("FromString() error = %v", err)
		}
		if id.Int64() != sampleTeamID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), sampleTeamID)
		}
	})

	t.Run("FromString invalid", func(t *testing.T) {
		if _, err := TeamIDFromString("notanumber"); err == nil {
			t.Error("FromString() should reject non-numeric input")
		}
	})

	t.Run("MustFromString valid", func(t *testing.T) {
		id := MustTeamIDFromString("10")
		if id.Int64() != sampleTeamID {
			t.Errorf("Int64() = %d, want %d", id.Int64(), sampleTeamID)
		}
	})

	t.Run("MustFromString panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustTeamIDFromString() did not panic on invalid input")
			}
		}()
		MustTeamIDFromString("notanumber")
	})
}
