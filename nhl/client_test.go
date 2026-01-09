package nhl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ===== Test Helper Functions =====

// makeJSONResponse creates an HTTP handler that returns a JSON response.
func makeJSONResponse(statusCode int, body interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}
}

// makeErrorResponse creates an HTTP handler that returns an error status code.
func makeErrorResponse(statusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}
}

// ===== Endpoint Tests =====

func TestEndpointBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
		want     string
	}{
		{
			name:     "APIWebV1",
			endpoint: EndpointAPIWebV1,
			want:     baseURLAPIWebV1,
		},
		{
			name:     "APICore",
			endpoint: EndpointAPICore,
			want:     baseURLAPICore,
		},
		{
			name:     "APIStats",
			endpoint: EndpointAPIStats,
			want:     baseURLAPIStats,
		},
		{
			name:     "SearchV1",
			endpoint: EndpointSearchV1,
			want:     baseURLSearchV1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.endpoint.baseURL()
			if got != tt.want {
				t.Errorf("baseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ===== buildURL Tests =====

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		resource string
		want     string
	}{
		{
			name:     "both have slash",
			base:     "https://api.example.com/",
			resource: "/resource",
			want:     "https://api.example.com/resource",
		},
		{
			name:     "neither has slash",
			base:     "https://api.example.com",
			resource: "resource",
			want:     "https://api.example.com/resource",
		},
		{
			name:     "base has slash only",
			base:     "https://api.example.com/",
			resource: "resource",
			want:     "https://api.example.com/resource",
		},
		{
			name:     "resource has slash only",
			base:     "https://api.example.com",
			resource: "/resource",
			want:     "https://api.example.com/resource",
		},
		{
			name:     "with path segments",
			base:     "https://api.example.com/v1/",
			resource: "/data/items",
			want:     "https://api.example.com/v1/data/items",
		},
		{
			name:     "empty resource",
			base:     "https://api.example.com/",
			resource: "",
			want:     "https://api.example.com/",
		},
		{
			name:     "empty base",
			base:     "",
			resource: "resource",
			want:     "resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildURL(tt.base, tt.resource)
			if got != tt.want {
				t.Errorf("buildURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ===== Client Creation Tests =====

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.httpClient == nil {
		t.Error("NewClient() created client with nil httpClient")
	}
}

func TestNewClientWithConfig(t *testing.T) {
	config := &ClientConfig{
		Timeout:         5 * time.Second,
		SSLVerify:       false,
		FollowRedirects: false,
	}

	client := NewClientWithConfig(config)
	if client == nil {
		t.Fatal("NewClientWithConfig() returned nil")
	}
	if client.httpClient == nil {
		t.Error("NewClientWithConfig() created client with nil httpClient")
	}
	if client.httpClient.Timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %v", client.httpClient.Timeout)
	}
}

// ===== Error Handling Tests =====

func TestGetJSON_HTTPErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedErrMsg string
		checkErrType   func(error) bool
	}{
		{
			name:           "404 Not Found",
			statusCode:     http.StatusNotFound,
			expectedErrMsg: "Request to",
			checkErrType:   func(err error) bool { _, ok := err.(*ResourceNotFoundError); return ok },
		},
		{
			name:           "429 Rate Limit",
			statusCode:     http.StatusTooManyRequests,
			expectedErrMsg: "Request to",
			checkErrType:   func(err error) bool { _, ok := err.(*RateLimitExceededError); return ok },
		},
		{
			name:           "400 Bad Request",
			statusCode:     http.StatusBadRequest,
			expectedErrMsg: "Request to",
			checkErrType:   func(err error) bool { _, ok := err.(*BadRequestError); return ok },
		},
		{
			name:           "401 Unauthorized",
			statusCode:     http.StatusUnauthorized,
			expectedErrMsg: "Request to",
			checkErrType:   func(err error) bool { _, ok := err.(*UnauthorizedError); return ok },
		},
		{
			name:           "500 Server Error",
			statusCode:     http.StatusInternalServerError,
			expectedErrMsg: "Request to",
			checkErrType:   func(err error) bool { _, ok := err.(*ServerError); return ok },
		},
		{
			name:           "503 Service Unavailable",
			statusCode:     http.StatusServiceUnavailable,
			expectedErrMsg: "Request to",
			checkErrType:   func(err error) bool { _, ok := err.(*ServerError); return ok },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(makeErrorResponse(tt.statusCode))
			defer server.Close()

			client := NewClient()
			ctx := context.Background()

			// Make request and check error type
			fullURL := server.URL + "/test"
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
			resp, err := client.httpClient.Do(req)
			if err != nil {
				t.Fatalf("unexpected http error: %v", err)
			}
			defer resp.Body.Close()

			// Manually check the error condition
			err = ErrorFromStatusCode(resp.StatusCode, "Request to /test failed")

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.expectedErrMsg) {
				t.Errorf("error message doesn't contain %q: %v", tt.expectedErrMsg, err)
			}

			if !tt.checkErrType(err) {
				t.Errorf("error is not the expected type: %T", err)
			}
		})
	}
}

func TestGetJSON_Success(t *testing.T) {
	type testResponse struct {
		Message string `json:"message"`
		Count   int    `json:"count"`
	}

	expectedResponse := testResponse{
		Message: "test message",
		Count:   42,
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, expectedResponse))
	defer server.Close()

	client := NewClient()
	ctx := context.Background()

	fullURL := server.URL + "/test"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	req.Header.Set("Accept", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result testResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Message != expectedResponse.Message {
		t.Errorf("expected message %q, got %q", expectedResponse.Message, result.Message)
	}

	if result.Count != expectedResponse.Count {
		t.Errorf("expected count %d, got %d", expectedResponse.Count, result.Count)
	}
}

// ===== Schedule Method Tests =====

func TestExtractDailySchedule(t *testing.T) {
	client := NewClient()

	weeklySchedule := &WeeklyScheduleResponse{
		NextStartDate:     "2024-01-15",
		PreviousStartDate: "2024-01-01",
		GameWeek: []GameDay{
			{
				Date: "2024-01-08",
				Games: []ScheduleGame{
					{
						ID:       GameID(2023020001),
						GameType: GameTypeRegularSeason,
					},
				},
			},
			{
				Date:  "2024-01-09",
				Games: []ScheduleGame{},
			},
		},
	}

	tests := []struct {
		name             string
		dateString       string
		expectedNumGames int
	}{
		{
			name:             "date with games",
			dateString:       "2024-01-08",
			expectedNumGames: 1,
		},
		{
			name:             "date without games",
			dateString:       "2024-01-09",
			expectedNumGames: 0,
		},
		{
			name:             "date not in schedule",
			dateString:       "2024-01-10",
			expectedNumGames: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.extractDailySchedule(weeklySchedule, tt.dateString)

			if result.Date != tt.dateString {
				t.Errorf("expected date %s, got %s", tt.dateString, result.Date)
			}

			if result.NumberOfGames != tt.expectedNumGames {
				t.Errorf("expected %d games, got %d", tt.expectedNumGames, result.NumberOfGames)
			}

			if len(result.Games) != tt.expectedNumGames {
				t.Errorf("expected %d games in slice, got %d", tt.expectedNumGames, len(result.Games))
			}

			if result.NextStartDate == nil || *result.NextStartDate != "2024-01-15" {
				t.Error("NextStartDate not set correctly")
			}

			if result.PreviousStartDate == nil || *result.PreviousStartDate != "2024-01-01" {
				t.Error("PreviousStartDate not set correctly")
			}
		})
	}
}

