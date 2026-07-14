package mainview

import (
	"strings"
	"testing"
	"txtamp/navidrome"

	tea "charm.land/bubbletea/v2"
)

func TestAddSelectedSongToQueue(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane
	m.selectedSong = 1

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected toast clear command")
	}
	if len(m.queue) != 1 {
		t.Fatalf("expected one queued song, got %d", len(m.queue))
	}
	if m.queue[0].ID != "song-2" {
		t.Fatalf("expected Sweet Emotion queued, got %+v", m.queue[0])
	}
	if !strings.Contains(m.toast, "Sweet Emotion") {
		t.Fatalf("expected toast to name queued song, got %q", m.toast)
	}
}

func TestToggleQueueShowsQueueView(t *testing.T) {
	m := loadedModel()
	m.queue = []navidrome.Song{m.songs[0]}

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected no command")
	}
	if m.contentMode != queueContent {
		t.Fatalf("expected queue content, got %v", m.contentMode)
	}

	content := m.renderMainArea(80, 12)
	if !strings.Contains(content, "Queue") || !strings.Contains(content, "Dream On") {
		t.Fatalf("expected queue view, got:\n%s", content)
	}
}

func TestRemoveSelectedQueueSong(t *testing.T) {
	m := loadedModel()
	m.contentMode = queueContent
	m.queue = []navidrome.Song{m.songs[0], m.songs[1]}
	m.selectedQueue = 0

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'd', Text: "d"})
	m = updated.(Model)

	if len(m.queue) != 1 {
		t.Fatalf("expected one queued song, got %d", len(m.queue))
	}
	if m.queue[0].ID != "song-2" {
		t.Fatalf("expected Sweet Emotion to remain, got %+v", m.queue[0])
	}
	if !strings.Contains(m.toast, "Dream On") {
		t.Fatalf("expected removal toast to name song, got %q", m.toast)
	}
}

func TestClearQueue(t *testing.T) {
	m := loadedModel()
	m.queue = []navidrome.Song{m.songs[0], m.songs[1]}
	m.selectedQueue = 1

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected save and toast command")
	}
	if len(m.queue) != 0 {
		t.Fatalf("expected queue to clear, got %+v", m.queue)
	}
	if m.selectedQueue != 0 {
		t.Fatalf("expected selected queue to reset, got %d", m.selectedQueue)
	}
	if !m.queueDirty {
		t.Fatal("expected queue to be marked dirty")
	}
	if !strings.Contains(m.toast, "Queue cleared") {
		t.Fatalf("expected queue clear toast, got %q", m.toast)
	}
}

func TestPlayQueueFromTopConsumesFirstSong(t *testing.T) {
	m := loadedModel()
	m.queue = []navidrome.Song{m.songs[0], m.songs[1]}
	m.selectedQueue = 1

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'P', Text: "P"})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected play command")
	}
	if len(m.queue) != 1 {
		t.Fatalf("expected first queued song to be consumed, got %+v", m.queue)
	}
	if m.queue[0].ID != "song-2" {
		t.Fatalf("expected second song to remain queued, got %+v", m.queue)
	}
	if m.playbackID != 1 {
		t.Fatalf("expected playback id to increment, got %d", m.playbackID)
	}
	if !strings.Contains(m.toast, "Playing queue") {
		t.Fatalf("expected playing queue toast, got %q", m.toast)
	}
}

func TestMoveQueuedSong(t *testing.T) {
	m := loadedModel()
	m.contentMode = queueContent
	m.queue = []navidrome.Song{m.songs[0], m.songs[1]}
	m.selectedQueue = 0

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'J', Text: "J"})
	m = updated.(Model)

	if m.selectedQueue != 1 {
		t.Fatalf("expected selected queue index 1, got %d", m.selectedQueue)
	}
	if m.queue[0].ID != "song-2" || m.queue[1].ID != "song-1" {
		t.Fatalf("expected queued songs to swap, got %+v", m.queue)
	}
}

