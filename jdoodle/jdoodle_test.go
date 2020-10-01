package jdoodle

import (
	"encoding/json"
	"testing"

	"github.com/Squwid/bytegolf/models"
)

func TestSendJdoodle(t *testing.T) {
	in := models.CompileInput{
		Code:         `print("Hello, world!`,
		Language:     "python3",
		VersionIndex: "3",
	}

	out, err := SendJdoodle(in)
	if err != nil {
		t.Logf("Error sending jdoodle: %v\n", err)
		t.FailNow()
	}

	t.Logf("Got output\n")

	bs, err := json.Marshal(out)
	if err != nil {
		t.Logf("Error marshalling output: %v\n", err)
		t.FailNow()
	}

	// Test store as well
	if err := store(in, *out, "abc-123"); err != nil {
		t.Logf("Error storing: %v\n", err)
		t.FailNow()
	}

	t.Logf("OUTPUT\n")
	t.Logf("%v\n", string(bs))
}
