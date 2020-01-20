package sess

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Squwid/bytegolf/firestore"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

const collection = "sessions"
const sessionID = "bg-session"

// Session represents a users session. A permission level change will make the current users
// session inactive
type Session struct {
	ID      string `json:"id"`
	BGID    string `json:"bg_id"`
	Timeout int64  `json:"timeout"`

	PermissionLevel string `json:"permission_level"`
}

// GetSession returns a session, if error is ErrNotFound that means the session
// does not exist at all
func GetSession(id string) (*Session, error) {
	return retreiveSess(id)
}

// ErrNotFound is an error that is returned when a doc is not found in firestore
var ErrNotFound = errors.New("session not found")

func retreiveSess(id string) (*Session, error) {
	ctx := context.Background()
	ref, err := firestore.Client.Collection(collection).Doc(id).Get(ctx)
	if err != nil && strings.Contains(err.Error(), "code = NotFound") {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	var sess Session
	err = mapstructure.Decode(ref.Data(), &sess)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// Put creates a new sesssion in firestore
func (sess Session) Put() error {
	ctx := context.Background()
	_, err := firestore.Client.Collection(collection).Doc(sess.ID).Set(ctx, sess)
	if err != nil {
		return err
	}
	return nil
}

// removeWhere removes all of a specific set of entries
func removeWhere(path, op string, value interface{}) error {
	ctx := context.Background()
	// iter := firestore.Client.Collection(collection).Where("Username", "==", username).Documents(ctx)
	iter := firestore.Client.Collection(collection).Where(path, op, value).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil
		}
		ctx2 := context.Background()
		_, err = doc.Ref.Delete(ctx2)
		if err != nil {
			return err
		}
	}
	return nil
}

// Login logs a user in, and returns the session and any other errors that occurred
// you need to check if you are logged in before using this function
func Login(bgID string) (*Session, error) {
	uid := uuid.New().String()
	s := Session{
		ID:      uid,
		BGID:    bgID,
		Timeout: time.Now().Local().Add(time.Hour * 24).Unix(),
	}
	if err := removeWhere("BGID", "==", bgID); err != nil {
		log.Errorf("error removing old bgID: %v", err)
	} // delete all old sessions
	if err := removeWhere("Timeout", "<", time.Now().Local()); err != nil {
		log.Errorf("error removing timeouts: %v", err)
	}
	// TODO: also should put a delete old sessions here to remove all the old junk
	// HAHA its there now
	return &s, s.Put() // add the session and return it
}

// LoggedIn checks if a user is logged in using the incoming request
// ONLY use the session here if they are logged in AND there are no errors
func LoggedIn(req *http.Request) (bool, *Session, error) {
	cookie, err := req.Cookie(sessionID)
	if err != nil {
		return false, nil, nil
	}
	s, err := retreiveSess(cookie.Value)
	if err == ErrNotFound {
		// the user was not found
		return false, nil, nil
	} else if err != nil {
		return false, nil, err
	}
	if s == nil || s.Timeout < time.Now().Local().Unix() {
		return false, nil, nil
	}
	return true, s, nil
}