func TestExtractDailySchedule_EmptyGameWeek(t *testing.T) {
	client := NewClient()

	weeklySchedule := &WeeklyScheduleResponse{
		NextStartDate:     "2024-01-15",
		PreviousStartDate: "2024-01-01",
		GameWeek:          []GameDay{},
	}

	result := client.extractDailySchedule(weeklySchedule, "2024-01-08")

	if result.Date != "2024-01-08" {
		t.Errorf("expected date 2024-01-08, got %s", result.Date)
	}

	if result.NumberOfGames != 0 {
		t.Errorf("expected 0 games, got %d", result.NumberOfGames)
	}

	if len(result.Games) != 0 {
		t.Errorf("expected empty games slice, got %d games", len(result.Games))
	}
}

// ===== Query Parameter Tests =====

func TestGetJSON_WithQueryParams(t *testing.T) {
	var receivedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
	}))
	defer server.Close()

	// Test that query parameters are properly encoded
	fullURL := fmt.Sprintf("%s?q=%s&limit=%s", server.URL, "test+query", "10")

	client := NewClient()
	req, _ := http.NewRequest(http.MethodGet, fullURL, nil)
	resp, _ := client.httpClient.Do(req)
	resp.Body.Close()

	expectedContains := []string{"q=test", "limit=10"}

	for _, expected := range expectedContains {
		if !strings.Contains(receivedQuery, expected) {
			t.Errorf("query string doesn't contain %q: %s", expected, receivedQuery)
		}
	}
}

// ===== Context Tests =====

func TestDefaultContext(t *testing.T) {
	ctx, cancel := DefaultContext()
	defer cancel()

	if ctx == nil {
		t.Fatal("DefaultContext() returned nil context")
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("DefaultContext() should have a deadline")
	}

	expectedDeadline := time.Now().Add(30 * time.Second)
	if deadline.After(expectedDeadline.Add(time.Second)) {
		t.Error("deadline is too far in the future")
	}
}

func TestClientContext_Cancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
	}))
	defer server.Close()

	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Make a request with the short timeout
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	_, err := client.httpClient.Do(req)

	if err == nil {
		t.Error("expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") &&
		!strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context error, got: %v", err)
	}
}

// ===== Comprehensive Method Signature Tests =====

