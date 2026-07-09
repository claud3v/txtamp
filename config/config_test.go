package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveCredentials(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".txtamp", "config.env")
	keyPath := filepath.Join(dir, ".txtamp", ".key")

	err := saveCredentials(envPath, keyPath, Credentials{
		Alias:    "home",
		Host:     "https://music.example.com",
		Username: "john",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	contents, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("expected env file: %v", err)
	}

	text := string(contents)
	if !strings.Contains(text, `alias="home"`) {
		t.Fatalf("expected alias in env file, got:\n%s", text)
	}
	if !strings.Contains(text, `host="https://music.example.com"`) {
		t.Fatalf("expected host in env file, got:\n%s", text)
	}
	if !strings.Contains(text, `username="john"`) {
		t.Fatalf("expected username in env file, got:\n%s", text)
	}
	if strings.Contains(text, "secret") {
		t.Fatalf("expected password to be encrypted, got:\n%s", text)
	}
	if !strings.Contains(text, `password="v1:`) {
		t.Fatalf("expected encrypted password in env file, got:\n%s", text)
	}

	key, err := os.ReadFile(keyPath)
	if err != nil {
		t.Fatalf("expected key file: %v", err)
	}
	if len(strings.TrimSpace(string(key))) != 64 {
		t.Fatalf("expected 32-byte hex key, got %q", string(key))
	}

	credentials, found, err := loadCredentials(envPath, keyPath)
	if err != nil {
		t.Fatalf("expected load to succeed, got %v", err)
	}
	if !found {
		t.Fatal("expected credentials to be found")
	}
	if credentials.Alias != "home" {
		t.Fatalf("expected alias home, got %q", credentials.Alias)
	}
	if credentials.Host != "https://music.example.com" {
		t.Fatalf("expected host, got %q", credentials.Host)
	}
	if credentials.Username != "john" {
		t.Fatalf("expected username john, got %q", credentials.Username)
	}
	if credentials.Password != "secret" {
		t.Fatalf("expected decrypted password secret, got %q", credentials.Password)
	}
}

func TestSaveCredentialsOmitsBlankAlias(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".txtamp", "config.env")
	keyPath := filepath.Join(dir, ".txtamp", ".key")

	err := saveCredentials(envPath, keyPath, Credentials{
		Host:     "https://music.example.com",
		Username: "john",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	contents, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("expected env file: %v", err)
	}

	if strings.Contains(string(contents), "alias=") {
		t.Fatalf("expected blank alias to be omitted, got:\n%s", string(contents))
	}
}

func TestLoadCredentialsReturnsNotFoundForMissingFile(t *testing.T) {
	dir := t.TempDir()

	credentials, found, err := loadCredentials(
		filepath.Join(dir, ".txtamp", "config.env"),
		filepath.Join(dir, ".txtamp", ".key"),
	)
	if err != nil {
		t.Fatalf("expected missing credentials to not error, got %v", err)
	}
	if found {
		t.Fatalf("expected credentials not to be found, got %+v", credentials)
	}
}
