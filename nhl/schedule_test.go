package nhl

import (
	"encoding/json"
	"testing"
)

// Helper function to create a pointer to a string.
func stringPtr(s string) *string {
	return &s
}

// Helper function to create a pointer to an int.
func intPtr(i int) *int {
	return &i
}

// teamBuilder provides a fluent interface for building test ScheduleTeam instances.
type teamBuilder struct {
	id        int64
	abbrev    string
	placeName *LocalizedString
	logo      string
	score     *int
}

// newTeamBuilder creates a new teamBuilder with sensible defaults.
func newTeamBuilder(abbrev string) *teamBuilder {
	return &teamBuilder{
		id:        1,
		abbrev:    abbrev,
		placeName: nil,
		logo:      "https://assets.nhle.com/logos/nhl/svg/" + abbrev + "_light.svg",
		score:     nil,
	}
}

// withID sets the team ID.
func (b *teamBuilder) withID(id int64) *teamBuilder {
	b.id = id
	return b
}

// withScore sets the team score.
func (b *teamBuilder) withScore(score int) *teamBuilder {
	b.score = intPtr(score)
	return b
}

// withPlaceName sets the team place name.
func (b *teamBuilder) withPlaceName(name string) *teamBuilder {
	b.placeName = &LocalizedString{Default: name}
	return b
}

// build returns the constructed ScheduleTeam.
func (b *teamBuilder) build() ScheduleTeam {
	return ScheduleTeam{
		ID:        b.id,
		Abbrev:    b.abbrev,
		PlaceName: b.placeName,
		Logo:      b.logo,
		Score:     b.score,
	}
}

// scheduleGameBuilder provides a fluent interface for building test ScheduleGame instances.
type scheduleGameBuilder struct {
	id           int64
	gameType     GameType
	gameDate     *string
	startTimeUTC string
	awayTeam     ScheduleTeam
	homeTeam     ScheduleTeam
	gameState    GameState
}

// newScheduleGameBuilder creates a new scheduleGameBuilder with sensible defaults.
func newScheduleGameBuilder(awayAbbrev, homeAbbrev string) *scheduleGameBuilder {
	return &scheduleGameBuilder{
		id:           2023020001,
		gameType:     GameTypeRegularSeason,
		gameDate:     nil,
		startTimeUTC: "23:00:00Z",
		awayTeam:     newTeamBuilder(awayAbbrev).withID(7).build(),
		homeTeam:     newTeamBuilder(homeAbbrev).withID(10).build(),
		gameState:    GameStateFuture,
	}
}

// withID sets the game ID.
func (b *scheduleGameBuilder) withID(id int64) *scheduleGameBuilder {
	b.id = id
	return b
}

// withGameDate sets the game date.
func (b *scheduleGameBuilder) withGameDate(date string) *scheduleGameBuilder {
	b.gameDate = stringPtr(date)
	return b
}

// withGameState sets the game state.
func (b *scheduleGameBuilder) withGameState(state GameState) *scheduleGameBuilder {
	b.gameState = state
	return b
}

// withAwayScore sets the away team score.
func (b *scheduleGameBuilder) withAwayScore(score int) *scheduleGameBuilder {
	b.awayTeam.Score = intPtr(score)
	return b
}

// withHomeScore sets the home team score.
func (b *scheduleGameBuilder) withHomeScore(score int) *scheduleGameBuilder {
	b.homeTeam.Score = intPtr(score)
	return b
}

// build returns the constructed ScheduleGame.
func (b *scheduleGameBuilder) build() ScheduleGame {
	return ScheduleGame{
		ID:           b.id,
		GameType:     b.gameType,
		GameDate:     b.gameDate,
		StartTimeUTC: b.startTimeUTC,
		AwayTeam:     b.awayTeam,
		HomeTeam:     b.homeTeam,
		GameState:    b.gameState,
	}
}

// gameScoreBuilder provides a fluent interface for building test GameScore instances.
type gameScoreBuilder struct {
	id        int64
	gameType  GameType
	gameState GameState
	awayTeam  ScheduleTeam
	homeTeam  ScheduleTeam
}

