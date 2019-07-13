package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // postgres driver
)

var logger *log.Logger

// DB is the database connection
var DB *sql.DB

// Database connection strings
var (
	dbUsername string
	dbPassword string
	dbHost     string
	dbName     string
	dbPort     string
)

func init() {
	logger = log.New(os.Stdout, "[db] ", log.Ldate|log.Ltime)
	dbUsername = os.Getenv("db_username")
	dbPassword = os.Getenv("db_password")
	dbHost = os.Getenv("db_host")
	dbName = os.Getenv("db_name")
	dbPort = os.Getenv("db_port")
	if dbUsername == "" {
		logger.Panic("db_username ENV NOT SET")
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

	dbName = "bytegolf"
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s", dbUsername, dbPassword, dbHost, dbName, dbPort)
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		// the database doesnt connect so panic ()
		panic(fmt.Sprintf("error connecting do database: %v", err))
	} else {
		DB = database
		logger.Printf("successfully connected to database\n")
	}
}

// InProd returns a boolean depending if the application is in the cloud or local
func InProd() bool {
	return os.Getenv("prod") == "true"
}
