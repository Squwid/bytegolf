package tests

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Squwid/bytegolf/api"
	"github.com/Squwid/bytegolf/sqldb"
	"github.com/Squwid/go-randomizer"
	"github.com/sirupsen/logrus"
)

const workers = 20

func TestPopulateHoles(t *testing.T) {
	const amount = 99999
	start := time.Now()
	ctx := context.Background()

	var holes = make(api.HolesDB, amount)
	for i := 0; i < amount; i++ {
		holes[i] = api.HoleDB{
			ID: fmt.Sprintf("%v_%v_%v",
				randomizer.Adjective(),
				randomizer.Noun(),
				randomizer.Number(0, 99999)),
			Name:       strings.Join(randomizer.Words(randomizer.Number(3, 8)), " "),
			Difficulty: uint8(randomizer.Number(1, 6)),
			Question:   strings.Join(randomizer.Words(randomizer.Number(10, 100)), " "),
			CreatedAt: randomizer.Date(
				// Any time in the last year
				time.Now().Add(-time.Hour*8760),
				time.Now()),
			Active:       randomizer.Number(0, 2) == 1,
			LanguageEnum: int64(randomizer.Number(1, 3)),
		}
	}

	var queue = make(chan api.HoleDB, workers)

	var wg = &sync.WaitGroup{}
	wg.Add(amount)

	for i := 0; i < workers; i++ {
		go func() {
			for hole := range queue {
				if _, err := sqldb.DB.NewInsert().
					Model(&hole).Exec(ctx); err != nil {
					logrus.WithError(err).Errorf("Error writing hole")
				}
				wg.Done()
			}
		}()
	}

	for i := 0; i < len(holes); i++ {
		queue <- holes[i]
	}

	wg.Wait()
	logrus.Infof("Wrote %v holes in %vms\n", amount, time.Since(start).Milliseconds())

	PopulateTests(ctx, holes)
}

func PopulateTests(ctx context.Context, holes api.HolesDB) {
	amount := len(holes) * 3
	start := time.Now()

	var tests = make(api.TestsDB, amount)
	for i := 0; i < amount; i++ {
		tests[i] = api.TestDB{
			Name:        strings.Join(randomizer.Words(randomizer.Number(3, 8)), " "),
			Hole:        holes[randomizer.Number(0, len(holes))].ID,
			Hidden:      randomizer.Number(0, 5) == 1,
			Active:      randomizer.Number(0, 5) == 1,
			Description: strings.Join(randomizer.Words(randomizer.Number(3, 10)), " "),
			Input:       strings.Join(randomizer.Words(randomizer.Number(3, 8)), " "),
			OutputRegex: strings.Join(randomizer.Words(randomizer.Number(3, 8)), " "),
			CreatedAt: randomizer.Date(
				// Any time in the last year
				time.Now().Add(-time.Hour*8760),
				time.Now()),
		}
	}

	var queue = make(chan api.TestDB, workers)

	var wg = &sync.WaitGroup{}
	wg.Add(amount)

	for i := 0; i < workers; i++ {
		go func() {
			for test := range queue {
				if _, err := sqldb.DB.NewInsert().
					Model(&test).Exec(ctx); err != nil {
					logrus.WithError(err).Errorf("Error writing test")
				}
				wg.Done()
			}
		}()
	}

	for i := 0; i < len(tests); i++ {
		queue <- tests[i]
	}

	wg.Wait()
	logrus.Infof("Wrote %v tests in %vms\n", amount, time.Since(start).Milliseconds())
}
