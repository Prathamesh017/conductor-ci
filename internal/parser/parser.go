package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"conductor-ci/internal/types"
	"gopkg.in/yaml.v3"
)

const (
	StatusPass Status = iota
	StatusFail
	StatusSkip
)

type Status int

type Check struct {
	Name   string
	Status Status
	Detail string
}

type Group struct {
	Name   string
	Checks []Check
}

type Report struct {
	Path   string
	Groups []Group
	Valid  bool
	Config *types.WorkflowConfig
}

var ValidWorkflowFiles = []string{"workflow.yml", "workflow.yaml"}

// Validate if workflow.yaml is in correct format.
func Validate(dir string) Report {
	report := Report{}

	//First Check if workflow.yaml is present or not and if it is in correct format or not.
	fileGroup, data, path, ok := runFileChecks(dir)
	report.Path = path
	report.Groups = append(report.Groups, fileGroup)
	if !ok {
		return report
	}

	//Checking if workflow file has name , task and execution section or not.
	structureGroup, ok := runStructureChecks(data)
	report.Groups = append(report.Groups, structureGroup)
	if !ok {
		return report
	}

	//check if there are duplicates and atleast one task and one execution section is present.
	contentGroup, wf, ok := runContentChecks(data)
	report.Groups = append(report.Groups, contentGroup)
	if !ok {
		return report
	}

	refGroup, ok := runReferenceChecks(dir, wf)
	report.Groups = append(report.Groups, refGroup)
	if !ok {
		return report
	}

	modeGroup, ok := runModeChecks(wf)
	report.Groups = append(report.Groups, modeGroup)
	if !ok {
		return report
	}

	report.Valid = true
	report.Config = wf.toConfig(filepath.Dir(path))
	return report
}

func passCheck(name string) Check {
	return Check{Name: name, Status: StatusPass}
}

func failCheck(name, detail string) Check {
	return Check{Name: name, Status: StatusFail, Detail: detail}
}

func groupPassed(g Group) bool {
	for _, c := range g.Checks {
		if c.Status == StatusFail {
			return false
		}
	}
	return true
}

func runFileChecks(dir string) (Group, []byte, string, bool) {
	group := Group{}

	path, err := FindWorkflowFile(dir)

	if err != nil {
		group.Checks = []Check{
			{Name: "File exists Check", Status: StatusFail, Detail: err.Error()},
		}
		return group, nil, "", false
	}

	group.Checks = append(group.Checks, Check{Name: "File exists", Status: StatusPass})

	data, err := os.ReadFile(path)
	if err != nil {
		group.Checks = append(group.Checks, Check{
			Name:   "YAML is valid",
			Status: StatusFail,
			Detail: fmt.Sprintf("could not read %s: %v", filepath.Base(path), err),
		})
		return group, nil, path, false
	}

	//checking if ymal file is in correct format or not
	var doc any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		group.Checks = append(group.Checks, Check{
			Name:   "YAML is valid",
			Status: StatusFail,
			Detail: err.Error(),
		})
		return group, data, path, false
	}

	group.Checks = append(group.Checks, Check{Name: "YAML is valid", Status: StatusPass})
	return group, data, path, true
}

func FindWorkflowFile(dir string) (string, error) {
	for _, name := range ValidWorkflowFiles {
		path := filepath.Join(dir, name)
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", fmt.Errorf("could not stat %s: %w", name, err)
		}
		if info.IsDir() {
			continue
		}
		return path, nil
	}

	return "", fmt.Errorf("expected workflow.yml or workflow.yaml in the current directory")
}
