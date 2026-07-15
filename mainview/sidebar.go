package mainview

import (
	"strings"
	"txtamp/navidrome"
	"txtamp/ui"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

type artistSidebarRow struct {
	header string
	artist indexedArtist
}

type albumSidebarRow struct {
	header string
	album  indexedAlbum
}

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
		lines = append(lines, emptyState("Loading..."))
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
	switch mode {
	case artistsMode:
		return "Artists"
	case albumsMode:
		return "Albums"
	default:
		return "Playlists"
	}
}

func (m Model) sidebarItemCount() int {
	switch m.mode {
	case artistsMode:
		return len(m.artists)
	case albumsMode:
		return len(m.albums)
	default:
		return len(m.playlists)
	}
}

func (m Model) sidebarItems(width, height int) []string {
	switch m.mode {
	case artistsMode:
		if m.err != nil {
			return nil
		}
		if len(m.artists) == 0 && !m.loading {
			return []string{emptyState("No artists found")}
		}

		artists := m.filteredArtists()
		if len(artists) == 0 && m.filterQueryFor(playlistsPane) != "" {
			return []string{emptyState("No matches")}
		}

		rows := artistSidebarRows(artists)
		selected := selectedArtistSidebarRow(rows, m.selectedArtist)
		start, end := groupedSidebarVisibleRange(rows, selected, m.visibleSidebarRows(height))
		lines := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			row := rows[i]
			if row.header != "" {
				lines = append(lines, sidebarGroupHeader(row.header, width))
				continue
			}

			line := m.sidebarSelectableLine(row.artist.artist.Name, row.artist.index == m.selectedArtist, width)
			lines = append(lines, line)
		}

		return lines
	case albumsMode:
		if m.err != nil {
			return nil
		}
		if len(m.albums) == 0 && !m.loading {
			return []string{emptyState("No albums found")}
		}

		albums := m.filteredAlbums()
		if len(albums) == 0 && m.filterQueryFor(playlistsPane) != "" {
			return []string{emptyState("No matches")}
		}

		rows := albumSidebarRows(albums)
		selected := selectedAlbumSidebarRow(rows, m.selectedAlbum)
		start, end := groupedSidebarVisibleRange(rows, selected, m.visibleSidebarRows(height))
		lines := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			row := rows[i]
			if row.header != "" {
				lines = append(lines, sidebarGroupHeader(row.header, width))
				continue
			}

			line := m.sidebarSelectableLine(formatSidebarAlbum(row.album.album), row.album.index == m.selectedAlbum, width)
			lines = append(lines, line)
		}

		return lines
	default:
		if m.err != nil {
			return nil
		}
		if len(m.playlists) == 0 && !m.loading {
			return []string{emptyState("No playlists found")}
		}

		playlists := m.filteredPlaylists()
		if len(playlists) == 0 && m.filterQueryFor(playlistsPane) != "" {
			return []string{emptyState("No matches")}
		}

		selected := m.selectedPlaylistPosition(playlists)
		start, end := visibleRange(selected, len(playlists), m.visibleSidebarRows(height))
		lines := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			playlist := playlists[i]
			line := m.sidebarSelectableLine(playlist.playlist.Name, playlist.index == m.selectedPlaylist, width)
			lines = append(lines, line)
		}

		return lines
	}
}

func formatSidebarAlbum(album navidrome.Album) string {
	if album.Artist != "" {
		return album.Name + " - " + album.Artist
	}

	return album.Name
}

func (m Model) visibleSidebarRows(height int) int {
	if m.filterQueryFor(playlistsPane) != "" || m.searching && m.searchPane == playlistsPane {
		return max(height-9, 1)
	}

	return max(height-8, 1)
}

func artistSidebarRows(artists []indexedArtist) []artistSidebarRow {
	rows := make([]artistSidebarRow, 0, len(artists))
	lastGroup := ""
	for _, artist := range artists {
		group := alphaGroup(artist.artist.Name)
		if group != lastGroup {
			rows = append(rows, artistSidebarRow{header: group})
			lastGroup = group
		}
		rows = append(rows, artistSidebarRow{artist: artist})
	}

	return rows
}

