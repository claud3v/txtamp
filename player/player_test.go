package player

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestStopBeforePlay(t *testing.T) {
	p := New()

	if err := p.Stop(); err != nil {
		t.Fatalf("expected stop before play to succeed, got %v", err)
	}
}

func TestSendWritesCommandToSocket(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "mpv.sock")
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("expected listener, got %v", err)
	}
	defer listener.Close()

	commands := make(chan []any, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			commands <- nil
			return
		}
		defer conn.Close()

		var message struct {
			Command []any `json:"command"`
		}
		if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&message); err != nil {
			commands <- nil
			return
		}

		if err := json.NewEncoder(conn).Encode(map[string]string{"error": "success"}); err != nil {
			commands <- nil
			return
		}

		commands <- message.Command
	}()

	p := &Player{socketPath: socketPath}
	if err := p.TogglePause(); err != nil {
		t.Fatalf("expected toggle pause to send command, got %v", err)
	}

	command := <-commands
	if len(command) != 2 {
		t.Fatalf("expected 2 command parts, got %#v", command)
	}
	if command[0] != "cycle" || command[1] != "pause" {
		t.Fatalf("unexpected command: %#v", command)
	}
}

func TestStatusReadsProperties(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "mpv.sock")
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("expected listener, got %v", err)
	}
	defer listener.Close()

	go serveStatusProperties(t, listener)

	p := &Player{socketPath: socketPath}
	status, err := p.Status()
	if err != nil {
		t.Fatalf("expected status, got %v", err)
	}

	if status.Elapsed != 42 {
		t.Fatalf("expected elapsed 42, got %d", status.Elapsed)
	}
	if status.Duration != 180 {
		t.Fatalf("expected duration 180, got %d", status.Duration)
	}
	if !status.Paused {
		t.Fatal("expected paused")
	}
}

func TestNewUsesTempSocketPath(t *testing.T) {
	p := New()

	if filepath.Clean(filepath.Dir(p.socketPath)) != filepath.Clean(os.TempDir()) {
		t.Fatalf("expected socket in temp dir, got %q", p.socketPath)
	}
	if filepath.Base(p.socketPath) == "" {
		t.Fatalf("expected socket filename, got %q", p.socketPath)
	}
}

func serveStatusProperties(t *testing.T, listener net.Listener) {
	t.Helper()

	responses := map[string]any{
		"time-pos": 42.0,
		"duration": 180.0,
		"pause":    true,
	}

	for range 3 {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		var message struct {
			Command []any `json:"command"`
		}
		if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&message); err != nil {
			conn.Close()
			return
		}

		property, _ := message.Command[1].(string)
		if err := json.NewEncoder(conn).Encode(map[string]any{
			"data":  responses[property],
			"error": "success",
		}); err != nil {
			conn.Close()
			return
		}

		conn.Close()
	}
}
