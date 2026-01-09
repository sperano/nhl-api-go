// Package nhl provides a client for interacting with the NHL Stats API.
package nhl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Endpoint represents different NHL API base URLs.
type Endpoint int

const (
	// EndpointAPIWebV1 is the primary web API endpoint.
	EndpointAPIWebV1 Endpoint = iota
	// EndpointAPICore is the core API endpoint.
	EndpointAPICore
	// EndpointAPIStats is the stats API endpoint.
	EndpointAPIStats
	// EndpointSearchV1 is the search API endpoint.
	EndpointSearchV1
)

const (
	baseURLAPIWebV1  = "https://api-web.nhle.com/v1/"
	baseURLAPICore   = "https://api.nhle.com/"
	baseURLAPIStats  = "https://api.nhle.com/stats/rest/"
	baseURLSearchV1  = "https://search.d3.nhle.com/api/v1/"
	defaultUserAgent = "nhl-api-go/1.0"
)

// baseURL returns the base URL for the given endpoint.
func (e Endpoint) baseURL() string {
	switch e {
	case EndpointAPIWebV1:
		return baseURLAPIWebV1
	case EndpointAPICore:
		return baseURLAPICore
	case EndpointAPIStats:
		return baseURLAPIStats
	case EndpointSearchV1:
		return baseURLSearchV1
	default:
		return baseURLAPIWebV1
	}
}

// Client is an HTTP client for the NHL Stats API.
type Client struct {
	httpClient      *http.Client
	baseURLOverride string
}

// NewClient creates a new NHL API client with default configuration.
func NewClient() *Client {
	config := DefaultClientConfig()
	return NewClientWithConfig(config)
}

// NewClientWithConfig creates a new NHL API client with the provided configuration.
func NewClientWithConfig(config *ClientConfig) *Client {
	return &Client{
		httpClient: config.ToHTTPClient(),
	}
}

// NewClientWithBaseURL creates a new NHL API client for testing with a custom base URL.
func NewClientWithBaseURL(baseURL string) *Client {
	return &Client{
		httpClient:      http.DefaultClient,
		baseURLOverride: baseURL,
	}
}

// buildURL constructs a full URL from a base URL and resource path.
// Handles proper slash normalization between base and resource.
func buildURL(base, resource string) string {
	if base == "" || resource == "" {
		return base + resource
	}

	baseEndsWithSlash := strings.HasSuffix(base, "/")
	resourceStartsWithSlash := strings.HasPrefix(resource, "/")

	if baseEndsWithSlash && resourceStartsWithSlash {
		return base + resource[1:]
	} else if !baseEndsWithSlash && !resourceStartsWithSlash {
		return base + "/" + resource
	}
	return base + resource
}

// getJSON performs an HTTP GET request and unmarshals the JSON response.
// Returns an appropriate error type based on HTTP status code.
func (c *Client) getJSON(ctx context.Context, endpoint Endpoint, resource string, queryParams map[string]string, result interface{}) error {
	var fullURL string
	if c.baseURLOverride != "" {
		fullURL = buildURL(c.baseURLOverride, resource)
	} else {
		fullURL = buildURL(endpoint.baseURL(), resource)
	}

	// Add query parameters if provided
	if len(queryParams) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return NewRequestError(fmt.Errorf("parsing URL %s: %w", fullURL, err))
		}
		q := u.Query()
		for key, value := range queryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return NewRequestError(fmt.Errorf("creating request: %w", err))
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NewRequestError(fmt.Errorf("executing request to %s: %w", fullURL, err))
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message := fmt.Sprintf("Request to %s failed", resource)
		return ErrorFromStatusCode(resp.StatusCode, message)
	}

	// Read and unmarshal response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewRequestError(fmt.Errorf("reading response body: %w", err))
	}

	if err := json.Unmarshal(body, result); err != nil {
		return NewJSONError(fmt.Errorf("unmarshaling response from %s: %w", fullURL, err))
	}

	return nil
}

// ===== Standings Methods =====

// CurrentLeagueStandings returns the current NHL standings.
func (c *Client) CurrentLeagueStandings(ctx context.Context) ([]Standing, error) {
	now := Now()
	return c.LeagueStandingsForDate(ctx, now)
}

// LeagueStandingsForDate returns league standings for a specific date.
func (c *Client) LeagueStandingsForDate(ctx context.Context, date GameDate) ([]Standing, error) {
	var response StandingsResponse
	resource := fmt.Sprintf("standings/%s", date.ToAPIString())
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return response.Standings, nil
}

