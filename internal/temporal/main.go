package temporal

import (
	"context"
	"log"

	"conductor-ci/internal/types"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	taskQueue            = "conductor-ci-queue"
	queryExecutionState  = "execution-state"
	approveSignal        = "approve"
	retrySignal          = "retry"
)

type silentLogger struct{}

func (silentLogger) Debug(string, ...any) {}
func (silentLogger) Info(string, ...any)  {}
func (silentLogger) Warn(string, ...any)  {}
func (silentLogger) Error(string, ...any) {}

var (
	temporalClient client.Client
	workflowRun    client.WorkflowRun
)

func CreateTemporalClient() (client.Client, error) {
	return client.Dial(client.Options{Logger: silentLogger{}})
}

func ApproveActivity() {
	if temporalClient == nil || workflowRun == nil {
		return
	}
	_ = temporalClient.SignalWorkflow(context.Background(), workflowRun.GetID(), workflowRun.GetRunID(), approveSignal, true)
}

func RetryActivity() {
	if temporalClient == nil || workflowRun == nil {
		return
	}
	_ = temporalClient.SignalWorkflow(context.Background(), workflowRun.GetID(), workflowRun.GetRunID(), retrySignal, true)
}

func StartTemporalServer(cfg types.WorkflowConfig) types.ExecutionState {
	c, err := CreateTemporalClient()
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(prWorkflow)
	w.RegisterActivity(runTaskActivity)
	if err := w.Start(); err != nil {
		log.Fatalln(err)
	}
	defer w.Stop()

	return startWorkflow(c, cfg)
}

func startWorkflow(c client.Client, cfg types.WorkflowConfig) types.ExecutionState {
	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        cfg.Name,
		TaskQueue: taskQueue,
	}, prWorkflow, cfg)
	if err != nil {
		state := types.QueuedState(cfg)
		state.WorkflowStatus = types.TaskFailed
		return state
	}

	temporalClient = c
	workflowRun = we

	var state types.ExecutionState
	if err := we.Get(context.Background(), &state); err != nil {
		state = types.QueuedState(cfg)
		state.WorkflowStatus = types.TaskFailed
		return state
	}
	return state
}

func PollExecutionState() types.ExecutionState {
	var state types.ExecutionState
	if temporalClient == nil || workflowRun == nil {
		return state
	}
	val, err := temporalClient.QueryWorkflow(context.Background(), workflowRun.GetID(), workflowRun.GetRunID(), queryExecutionState)
	if err != nil {
		return state
	}
	_ = val.Get(&state)
	return state
}
