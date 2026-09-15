package types

type TaskStatus string

const (
	TaskQueued  TaskStatus = "queued"
	TaskRunning TaskStatus = "running"
	TaskPassed  TaskStatus = "passed"
	TaskFailed  TaskStatus = "failed"
)