// LeagueStandingsForSeason returns league standings for a specific season.
func (c *Client) LeagueStandingsForSeason(ctx context.Context, season Season) ([]Standing, error) {
	seasons, err := c.SeasonStandingManifest(ctx)
	if err != nil {
		return nil, err
	}

	// Find the season info for the requested season
	seasonID := season.ToInt64()
	var seasonInfo *SeasonInfo
	for i := range seasons {
		if seasons[i].ID.ToInt64() == seasonID {
			seasonInfo = &seasons[i]
			break
		}
	}

	if seasonInfo == nil {
		return nil, fmt.Errorf("invalid season: %s", season.String())
	}

	// Fetch standings for the end date of the season
	var response StandingsResponse
	resource := fmt.Sprintf("standings/%s", seasonInfo.StandingsEnd)
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return response.Standings, nil
}

// SeasonStandingManifest returns metadata for all NHL seasons.
func (c *Client) SeasonStandingManifest(ctx context.Context) ([]SeasonInfo, error) {
	var response SeasonsResponse
	if err := c.getJSON(ctx, EndpointAPIWebV1, "standings-season", nil, &response); err != nil {
		return nil, err
	}
	return response.Seasons, nil
}

// Teams returns all NHL teams for a specific date.
func (c *Client) Teams(ctx context.Context, date GameDate) ([]Team, error) {
	standingsResponse, err := c.LeagueStandingsForDate(ctx, date)
	if err != nil {
		return nil, err
	}

	teams := make([]Team, len(standingsResponse))
	for i, standing := range standingsResponse {
		teams[i] = standing.ToTeam()
	}
	return teams, nil
}

// ===== Schedule Methods =====

// DailySchedule returns the schedule for a specific date.
func (c *Client) DailySchedule(ctx context.Context, date GameDate) (*DailySchedule, error) {
	dateString := date.ToAPIString()
	weeklySchedule, err := c.fetchWeeklySchedule(ctx, dateString)
	if err != nil {
		return nil, err
	}

	return c.extractDailySchedule(weeklySchedule, dateString), nil
}

// WeeklySchedule returns the schedule for a week starting from the specified date.
func (c *Client) WeeklySchedule(ctx context.Context, date GameDate) (*WeeklyScheduleResponse, error) {
	return c.fetchWeeklySchedule(ctx, date.ToAPIString())
}

