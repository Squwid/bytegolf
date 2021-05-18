package globals

import (
	"os"
	"strings"
)

// Possible environments
const (
	EnvDev  = "dev"
	EnvProd = "prod"
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

func Addr() string {
	return "http://10.218.67.120"
}

func Env() string {
	env := strings.ToLower(os.Getenv("BG_ENV"))
	switch env {
	case EnvProd:
		return env
	default:
		return EnvDev
	}
}
