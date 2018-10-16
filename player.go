package main

import "github.com/Squwid/bytegolf/aws"

// Player holds the information about each player and each of the holes that theyve golfed
type Player struct {
	User   *aws.User     `json:"user"`
	Scores map[int]int64 `json:"scores"` // holds scores for each hole
}

// Score is the structure that gets returned from the score method and holds data about the players score
type Score struct {
	TotalScore   int64
	HolesCorrect int
}

// NewPlayerFromUser returns a new player from the aws user
func NewPlayerFromUser(user *aws.User) *Player {
	return &Player{
		User:   user,
		Scores: make(map[int]int64),
	}
}

// NewPlayerFromEmail returns a new player after getting the user from aws hence the error
func NewPlayerFromEmail(email string) (*Player, error) {
	user, err := getAwsUser(email)
	if err != nil {
		return nil, err
	}
	return NewPlayerFromUser(user), nil
}

// Score counts the total score amount of a player
func (p *Player) Score() *Score {
	var totalScore int64
	var holesCorrect int
	for _, score := range p.Scores {
		totalScore += score
		holesCorrect++
	}
	return &Score{totalScore, holesCorrect}
}

// HoleCorrect checks to see if a player got a specific hole correct or not
func (p *Player) HoleCorrect(hole int) bool {
	if _, ok := p.Scores[hole]; ok {
		return true
	}
	return false
}

// HoleScore lets you check what a player got on a specific hole
func (p *Player) HoleScore(hole int) int64 {
	if score, ok := p.Scores[hole]; ok {
		return score
	}
	return 0
}
