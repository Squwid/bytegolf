package api

import "github.com/uptrace/bun"

type JobOutputDB struct {
	bun.BaseModel `bun:"table:jobs,alias:jo"`

	ID           int64  `bun:"id,pk,autoincrement,notnull"`
	SubmissionID string `bun:"submission_id,notnull"`
	TestID       int64  `bun:"test_id,notnull"`

	StdOut   string `bun:"stdout,notnull"`
	StdErr   string `bun:"stderr,notnull"`
	ExitCode int    `bun:"exit_code,notnull"`
	Duration int64  `bun:"duration,notnull"`
	Memory   int64  `bun:"memory,notnull"`
	CPU      int64  `bun:"cpu,notnull"`

	Correct bool `bun:"correct,notnull"`
}
