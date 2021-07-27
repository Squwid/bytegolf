package holes

import (
	"regexp"

	"github.com/Squwid/bytegolf/db"
	"github.com/Squwid/bytegolf/models"

	"cloud.google.com/go/firestore"
	"github.com/mitchellh/mapstructure"
)

func GetTests(hole string) (Tests, error) {
	query := db.TestSubCollection(hole).Where("Active", "==", true).OrderBy("CreatedAt", firestore.Desc).Limit(100)
	docs, err := db.Query(models.NewQuery(query, nil))
	if err != nil {
		return nil, err
	}

	var tests Tests
	if err := mapstructure.Decode(docs, &tests); err != nil {
		return nil, err
	}

	return tests, nil
}

// Check is a function to see if the output given is correct in the scope of the test
func (test Test) Check(output string) (bool, error) {
	r, err := regexp.Compile(test.OutputRegex)
	if err != nil {
		return false, err
	}

	matches := r.MatchString(output)
	return matches, nil
}
