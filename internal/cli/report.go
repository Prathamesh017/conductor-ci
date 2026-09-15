package cli

import (
	"fmt"
	"path/filepath"
	"strings"
	"conductor-ci/internal/parser"
	"conductor-ci/internal/theme"
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
