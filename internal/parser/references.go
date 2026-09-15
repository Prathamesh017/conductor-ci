package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func runReferenceChecks(dir string, wf Workflow) (Group, bool) {
	group := Group{Name: "References Checks"}

	group.Checks = append(group.Checks, checkStageTasksInCatalog(wf))
	group.Checks = append(group.Checks, checkScripts(dir, wf))

	return group, groupPassed(group)
}

// check if all the tasks in the stage are present in the Tasks section or not.
func checkStageTasksInCatalog(wf Workflow) Check {
	const name = "All stage tasks exist in catalog"
	var catalog []string
	for _, task := range wf.Tasks {
		taskName := strings.TrimSpace(task.Name)
		if taskName == "" {
			continue
		}
		catalog = append(catalog, taskName)
	}

	var missing []string
	for _, stage := range wf.Execution {
		stageName := strings.TrimSpace(stage.Name)
		for _, taskName := range stage.Tasks {
			taskName = strings.TrimSpace(taskName)
			if taskName == "" {
				missing = append(missing, fmt.Sprintf("stage %q references an empty task name", stageName))
				continue
			}
			if !slices.Contains(catalog, taskName) {
				missing = append(missing, fmt.Sprintf("task %q in stage %q is not in the tasks catalog", taskName, stageName))
			}
		}
	}

	if len(missing) > 0 {
		return failCheck(name, strings.Join(missing, "\n"))
	}
	return passCheck(name)
}

func checkScripts(dir string, wf Workflow) Check {
	const name = "All scripts exist and are executable"

	var problems []string
	for _, task := range wf.Tasks {
		taskName := strings.TrimSpace(task.Name)
		if taskName == "" {
			taskName = "(unnamed)"
		}
		if err := validateTaskScript(dir, task.Script); err != nil {
			problems = append(problems, fmt.Sprintf("task %q: %s", taskName, err.Error()))
		}
	}

	if len(problems) > 0 {
		return failCheck(name, strings.Join(problems, "\n"))
	}
	return passCheck(name)
}

func validateTaskScript(dir, script string) error {
	script = strings.TrimSpace(script)
	if script == "" {
		return fmt.Errorf("script is empty")
	}

	path, isFile := checkScriptFileExists(dir, script)
	if !isFile {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", displayScript(script))
		}
		return fmt.Errorf("could not stat %s: %w", displayScript(script), err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", displayScript(script))
	}
	if info.Mode()&0o111 == 0 {
		return fmt.Errorf("%s is not executable (chmod +x)", displayScript(script))
	}
	return nil
}

// Check if  the script mentions a file , check if the file exists or not.
// if it mentions a command like npm run dev , it will be skipped only direct file checks will be done.
func checkScriptFileExists(dir, script string) (string, bool) {
	first := strings.Fields(script)[0]
	path := first
	if !filepath.IsAbs(first) {
		path = filepath.Join(dir, first)
	}

	if strings.HasPrefix(first, ".") || filepath.IsAbs(first) || strings.Contains(first, "/") {
		return path, true
	}

	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return path, true
	}
	return "", false
}

func displayScript(script string) string {
	return strings.Fields(script)[0]
}
