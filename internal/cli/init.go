package cli

import (
	"fmt"

	"conductor-ci/internal/parser"
	"conductor-ci/internal/temporal"
	"conductor-ci/internal/theme"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	validateWorkflowCommand    = "validate-workflow"
	startTemporalServerCommand = "start-temporal-server"
)

type screen int

const (
	screenMenu screen = iota
	screenReport
	screenWorkflow
)

type model struct {
	commands        []string
	cursor          int
	selectedCommand string
	theme           theme.Theme
	width           int
	height          int
	screen          screen
	report          parser.Report
}

func InitialModel() model {
	return model{
		commands: []string{validateWorkflowCommand, startTemporalServerCommand},
		theme:    theme.Default(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) onMenu() bool {
	return m.screen == screenMenu
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if !m.onMenu() {
				return m, nil
			}
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if !m.onMenu() {
				return m, nil
			}
			if m.cursor < len(m.commands)-1 {
				m.cursor++
			}
		case "enter":
			if !m.onMenu() {
				m.screen = screenMenu
				m.selectedCommand = ""
				return m, nil
			}
			m.selectedCommand = m.commands[m.cursor]
			switch m.selectedCommand {
			case validateWorkflowCommand:
				m.report = parser.Validate(".")
				m.screen = screenReport
			case startTemporalServerCommand:
				m.report = parser.Validate(".")
				if !m.report.Valid {
					m.screen = screenReport
					break
				}
				m.screen = screenWorkflow
				cfg := *m.report.Config
				return m, func() tea.Msg {
					temporal.StartTemporalServer(cfg)
					return nil
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var body string
	switch m.screen {
	case screenReport:
		body = renderReport(m.theme, m.report)
	case screenWorkflow:
		body = renderWorkflow(m.theme)
	default:
		body = renderMenu(m)
	}
	return m.theme.RenderScreen(m.width, m.height, body)
}

func renderMenu(m model) string {
	s := m.theme.Primary.Render("Welcome to the Conductor CI!") + "\n\n"

	for i, choice := range m.commands {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		checkStyle := m.theme.Secondary
		if m.selectedCommand == choice {
			checked = "✓"
			checkStyle = m.theme.Success
		}

		line := fmt.Sprintf("%s [%s] %s", cursor, checked, choice)
		if m.cursor == i {
			s += m.theme.Highlight.Render(line) + "\n"
		} else {
			s += checkStyle.Render(line) + "\n"
		}
	}

	s += "\n" + m.theme.Subtle.Render("Press q to quit.")
	return s
}
