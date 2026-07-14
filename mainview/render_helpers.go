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
	return styledLine(rowPrefix(selected, false, ""), text, selected, focused, false, width)
}

func searchLine(query string) string {
	cursor := ""
	if query == "" {
		cursor = "_"
	}

	return ui.Subtitle.Render("Filter: " + query + cursor)
}

func globalSearchLine(query string) string {
	cursor := ""
	if query == "" {
		cursor = "_"
	}

	return ui.Subtitle.Render("Search: " + query + cursor)
}

func songLine(text string, index int, m Model, width int) string {
	selected := index == m.selectedSong
	playing := m.currentSong != nil && index < len(m.songs) && m.songs[index].ID == m.currentSong.ID

	return styledLine(rowPrefix(selected, playing, ""), text, selected, m.focused == songsPane, playing, width)
}

func nestedSongLine(text string, index int, selected bool, m Model, width int) string {
	playing := m.currentSong != nil && index < len(m.songs) && m.songs[index].ID == m.currentSong.ID

	return styledLine(rowPrefix(selected, playing, "  "), text, selected, m.focused == songsPane, playing, width)
}

func queueSongLine(text string, index int, m Model, width int) string {
	selected := index == m.selectedQueue
	playing := m.currentSong != nil && index < len(m.queue) && m.queue[index].ID == m.currentSong.ID

	return styledLine(rowPrefix(selected, playing, ""), text, selected, m.focused == songsPane, playing, width)
}

func albumHeaderLine(prefix, text string, expanded, selected, focused bool, width int) string {
	line := prefix + ui.Truncate(text, max(width-len(prefix), 1))
	if selected && focused {
		return ui.SelectedRowFocused.Width(width).Render(line)
	}
	if selected {
		return ui.SelectedRow.Width(width).Render(line)
	}
	if expanded {
		return ui.AlbumExpanded.Width(width).Render(line)
	}

	return ui.AlbumCollapsed.Width(width).Render(line)
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

func rowPrefix(selected, playing bool, indent string) string {
	switch {
	case selected:
		return indent + "> "
	case playing:
		return indent + "* "
	default:
		return indent + "  "
	}
}

func emptyState(text string) string {
	return ui.EmptyState.Render(text)
}

func sectionHeader(text string) string {
	return ui.SectionHeader.Render(text)
}

func formatDuration(seconds int) string {
	if seconds <= 0 {
		return "--:--"
	}

	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}
