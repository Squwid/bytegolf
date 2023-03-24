package processor

import (
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types"
)

type Stats struct {
	CPU int64
	Mem int64
}

func waitAndGetContainerStats(reader io.ReadCloser, job *Job) {
	defer job.wg.Done()
	defer reader.Close()

	dec := json.NewDecoder(reader)
	var stats = &Stats{}
	for {
		var v *types.StatsJSON
		if err := dec.Decode(&v); err != nil {
			job.Stats = stats
			if err == io.EOF {
				return
			}
			job.logger.WithError(err).Errorf("Error decoding stats")
			return
		}

		if stats.CPU < int64(v.CPUStats.CPUUsage.TotalUsage) {
			stats.CPU = int64(v.CPUStats.CPUUsage.TotalUsage)
		}
		if stats.Mem < int64(v.MemoryStats.Usage) {
			stats.Mem = int64(v.MemoryStats.Usage)
		}
	}
}
