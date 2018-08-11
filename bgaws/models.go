package bgaws

// User TODO:
type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // todo: change this back to a byte slice + encypt
	Role     string `json:"role"`
}

// Question TODO:
type Question struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	Difficulty string `json:"difficulty"`
	Source     string `json:"source"`
}
