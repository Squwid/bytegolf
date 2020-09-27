package globals

import "os"

// Environment Constants
const (
	EnvDev  = "DEV"
	EnvProd = "PROD"
)

// Port gets the port for the application
func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}

// ENV Gets the current environment
func ENV() string {
	env := os.Getenv("BG_ENV")
	switch env {
	case EnvDev, EnvProd:
		return env
	default:
		return EnvDev
	}
}

func Addr() string {
	return "http://192.168.1.158"
}
