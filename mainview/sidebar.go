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
	}
	if m.filterQueryFor(playlistsPane) != "" || m.searching && m.searchPane == playlistsPane {
		lines = append(lines, searchLine(m.searchQuery))
	}
	lines = append(lines, paneTitle(m.sidebarTitle(), m.focused == playlistsPane))

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
	if mode == artistsMode {
		return "Artists"
	}

	return "Playlists"
}

func (m Model) sidebarItemCount() int {
	if m.mode == artistsMode {
		return len(m.artists)
	}

	return len(m.playlists)
}

func (m Model) sidebarItems(width, height int) []string {
	switch m.mode {
	case artistsMode:
		if m.err != nil {
			return nil
		}
		if len(m.artists) == 0 && !m.loading {
			return []string{ui.Subtitle.Render("No artists found")}
		}

		artists := m.filteredArtists()
		if len(artists) == 0 && m.filterQueryFor(playlistsPane) != "" {
			return []string{ui.Subtitle.Render("No matches")}
		}

		selected := m.selectedArtistPosition(artists)
		start, end := visibleRange(selected, len(artists), m.visibleSidebarRows(height))
		lines := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			artist := artists[i]
			line := selectableLine(artist.artist.Name, artist.index == m.selectedArtist, m.focused == playlistsPane, width)
			lines = append(lines, line)
		}

		return lines
	default:
		if m.err != nil {
			return nil
		}
		if len(m.playlists) == 0 && !m.loading {
			return []string{ui.Subtitle.Render("No playlists found")}
		}

		playlists := m.filteredPlaylists()
		if len(playlists) == 0 && m.filterQueryFor(playlistsPane) != "" {
			return []string{ui.Subtitle.Render("No matches")}
		}

		selected := m.selectedPlaylistPosition(playlists)
		start, end := visibleRange(selected, len(playlists), m.visibleSidebarRows(height))
		lines := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			playlist := playlists[i]
			line := selectableLine(playlist.playlist.Name, playlist.index == m.selectedPlaylist, m.focused == playlistsPane, width)
			lines = append(lines, line)
		}

		return lines
	}
}

func (m Model) visibleSidebarRows(height int) int {
	if m.filterQueryFor(playlistsPane) != "" || m.searching && m.searchPane == playlistsPane {
		return max(height-9, 1)
	}

	return max(height-8, 1)
}
