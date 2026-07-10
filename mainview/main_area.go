package mainview

import (
	"fmt"
	"strings"
	"txtamp/ui"
)

func (m Model) renderSongs(width, height int) string {
	title := "Songs"
	if len(m.playlists) > 0 {
		title = m.playlists[m.selectedPlaylist].Name
	}

	lines := []string{
		paneTitle(title, m.focused == songsPane),
		ui.Subtitle.Render(fmt.Sprintf("%d songs", len(m.songs))),
		"",
	}

	if m.err != nil {
		lines = append(lines, ui.Error.Render(m.err.Error()))
	} else if m.loading && len(m.songs) == 0 {
		lines = append(lines, ui.Subtitle.Render("Loading..."))
	} else if len(m.playlists) == 0 {
		lines = append(lines, ui.Subtitle.Render("No playlists found"))
	}

	for i, song := range m.songs {
		titleWidth := max(width-18, 10)
		title := ui.Truncate(song.Title, titleWidth)
		line := fmt.Sprintf("%-*"+"s %5s", titleWidth, title, formatDuration(song.Duration))
		line = selectableLine(line, i == m.selectedSong, m.focused == songsPane, width-4)
		lines = append(lines, line)
	}

	return ui.MainPane.
		Width(width).
		Height(height).
		Render(strings.Join(lines, "\n"))
}
