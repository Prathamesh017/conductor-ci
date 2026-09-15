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
	if err := workflow.SetQueryHandler(ctx, queryExecutionState, func() (types.ExecutionState, error) {
		return state, nil
	}); err != nil {
		return state, err
	}

	for _, stage := range cfg.Execution {
		if stage.RequiresApproval {
			state.StageStatus[stage.Name] = types.TaskAwaiting
			workflow.GetSignalChannel(ctx, approveSignal).Receive(ctx, nil)
			state.StageStatus[stage.Name] = types.TaskQueued
		}

		var names []string
		var futures []workflow.Future

		for _, taskName := range stage.Tasks {
			task, ok := cfg.Tasks[taskName]
			if !ok {
				state.Fail(stage.Name, taskName)
				return state, nil
			}

			state.TaskStatus[taskName] = types.TaskRunning
			future := startTask(ctx, cfg, task)
			if stage.Mode == "parallel" {
				names = append(names, taskName)
				futures = append(futures, future)
				continue
			}

			if err := finishTask(ctx, &state, taskName, future); err != nil {
				state.Fail(stage.Name, taskName)
				return state, nil
			}
		}

		failed := false
		for i, future := range futures {
			if err := finishTask(ctx, &state, names[i], future); err != nil {
				failed = true
			}
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

func finishTask(ctx workflow.Context, state *types.ExecutionState, taskName string, future workflow.Future) error {
	var d time.Duration
	err := future.Get(ctx, &d)
	state.TaskDuration[taskName] = d
	if err != nil {
		state.TaskStatus[taskName] = types.TaskFailed
		return err
	}
	state.TaskStatus[taskName] = types.TaskPassed
	return nil
}

func runTaskActivity(ctx context.Context, script string, workDir string) (time.Duration, error) {
	start := time.Now()
	cmd := exec.CommandContext(ctx, "sh", "-c", script)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	d := time.Since(start)
	if err != nil {
		if out := strings.TrimSpace(string(output)); out != "" {
			return d, fmt.Errorf("%w: %s", err, out)
		}
		return d, err
	}
	return d, nil
}
