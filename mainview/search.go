package mainview

import (
	"context"
	"strings"
	"txtamp/navidrome"

	tea "charm.land/bubbletea/v2"
)

type indexedPlaylist struct {
	index    int
	playlist navidrome.Playlist
}

type indexedArtist struct {
	index  int
	artist navidrome.Artist
}

type indexedSong struct {
	index int
	song  navidrome.Song
}

type globalSearchResultKind int

const (
	searchArtistResult globalSearchResultKind = iota
	searchAlbumResult
	searchSongResult
)

type globalSearchRow struct {
	kind   globalSearchResultKind
	artist navidrome.Artist
	album  navidrome.Album
	song   navidrome.Song
}

func (m *Model) startSearch() {
	if m.focused == modeSelectorPane {
		return
	}

	if m.searchPane != m.focused {
		m.searchQuery = ""
	}
	m.searching = true
	m.searchPane = m.focused
}

func (m *Model) startGoToArtist() {
	if m.mode != artistsMode || m.focused != playlistsPane || len(m.filteredArtists()) == 0 {
		return
	}

	m.goToArtistPending = true
}

func (m *Model) handleGoToArtistKey(msg tea.KeyMsg) tea.Cmd {
	m.goToArtistPending = false
	action, ok := actionForKey(msg.String())
	if ok && action == actionCloseDialog {
		return nil
	}
	if msg.String() == "" || len([]rune(msg.String())) != 1 {
		return nil
	}

	return m.goToArtistLetter([]rune(msg.String())[0])
}

func (m *Model) goToArtistLetter(letter rune) tea.Cmd {
	target := artistGroup(string(letter))
	for _, artist := range m.filteredArtists() {
		if artistGroup(artist.artist.Name) == target {
			if m.selectedArtist == artist.index {
				return nil
			}

			m.selectedArtist = artist.index
			m.selectedArtistRow = 0
			m.collapsedAlbums = nil
			m.selectedSong = 0
			m.sidebarMarqueeOffset = 0
			return tea.Batch(m.loadSelectedArtist(), tickSidebarMarquee())
		}
	}

	return nil
}

func (m *Model) startGlobalSearch() {
	m.contentMode = globalSearchContent
	m.focused = songsPane
	m.globalSearching = true
	m.globalSearchQuery = ""
	m.globalSearchSubmittedQuery = ""
	m.globalSearchErr = nil
	m.globalSearchResult = navidrome.SearchResult{}
	m.selectedSearchResult = 0
}

func (m *Model) clearSearch() {
	m.searching = false
	m.searchQuery = ""
}

func (m *Model) handleGlobalSearchKey(msg tea.KeyMsg) tea.Cmd {
	action, ok := actionForKey(msg.String())
	if ok {
		switch action {
		case actionQuit:
			m.player.Stop()
			return tea.Quit
		case actionCloseDialog:
			m.globalSearching = false
			if strings.TrimSpace(m.globalSearchQuery) == "" && m.globalSearchResultCount() == 0 {
				m.contentMode = libraryContent
			}
			return nil
		case actionActivate:
			m.globalSearching = false
			return m.runGlobalSearch()
		}
	}

	switch msg.String() {
	case "backspace":
		m.globalSearchQuery = dropLastRune(m.globalSearchQuery)
	case "space":
		m.globalSearchQuery += " "
	default:
		if msg.String() == "" || len([]rune(msg.String())) != 1 {
			return nil
		}

		m.globalSearchQuery += msg.String()
	}

	return nil
}

func (m *Model) runGlobalSearch() tea.Cmd {
	query := strings.TrimSpace(m.globalSearchQuery)
	if query == "" {
		m.globalSearchResult = navidrome.SearchResult{}
		m.globalSearchErr = nil
		m.globalSearchLoading = false
		m.globalSearchSubmittedQuery = ""
		m.selectedSearchResult = 0
		return nil
	}

	m.globalSearchSubmittedQuery = query
	m.globalSearchLoading = true
	m.globalSearchErr = nil

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), loadTimeout)
		defer cancel()

		result, err := m.client.Search(ctx, query)
		return globalSearchLoadedMsg{query: query, result: result, err: err}
	}
}

func (m Model) globalSearchRows() []globalSearchRow {
	result := m.globalSearchResult
	rows := make([]globalSearchRow, 0, len(result.Artists)+len(result.Albums)+len(result.Songs))
	for _, artist := range result.Artists {
		rows = append(rows, globalSearchRow{kind: searchArtistResult, artist: artist})
	}
	for _, album := range result.Albums {
		rows = append(rows, globalSearchRow{kind: searchAlbumResult, album: album})
	}
	for _, song := range result.Songs {
		rows = append(rows, globalSearchRow{kind: searchSongResult, song: song})
	}

	return rows
}

func (m Model) globalSearchResultCount() int {
	return len(m.globalSearchResult.Artists) + len(m.globalSearchResult.Albums) + len(m.globalSearchResult.Songs)
}