func TestEnterPlaysSelectedQueueSongWithoutRemovingIt(t *testing.T) {
	m := loadedModel()
	m.contentMode = queueContent
	m.focused = songsPane
	m.queue = []navidrome.Song{m.songs[0], m.songs[1]}
	m.selectedQueue = 1

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected play command")
	}
	if len(m.queue) != 2 {
		t.Fatalf("expected queue to keep both songs after play, got %d", len(m.queue))
	}
	if m.queue[1].ID != "song-2" {
		t.Fatalf("expected selected queue song to remain, got %+v", m.queue)
	}
	if m.playbackID != 1 {
		t.Fatalf("expected playback id to increment, got %d", m.playbackID)
	}
	if m.playbackSource != "Queue" {
		t.Fatalf("expected queue playback source, got %q", m.playbackSource)
	}
}

func TestNextSongConsumesQueueBeforePlaylist(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.currentSongIndex = 0
	m.queue = []navidrome.Song{
		{ID: "queued-song", Title: "Queued Song", Duration: 200},
	}

	cmd := m.playNextSong()
	if cmd == nil {
		t.Fatal("expected queued song play command")
	}
	if len(m.queue) != 0 {
		t.Fatalf("expected queue to be consumed, got %d songs", len(m.queue))
	}
	if m.playbackID != 1 {
		t.Fatalf("expected playback id to increment, got %d", m.playbackID)
	}
	if m.selectedSong != 0 {
		t.Fatalf("expected playlist selection to stay put, got %d", m.selectedSong)
	}
}

func TestQueueLoadedStoresSavedQueue(t *testing.T) {
	m := loadedModel()
	savedQueue := []navidrome.Song{{ID: "saved-song", Title: "Saved Song"}}

	updated, _ := m.Update(queueLoadedMsg{songs: savedQueue, found: true})
	m = updated.(Model)

	if len(m.queue) != 1 || m.queue[0].ID != "saved-song" {
		t.Fatalf("expected saved queue to load, got %+v", m.queue)
	}
}

func TestQueueLoadedDoesNotOverwriteDirtyQueue(t *testing.T) {
	m := loadedModel()
	m.queue = []navidrome.Song{{ID: "local-song", Title: "Local Song"}}
	m.queueDirty = true

	updated, _ := m.Update(queueLoadedMsg{
		songs: []navidrome.Song{{ID: "saved-song", Title: "Saved Song"}},
		found: true,
	})
	m = updated.(Model)

	if len(m.queue) != 1 || m.queue[0].ID != "local-song" {
		t.Fatalf("expected local queue to remain, got %+v", m.queue)
	}
}

func TestAddGlobalSearchSongToQueue(t *testing.T) {
	m := loadedModel()
	m.contentMode = globalSearchContent
	m.selectedSearchResult = 0
	m.globalSearchResult = navidrome.SearchResult{
		Songs: []navidrome.Song{{ID: "search-song", Title: "Victory", Duration: 240}},
	}

	m.addSelectedSongToQueue()

	if len(m.queue) != 1 || m.queue[0].ID != "search-song" {
		t.Fatalf("expected search song queued, got %+v", m.queue)
	}
}

func TestToastClearUsesLatestToastID(t *testing.T) {
	m := loadedModel()
	m.showToast("first")
	firstID := m.toastID
	m.showToast("second")

	updated, _ := m.Update(clearToastMsg{toastID: firstID})
	m = updated.(Model)
	if m.toast != "second" {
		t.Fatalf("expected stale toast clear to be ignored, got %q", m.toast)
	}

	updated, _ = m.Update(clearToastMsg{toastID: m.toastID})
	m = updated.(Model)
	if m.toast != "" {
		t.Fatalf("expected current toast to clear, got %q", m.toast)
	}
}

func TestViewOverlaysToast(t *testing.T) {
	m := loadedModel()
	m.width = 120
	m.height = 30
	m.showToast("Added to queue: Sweet Emotion")

	view := m.View()
	if !strings.Contains(view.Content, "Added to queue: Sweet Emotion") {
		t.Fatalf("expected toast overlay, got:\n%s", view.Content)
	}
	if !strings.Contains(view.Content, "Dream On") {
		t.Fatalf("expected toast to overlay existing content, got:\n%s", view.Content)
	}
}

func TestStatusBarShowsQueueCount(t *testing.T) {
	m := loadedModel()
	m.queue = []navidrome.Song{m.songs[0], m.songs[1]}

	content := m.renderStatusBar(100)
	if !strings.Contains(content, "Queue 2") {
		t.Fatalf("expected queue count in status bar, got:\n%s", content)
	}
}
