# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -v ./nhl -run TestClientMethodSignatures

# Run benchmarks
go test -bench=. ./nhl

# Build (verify compilation)
go build ./...

# Format code
gofmt -w .

# Vet code
go vet ./...
```

## Architecture

This is a Go client library for the NHL Stats API. All code lives in the `nhl` package.

### Core Components

**Client (`client.go`)**: The main API client that wraps HTTP requests to NHL endpoints. Uses `NewClientWithBaseURL()` for testing with mock servers.

**Endpoints**: The client communicates with four NHL API endpoints:
- `api-web.nhle.com/v1/` - Primary web API (standings, schedules, boxscores, players)
- `api.nhle.com/` - Core API
- `api.nhle.com/stats/rest/` - Stats API (shift charts, franchises)
- `search.d3.nhle.com/api/v1/` - Search API (player search)

### Type System

**Strongly-typed enums**: `GameType`, `GameState`, `Position`, `Handedness`, `PeriodType`, `HomeRoad`, `ZoneCode`, `PlayEventType`, `GameScheduleState`, `DefendingSide` - all implement custom JSON marshaling/unmarshaling and validation.

**ID wrapper types** (prevent mixing up different identifier types):
- `GameID` (`game_id.go`): 10-digit game identifiers encoding season, game type, and game number. Use `GameID(2024020001)`.
- `PlayerID` (`player_id.go`): Player identifiers. Unmarshals from int or string JSON. Use `PlayerID(8478402)`.
- `TeamID` (`team_id.go`): Team identifiers. Use `TeamID(10)`.
- `Season` (`season.go`): Season values like 20232024. Use `NewSeason(2023)` for the 2023-2024 season. Unmarshals from int, int64, or string JSON. `String()` returns `"2023-2024"` format.

**Date handling**: `GameDate` handles NHL-specific date format (YYYY-MM-DD).

**LocalizedString**: Handles NHL API's `{"default": "value"}` format for internationalized strings.

### API Response Types

Response types match NHL API structure:
- `Standing`, `StandingsResponse` - Team standings
- `ScheduleGame`, `DailySchedule`, `WeeklyScheduleResponse` - Game schedules
- `Boxscore`, `PlayByPlay`, `GameMatchup` - Game data
- `PlayerLanding`, `PlayerGameLog`, `PlayerSearchResult` - Player data
- `ClubStats`, `ClubSkaterStats`, `ClubGoalieStats` - Team statistics
- `Roster`, `RosterPlayer` - Team rosters

### Error Handling

Custom error types in `errors.go` map to HTTP status codes:
- `ResourceNotFoundError` (404)
- `RateLimitExceededError` (429)
- `BadRequestError` (400)
- `ServerError` (5xx)
- `RequestError`, `JSONError` - Wrap underlying errors

### Testing Pattern

Tests use `httptest.NewServer` with `NewClientWithBaseURL()` for mocking. Helper functions `makeJSONResponse()` and `makeErrorResponse()` simplify test setup.
