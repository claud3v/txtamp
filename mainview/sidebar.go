package mainview

import (
	"strings"
	"txtamp/ui"
)

func (m Model) renderPlaylists(width, height int) string {
	lines := []string{
		ui.Title.Render("TxtAmp"),
		ui.Subtitle.Render("Connected: " + m.connectedTo),
		"",
		paneTitle("Playlists", m.focused == playlistsPane),
	}

	if m.loading && len(m.playlists) == 0 {
		lines = append(lines, ui.Subtitle.Render("Loading..."))
	}

	for i, playlist := range m.playlists {
		line := selectableLine(playlist.Name, i == m.selectedPlaylist, m.focused == playlistsPane, width-4)
		lines = append(lines, line)
	}

	return ui.Sidebar.
		Width(width).
		Height(height).
		Render(strings.Join(lines, "\n"))
}
