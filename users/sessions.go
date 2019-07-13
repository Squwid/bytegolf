package users

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	db "github.com/Squwid/bytegolf/database"
	"golang.org/x/crypto/bcrypt"
)

var (
	sessions    = map[string]*session{}
	sessionLock = &sync.RWMutex{}
)

type session struct {
	User         *BGUser
	lastActivity time.Time
}

// LoggedIn checks to see if a user is already logged in
func LoggedIn(req *http.Request) bool {
	cookie, err := req.Cookie("bgsession")
	if err != nil {
		return false
	}

	var ok bool
	sessionLock.RLock()
	defer sessionLock.RUnlock()

	if session, ok := sessions[cookie.Value]; ok {
		session.lastActivity = time.Now()
		return true
	}
	return ok
}

// logged will take a github user and log in the bytegolf user after hitting the database and
// setting a session
func (user *GithubUser) login(w http.ResponseWriter, req *http.Request) {
	bgu, err := user.exchange()
	if err != nil {
		logger.Printf("error exchaning the user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// hash the github login id to not store it raw
	uid, err := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(user.ID)), bcrypt.MinCost)
	if err != nil {
		logger.Printf("error generating cookie: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "bgsession",
		Value: string(uid),
		Path:  "/",
	}

	http.SetCookie(w, cookie)
	req.AddCookie(cookie)

	s := &session{
		User:         bgu,
		lastActivity: time.Now(),
	}

	// add the session to the list of sessions
	sessionLock.Lock()
	sessions[string(uid)] = s
	sessionLock.Unlock()

	// run a function to remove the session in a day
	go func(id string) {
		time.Sleep(time.Hour * 24)
		sessionLock.Lock()
		delete(sessions, id)
		sessionLock.Unlock()
	}(string(uid))
}

// Exchange exchanges a github user for a bytegolf user and updates the table accordingly
func (user *GithubUser) exchange() (*BGUser, error) {
	stmt, err := db.DB.Prepare("SELECT total_submissions, total_correct, total_bytes FROM user WHERE id=$1;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(user.ID)
	if err != nil {
		return nil, err
	}

	// check how many rows are returned, which should be 1 or 0 unless there is a mistake
	// also gets the bytegolf user while updating the github user
	var rs int
	var bgu = BGUser{GithubUser: user}
	for rows.Next() {
		rs++
		err = rows.Scan(&bgu.TotalSubmissions, &bgu.TotalCorrect, &bgu.TotalBytes)
		if err != nil {
			return nil, err
		}
	}
	// TODO: fix the possibility that there are more than 1 row here

	// if the player does not exist in the database add them
	if rs == 0 {
		stmt2, err := db.DB.Prepare(`INSERT INTO 
		user(id, username, picture_uri, github_uri, name, total_submissions, total_correct, total_bytes)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8);`)
		if err != nil {
			return nil, err
		}
		defer stmt2.Close()

		_, err = stmt2.Exec(user.ID, user.Username, user.PictureURI, user.GithubURI, user.Name, 0, 0, 0)
		if err != nil {
			return nil, err
		}
		bgu.TotalBytes = 0
		bgu.TotalCorrect = 0
		bgu.TotalSubmissions = 0
	} else {
		// the user does exist, but the github information needs to be updated incase it changed
		stmt2, err := db.DB.Prepare(`UPDATE user SET username=$1, picture_uri=$2, github_uri=$3, name=$4 WHERE id=$5;`)
		if err != nil {
			return nil, err
		}
		defer stmt2.Close()

		_, err = stmt2.Exec(user.Username, user.PictureURI, user.GithubURI, user.Name, user.ID)
		if err != nil {
			return nil, err
		}
	}
	return &bgu, nil
}

// GetUser gets a bytegolf user using a writer and
func GetUser(req *http.Request) (*BGUser, error) {
	if !LoggedIn(req) {
		// FIX: why does this return two nils
		return nil, nil
	}

	cookie, err := req.Cookie("bgsession")
	if err != nil {
		return nil, err
	}

	sessionLock.RLock()
	defer sessionLock.RUnlock()
	return sessions[cookie.Value].User, nil
}
