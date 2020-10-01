package jdoodle

import (
	"context"
	"regexp"
	"time"

	fs "github.com/Squwid/bytegolf/firestore"
	"github.com/Squwid/bytegolf/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const compileCollection = "compiles"

// store will create a compile object and store it in the database
func store(input models.CompileInput, output models.CompileOutput, bgid string) error {
	id := uuid.New().String()

	log := logrus.WithFields(logrus.Fields{
		"BGID":   bgid,
		"Action": "StoreCompile",
		"ID":     id,
	})

	log.Infof("Request to store compile")

	// Create a new Database Object
	db := models.CompileDB{
		ID:        id,
		Input:     input,
		Output:    output,
		BGID:      bgid,
		CreatedAt: time.Now().UTC(),
	}

	_, err := fs.Client.Collection(compileCollection).Doc(db.ID).Set(context.Background(), db)
	if err != nil {
		log.WithError(err).Errorf("Error storing compile")
	}
	return err
}

// MapLanguage takes an incoming submission and maps it to a language and version index
func MapLanguage(lang string) (language string, versionIndex string, err error) {
	// TODO: Map a language to a language and version index
	return "python3", "3", nil
}

// CheckOutput will take the input of the test and the output of the compile and match the regex
func CheckOutput(testInput models.TestCaseInput, compileOutput models.CompileOutput) (*models.TestCaseOutput, error) {
	r, err := regexp.Compile(testInput.Solution)
	if err != nil {
		return nil, err
	}

	correct := r.MatchString(compileOutput.Output)

	return &models.TestCaseOutput{
		ID:      testInput.ID,
		Output:  compileOutput.Output,
		Correct: correct,
	}, nil
}
