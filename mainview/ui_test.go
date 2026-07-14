package mainview

import (
	"errors"
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

func TestSidebarDoesNotShowEmptyStateWhenLoadFailed(t *testing.T) {
	m := loadedModel()
	m.playlists = nil
	m.songs = nil
	m.err = errors.New("sending request: context deadline exceeded")
	m.loading = false

	content := m.renderSidebar(32, 18)
	if strings.Contains(content, "No playlists found") {
		t.Fatalf("expected sidebar not to show empty state on error, got:\n%s", content)
	}
}

func TestEnterOnLoadedSidebarItemOnlyMovesFocus(t *testing.T) {
	m := loadedModel()
	m.focused = playlistsPane
	m.loadedPlaylistID = "playlist-1"

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected no reload command")
	}
	if m.focused != songsPane {
		t.Fatalf("expected songs pane focus, got %v", m.focused)
	}
	if m.selectedSong != 0 {
		t.Fatalf("expected first song selected, got %d", m.selectedSong)
	}
}

func TestEnterOnUnloadedSidebarItemLoadsSelection(t *testing.T) {
	m := loadedModel()
	m.focused = playlistsPane
	m.loadedPlaylistID = "playlist-2"

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected reload command")
	}
	if m.focused != songsPane {
		t.Fatalf("expected songs pane focus, got %v", m.focused)
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

func TestPlaylistPlaybackSetsSource(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane
	m.selectedPlaylist = 0

	cmd := m.playSongAt(0)
	if cmd == nil {
		t.Fatal("expected play command")
	}
	if m.playbackSource != "Playlist: Favorites" {
		t.Fatalf("expected playlist source, got %q", m.playbackSource)
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

func TestSeekForwardKeyReturnsSeekCommand(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.elapsed = 20
	m.duration = 100

	updated, cmd := m.Update(tea.KeyPressMsg{Code: '.', Text: "."})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected seek command")
	}
	if m.elapsed != 20 {
		t.Fatalf("expected elapsed to wait for seek result, got %d", m.elapsed)
	}
	if m.toast != "+10s" {
		t.Fatalf("expected seek toast, got %q", m.toast)
	}
}

func TestSeekBackwardKeyReturnsSeekCommand(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.elapsed = 20
	m.duration = 100

	updated, cmd := m.Update(tea.KeyPressMsg{Code: ',', Text: ","})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected seek command")
	}
	if m.toast != "-10s" {
		t.Fatalf("expected seek toast, got %q", m.toast)
	}
}

func TestSeekKeyWithoutCurrentSongDoesNothing(t *testing.T) {
	m := loadedModel()

	updated, cmd := m.Update(tea.KeyPressMsg{Code: '.', Text: "."})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected no seek command")
	}
	if m.elapsed != 0 {
		t.Fatalf("expected elapsed to stay 0, got %d", m.elapsed)
	}
}

func TestSeekMessageUpdatesElapsed(t *testing.T) {
	m := loadedModel()
	m.elapsed = 20

	updated, _ := m.Update(seekMsg{elapsed: 30})
	m = updated.(Model)

	if m.elapsed != 30 {
		t.Fatalf("expected elapsed 30, got %d", m.elapsed)
	}
}

func TestSeekRelativeClampsElapsed(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.elapsed = 5
	m.duration = 100

	cmd := m.seekRelative(-seekStep)
	if cmd == nil {
		t.Fatal("expected seek command")
	}
	msg := cmd().(seekMsg)
	if msg.elapsed != 0 {
		t.Fatalf("expected backward seek to clamp to 0, got %d", msg.elapsed)
	}

	m.elapsed = 95
	cmd = m.seekRelative(seekStep)
	msg = cmd().(seekMsg)
	if msg.elapsed != 100 {
		t.Fatalf("expected forward seek to clamp to duration, got %d", msg.elapsed)
	}
}

