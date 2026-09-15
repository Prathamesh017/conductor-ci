package parser

type Workflow struct {
	Name      string  `yaml:"name"`
	Tasks     []Task  `yaml:"tasks"`
	Execution []Stage `yaml:"execution"`
}

type Task struct {
	Name   string `yaml:"name"`
	Script string `yaml:"script"`
}

type Stage struct {
	Name             string   `yaml:"stage"`
	Description      string   `yaml:"description"`
	Tasks            []string `yaml:"tasks"`
	Mode             string   `yaml:"mode"`
	RequiresApproval bool     `yaml:"requires_approval"`
}
