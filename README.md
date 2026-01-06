# nhl-api-go

A Go client library for the NHL Stats API.

## Installation

```bash
go get github.com/ericblue/nhl-api-go
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/ericblue/nhl-api-go/nhl"
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

- **Standings**: `CurrentLeagueStandings`, `LeagueStandingsForDate`, `LeagueStandingsForSeason`
- **Schedule**: `DailySchedule`, `WeeklySchedule`, `TeamWeeklySchedule`, `DailyScores`
- **Games**: `Boxscore`, `PlayByPlay`, `Landing`, `GameStory`, `SeasonSeries`, `ShiftChart`
- **Players**: `PlayerLanding`, `PlayerGameLog`, `SearchPlayer`
- **Teams**: `Teams`, `Franchises`, `RosterCurrent`, `RosterSeason`, `ClubStats`

## License

MIT
