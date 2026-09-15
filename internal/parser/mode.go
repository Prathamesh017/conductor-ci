package parser

import (
	"fmt"
	"strings"
)

var validModes = map[string]struct{}{
	"parallel":   {},
	"sequential": {},
}

func runModeChecks(wf Workflow) (Group, bool) {
	group := Group{Name: "Mode Checks"}
	group.Checks = append(group.Checks, checkStageModes(wf))
	return group, groupPassed(group)
}

func checkStageModes(wf Workflow) Check {
	const name = "All stages have valid mode (parallel/sequential)"

	var problems []string
	for _, stage := range wf.Execution {
		stageName := strings.TrimSpace(stage.Name)
		if stageName == "" {
			stageName = "(unnamed)"
		}
		mode := strings.TrimSpace(stage.Mode)
		if mode == "" {
			problems = append(problems, fmt.Sprintf("stage %q is missing a mode", stageName))
			continue
		}
		if _, ok := validModes[mode]; !ok {
			problems = append(problems, fmt.Sprintf("stage %q has mode %q; want parallel or sequential", stageName, mode))
		}
	}

	if len(problems) > 0 {
		return failCheck(name, strings.Join(problems, "\n"))
	}
	return passCheck(name)
}
