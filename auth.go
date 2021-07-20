package main

import (
	"context"
	"net/http"

	"github.com/Squwid/bytegolf/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

func parseCookie(r *http.Request) *auth.Claims {
	c, err := r.Cookie(auth.CookieName)
	if err != nil {
		return nil
	}

	var claims auth.Claims
	token, err := jwt.ParseWithClaims(c.Value, &claims, auth.JWT)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil
		}
		logrus.WithError(err).Errorf("Error parsing cookie")
		return nil
	}

	if !token.Valid {
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
