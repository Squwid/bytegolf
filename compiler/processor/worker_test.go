package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	Init()
}

func TestInitStartsWorkers(t *testing.T) {
	assert.Equal(t, workerCount, len(WorkerPool))
}

func TestNewWorker(t *testing.T) {
	worker := NewWorker(1)

	assert.Equal(t, "w1", worker.ID)
}
