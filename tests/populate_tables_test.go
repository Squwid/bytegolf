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

const workers = 25

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
			Active: randomizer.Number(0, 2) == 1,
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
	logrus.Infof("Wrote %v in %vms\n", amount, time.Since(start).Milliseconds())
}
