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

	start, end := visibleRange(m.selectedSong, len(m.songs), visibleSongRows(height))
	for i := start; i < end; i++ {
		song := m.songs[i]
		titleWidth := max(width-18, 10)
		title := ui.Truncate(song.Title, titleWidth)
		line := fmt.Sprintf("%-*"+"s %5s", titleWidth, title, formatDuration(song.Duration))
		line = songLine(line, i, m, width-4)
		lines = append(lines, line)
	}

	return ui.MainPane.
		Width(width).
		Height(height).
		MaxHeight(height).
		Render(strings.Join(lines, "\n"))
}

func visibleSongRows(height int) int {
	return max(height-5, 1)
}
