package navidrome

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/ping.view" {
			t.Fatalf("expected /rest/ping.view, got %q", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("u") != "john" {
			t.Fatalf("expected username john, got %q", query.Get("u"))
		}
		if query.Get("t") == "" {
			t.Fatal("expected auth token")
		}
		if query.Get("s") == "" {
			t.Fatal("expected salt")
		}
		if query.Get("v") != apiVersion {
			t.Fatalf("expected api version %q, got %q", apiVersion, query.Get("v"))
		}
		if query.Get("c") != clientName {
			t.Fatalf("expected client name %q, got %q", clientName, query.Get("c"))
		}
		if query.Get("f") != "json" {
			t.Fatalf("expected json format, got %q", query.Get("f"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"subsonic-response":{"status":"ok","version":"1.16.1"}}`))
	}))
	defer server.Close()

	client := Client{
		BaseURL:  server.URL,
		Username: "john",
		Password: "secret",
	}

	if err := client.Ping(context.Background()); err != nil {
		t.Fatalf("expected ping to succeed, got %v", err)
	}
}

func TestPingReturnsSubsonicError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"subsonic-response":{"status":"failed","error":{"code":40,"message":"Wrong username or password"}}}`))
	}))
	defer server.Close()

	client := Client{
		BaseURL:  server.URL,
		Username: "john",
		Password: "bad-password",
	}

	err := client.Ping(context.Background())
	if err == nil {
		t.Fatal("expected ping to fail")
	}

	if err.Error() != "navidrome ping failed: Wrong username or password" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListPlaylists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/getPlaylists.view" {
			t.Fatalf("expected /rest/getPlaylists.view, got %q", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"subsonic-response": {
				"status": "ok",
				"playlists": {
					"playlist": [
						{"id": "playlist-1", "name": "Favorites", "songCount": 2, "duration": 360}
					]
				}
			}
		}`))
	}))
	defer server.Close()

	client := Client{
		BaseURL:  server.URL,
		Username: "john",
		Password: "secret",
	}

	playlists, err := client.ListPlaylists(context.Background())
	if err != nil {
		t.Fatalf("expected playlists to load, got %v", err)
	}
	if len(playlists) != 1 {
		t.Fatalf("expected 1 playlist, got %d", len(playlists))
	}
	if playlists[0].ID != "playlist-1" || playlists[0].Name != "Favorites" {
		t.Fatalf("unexpected playlist: %+v", playlists[0])
	}
}

func TestGetPlaylist(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/getPlaylist.view" {
			t.Fatalf("expected /rest/getPlaylist.view, got %q", r.URL.Path)
		}
		if r.URL.Query().Get("id") != "playlist-1" {
			t.Fatalf("expected playlist id, got %q", r.URL.Query().Get("id"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"subsonic-response": {
				"status": "ok",
				"playlist": {
					"id": "playlist-1",
					"name": "Favorites",
					"entry": [
						{"id": "song-1", "title": "Dream On", "artist": "Aerosmith", "album": "Aerosmith", "track": 3, "duration": 268}
					]
				}
			}
		}`))
	}))
	defer server.Close()

	client := Client{
		BaseURL:  server.URL,
		Username: "john",
		Password: "secret",
	}

	songs, err := client.GetPlaylist(context.Background(), "playlist-1")
	if err != nil {
		t.Fatalf("expected songs to load, got %v", err)
	}
	if len(songs) != 1 {
		t.Fatalf("expected 1 song, got %d", len(songs))
	}
	if songs[0].ID != "song-1" || songs[0].Title != "Dream On" {
		t.Fatalf("unexpected song: %+v", songs[0])
	}
	if songs[0].Album != "Aerosmith" || songs[0].Track != 3 {
		t.Fatalf("expected album and track metadata, got %+v", songs[0])
	}
}

func TestListArtists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/getArtists.view" {
			t.Fatalf("expected /rest/getArtists.view, got %q", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"subsonic-response": {
				"status": "ok",
				"artists": {
					"index": [
						{"name": "A", "artist": [{"id": "artist-1", "name": "Aerosmith", "albumCount": 2}]},
						{"name": "B", "artist": [{"id": "artist-2", "name": "Black Sabbath", "albumCount": 3}]}
					]
				}
			}
		}`))
	}))
	defer server.Close()

	client := Client{
		BaseURL:  server.URL,
		Username: "john",
		Password: "secret",
	}

	artists, err := client.ListArtists(context.Background())
	if err != nil {
		t.Fatalf("expected artists to load, got %v", err)
	}
	if len(artists) != 2 {
		t.Fatalf("expected 2 artists, got %d", len(artists))
	}
	if artists[0].ID != "artist-1" || artists[0].Name != "Aerosmith" || artists[0].AlbumCount != 2 {
		t.Fatalf("unexpected artist: %+v", artists[0])
	}
}

func TestGetArtistAlbums(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/getArtist.view" {
			t.Fatalf("expected /rest/getArtist.view, got %q", r.URL.Path)
		}
		if r.URL.Query().Get("id") != "artist-1" {
			t.Fatalf("expected artist id, got %q", r.URL.Query().Get("id"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"subsonic-response": {
				"status": "ok",
				"artist": {
					"id": "artist-1",
					"name": "Aerosmith",
					"album": [
						{"id": "album-1", "name": "Toys in the Attic", "artist": "Aerosmith", "songCount": 9, "duration": 2240}
					]
				}
			}
		}`))
	}))
	defer server.Close()

	client := Client{
		BaseURL:  server.URL,
		Username: "john",
		Password: "secret",
	}

	albums, err := client.GetArtistAlbums(context.Background(), "artist-1")
	if err != nil {
		t.Fatalf("expected albums to load, got %v", err)
	}
	if len(albums) != 1 {
		t.Fatalf("expected 1 album, got %d", len(albums))
	}
	if albums[0].ID != "album-1" || albums[0].Name != "Toys in the Attic" {
		t.Fatalf("unexpected album: %+v", albums[0])
	}
}

