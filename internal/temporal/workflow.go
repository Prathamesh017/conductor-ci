package temporal

import (
	"context"
	"fmt"
)

func  prValidationWorkflow(ctx context.Context, prNumber int) {
	fmt.Println("PR Validation Workflow started")
	fmt.Println("PR Number:", prNumber)
}