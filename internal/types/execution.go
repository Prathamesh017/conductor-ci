package types

import "time"

type RunResult struct {
	TaskName string
	Success  bool
	Output   string
	ExitCode int
	Duration time.Duration
	Err      string
}

type ExecutionState struct {
	WorkflowStatus TaskStatus
	StageStatus    map[string]TaskStatus
	TaskStatus     map[string]TaskStatus
	TaskDuration   map[string]time.Duration
}

func QueuedState(cfg WorkflowConfig) ExecutionState {
	state := ExecutionState{
		WorkflowStatus: TaskRunning,
		StageStatus:    make(map[string]TaskStatus),
		TaskStatus:     make(map[string]TaskStatus),
		TaskDuration:   make(map[string]time.Duration),
	}
	for _, stage := range cfg.Execution {
		state.StageStatus[stage.Name] = TaskQueued
		for _, taskName := range stage.Tasks {
			state.TaskStatus[taskName] = TaskQueued
		}
	}
	return state
}

func (s *ExecutionState) Fail(stage, task string) {
	if task != "" {
		s.TaskStatus[task] = TaskFailed
	}
	s.StageStatus[stage] = TaskFailed
	s.WorkflowStatus = TaskFailed
}
