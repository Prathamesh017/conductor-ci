package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"conductor-ci/internal/types"

	"gopkg.in/yaml.v3"
)

// ParseWorkflow reads filePath into a WorkflowConfig.
// It does not validate; call Validate first if you need checks.
func ParseWorkflow(filePath string) (*types.WorkflowConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read workflow file: %w", err)
	}
	return parseWorkflowBytes(data, filepath.Dir(filePath))
}

func parseWorkflowBytes(data []byte, root string) (*types.WorkflowConfig, error) {
	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("parse workflow yaml: %w", err)
	}
	return wf.toConfig(root), nil
}

func (w Workflow) toConfig(root string) *types.WorkflowConfig {
	tasks := make(map[string]types.Task, len(w.Tasks))
	for _, task := range w.Tasks {
		name := task.Name
		tasks[name] = types.Task{
			Name:    name,
			Script:  task.Script,
			Timeout: task.Timeout,
			Retries: task.Retries,
		}
	}

	stages := make([]types.Stage, 0, len(w.Execution))
	for _, stage := range w.Execution {
		stages = append(stages, types.Stage{
			Name:             stage.Name,
			Description:      stage.Description,
			Tasks:            stage.Tasks,
			Mode:             stage.Mode,
			RequiresApproval: stage.RequiresApproval,
		})
	}

	return &types.WorkflowConfig{
		Name:      w.Name,
		Root:      root,
		Tasks:     tasks,
		Execution: stages,
	}
}
