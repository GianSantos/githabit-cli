package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Bold warning for nudge when today's score is 0
	NudgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ff6b6b")).
			Padding(0, 1)

	// Streak grid styles
	EmptyCell  = lipgloss.NewStyle().Foreground(lipgloss.Color("#4a4a4a")).SetString("·")
	Level0     = lipgloss.NewStyle().Foreground(lipgloss.Color("#161b22"))
	Level1     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0e4429"))
	Level2     = lipgloss.NewStyle().Foreground(lipgloss.Color("#006d32"))
	Level3     = lipgloss.NewStyle().Foreground(lipgloss.Color("#26a641"))
	Level4     = lipgloss.NewStyle().Foreground(lipgloss.Color("#39d353"))
	TitleStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1)
)

// StreakCell returns a styled character for the contribution level (0-4).
func StreakCell(level int) string {
	switch level {
	case 0:
		return Level0.SetString("·").String()
	case 1:
		return Level1.SetString("▁").String()
	case 2:
		return Level2.SetString("▂").String()
	case 3:
		return Level3.SetString("▃").String()
	default:
		return Level4.SetString("▄").String()
	}
}
