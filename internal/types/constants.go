package types

type TaskStatus string

const (
	TaskQueued   TaskStatus = "queued"
	TaskRunning  TaskStatus = "running"
	TaskAwaiting TaskStatus = "awaiting"
	TaskPassed   TaskStatus = "passed"
	TaskFailed   TaskStatus = "failed"
)
