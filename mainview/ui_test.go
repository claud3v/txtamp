package mainview

import (
	"testing"
	"txtamp/navidrome"

	tea "charm.land/bubbletea/v2"
)

func TestPlaylistSelectionResetsSongSelection(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane
	m.selectedSong = 2
	m.focused = playlistsPane

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(Model)

	if m.selectedPlaylist != 1 {
		t.Fatalf("expected second playlist to be selected, got %d", m.selectedPlaylist)
	}
	if m.selectedSong != 0 {
		t.Fatalf("expected song selection to reset, got %d", m.selectedSong)
	}
}

func TestSongNavigation(t *testing.T) {
	m := loadedModel()

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	m = updated.(Model)

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(Model)

	if m.focused != songsPane {
		t.Fatalf("expected songs pane to be focused, got %v", m.focused)
	}
	if m.selectedSong != 1 {
		t.Fatalf("expected second song to be selected, got %d", m.selectedSong)
	}
}

func TestSpaceStartsAndTogglesPlayback(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	m = updated.(Model)

	if m.currentSong == nil {
		t.Fatal("expected selected song to start playing")
	}
	if m.paused {
		t.Fatal("expected playback to start unpaused")
	}

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	m = updated.(Model)

	if !m.paused {
		t.Fatal("expected playback to pause")
	}
}

func TestPlaylistsLoadedLoadsFirstPlaylist(t *testing.T) {
	m := New("home", navidrome.Client{})

	updated, cmd := m.Update(playlistsLoadedMsg{
		playlists: []navidrome.Playlist{
			{ID: "playlist-1", Name: "Favorites"},
		},
	})
	m = updated.(Model)

	if m.err != nil {
		t.Fatalf("expected no error, got %v", m.err)
	}
	if len(m.playlists) != 1 {
		t.Fatalf("expected playlists to be stored, got %d", len(m.playlists))
	}
	if cmd == nil {
		t.Fatal("expected first playlist load command")
	}
}

func loadedModel() Model {
	m := New("home", navidrome.Client{})
	m.playlists = []navidrome.Playlist{
		{ID: "playlist-1", Name: "Favorites"},
		{ID: "playlist-2", Name: "Road Trip"},
	}
	m.songs = []navidrome.Song{
		{ID: "song-1", Title: "Dream On", Artist: "Aerosmith", Duration: 268},
		{ID: "song-2", Title: "Sweet Emotion", Artist: "Aerosmith", Duration: 274},
		{ID: "song-3", Title: "Walk This Way", Artist: "Aerosmith", Duration: 220},
	}
	m.loading = false

	return m
}
