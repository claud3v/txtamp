package mainview

import (
	"fmt"
	"strings"
	"txtamp/navidrome"
	"txtamp/ui"

	tea "charm.land/bubbletea/v2"
)

type artistRow struct {
	album      navidrome.Album
	albumIndex int
	song       navidrome.Song
	songIndex  int
}

type renderedSearchRow struct {
	text        string
	resultIndex int
}

func (m Model) renderMainArea(width, height int) string {
	if m.contentMode == queueContent {
		return m.renderQueue(width, height)
	}

	if m.contentMode == globalSearchContent {
		return m.renderGlobalSearch(width, height)
	}

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

func (m Model) renderGlobalSearch(width, height int) string {
	title := "Search"
	lines := []string{
		paneTitle(title, m.focused == songsPane),
		globalSearchLine(m.globalSearchQuery),
		"",
	}

	if m.globalSearchErr != nil {
		lines = append(lines, ui.Error.Render(m.globalSearchErr.Error()))
	} else if m.globalSearchLoading {
		lines = append(lines, ui.Subtitle.Render("Searching..."))
	} else if m.globalSearching && strings.TrimSpace(m.globalSearchQuery) != "" {
		lines = append(lines, ui.Subtitle.Render("Press enter to search"))
	} else if strings.TrimSpace(m.globalSearchQuery) == "" {
		lines = append(lines, ui.Subtitle.Render("Type a query and press enter"))
	} else if m.globalSearchSubmittedQuery != "" && m.globalSearchResultCount() == 0 {
		lines = append(lines, ui.Subtitle.Render("No matches"))
	}

	rows := m.renderGlobalSearchRows(width - 4)
	selectedRow := selectedSearchRenderRow(rows, m.selectedSearchResult)
	start, end := visibleRange(selectedRow, len(rows), m.visibleSongRows(height))
	for _, row := range rows[start:end] {
		lines = append(lines, row.text)
	}

	return ui.MainPane.
		Width(width).
		Height(height).
		MaxHeight(height).
		Render(strings.Join(lines, "\n"))
}

func (m Model) renderGlobalSearchRows(width int) []renderedSearchRow {
	rows := make([]renderedSearchRow, 0, m.globalSearchResultCount()+3)
	resultIndex := 0

	if len(m.globalSearchResult.Artists) > 0 {
		rows = append(rows, renderedSearchRow{text: ui.PaneTitle.Render("Artists"), resultIndex: -1})
		for _, artist := range m.globalSearchResult.Artists {
			rows = append(rows, renderedSearchRow{
				text:        selectableLine(artist.Name, resultIndex == m.selectedSearchResult, m.focused == songsPane, width),
				resultIndex: resultIndex,
			})
			resultIndex++
		}
	}

	if len(m.globalSearchResult.Albums) > 0 {
		if len(rows) > 0 {
			rows = append(rows, renderedSearchRow{text: "", resultIndex: -1})
		}
		rows = append(rows, renderedSearchRow{text: ui.PaneTitle.Render("Albums"), resultIndex: -1})
		for _, album := range m.globalSearchResult.Albums {
			text := album.Name
			if album.Artist != "" {
				text += " - " + album.Artist
			}
			rows = append(rows, renderedSearchRow{
				text:        selectableLine(text, resultIndex == m.selectedSearchResult, m.focused == songsPane, width),
				resultIndex: resultIndex,
			})
			resultIndex++
		}
	}

	if len(m.globalSearchResult.Songs) > 0 {
		if len(rows) > 0 {
			rows = append(rows, renderedSearchRow{text: "", resultIndex: -1})
		}
		rows = append(rows, renderedSearchRow{text: ui.PaneTitle.Render("Songs"), resultIndex: -1})
		for _, song := range m.globalSearchResult.Songs {
			titleWidth := max(width-18, 10)
			title := ui.Truncate(formatSongSearchResult(song), titleWidth)
			line := fmt.Sprintf("%-*"+"s %5s", titleWidth, title, formatDuration(song.Duration))
			rows = append(rows, renderedSearchRow{
				text:        styledLine("  ", line, resultIndex == m.selectedSearchResult, m.focused == songsPane, false, width),
				resultIndex: resultIndex,
			})
			resultIndex++
		}
	}

	return rows
}

func selectedSearchRenderRow(rows []renderedSearchRow, selectedResult int) int {
	for i, row := range rows {
		if row.resultIndex == selectedResult {
			return i
		}
	}

	return 0
}

func formatSongSearchResult(song navidrome.Song) string {
	if song.Artist != "" && song.Album != "" {
		return song.Artist + " - " + song.Album + " - " + song.Title
	}
	if song.Artist != "" {
		return song.Artist + " - " + song.Title
	}

	return song.Title
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

	selected := clamp(m.selectedArtistRow, 0, max(len(rows)-1, 0))
	start, end := visibleRange(selected, len(rows), m.visibleSongRows(height))
	for i := start; i < end; i++ {
		row := rows[i]
		if row.songIndex < 0 {
			lines = append(lines, albumLine(row.album, !m.albumCollapsed(row.albumIndex), i == selected, m.focused == songsPane, width-4))
			continue
		}

		titleWidth := max(width-22, 10)
		title := ui.Truncate(row.song.Title, titleWidth)
		line := fmt.Sprintf("%-*"+"s %5s", titleWidth, title, formatDuration(row.song.Duration))
		line = nestedSongLine(line, row.songIndex, i == selected, m, width-4)
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
	for albumIndex, group := range m.albums {
		albumStart := len(rows)
		rows = append(rows, artistRow{album: group.album, albumIndex: albumIndex, songIndex: -1})
		if m.albumCollapsed(albumIndex) && query == "" {
			songIndex += len(group.songs)
			continue
		}
		for _, song := range group.songs {
			if query == "" || songMatches(song, query) {
				rows = append(rows, artistRow{albumIndex: albumIndex, song: song, songIndex: songIndex})
			}
			songIndex++
		}
		if len(rows) == albumStart+1 {
			rows = rows[:albumStart]
		}
	}

	return rows
}

func (m *Model) moveArtistRowSelection(delta int) {
	rows := m.artistRows()
	if len(rows) == 0 {
		return
	}

	m.selectedArtistRow = clamp(m.selectedArtistRow+delta, 0, len(rows)-1)
	m.syncSelectedSongToArtistRow(rows[m.selectedArtistRow])
}

func (m *Model) activateArtistRow() tea.Cmd {
	rows := m.artistRows()
	if len(rows) == 0 {
		return nil
	}

	m.selectedArtistRow = clamp(m.selectedArtistRow, 0, len(rows)-1)
	row := rows[m.selectedArtistRow]
	if row.songIndex >= 0 {
		m.selectedSong = row.songIndex
		return m.playSongAt(row.songIndex)
	}

	m.toggleAlbum(row.albumIndex)
	return nil
}

func (m *Model) syncSelectedSongToArtistRow(row artistRow) {
	if row.songIndex >= 0 {
		m.selectedSong = row.songIndex
	}
}

func (m Model) selectedArtistAlbumRow() *artistRow {
	if m.contentMode != libraryContent || m.mode != artistsMode {
		return nil
	}

	rows := m.artistRows()
	if len(rows) == 0 {
		return nil
	}

	row := rows[clamp(m.selectedArtistRow, 0, len(rows)-1)]
	if row.songIndex >= 0 {
		return nil
	}

	return &row
}

func (m Model) albumCollapsed(albumIndex int) bool {
	return m.collapsedAlbums != nil && m.collapsedAlbums[albumIndex]
}

func (m *Model) toggleAlbum(albumIndex int) {
	if m.collapsedAlbums == nil {
		m.collapsedAlbums = map[int]bool{}
	}

	m.collapsedAlbums[albumIndex] = !m.collapsedAlbums[albumIndex]
	m.clampSelectedArtistRow()
}

func (m *Model) expandAllAlbums() {
	if m.contentMode != libraryContent || m.mode != artistsMode {
		return
	}

	m.collapsedAlbums = nil
	m.clampSelectedArtistRow()
}

func (m *Model) collapseAllAlbums() {
	if m.contentMode != libraryContent || m.mode != artistsMode || len(m.albums) == 0 {
		return
	}

	m.collapsedAlbums = make(map[int]bool, len(m.albums))
	for i := range m.albums {
		m.collapsedAlbums[i] = true
	}
	m.clampSelectedArtistRow()
}

func (m *Model) clampSelectedArtistRow() {
	rows := m.artistRows()
	if len(rows) == 0 {
		m.selectedArtistRow = 0
		return
	}

	m.selectedArtistRow = clamp(m.selectedArtistRow, 0, len(rows)-1)
}

func albumLine(album navidrome.Album, expanded, selected, focused bool, width int) string {
	prefix := "> "
	if expanded {
		prefix = "v "
	}

	return albumHeaderLine(prefix, formatAlbumRow(album, max(width-len(prefix), 1)), expanded, selected, focused, width)
}

func formatAlbumTitle(album navidrome.Album) string {
	if album.Year > 0 {
		return fmt.Sprintf("%s (%d)", album.Name, album.Year)
	}

	return album.Name
}

func formatAlbumRow(album navidrome.Album, width int) string {
	title := formatAlbumTitle(album)
	metadata := formatAlbumMetadata(album)
	if metadata == "" {
		return ui.Truncate(title, width)
	}

	titleWidth := max(width-len(metadata)-2, 8)
	return fmt.Sprintf("%-*s  %s", titleWidth, ui.Truncate(title, titleWidth), metadata)
}

func formatAlbumMetadata(album navidrome.Album) string {
	parts := []string{}
	if album.SongCount > 0 {
		label := "songs"
		if album.SongCount == 1 {
			label = "song"
		}
		parts = append(parts, fmt.Sprintf("%d %s", album.SongCount, label))
	}
	if album.Duration > 0 {
		parts = append(parts, formatDuration(album.Duration))
	}

	return strings.Join(parts, "  ")
}

func (m Model) visibleSongRows(height int) int {
	if m.filterQueryFor(songsPane) != "" || m.searching && m.searchPane == songsPane {
		return max(height-6, 1)
	}

	return max(height-5, 1)
}