// newGameScoreBuilder creates a new gameScoreBuilder with sensible defaults.
func newGameScoreBuilder(awayAbbrev, homeAbbrev string) *gameScoreBuilder {
	return &gameScoreBuilder{
		id:        2023020001,
		gameType:  GameTypeRegularSeason,
		gameState: GameStateFuture,
		awayTeam:  newTeamBuilder(awayAbbrev).withID(7).build(),
		homeTeam:  newTeamBuilder(homeAbbrev).withID(10).build(),
	}
}

// withID sets the game ID.
func (b *gameScoreBuilder) withID(id int64) *gameScoreBuilder {
	b.id = id
	return b
}

// withGameState sets the game state.
func (b *gameScoreBuilder) withGameState(state GameState) *gameScoreBuilder {
	b.gameState = state
	return b
}

// withAwayScore sets the away team score.
func (b *gameScoreBuilder) withAwayScore(score int) *gameScoreBuilder {
	b.awayTeam.Score = intPtr(score)
	return b
}

// withHomeScore sets the home team score.
func (b *gameScoreBuilder) withHomeScore(score int) *gameScoreBuilder {
	b.homeTeam.Score = intPtr(score)
	return b
}

// build returns the constructed GameScore.
func (b *gameScoreBuilder) build() GameScore {
	return GameScore{
		ID:        b.id,
		GameType:  b.gameType,
		GameState: b.gameState,
		AwayTeam:  b.awayTeam,
		HomeTeam:  b.homeTeam,
	}
}

func TestDailyScheduleWithNoGames(t *testing.T) {
	schedule := DailySchedule{
		NextStartDate:     stringPtr("2024-10-20"),
		PreviousStartDate: stringPtr("2024-10-18"),
		Date:              "2024-10-19",
		Games:             []ScheduleGame{},
		NumberOfGames:     0,
	}

	if len(schedule.Games) != 0 {
		t.Errorf("expected 0 games, got %d", len(schedule.Games))
	}
	if schedule.NumberOfGames != 0 {
		t.Errorf("expected NumberOfGames = 0, got %d", schedule.NumberOfGames)
	}
}

func TestDailyScoresDeserialization(t *testing.T) {
	jsonData := `{
		"prevDate": "2024-10-18",
		"currentDate": "2024-10-19",
		"nextDate": "2024-10-20",
		"games": []
	}`

	var scores DailyScores
	if err := json.Unmarshal([]byte(jsonData), &scores); err != nil {
		t.Fatalf("failed to unmarshal DailyScores: %v", err)
	}

	if scores.CurrentDate != "2024-10-19" {
		t.Errorf("expected CurrentDate = 2024-10-19, got %s", scores.CurrentDate)
	}
	if len(scores.Games) != 0 {
		t.Errorf("expected 0 games, got %d", len(scores.Games))
	}
}

func TestScheduleGameDisplayWithDate(t *testing.T) {
	game := newScheduleGameBuilder("BUF", "TOR").
		withGameDate("2023-10-10").
		build()

	expected := "BUF @ TOR on 2023-10-10 [FUT]"
	if game.String() != expected {
		t.Errorf("expected %q, got %q", expected, game.String())
	}
}

func TestScheduleGameDisplayWithoutDate(t *testing.T) {
	game := newScheduleGameBuilder("BUF", "TOR").build()

	expected := "BUF @ TOR [FUT]"
	if game.String() != expected {
		t.Errorf("expected %q, got %q", expected, game.String())
	}
}

