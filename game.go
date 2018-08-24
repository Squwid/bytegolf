package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Squwid/bytegolf/bgaws"
	"github.com/Squwid/bytegolf/runner"
	uuid "github.com/satori/go.uuid"
)

//TODO: Games currently dont end

// Errors
var (
	ErrGameFull       = errors.New("game already has maximum amount of players")
	ErrHoleNotFound   = errors.New("that hole was not found")
	ErrPlayerNotFound = errors.New("player not found in this game")
)

// CreateNewGame todo
func CreateNewGame(w http.ResponseWriter, req *http.Request) (*Game, error) {
	if len(strings.TrimSpace(req.FormValue("gamename"))) == 0 {
		return nil, errors.New("blank game name not allowed")
	}

	gameID, _ := uuid.NewV4()
	holes, err := strconv.Atoi(req.FormValue("holes"))
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	holes = 3 //todo: remove this. It is hardset to 3 right now because there are not enough

	// questions in the bank
	max, err := strconv.Atoi(req.FormValue("maxplayers"))
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// diff := req.FormValue("difficulty") // TODO: this is not currently an option
	diff := "medium"
	logger.Printf("new game requested with %v holes at %s difficulty\n", holes, diff)
	qs, err := bgaws.GetQuestions(diff, holes)
	if err != nil {
		return nil, err
	}

	user := getUser(w, req)
	player := createPlayer(user)
	g := Game{
		ID:             gameID.String(),
		Name:           req.FormValue("gamename"),
		Password:       req.FormValue("password"),
		CurrentPlayers: 0,
		MaxPlayers:     max,
		Holes:          holes,
		Difficulty:     diff,
		StartedTime:    time.Now(),
		Started:        true,
		Questions:      qs,
		Owner:          player,
	}
	err = g.AddGamePlayer(player)
	if err != nil {
		logger.Println(err)
		return &g, err
	}
	logger.Printf("%s created new game %s\n", user.Username, CurrentGame.Name)
	return &g, nil
}

// AddGameUser adds a user to the specified game
func (game *Game) AddGameUser(user *bgaws.User) error {
	if game.MaxPlayers == game.CurrentPlayers {
		return ErrGameFull
	}
	player := createPlayer(user)
	game.CurrentPlayers++ // add the player to list of players
	game.Players = append(game.Players, player)

	logger.Printf("%s added to game %s\n", user.Username, game.Name)
	logger.Printf("there are now %v people in game %s\n", game.CurrentPlayers, game.Name)
	games[player.User.Username] = game
	game.update()
	return nil
}

// AddGamePlayer adds a player to the specified game
func (game *Game) AddGamePlayer(player *Player) error {
	if game.MaxPlayers == game.CurrentPlayers {
		return ErrGameFull
	}
	game.CurrentPlayers++
	game.Players = append(game.Players, player)

	logger.Printf("%s added to game %s\n", player.User.Username, game.Name)
	logger.Printf("there are now %v people in game %s\n", game.CurrentPlayers, game.Name)
	games[player.User.Username] = game
	game.update()
	return nil
}

// GetPlayer gets a player from the game
func (game *Game) GetPlayer(user *bgaws.User) (*Player, error) {
	var player *Player
	for _, g := range game.Players {
		if g.User.Username == user.Username {
			return g, nil
		}
	}
	return player, ErrPlayerNotFound
}

// UserInGame checks to see if a user is in a specific game
func (game *Game) UserInGame(user *bgaws.User) bool {
	for _, p := range game.Players {
		if p.User.Username == user.Username {
			return true
		}
	}
	return false
}

// PlayerInGame checks to see if a certain player is in a game
func (game *Game) PlayerInGame(player *Player) bool {
	for _, p := range game.Players {
		if p.User.Username == player.User.Username {
			return true
		}
	}
	return false
}

// Score adds the score that is in the players submission to the scoreboard, however this function
// assumes that the submission is already correct
func (game *Game) Score(p *Player, hole int, sub *runner.CodeSubmission, resp *runner.CodeResponse) error {
	q, ok := game.Questions[hole]
	if !ok {
		return ErrHoleNotFound
	}

	points := score(sub, &q)
	p.Correct[hole] = true
	if p.Scores[hole] < points && p.Scores[hole] != 0 {
		return nil // this means that the previous score was already better
	}
	p.Scores[hole] = points
	p.Output[hole] = resp.Output
	// Add up each of the points each time the question is correct to not have to deal with
	// odd situations such as keep on submitting the same code over and over again
	var totalScore int64
	var totalHoles int
	for _, val := range p.Correct {
		if val {
			totalHoles++
		}
	}
	for _, val := range p.Scores {
		totalScore += val
	}
	p.TotalScore = totalScore
	p.HolesCorrect = totalHoles
	p.Average = float64(p.TotalScore) / float64(p.HolesCorrect)
	logger.Printf("%s now has %v total holes correct at %v points", p.User.Username, p.HolesCorrect, p.TotalScore)
	game.update()
	return nil
}

// Check checks to see if a CodeRunner response is correct based on the hole inside of that game
func (game *Game) Check(resp *runner.CodeResponse, hole int) (bool, error) {
	q, ok := game.Questions[hole]
	if !ok {
		return false, ErrHoleNotFound
	}
	if strings.TrimSpace(strings.ToLower(resp.Output)) == strings.TrimSpace(strings.ToLower(q.Answer)) {
		return true, nil
	}
	return false, nil
}

// update updates the game since the page is static, this method will be applied to every other
// game method when something changes inside of the game
func (game *Game) update() {
	var winning Player
	for _, p := range game.Players {
		if winning.Average == 0 {
			winning = *p
			continue
		}
		if p.Average < winning.Average && p.Average != 0.0 {
			winning = *p
			continue
		}
	}
	var others []*Player
	for _, p := range game.Players {
		if p.User.Username != winning.User.Username {
			others = append(others, p)
		}
	}
	//TODO: sort here
	game.Leaderboard.Winning = &winning
	game.Leaderboard.OtherPlayers = others
}
