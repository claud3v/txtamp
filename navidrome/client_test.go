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