func TestGetAlbumSongs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/getAlbum.view" {
			t.Fatalf("expected /rest/getAlbum.view, got %q", r.URL.Path)
		}
		if r.URL.Query().Get("id") != "album-1" {
			t.Fatalf("expected album id, got %q", r.URL.Query().Get("id"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"subsonic-response": {
				"status": "ok",
				"album": {
					"id": "album-1",
					"name": "Toys in the Attic",
					"song": [
						{"id": "song-1", "title": "Sweet Emotion", "artist": "Aerosmith", "album": "Toys in the Attic", "track": 1, "duration": 274}
					]
				}
			}
		}`))
	}))
	defer server.Close()

	client := Client{
		BaseURL:  server.URL,
		Username: "john",
		Password: "secret",
	}

	songs, err := client.GetAlbumSongs(context.Background(), "album-1")
	if err != nil {
		t.Fatalf("expected songs to load, got %v", err)
	}
	if len(songs) != 1 {
		t.Fatalf("expected 1 song, got %d", len(songs))
	}
	if songs[0].ID != "song-1" || songs[0].Title != "Sweet Emotion" {
		t.Fatalf("unexpected song: %+v", songs[0])
	}
}

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/search3.view" {
			t.Fatalf("expected /rest/search3.view, got %q", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "victory" {
			t.Fatalf("expected search query, got %q", r.URL.Query().Get("query"))
		}
		if r.URL.Query().Get("artistCount") != "20" {
			t.Fatalf("expected artist count, got %q", r.URL.Query().Get("artistCount"))
		}
		if r.URL.Query().Get("albumCount") != "20" {
			t.Fatalf("expected album count, got %q", r.URL.Query().Get("albumCount"))
		}
		if r.URL.Query().Get("songCount") != "50" {
			t.Fatalf("expected song count, got %q", r.URL.Query().Get("songCount"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"subsonic-response": {
				"status": "ok",
				"searchResult3": {
					"artist": [
						{"id": "artist-1", "name": "Victory", "albumCount": 2}
					],
					"album": [
						{"id": "album-1", "name": "Victory Songs", "artist": "Various Artists", "songCount": 10, "duration": 2500}
					],
					"song": [
						{"id": "song-1", "title": "Victory", "artist": "H.E.A.T", "album": "Force Majeure", "track": 4, "duration": 240}
					]
				}
			}
		}`))
	}))
	defer server.Close()

	client := Client{
		BaseURL:  server.URL,
		Username: "john",
		Password: "secret",
	}

	result, err := client.Search(context.Background(), "victory")
	if err != nil {
		t.Fatalf("expected search to load, got %v", err)
	}
	if len(result.Artists) != 1 || result.Artists[0].Name != "Victory" {
		t.Fatalf("unexpected artists: %+v", result.Artists)
	}
	if len(result.Albums) != 1 || result.Albums[0].Name != "Victory Songs" {
		t.Fatalf("unexpected albums: %+v", result.Albums)
	}
	if len(result.Songs) != 1 || result.Songs[0].Title != "Victory" {
		t.Fatalf("unexpected songs: %+v", result.Songs)
	}
}

func TestStreamURL(t *testing.T) {
	client := Client{
		BaseURL:  "https://music.example.com",
		Username: "john",
		Password: "secret",
	}

	streamURL, err := client.StreamURL("song-1")
	if err != nil {
		t.Fatalf("expected stream URL, got %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, streamURL, nil)
	if err != nil {
		t.Fatalf("expected valid URL, got %v", err)
	}

	if req.URL.Path != "/rest/stream.view" {
		t.Fatalf("expected stream endpoint, got %q", req.URL.Path)
	}
	if req.URL.Query().Get("id") != "song-1" {
		t.Fatalf("expected song id, got %q", req.URL.Query().Get("id"))
	}
	if req.URL.Query().Get("u") != "john" {
		t.Fatalf("expected username, got %q", req.URL.Query().Get("u"))
	}
	if req.URL.Query().Get("t") == "" {
		t.Fatal("expected auth token")
	}
}
