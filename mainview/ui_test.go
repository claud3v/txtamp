package mainview

import (
	"testing"
	"txtamp/navidrome"
	"txtamp/player"

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

func TestSpaceReturnsPlayCommand(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected play command")
	}
	if m.currentSong != nil {
		t.Fatal("expected playback state to wait for command result")
	}
}

func TestPlaybackMessageUpdatesCurrentSong(t *testing.T) {
	m := loadedModel()
	song := m.songs[0]

	updated, cmd := m.Update(playbackMsg{song: &song})
	m = updated.(Model)

	if m.currentSong == nil {
		t.Fatal("expected current song")
	}
	if m.currentSong.ID != "song-1" {
		t.Fatalf("expected song-1, got %q", m.currentSong.ID)
	}
	if m.paused {
		t.Fatal("expected playback to be unpaused")
	}
	if m.duration != song.Duration {
		t.Fatalf("expected duration %d, got %d", song.Duration, m.duration)
	}
	if cmd == nil {
		t.Fatal("expected status polling command")
	}
}

func TestSpaceReturnsPauseCommand(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane
	m.currentSong = &m.songs[0]

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected pause command")
	}
	if m.paused {
		t.Fatal("expected pause state to wait for command result")
	}
}

func TestPlayerStatusMessageUpdatesProgress(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]

	updated, cmd := m.Update(playerStatusMsg{
		status: player.Status{
			Elapsed:  42,
			Duration: 268,
			Paused:   true,
		},
	})
	m = updated.(Model)

	if m.elapsed != 42 {
		t.Fatalf("expected elapsed 42, got %d", m.elapsed)
	}
	if m.duration != 268 {
		t.Fatalf("expected duration 268, got %d", m.duration)
	}
	if !m.paused {
		t.Fatal("expected paused")
	}
	if cmd == nil {
		t.Fatal("expected next status tick")
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
		{ID: "song-1", Title: "Dream On", Artist: "Aerosmith", Album: "Aerosmith", Track: 3, Duration: 268},
		{ID: "song-2", Title: "Sweet Emotion", Artist: "Aerosmith", Album: "Toys in the Attic", Track: 1, Duration: 274},
		{ID: "song-3", Title: "Walk This Way", Artist: "Aerosmith", Album: "Toys in the Attic", Track: 4, Duration: 220},
	}
	m.loading = false

	return m
}
