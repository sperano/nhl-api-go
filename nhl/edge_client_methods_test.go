package nhl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// edgeClientCase enumerates one Client.Edge* method: its name, the URL path
// the client is expected to GET, and a closure that calls it with shared
// fixture IDs. Two tests share this table so the URL-path contract and the
// error-propagation contract stay in lockstep across all 22 Edge methods.
type edgeClientCase struct {
	name string
	path string
	do   func(*Client) (any, error)
}

// Shared fixtures: 2024-2025 regular season → APIString "20242025", gameType.Int() 2.
const (
	edgeTestPlayerIDStr = "8478402"
	edgeTestTeamIDStr   = "22"
	edgeTestSeasonStr   = "20242025"
	edgeTestGameTypeStr = "2"
)

var (
	edgeTestPlayerID = PlayerID(8478402)
	edgeTestTeamID   = TeamID(22)
	edgeTestSeason   = NewSeason(2024)
	edgeTestGameType = GameTypeRegularSeason
)

func edgeClientCases() []edgeClientCase {
	playerPath := func(slug string) string {
		return "/edge/" + slug + "/" + edgeTestPlayerIDStr + "/" + edgeTestSeasonStr + "/" + edgeTestGameTypeStr
	}
	teamPath := func(slug string) string {
		return "/edge/" + slug + "/" + edgeTestTeamIDStr + "/" + edgeTestSeasonStr + "/" + edgeTestGameTypeStr
	}
	landingPath := func(slug string) string {
		return "/edge/" + slug + "/" + edgeTestSeasonStr + "/" + edgeTestGameTypeStr
	}

	ctx := context.Background()
	return []edgeClientCase{
		// --- Skater per-player ---
		{"EdgeSkaterDetail", playerPath("skater-detail"),
			func(c *Client) (any, error) {
				return c.EdgeSkaterDetail(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeSkaterSpeedDetail", playerPath("skater-skating-speed-detail"),
			func(c *Client) (any, error) {
				return c.EdgeSkaterSpeedDetail(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeSkaterDistanceDetail", playerPath("skater-skating-distance-detail"),
			func(c *Client) (any, error) {
				return c.EdgeSkaterDistanceDetail(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeSkaterShotSpeedDetail", playerPath("skater-shot-speed-detail"),
			func(c *Client) (any, error) {
				return c.EdgeSkaterShotSpeedDetail(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeSkaterShotLocationDetail", playerPath("skater-shot-location-detail"),
			func(c *Client) (any, error) {
				return c.EdgeSkaterShotLocationDetail(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeSkaterZoneTime", playerPath("skater-zone-time"),
			func(c *Client) (any, error) {
				return c.EdgeSkaterZoneTime(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeSkaterComparison", playerPath("skater-comparison"),
			func(c *Client) (any, error) {
				return c.EdgeSkaterComparison(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},

		// --- Goalie per-player ---
		{"EdgeGoalieDetail", playerPath("goalie-detail"),
			func(c *Client) (any, error) {
				return c.EdgeGoalieDetail(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeGoalie5v5Detail", playerPath("goalie-5v5-detail"),
			func(c *Client) (any, error) {
				return c.EdgeGoalie5v5Detail(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeGoalieShotLocationDetail", playerPath("goalie-shot-location-detail"),
			func(c *Client) (any, error) {
				return c.EdgeGoalieShotLocationDetail(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeGoalieSavePctgDetail", playerPath("goalie-save-percentage-detail"),
			func(c *Client) (any, error) {
				return c.EdgeGoalieSavePctgDetail(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeGoalieComparison", playerPath("goalie-comparison"),
			func(c *Client) (any, error) {
				return c.EdgeGoalieComparison(ctx, edgeTestPlayerID, edgeTestSeason, edgeTestGameType)
			}},

		// --- Team per-team ---
		{"EdgeTeamDetail", teamPath("team-detail"),
			func(c *Client) (any, error) {
				return c.EdgeTeamDetail(ctx, edgeTestTeamID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeTeamSpeedDetail", teamPath("team-skating-speed-detail"),
			func(c *Client) (any, error) {
				return c.EdgeTeamSpeedDetail(ctx, edgeTestTeamID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeTeamDistanceDetail", teamPath("team-skating-distance-detail"),
			func(c *Client) (any, error) {
				return c.EdgeTeamDistanceDetail(ctx, edgeTestTeamID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeTeamShotSpeedDetail", teamPath("team-shot-speed-detail"),
			func(c *Client) (any, error) {
				return c.EdgeTeamShotSpeedDetail(ctx, edgeTestTeamID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeTeamShotLocationDetail", teamPath("team-shot-location-detail"),
			func(c *Client) (any, error) {
				return c.EdgeTeamShotLocationDetail(ctx, edgeTestTeamID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeTeamZoneTimeDetails", teamPath("team-zone-time-details"),
			func(c *Client) (any, error) {
				return c.EdgeTeamZoneTimeDetails(ctx, edgeTestTeamID, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeTeamComparison", teamPath("team-comparison"),
			func(c *Client) (any, error) {
				return c.EdgeTeamComparison(ctx, edgeTestTeamID, edgeTestSeason, edgeTestGameType)
			}},

		// --- League-wide landings ---
		{"EdgeSkaterLanding", landingPath("skater-landing"),
			func(c *Client) (any, error) {
				return c.EdgeSkaterLanding(ctx, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeGoalieLanding", landingPath("goalie-landing"),
			func(c *Client) (any, error) {
				return c.EdgeGoalieLanding(ctx, edgeTestSeason, edgeTestGameType)
			}},
		{"EdgeTeamLanding", landingPath("team-landing"),
			func(c *Client) (any, error) {
				return c.EdgeTeamLanding(ctx, edgeTestSeason, edgeTestGameType)
			}},
	}
}

// TestEdge_AllClientMethods_PathContract verifies every Edge client method
// requests the documented URL path. Each method is run against a stub server
// that records the path and returns an empty JSON object — empty {} unmarshals
// cleanly into any of the response struct types (all fields stay zero-valued).
//
// This catches URL-formatter regressions (typos, swapped IDs, wrong endpoint
// suffix). For field-level deserialization the per-type *_Deserialization and
// *_RealAPIStructure tests in edge_test.go are authoritative.
func TestEdge_AllClientMethods_PathContract(t *testing.T) {
	for _, tc := range edgeClientCases() {
		t.Run(tc.name, func(t *testing.T) {
			var seenPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seenPath = r.URL.Path
				fmt.Fprint(w, `{}`)
			}))
			defer server.Close()

			client := NewClientWithBaseURL(server.URL)
			if _, err := tc.do(client); err != nil {
				t.Fatalf("%s returned error: %v", tc.name, err)
			}
			if seenPath != tc.path {
				t.Errorf("%s requested path %q, want %q", tc.name, seenPath, tc.path)
			}
		})
	}
}

// TestEdge_AllClientMethods_404PropagatesError verifies every Edge client
// method propagates a 404 from getJSON as ErrNotFound. Without this the error
// branch ("if err := getJSON; err != nil { return nil, err }") goes
// uncovered in 19 of 22 methods — making it cheap to silently swallow errors
// in a future refactor that adds, say, an err == ErrNotFound special-case.
func TestEdge_AllClientMethods_404PropagatesError(t *testing.T) {
	for _, tc := range edgeClientCases() {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()

			client := NewClientWithBaseURL(server.URL)
			_, err := tc.do(client)
			if err == nil {
				t.Fatalf("%s expected error for 404, got nil", tc.name)
			}
			if !errors.Is(err, ErrNotFound) {
				t.Errorf("%s expected ErrNotFound, got %T: %v", tc.name, err, err)
			}
		})
	}
}
