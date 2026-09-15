package temporal

import (
	"go.temporal.io/sdk/workflow"
)

func prWorkflow(ctx workflow.Context, workflowConfig int) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("workflow started", "name", workflowConfig)
	// logger.Info("workflow started", "name", workflowConfig.Name, "stages", len(workflowConfig.Execution), "tasks", len(workflowConfig.Tasks))
	// fmt.Println("workflow started", workflowConfig)
	return nil
}
