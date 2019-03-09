package runner

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
)

// GetHoleLB returns 3 leaderboard positions using the LbOverview struct for the play page
func GetHoleLB(holeID string) map[int]LbOverview {
	// var logger = log.New(os.Stdout, "[debug] ", log.Ltime)
	// logger.Printf("starting GetHoleLB\n")
	// defer logger.Printf("ending GetHoleLB\n")

	var lb = map[int]LbOverview{
		1: LbOverview{},
		2: LbOverview{},
		3: LbOverview{},
	}
	var p = path.Join("localfiles", "codefiles", holeID)

	folders, err := ioutil.ReadDir(p)
	if err != nil {
		// logger.Printf("error getting first folders : %v\n", err)
		return lb
	}
	for _, folder := range folders {
		files, err := ioutil.ReadDir(path.Join(p, folder.Name()))
		if err != nil {
			// logger.Printf("error getting second files : %v\n", err)
			return lb
		}
		if len(files) == 0 {
			// logger.Printf("no files found, continuing\n")
			continue
		}
		if !files[0].Mode().IsRegular() {
			// logger.Printf("%s file is not regular\n", files[0].Name())
			continue
		}

		// get the contents from the file
		contents, err := ioutil.ReadFile(path.Join(p, folder.Name(), files[0].Name()))
		if err != nil {
			// logger.Printf("error reading file %s : %v\n", files[0].Name(), err)
			continue
		}

		var cf CodeFile
		err = json.Unmarshal(contents, &cf)
		if err != nil {
			log.Printf("error unmarshalling codefile for %s : %v\n", files[0].Name(), err)
			continue
		}
		var newLB = LbOverview{
			Username: cf.Submission.Info.User,
			Language: cf.Submission.Language,
			Score:    cf.Length,
		}
		/* LEADERBOARD LOGIC */
		if lb[1].Score > newLB.Score || lb[1].Score == 0 {
			lb[3] = lb[2]
			lb[2] = lb[1]
			lb[1] = newLB
			continue
		}
		if lb[2].Score > newLB.Score || lb[2].Score == 0 {
			lb[3] = lb[2]
			lb[2] = newLB
			continue
		}
		if lb[3].Score > newLB.Score || lb[3].Score == 0 {
			lb[3] = newLB
			continue
		}
	}

	return lb
}

// PreviouslyAnswered checks to see the players previously answered question
func PreviouslyAnswered(holeID, user string) PrevAnswered {
	var prev PrevAnswered
	var p = path.Join("localfiles", "codefiles", holeID, user)

	// get the files in the folder
	files, err := ioutil.ReadDir(p)
	// the folder does not exist therefore they have never answered the question before
	if err != nil {
		return prev
	}

	if len(files) == 0 {
		// the folder exists but they have no solutions (dont know how this could even happen)
		return prev
	}

	if !files[0].Mode().IsRegular() {
		log.Printf("file is not regular!\n")
		return prev
	}
	contents, err := ioutil.ReadFile(path.Join(p, files[0].Name()))
	if err != nil {
		return prev
	}

	var cf CodeFile
	err = json.Unmarshal(contents, &cf)
	if err != nil {
		log.Printf("error unmarshalling codefile for %s : %v\n", user, err)
		return prev
	}

	// transfer the codefile to the prev
	prev.Language = cf.Submission.Language
	prev.Correct = cf.Correct
	prev.Score = cf.Length
	return prev
}
