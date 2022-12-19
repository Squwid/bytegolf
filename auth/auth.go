package auth

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/sqldb"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

// LoggedIn looks at the context value Claims and returns the
// Claims if the user is logged in, otherwise nil.
func LoggedIn(r *http.Request) *Claims {
	claims, ok := r.Context().Value("Claims").(*Claims)
	if !ok {
		return nil
	}
	return claims
}

func JWT(token *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}

// createOrGetDBUser returns a bytegolf user based on the github user if exists. Otherwise,
// a new user will be created, inserted into the db and returned.
func createOrGetDBUser(ctx context.Context, ghu *GithubUser) (*BytegolfUserDB, error) {
	user, err := fetchUser(ctx, ghu.GithubID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		logrus.WithField("GithubID", ghu.GithubID).
			Infof("Bytegolf user did not exist, created one.")

		user = NewBytegolfUser(*ghu)
		if _, err := sqldb.DB.NewInsert().Model(user).Exec(ctx); err != nil {
			return nil, err
		}
	} else {
		logrus.WithField("Github", ghu.GithubID).
			Infof("Found existing bytegolf user. Updating.")

		user.LastUpdatedTime = time.Now().UTC()
		user.GithubUser.Login = ghu.Login // Update to newest GitHub name.
		user.GithubUser.URL = ghu.URL     // Update github_url

		if _, err := sqldb.DB.NewUpdate().Model(user).WherePK().
			Column("updated_time").
			Column("login").
			Column("github_url").
			Exec(ctx); err != nil {
			return nil, err
		}
	}
	return user, nil
}

// fetchUser will get a BytegolfUser from the table using a gitID, will return nil, nil if one
// does not exist.
func fetchUser(ctx context.Context, gitID int64) (*BytegolfUserDB, error) {
	var user = &BytegolfUserDB{}
	err := sqldb.DB.NewSelect().Model(user).Where("id = ?", gitID).Scan(ctx)
	if err == sql.ErrNoRows {
		err = nil
		user = nil
	}
	return user, err
}

func writeJWT(w http.ResponseWriter, user *BytegolfUserDB) error {
	timeout := time.Hour * 48
	expires := time.Now().Add(timeout)
	claims := Claims{
		BGID: user.BGID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return err
	}

	// TODO: Look into making this cookie secure only for
	// non-dev environments.
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    signedToken,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}
