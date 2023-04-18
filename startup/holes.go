package main

import (
	"context"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

var codeOutline = `
// You are given a scanner that scans a file line by line.
// Your task is to write a function that returns the nth most frequently occurring word in the file,
// and how many times it occurs in the file.
func mostFrequentWords(scanner *bufio.Scanner, top int) (string, int) {
    return "", 0
}
`

func populateHoles() {
	var holes = []api.HoleDB{
		{
			ID:         "word_counter",
			Name:       "Word Count",
			Difficulty: 1,
			Question: `You are given a large text file, called input.txt, which contains a large amount of text. Your task is to write a program that reads in the contents of the file and finds the most frequently occurring word in the text.

			Your program should handle punctuation and capitalization appropriately. For example, "Hello" and "hello" should be treated as the same word.`,
			CreatedAt:    time.Now().UTC(),
			Active:       true,
			LanguageEnum: 1,
			CodeOutline:  codeOutline,
		},
	}
	ctx := context.Background()

	for _, h := range holes {
		if _, err := sqldb.DB.NewInsert().
			Model(&h).
			On("CONFLICT (id) DO UPDATE").
			Exec(ctx); err != nil {
			logrus.WithError(err).Fatal("failed to insert hole")
		}
	}
	logrus.Infof("Successfully inserted holes.")
}

var goBoilerplate = `
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	const input = "input.txt"

	file, err := os.Open(input)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	fifthWord, i := mostFrequentWords(scanner, 5)
	fmt.Printf("%v:%v\n", fifthWord, i)
}

{{.UserSolution}}
`

func populateTests() {
	var tests = []api.TestDB{
		{
			Name:        "The Corgs",
			Hole:        "word_counter",
			Hidden:      false,
			Description: "N/A",
			Benchmark:   false,
			Active:      true,
			InputFile:   "input1.txt",
			OutputRegex: `\bfriends:5\b`,
			CreatedAt:   time.Now().UTC(),
			Boilerplate: goBoilerplate,
		},
		{
			Name:        "Queen of Corgis",
			Hole:        "word_counter",
			Hidden:      true,
			Description: "N/A",
			Benchmark:   false,
			Active:      true,
			InputFile:   "input2.txt",
			OutputRegex: `\bdaisy:16\b`,
			CreatedAt:   time.Now().UTC(),
			Boilerplate: goBoilerplate,
		},
		{
			Name:        "The Big Corgi",
			Hole:        "word_counter",
			Hidden:      true,
			Description: "N/A",
			Benchmark:   true,
			Active:      true,
			InputFile:   "input3.txt",
			OutputRegex: `\bhungry:18753\b`,
			CreatedAt:   time.Now().UTC(),
			Boilerplate: goBoilerplate,
		},
	}
	ctx := context.Background()

	for _, t := range tests {
		// Get test by name and hole to see if it already exists
		exists, err := sqldb.DB.NewSelect().
			Model((*api.TestDB)(nil)).
			Where("name = ?", t.Name).
			Where("hole = ?", t.Hole).
			Exists(ctx)
		if err != nil {
			logrus.WithError(err).Fatal("failed to check if test exists")
		}
		if exists {
			continue
		}

		if _, err := sqldb.DB.NewInsert().
			Model(&t).
			Exec(ctx); err != nil {
			logrus.WithError(err).Fatal("failed to insert test")
		}
	}
	logrus.Infof("Successfully inserted %v tests.", len(tests))
}