func TestStopPlaybackKeyReturnsStopCommand(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.playbackID = 4

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected stop command")
	}
	if m.playbackID != 5 {
		t.Fatalf("expected playback id to increment, got %d", m.playbackID)
	}
	if m.toast != "Stopped" {
		t.Fatalf("expected stop toast, got %q", m.toast)
	}
}

func TestStopMessageClearsPlaybackState(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.elapsed = 42
	m.duration = 268
	m.paused = true

	updated, _ := m.Update(stopMsg{})
	m = updated.(Model)

	if m.currentSong != nil {
		t.Fatalf("expected current song to clear, got %+v", m.currentSong)
	}
	if m.elapsed != 0 || m.duration != 0 || m.paused {
		t.Fatalf("expected playback state to clear, got elapsed=%d duration=%d paused=%v", m.elapsed, m.duration, m.paused)
	}
}

func TestStaleStatusAfterStopDoesNotAdvance(t *testing.T) {
	m := loadedModel()
	m.currentSong = nil
	m.playbackID = 2

	updated, cmd := m.Update(playerStatusMsg{
		playbackID: 1,
		err:        player.ErrNotRunning,
	})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected stale stopped status to be ignored")
	}
	if m.selectedSong != 0 {
		t.Fatalf("expected selected song to stay put, got %d", m.selectedSong)
	}
}

func TestVolumeKeysReturnVolumeCommands(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]

	updated, cmd := m.Update(tea.KeyPressMsg{Code: '+', Text: "+"})
	m = updated.(Model)
	if cmd == nil {
		t.Fatal("expected volume up command")
	}
	if m.toast != "Volume Up" {
		t.Fatalf("expected volume up toast, got %q", m.toast)
	}

	updated, cmd = m.Update(tea.KeyPressMsg{Code: '-', Text: "-"})
	m = updated.(Model)
	if cmd == nil {
		t.Fatal("expected volume down command")
	}
	if m.toast != "Volume Down" {
		t.Fatalf("expected volume down toast, got %q", m.toast)
	}
}

