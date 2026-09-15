package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"conductor-ci/internal/parser"
	"conductor-ci/internal/theme"
	"conductor-ci/internal/types"

	"github.com/charmbracelet/lipgloss"
)

func renderReport(t theme.Theme, report parser.Report) string {
	var b strings.Builder

	target := "workflow.yml"
	if report.Path != "" {
		target = filepath.Base(report.Path)
	}
	b.WriteString(t.Primary.Render("Validating "+target+"...") + "\n")

	passed, failed, skipped := 0, 0, 0
	for _, group := range report.Groups {
		b.WriteString("\n")
		b.WriteString(t.Info.Render(group.Name) + "\n")
		for _, check := range group.Checks {
			switch check.Status {
			case parser.StatusPass:
				passed++
			case parser.StatusFail:
				failed++
			default:
				skipped++
			}
			b.WriteString(renderCheck(t, check) + "\n")
			if check.Detail != "" && check.Status != parser.StatusPass {
				for _, line := range strings.Split(check.Detail, "\n") {
					b.WriteString("         " + t.Subtle.Render(line) + "\n")
				}
			}
		}
	}

	b.WriteString("\n")
	if report.Valid {
		b.WriteString(t.Success.Render("✓ Workflow is valid") + "\n")
	} else {
		b.WriteString(t.Error.Render("✗ Workflow is invalid") + "\n")
	}
	b.WriteString(t.Subtle.Render(fmt.Sprintf("%d passed · %d failed · %d skipped", passed, failed, skipped)) + "\n")

	b.WriteString("\n" + t.Subtle.Render("Press enter to go back · q to quit."))
	return b.String()
}

func renderWorkflow(t theme.Theme, cfg types.WorkflowConfig, state types.ExecutionState) string {
	s := t.Primary.Render(cfg.Name) + "\n\n"

	for _, stage := range cfg.Execution {
		icon, _ := statusLook(t, state.StageStatus[stage.Name])
		stageLine := fmt.Sprintf("%s %s", icon, stage.Name)
		if state.StageStatus[stage.Name] == types.TaskAwaiting {
			stageLine += " (awaiting approval...)"
		}
		s += t.Secondary.Render(stageLine) + "\n"

		for _, taskName := range stage.Tasks {
			task := cfg.Tasks[taskName]
			icon, style := statusLook(t, state.TaskStatus[taskName])
			duration := state.TaskDuration[taskName]

			line := fmt.Sprintf("  %s %s", icon, task.Name)
			if state.TaskStatus[taskName] == types.TaskRunning {
				line += " (running...)"
			} else if duration > 0 {
				line += fmt.Sprintf(" (%s)", duration)
			}
			s += style.Render(line) + "\n"
		}
		s += "\n"
	}

	hint := "Press enter to go back · q to quit."
	for _, stage := range cfg.Execution {
		if state.StageStatus[stage.Name] == types.TaskAwaiting {
			hint = "Press a to approve · enter to go back · q to quit."
			break
		}
	}
	s += t.Subtle.Render(hint)
	return s
}

func statusLook(t theme.Theme, status types.TaskStatus) (string, lipgloss.Style) {
	switch status {
	case types.TaskPassed:
		return "✓", t.Success
	case types.TaskFailed:
		return "✗", t.Error
	case types.TaskRunning:
		return "⟳", t.Warning
	case types.TaskAwaiting:
		return "⏸", t.Warning
	default:
		return "⏳", t.Info
	}
}

func renderCheck(t theme.Theme, check parser.Check) string {
	switch check.Status {
	case parser.StatusPass:
		return "  " + t.Success.Render("PASS") + "  " + t.Success.Render(check.Name)
	case parser.StatusFail:
		return "  " + t.Error.Render("FAIL") + "  " + t.Error.Render(check.Name)
	default:
		return "  " + t.Subtle.Render("SKIP") + "  " + t.Subtle.Render(check.Name)
	}
}