func TestScheduleGameDeserializationWithoutGameDate(t *testing.T) {
	jsonData := `{
		"id": 2024020001,
		"gameType": 2,
		"startTimeUTC": "23:00:00Z",
		"awayTeam": {
			"id": 7,
			"abbrev": "BUF",
			"logo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg"
		},
		"homeTeam": {
			"id": 10,
			"abbrev": "TOR",
			"logo": "https://assets.nhle.com/logos/nhl/svg/TOR_light.svg"
		},
		"gameState": "FUT"
	}`

	var game ScheduleGame
	if err := json.Unmarshal([]byte(jsonData), &game); err != nil {
		t.Fatalf("failed to unmarshal ScheduleGame: %v", err)
	}

	if game.ID != 2024020001 {
		t.Errorf("expected ID = 2024020001, got %d", game.ID)
	}
	if game.GameDate != nil {
		t.Errorf("expected GameDate = nil, got %v", game.GameDate)
	}
	if game.AwayTeam.Abbrev != "BUF" {
		t.Errorf("expected AwayTeam.Abbrev = BUF, got %s", game.AwayTeam.Abbrev)
	}
	if game.HomeTeam.Abbrev != "TOR" {
		t.Errorf("expected HomeTeam.Abbrev = TOR, got %s", game.HomeTeam.Abbrev)
	}
}

func TestGameScoreDisplayWithScores(t *testing.T) {
	game := newGameScoreBuilder("BUF", "TOR").
		withAwayScore(3).
		withHomeScore(2).
		withGameState(GameStateFinal).
		build()

	expected := "BUF 3 @ TOR 2 [FINAL]"
	if game.String() != expected {
		t.Errorf("expected %q, got %q", expected, game.String())
	}
}

func TestGameScoreDisplayWithNoScores(t *testing.T) {
	game := newGameScoreBuilder("BUF", "TOR").build()

	expected := "BUF - @ TOR - [FUT]"
	if game.String() != expected {
		t.Errorf("expected %q, got %q", expected, game.String())
	}
}

func TestGameScoreDisplayWithPartialScore(t *testing.T) {
	game := newGameScoreBuilder("BUF", "TOR").
		withAwayScore(1).
		withGameState(GameStateLive).
		build()

	expected := "BUF 1 @ TOR - [LIVE]"
	if game.String() != expected {
		t.Errorf("expected %q, got %q", expected, game.String())
	}
}

func TestGameScoreDisplayWithZeroScores(t *testing.T) {
	game := newGameScoreBuilder("BUF", "TOR").
		withAwayScore(0).
		withHomeScore(0).
		withGameState(GameStateLive).
		build()

	expected := "BUF 0 @ TOR 0 [LIVE]"
	if game.String() != expected {
		t.Errorf("expected %q, got %q", expected, game.String())
	}
}

func TestScheduleGameSerialization(t *testing.T) {
	game := newScheduleGameBuilder("BUF", "TOR").
		withGameDate("2023-10-10").
		withID(2023020050).
		build()

	data, err := json.Marshal(game)
	if err != nil {
		t.Fatalf("failed to marshal ScheduleGame: %v", err)
	}

	var unmarshaled ScheduleGame
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal ScheduleGame: %v", err)
	}

	if unmarshaled.ID != game.ID {
		t.Errorf("expected ID = %d, got %d", game.ID, unmarshaled.ID)
	}
	if unmarshaled.AwayTeam.Abbrev != game.AwayTeam.Abbrev {
		t.Errorf("expected AwayTeam.Abbrev = %s, got %s", game.AwayTeam.Abbrev, unmarshaled.AwayTeam.Abbrev)
	}
}