func TestVolumeKeyWithoutCurrentSongDoesNothing(t *testing.T) {
	m := loadedModel()

	updated, cmd := m.Update(tea.KeyPressMsg{Code: '+', Text: "+"})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected no volume command")
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

func TestSongsLoadedStoresLoadedPlaylist(t *testing.T) {
	m := loadedModel()
	m.selectedPlaylist = 1

	updated, _ := m.Update(songsLoadedMsg{
		playlistID: "playlist-2",
		songs:      []navidrome.Song{{ID: "song-4", Title: "Road Song"}},
	})
	m = updated.(Model)

	if m.loadedPlaylistID != "playlist-2" {
		t.Fatalf("expected loaded playlist-2, got %q", m.loadedPlaylistID)
	}
	if len(m.songs) != 1 || m.songs[0].ID != "song-4" {
		t.Fatalf("expected new songs to be stored, got %+v", m.songs)
	}
}

func TestStaleSongsLoadedIsIgnored(t *testing.T) {
	m := loadedModel()
	m.selectedPlaylist = 1
	originalSongs := m.songs

	updated, _ := m.Update(songsLoadedMsg{
		playlistID: "playlist-1",
		songs:      []navidrome.Song{{ID: "stale", Title: "Stale Song"}},
	})
	m = updated.(Model)

	if m.loadedPlaylistID == "playlist-1" {
		t.Fatal("expected stale playlist result to be ignored")
	}
	if len(m.songs) != len(originalSongs) || m.songs[0].ID != originalSongs[0].ID {
		t.Fatalf("expected songs to stay unchanged, got %+v", m.songs)
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

func TestSearchFiltersSidebarAndSelectsOriginalPlaylist(t *testing.T) {
	m := loadedModel()
	m.focused = playlistsPane

	updated, _ := m.Update(tea.KeyPressMsg{Code: '/', Text: "/"})
	m = updated.(Model)

	var cmd tea.Cmd
	sawCmd := false
	for _, char := range "road" {
		updated, cmd = m.Update(tea.KeyPressMsg{Code: char, Text: string(char)})
		m = updated.(Model)
		if cmd != nil {
			sawCmd = true
		}
	}

	if !m.searching {
		t.Fatal("expected search to stay active")
	}
	if m.searchQuery != "road" {
		t.Fatalf("expected search query road, got %q", m.searchQuery)
	}
	if m.selectedPlaylist != 1 {
		t.Fatalf("expected Road Trip original playlist index 1, got %d", m.selectedPlaylist)
	}
	if !sawCmd {
		t.Fatal("expected selected playlist load command")
	}

	content := m.renderSidebar(32, 18)
	if !strings.Contains(content, "Road Trip") {
		t.Fatalf("expected filtered playlist to render, got:\n%s", content)
	}
	if strings.Contains(content, "Favorites") {
		t.Fatalf("expected non-matching playlist to be hidden, got:\n%s", content)
	}
}

func TestArtistSidebarGroupsByLetter(t *testing.T) {
	m := loadedModel()
	m.mode = artistsMode
	m.artists = []navidrome.Artist{
		{ID: "artist-1", Name: "AC/DC"},
		{ID: "artist-2", Name: "Aerosmith"},
		{ID: "artist-3", Name: "Black Sabbath"},
	}

	content := m.renderSidebar(36, 20)
	for _, expected := range []string{"A -", "AC/DC", "Aerosmith", "B -", "Black Sabbath"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("expected %q in grouped artist sidebar, got:\n%s", expected, content)
		}
	}
}

func TestArtistSidebarScrollsSelectedArtistIntoViewWithGroups(t *testing.T) {
	m := loadedModel()
	m.mode = artistsMode
	m.selectedArtist = 4
	m.artists = []navidrome.Artist{
		{ID: "artist-1", Name: "AC/DC"},
		{ID: "artist-2", Name: "Aerosmith"},
		{ID: "artist-3", Name: "Black Sabbath"},
		{ID: "artist-4", Name: "Cream"},
		{ID: "artist-5", Name: "Dio"},
	}

	content := m.renderSidebar(36, 12)
	if !strings.Contains(content, "> Dio") {
		t.Fatalf("expected selected artist to stay visible, got:\n%s", content)
	}
}

func TestArtistSidebarStartsAtGroupHeaderWhenScrolledIntoGroup(t *testing.T) {
	m := loadedModel()
	m.mode = artistsMode
	m.selectedArtist = 4
	m.artists = []navidrome.Artist{
		{ID: "artist-1", Name: "AC/DC"},
		{ID: "artist-2", Name: "Aerosmith"},
		{ID: "artist-3", Name: "Black Sabbath"},
		{ID: "artist-4", Name: "Derek & The Dominos"},
		{ID: "artist-5", Name: "Eclipse"},
		{ID: "artist-6", Name: "Ed Sheeran"},
		{ID: "artist-7", Name: "Eric Clapton"},
	}

	content := m.renderSidebar(36, 12)
	if strings.Contains(content, "Derek & The Dominos") && !strings.Contains(content, "D -") {
		t.Fatalf("expected visible artist group to include its header, got:\n%s", content)
	}
	if !strings.Contains(content, "E -") {
		t.Fatalf("expected selected artist group header, got:\n%s", content)
	}
}

func TestSelectedArtistMarqueesWhenFocused(t *testing.T) {
	m := loadedModel()
	m.mode = artistsMode
	m.focused = playlistsPane
	m.sidebarMarqueeOffset = 5
	m.artists = []navidrome.Artist{
		{ID: "artist-1", Name: "Eric Clapton & B.B. King"},
	}

	content := m.renderSidebar(24, 14)
	if strings.Contains(content, "Eric Clapton & B.B.") {
		t.Fatalf("expected selected long artist to scroll, got:\n%s", content)
	}
	if !strings.Contains(content, "Clapton") {
		t.Fatalf("expected marquee window to include artist text, got:\n%s", content)
	}
}

func TestSidebarSelectionChangeRestartsMarqueeTick(t *testing.T) {
	m := loadedModel()
	m.mode = artistsMode
	m.focused = playlistsPane
	m.artists = []navidrome.Artist{
		{ID: "artist-1", Name: "AC/DC"},
		{ID: "artist-2", Name: "Eric Clapton & B.B. King"},
	}

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected artist load and marquee tick command")
	}
	if m.sidebarMarqueeOffset != 0 {
		t.Fatalf("expected marquee offset to reset, got %d", m.sidebarMarqueeOffset)
	}
}

func TestMarqueeTextLoops(t *testing.T) {
	if got := marqueeText("abcdef", 3, 0); got != "abc" {
		t.Fatalf("expected abc, got %q", got)
	}
	if got := marqueeText("abcdef", 3, marqueePauseTicks-1); got != "abc" {
		t.Fatalf("expected pause at start, got %q", got)
	}
	if got := marqueeText("abcdef", 3, marqueePauseTicks+2); got != "cde" {
		t.Fatalf("expected cde, got %q", got)
	}
	if got := marqueeText("abcdef", 3, marqueePauseTicks+7); got != "  a" {
		t.Fatalf("expected wrapped text, got %q", got)
	}
}

func TestSearchFiltersSongsAndKeepsOriginalSongIndex(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane

	updated, _ := m.Update(tea.KeyPressMsg{Code: '/', Text: "/"})
	m = updated.(Model)

	for _, char := range "sweet" {
		updated, _ = m.Update(tea.KeyPressMsg{Code: char, Text: string(char)})
		m = updated.(Model)
	}

	if m.selectedSong != 1 {
		t.Fatalf("expected Sweet Emotion original song index 1, got %d", m.selectedSong)
	}

	content := m.renderSongs(80, 12)
	if !strings.Contains(content, "Sweet Emotion") {
		t.Fatalf("expected matching song to render, got:\n%s", content)
	}
	if strings.Contains(content, "Dream On") {
		t.Fatalf("expected non-matching song to be hidden, got:\n%s", content)
	}
}

func TestSearchNewPaneStartsWithEmptyQuery(t *testing.T) {
	m := loadedModel()
	m.searchPane = playlistsPane
	m.searchQuery = "aero"
	m.focused = songsPane

	updated, _ := m.Update(tea.KeyPressMsg{Code: '/', Text: "/"})
	m = updated.(Model)

	if !m.searching {
		t.Fatal("expected search to start")
	}
	if m.searchPane != songsPane {
		t.Fatalf("expected song pane search, got %v", m.searchPane)
	}
	if m.searchQuery != "" {
		t.Fatalf("expected fresh song filter, got %q", m.searchQuery)
	}
}

func TestSearchEscapeClearsFilter(t *testing.T) {
	m := loadedModel()
	m.focused = songsPane
	m.searching = true
	m.searchPane = songsPane
	m.searchQuery = "sweet"

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	m = updated.(Model)

	if m.searching {
		t.Fatal("expected search to stop")
	}
	if m.searchQuery != "" {
		t.Fatalf("expected search query to clear, got %q", m.searchQuery)
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
			album: navidrome.Album{ID: "album-1", Name: "Aerosmith", Year: 1973, SongCount: 1, Duration: 268},
			songs: []navidrome.Song{
				{ID: "song-1", Title: "Dream On", Duration: 268},
			},
		},
		{
			album: navidrome.Album{ID: "album-2", Name: "Toys in the Attic", Year: 1975, SongCount: 1, Duration: 274},
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
	for _, expected := range []string{"v Aerosmith (1973)", "1 song", "4:28", "Dream On", "v Toys in the Attic (1975)", "Sweet Emotion"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("expected %q in artist view, got:\n%s", expected, content)
		}
	}
}

func TestArtistSearchKeepsMatchingAlbumGroups(t *testing.T) {
	m := loadedModel()
	m.mode = artistsMode
	m.focused = songsPane
	m.searchPane = songsPane
	m.searchQuery = "emotion"
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
	if !strings.Contains(content, "v Toys in the Attic") || !strings.Contains(content, "Sweet Emotion") {
		t.Fatalf("expected matching album group to render, got:\n%s", content)
	}
	if strings.Contains(content, "v Aerosmith") || strings.Contains(content, "Dream On") {
		t.Fatalf("expected non-matching album group to be hidden, got:\n%s", content)
	}
}

func TestArtistNavigationSelectsAlbumRows(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(Model)
	if m.selectedArtistRow != 1 {
		t.Fatalf("expected first song row selected, got %d", m.selectedArtistRow)
	}
	if m.selectedSong != 0 {
		t.Fatalf("expected first song selected, got %d", m.selectedSong)
	}

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(Model)
	if m.selectedArtistRow != 2 {
		t.Fatalf("expected second album row selected, got %d", m.selectedArtistRow)
	}
	if m.selectedSong != 0 {
		t.Fatalf("expected selected song to stay on previous song for album row, got %d", m.selectedSong)
	}
}

func TestEnterOnArtistAlbumTogglesCollapse(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane
	m.selectedArtistRow = 2

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected album activation to toggle without playing")
	}
	if !m.albumCollapsed(1) {
		t.Fatal("expected album to collapse")
	}
	if m.selectedArtistRow != 2 {
		t.Fatalf("expected album row to stay selected, got %d", m.selectedArtistRow)
	}
}

