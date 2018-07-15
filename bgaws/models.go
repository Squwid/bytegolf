package bgaws

import "time"

// User todo
type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // todo: change this back to a byte slice + encypt
	Role     string `json:"role"`
}

// Game struct
type Game struct {
	ID          string
	Name        string
	StartedTime time.Time
	Started     bool
}
