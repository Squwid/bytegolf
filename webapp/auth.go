package main

import (
	"context"
	"net/http"

	"github.com/Squwid/bytegolf/lib/auth"
	"github.com/Squwid/bytegolf/lib/log"
	"github.com/golang-jwt/jwt"
)

func parseCookie(r *http.Request) *auth.Claims {
	c, err := r.Cookie(auth.CookieName)
	if err != nil {
		log.GetLogger().WithField("CookieName", auth.CookieName).
			Debugf("Cookie doesnt exist")
		return nil
	}

	var claims auth.Claims
	token, err := jwt.ParseWithClaims(c.Value, &claims, auth.JWT)
	if err != nil {
		log.GetLogger().WithError(err).Debugf("Couldnt parse cookie")
		if err == jwt.ErrSignatureInvalid {
			return nil
		}
		log.GetLogger().WithError(err).Errorf("Error parsing cookie")
		return nil
	}

	if !token.Valid {
		log.GetLogger().WithError(err).Debugf("Invalid cookie")
		return nil
	}

	return &claims
}

func loggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := parseCookie(r)
		r = r.WithContext(context.WithValue(r.Context(), "Claims", claims))
		next.ServeHTTP(w, r)
	})
}
