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
						{"id": "song-1", "title": "Dream On", "artist": "Aerosmith", "duration": 268}
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
}
