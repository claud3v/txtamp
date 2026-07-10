package mainview

import (
	"strings"
	"txtamp/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderSidebar(width, height int) string {
	lines := []string{
		ui.Title.Render("TxtAmp"),
		ui.Subtitle.Render("Connected: " + m.connectedTo),
		"",
		modeSelector(width-4, m.mode, m.focused == modeSelectorPane),
		"",
		paneTitle(m.sidebarTitle(), m.focused == playlistsPane),
	}

	if m.loading && m.sidebarItemCount() == 0 {
		lines = append(lines, ui.Subtitle.Render("Loading..."))
	}

	lines = append(lines, m.sidebarItems(width-4, height)...)

	return ui.Sidebar.
		Width(width).
		Height(height).
		MaxHeight(height).
		Render(strings.Join(lines, "\n"))
}

func modeSelector(width int, mode sidebarMode, focused bool) string {
	label := modeLabel(mode) + " v"
	if focused {
		label = ui.ModeSelectorActive.Render(label)
	} else {
		label = ui.ModeSelector.Render(label)
	}

	return lipgloss.PlaceHorizontal(width, lipgloss.Center, label)
}

func (m Model) sidebarTitle() string {
	return modeLabel(m.mode)
}

func modeLabel(mode sidebarMode) string {
	if mode == bandsMode {
		return "Bands"
	}

	return "Playlists"
}

func (m Model) sidebarItemCount() int {
	if m.mode == bandsMode {
		return len(m.artists)
	}

	return len(m.playlists)
}

func (m Model) sidebarItems(width, height int) []string {
	switch m.mode {
	case bandsMode:
		if len(m.artists) == 0 && !m.loading {
			return []string{ui.Subtitle.Render("No bands found")}
		}

		start, end := visibleRange(m.selectedArtist, len(m.artists), visibleSidebarRows(height))
		lines := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			artist := m.artists[i]
			line := selectableLine(artist.Name, i == m.selectedArtist, m.focused == playlistsPane, width)
			lines = append(lines, line)
		}

		return lines
	default:
		if len(m.playlists) == 0 && !m.loading {
			return []string{ui.Subtitle.Render("No playlists found")}
		}

		start, end := visibleRange(m.selectedPlaylist, len(m.playlists), visibleSidebarRows(height))
		lines := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			playlist := m.playlists[i]
			line := selectableLine(playlist.Name, i == m.selectedPlaylist, m.focused == playlistsPane, width)
			lines = append(lines, line)
		}

		return lines
	}
}

func visibleSidebarRows(height int) int {
	return max(height-8, 1)
}
