package mainview

import (
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

func TestProgressBar(t *testing.T) {
	got := progressBar(30, 120, 10)
	want := "[==------]"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
