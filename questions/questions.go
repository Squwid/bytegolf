package questions

import (
	"log"
	"os"

	db "github.com/Squwid/bytegolf/database"
)

// DefaultRegion is the default aws region for this package
const DefaultRegion = "us-east-1"

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "[questions] ", log.Ldate|log.Ltime)
}

// Question is the type that gets stored as a question in dynamodb
type Question struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Question   string `json:"question"`
	Input      string `json:"input"`
	Answer     string `json:"answer"`
	Difficulty string `json:"difficulty"`
	Source     string `json:"source"`

	// Information regarding whether or not the hole is live or not and what number the hole is
	Live bool `json:"live"`
}

// NewQuestion creates a new question with a UUID
func NewQuestion(name, question, input, answer, difficulty, source string, live bool) *Question {
	return &Question{
		Name:       name,
		Question:   question,
		Input:      input,
		Answer:     answer,
		Difficulty: difficulty,
		Source:     source,
		Live:       live,
	}
}

// Store will store the question after it is created
func (q *Question) Store() error {
	// TODO: check to see if a question already exists, update if it is instead
	stmt, err := db.DB.Prepare(`INSERT INTO question(name, question, input, answer, difficulty, source, live)
	VALUES ($1, $2, $3, $4, $5, $6, $7);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(q.Name, q.Question, q.Input, q.Answer, q.Difficulty, q.Source, q.Live)
	if err != nil {
		return err
	}

	logger.Printf("successfully stored %s as a question\n", q.Name)
	return nil
}

// RemoveQuestion removes a question by id rather than the question itself
func RemoveQuestion(id string) error {
	return removeQ(id)
}

// Remove removes the question from the database
func (q *Question) Remove() error {
	return removeQ(q.ID)
}

func removeQ(id string) error {
	stmt, err := db.DB.Prepare("DELETE FROM question WHERE id=$1;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	logger.Printf("successfully deleted question with id: %s\n", id)
	return nil
}

// ArchiveQuestion archives a question
func ArchiveQuestion(id string) error {
	return archive(id)
}

// Archive archives a question
func (q *Question) Archive() error {
	return archive(q.ID)
}

// MakeLive makes an archived question live
func (q *Question) MakeLive() error {
	return makeLive(q.ID)
}

// MakeLive makes a question live using an id rather than a question
func MakeLive(id string) error {
	return makeLive(id)
}

func makeLive(id string) error {
	stmt, err := db.DB.Prepare("UPDATE question SET live=true WHERE id=$1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	logger.Printf("successfully made live question with id: %s\n", id)
	return nil
}

func archive(id string) error {
	stmt, err := db.DB.Prepare("UPDATE question SET live=false WHERE id=$1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	logger.Printf("successfully archived question with id: %s\n", id)
	return nil
}

// GetActiveQuestions Retreive all of the live questions that are active
func GetActiveQuestions() ([]Question, error) {
	stmt, err := db.DB.Prepare("SELECT * FROM question WHERE live=true;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var qs []Question
	for rows.Next() {
		var q Question
		err = rows.Scan(&q.ID, &q.Name, &q.Question, &q.Input, &q.Answer, &q.Difficulty, &q.Source, &q.Live)
		if err != nil {
			return nil, err
		}
		qs = append(qs, q)
	}

	return qs, nil
}

// GetAllQuestions gets all of the questions from the database without using any queries
func GetAllQuestions() ([]Question, error) {
	stmt, err := db.DB.Prepare("SELECT * FROM question;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var qs []Question
	for rows.Next() {
		var q Question
		err = rows.Scan(&q.ID, &q.Name, &q.Question, &q.Input, &q.Answer, &q.Difficulty, &q.Source, &q.Live)
		if err != nil {
			return nil, err
		}
		qs = append(qs, q)
	}

	return qs, nil
}

// GetQuestionByID gets a question by using its ID
func GetQuestionByID(id string) (*Question, error) {
	stmt, err := db.DB.Prepare("SELECT * FROM question WHERE id=$1;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// rows in this case should just be 1

	var q Question
	for rows.Next() {
		err = rows.Scan(&q.ID, &q.Name, &q.Question, &q.Input, &q.Answer, &q.Difficulty, &q.Source, &q.Live)
		if err != nil {
			return nil, err
		}
	}
	return &q, nil
}
