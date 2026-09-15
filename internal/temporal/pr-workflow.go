package temporal

import (
	"conductor-ci/internal/types"
	"go.temporal.io/sdk/workflow"
)

func prWorkflow(ctx workflow.Context, cfg types.WorkflowConfig) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("workflow started", "name", cfg.Name, "stages", len(cfg.Execution), "tasks", len(cfg.Tasks))
	return nil
}
