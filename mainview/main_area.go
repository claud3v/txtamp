package mainview

import (
	"fmt"
	"strings"
	"txtamp/navidrome"
	"txtamp/ui"
)

type artistRow struct {
	albumTitle string
	song       navidrome.Song
	songIndex  int
}

func (m Model) renderMainArea(width, height int) string {
	if m.mode == artistsMode {
		return m.renderArtists(width, height)
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
	if m.filterQueryFor(songsPane) != "" || m.searching && m.searchPane == songsPane {
		lines = append(lines, searchLine(m.searchQuery))
	}

	if m.err != nil {
		lines = append(lines, ui.Error.Render(m.err.Error()))
	} else if m.loading && len(m.songs) == 0 {
		lines = append(lines, ui.Subtitle.Render("Loading..."))
	} else if len(m.playlists) == 0 {
		lines = append(lines, ui.Subtitle.Render("No playlists found"))
	}

	songs := m.filteredSongs()
	if len(songs) == 0 && m.filterQueryFor(songsPane) != "" {
		lines = append(lines, ui.Subtitle.Render("No matches"))
	}

	selected := m.selectedSongPosition(songs)
	start, end := visibleRange(selected, len(songs), m.visibleSongRows(height))
	for i := start; i < end; i++ {
		song := songs[i]
		titleWidth := max(width-18, 10)
		title := ui.Truncate(song.song.Title, titleWidth)
		line := fmt.Sprintf("%-*"+"s %5s", titleWidth, title, formatDuration(song.song.Duration))
		line = songLine(line, song.index, m, width-4)
		lines = append(lines, line)
	}

	return ui.MainPane.
		Width(width).
		Height(height).
		MaxHeight(height).
		Render(strings.Join(lines, "\n"))
}

func (m Model) renderArtists(width, height int) string {
	title := "Artists"
	if len(m.artists) > 0 {
		title = m.artists[m.selectedArtist].Name
	}

	lines := []string{
		paneTitle(title, m.focused == songsPane),
		ui.Subtitle.Render(fmt.Sprintf("%d albums, %d songs", len(m.albums), len(m.songs))),
		"",
	}
	if m.filterQueryFor(songsPane) != "" || m.searching && m.searchPane == songsPane {
		lines = append(lines, searchLine(m.searchQuery))
	}

	if m.err != nil {
		lines = append(lines, ui.Error.Render(m.err.Error()))
	} else if m.loading && len(m.songs) == 0 {
		lines = append(lines, ui.Subtitle.Render("Loading..."))
	} else if len(m.artists) == 0 {
		lines = append(lines, ui.Subtitle.Render("No artists found"))
	}

	rows := m.artistRows()
	if len(rows) == 0 && m.filterQueryFor(songsPane) != "" {
		lines = append(lines, ui.Subtitle.Render("No matches"))
	}

	selectedRow := selectedArtistRow(rows, m.selectedSong)
	start, end := visibleRange(selectedRow, len(rows), m.visibleSongRows(height))
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

func (m Model) artistRows() []artistRow {
	rows := make([]artistRow, 0, len(m.albums)+len(m.songs))
	songIndex := 0
	query := m.filterQueryFor(songsPane)
	for _, group := range m.albums {
		albumStart := len(rows)
		rows = append(rows, artistRow{albumTitle: group.album.Name, songIndex: -1})
		for _, song := range group.songs {
			if query == "" || songMatches(song, query) {
				rows = append(rows, artistRow{song: song, songIndex: songIndex})
			}
			songIndex++
		}
		if len(rows) == albumStart+1 {
			rows = rows[:albumStart]
		}
	}

	return rows
}

func selectedArtistRow(rows []artistRow, selectedSong int) int {
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

func (m Model) visibleSongRows(height int) int {
	if m.filterQueryFor(songsPane) != "" || m.searching && m.searchPane == songsPane {
		return max(height-6, 1)
	}

	return max(height-5, 1)
}
