package main

import (
	"net/http"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/auth"
	"github.com/Squwid/bytegolf/lib/globals"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := sqldb.Open(); err != nil {
		logrus.WithError(err).Fatalf("Error connecting to db")
	}
	defer func() {
		if err := sqldb.Close(); err != nil {
			logrus.WithError(err).Errorf("")
		}
	}()

	port := globals.Port()
	env := globals.Env()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter,
		r *http.Request) {
		w.WriteHeader(200)
	})

	r.HandleFunc("/login/check", auth.CallbackHandler)
	r.HandleFunc("/login", auth.LoginHandler)

	r.HandleFunc("/api/holes", api.ListHolesHandler).Methods("GET")
	r.HandleFunc("/api/hole/{hole}", api.GetHoleHandler).Methods("GET")
	r.HandleFunc("/api/languages", api.ListLanguagesHandler).Methods("GET")
	r.HandleFunc("/api/submission", api.PostSubmissionHandler).Methods("POST")

	r.HandleFunc("/api/claims", auth.ShowClaims).Methods("GET") // Returns a user's claims and see if they are logged in

	logrus.WithField("Env", env).Infof("Starting container on port :%s", port)
	http.ListenAndServe(":"+port, loggedIn(cors(r)))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", globals.FrontendAddr())
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
			return
		}
		next.ServeHTTP(w, r)
	})
}
