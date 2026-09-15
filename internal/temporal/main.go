package temporal

import (
	"context"
	"log"

	"conductor-ci/internal/types"
	"go.temporal.io/sdk/worker"

	"go.temporal.io/sdk/client"
)

const taskQueue = "conductor-ci-queue"

type silentLogger struct{}

func (silentLogger) Debug(string, ...any) {}
func (silentLogger) Info(string, ...any)  {}
func (silentLogger) Warn(string, ...any)  {}
func (silentLogger) Error(string, ...any) {}

func CreateTemporalClient() (client.Client, error) {
	c, err := client.Dial(client.Options{Logger: silentLogger{}})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	return c, nil
}

func StartTemporalServer(cfg types.WorkflowConfig) {
	c, err := CreateTemporalClient()
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	w := worker.New(c, taskQueue, worker.Options{})
	w.RegisterWorkflow(prWorkflow)
	if err := w.Start(); err != nil {
		log.Fatalln(err)
	}
	defer w.Stop()

	startWorkflow(c, cfg.Name, taskQueue, 5)
}

func startWorkflow(c client.Client, workflowName, queueName string, cfg int) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowName,
		TaskQueue: queueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, prWorkflow, cfg)
	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}

	err = we.Get(context.Background(), nil)
	if err != nil {
		log.Fatalln("Unable to get workflow result", err)
	}
}
