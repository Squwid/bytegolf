package globals

import "os"

// Returns true if mode is debug.
func IsDebugMode() bool {
	return os.Getenv("BG_LEVEL") == "debug"
}
