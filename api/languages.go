package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/sqldb"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type LanguageDB struct {
	bun.BaseModel `bun:"table:languages,alias:l"`

	ID        int64  `bun:"id,pk,autoincrement"`
	Language  string `bun:"language,notnull"`
	Version   string `bun:"version,notnull"`
	Image     string `bun:"image,notnull"`
	Active    bool   `bun:"active,notnull"`
	Cmd       string `bun:"cmd,notnull"`
	Extension string `bun:"extension,notnull"`
}

type LanguageClient struct {
	Language string `json:"Language"`
	Version  string `json:"Version"`
}

func (ldb LanguageDB) toClient() LanguageClient {
	return LanguageClient{
		Language: ldb.Language,
		Version:  ldb.Version,
	}
}

// Rest handler to return LanguageClient struct for all active
// languages in the database.
func ListLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	logger := logrus.WithField("Action", "ListLanguages")

	claims := auth.LoggedIn(r)
	if claims != nil {
		logger = logger.WithField("User", claims.BGID)
	}

	var langs []LanguageDB
	if err := sqldb.DB.NewSelect().Model(&langs).
		Where("active = true").
		Scan(ctx); err != nil {
		logger.WithError(err).Error("Error retrieving languages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var clientLangs = make([]LanguageClient, len(langs))
	for i := range langs {
		clientLangs[i] = langs[i].toClient()
	}

	bs, _ := json.Marshal(clientLangs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)

	logger.Infof("Fetched %v languages.", len(clientLangs))
}
