package mainview

import (
	"fmt"
	"txtamp/ui"

	"github.com/charmbracelet/lipgloss"
)

func paneTitle(title string, focused bool) string {
	if focused {
		return ui.PaneTitleFocused.Render(title)
	}

	return ui.PaneTitle.Render(title)
}

func selectableLine(text string, selected, focused bool, width int) string {
	prefix := "  "
	if selected {
		prefix = "> "
	}

	line := prefix + ui.Truncate(text, max(width-2, 1))
	if selected && focused {
		return ui.SelectedRowFocused.Width(width).Render(line)
	}
	if selected {
		return ui.SelectedRow.Width(width).Render(line)
	}

	return lipgloss.NewStyle().Width(width).Render(line)
}

func formatDuration(seconds int) string {
	if seconds <= 0 {
		return "--:--"
	}

	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}
