package types

import "time"

// WorkflowConfig is the parsed workflow.yml. Treat it as read-only after parse.
type WorkflowConfig struct {
	Name      string
	Root      string
	Tasks     map[string]Task
	Execution []Stage
}

type Task struct {
	Name    string
	Script  string
	Timeout string `yaml:"timeout,omitempty"`
	Retries int    `yaml:"retries,omitempty"`
}

type Stage struct {
	Name             string
	Description      string
	Tasks            []string
	Mode             string
	RequiresApproval bool
}

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
