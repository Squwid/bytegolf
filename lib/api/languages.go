package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Squwid/bytegolf/lib/auth"
	"github.com/Squwid/bytegolf/lib/sqldb"
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

func GetLanguage(ctx context.Context, id string) (*LanguageDB, error) {
	// TODO: Could be out of sync between compiler and bytegolf
	// grabbing an active language. Make it an arg.
	var lang = &LanguageDB{}
	err := sqldb.DB.NewSelect().
		Model(lang).
		Where("l.id = ?", id).
		Where("l.active = true").
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return lang, nil
}
