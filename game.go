package main

import (
	"fmt"
	"time"

	"github.com/Squwid/bytegolf/aws"

	uuid "github.com/satori/go.uuid"
)

// Game holds the game type
type Game struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Password    string               `json:"password"`
	Difficulty  string               `json:"difficulty"`
	Holes       int                  `json:"holes"`
	MaxPlayers  int                  `json:"maxPlayers"`
	CreatedTime time.Time            `json:"createdTime"`
	Questions   map[int]aws.Question `json:"questions"`

	CurrentPlayers int      `json:"currentPlayers"`
	Players        []Player `json:"players"`

	LastModified time.Time `json:"lastModified"`
	Started      bool      `json:"started"`
	StartTime    time.Time `json:"startTime"`
	Ended        bool      `json:"ended"`
	EndTime      time.Time `json:"endTime"`
}

// NewGame returns a pointer to a game and an error which would come either if server is out of mem or aws is down.
func NewGame(name, password, difficulty string, holes, maxPlayers int) (Game, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return Game{}, err
	}
	questions, err := aws.GetQuestionsDynamo(holes, difficulty)
	if err != nil {
		return Game{}, err
	}
	return Game{
		ID:             uuid.String(),
		Name:           name,
		Password:       password,
		Difficulty:     difficulty,
		Holes:          holes,
		MaxPlayers:     maxPlayers,
		CreatedTime:    time.Now(),
		Questions:      questions,
		CurrentPlayers: 0,
		Players:        []Player{},
		LastModified:   time.Now(),
		Started:        false,
		Ended:          false,
	}, nil
}

// Start starts a game with no time amount
func (game Game) Start() error {
	if game.Started {
		return fmt.Errorf("game %v was already started", game.ID)
	}
	game.Started = true
	game.StartTime = time.Now()
	game.LastModified = time.Now()
	return nil
}

// End ends a game and returns an error if the game had not started yet
func (game Game) End() error {
	if !game.Started {
		return fmt.Errorf("game %v has not started yet", game.ID)
	}
	game.Ended = true
	game.EndTime = time.Now()
	game.LastModified = time.Now()
	return nil
}

// CheckSubmission will check a submission of a specific hole against the correct output using
// explicit string comparison
func (game Game) CheckSubmission(hole int, submission string) bool {
	// TODO: this should be a REGEX instead of a string answer to allow for cool and complex answers
	expected := game.Questions[hole].Answer
	if expected == submission {
		return true
	}
	return false
}

// Contains checks to see if a game contains a player
func (game Game) Contains(player *Player) bool {
	for _, player := range game.Players {
		if player.User.Email == player.User.Email {
			return true
		}
	}
	return false
}

// Add adds a player to a game, after checking if they are in the game
func (game Game) Add(user aws.User) error {
	player := NewPlayerFromUser(user)
	if game.Contains(player) {
		return fmt.Errorf("player %s is already in this game", player.User.Email)
	}
	if !game.InProgress() {
		return fmt.Errorf("error adding player %s to not inprogress game %s", player.User.Email, game.ID)
	}
	if game.CurrentPlayers >= game.MaxPlayers {
		return fmt.Errorf("game is already full")
	}
	game.Players = append(game.Players, *player)
	game.CurrentPlayers++
	game.LastModified = time.Now()
	return nil
}

// GetPlayer gets a specific player from a specific game. It returns an error if the
// player is not in the game
func (game Game) GetPlayer(email string) (*Player, error) {
	for _, player := range game.Players {
		if player.User.Email == email {
			return &player, nil
		}
	}
	return nil, fmt.Errorf("could not find player %s", email)
}

// InProgress checks to see if a game is in progress by looking at the started and ended bools
func (game Game) InProgress() bool {
	if game.Started && !game.Ended {
		return true
	}
	return false
}
