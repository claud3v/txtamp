package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"txtamp/navidrome"
)

const queueFileName = "queue.json"

type queueFile struct {
	Songs []queueSong `json:"songs"`
}

type queueSong struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist,omitempty"`
	Album    string `json:"album,omitempty"`
	Track    int    `json:"track,omitempty"`
	Duration int    `json:"duration,omitempty"`
}

func SaveQueue(songs []navidrome.Song) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("finding home directory: %w", err)
	}

	return saveQueue(filepath.Join(home, configDirName, queueFileName), songs)
}

func LoadQueue() ([]navidrome.Song, bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, false, fmt.Errorf("finding home directory: %w", err)
	}

	return loadQueue(filepath.Join(home, configDirName, queueFileName))
}

func saveQueue(path string, songs []navidrome.Song) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(queueFile{Songs: queueSongsFromNavidrome(songs)}, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding queue: %w", err)
	}

	if err := os.WriteFile(path, append(data, '\n'), 0600); err != nil {
		return fmt.Errorf("writing queue: %w", err)
	}

	return nil
}

func loadQueue(path string) ([]navidrome.Song, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("reading queue: %w", err)
	}

	var queue queueFile
	if err := json.Unmarshal(data, &queue); err != nil {
		return nil, true, fmt.Errorf("decoding queue: %w", err)
	}

	return navidromeSongsFromQueue(queue.Songs), true, nil
}

func queueSongsFromNavidrome(songs []navidrome.Song) []queueSong {
	queueSongs := make([]queueSong, 0, len(songs))
	for _, song := range songs {
		queueSongs = append(queueSongs, queueSong{
			ID:       song.ID,
			Title:    song.Title,
			Artist:   song.Artist,
			Album:    song.Album,
			Track:    song.Track,
			Duration: song.Duration,
		})
	}

	return queueSongs
}

func navidromeSongsFromQueue(songs []queueSong) []navidrome.Song {
	navidromeSongs := make([]navidrome.Song, 0, len(songs))
	for _, song := range songs {
		navidromeSongs = append(navidromeSongs, navidrome.Song{
			ID:       song.ID,
			Title:    song.Title,
			Artist:   song.Artist,
			Album:    song.Album,
			Track:    song.Track,
			Duration: song.Duration,
		})
	}

	return navidromeSongs
}
