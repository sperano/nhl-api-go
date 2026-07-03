# nhl-api-go

A Go client library for the NHL Stats API.

## Installation

```bash
go get github.com/sperano/nhl-api-go
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/sperano/nhl-api-go/nhl"
)

func main() {
    client := nhl.NewClient()
    ctx := context.Background()

    // Get current standings
    standings, _ := client.CurrentLeagueStandings(ctx)
    fmt.Printf("Found %d teams\n", len(standings))

    // Get today's schedule
    schedule, _ := client.DailySchedule(ctx, nhl.Now())
    fmt.Printf("Games today: %d\n", schedule.NumberOfGames)

    // Search for a player
    players, _ := client.SearchPlayer(ctx, "McDavid", nil)
    fmt.Printf("Found: %s\n", players[0].Name)

    // Get game boxscore
    boxscore, _ := client.Boxscore(ctx, 2024020001)
    fmt.Printf("%s vs %s\n", boxscore.AwayTeam.Name.Default, boxscore.HomeTeam.Name.Default)
}
```

## Available Methods

- **Standings**: `CurrentLeagueStandings`, `LeagueStandingsForDate`, `LeagueStandingsForSeason`, `SeasonStandingManifest`
- **Schedule**: `DailySchedule`, `WeeklySchedule`, `TeamWeeklySchedule`, `ClubScheduleSeason`, `DailyScores`
- **Games**: `Boxscore`, `PlayByPlay`, `Landing`, `GameStory`, `SeasonSeries`, `ShiftChart`
- **Players**: `PlayerLanding`, `PlayerGameLog`, `SearchPlayer`
- **Teams**: `Teams`, `Franchises`, `RosterCurrent`, `RosterSeason`, `ClubStats`, `ClubStatsSeason`
- **Edge (skater)**: `EdgeSkaterDetail`, `EdgeSkaterSpeedDetail`, `EdgeSkaterDistanceDetail`, `EdgeSkaterShotSpeedDetail`, `EdgeSkaterShotLocationDetail`, `EdgeSkaterZoneTime`, `EdgeSkaterComparison`, `EdgeSkaterLanding`
- **Edge (goalie)**: `EdgeGoalieDetail`, `EdgeGoalie5v5Detail`, `EdgeGoalieShotLocationDetail`, `EdgeGoalieSavePctgDetail`, `EdgeGoalieComparison`, `EdgeGoalieLanding`
- **Edge (team)**: `EdgeTeamDetail`, `EdgeTeamSpeedDetail`, `EdgeTeamDistanceDetail`, `EdgeTeamShotSpeedDetail`, `EdgeTeamShotLocationDetail`, `EdgeTeamZoneTimeDetails`, `EdgeTeamComparison`, `EdgeTeamLanding`

## Error Handling

Non-2xx responses return an `*APIError` carrying the HTTP status code. Match well-known statuses with `errors.Is`:

```go
standings, err := client.CurrentLeagueStandings(ctx)
if errors.Is(err, nhl.ErrRateLimited) {
    // back off and retry
}
```

Sentinels: `ErrBadRequest` (400), `ErrUnauthorized` (401), `ErrNotFound` (404), `ErrRateLimited` (429), `ErrServerError` (500, and matches any 5xx). Transport failures are wrapped in `RequestError` and decode failures in `JSONError`. When the API returns an unrecognized enum value, the decode fails with a typed `*UnknownEnumValueError` (recoverable via `errors.As` to learn which enum type and value were unknown).

## License

MIT
