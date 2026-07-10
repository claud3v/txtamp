package mainview

import (
	"fmt"
	"strings"
	"txtamp/navidrome"
	"txtamp/ui"
)

type bandRow struct {
	albumTitle string
	song       navidrome.Song
	songIndex  int
}

func (m Model) renderMainArea(width, height int) string {
	if m.mode == bandsMode {
		return m.renderBands(width, height)
	}

	return m.renderSongs(width, height)
}

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

func (m Model) renderBands(width, height int) string {
	title := "Bands"
	if len(m.artists) > 0 {
		title = m.artists[m.selectedArtist].Name
	}

	lines := []string{
		paneTitle(title, m.focused == songsPane),
		ui.Subtitle.Render(fmt.Sprintf("%d albums, %d songs", len(m.albums), len(m.songs))),
		"",
	}

	if m.err != nil {
		lines = append(lines, ui.Error.Render(m.err.Error()))
	} else if m.loading && len(m.songs) == 0 {
		lines = append(lines, ui.Subtitle.Render("Loading..."))
	} else if len(m.artists) == 0 {
		lines = append(lines, ui.Subtitle.Render("No bands found"))
	}

	rows := m.bandRows()
	selectedRow := selectedBandRow(rows, m.selectedSong)
	start, end := visibleRange(selectedRow, len(rows), visibleSongRows(height))
	for i := start; i < end; i++ {
		row := rows[i]
		if row.songIndex < 0 {
			lines = append(lines, albumLine(row.albumTitle, width-4))
			continue
		}

		titleWidth := max(width-22, 10)
		title := ui.Truncate(row.song.Title, titleWidth)
		line := fmt.Sprintf("%-*"+"s %5s", titleWidth, title, formatDuration(row.song.Duration))
		line = nestedSongLine(line, row.songIndex, m, width-4)
		lines = append(lines, line)
	}

	return ui.MainPane.
		Width(width).
		Height(height).
		MaxHeight(height).
		Render(strings.Join(lines, "\n"))
}

func (m Model) bandRows() []bandRow {
	rows := make([]bandRow, 0, len(m.albums)+len(m.songs))
	songIndex := 0
	for _, group := range m.albums {
		rows = append(rows, bandRow{albumTitle: group.album.Name, songIndex: -1})
		for _, song := range group.songs {
			rows = append(rows, bandRow{song: song, songIndex: songIndex})
			songIndex++
		}
	}

	return rows
}

func selectedBandRow(rows []bandRow, selectedSong int) int {
	for i, row := range rows {
		if row.songIndex == selectedSong {
			return i
		}
	}

	return 0
}

func albumLine(title string, width int) string {
	return ui.PaneTitle.Width(width).Render("> " + ui.Truncate(title, max(width-2, 1)))
}

func visibleSongRows(height int) int {
	return max(height-5, 1)
}
