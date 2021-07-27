package globals

import (
	"os"
	"strings"
)

// Possible environments
const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
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
	env := strings.ToLower(os.Getenv("BG_ENV"))
	switch env {
	case EnvProd, EnvDev:
		return env
	default:
		return EnvLocal
	}
}
