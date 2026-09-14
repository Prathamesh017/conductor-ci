package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	commands []string
	cursor int
	selectedCommand string
}

func  InitialModel() model {
	return model{
		commands: []string{"Lint","Format","Test","Build","Run"},
		cursor: 0,
		selectedCommand: "",
	}
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    //checking if the message sent is a key press
    keyMsg, ok := msg.(tea.KeyMsg)
    if !ok {
        return m, nil
    }

    switch keyMsg.String() {
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
            //select unselect flow
            if m.selectedCommand == "" {
                m.selectedCommand = m.commands[m.cursor]
            } else {
                m.selectedCommand = ""
            }
            
    }
    return m, nil
}

func (m model) View() string {
    // The header
    s := "What should we buy at the market?\n\n"

    // Iterate over our choices
    for i, choice := range m.commands {

        // Is the cursor pointing at this choice?
        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        // Is this choice selected?
        checked := " " // not selected
        if m.selectedCommand == choice {
            checked = "✓" // selected!
        }

        // Render the row
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

    // The footer
    s += "\nPress q to quit.\n"

    // Send the UI for rendering
    return s
}