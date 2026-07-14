package ui

import "github.com/charmbracelet/lipgloss"

type ShellLayout struct {
	Width        int
	Height       int
	BodyHeight   int
	SidebarWidth int
	MainWidth    int
}

func Center(width, height int, content string) string {
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func NewShellLayout(width, height int) ShellLayout {
	contentWidth := max(width, 80)
	contentHeight := max(height, 24)
	sidebarWidth := min(max(contentWidth/4, 24), 34)

	return ShellLayout{
		Width:        contentWidth,
		Height:       contentHeight,
		BodyHeight:   max(contentHeight-7, 8),
		SidebarWidth: sidebarWidth,
		MainWidth:    max(contentWidth-sidebarWidth-1, 30),
	}
}

func Truncate(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if len(value) <= width {
		return value
	}
	if width == 1 {
		return value[:1]
	}

	return value[:width-1] + "."
}