func TestDailyScheduleDeserialization(t *testing.T) {
	jsonData := `{
		"nextStartDate": "2024-10-20",
		"previousStartDate": "2024-10-18",
		"date": "2024-10-19",
		"numberOfGames": 2,
		"games": [
			{
				"id": 2024020001,
				"gameType": 2,
				"gameDate": "2024-10-19",
				"startTimeUTC": "23:00:00Z",
				"awayTeam": {
					"id": 7,
					"abbrev": "BUF",
					"logo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg"
				},
				"homeTeam": {
					"id": 10,
					"abbrev": "TOR",
					"logo": "https://assets.nhle.com/logos/nhl/svg/TOR_light.svg"
				},
				"gameState": "FUT"
			},
			{
				"id": 2024020002,
				"gameType": 2,
				"gameDate": "2024-10-19",
				"startTimeUTC": "00:00:00Z",
				"awayTeam": {
					"id": 1,
					"abbrev": "NJD",
					"logo": "https://assets.nhle.com/logos/nhl/svg/NJD_light.svg"
				},
				"homeTeam": {
					"id": 2,
					"abbrev": "NYI",
					"logo": "https://assets.nhle.com/logos/nhl/svg/NYI_light.svg"
				},
				"gameState": "LIVE"
			}
		]
	}`

	var schedule DailySchedule
	if err := json.Unmarshal([]byte(jsonData), &schedule); err != nil {
		t.Fatalf("failed to unmarshal DailySchedule: %v", err)
	}

	if schedule.Date != "2024-10-19" {
		t.Errorf("expected Date = 2024-10-19, got %s", schedule.Date)
	}
	if schedule.NumberOfGames != 2 {
		t.Errorf("expected NumberOfGames = 2, got %d", schedule.NumberOfGames)
	}
	if len(schedule.Games) != 2 {
		t.Fatalf("expected 2 games, got %d", len(schedule.Games))
	}

	if schedule.Games[0].AwayTeam.Abbrev != "BUF" {
		t.Errorf("expected first game AwayTeam.Abbrev = BUF, got %s", schedule.Games[0].AwayTeam.Abbrev)
	}
	if schedule.Games[1].GameState != GameStateLive {
		t.Errorf("expected second game GameState = LIVE, got %s", schedule.Games[1].GameState)
	}
}

func TestWeeklyScheduleResponseDeserialization(t *testing.T) {
	jsonData := `{
		"nextStartDate": "2024-10-27",
		"previousStartDate": "2024-10-13",
		"gameWeek": [
			{
				"date": "2024-10-20",
				"games": [
					{
						"id": 2024020001,
						"gameType": 2,
						"startTimeUTC": "23:00:00Z",
						"awayTeam": {
							"id": 7,
							"abbrev": "BUF",
							"logo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg"
						},
						"homeTeam": {
							"id": 10,
							"abbrev": "TOR",
							"logo": "https://assets.nhle.com/logos/nhl/svg/TOR_light.svg"
						},
						"gameState": "FUT"
					}
				]
			},
			{
				"date": "2024-10-21",
				"games": []
			}
		]
	}`

	var response WeeklyScheduleResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal WeeklyScheduleResponse: %v", err)
	}

	if response.NextStartDate != "2024-10-27" {
		t.Errorf("expected NextStartDate = 2024-10-27, got %s", response.NextStartDate)
	}
	if response.PreviousStartDate != "2024-10-13" {
		t.Errorf("expected PreviousStartDate = 2024-10-13, got %s", response.PreviousStartDate)
	}
	if len(response.GameWeek) != 2 {
		t.Fatalf("expected 2 game days, got %d", len(response.GameWeek))
	}
	if response.GameWeek[0].Date != "2024-10-20" {
		t.Errorf("expected first day date = 2024-10-20, got %s", response.GameWeek[0].Date)
	}
	if len(response.GameWeek[0].Games) != 1 {
		t.Errorf("expected 1 game on first day, got %d", len(response.GameWeek[0].Games))
	}
	if len(response.GameWeek[1].Games) != 0 {
		t.Errorf("expected 0 games on second day, got %d", len(response.GameWeek[1].Games))
	}
}

func TestTeamScheduleResponseDeserialization(t *testing.T) {
	jsonData := `{
		"games": [
			{
				"id": 2024020001,
				"gameType": 2,
				"gameDate": "2024-10-19",
				"startTimeUTC": "23:00:00Z",
				"awayTeam": {
					"id": 7,
					"abbrev": "BUF",
					"logo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg"
				},
				"homeTeam": {
					"id": 10,
					"abbrev": "TOR",
					"logo": "https://assets.nhle.com/logos/nhl/svg/TOR_light.svg"
				},
				"gameState": "FUT"
			}
		]
	}`

	var response TeamScheduleResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal TeamScheduleResponse: %v", err)
	}

	if len(response.Games) != 1 {
		t.Fatalf("expected 1 game, got %d", len(response.Games))
	}
	if response.Games[0].ID != 2024020001 {
		t.Errorf("expected game ID = 2024020001, got %d", response.Games[0].ID)
	}
}

