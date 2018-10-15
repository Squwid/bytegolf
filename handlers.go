package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Squwid/bytegolf/aws"
	"github.com/Squwid/bytegolf/runner"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, req *http.Request) {
	user := getUser(w, req)

	// If there is a post on index that means that the user is creating a new game
	if req.Method == http.MethodPost {
		if !currentlyLoggedIn(w, req) {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		// if there is already a current game
		if CurrentGame.Started && !CurrentGame.Over {
			logger.Printf("%s tried to start a new game but one already exists\n", user.Email)
			http.Redirect(w, req, "/currentgame/1", http.StatusSeeOther)
			return
		}
		game, err := CreateNewGame(w, req)
		CurrentGame = *game
		if err != nil {
			logger.Printf("an error occured creating a new game: %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, req, "/currentgame/1", http.StatusSeeOther)
		return
	}
	// if they send the get method
	tpl.ExecuteTemplate(w, "index.html", struct {
		User     *aws.User
		Name     string
		Game     Game
		LoggedIn bool
	}{
		User:     user,
		Name:     user.Email,
		Game:     CurrentGame,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

// current holds the code for the submission along with joining a current game
func current(w http.ResponseWriter, req *http.Request) {
	// if the player isnt logged in send them to the login screen
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	// the user is logged in so get the user
	user := getUser(w, req)

	// if the user is not in a game check to see if there is a current game
	if !CurrentGame.UserInGame(user) {
		if !CurrentGame.Started {
			logger.Printf("%s tried to join a game but there is not one\n", user.Email)
			http.Error(w, "there is not a current game", http.StatusNoContent)
			return
		}
		// if the user is not in a game and there is a current game add them to it
		err := CurrentGame.AddGameUser(user)
		if err != nil {
			logger.Printf("an error occured adding a user to a game: %v\n", err)
			http.Error(w, "an internal error occurred", http.StatusInternalServerError)
			return
		}
	}

	// the player is in the game now
	// hole logic
	hole := 1
	h := strings.TrimPrefix(req.URL.Path, "/currentgame/")
	// possible holes are only 1-9 for ~now~
	if len(h) == 1 {
		i, err := strconv.Atoi(h)
		hole = i
		if err != nil {
			hole = 1
		}
	} else {
		hole = 1
	}
	if hole > CurrentGame.Holes {
		http.Redirect(w, req, fmt.Sprintf("/currentgame/%v", CurrentGame.Holes), http.StatusSeeOther)
	}

	var currentCode code
	currentCode.Show = false
	// if the user is submitting a file
	if req.Method == http.MethodPost {
		// open submitted file
		lang := req.FormValue("language")
		file, fileHead, err := req.FormFile("codefile")
		if err != nil {
			logger.Printf("an error occurred opening a form file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		logger.Printf("%s uploading file to server %s\n", user.Email, fileHead.Filename)

		// read the file submitted
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Println("an error occured reading file", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Run the code given from the submission
		client := runner.NewClient()
		config := runner.NewConfiguration(Config.Storage.SaveLogs, Config.Storage.SaveSubmissions)
		sub := runner.NewCodeSubmission(user.Email, CurrentGame.Name, CurrentGame.ID, fileHead.Filename, lang, string(bs), client, config)
		resp, err := sub.Send()
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, "an unexpected error has occured", http.StatusInternalServerError)
			return
		}
		correct, err := CurrentGame.Check(resp, hole)
		if err != nil {
			logger.Printf("error checking question %v for %s : %v\n", hole, user.Email, err)
			http.Error(w, err.Error(), http.StatusBadRequest) // this has to be the users fault if the error is not nil
			return
		}
		player, err := CurrentGame.GetPlayer(user)
		if err != nil {
			logger.Printf("an error occurred getting %s from game %s: %v\n", user.Email, CurrentGame.Name, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !correct {
			currentCode.Output = resp.Output
			currentCode.Show = true
			logger.Printf("%s got hole %v incorrect\n", user.Email, hole)
		} else {
			err = CurrentGame.Score(player, hole, sub, resp)
			if err != nil {
				logger.Printf("an error occurred scoring %s submission \n", player.User.Email)
				http.Error(w, "an error occurred scoring your submission", http.StatusInternalServerError)
				return
			}
			currentCode.Output = resp.Output
			currentCode.Show = true
			currentCode.Correct = true
			logger.Printf("%s got hole %v correct!\n", user.Email, hole)
		}
	}

	q, ok := CurrentGame.Questions[hole]
	if !ok {
		http.Error(w, "unable to find that hole", http.StatusInternalServerError)
		return
	}

	player, err := CurrentGame.GetPlayer(user)
	if err != nil {
		logger.Printf("an error occurred getting %s from game %s: %v\n", user.Email, CurrentGame.Name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	oc := checkCorrect(hole, player) // get the overall code
	CurrentGame.update()
	tpl.ExecuteTemplate(w, "currentgame.html", struct {
		User        *aws.User
		Name        string
		Game        Game
		Hole        int
		Question    aws.Question
		LoggedIn    bool
		OverallCode code
		CurrentCode code
		GameOver    bool
	}{
		User:        user,
		Name:        user.Email,
		Game:        CurrentGame,
		Hole:        hole,
		Question:    q,
		LoggedIn:    currentlyLoggedIn(w, req),
		OverallCode: oc,
		CurrentCode: currentCode,
		GameOver:    CurrentGame.Over,
	})
}

func signup(w http.ResponseWriter, req *http.Request) {
	// if the user is already logged in then send them to the home screen
	if currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		reqName := req.FormValue("username")
		reqPassword := req.FormValue("password")

		// Check if username is already taken
		exist, _ := aws.UserExist(reqName)
		if exist {
			logger.Printf("user tried to register with %s but it was taken\n", reqName)
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		// username is not taken so sign them up and create session
		bs, err := bcrypt.GenerateFromPassword([]byte(reqPassword), bcrypt.MinCost)
		if err != nil {
			logger.Printf("an error occurred generating password: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		newUser := &aws.User{
			Email:    reqName,
			Password: string(bs),
			Role:     "default",
		}
		err = newUser.Store()
		if err != nil {
			logger.Println(err.Error())
			http.Error(w, "an unexpected error occured", http.StatusInternalServerError)
			return
		}

		sessionID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sessionID.String(),
		}

		http.SetCookie(w, c)
		currentSessions[c.Value] = session{
			Username:     newUser.Email,
			lastActivity: time.Now(),
		}
		logger.Printf("%s successfully signed up\n", newUser.Email)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "signup.html", struct {
		Game     Game
		LoggedIn bool
	}{
		Game:     CurrentGame,
		LoggedIn: false,
	})
}

func master(w http.ResponseWriter, req *http.Request) {
	// if the user is already logged in then send them to the home screen
	canEnd := false
	gameOver := false
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	user := getUser(w, req)

	// if the user is not in a game check to see if there is a current game
	if !CurrentGame.UserInGame(user) {
		http.Error(w, "you must be in a game", http.StatusNoContent)
		return
	}
	player, err := CurrentGame.GetPlayer(user)
	if err != nil {
		http.Error(w, "an internal server error occurred", http.StatusInternalServerError)
		logger.Printf("an error occurred getting player %s: %v", user.Email, err)
		return
	}

	if CurrentGame.Owner.User.Email == player.User.Email {
		canEnd = true
	}
	// They are the owner of the CurrentGame if the code gets here
	if req.Method == http.MethodPost {
		// TODO: change winner logic
		canEnd = false
		gameOver = true
		CurrentGame.Over = true // Ends the game in the Game objects
		CurrentGame.update()
	} else if CurrentGame.Over {
		canEnd = false
		CurrentGame.update()
	}

	tpl.ExecuteTemplate(w, "master.html", struct {
		User     *aws.User
		Name     string
		Game     Game
		LoggedIn bool
		CanEnd   bool
		GameOver bool
	}{
		User:     user,
		Name:     user.Email,
		Game:     CurrentGame,
		LoggedIn: currentlyLoggedIn(w, req),
		CanEnd:   canEnd,
		GameOver: gameOver,
	})
}

func login(w http.ResponseWriter, req *http.Request) {
	if currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/profile", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		reqName := req.FormValue("username")
		reqPass := req.FormValue("password")

		exist, user := aws.UserExist(reqName)
		if !exist {
			logger.Printf("user tried to login with %s but it does not exist\n", reqName)
			http.Error(w, "That user does not exist", http.StatusForbidden)
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqPass))
		if err != nil {
			logger.Printf("%s tried to login with incorrect password\n", user.Email)
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		// Password entered was correct so create new session
		sessionID, _ := uuid.NewV4()
		cookie := &http.Cookie{
			Name:  "session",
			Value: sessionID.String(),
		}

		http.SetCookie(w, cookie)
		currentSessions[cookie.Value] = session{
			Username:     reqName,
			lastActivity: time.Now(),
		}
		logger.Printf("%s successfully logged in\n", user.Email)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.html", struct {
		Game     Game
		LoggedIn bool
	}{
		Game:     CurrentGame,
		LoggedIn: false,
	})
}

func logout(w http.ResponseWriter, req *http.Request) {
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	user := getUser(w, req)
	sessionCookie, err := req.Cookie("session")
	if err != nil {
		http.Error(w, "an unexpected error has occured", http.StatusInternalServerError)
		return
	}
	delete(currentSessions, sessionCookie.Value)
	sessionCookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	logger.Printf("%s successfully logged out\n", user.Email)
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func profile(w http.ResponseWriter, req *http.Request) {
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	user := getUser(w, req)
	if req.Method == http.MethodPost {
		http.Redirect(w, req, "/logout", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "profile.html", struct {
		User     *aws.User
		Name     string
		LoggedIn bool
	}{
		User:     user,
		Name:     user.Email,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

func rules(w http.ResponseWriter, req *http.Request) {
	user := getUser(w, req)
	tpl.ExecuteTemplate(w, "rules.html", struct {
		User     *aws.User
		Name     string
		LoggedIn bool
	}{
		User:     user,
		Name:     user.Email,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

func leaderboard(w http.ResponseWriter, req *http.Request) {
	user := getUser(w, req)
	tpl.ExecuteTemplate(w, "leaderboards.html", struct {
		User     *aws.User
		Name     string
		LoggedIn bool
	}{
		User:     user,
		Name:     user.Email,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}

func admin(w http.ResponseWriter, req *http.Request) {
	if !currentlyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	user := getUser(w, req)

	if user.Role != "admin" {
		logger.Printf("%s tried to see /admin but has permission %s", user.Email, user.Role)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(w, "admin.html", struct {
		User     *aws.User
		Name     string
		LoggedIn bool
	}{
		User:     user,
		Name:     user.Email,
		LoggedIn: currentlyLoggedIn(w, req),
	})
}
