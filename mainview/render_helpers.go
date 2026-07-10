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

	return styledLine(prefix, text, selected, focused, false, width)
}

func songLine(text string, index int, m Model, width int) string {
	selected := index == m.selectedSong
	playing := m.currentSong != nil && index < len(m.songs) && m.songs[index].ID == m.currentSong.ID
	prefix := "  "
	if selected {
		prefix = "> "
	} else if playing {
		prefix = "* "
	}

	return styledLine(prefix, text, selected, m.focused == songsPane, playing, width)
}

func styledLine(prefix, text string, selected, focused, playing bool, width int) string {
	line := prefix + ui.Truncate(text, max(width-len(prefix), 1))
	if selected && focused {
		return ui.SelectedRowFocused.Width(width).Render(line)
	}
	if playing {
		return ui.PlayingRow.Width(width).Render(line)
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
