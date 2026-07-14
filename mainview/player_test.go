package mainview

import (
	"strings"
	"testing"
	"txtamp/navidrome"
)

func TestFormatNowPlayingUsesAvailableMetadata(t *testing.T) {
	got := formatNowPlaying(navidrome.Song{
		Artist: "AC/DC",
		Album:  "Back in Black",
		Title:  "Hells Bells",
	})

	want := "AC/DC - Back in Black - Hells Bells"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestFormatNowPlayingOmitsMissingAlbum(t *testing.T) {
	got := formatNowPlaying(navidrome.Song{
		Artist: "AC/DC",
		Title:  "Hells Bells",
	})

	want := "AC/DC - Hells Bells"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRenderPlayerShowsUpNext(t *testing.T) {
	m := loadedModel()
	m.queue = []navidrome.Song{{
		Artist: "AC/DC",
		Album:  "Back in Black",
		Title:  "Hells Bells",
	}}

	content := m.renderPlayer(100)
	if !strings.Contains(content, "Up next: Hells Bells") {
		t.Fatalf("expected up next song, got:\n%s", content)
	}
}

func TestRenderPlayerShowsPlaybackSource(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.playbackSource = "Playlist: Favorites"

	content := m.renderPlayer(100)
	if !strings.Contains(content, "Source: Playlist: Favorites") {
		t.Fatalf("expected playback source, got:\n%s", content)
	}
}

func TestRenderPlayerHidesRedundantArtistSource(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.playbackSource = "Artist: Aerosmith"

	content := m.renderPlayer(100)
	if strings.Contains(content, "Source: Artist: Aerosmith") {
		t.Fatalf("expected redundant artist source to be hidden, got:\n%s", content)
	}
}

func TestRenderPlayerHidesSearchSource(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.playbackSource = "Search: aerosmith"

	content := m.renderPlayer(100)
	if strings.Contains(content, "Source: Search: aerosmith") {
		t.Fatalf("expected search source to be hidden, got:\n%s", content)
	}
	if !strings.Contains(content, "Aerosmith - Aerosmith") {
		t.Fatalf("expected song metadata to remain visible, got:\n%s", content)
	}
}

func TestRenderPlayerShowsNextLoadedSongWhenQueueIsEmpty(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.currentSongIndex = 0

	content := m.renderPlayer(100)
	if !strings.Contains(content, "Up next: Sweet Emotion") {
		t.Fatalf("expected next loaded song, got:\n%s", content)
	}
}

func TestRenderPlayerQueueTakesPriorityOverNextLoadedSong(t *testing.T) {
	m := loadedModel()
	m.currentSong = &m.songs[0]
	m.currentSongIndex = 0
	m.queue = []navidrome.Song{{
		Artist: "AC/DC",
		Album:  "Back in Black",
		Title:  "Hells Bells",
	}}

	content := m.renderPlayer(100)
	if !strings.Contains(content, "Up next: Hells Bells") {
		t.Fatalf("expected queued song to win, got:\n%s", content)
	}
}

func TestRenderPlayerShowsEmptyUpNext(t *testing.T) {
	m := loadedModel()

	content := m.renderPlayer(100)
	if !strings.Contains(content, "Up next: -") {
		t.Fatalf("expected empty up next placeholder, got:\n%s", content)
	}
}

func TestProgressBar(t *testing.T) {
	got := progressBar(30, 120, 10)
	if !strings.Contains(got, "━") {
		t.Fatalf("expected filled progress bar, got %q", got)
	}
	if !strings.Contains(got, "─") {
		t.Fatalf("expected empty progress bar, got %q", got)
	}
}
