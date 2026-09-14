package main 

import (
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"conductor-ci/internal/cli"
)

func main(){	
	p := tea.NewProgram(cli.InitialModel(), tea.WithAltScreen())
 	if _, err := p.Run(); err != nil {
 		fmt.Println("Error running program:", err)
 		os.Exit(1)
 	}
}