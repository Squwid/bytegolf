package api

import (
	"context"
	"database/sql"
	"io"
	"net/http"

	"github.com/Squwid/bytegolf/lib/auth"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

const maxBodySize int = 4096

type SubmissionDB struct {
	bun.BaseModel `bun:"table:submissions,alias:ss"`

	ID     string `bun:"id,pk,notnull"`
	Script string `bun:"script,notnull"`
	Hole   string `bun:"hole,notnull"`
	BGID   string `bun:"bgid,notnull"` // Player ID.
}

func (sdb SubmissionDB) Store(ctx context.Context) error {
	_, err := sqldb.DB.NewInsert().Model(&sdb).Exec(ctx)
	return err
}

// Submit submits a submission id to the compiler.
func (sdb SubmissionDB) Submit(ctx context.Context) error {
	return sqldb.PubSubClient.Publish(ctx, []byte(sdb.ID))
}

// PostSubmissionHandler is the handler to take a submission from a player, parse it, and
// send it to the compiler.
func PostSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	logger := logrus.WithField("Action", "PostSubmission")

	// claims := auth.LoggedIn(r)
	// if claims == nil {
	// 	http.Error(w, "Not logged in", http.StatusUnauthorized)
	// 	logger.Infof("Unauthorized submission from %s", r.RemoteAddr)
	// 	return
	// }
	claims := &auth.Claims{BGID: "test"}

	holeID := mux.Vars(r)["hole"]
	id := RandomString(6)
	logger = logger.WithFields(logrus.Fields{
		"User": claims.BGID,
		"Hole": holeID,
		"ID":   id,
	})
	logger.Infof("New submission received")

	// TODO: Come up with a better way to handle this.
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.WithError(err).Error("Error reading submission")
		return
	}
	defer r.Body.Close()

	// Make sure that the hole exists.
	hole, err := GetHole(ctx, holeID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.WithError(err).Error("Error retrieving hole")
		return
	} else if hole == nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Error("Hole not found")
		return
	}

	// Create a new submission in the database and submit to compiler.
	sub := SubmissionDB{
		ID:     id,
		Script: string(bs),
		Hole:   hole.ID,
		BGID:   claims.BGID,
	}
	if err := sub.Store(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.WithError(err).Error("Error storing submission")
		return
	}
	if err := sub.Submit(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.WithError(err).Error("Error submitting to pubsub")
		return
	}
}

// GetSubmission returns a submission by id. This doesnt do a BGID check yet.
func GetSubmission(ctx context.Context, id string) (*SubmissionDB, error) {
	sub := &SubmissionDB{}
	err := sqldb.DB.NewSelect().Model(sub).Where("id = ?", id).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return sub, nil
}
