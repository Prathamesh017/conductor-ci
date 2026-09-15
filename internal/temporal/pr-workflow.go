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

func prWorkflow(ctx workflow.Context, cfg types.WorkflowConfig) (types.ExecutionState, error) {
	state := types.QueuedState(cfg)

	for _, stage := range cfg.Execution {
		var names []string
		var futures []workflow.Future

		for _, taskName := range stage.Tasks {
			task, ok := cfg.Tasks[taskName]
			if !ok {
				state.Fail(stage.Name, taskName)
				return state, nil
			}

			future := startTask(ctx, cfg, task)
			if stage.Mode == "parallel" {
				names = append(names, taskName)
				futures = append(futures, future)
				continue
			}

			if err := future.Get(ctx, nil); err != nil {
				state.Fail(stage.Name, taskName)
				return state, nil
			}
			state.TaskStatus[taskName] = types.TaskPassed
		}

		failed := false
		for i, future := range futures {
			if err := future.Get(ctx, nil); err != nil {
				state.TaskStatus[names[i]] = types.TaskFailed
				failed = true
				continue
			}
			state.TaskStatus[names[i]] = types.TaskPassed
		}
		if failed {
			state.Fail(stage.Name, "")
			return state, nil
		}
		state.StageStatus[stage.Name] = types.TaskPassed
	}

	state.WorkflowStatus = types.TaskPassed
	return state, nil
}

func startTask(ctx workflow.Context, cfg types.WorkflowConfig, task types.Task) workflow.Future {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: int32(max(task.Retries, 0) + 1),
		},
	})
	return workflow.ExecuteActivity(ctx, runTaskActivity, task.Script, cfg.Root)
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
