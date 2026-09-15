package temporal

import (
	"context"
	"log"

	"conductor-ci/internal/types"

	"go.temporal.io/sdk/client"
)

func CreateTemporalClient() (client.Client, error) {
	c, err := client.Dial(client.Options{})
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

	startWorkflow(c, cfg.Name, "conductor-ci-queue", cfg)
}

func startWorkflow(c client.Client, workflowName, queueName string, cfg types.WorkflowConfig) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowName,
		TaskQueue: queueName,
	}

	c.ExecuteWorkflow(context.Background(), workflowOptions, prWorkflow, cfg)
}
