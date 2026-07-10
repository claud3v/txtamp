package mainview

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestPlaylistSelectionResetsSongSelection(t *testing.T) {
	m := New("home")
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
	m := New("home")

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
	m := New("home")
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
