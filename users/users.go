package users

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // postgres driver
)

// variables that are needed for the database
var (
	dbUsername string
	dbPassword string
	dbHost     string
	dbName     string
	dbPort     string
)

var db *sql.DB
var logger *log.Logger

// initializes the database using connection strings
func init() {
	logger = log.New(os.Stdout, "[users] ", log.Ldate|log.Ltime)
	dbUsername = os.Getenv("db_username")
	dbPassword = os.Getenv("db_password")
	dbHost = os.Getenv("db_host")
	dbName = os.Getenv("db_name")
	dbPort = os.Getenv("db_port")
	if dbUsername == "" {
		logger.Panic("db_username ENV NOT SET")
		logger.p
	}
	if dbPassword == "" {
		logger.Panic("db_password ENV NOT SET")
	}
	if dbHost == "" {
		logger.Panic("db_host ENV NOT SET")
	}
	if dbName == "" {
		logger.Panic("db_name ENV NOT SET")
	}
	if dbPort == "" {
		dbPort = "5432" // default postgres port
	}

	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s", dbUsername, dbPassword, dbHost, dbName, dbPort)
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Panic(err)
	}
	db = database
}

// GithubUser is the user structure that uses github
type GithubUser struct {
	ID         int    `json:"id"`
	Username   string `json:"login"`
	PictureURI string `json:"avatar_url"`
	GithubURI  string `json:"html_url"`
	Name       string `json:"name"`
}

// BGUser is the type that gets stored in the database and used for the profile page
type BGUser struct {
	GithubUser       *GithubUser `json:"githubUser"`
	TotalSubmissions int         `json:"totalSubmissions"`
	TotalCorrect     int         `json:"totalCorrect"`
	TotalBytes       int         `json:"totalBytes"`
}

// NewBGUser creates a new BytegolfUser from the Github User
func NewBGUser(ghu *GithubUser) *BGUser {
	// TODO: this needs to grab the user from the database rather than just create a new one
	return &BGUser{
		GithubUser:       ghu,
		TotalSubmissions: 0,
		TotalCorrect:     0,
		TotalBytes:       0,
	}
}

// Exists checks to see if a github user already exists in the database
func (ghu *GithubUser) Exists() (bool, error) {
	stmt, err := db.Prepare("SELECT * FROM users WHERE user_id=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(ghu.ID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return true, nil
}
