package main

import (
	"context"
	"time"

	"github.com/Squwid/bytegolf/lib/api"
	"github.com/Squwid/bytegolf/lib/sqldb"
	"github.com/sirupsen/logrus"
)

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

func populateTests() {
	var tests = []api.TestDB{
		{
			Name:        "The Corgs",
			Hole:        "word_counter",
			Hidden:      false,
			Description: "N/A",
			Benchmark:   false,
			Active:      true,
			Input:       "input1.txt",
			OutputRegex: `\bfriends:5\b`,
			CreatedAt:   time.Now().UTC(),
		},
		{
			Name:        "Queen of Corgis",
			Hole:        "word_counter",
			Hidden:      true,
			Description: "N/A",
			Benchmark:   false,
			Active:      true,
			Input:       "input2.txt",
			OutputRegex: `\bdaisy:16\b`,
			CreatedAt:   time.Now().UTC(),
		},
		{
			Name:        "The Big Corgi",
			Hole:        "word_counter",
			Hidden:      true,
			Description: "N/A",
			Benchmark:   true,
			Active:      true,
			Input:       "input3.txt",
			OutputRegex: `\bhungry:18753\b`,
			CreatedAt:   time.Now().UTC(),
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
