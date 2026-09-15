package temporal

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"conductor-ci/internal/types"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func prWorkflow(ctx workflow.Context, workflowConfig types.WorkflowConfig) (string, error) {
	var results []string
	for _, stage := range workflowConfig.Execution {
		for _, taskName := range stage.Tasks {
			taskDef, ok := workflowConfig.Tasks[taskName]
			if !ok {
				return "", fmt.Errorf("unknown task %q", taskName)
			}

			activityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 5 * time.Minute,
				RetryPolicy: &temporal.RetryPolicy{
					MaximumAttempts: retryAttempts(taskDef.Retries),
				},
			})

			var result string
			err := workflow.ExecuteActivity(activityCtx, runTaskActivity, taskDef.Script, workflowConfig.Root).Get(activityCtx, &result)
			if err != nil {
				return "", err
			}
			results = append(results, result)
		}
	}
	return strings.Join(results, "\n"), nil
}

func retryAttempts(retries int) int32 {
	if retries < 0 {
		retries = 0
	}
	return int32(retries) + 1
}

func runTaskActivity(ctx context.Context, script string, workDir string) (string, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", script)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(output))
	if err != nil {
		if out != "" {
			return "", fmt.Errorf("%w: %s", err, out)
		}
		return "", err
	}
	return out, nil
}
