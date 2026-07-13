package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"txtamp/navidrome"
)

func TestSaveQueue(t *testing.T) {
	dir := t.TempDir()
	queuePath := filepath.Join(dir, ".txtamp", "queue.json")

	err := saveQueue(queuePath, []navidrome.Song{
		{
			ID:       "song-1",
			Title:    "Dream On",
			Artist:   "Aerosmith",
			Album:    "Aerosmith",
			Track:    3,
			Duration: 268,
		},
	})
	if err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	contents, err := os.ReadFile(queuePath)
	if err != nil {
		t.Fatalf("expected queue file: %v", err)
	}

	text := string(contents)
	for _, expected := range []string{`"songs"`, `"id": "song-1"`, `"title": "Dream On"`, `"artist": "Aerosmith"`} {
		if !strings.Contains(text, expected) {
			t.Fatalf("expected %q in queue file, got:\n%s", expected, text)
		}
	}

	songs, found, err := loadQueue(queuePath)
	if err != nil {
		t.Fatalf("expected load to succeed, got %v", err)
	}
	if !found {
		t.Fatal("expected queue to be found")
	}
	if len(songs) != 1 {
		t.Fatalf("expected one song, got %d", len(songs))
	}
	if songs[0].ID != "song-1" || songs[0].Title != "Dream On" || songs[0].Duration != 268 {
		t.Fatalf("unexpected loaded song: %+v", songs[0])
	}
}

func TestLoadQueueReturnsNotFoundForMissingFile(t *testing.T) {
	dir := t.TempDir()

	songs, found, err := loadQueue(filepath.Join(dir, ".txtamp", "queue.json"))
	if err != nil {
		t.Fatalf("expected missing queue to not error, got %v", err)
	}
	if found {
		t.Fatalf("expected queue not to be found, got %+v", songs)
	}
}
