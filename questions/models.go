package questions

import (
	"encoding/json"
	"io/ioutil"
	"log"

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
	Link       string `json:"link"`
	Created    string `json:"created"`
}

// NewQuestion creates a new question with a UUID
func NewQuestion(name, question, answer, difficulty, source, link string) *Question {
	uuid, _ := uuid.NewV4()
	return &Question{
		ID:         uuid.String(),
		Name:       name,
		Question:   question,
		Answer:     answer,
		Difficulty: difficulty,
		Source:     source,
		Link:       link,
	}
}

// GetLocalQuestions returns a list of questions that are retrieved from the local file system
func GetLocalQuestions() []Question {
	var qs = []Question{}
	filelist, err := ioutil.ReadDir("./questions/questions/")
	if err != nil {
		log.Fatal(err)
	}

	for _, fileinfo := range filelist {
		if fileinfo.Mode().IsRegular() {
			contents, err := ioutil.ReadFile("./questions/questions/" + fileinfo.Name())
			if err != nil {
				log.Fatalln("fatal err:", err)
			}
			// fmt.Println("Bytes read: ", len(bytes))
			// fmt.Println("String read: ", string(contents))
			var q Question
			err = json.Unmarshal(contents, &q)
			if err != nil {
				log.Fatalln("Error unmarshaling questions:", err)
			}
			qs = append(qs, q)
		}
	}
	return qs
}

// ToMap takes a list of questions and changes it to a map
func ToMap(qs []Question) map[int]Question {
	m := make(map[int]Question)
	for i, q := range qs {
		m[i+1] = q
	}
	return m
}