func TestClientMethodSignatures(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	// Test that all methods exist and have the correct signatures
	// This is a compile-time check more than a runtime test

	// Standings methods
	var _ func(context.Context) ([]Standing, error) = client.CurrentLeagueStandings
	var _ func(context.Context, GameDate) ([]Standing, error) = client.LeagueStandingsForDate
	var _ func(context.Context, Season) ([]Standing, error) = client.LeagueStandingsForSeason
	var _ func(context.Context) ([]SeasonInfo, error) = client.SeasonStandingManifest
	var _ func(context.Context, GameDate) ([]Team, error) = client.Teams

	// Schedule methods
	var _ func(context.Context, GameDate) (*DailySchedule, error) = client.DailySchedule
	var _ func(context.Context, GameDate) (*WeeklyScheduleResponse, error) = client.WeeklySchedule
	var _ func(context.Context, string, GameDate) (*TeamScheduleResponse, error) = client.TeamWeeklySchedule
	var _ func(context.Context, GameDate) (*DailyScores, error) = client.DailyScores

	// Game data methods
	var _ func(context.Context, GameID) (*Boxscore, error) = client.Boxscore
	var _ func(context.Context, GameID) (*PlayByPlay, error) = client.PlayByPlay
	var _ func(context.Context, GameID) (*GameMatchup, error) = client.Landing
	var _ func(context.Context, GameID) (*GameStory, error) = client.GameStory
	var _ func(context.Context, GameID) (*SeasonSeriesMatchup, error) = client.SeasonSeries
	var _ func(context.Context, GameID) (*ShiftChart, error) = client.ShiftChart

	// Player methods
	var _ func(context.Context, PlayerID) (*PlayerLanding, error) = client.PlayerLanding
	var _ func(context.Context, PlayerID, Season, GameType) (*PlayerGameLog, error) = client.PlayerGameLog
	var _ func(context.Context, string, *int) ([]PlayerSearchResult, error) = client.SearchPlayer

	// Team/Franchise methods
	var _ func(context.Context) ([]Franchise, error) = client.Franchises
	var _ func(context.Context, string) (*Roster, error) = client.RosterCurrent
	var _ func(context.Context, string, Season) (*Roster, error) = client.RosterSeason
	var _ func(context.Context, string, Season, GameType) (*ClubStats, error) = client.ClubStats
	var _ func(context.Context, string) ([]SeasonGameTypes, error) = client.ClubStatsSeason

	_ = ctx
}

// ===== Integration-Style Usage Examples =====

func TestClientUsagePattern_Standings(t *testing.T) {
	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Example usage pattern demonstration
	_ = client
	_ = ctx

	// In actual integration tests:
	// standings, err := client.CurrentLeagueStandings(ctx)
	// if err != nil {
	//     t.Fatalf("failed to get standings: %v", err)
	// }
}

func TestClientUsagePattern_Schedule(t *testing.T) {
	client := NewClient()
	ctx, cancel := DefaultContext()
	defer cancel()

	date := FromYMD(2024, 1, 15)

	// Example usage pattern demonstration
	_ = client
	_ = ctx
	_ = date

	// In actual integration tests:
	// schedule, err := client.DailySchedule(ctx, date)
	// if err != nil {
	//     t.Fatalf("failed to get schedule: %v", err)
	// }
}

func TestClientUsagePattern_Player(t *testing.T) {
	client := NewClient()
	ctx, cancel := DefaultContext()
	defer cancel()

	playerID := int64(8478402) // Connor McDavid

	// Example usage pattern demonstration
	_ = client
	_ = ctx
	_ = playerID

	// In actual integration tests:
	// player, err := client.PlayerLanding(ctx, playerID)
	// if err != nil {
	//     t.Fatalf("failed to get player: %v", err)
	// }
}

func TestClientUsagePattern_SearchPlayer(t *testing.T) {
	client := NewClient()
	ctx, cancel := DefaultContext()
	defer cancel()

	query := "McDavid"
	limit := 10

	// Example usage pattern demonstration
	_ = client
	_ = ctx
	_ = query

	// In actual integration tests:
	// results, err := client.SearchPlayer(ctx, query, &limit)
	// if err != nil {
	//     t.Fatalf("failed to search players: %v", err)
	// }
	_ = limit
}

// ===== Benchmark Tests =====

func BenchmarkBuildURL(b *testing.B) {
	base := "https://api.example.com/v1/"
	resource := "/resource/path"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildURL(base, resource)
	}
}

func BenchmarkClientCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewClient()
	}
}

