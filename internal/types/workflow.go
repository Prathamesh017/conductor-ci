package types

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