func TestCollapsedArtistAlbumHidesSongs(t *testing.T) {
	m := artistAlbumModel()
	m.collapsedAlbums = map[int]bool{1: true}

	content := m.renderMainArea(80, 12)
	if !strings.Contains(content, "> Toys in the Attic (1975)") {
		t.Fatalf("expected collapsed album marker, got:\n%s", content)
	}
	if strings.Contains(content, "Sweet Emotion") {
		t.Fatalf("expected collapsed album songs to be hidden, got:\n%s", content)
	}
}

func TestFormatAlbumTitleIncludesYearWhenPresent(t *testing.T) {
	if got := formatAlbumTitle(navidrome.Album{Name: "Toys in the Attic", Year: 1975}); got != "Toys in the Attic (1975)" {
		t.Fatalf("unexpected album title: %q", got)
	}
	if got := formatAlbumTitle(navidrome.Album{Name: "Aerosmith"}); got != "Aerosmith" {
		t.Fatalf("unexpected album title without year: %q", got)
	}
}

func TestFormatAlbumRowIncludesMetadata(t *testing.T) {
	album := navidrome.Album{Name: "Toys in the Attic", Year: 1975, SongCount: 9, Duration: 2240}
	row := formatAlbumRow(album, 80)
	for _, expected := range []string{"Toys in the Attic (1975)", "9 songs", "37:20"} {
		if !strings.Contains(row, expected) {
			t.Fatalf("expected %q in album row, got %q", expected, row)
		}
	}
}