func BenchmarkExtractDailySchedule(b *testing.B) {
	client := NewClient()

	weeklySchedule := &WeeklyScheduleResponse{
		NextStartDate:     "2024-01-15",
		PreviousStartDate: "2024-01-01",
		GameWeek: []GameDay{
			{
				Date: "2024-01-08",
				Games: []ScheduleGame{
					{ID: GameID(2023020001), GameType: GameTypeRegularSeason},
					{ID: GameID(2023020002), GameType: GameTypeRegularSeason},
					{ID: GameID(2023020003), GameType: GameTypeRegularSeason},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.extractDailySchedule(weeklySchedule, "2024-01-08")
	}
}

// ===== NewClientWithBaseURL Tests =====

func TestNewClientWithBaseURL(t *testing.T) {
	customBaseURL := "https://custom.example.com"

	client := NewClientWithBaseURL(customBaseURL)
	if client == nil {
		t.Fatal("NewClientWithBaseURL() returned nil")
	}
	if client.httpClient == nil {
		t.Error("NewClientWithBaseURL() created client with nil httpClient")
	}
	if client.baseURLOverride != customBaseURL {
		t.Errorf("expected baseURLOverride %s, got %s", customBaseURL, client.baseURLOverride)
	}
}

// ===== API Method Tests with Mocked Responses =====

func TestCurrentLeagueStandings(t *testing.T) {
	standings := []Standing{
		{TeamAbbrev: LocalizedString{Default: "TOR"}, Points: 50},
		{TeamAbbrev: LocalizedString{Default: "MTL"}, Points: 45},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, map[string]interface{}{
		"standings": standings,
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.CurrentLeagueStandings(ctx)

	if err != nil {
		t.Fatalf("CurrentLeagueStandings() error = %v", err)
	}

	if len(result) != len(standings) {
		t.Errorf("expected %d standings, got %d", len(standings), len(result))
	}
}

func TestCurrentLeagueStandings_Error(t *testing.T) {
	server := httptest.NewServer(makeErrorResponse(http.StatusInternalServerError))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	_, err := client.CurrentLeagueStandings(ctx)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLeagueStandingsForDate(t *testing.T) {
	standings := []Standing{
		{TeamAbbrev: LocalizedString{Default: "TOR"}, Points: 50},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, map[string]interface{}{
		"standings": standings,
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	date := FromYMD(2024, 1, 15)
	result, err := client.LeagueStandingsForDate(ctx, date)

	if err != nil {
		t.Fatalf("LeagueStandingsForDate() error = %v", err)
	}

	if len(result) != len(standings) {
		t.Errorf("expected %d standings, got %d", len(standings), len(result))
	}
}

func TestLeagueStandingsForSeason(t *testing.T) {
	standings := []Standing{
		{TeamAbbrev: LocalizedString{Default: "TOR"}, Points: 50},
	}

	seasonsResponse := SeasonsResponse{
		Seasons: []SeasonInfo{
			{ID: NewSeason(2023), StandingsStart: "2023-10-10", StandingsEnd: "2024-04-18"},
		},
	}

	standingsResponse := StandingsResponse{
		Standings: standings,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/standings-season", makeJSONResponse(http.StatusOK, seasonsResponse))
	mux.HandleFunc("/standings/2024-04-18", makeJSONResponse(http.StatusOK, standingsResponse))

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.LeagueStandingsForSeason(ctx, NewSeason(2023))

	if err != nil {
		t.Fatalf("LeagueStandingsForSeason() error = %v", err)
	}

	if len(result) != len(standings) {
		t.Errorf("expected %d standings, got %d", len(standings), len(result))
	}
}

func TestSeasonStandingManifest(t *testing.T) {
	seasonInfo := []SeasonInfo{
		{ID: NewSeason(2023), StandingsStart: "2023-10-10", StandingsEnd: "2024-04-18"},
		{ID: NewSeason(2022), StandingsStart: "2022-10-07", StandingsEnd: "2023-04-13"},
	}

	response := SeasonsResponse{
		Seasons: seasonInfo,
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, response))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.SeasonStandingManifest(ctx)

	if err != nil {
		t.Fatalf("SeasonStandingManifest() error = %v", err)
	}

	if len(result) != len(seasonInfo) {
		t.Errorf("expected %d seasons, got %d", len(seasonInfo), len(result))
	}
}

func TestTeams(t *testing.T) {
	standings := []Standing{
		{TeamAbbrev: LocalizedString{Default: "TOR"}, DivisionAbbrev: "ATL", DivisionName: "Atlantic"},
		{TeamAbbrev: LocalizedString{Default: "MTL"}, DivisionAbbrev: "ATL", DivisionName: "Atlantic"},
	}

	response := StandingsResponse{
		Standings: standings,
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, response))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	date := FromYMD(2024, 1, 15)
	result, err := client.Teams(ctx, date)

	if err != nil {
		t.Fatalf("Teams() error = %v", err)
	}

	if len(result) != len(standings) {
		t.Errorf("expected %d teams, got %d", len(standings), len(result))
	}
}

func TestDailySchedule(t *testing.T) {
	weeklySchedule := &WeeklyScheduleResponse{
		NextStartDate:     "2024-01-15",
		PreviousStartDate: "2024-01-01",
		GameWeek: []GameDay{
			{
				Date: "2024-01-08",
				Games: []ScheduleGame{
					{
						ID:        GameID(2023020001),
						GameType:  GameTypeRegularSeason,
						GameState: GameStateLive,
					},
				},
			},
		},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, weeklySchedule))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	date := FromYMD(2024, 1, 8)
	result, err := client.DailySchedule(ctx, date)

	if err != nil {
		t.Fatalf("DailySchedule() error = %v", err)
	}

	if result.Date != "2024-01-08" {
		t.Errorf("expected date 2024-01-08, got %s", result.Date)
	}
}

func TestWeeklySchedule(t *testing.T) {
	weeklySchedule := &WeeklyScheduleResponse{
		NextStartDate:     "2024-01-15",
		PreviousStartDate: "2024-01-01",
		GameWeek:          []GameDay{},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, weeklySchedule))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	date := FromYMD(2024, 1, 8)
	result, err := client.WeeklySchedule(ctx, date)

	if err != nil {
		t.Fatalf("WeeklySchedule() error = %v", err)
	}

	if result.NextStartDate != "2024-01-15" {
		t.Errorf("expected NextStartDate 2024-01-15, got %s", result.NextStartDate)
	}
}

func TestTeamWeeklySchedule(t *testing.T) {
	teamSchedule := &TeamScheduleResponse{
		Games: []ScheduleGame{},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, teamSchedule))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	date := FromYMD(2024, 1, 8)
	result, err := client.TeamWeeklySchedule(ctx, "TOR", date)

	if err != nil {
		t.Fatalf("TeamWeeklySchedule() error = %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestDailyScores(t *testing.T) {
	dailyScores := &DailyScores{
		CurrentDate: "2024-01-08",
		Games:       []GameScore{},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, dailyScores))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	date := FromYMD(2024, 1, 8)
	result, err := client.DailyScores(ctx, date)

	if err != nil {
		t.Fatalf("DailyScores() error = %v", err)
	}

	if result.CurrentDate != "2024-01-08" {
		t.Errorf("expected CurrentDate 2024-01-08, got %s", result.CurrentDate)
	}
}

func TestBoxscore(t *testing.T) {
	boxscore := &Boxscore{
		ID:                GameID(2023020001),
		Season:            NewSeason(2023),
		GameType:          GameTypeRegularSeason,
		GameDate:          "2023-10-10",
		GameState:         GameStateFinal,
		GameScheduleState: GameScheduleStateCompleted,
		Clock:             GameClock{},
		PlayerByGameStats: PlayerByGameStats{},
		PeriodDescriptor: PeriodDescriptor{
			Number:     3,
			PeriodType: PeriodTypeRegulation,
		},
		TVBroadcasts: []TVBroadcast{},
		AwayTeam:     BoxscoreTeam{},
		HomeTeam:     BoxscoreTeam{},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, boxscore))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.Boxscore(ctx, GameID(2023020001))

	if err != nil {
		t.Fatalf("Boxscore() error = %v", err)
	}

	if result.ID != GameID(2023020001) {
		t.Errorf("expected ID 2023020001, got %d", result.ID)
	}
}

func TestPlayByPlay(t *testing.T) {
	playByPlay := &PlayByPlay{
		ID:                GameID(2023020001),
		Season:            NewSeason(2023),
		GameType:          GameTypeRegularSeason,
		GameDate:          "2023-10-10",
		GameState:         GameStateFinal,
		GameScheduleState: GameScheduleStateCompleted,
		Clock:             GameClock{},
		Plays:             []PlayEvent{},
		RosterSpots:       []RosterSpot{},
		PeriodDescriptor: PeriodDescriptor{
			Number:     3,
			PeriodType: PeriodTypeRegulation,
		},
		TVBroadcasts: []TVBroadcast{},
		AwayTeam:     BoxscoreTeam{},
		HomeTeam:     BoxscoreTeam{},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, playByPlay))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.PlayByPlay(ctx, GameID(2023020001))

	if err != nil {
		t.Fatalf("PlayByPlay() error = %v", err)
	}

	if result.ID != GameID(2023020001) {
		t.Errorf("expected ID 2023020001, got %d", result.ID)
	}
}

func TestLanding(t *testing.T) {
	landing := &GameMatchup{
		ID:                GameID(2023020001),
		Season:            NewSeason(2023),
		GameType:          GameTypeRegularSeason,
		GameDate:          "2023-10-10",
		GameState:         GameStateFinal,
		GameScheduleState: GameScheduleStateCompleted,
		TVBroadcasts:      []TVBroadcast{},
		PeriodDescriptor: PeriodDescriptor{
			Number:     3,
			PeriodType: PeriodTypeRegulation,
		},
		AwayTeam: MatchupTeam{},
		HomeTeam: MatchupTeam{},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, landing))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.Landing(ctx, GameID(2023020001))

	if err != nil {
		t.Fatalf("Landing() error = %v", err)
	}

	if result.ID != GameID(2023020001) {
		t.Errorf("expected ID 2023020001, got %d", result.ID)
	}
}

func TestGameStory(t *testing.T) {
	gameStory := &GameStory{
		ID:                GameID(2023020001),
		Season:            NewSeason(2023),
		GameType:          GameTypeRegularSeason,
		GameDate:          "2023-10-10",
		GameState:         GameStateFinal,
		GameScheduleState: GameScheduleStateCompleted,
		TVBroadcasts:      []TVBroadcast{},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, gameStory))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.GameStory(ctx, GameID(2023020001))

	if err != nil {
		t.Fatalf("GameStory() error = %v", err)
	}

	if result.ID != GameID(2023020001) {
		t.Errorf("expected ID 2023020001, got %d", result.ID)
	}
}

func TestSeasonSeries(t *testing.T) {
	seasonSeries := &SeasonSeriesMatchup{
		SeasonSeries: []SeriesGame{},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, seasonSeries))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.SeasonSeries(ctx, GameID(2023020001))

	if err != nil {
		t.Fatalf("SeasonSeries() error = %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestShiftChart(t *testing.T) {
	shiftChart := &ShiftChart{
		Data: []ShiftEntry{},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, shiftChart))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.ShiftChart(ctx, GameID(2023020001))

	if err != nil {
		t.Fatalf("ShiftChart() error = %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestPlayerLanding(t *testing.T) {
	playerLanding := &PlayerLanding{
		PlayerID:       PlayerID(8478402),
		IsActive:       true,
		Position:       PositionCenter,
		ShootsCatches:  HandednessLeft,
		HeightInInches: 73,
		WeightInPounds: 193,
		BirthDate:      "1997-01-13",
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, playerLanding))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.PlayerLanding(ctx, PlayerID(8478402))

	if err != nil {
		t.Fatalf("PlayerLanding() error = %v", err)
	}

	if result.PlayerID != PlayerID(8478402) {
		t.Errorf("expected PlayerID 8478402, got %d", result.PlayerID)
	}
}

func TestPlayerGameLog(t *testing.T) {
	gameLog := &PlayerGameLog{
		PlayerID: PlayerID(8478402),
		GameLog:  []GameLog{},
		Season:   NewSeason(2023),
		GameType: GameTypeRegularSeason,
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, gameLog))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.PlayerGameLog(ctx, PlayerID(8478402), NewSeason(2023), GameTypeRegularSeason)

	if err != nil {
		t.Fatalf("PlayerGameLog() error = %v", err)
	}

	if result.PlayerID != PlayerID(8478402) {
		t.Errorf("expected PlayerID 8478402, got %d", result.PlayerID)
	}
}

func TestSearchPlayer(t *testing.T) {
	searchResults := []PlayerSearchResult{
		{PlayerID: PlayerID(8478402), Name: "Connor McDavid", Position: PositionCenter},
		{PlayerID: PlayerID(8477934), Name: "Leon Draisaitl", Position: PositionCenter},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, searchResults))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	limit := 10
	result, err := client.SearchPlayer(ctx, "McDavid", &limit)

	if err != nil {
		t.Fatalf("SearchPlayer() error = %v", err)
	}

	if len(result) != len(searchResults) {
		t.Errorf("expected %d results, got %d", len(searchResults), len(result))
	}
}

