package main

/*
	Scoring System:
	Each character that is not white space or newline is 1 point
	Max score per problem is 300 (Maybe this needs to change)



*/

// Leaderboard holds the leaderboards for the game
type Leaderboard struct {
	Leaderboards map[int]Player `json:"leaderboards"`
}

// NewLeaderboard returns a new leaderboard with an initiated map
func NewLeaderboard() Leaderboard {
	return Leaderboard{
		Leaderboards: make(map[int]Player),
	}
}