func TestAlbumHeaderStylesDifferForExpandedAndCollapsed(t *testing.T) {
	expanded := albumHeaderLine("v ", "Toys in the Attic", true, false, false, 40)
	collapsed := albumHeaderLine("> ", "Toys in the Attic", false, false, false, 40)
	if expanded == collapsed {
		t.Fatal("expected expanded and collapsed album rows to render differently")
	}

	selected := albumHeaderLine("v ", "Toys in the Attic", true, true, true, 40)
	if !strings.Contains(selected, "Toys in the Attic") {
		t.Fatalf("expected selected album row to keep title, got %q", selected)
	}
}

func TestStatusBarShowsAlbumActionsForSelectedAlbum(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane
	m.selectedArtistRow = 0

	content := m.renderStatusBar(120)
	for _, expected := range []string{"Space Play Album", "Enter Toggle", "a Queue Album"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("expected %q in status bar, got:\n%s", expected, content)
		}
	}
}

func TestArtistNavigationSkipsCollapsedAlbumSongs(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane
	m.collapsedAlbums = map[int]bool{0: true}

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(Model)

	if m.selectedArtistRow != 1 {
		t.Fatalf("expected second album row selected, got %d", m.selectedArtistRow)
	}
}

func TestCollapseAllAlbums(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane
	m.selectedArtistRow = 3

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'C', Text: "C"})
	m = updated.(Model)

	if !m.albumCollapsed(0) || !m.albumCollapsed(1) {
		t.Fatalf("expected all albums collapsed, got %+v", m.collapsedAlbums)
	}
	if m.selectedArtistRow != 1 {
		t.Fatalf("expected selected row to clamp to visible album rows, got %d", m.selectedArtistRow)
	}

	content := m.renderMainArea(80, 12)
	if strings.Contains(content, "Dream On") || strings.Contains(content, "Sweet Emotion") {
		t.Fatalf("expected collapsed songs to be hidden, got:\n%s", content)
	}
}