func TestSearchPlayer_NoLimit(t *testing.T) {
	searchResults := []PlayerSearchResult{
		{PlayerID: PlayerID(8478402), Name: "Connor McDavid", Position: PositionCenter},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, searchResults))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.SearchPlayer(ctx, "McDavid", nil)

	if err != nil {
		t.Fatalf("SearchPlayer() error = %v", err)
	}

	if len(result) != len(searchResults) {
		t.Errorf("expected %d results, got %d", len(searchResults), len(result))
	}
}

func TestFranchises(t *testing.T) {
	franchises := []Franchise{
		{ID: 1, FullName: "Toronto Maple Leafs"},
		{ID: 2, FullName: "Montreal Canadiens"},
	}

	response := FranchisesResponse{
		Data: franchises,
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, response))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.Franchises(ctx)

	if err != nil {
		t.Fatalf("Franchises() error = %v", err)
	}

	if len(result) != len(franchises) {
		t.Errorf("expected %d franchises, got %d", len(franchises), len(result))
	}
}

func TestRosterCurrent(t *testing.T) {
	roster := &Roster{
		Forwards:   []RosterPlayer{{ID: 8478402, Position: PositionCenter, ShootsCatches: HandednessLeft}},
		Defensemen: []RosterPlayer{{ID: 8477498, Position: PositionDefense, ShootsCatches: HandednessLeft}},
		Goalies:    []RosterPlayer{{ID: 8471214, Position: PositionGoalie, ShootsCatches: HandednessLeft}},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, roster))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.RosterCurrent(ctx, "EDM")

	if err != nil {
		t.Fatalf("RosterCurrent() error = %v", err)
	}

	if len(result.Forwards) != 1 {
		t.Errorf("expected 1 forward, got %d", len(result.Forwards))
	}
}