func TestScheduleTeamWithPlaceName(t *testing.T) {
	jsonData := `{
		"id": 7,
		"abbrev": "BUF",
		"placeName": {"default": "Buffalo"},
		"logo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg"
	}`

	var team ScheduleTeam
	if err := json.Unmarshal([]byte(jsonData), &team); err != nil {
		t.Fatalf("failed to unmarshal ScheduleTeam: %v", err)
	}

	if team.ID != 7 {
		t.Errorf("expected ID = 7, got %d", team.ID)
	}
	if team.Abbrev != "BUF" {
		t.Errorf("expected Abbrev = BUF, got %s", team.Abbrev)
	}
	if team.PlaceName == nil {
		t.Fatal("expected PlaceName to be non-nil")
	}
	if team.PlaceName.Default != "Buffalo" {
		t.Errorf("expected PlaceName.Default = Buffalo, got %s", team.PlaceName.Default)
	}
}

func TestScheduleTeamWithScore(t *testing.T) {
	jsonData := `{
		"id": 7,
		"abbrev": "BUF",
		"logo": "https://assets.nhle.com/logos/nhl/svg/BUF_light.svg",
		"score": 5
	}`

	var team ScheduleTeam
	if err := json.Unmarshal([]byte(jsonData), &team); err != nil {
		t.Fatalf("failed to unmarshal ScheduleTeam: %v", err)
	}

	if team.Score == nil {
		t.Fatal("expected Score to be non-nil")
	}
	if *team.Score != 5 {
		t.Errorf("expected Score = 5, got %d", *team.Score)
	}
}

func TestGameScoreSerialization(t *testing.T) {
	game := newGameScoreBuilder("BUF", "TOR").
		withAwayScore(3).
		withHomeScore(2).
		withGameState(GameStateFinal).
		build()

	data, err := json.Marshal(game)
	if err != nil {
		t.Fatalf("failed to marshal GameScore: %v", err)
	}

	var unmarshaled GameScore
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal GameScore: %v", err)
	}

	if unmarshaled.ID != game.ID {
		t.Errorf("expected ID = %d, got %d", game.ID, unmarshaled.ID)
	}
	if unmarshaled.AwayTeam.Score == nil || *unmarshaled.AwayTeam.Score != 3 {
		t.Error("expected AwayTeam.Score = 3")
	}
	if unmarshaled.HomeTeam.Score == nil || *unmarshaled.HomeTeam.Score != 2 {
		t.Error("expected HomeTeam.Score = 2")
	}
}

func TestDailyScheduleNilNavigationDates(t *testing.T) {
	jsonData := `{
		"date": "2024-10-19",
		"numberOfGames": 0,
		"games": []
	}`

	var schedule DailySchedule
	if err := json.Unmarshal([]byte(jsonData), &schedule); err != nil {
		t.Fatalf("failed to unmarshal DailySchedule: %v", err)
	}

	if schedule.NextStartDate != nil {
		t.Errorf("expected NextStartDate = nil, got %v", schedule.NextStartDate)
	}
	if schedule.PreviousStartDate != nil {
		t.Errorf("expected PreviousStartDate = nil, got %v", schedule.PreviousStartDate)
	}
}

func TestGameDayEmptyGames(t *testing.T) {
	day := GameDay{
		Date:  "2024-10-19",
		Games: []ScheduleGame{},
	}

	if len(day.Games) != 0 {
		t.Errorf("expected 0 games, got %d", len(day.Games))
	}
}

func TestDailyScoresSerialization(t *testing.T) {
	scores := DailyScores{
		PrevDate:    "2024-10-18",
		CurrentDate: "2024-10-19",
		NextDate:    "2024-10-20",
		Games:       []GameScore{},
	}

	data, err := json.Marshal(scores)
	if err != nil {
		t.Fatalf("failed to marshal DailyScores: %v", err)
	}

	var unmarshaled DailyScores
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal DailyScores: %v", err)
	}

	if unmarshaled.CurrentDate != scores.CurrentDate {
		t.Errorf("expected CurrentDate = %s, got %s", scores.CurrentDate, unmarshaled.CurrentDate)
	}
}
