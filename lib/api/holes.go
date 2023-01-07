package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Squwid/bytegolf/lib/auth"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

var difficulties = map[uint8]string{
	1: "BASIC",
	2: "EASY",
	3: "NORMAL",
	4: "HARD",
	5: "IMPOSSIBLE",
}

type HolesDB []HoleDB

type HoleDB struct {
	bun.BaseModel `bun:"table:holes,alias:h"`

	LanguageDB *LanguageDB `bun:"rel:has-one,join:language_enum=id"`
	TestsDB    []*TestDB   `bun:"rel:has-many,join:id=hole"`

	ID           string    `bun:"id,pk,notnull"`
	Name         string    `bun:"name,notnull"`
	Difficulty   uint8     `bun:"difficulty,notnull"`
	Question     string    `bun:"question,notnull"`
	CreatedAt    time.Time `bun:"created_at,notnull"`
	Active       bool      `bun:"active,notnull"`
	LanguageEnum int64     `bun:"language_enum,notnull"`
}

type HoleClient struct {
	ID         string    `json:"ID"`
	Name       string    `json:"Name"`
	Difficulty string    `json:"Difficulty"`
	Question   string    `json:"Question"`
	CreatedAt  time.Time `json:"CreatedAt"`
	Active     bool      `json:"Active"`

	Tests []TestClient `json:"Tests,omitempty"`

	Language LanguageClient `json:"Language"`
}

func (hdb HoleDB) toClient() HoleClient {
	// Gather test cases suitable for the client.
	var tests = []TestClient{}
	if hdb.TestsDB != nil {
		for _, tdb := range hdb.TestsDB {
			if tdb.Active {
				tests = append(tests, tdb.toClient())
			}
		}
	}

	return HoleClient{
		ID:         hdb.ID,
		Name:       hdb.Name,
		Difficulty: difficulties[hdb.Difficulty],
		Question:   hdb.Question,
		CreatedAt:  hdb.CreatedAt,
		Active:     hdb.Active,
		Language:   hdb.LanguageDB.toClient(),
		Tests:      tests,
	}
}

func ListHolesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	logger := logrus.WithField("Action", "ListHoles")

	claims := auth.LoggedIn(r)
	if claims != nil {
		logger = logger.WithField("User", claims.BGID)
	}

	var offset = 0
	if i, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
		offset = i
		logger = logger.WithField("Offset", offset)
	}

	var holes = HolesDB{}
	if err := sqldb.DB.NewSelect().Model(&holes).
		Column("h.*").
		Relation("LanguageDB").
		Limit(20).
		Order("created_at DESC").
		Where("h.active = true").
		Offset(offset).
		Scan(ctx); err != nil {
		logger.WithError(err).Errorf("Error getting holes")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var clientHoles = make([]HoleClient, len(holes))
	for i := range holes {
		clientHoles[i] = holes[i].toClient()
	}

	bs, _ := json.Marshal(clientHoles)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)

	logger.Infof("Fetched %v holes.", len(clientHoles))
}

func GetHoleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	holeID := mux.Vars(r)["hole"]
	logger := logrus.WithField("Action", "GetHole").
		WithField("HoleID", holeID)

	claims := auth.LoggedIn(r)
	if claims != nil {
		logger = logger.WithField("User", claims.BGID)
	}

	hole, err := GetHole(ctx, holeID)
	if err != nil {
		logger.WithError(err).Errorf("Error getting hole")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if hole == nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Warnf("Hole not found")
		return
	}

	bs, _ := json.Marshal(hole.toClient())
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)

	logger.Infof("Fetched hole.")
}

// GetHole returns a hole from the database, or nil if it doesn't exist.
func GetHole(ctx context.Context, id string) (*HoleDB, error) {
	var hole = &HoleDB{}
	err := sqldb.DB.NewSelect().
		Model(hole).
		Where("h.id = ?", id).
		Where("h.active = true").
		Relation("LanguageDB").
		Relation("TestsDB").
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// TODO: Eventually alllow fetching of non-active
	// tests and holes for admins to enable.
	var tests = []*TestDB{}
	for _, test := range hole.TestsDB {
		if test.Active {
			tests = append(tests, test)
		}
	}
	hole.TestsDB = tests

	return hole, nil
}
