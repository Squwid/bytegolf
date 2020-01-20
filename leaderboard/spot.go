package leaderboard

// Spot represents a single leaderboard spot
type Spot struct {
	BGID       string `json:"bgid"`
	ExecuteID  string `json:"execute_id"` // how to find specific execute
	QuestionID string `json:"question_id"`
	Score      int    `json:"score"`
	Language   string `json:"language"`
}

func selectTopScores(hole string, max int) ([]Spot, error) {
	return nil, nil
}