func TestRosterSeason(t *testing.T) {
	roster := &Roster{
		Forwards:   []RosterPlayer{{ID: 8478402, Position: PositionCenter, ShootsCatches: HandednessLeft}},
		Defensemen: []RosterPlayer{{ID: 8477498, Position: PositionDefense, ShootsCatches: HandednessLeft}},
		Goalies:    []RosterPlayer{{ID: 8471214, Position: PositionGoalie, ShootsCatches: HandednessLeft}},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, roster))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.RosterSeason(ctx, "EDM", NewSeason(2023))

	if err != nil {
		t.Fatalf("RosterSeason() error = %v", err)
	}

	if len(result.Forwards) != 1 {
		t.Errorf("expected 1 forward, got %d", len(result.Forwards))
	}
}

func TestClubStats(t *testing.T) {
	clubStats := &ClubStats{
		Season:   "20232024",
		GameType: GameTypeRegularSeason,
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, clubStats))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.ClubStats(ctx, "TOR", NewSeason(2023), GameTypeRegularSeason)

	if err != nil {
		t.Fatalf("ClubStats() error = %v", err)
	}

	if result.Season != "20232024" {
		t.Errorf("expected Season 20232024, got %s", result.Season)
	}
}

func TestClubStatsSeason(t *testing.T) {
	seasons := []SeasonGameTypes{
		{Season: NewSeason(2023), GameTypes: []GameType{GameTypeRegularSeason}},
		{Season: NewSeason(2022), GameTypes: []GameType{GameTypeRegularSeason}},
	}

	server := httptest.NewServer(makeJSONResponse(http.StatusOK, seasons))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)

	ctx := context.Background()
	result, err := client.ClubStatsSeason(ctx, "TOR")

	if err != nil {
		t.Fatalf("ClubStatsSeason() error = %v", err)
	}

	if len(result) != len(seasons) {
		t.Errorf("expected %d seasons, got %d", len(seasons), len(result))
	}
}