func (m *Model) moveGlobalSearchSelection(delta int) {
	rows := m.globalSearchRows()
	if len(rows) == 0 {
		return
	}

	m.selectedSearchResult = clamp(m.selectedSearchResult+delta, 0, len(rows)-1)
}

func (m *Model) activateGlobalSearchResult() tea.Cmd {
	rows := m.globalSearchRows()
	if len(rows) == 0 {
		return nil
	}

	row := rows[clamp(m.selectedSearchResult, 0, len(rows)-1)]
	if row.kind != searchSongResult {
		return nil
	}

	return m.playSongFromList(row.song, m.globalSearchSongIndex(row.song), m.globalSearchResult.Songs, "Search: "+m.globalSearchSubmittedQuery)
}

func (m Model) globalSearchSongIndex(song navidrome.Song) int {
	for i, result := range m.globalSearchResult.Songs {
		if result.ID == song.ID {
			return i
		}
	}

	return -1
}

func (m *Model) handleSearchKey(msg tea.KeyMsg) tea.Cmd {
	action, ok := actionForKey(msg.String())
	if ok {
		switch action {
		case actionQuit:
			m.player.Stop()
			return tea.Quit
		case actionCloseDialog:
			m.clearSearch()
			return nil
		case actionActivate:
			m.searching = false
			return nil
		}
	}

	switch msg.String() {
	case "backspace":
		m.searchQuery = dropLastRune(m.searchQuery)
	case "space":
		m.searchQuery += " "
	default:
		if msg.String() == "" || len([]rune(msg.String())) != 1 {
			return nil
		}

		m.searchQuery += msg.String()
	}

	selectionChanged := m.selectFirstFilteredItem()
	if selectionChanged && m.searchPane == playlistsPane {
		return m.loadSelectedSidebarItem()
	}

	return nil
}

func (m *Model) selectFirstFilteredItem() bool {
	switch m.searchPane {
	case playlistsPane:
		switch m.mode {
		case playlistsMode:
			playlists := m.filteredPlaylists()
			if len(playlists) > 0 && m.selectedPlaylist != playlists[0].index {
				m.selectedPlaylist = playlists[0].index
				return true
			}
		case artistsMode:
			artists := m.filteredArtists()
			if len(artists) > 0 && m.selectedArtist != artists[0].index {
				m.selectedArtist = artists[0].index
				return true
			}
		}
	case songsPane:
		songs := m.filteredSongs()
		if len(songs) > 0 && m.selectedSong != songs[0].index {
			m.selectedSong = songs[0].index
			return true
		}
	}

	return false
}

func dropLastRune(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return ""
	}

	return string(runes[:len(runes)-1])
}

func (m Model) filteredPlaylists() []indexedPlaylist {
	query := m.filterQueryFor(playlistsPane)
	playlists := make([]indexedPlaylist, 0, len(m.playlists))
	for i, playlist := range m.playlists {
		if query == "" || containsFold(playlist.Name, query) {
			playlists = append(playlists, indexedPlaylist{index: i, playlist: playlist})
		}
	}

	return playlists
}

func (m Model) filteredArtists() []indexedArtist {
	query := m.filterQueryFor(playlistsPane)
	artists := make([]indexedArtist, 0, len(m.artists))
	for i, artist := range m.artists {
		if query == "" || containsFold(artist.Name, query) {
			artists = append(artists, indexedArtist{index: i, artist: artist})
		}
	}

	return artists
}

func (m Model) filteredSongs() []indexedSong {
	query := m.filterQueryFor(songsPane)
	songs := make([]indexedSong, 0, len(m.songs))
	for i, song := range m.songs {
		if query == "" || songMatches(song, query) {
			songs = append(songs, indexedSong{index: i, song: song})
		}
	}

	return songs
}

func songMatches(song navidrome.Song, query string) bool {
	return containsFold(song.Title, query) ||
		containsFold(song.Artist, query) ||
		containsFold(song.Album, query)
}

func containsFold(value, query string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(query))
}

func (m Model) filterQueryFor(pane focusPane) string {
	if m.searchPane != pane {
		return ""
	}

	return strings.TrimSpace(m.searchQuery)
}

func (m Model) selectedSidebarPosition() int {
	switch m.mode {
	case artistsMode:
		return m.selectedArtistPosition(m.filteredArtists())
	default:
		return m.selectedPlaylistPosition(m.filteredPlaylists())
	}
}

func (m Model) selectedPlaylistPosition(playlists []indexedPlaylist) int {
	for i, playlist := range playlists {
		if playlist.index == m.selectedPlaylist {
			return i
		}
	}

	return 0
}

func (m Model) selectedArtistPosition(artists []indexedArtist) int {
	for i, artist := range artists {
		if artist.index == m.selectedArtist {
			return i
		}
	}

	return 0
}

func (m Model) selectedSongPosition(songs []indexedSong) int {
	for i, song := range songs {
		if song.index == m.selectedSong {
			return i
		}
	}

	return 0
}