func albumSidebarRows(albums []indexedAlbum) []albumSidebarRow {
	rows := make([]albumSidebarRow, 0, len(albums))
	lastGroup := ""
	for _, album := range albums {
		group := alphaGroup(album.album.Name)
		if group != lastGroup {
			rows = append(rows, albumSidebarRow{header: group})
			lastGroup = group
		}
		rows = append(rows, albumSidebarRow{album: album})
	}

	return rows
}

func alphaGroup(name string) string {
	for _, r := range strings.TrimSpace(name) {
		r = unicode.ToUpper(r)
		if r >= 'A' && r <= 'Z' {
			return string(r)
		}
		return "#"
	}

	return "#"
}

func artistGroup(name string) string {
	return alphaGroup(name)
}

func selectedArtistSidebarRow(rows []artistSidebarRow, selectedArtist int) int {
	for i, row := range rows {
		if row.header == "" && row.artist.index == selectedArtist {
			return i
		}
	}

	return 0
}

func selectedAlbumSidebarRow(rows []albumSidebarRow, selectedAlbum int) int {
	for i, row := range rows {
		if row.header == "" && row.album.index == selectedAlbum {
			return i
		}
	}

	return 0
}

type groupedSidebarRow interface {
	artistSidebarRow | albumSidebarRow
}

func groupedSidebarVisibleRange[T groupedSidebarRow](rows []T, selected, height int) (int, int) {
	if len(rows) == 0 || height <= 0 {
		return 0, 0
	}

	if height >= len(rows) {
		return 0, len(rows)
	}

	selected = clamp(selected, 0, len(rows)-1)
	bottomPadding := min(2, height/3)
	start := selected - height + 1 + bottomPadding
	start = clamp(start, 0, max(len(rows)-height, 0))
	end := min(start+height, len(rows))

	if start <= 0 || sidebarRowHeader(rows[start]) != "" {
		return start, end
	}

	header := start - 1
	for header >= 0 && sidebarRowHeader(rows[header]) == "" {
		header--
	}
	if header < 0 {
		return start, end
	}
	if selected >= header+height {
		return start, end
	}

	start = header
	end = min(start+height, len(rows))
	return start, end
}

func sidebarRowHeader(row any) string {
	switch row := row.(type) {
	case artistSidebarRow:
		return row.header
	case albumSidebarRow:
		return row.header
	default:
		return ""
	}
}

func sidebarGroupHeader(group string, width int) string {
	line := group + " " + strings.Repeat("-", min(max(width-len(group)-1, 0), 12))
	return ui.SectionHeader.Width(width).Render(ui.Truncate(line, width))
}

func (m Model) sidebarSelectableLine(text string, selected bool, width int) string {
	displayText := text
	if selected && m.focused == playlistsPane {
		displayText = marqueeText(text, max(width-2, 1), m.sidebarMarqueeOffset)
	}

	return selectableLine(displayText, selected, m.focused == playlistsPane, width)
}

func (m Model) sidebarMarqueeActive() bool {
	if m.focused != playlistsPane || m.searching || m.modeDialogOpen || m.helpOpen {
		return false
	}

	text, ok := m.selectedSidebarText()
	if !ok {
		return false
	}

	layout := ui.NewShellLayout(m.width, m.height)
	return len([]rune(text)) > max(layout.SidebarWidth-6, 1)
}

func (m Model) selectedSidebarText() (string, bool) {
	switch m.mode {
	case artistsMode:
		if len(m.artists) == 0 {
			return "", false
		}
		artist := m.artists[clamp(m.selectedArtist, 0, len(m.artists)-1)]
		return artist.Name, true
	case albumsMode:
		if len(m.albums) == 0 {
			return "", false
		}
		album := m.albums[clamp(m.selectedAlbum, 0, len(m.albums)-1)]
		return formatSidebarAlbum(album), true
	default:
		if len(m.playlists) == 0 {
			return "", false
		}
		playlist := m.playlists[clamp(m.selectedPlaylist, 0, len(m.playlists)-1)]
		return playlist.Name, true
	}
}

func marqueeText(text string, width, offset int) string {
	runes := []rune(text)
	if width <= 0 || len(runes) <= width {
		return text
	}

	padding := []rune("   ")
	loop := append(append([]rune{}, runes...), padding...)
	loop = append(loop, runes...)
	cycleLength := len(runes) + len(padding) + marqueePauseTicks
	position := offset % cycleLength
	start := max(position-marqueePauseTicks, 0)
	window := loop[start : start+width]

	return string(window)
}
