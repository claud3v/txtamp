package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestConnectServer(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"subsonic-response":{"status":"ok","version":"1.16.1"}}`))
	}))
	defer server.Close()

	msg := connectServer("home", server.URL, "john", "secret")()
	result, ok := msg.(connectResultMsg)
	if !ok {
		t.Fatalf("expected connectResultMsg, got %T", msg)
	}

	if result.err != nil {
		t.Fatalf("expected connect to succeed, got %v", result.err)
	}

	if result.connectedTo != "home" {
		t.Fatalf("expected connected name home, got %q", result.connectedTo)
	}

	if _, err := os.Stat(filepath.Join(home, ".txtamp", "config.env")); err != nil {
		t.Fatalf("expected credentials file: %v", err)
	}
}

func TestConnectServerReturnsError(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"subsonic-response":{"status":"failed","error":{"code":40,"message":"Wrong username or password"}}}`))
	}))
	defer server.Close()

	msg := connectServer("home", server.URL, "john", "bad-password")()
	result, ok := msg.(connectResultMsg)
	if !ok {
		t.Fatalf("expected connectResultMsg, got %T", msg)
	}

	if result.err == nil {
		t.Fatal("expected connect to fail")
	}

	if _, err := os.Stat(filepath.Join(home, ".txtamp", "config.env")); !os.IsNotExist(err) {
		t.Fatalf("expected no credentials file, got %v", err)
	}
}
