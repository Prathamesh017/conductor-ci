package temporal

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

func CreateTemporalClient() (client.Client, error) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	return c, nil
}

func StartTemporalServer() {
	c, err := CreateTemporalClient()
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	createWorkflow(c, "pr-validation", "conductor-ci-queue")
}

func createWorkflow(c client.Client, workflowName string, queueName string) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowName,
		TaskQueue: queueName,
	}

	c.ExecuteWorkflow(context.Background(), workflowOptions, prValidationWorkflow, 1)
}
