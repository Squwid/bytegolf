package submissions

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	_ "github.com/Squwid/bytegolf/firestore"
	"github.com/google/uuid"
)

func TestStoreSubmissions(t *testing.T) {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	i1 := r1.Intn(1000) + 1

	fmt.Println("rand:", i1)

	fullSub := FullSubmission{
		Short: ShortSubmission{
			UUID:          uuid.New().String(),
			BGID:          "3",
			Correct:       true,
			HoleID:        "1",
			Language:      "Golang",
			Length:        i1,
			SubmittedTime: time.Now().UTC(),
		},
	}

	if err := fullSub.Store(); err != nil {
		t.Errorf("Error: %v", err)
		t.Fail()
	}

	fmt.Println("STORED")
}
