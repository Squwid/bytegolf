package api

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/lib/auth"
	"github.com/Squwid/bytegolf/lib/comms"
	"github.com/Squwid/bytegolf/lib/log"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type SubmissionDB struct {
	bun.BaseModel `bun:"table:submissions,alias:ss"`

	ID           string    `bun:"id,pk,notnull"`
	Script       string    `bun:"script,notnull"`
	Hole         string    `bun:"hole,notnull"`
	BGID         string    `bun:"bgid,notnull"` // Player ID.
	ScriptHash   string    `bun:"hash,notnull"`
	Status       int       `bun:"status,notnull"`
	CreatedTime  time.Time `bun:"created_time,notnull"`
	CompiledTime time.Time `bun:"compiled_time"`
	Length       int       `bun:"length,notnull"`

	// Averages
	AvgDur int64 `bun:"avg_dur"`
	AvgCPU int64 `bun:"avg_cpu"`
	AvgMem int64 `bun:"avg_mem"`
	Passed bool  `bun:"passed"`
}

const (
	StatusPending = 0
	StatusRunning = 1
	StatusSuccess = 2
	StatusFailure = 3

	StatusPendingStr = "PENDING"
	StatusRunningStr = "RUNNING"
	StatusSuccessStr = "SUCCESS"
	StatusFailureStr = "FAILURE"
)

// Map of submission status to string.
var StatusMap = map[int]string{
	StatusPending: StatusPendingStr,
	StatusRunning: StatusRunningStr,
	StatusSuccess: StatusSuccessStr,
	StatusFailure: StatusFailureStr,
}

func UpdateSubmissionStatus(ctx context.Context, id string, status int) error {
	_, err := sqldb.DB.NewUpdate().
		Model(&SubmissionDB{}).
		Set("status = ?", status).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func CompleteSubmission(ctx context.Context, id string, passed bool,
	avgs SubmissionAverages) error {
	_, err := sqldb.DB.NewUpdate().
		Model(&SubmissionDB{}).
		Set("status = ?", StatusSuccess).
		Set("passed = ?", passed).
		Set("avg_cpu = ?", avgs.AvgCPU).
		Set("avg_mem = ?", avgs.AvgMem).
		Set("avg_dur = ?", avgs.AvgDur).
		Set("compiled_time = ?", time.Now().UTC()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (sdb SubmissionDB) Store(ctx context.Context) error {
	_, err := sqldb.DB.NewInsert().Model(&sdb).Exec(ctx)
	return err
}

// Submit submits a submission id to the compiler.
func (sdb SubmissionDB) Submit(ctx context.Context) error {
	return comms.PublisherImpl.Publish(ctx, []byte(sdb.ID))
}

// PostSubmissionHandler is the handler to take a submission from a player, parse it, and
// send it to the compiler.
func PostSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	logger := log.GetLogger().WithField("Action", "PostSubmission")

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
		ID:          id,
		Script:      string(bs),
		Hole:        hole.ID,
		BGID:        claims.BGID,
		ScriptHash:  hash(string(bs)),
		CreatedTime: time.Now().UTC(),
		Length:      length(bs),
		Status:      0,
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

func hash(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func length(bs []byte) int {
	return len(bs)
}