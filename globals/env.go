package globals

import (
	"os"
	"strings"
)

// TODO: Use user tokens instead of this hardcoded BGID
const BGID = "9581d9ef-d998-4903-b88c-5345e980770f"

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

func FrontendAddr() string {
	return "http://10.218.67.120:3000"
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
