package mainview

import (
	"strings"
	"testing"
	"txtamp/navidrome"
	"txtamp/player"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
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

	updated, cmd := m.Update(playbackMsg{song: &song, playbackID: 1})
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
	if m.playbackID != 1 {
		t.Fatalf("expected playback ID 1, got %d", m.playbackID)
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
	m.playbackID = 1

	updated, cmd := m.Update(playerStatusMsg{
		playbackID: 1,
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

func TestFinishedStatusPlaysNextSong(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.currentSongIndex = 0
	m.playbackID = 1

	updated, cmd := m.Update(playerStatusMsg{
		playbackID: 1,
		status: player.Status{
			Elapsed:  267,
			Duration: 268,
		},
	})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected next song play command")
	}
	if m.selectedSong != 1 {
		t.Fatalf("expected second song to be selected, got %d", m.selectedSong)
	}
}

func TestStaleFinishedStatusIsIgnored(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[1]
	m.currentSongIndex = 1
	m.selectedSong = 1
	m.playbackID = 2

	updated, cmd := m.Update(playerStatusMsg{
		playbackID: 1,
		err:        player.ErrNotRunning,
	})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected stale status to be ignored")
	}
	if m.selectedSong != 1 {
		t.Fatalf("expected selected song to stay at 1, got %d", m.selectedSong)
	}
}

func TestNextKeyPlaysNextSong(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.currentSongIndex = 0

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected next song play command")
	}
	if m.selectedSong != 1 {
		t.Fatalf("expected second song to be selected, got %d", m.selectedSong)
	}
}

func TestNextAliasPlaysNextSong(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.currentSongIndex = 0

	updated, cmd := m.Update(tea.KeyPressMsg{Code: ']', Text: "]"})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected next song play command")
	}
	if m.selectedSong != 1 {
		t.Fatalf("expected second song to be selected, got %d", m.selectedSong)
	}
}

func TestPreviousKeyRestartsCurrentSong(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[1]
	m.currentSongIndex = 1
	m.selectedSong = 1
	m.elapsed = 10

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected seek command")
	}
	if m.selectedSong != 1 {
		t.Fatalf("expected selection to stay unchanged, got %d", m.selectedSong)
	}
}

func TestPreviousKeyQuickPressPlaysPreviousSong(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[1]
	m.currentSongIndex = 1
	m.selectedSong = 1
	m.elapsed = 1

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'p', Text: "p"})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected previous song play command")
	}
	if m.selectedSong != 0 {
		t.Fatalf("expected first song to be selected, got %d", m.selectedSong)
	}
}

func TestPreviousAliasPlaysPreviousSong(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[1]
	m.currentSongIndex = 1
	m.selectedSong = 1
	m.elapsed = 1

	updated, cmd := m.Update(tea.KeyPressMsg{Code: '[', Text: "["})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected previous song play command")
	}
	if m.selectedSong != 0 {
		t.Fatalf("expected first song to be selected, got %d", m.selectedSong)
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

func TestSwitchToArtistsLoadsArtists(t *testing.T) {
	m := loadedModel()

	updated, cmd := m.Update(tea.KeyPressMsg{Code: '1', Text: "1"})
	m = updated.(Model)

	if m.mode != artistsMode {
		t.Fatalf("expected artists mode, got %v", m.mode)
	}
	if m.focused != playlistsPane {
		t.Fatalf("expected sidebar focus, got %v", m.focused)
	}
	if cmd == nil {
		t.Fatal("expected artist load command")
	}
}

func TestDropdownOpensModeDialog(t *testing.T) {
	m := loadedModel()
	m.focused = modeSelectorPane

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected no command when opening mode dialog")
	}
	if !m.modeDialogOpen {
		t.Fatal("expected mode dialog to open")
	}
	if m.selectedMode != playlistsMode {
		t.Fatalf("expected current mode to be selected, got %v", m.selectedMode)
	}
}

func TestModeDialogAppliesSelectedMode(t *testing.T) {
	m := loadedModel()
	m.focused = modeSelectorPane
	m.modeDialogOpen = true
	m.selectedMode = playlistsMode

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = updated.(Model)

	if m.selectedMode != artistsMode {
		t.Fatalf("expected artists to be selected, got %v", m.selectedMode)
	}

	updated, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if m.modeDialogOpen {
		t.Fatal("expected mode dialog to close")
	}
	if m.mode != artistsMode {
		t.Fatalf("expected artists mode, got %v", m.mode)
	}
	if cmd == nil {
		t.Fatal("expected artist load command")
	}
}

func TestModeDialogViewShowsPicker(t *testing.T) {
	m := loadedModel()
	m.width = 100
	m.height = 30
	m.modeDialogOpen = true
	m.selectedMode = artistsMode

	view := m.View()
	if !strings.Contains(view.Content, "Artists") || !strings.Contains(view.Content, "Playlists") {
		t.Fatalf("expected mode picker options, got:\n%s", view.Content)
	}
	if !strings.Contains(view.Content, "Dream On") {
		t.Fatalf("expected dialog to overlay the existing view, got:\n%s", view.Content)
	}
}

func TestUpFromFirstSidebarItemFocusesDropdown(t *testing.T) {
	m := loadedModel()
	m.focused = playlistsPane
	m.selectedPlaylist = 0

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected no command when moving focus to dropdown")
	}
	if m.focused != modeSelectorPane {
		t.Fatalf("expected mode selector focus, got %v", m.focused)
	}
}