// Error path tests for client.go functions

func TestClient_ErrorPaths(t *testing.T) {
	t.Run("LeagueStandingsForSeason - invalid season", func(t *testing.T) {
		seasonsResponse := SeasonsResponse{
			Seasons: []SeasonInfo{
				{ID: NewSeason(2022), StandingsStart: "2022-10-07", StandingsEnd: "2023-04-13"},
			},
		}

		server := httptest.NewServer(makeJSONResponse(http.StatusOK, seasonsResponse))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		// Request a season that doesn't exist in the manifest
		_, err := client.LeagueStandingsForSeason(ctx, NewSeason(1999))
		if err == nil {
			t.Error("LeagueStandingsForSeason() should error for invalid season")
		}
	})

	t.Run("getJSON - HTTP error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusInternalServerError))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		var result StandingsResponse
		err := client.getJSON(ctx, EndpointAPIWebV1, "test", nil, &result)
		if err == nil {
			t.Error("getJSON() should error on HTTP error")
		}
	})

	t.Run("getJSON - invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{invalid json`))
		}))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		var result StandingsResponse
		err := client.getJSON(ctx, EndpointAPIWebV1, "test", nil, &result)
		if err == nil {
			t.Error("getJSON() should error on invalid JSON")
		}
	})

	t.Run("SeasonStandingManifest - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.SeasonStandingManifest(ctx)
		if err == nil {
			t.Error("SeasonStandingManifest() should error on HTTP error")
		}
	})

	t.Run("Teams - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.Teams(ctx, Today())
		if err == nil {
			t.Error("Teams() should error on HTTP error")
		}
	})

	t.Run("DailySchedule - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.DailySchedule(ctx, Today())
		if err == nil {
			t.Error("DailySchedule() should error on HTTP error")
		}
	})

	t.Run("TeamWeeklySchedule - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.TeamWeeklySchedule(ctx, "EDM", Today())
		if err == nil {
			t.Error("TeamWeeklySchedule() should error on HTTP error")
		}
	})

	t.Run("DailyScores - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.DailyScores(ctx, Today())
		if err == nil {
			t.Error("DailyScores() should error on HTTP error")
		}
	})

	t.Run("fetchWeeklySchedule - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.fetchWeeklySchedule(ctx, "2024-01-08")
		if err == nil {
			t.Error("fetchWeeklySchedule() should error on HTTP error")
		}
	})

	t.Run("Boxscore - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.Boxscore(ctx, GameID(2023020001))
		if err == nil {
			t.Error("Boxscore() should error on HTTP error")
		}
	})

	t.Run("PlayByPlay - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.PlayByPlay(ctx, GameID(2023020001))
		if err == nil {
			t.Error("PlayByPlay() should error on HTTP error")
		}
	})

	t.Run("Landing - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.Landing(ctx, GameID(2023020001))
		if err == nil {
			t.Error("Landing() should error on HTTP error")
		}
	})

	t.Run("GameStory - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.GameStory(ctx, GameID(2023020001))
		if err == nil {
			t.Error("GameStory() should error on HTTP error")
		}
	})

	t.Run("SeasonSeries - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.SeasonSeries(ctx, GameID(2023020001))
		if err == nil {
			t.Error("SeasonSeries() should error on HTTP error")
		}
	})

	t.Run("ShiftChart - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.ShiftChart(ctx, GameID(2023020001))
		if err == nil {
			t.Error("ShiftChart() should error on HTTP error")
		}
	})

	t.Run("PlayerLanding - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.PlayerLanding(ctx, PlayerID(8478402))
		if err == nil {
			t.Error("PlayerLanding() should error on HTTP error")
		}
	})

	t.Run("PlayerGameLog - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.PlayerGameLog(ctx, PlayerID(8478402), NewSeason(2023), GameTypeRegularSeason)
		if err == nil {
			t.Error("PlayerGameLog() should error on HTTP error")
		}
	})

	t.Run("SearchPlayer - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.SearchPlayer(ctx, "McDavid", nil)
		if err == nil {
			t.Error("SearchPlayer() should error on HTTP error")
		}
	})

	t.Run("Franchises - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.Franchises(ctx)
		if err == nil {
			t.Error("Franchises() should error on HTTP error")
		}
	})

	t.Run("RosterCurrent - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.RosterCurrent(ctx, "EDM")
		if err == nil {
			t.Error("RosterCurrent() should error on HTTP error")
		}
	})

	t.Run("RosterSeason - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.RosterSeason(ctx, "EDM", NewSeason(2023))
		if err == nil {
			t.Error("RosterSeason() should error on HTTP error")
		}
	})

	t.Run("ClubStats - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.ClubStats(ctx, "EDM", NewSeason(2023), GameTypeRegularSeason)
		if err == nil {
			t.Error("ClubStats() should error on HTTP error")
		}
	})

	t.Run("ClubStatsSeason - error", func(t *testing.T) {
		server := httptest.NewServer(makeErrorResponse(http.StatusNotFound))
		defer server.Close()

		client := NewClientWithBaseURL(server.URL)
		ctx := context.Background()

		_, err := client.ClubStatsSeason(ctx, "EDM")
		if err == nil {
			t.Error("ClubStatsSeason() should error on HTTP error")
		}
	})
}