// TeamWeeklySchedule returns the weekly schedule for a specific team.
// The teamAbbr should be a team abbreviation like "MTL", "TOR", etc.
func (c *Client) TeamWeeklySchedule(ctx context.Context, teamAbbr string, date GameDate) (*TeamScheduleResponse, error) {
	var response TeamScheduleResponse
	resource := fmt.Sprintf("club-schedule/%s/week/%s", teamAbbr, date.ToAPIString())
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// DailyScores returns game scores for a specific date.
func (c *Client) DailyScores(ctx context.Context, date GameDate) (*DailyScores, error) {
	var response DailyScores
	resource := fmt.Sprintf("score/%s", date.ToAPIString())
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// fetchWeeklySchedule is a helper to fetch weekly schedule data.
func (c *Client) fetchWeeklySchedule(ctx context.Context, dateString string) (*WeeklyScheduleResponse, error) {
	var response WeeklyScheduleResponse
	resource := fmt.Sprintf("schedule/%s", dateString)
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// extractDailySchedule extracts a single day's schedule from weekly schedule data.
func (c *Client) extractDailySchedule(weeklySchedule *WeeklyScheduleResponse, dateString string) *DailySchedule {
	var games []ScheduleGame

	// Find the games for the requested date
	for _, day := range weeklySchedule.GameWeek {
		if day.Date == dateString {
			games = day.Games
			break
		}
	}

	// If no games found, initialize empty slice
	if games == nil {
		games = []ScheduleGame{}
	}

	return &DailySchedule{
		NextStartDate:     &weeklySchedule.NextStartDate,
		PreviousStartDate: &weeklySchedule.PreviousStartDate,
		Date:              dateString,
		Games:             games,
		NumberOfGames:     len(games),
	}
}

// ===== Game Data Methods =====

// Boxscore returns detailed boxscore data for a game.
func (c *Client) Boxscore(ctx context.Context, gameID GameID) (*Boxscore, error) {
	var response Boxscore
	if err := c.fetchGamecenter(ctx, gameID, "boxscore", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// PlayByPlay returns play-by-play data for a game.
func (c *Client) PlayByPlay(ctx context.Context, gameID GameID) (*PlayByPlay, error) {
	var response PlayByPlay
	if err := c.fetchGamecenter(ctx, gameID, "play-by-play", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Landing returns game landing/matchup data (lighter than play-by-play).
func (c *Client) Landing(ctx context.Context, gameID GameID) (*GameMatchup, error) {
	var response GameMatchup
	if err := c.fetchGamecenter(ctx, gameID, "landing", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GameStory returns narrative game story content.
func (c *Client) GameStory(ctx context.Context, gameID GameID) (*GameStory, error) {
	var response GameStory
	resource := fmt.Sprintf("wsc/game-story/%s", gameID.String())
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// SeasonSeries returns season series matchup data including head-to-head records.
func (c *Client) SeasonSeries(ctx context.Context, gameID GameID) (*SeasonSeriesMatchup, error) {
	var response SeasonSeriesMatchup
	if err := c.fetchGamecenter(ctx, gameID, "right-rail", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ShiftChart returns shift chart data for a game.
func (c *Client) ShiftChart(ctx context.Context, gameID GameID) (*ShiftChart, error) {
	cayenneExpr := fmt.Sprintf(
		"gameId=%s and ((duration != '00:00' and typeCode = 517) or typeCode != 517)",
		gameID.String(),
	)

	params := map[string]string{
		"cayenneExp": cayenneExpr,
		"exclude":    "eventDetails",
	}

	var response ShiftChart
	if err := c.getJSON(ctx, EndpointAPIStats, "en/shiftcharts", params, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// fetchGamecenter is a helper to fetch data from gamecenter endpoints.
func (c *Client) fetchGamecenter(ctx context.Context, gameID GameID, resource string, result interface{}) error {
	fullResource := fmt.Sprintf("gamecenter/%s/%s", gameID.String(), resource)
	return c.getJSON(ctx, EndpointAPIWebV1, fullResource, nil, result)
}

// ===== Player Methods =====

// PlayerLanding returns comprehensive player profile data.
// The playerID should be an NHL player ID (e.g., 8478402 for Connor McDavid).
func (c *Client) PlayerLanding(ctx context.Context, playerID PlayerID) (*PlayerLanding, error) {
	var response PlayerLanding
	resource := fmt.Sprintf("player/%s/landing", playerID.String())
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// PlayerGameLog returns a game-by-game log for a player's season.
func (c *Client) PlayerGameLog(ctx context.Context, playerID PlayerID, season Season, gameType GameType) (*PlayerGameLog, error) {
	var response PlayerGameLog
	resource := fmt.Sprintf("player/%s/game-log/%s/%d", playerID.String(), season.ToAPIString(), gameType.ToInt())
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	// The API doesn't include player_id in the response, so we set it from the parameter
	response.PlayerID = playerID
	return &response, nil
}

// SearchPlayer searches for players by name.
// The limit parameter is optional; if nil, defaults to 20.
func (c *Client) SearchPlayer(ctx context.Context, query string, limit *int) ([]PlayerSearchResult, error) {
	limitValue := 20
	if limit != nil {
		limitValue = *limit
	}

	params := map[string]string{
		"culture": "en-us",
		"q":       query,
		"limit":   fmt.Sprintf("%d", limitValue),
	}

	var response []PlayerSearchResult
	if err := c.getJSON(ctx, EndpointSearchV1, "search/player", params, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// ===== Teams/Franchises Methods =====

// FranchisesResponse represents the API response for franchises.
type FranchisesResponse struct {
	Data []Franchise `json:"data"`
}

// Franchises returns a list of all NHL franchises (past and current).
func (c *Client) Franchises(ctx context.Context) ([]Franchise, error) {
	var response FranchisesResponse
	if err := c.getJSON(ctx, EndpointAPIStats, "en/franchise", nil, &response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

// RosterCurrent returns the current roster for a team.
// The teamAbbr should be a team abbreviation like "MTL", "TOR", etc.
func (c *Client) RosterCurrent(ctx context.Context, teamAbbr string) (*Roster, error) {
	var response Roster
	resource := fmt.Sprintf("roster/%s/current", teamAbbr)
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// RosterSeason returns the roster for a team in a specific season.
func (c *Client) RosterSeason(ctx context.Context, teamAbbr string, season Season) (*Roster, error) {
	var response Roster
	resource := fmt.Sprintf("roster/%s/%s", teamAbbr, season.ToAPIString())
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ClubStats returns player statistics for a team in a specific season.
func (c *Client) ClubStats(ctx context.Context, teamAbbr string, season Season, gameType GameType) (*ClubStats, error) {
	var response ClubStats
	resource := fmt.Sprintf("club-stats/%s/%s/%d", teamAbbr, season.ToAPIString(), gameType.ToInt())
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// ClubStatsSeason returns available seasons and game types for a team.
func (c *Client) ClubStatsSeason(ctx context.Context, teamAbbr string) ([]SeasonGameTypes, error) {
	var response []SeasonGameTypes
	resource := fmt.Sprintf("club-stats-season/%s", teamAbbr)
	if err := c.getJSON(ctx, EndpointAPIWebV1, resource, nil, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// ===== Helper Types and Methods =====

// DefaultContext returns a context with a default timeout.
// This is a convenience function for users who don't want to manage contexts.
func DefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}