func TestExpandAllAlbums(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane
	m.collapsedAlbums = map[int]bool{0: true, 1: true}

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'E', Text: "E"})
	m = updated.(Model)

	if m.collapsedAlbums != nil {
		t.Fatalf("expected collapsed album state to clear, got %+v", m.collapsedAlbums)
	}

	content := m.renderMainArea(80, 12)
	if !strings.Contains(content, "Dream On") || !strings.Contains(content, "Sweet Emotion") {
		t.Fatalf("expected expanded songs to render, got:\n%s", content)
	}
}

func TestStatusBarDoesNotShowExpandCollapseAll(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane

	content := m.renderStatusBar(120)
	if strings.Contains(content, "Expand") || strings.Contains(content, "Collapse") {
		t.Fatalf("expected expand/collapse all to stay out of status bar, got:\n%s", content)
	}
}

func TestEnterOnArtistSongPlaysSong(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane
	m.selectedArtistRow = 3

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected song play command")
	}
	if m.selectedSong != 1 {
		t.Fatalf("expected Sweet Emotion selected, got %d", m.selectedSong)
	}
}

func TestAddArtistAlbumToQueue(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane
	m.selectedArtistRow = 2

	updated, cmd := m.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected save and toast command")
	}
	if len(m.queue) != 1 {
		t.Fatalf("expected one album song queued, got %+v", m.queue)
	}
	if m.queue[0].ID != "song-2" {
		t.Fatalf("expected Sweet Emotion queued, got %+v", m.queue)
	}
	if !strings.Contains(m.toast, "Added album to queue") {
		t.Fatalf("expected album queue toast, got %q", m.toast)
	}
}

func TestSpaceOnArtistAlbumUsesAlbumAsPlaybackContext(t *testing.T) {
	m := artistAlbumModel()
	m.focused = songsPane
	m.selectedArtistRow = 2
	m.albums[1].songs = []navidrome.Song{
		{ID: "song-2", Title: "Sweet Emotion", Duration: 274},
		{ID: "song-3", Title: "Walk This Way", Duration: 220},
	}
	m.songs = []navidrome.Song{
		{ID: "song-1", Title: "Dream On", Duration: 268},
		{ID: "song-2", Title: "Sweet Emotion", Duration: 274},
		{ID: "song-3", Title: "Walk This Way", Duration: 220},
	}
	m.queue = []navidrome.Song{{ID: "existing", Title: "Existing Queue Song"}}

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	m = updated.(Model)

	if cmd == nil {
		t.Fatal("expected album play command")
	}
	if m.playbackID != 1 {
		t.Fatalf("expected playback id to increment, got %d", m.playbackID)
	}
	if len(m.queue) != 1 || m.queue[0].ID != "existing" {
		t.Fatalf("expected existing queue to stay unchanged, got %+v", m.queue)
	}
	if len(m.playbackSongs) != 2 {
		t.Fatalf("expected album playback context, got %+v", m.playbackSongs)
	}
	if m.playbackSongs[0].ID != "song-2" || m.playbackSongs[1].ID != "song-3" {
		t.Fatalf("expected album songs in playback context, got %+v", m.playbackSongs)
	}
	if m.playbackSource != "Album: Toys in the Attic (1975)" {
		t.Fatalf("expected album playback source, got %q", m.playbackSource)
	}
	if !strings.Contains(m.toast, "Playing album") {
		t.Fatalf("expected playing album toast, got %q", m.toast)
	}
}

