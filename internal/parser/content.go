package parser

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func runContentChecks(data []byte) (Group, Workflow, bool) {
	group := Group{Name: "Content Checks"}

	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		group.Checks = []Check{
			failCheck("At least 1 task defined", "could not parse tasks/execution: "+err.Error()),
		}
		return group, wf, false
	}

	group.Checks = append(group.Checks, checkAtLeastOne("At least 1 task defined", len(wf.Tasks), "tasks list is empty"))
	group.Checks = append(group.Checks, checkAtLeastOne("At least 1 stage defined", len(wf.Execution), "execution list is empty"))
	group.Checks = append(group.Checks, checkUniqueNames("No duplicate task names", taskNames(wf), "task"))
	group.Checks = append(group.Checks, checkUniqueNames("No duplicate stage names", stageNames(wf), "stage"))

	return group, wf, groupPassed(group)
}

func checkAtLeastOne(name string, count int, emptyDetail string) Check {
	if count < 1 {
		return failCheck(name, emptyDetail)
	}
	return passCheck(name)
}

func checkUniqueNames(checkName string, names []string, kind string) Check {
	var problems []string
	seen := map[string]int{}
	empty := 0
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			empty++
			continue
		}
		seen[name]++
	}
	if empty > 0 {
		problems = append(problems, fmt.Sprintf("one or more %ss are missing a name", kind))
	}

	var dups []string
	for name, count := range seen {
		if count > 1 {
			dups = append(dups, name)
		}
	}
	sort.Strings(dups)
	if len(dups) > 0 {
		problems = append(problems, fmt.Sprintf("duplicate %s names: %s", kind, strings.Join(dups, ", ")))
	}

	if len(problems) > 0 {
		return failCheck(checkName, strings.Join(problems, "\n"))
	}
	return passCheck(checkName)
}

func taskNames(wf Workflow) []string {
	names := make([]string, 0, len(wf.Tasks))
	for _, task := range wf.Tasks {
		names = append(names, task.Name)
	}
	return names
}

func stageNames(wf Workflow) []string {
	names := make([]string, 0, len(wf.Execution))
	for _, stage := range wf.Execution {
		names = append(names, stage.Name)
	}
	return names
}
