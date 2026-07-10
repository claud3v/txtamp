package mainview

import (
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

func (m *Model) clearSearch() {
	m.searching = false
	m.searchQuery = ""
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