func TestUpNextUsesPlaybackContextAfterVisibleSongsChange(t *testing.T) {
	m := loadedModel()
	playing := m.songs[0]
	originalContext := append([]navidrome.Song(nil), m.songs...)
	m.currentSong = &playing
	m.currentSongIndex = 0
	m.playbackSongs = originalContext
	m.songs = []navidrome.Song{
		{ID: "other-1", Title: "Other First", Duration: 100},
		{ID: "other-2", Title: "Other Second", Duration: 100},
	}

	if got := m.upNextText(); got != "Sweet Emotion" {
		t.Fatalf("expected up next from original playback context, got %q", got)
	}
}

func TestGlobalSearchMainAreaRendersGroupedResults(t *testing.T) {
	m := loadedModel()
	m.contentMode = globalSearchContent
	m.focused = songsPane
	m.globalSearchQuery = "victory"
	m.globalSearchResult = navidrome.SearchResult{
		Artists: []navidrome.Artist{{ID: "artist-1", Name: "Victory"}},
		Albums:  []navidrome.Album{{ID: "album-1", Name: "Victory Songs", Artist: "Various Artists"}},
		Songs:   []navidrome.Song{{ID: "song-1", Title: "Victory", Artist: "H.E.A.T", Album: "Force Majeure", Duration: 240}},
	}

	content := m.renderMainArea(100, 18)
	for _, expected := range []string{"Artists", "Victory", "Albums", "Victory Songs", "Songs", "H.E.A.T - Force Majeure - Victory"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("expected %q in global search view, got:\n%s", expected, content)
		}
	}
}

func TestGlobalSearchTypingPromptsForEnter(t *testing.T) {
	m := loadedModel()
	m.contentMode = globalSearchContent
	m.globalSearching = true
	m.globalSearchQuery = "iron"

	content := m.renderMainArea(100, 18)
	if !strings.Contains(content, "Press enter to search") {
		t.Fatalf("expected submit prompt, got:\n%s", content)
	}
	if strings.Contains(content, "No matches") {
		t.Fatalf("expected not to show no matches before submit, got:\n%s", content)
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

func artistAlbumModel() Model {
	m := loadedModel()
	m.mode = artistsMode
	m.artists = []navidrome.Artist{{ID: "artist-1", Name: "Aerosmith"}}
	m.albums = []albumGroup{
		{
			album: navidrome.Album{ID: "album-1", Name: "Aerosmith", Year: 1973, SongCount: 1, Duration: 268},
			songs: []navidrome.Song{
				{ID: "song-1", Title: "Dream On", Duration: 268},
			},
		},
		{
			album: navidrome.Album{ID: "album-2", Name: "Toys in the Attic", Year: 1975, SongCount: 1, Duration: 274},
			songs: []navidrome.Song{
				{ID: "song-2", Title: "Sweet Emotion", Duration: 274},
			},
		},
	}
	m.songs = []navidrome.Song{
		{ID: "song-1", Title: "Dream On", Duration: 268},
		{ID: "song-2", Title: "Sweet Emotion", Duration: 274},
	}

	return m
}
