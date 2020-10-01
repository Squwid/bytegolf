package holes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Squwid/bytegolf/auth"
	"github.com/Squwid/bytegolf/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetHoleHandler gets a single hole by its id
// User doesnt have to be signed in
func GetHoleHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// TODO: Check if user is logged in to show inactive holes?

	log := logrus.WithFields(logrus.Fields{
		"ID":     id,
		"Action": "GetHole",
		"IP":     r.RemoteAddr,
	})

	// Get the hole using the id
	hole, err := GetHole(id, false)
	if err != nil {
		log.WithError(err).Errorf("Error getting hole")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if hole == nil {
		// Hole wasnt found
		log.Warnf("Hole was not found")
		http.Error(w, "Hole not found", http.StatusNotFound)
		return
	}

	// Marshal hole and return to user
	bs, err := json.Marshal(*hole)
	if err != nil {
		log.WithError(err).Errorf("Error marshalling hole")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Infof("Got hole")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

// ListHolesHandler gets all of the active holes
// User doesnt have to be signed in
func ListHolesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if user is logged in to show inactive holes?
	var onlyActive = true

	log := logrus.WithFields(logrus.Fields{
		"Action":     "ListHoles",
		"IP":         r.RemoteAddr,
		"OnlyActive": onlyActive,
	})

	// Get the holes
	holes, err := GetHoles(onlyActive)
	if err != nil {
		log.WithError(err).Errorf("Error getting holes")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Count", fmt.Sprintf("%v", len(holes)))

	// Marshal holes
	bs, err := json.Marshal(holes)
	if err != nil {
		log.WithError(err).Errorf("Error marshalling holes")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Infof("Listed %v holes", len(holes))
	w.Write(bs)
}

// StoreHoleHandler takes a hole from the post request body and stores it if you have admin perms
// User has to be signed in to check their role
func StoreHoleHandler(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"Action": "StoreHole",
		"IP":     r.RemoteAddr,
	})

	loggedIn, claims := auth.LoggedIn(r)
	if !loggedIn {
		log.Warnf("User not logged in")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log = log.WithFields(logrus.Fields{
		"BGID": claims.BGID,
		"Role": claims.Role,
	})

	// Check if user has access
	if !claims.Role.CanCreateHole() {
		log.Warnf("User does not have sufficient permissions")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log.Infof("User adding new hole")

	// parse hole
	var hole models.HoleDB
	if err := json.NewDecoder(r.Body).Decode(&hole); err != nil {
		log.WithError(err).Errorf("Error unmarshalling request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate hole a little
	if hole.Hole.Difficulty == "" {
		hole.Hole.Difficulty = "easy"
	}
	if hole.Hole.Name == "" || hole.Hole.Question == "" {
		log.Warnf("Hole name or question were blank")
		http.Error(w, "Hole name or question cannot be null", http.StatusBadRequest)
		return
	}

	// Check if hole has an id
	if hole.Hole.ID == "" {
		hole.Hole.ID = models.HoleTitle(hole.Hole.Name)
	}

	// See if id is taken already, if it is, give a uuid
	checkHole, err := getDBHole(hole.Hole.ID)
	if err != nil {
		log.WithError(err).Errorf("Error checking if hole %s is taken", hole.Hole.ID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hole already exists, set to uuid
	if checkHole != nil {
		hole.Hole.ID = uuid.New().String()
	}

	hole.Active = true
	hole.CreatedAt = time.Now().UTC()
	hole.CreatedBy = claims.BGID

	// Iterate over tests and add an id to each
	for i := range hole.Tests {
		hole.Tests[i].ID = uuid.New().String()
	}

	// Store the hole
	if err := storeDBHole(&hole); err != nil {
		log.WithError(err).Errorf("Error storing hole db object")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof("User successfully stored hole %s (%s)", hole.Hole.Name, hole.Hole.ID)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"success": "stored %s"}`, hole.Hole.ID)))
}

func AddTestCaseHandler(w http.ResponseWriter, r *http.Request) {
	hole := mux.Vars(r)["hole"]

	log := logrus.WithFields(logrus.Fields{
		"Hole":   hole,
		"Action": "AddTest",
		"IP":     r.RemoteAddr,
	})

	loggedIn, claims := auth.LoggedIn(r)
	if !loggedIn {
		log.Warnf("User tried to add a test case but not logged in")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log = log.WithFields(logrus.Fields{
		"BGID": claims.BGID,
		"Role": claims.Role,
	})

	if !claims.Role.CanCreateHole() {
		log.Warnf("User does not have sufficient permissions")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log.Infof("Adding new test case")

	// Parse test input
	var test models.TestCaseInput
	if err := json.NewDecoder(r.Body).Decode(&test); err != nil {
		log.WithError(err).Errorf("Error decoding test input")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Dont validate test input, because some input ~could~ be blank
	test.ID = uuid.New().String() // Overwrite ID
	// TODO: This
}

// EditHoleHandler allows for the editing of a hole
// User has to be signed in to check their role
func EditHoleHandler(w http.ResponseWriter, r *http.Request) {

}

func AdminListHolesDB(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"Action": "AdminListHoles",
		"IP":     r.RemoteAddr,
	})

	loggedIn, claims := auth.LoggedIn(r)
	if !loggedIn {
		log.Warnf("User not logged in")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log = log.WithFields(logrus.Fields{
		"BGID": claims.BGID,
		"Role": claims.Role,
	})

	// Check if user has access
	if !claims.Role.CanListAdminHoles() {
		log.Warnf("User has insufficient permissions")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log.Infof("Listing admin holes")

	holes, err := getAllDBHoles(false)
	if err != nil {
		log.WithError(err).Errorf("Error getting all db holes")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(holes)
	if err != nil {
		log.WithError(err).Errorf("error marshalling holes")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Count", fmt.Sprintf("%v", len(holes)))
	w.Write(bs)
}
