package aws

// User TODO:
type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // todo: change this back to a byte slice + encypt
	Role     string `json:"role"`
}