func TestSwitchToPlaylistsLoadsSelectedPlaylist(t *testing.T) {
	m := loadedModel()
	m.mode = artistsMode
	m.artists = []navidrome.Artist{{ID: "artist-1", Name: "Aerosmith"}}

	updated, cmd := m.Update(tea.KeyPressMsg{Code: '2', Text: "2"})
	m = updated.(Model)

	if m.mode != playlistsMode {
		t.Fatalf("expected playlists mode, got %v", m.mode)
	}
	if cmd == nil {
		t.Fatal("expected playlist load command")
	}
}

func TestArtistsLoadedLoadsFirstArtist(t *testing.T) {
	m := New("home", navidrome.Client{})
	m.mode = artistsMode

	updated, cmd := m.Update(artistsLoadedMsg{
		artists: []navidrome.Artist{
			{ID: "artist-1", Name: "Aerosmith"},
		},
	})
	m = updated.(Model)

	if len(m.artists) != 1 {
		t.Fatalf("expected artists to be stored, got %d", len(m.artists))
	}
	if m.selectedArtist != 0 {
		t.Fatalf("expected first artist to be selected, got %d", m.selectedArtist)
	}
	if cmd == nil {
		t.Fatal("expected first artist load command")
	}
}

func TestStalePlaylistLoadDoesNotOverwriteArtistsMode(t *testing.T) {
	m := New("home", navidrome.Client{})
	m.mode = artistsMode
	m.loading = true

	updated, cmd := m.Update(playlistsLoadedMsg{
		playlists: []navidrome.Playlist{
			{ID: "playlist-1", Name: "Favorites"},
		},
	})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected stale playlist load not to trigger a song load")
	}
	if m.mode != artistsMode {
		t.Fatalf("expected artists mode, got %v", m.mode)
	}
	if !m.loading {
		t.Fatal("expected current artists load state to be preserved")
	}
	if len(m.playlists) != 1 {
		t.Fatalf("expected playlists to be cached, got %d", len(m.playlists))
	}
}

func TestArtistMainAreaRendersAlbumGroups(t *testing.T) {
	m := loadedModel()
	m.mode = artistsMode
	m.artists = []navidrome.Artist{{ID: "artist-1", Name: "Aerosmith"}}
	m.albums = []albumGroup{
		{
			album: navidrome.Album{ID: "album-1", Name: "Aerosmith"},
			songs: []navidrome.Song{
				{ID: "song-1", Title: "Dream On", Duration: 268},
			},
		},
		{
			album: navidrome.Album{ID: "album-2", Name: "Toys in the Attic"},
			songs: []navidrome.Song{
				{ID: "song-2", Title: "Sweet Emotion", Duration: 274},
			},
		},
	}
	m.songs = []navidrome.Song{
		{ID: "song-1", Title: "Dream On", Duration: 268},
		{ID: "song-2", Title: "Sweet Emotion", Duration: 274},
	}

	content := m.renderMainArea(80, 12)
	for _, expected := range []string{"> Aerosmith", "Dream On", "> Toys in the Attic", "Sweet Emotion"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("expected %q in artist view, got:\n%s", expected, content)
		}
	}
}

func TestViewFitsTerminalHeight(t *testing.T) {
	m := loadedModel()
	for i := range 100 {
		m.songs = append(m.songs, navidrome.Song{
			ID:       "extra",
			Title:    "Extra Song",
			Artist:   "Artist",
			Duration: 180 + i,
		})
	}
	m.width = 120
	m.height = 30

	view := m.View()
	if got := lipgloss.Height(view.Content); got > m.height {
		t.Fatalf("expected view height <= %d, got %d", m.height, got)
	}
}

func TestViewFillsTerminalHeightWithShortSongList(t *testing.T) {
	m := loadedModel()
	m.width = 120
	m.height = 30

	view := m.View()
	if got := lipgloss.Height(view.Content); got != m.height {
		t.Fatalf("expected view height %d, got %d", m.height, got)
	}
}

func TestSongListScrollsToSelectedSong(t *testing.T) {
	m := loadedModel()
	m.songs = nil
	for i := range 20 {
		m.songs = append(m.songs, navidrome.Song{
			ID:       "song",
			Title:    "Song " + string(rune('A'+i)),
			Artist:   "Artist",
			Duration: 180,
		})
	}
	m.selectedSong = 15

	content := m.renderSongs(80, 10)
	if !strings.Contains(content, "Song P") {
		t.Fatalf("expected selected song to be visible, got:\n%s", content)
	}
	if strings.Contains(content, "Song A") {
		t.Fatalf("expected top songs to scroll out, got:\n%s", content)
	}
}

func TestPlayingSongIsMarked(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[1]
	m.currentSongIndex = 1
	m.selectedSong = 0

	content := m.renderSongs(80, 12)
	if !strings.Contains(content, "* Sweet Emotion") {
		t.Fatalf("expected playing song to be marked, got:\n%s", content)
	}
}

func TestSameIndexDifferentSongIsNotMarkedPlaying(t *testing.T) {
	m := loadedModel()
	playingSong := navidrome.Song{ID: "other-playlist-song", Title: "Other Playlist Song"}
	m.currentSong = &playingSong
	m.currentSongIndex = 1
	m.selectedSong = 0

	content := m.renderSongs(80, 12)
	if strings.Contains(content, "* Sweet Emotion") {
		t.Fatalf("expected same index different song not to be marked, got:\n%s", content)
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
