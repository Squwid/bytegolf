package main

type SubmissionAverages struct {
	AvgDur int64
	AvgCPU int64
	AvgMem int64
}

/**
func calculateAverages(jobs []*Job) SubmissionAverages {
	if len(jobs) == 0 {
		return SubmissionAverages{}
	}

	avgDur := int64(0)
	avgCPU := int64(0)
	avgMem := int64(0)
	count := 0

	for _, job := range jobs {
		if job.Test.Benchmark {
			count++
			avgDur += job.output.Duration
			avgCPU += job.stats.CPU
			avgMem += job.stats.Mem
		}
	}
	avgDur /= int64(len(jobs))
	avgCPU /= int64(len(jobs))
	avgMem /= int64(len(jobs))
	return SubmissionAverages{
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
*/
