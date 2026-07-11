package mainview

import (
	"testing"
	"txtamp/navidrome"

	tea "charm.land/bubbletea/v2"
)

func TestFilteredPlaylistsKeepsOriginalIndexes(t *testing.T) {
	m := loadedModel()
	m.searchPane = playlistsPane
	m.searchQuery = "road"

	playlists := m.filteredPlaylists()
	if len(playlists) != 1 {
		t.Fatalf("expected 1 playlist, got %d", len(playlists))
	}
	if playlists[0].index != 1 {
		t.Fatalf("expected original playlist index 1, got %d", playlists[0].index)
	}
	if playlists[0].playlist.Name != "Road Trip" {
		t.Fatalf("expected Road Trip, got %q", playlists[0].playlist.Name)
	}
}

func TestFilteredArtistsKeepsOriginalIndexes(t *testing.T) {
	m := loadedModel()
	m.mode = artistsMode
	m.artists = []navidrome.Artist{
		{ID: "artist-1", Name: "Aerosmith"},
		{ID: "artist-2", Name: "H.E.A.T"},
		{ID: "artist-3", Name: "Heart"},
	}
	m.searchPane = playlistsPane
	m.searchQuery = "h.e.a"

	artists := m.filteredArtists()
	if len(artists) != 1 {
		t.Fatalf("expected 1 artist, got %d", len(artists))
	}
	if artists[0].index != 1 {
		t.Fatalf("expected original artist index 1, got %d", artists[0].index)
	}
	if artists[0].artist.Name != "H.E.A.T" {
		t.Fatalf("expected H.E.A.T, got %q", artists[0].artist.Name)
	}
}

func TestFilteredSongsMatchesTitleArtistAndAlbum(t *testing.T) {
	m := loadedModel()
	m.songs = []navidrome.Song{
		{ID: "song-1", Title: "Dream On", Artist: "Aerosmith", Album: "Aerosmith"},
		{ID: "song-2", Title: "Victory", Artist: "H.E.A.T", Album: "Force Majeure"},
		{ID: "song-3", Title: "Back in Black", Artist: "AC/DC", Album: "Back in Black"},
	}
	m.searchPane = songsPane

	tests := []struct {
		name          string
		query         string
		expectedIndex int
	}{
		{name: "title", query: "victory", expectedIndex: 1},
		{name: "artist", query: "ac/dc", expectedIndex: 2},
		{name: "album", query: "force", expectedIndex: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.searchQuery = test.query
			songs := m.filteredSongs()
			if len(songs) != 1 {
				t.Fatalf("expected 1 song, got %d", len(songs))
			}
			if songs[0].index != test.expectedIndex {
				t.Fatalf("expected original song index %d, got %d", test.expectedIndex, songs[0].index)
			}
		})
	}
}

func TestStartSearchKeepsQueryForSamePane(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane
	m.searchPane = songsPane
	m.searchQuery = "victory"

	m.startSearch()

	if !m.searching {
		t.Fatal("expected search to start")
	}
	if m.searchQuery != "victory" {
		t.Fatalf("expected existing query to be kept, got %q", m.searchQuery)
	}
}

func TestStartSearchClearsQueryForDifferentPane(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane
	m.searchPane = playlistsPane
	m.searchQuery = "h.e.a"

	m.startSearch()

	if !m.searching {
		t.Fatal("expected search to start")
	}
	if m.searchQuery != "" {
		t.Fatalf("expected query to clear, got %q", m.searchQuery)
	}
	if m.searchPane != songsPane {
		t.Fatalf("expected songs pane search, got %v", m.searchPane)
	}
}

func TestDropLastRune(t *testing.T) {
	if got := dropLastRune("abc"); got != "ab" {
		t.Fatalf("expected ab, got %q", got)
	}
	if got := dropLastRune(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestStartGlobalSearch(t *testing.T) {
	m := loadedModel()

	m.startGlobalSearch()

	if m.contentMode != globalSearchContent {
		t.Fatalf("expected global search content, got %v", m.contentMode)
	}
	if !m.globalSearching {
		t.Fatal("expected global search input to be active")
	}
	if m.focused != songsPane {
		t.Fatalf("expected main pane focus, got %v", m.focused)
	}
}

func TestGlobalSearchInputRunsSearchOnEnter(t *testing.T) {
	m := loadedModel()
	m.startGlobalSearch()

	var cmd tea.Cmd
	for _, char := range "victory" {
		cmd = m.handleGlobalSearchKey(tea.KeyPressMsg{Code: char, Text: string(char)})
		if cmd != nil {
			t.Fatal("expected no command while typing")
		}
	}

	cmd = m.handleGlobalSearchKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected search command")
	}
	if m.globalSearching {
		t.Fatal("expected global search input to close")
	}
	if !m.globalSearchLoading {
		t.Fatal("expected global search loading")
	}
	if m.globalSearchSubmittedQuery != "victory" {
		t.Fatalf("expected submitted query victory, got %q", m.globalSearchSubmittedQuery)
	}
}

func TestGlobalSearchRowsAreGrouped(t *testing.T) {
	m := loadedModel()
	m.contentMode = globalSearchContent
	m.globalSearchResult = navidrome.SearchResult{
		Artists: []navidrome.Artist{{ID: "artist-1", Name: "Victory"}},
		Albums:  []navidrome.Album{{ID: "album-1", Name: "Victory Songs", Artist: "Various Artists"}},
		Songs:   []navidrome.Song{{ID: "song-1", Title: "Victory", Artist: "H.E.A.T", Album: "Force Majeure", Duration: 240}},
	}

	rows := m.globalSearchRows()
	if len(rows) != 3 {
		t.Fatalf("expected 3 selectable rows, got %d", len(rows))
	}
	if rows[0].kind != searchArtistResult || rows[1].kind != searchAlbumResult || rows[2].kind != searchSongResult {
		t.Fatalf("unexpected result rows: %+v", rows)
	}
}

func TestActivateGlobalSearchSongPlaysSongOutsideCurrentList(t *testing.T) {
	m := loadedModel()
	m.contentMode = globalSearchContent
	m.focused = songsPane
	m.selectedSearchResult = 0
	m.globalSearchResult = navidrome.SearchResult{
		Songs: []navidrome.Song{{ID: "search-song", Title: "Victory", Duration: 240}},
	}

	cmd := m.activateGlobalSearchResult()
	if cmd == nil {
		t.Fatal("expected play command")
	}
	if m.playbackID != 1 {
		t.Fatalf("expected playback id to increment, got %d", m.playbackID)
	}
}