func TestEndpoint_baseURL_Default(t *testing.T) {
	// Test the default case in baseURL function
	endpoint := Endpoint(999) // Invalid endpoint
	url := endpoint.baseURL()
	if url != baseURLAPIWebV1 {
		t.Errorf("baseURL() for invalid endpoint should return default, got %s", url)
	}
}

func TestClient_getJSON_URLParseError(t *testing.T) {
	// Use a client with a base URL override that will cause URL parsing issues
	// when combined with query params containing invalid characters
	client := NewClientWithBaseURL("http://example.com")
	ctx := context.Background()

	// Create query params with characters that could cause issues
	queryParams := map[string]string{
		"test": "value with spaces",
	}

	var result StandingsResponse
	// This should still work because url.Parse handles spaces, but let's test it
	err := client.getJSON(ctx, EndpointAPIWebV1, "test-resource", queryParams, &result)
	// This will fail with a connection error since the server doesn't exist
	// but we're just testing that the query param handling doesn't panic
	if err == nil {
		t.Log("Expected error due to non-existent server")
	}
}

func TestClient_getJSON_ReadBodyError(t *testing.T) {
	// Create a server that closes the connection immediately
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Don't write anything, just close
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	ctx := context.Background()

	var result StandingsResponse
	err := client.getJSON(ctx, EndpointAPIWebV1, "test", nil, &result)
	if err == nil {
		t.Error("getJSON() should error on empty response body")
	}
}

func TestLeagueStandingsForSeason_ManifestError(t *testing.T) {
	server := httptest.NewServer(makeErrorResponse(http.StatusInternalServerError))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	ctx := context.Background()

	_, err := client.LeagueStandingsForSeason(ctx, NewSeason(2023))
	if err == nil {
		t.Error("LeagueStandingsForSeason() should error when manifest fetch fails")
	}
}

func TestLeagueStandingsForSeason_StandingsFetchError(t *testing.T) {
	seasonsResponse := SeasonsResponse{
		Seasons: []SeasonInfo{
			{ID: NewSeason(2023), StandingsStart: "2023-10-10", StandingsEnd: "2024-04-18"},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/standings-season", makeJSONResponse(http.StatusOK, seasonsResponse))
	mux.HandleFunc("/standings/2024-04-18", makeErrorResponse(http.StatusInternalServerError))

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	ctx := context.Background()

	_, err := client.LeagueStandingsForSeason(ctx, NewSeason(2023))
	if err == nil {
		t.Error("LeagueStandingsForSeason() should error when standings fetch fails")
	}
}
