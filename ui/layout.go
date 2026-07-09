package ui

import "github.com/charmbracelet/lipgloss"

func Center(width, height int, content string) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
