package mainview

import (
	"strings"
	"txtamp/ui"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

type artistSidebarRow struct {
	header string
	artist indexedArtist
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

		rows := artistSidebarRows(artists)
		selected := selectedArtistSidebarRow(rows, m.selectedArtist)
		start, end := artistSidebarVisibleRange(rows, selected, m.visibleSidebarRows(height))
		lines := make([]string, 0, end-start)
		for i := start; i < end; i++ {
			row := rows[i]
			if row.header != "" {
				lines = append(lines, artistGroupHeader(row.header, width))
				continue
			}

			line := m.sidebarSelectableLine(row.artist.artist.Name, row.artist.index == m.selectedArtist, width)
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
			line := m.sidebarSelectableLine(playlist.playlist.Name, playlist.index == m.selectedPlaylist, width)
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

func artistSidebarRows(artists []indexedArtist) []artistSidebarRow {
	rows := make([]artistSidebarRow, 0, len(artists))
	lastGroup := ""
	for _, artist := range artists {
		group := artistGroup(artist.artist.Name)
		if group != lastGroup {
			rows = append(rows, artistSidebarRow{header: group})
			lastGroup = group
		}
		rows = append(rows, artistSidebarRow{artist: artist})
	}

	return rows
}

func artistGroup(name string) string {
	for _, r := range strings.TrimSpace(name) {
		r = unicode.ToUpper(r)
		if r >= 'A' && r <= 'Z' {
			return string(r)
		}
		return "#"
	}

	return "#"
}

func selectedArtistSidebarRow(rows []artistSidebarRow, selectedArtist int) int {
	for i, row := range rows {
		if row.header == "" && row.artist.index == selectedArtist {
			return i
		}
	}

	return 0
}

func artistSidebarVisibleRange(rows []artistSidebarRow, selected, height int) (int, int) {
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

	if start <= 0 || rows[start].header != "" {
		return start, end
	}

	header := start - 1
	for header >= 0 && rows[header].header == "" {
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

func artistGroupHeader(group string, width int) string {
	line := group + " " + strings.Repeat("-", max(width-len(group)-1, 0))
	return ui.AlbumCollapsed.Width(width).Render(ui.Truncate(line, width))
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
