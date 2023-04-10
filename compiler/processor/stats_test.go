package processor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const CPU_FILE = `
usage_usec 40692920
user_usec 33926287
system_usec 6766633
nr_periods 0
nr_throttled 0
throttled_usec 0
nr_bursts 0
burst_usec 0
`

const MEM_FILE = `
430080
`

func TestGetFieldFromFile(t *testing.T) {
	f, err := os.CreateTemp("", "CPU_FILE")
	if err != nil {
		t.Logf("Error creating temp file: %v\n", err)
		t.FailNow()
	}
	defer os.Remove(f.Name())

	_, err = f.WriteString(CPU_FILE)
	if err != nil {
		t.Logf("Error writing to temp file: %v\n", err)
		t.FailNow()
	}
	f.Close()

	cpu, err := getFieldFromFile(f.Name(), "usage_usec")
	if err != nil {
		t.Logf("Error getting field from file: %v\n", err)
		t.FailNow()
	}

	assert.Equal(t, int64(40692920), cpu)
}

func TestGetMemoryFieldFromFile(t *testing.T) {
	f, err := os.CreateTemp("", "MEM_FILE")
	if err != nil {
		t.Logf("Error creating temp file: %v\n", err)
		t.FailNow()
	}
	defer os.Remove(f.Name())

	_, err = f.WriteString(MEM_FILE)
	if err != nil {
		t.Logf("Error writing to temp file: %v\n", err)
		t.FailNow()
	}
	f.Close()

	mem, err := getMemoryFieldFromFile(f.Name())
	if err != nil {
		t.Logf("Error getting field from file: %v\n", err)
		t.FailNow()
	}

	assert.Equal(t, int64(430080), mem)
}

func TestGetFieldFromNonExistantFile(t *testing.T) {
	cpu, err := getFieldFromFile("non-existant-file", "usage_usec")
	assert.Equal(t, int64(0), cpu)
	assert.Error(t, err)
}

func TestGetMemoryFieldFromNonExistantFile(t *testing.T) {
	mem, err := getMemoryFieldFromFile("non-existant-file")
	assert.Equal(t, int64(0), mem)
	assert.Error(t, err)
}
