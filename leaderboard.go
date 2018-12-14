package main

/*
	Leaderboards are going to be stored in memory for now to make sure that they
	work and in the future they will be stored on AWS
*/

// holeScores is a dictionary of hole to list of top scores. Each hole list can only
// hold the top of each player submission and currently does not save on exit
var holeScores = map[string][]LBSingleScore{}

// LBSingleScore holds a single e lb entry for a single user on a single hole
type LBSingleScore struct {
	Username string
	Language string
	Score    int
}

func addScore(hole string, score LBSingleScore) {
	// fmt.Printf("ADDING %s's SUBMISSION TO SCORE\n", score.Username)
	currentScore, index, exist := userHasSubmission(hole, score.Username)
	if exist {
		// the user already has a submission
		if score.Score >= currentScore {
			// no need to add a score if its larger or the same as the current score
			return
		}
		holeScores[hole] = append(holeScores[hole][:index], holeScores[hole][index+1:]...)
	}
	// fmt.Println("HOLE:", holeScores[hole])
	holeScores[hole] = append(holeScores[hole], score)
	logger.Printf("added %s's submission to leaderboard\n", score.Username)
	// fmt.Println("HOLE:", holeScores[hole])
	return
}

func getTopThree(hole string) (*LBSingleScore, *LBSingleScore, *LBSingleScore) {
	var first, second, third LBSingleScore
	for _, score := range holeScores[hole] {
		// fmt.Printf("PREV SCORE: %v CURRENT HIGH SCORE: %v\n", score.Score, first.Score)
		if score.Score < first.Score || first.Score == 0 {
			third = second
			second = first
			first = score
		} else if score.Score < second.Score || second.Score == 0 {
			third = second
			second = score
		} else if score.Score < third.Score || third.Score == 0 {
			third = score
		}
	}
	return &first, &second, &third
}

func userHasSubmission(hole string, username string) (int, int, bool) {
	for i, score := range holeScores[hole] {
		if score.Username == username {
			return score.Score, i, true
		}
	}
	return 0, 0, false
}
