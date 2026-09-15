package parser

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func runStructureChecks(data []byte) (Group, bool) {
	group := Group{Name: "Structure Checks"}

	raw := map[string]any{}
	if err := yaml.Unmarshal(data, &raw); err != nil || raw == nil {
		raw = map[string]any{}
	}

	group.Checks = append(group.Checks, checkNamePresent(raw))
	group.Checks = append(group.Checks, checkKeyPresent(raw, "tasks", "tasks section present", "workflow is missing a tasks section"))
	group.Checks = append(group.Checks, checkKeyPresent(raw, "execution", "execution section present", "workflow is missing an execution section"))

	return group, groupPassed(group)
}

func checkNamePresent(raw map[string]any) Check {
	const name = "name field present"
	value, ok := raw["name"]
	if !ok || value == nil {
		return failCheck(name, "workflow is missing a name field")
	}
	s, isString := value.(string)
	if !isString || strings.TrimSpace(s) == "" {
		return failCheck(name, "name must be a non-empty string")
	}
	return passCheck(name)
}

func checkKeyPresent(raw map[string]any, key, checkName, missingDetail string) Check {
	if _, ok := raw[key]; !ok {
		return failCheck(checkName, missingDetail)
	}
	return passCheck(checkName)
}
