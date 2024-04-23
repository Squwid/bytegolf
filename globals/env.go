package globals

import (
	"os"
)

const (
	EnvLocal   = "local"
	EnvDev     = "dev"
	EnvStaging = "staging"
	EnvProd    = "prod"
)

var ENV string

func init() {
	ENV = Env()
}

// Port gets the port for the application
func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}

// Addr is the self address of the backend
func Addr() string {
	return os.Getenv("BG_BACKEND_ADDR")
}

func FrontendAddr() string {
	return os.Getenv("BG_FRONTEND_ADDR")
}

func Env() string {
	// TODO: Add additional environments for development.
	return EnvProd
}
