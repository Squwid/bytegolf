package processor

import (
	"github.com/Squwid/bytegolf/lib/api"
)

func calculateAverages(jobs []*Job) api.SubmissionAverages {
	if len(jobs) == 0 {
		return api.SubmissionAverages{}
	}

	avgDur := int64(0)
	avgCPU := int64(0)
	avgMem := int64(0)
	count := 0

	for _, job := range jobs {
		if job.Test.Benchmark {
			count++
			avgDur += job.Output.Duration
			avgCPU += job.Output.CPU
			avgMem += job.Output.Memory
		}
	}
	avgDur /= int64(len(jobs))
	avgCPU /= int64(len(jobs))
	avgMem /= int64(len(jobs))
	return api.SubmissionAverages{
		AvgDur: avgDur,
		AvgCPU: avgCPU,
		AvgMem: avgMem,
	}
}

// checkIfPassed returns true if all jobs are correct,
// and the amount of correct jobs.
func checkIfPassed(jobs []*Job) (bool, int) {
	var correct = 0
	for _, job := range jobs {
		if job.correct {
			correct++
		}
	}
	return correct == len(jobs), correct
}
