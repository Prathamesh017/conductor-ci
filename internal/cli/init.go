package cli

import (
	"fmt"

	"conductor-ci/internal/theme"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	commands        []string
	cursor          int
	selectedCommand string
	theme           theme.Theme
	width           int
	height          int
}

func InitialModel() model {
	return model{
		commands:        []string{"Lint", "Format", "Test", "Build", "Run"},
		cursor:          0,
		selectedCommand: "",
		theme:           theme.Default(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
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
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.commands)-1 {
				m.cursor++
			}
		case "enter":
			if m.selectedCommand == "" {
				m.selectedCommand = m.commands[m.cursor]
			} else {
				m.selectedCommand = ""
			}
		}
	}
	return m, nil
}

func (m model) View() string {
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

	return m.theme.RenderScreen(m.width, m.height, s)
}