package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	// Screen: Full terminal background
	Screen lipgloss.Style
	// Primary: Headers, titles
	Primary lipgloss.Style
	// Secondary: Supporting text
	Secondary lipgloss.Style
	// Success: Passed tasks, checkmarks
	Success lipgloss.Style
	// Error: Failed tasks, errors
	Error lipgloss.Style
	// Warning: Running tasks, in progress
	Warning lipgloss.Style
	// Info: Queued tasks, waiting
	Info lipgloss.Style
	// Highlight: Selected items, focus
	Highlight lipgloss.Style
	// Subtle: Less important text
	Subtle lipgloss.Style
	// Border: Dividers, borders
	Border lipgloss.Style
}

// Dracula color palette
const (
	DraculaBackground  = "#282a36"
	DraculaCurrentLine = "#44475a"
	DraculaForeground  = "#f8f8f2"
	DraculaComment     = "#6272a4"
	DraculaRed         = "#ff5555"
	DraculaOrange      = "#ffb86c"
	DraculaYellow      = "#f1fa8c"
	DraculaGreen       = "#50fa7b"
	DraculaPurple      = "#bd93f9"
	DraculaCyan        = "#8be9fd"
	DraculaPink        = "#ff79c6"
)

func Dracula() Theme {
	return Theme{
		Screen: lipgloss.NewStyle().
			Background(lipgloss.Color(DraculaBackground)).
			Foreground(lipgloss.Color(DraculaForeground)).
			Padding(1, 2),

		Primary: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(DraculaCyan)),

		Secondary: lipgloss.NewStyle().
			Foreground(lipgloss.Color(DraculaForeground)),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color(DraculaGreen)).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(DraculaRed)).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(DraculaOrange)).
			Bold(true),

		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color(DraculaPurple)),

		Highlight: lipgloss.NewStyle().
			Background(lipgloss.Color(DraculaPurple)).
			Foreground(lipgloss.Color(DraculaBackground)).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1),

		Subtle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(DraculaComment)),

		Border: lipgloss.NewStyle().
			Foreground(lipgloss.Color(DraculaCurrentLine)),
	}
}

// RenderScreen paints content onto a full-terminal Dracula background.
func (t Theme) RenderScreen(width, height int, content string) string {
	body := t.Screen.Render(content)
	if width <= 0 || height <= 0 {
		return body
	}
	return lipgloss.Place(
		width,
		height,
		lipgloss.Top,
		lipgloss.Left,
		body,
		lipgloss.WithWhitespaceBackground(lipgloss.Color(DraculaBackground)),
	)
}

// Default returns Dracula theme
func Default() Theme {
	return Dracula()
}