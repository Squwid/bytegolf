package questions

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	uuid "github.com/satori/go.uuid"
)

// DefaultRegion is the default aws region for this package
const DefaultRegion = "us-east-1"

// Questions Constants
const (
	questionsTable = "bytegolf-questions"
)

// Question is the type that gets stored as a question in dynamodb
type Question struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	Difficulty string `json:"difficulty"`
	Source     string `json:"source"`
	Created    string `json:"created"`

	// Information regarding whether or not the hole is live or not and what number the hole is
	Live bool `json:"live"`
}

// NewQuestion creates a new question with a UUID
func NewQuestion(name, question, answer, difficulty, source string, live bool) *Question {
	uuid, _ := uuid.NewV4()
	return &Question{
		ID:         uuid.String(),
		Name:       name,
		Question:   question,
		Answer:     answer,
		Difficulty: difficulty,
		Source:     source,
		Live:       live,
	}
}

// Store stores a question locally, however it does not make the question live
func (q *Question) Store() error {
	var p = path.Join("localfiles", "qs")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.MkdirAll(p, os.ModePerm)
	}

	filePath := path.Join(p, q.ID+".json")
	os.Remove(filePath) // remove the file before removing it
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	bs, err := json.Marshal(*q)
	if err != nil {
		return err
	}

	_, err = f.Write(bs)
	if err != nil {
		return err
	}
	return f.Sync()
}

// Deploy deploys a question to live, if a number already exists that is requested that hole will be removed
func (q *Question) Deploy() error {
	q.Live = true
	return q.Store()
}

// RemoveLive removes q question from being live and moves it to the
func (q *Question) RemoveLive() error {
	// remove the live information and restore it
	q.Live = false
	return q.Store()
}

// GetLocalQuestions returns a list of questions that are retrieved from the local file system
func GetLocalQuestions() ([]Question, error) {
	var qs = []Question{}
	filelist, err := ioutil.ReadDir(path.Join("localfiles", "qs"))
	if err != nil {
		if os.IsNotExist(err) {
			// create the folder if it doesnt exist
			os.MkdirAll(path.Join("localfiles", "qs"), os.ModePerm)
			return qs, nil
		}
		return nil, err
	}

	for _, fileinfo := range filelist {
		if fileinfo.Mode().IsRegular() {
			contents, err := ioutil.ReadFile(path.Join("localfiles", "qs", fileinfo.Name()))
			if err != nil {
				return nil, err
			}

			var q Question
			err = json.Unmarshal(contents, &q)
			if err != nil {
				return nil, err
			}
			qs = append(qs, q)
		}
	}
	return qs, nil
}

// GetLiveQuestions gets a list of live questions
func GetLiveQuestions() ([]Question, error) {
	var live = []Question{}
	qs, err := GetLocalQuestions()
	if err != nil {
		return nil, err
	}
	for _, q := range qs {
		if q.Live {
			live = append(live, q)
		}
	}
	return live, nil
}

// MapLiveQuestions creates a map of all live questions in a fasion of hole number -> question
func MapLiveQuestions() (map[int]Question, error) {
	m := make(map[int]Question)
	qs, err := GetLiveQuestions()
	if err != nil {
		return nil, err
	}

	for i := range qs {
		m[i] = qs[i]
	}
	return m, nil
}
