package questions

import (
	"fmt"
	"testing"

	db "github.com/Squwid/bytegolf/database"
	"github.com/kr/pretty"
	_ "github.com/lib/pq" // postgres driver
)

func TestDeleteQuestion(t *testing.T) {
	t.Skip("this should only be run by hand")
	err := RemoveQuestion("1")
	if err != nil {
		t.Errorf("error removing question: %v\n", err)
		return
	}
}

func TestGetQuestionByID(t *testing.T) {
	qs, err := GetAllQuestions()
	if err != nil {
		t.Errorf("error getting all qs: %v\n", err)
		return
	}

	if len(qs) == 0 {
		t.Errorf("expected more than a question but got none")
		return
	}

	q, err := GetQuestionByID(qs[0].ID)
	if err != nil {
		t.Errorf("error getting question by id: %v\n", err)
		return
	}

	pretty.Println(*q)
}

func TestGetLiveQuestions(t *testing.T) {
	qs, err := GetActiveQuestions()
	if err != nil {
		t.Errorf("error getting live qs: %v\n", err)
		return
	}

	if len(qs) == 0 {
		t.Errorf("expected more than a question but got none")
		return
	}
	pretty.Println(qs)
}

func TestGetAllQuestions(t *testing.T) {
	qs, err := GetAllQuestions()
	if err != nil {
		t.Errorf("error getting all qs: %v\n", err)
		return
	}

	if len(qs) == 0 {
		t.Errorf("expected more than a question but got none")
		return
	}
	pretty.Println(qs)
}

func TestGetDatabase(t *testing.T) {
	stmt, err := db.DB.Prepare("SELECT * FROM question;")
	if err != nil {
		t.Errorf("could not prepare statement: %v\n", err)
		return
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		t.Errorf("could not query: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("*** ROWS:")
	fmt.Println(rows)
}

func TestCreateTable(t *testing.T) {
	t.Skip("Should only be run by hand")
	stmt, err := db.DB.Prepare(`
	CREATE TABLE IF NOT EXISTS question (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		question TEXT NOT NULL,
		input VARCHAR(200) NOT NULL,
		answer VARCHAR(200) NOT NULL,
		difficulty VARCHAR(100) NOT NULL,
		source VARCHAR(200) NULL,
		live BOOLEAN NOT NULL
	);`)
	if err != nil {
		t.Errorf("could not prepare: %v\n", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Query()
	if err != nil {
		t.Errorf("could not query: %v\n", err)
		return
	}
}

func TestStoreQuestions(t *testing.T) {
	t.Skip("Should only be run by hand")
	q := NewQuestion("name", "question", "input", "answer", "difficulty", "source", true)
	err := q.Store()
	if err != nil {
		t.Errorf("error storing a test question: %v\n", err)
	}
}

func TestDeleteTable(t *testing.T) {
	t.Skip("Should only be run by hand")
	stmt, err := db.DB.Prepare("DROP TABLE IF EXISTS question;")
	if err != nil {
		t.Errorf("could not prepare statment: %v\n", err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		t.Errorf("error deleting table: %v\n", err)
		return
	}
}
