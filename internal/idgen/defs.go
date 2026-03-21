package main

// IDDef defines a numeric ID wrapper type for code generation.
type IDDef struct {
	TypeName   string // e.g., "GameID"
	Doc        string // type doc comment (may contain embedded newlines with // prefixes)
	ErrorLabel string // label for error messages (e.g., "game ID")
	Receiver   string // receiver variable name (e.g., "g" for GameID)
}

var ids = []IDDef{
	{
		TypeName:   "GameID",
		Doc:        "GameID is a wrapper type for NHL game identifiers.\n// Game IDs are 10-digit integers in the format: SSSGTNNNN where:\n// - SSS is the first 4 digits of the season (e.g., 2023 for 2023-2024)\n// - GT is the game type (01=preseason, 02=regular, 03=playoffs, 04=all-star, 12=PWHL showcase)\n// - NNNN is the game number",
		ErrorLabel: "game ID",
		Receiver:   "g",
	},
	{
		TypeName:   "PlayerID",
		Doc:        "PlayerID is a wrapper type for NHL player identifiers.\n// Player IDs are numeric identifiers assigned to each player (e.g., 8478402 for Connor McDavid).",
		ErrorLabel: "player ID",
		Receiver:   "p",
	},
	{
		TypeName:   "TeamID",
		Doc:        "TeamID is a wrapper type for NHL team identifiers.\n// Team IDs are numeric identifiers assigned to each team (e.g., 10 for Toronto Maple Leafs).",
		ErrorLabel: "team ID",
		Receiver:   "t",
	},
}
