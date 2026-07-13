package player

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"os/exec"
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
	startFakeProcess(t, p)
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

func TestSeekSendsRelativeCommand(t *testing.T) {
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
	startFakeProcess(t, p)
	if err := p.Seek(10); err != nil {
		t.Fatalf("expected seek to send command, got %v", err)
	}

	command := <-commands
	if len(command) != 3 {
		t.Fatalf("expected 3 command parts, got %#v", command)
	}
	if command[0] != "seek" || command[1] != float64(10) || command[2] != "relative" {
		t.Fatalf("unexpected command: %#v", command)
	}
}

func TestAdjustVolumeSendsAddVolumeCommand(t *testing.T) {
	socketPath := filepath.Join(os.TempDir(), "txtamp-test-volume.sock")
	os.Remove(socketPath)
	t.Cleanup(func() {
		os.Remove(socketPath)
	})

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
	startFakeProcess(t, p)
	if err := p.AdjustVolume(5); err != nil {
		t.Fatalf("expected volume adjustment to send command, got %v", err)
	}

	command := <-commands
	if len(command) != 3 {
		t.Fatalf("expected 3 command parts, got %#v", command)
	}
	if command[0] != "add" || command[1] != "volume" || command[2] != float64(5) {
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
	startFakeProcess(t, p)
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

func TestSendReturnsNotRunningWhenPlayerIsStopped(t *testing.T) {
	p := New()

	err := p.TogglePause()
	if err == nil {
		t.Fatal("expected error")
	}
	if err != ErrNotRunning {
		t.Fatalf("expected ErrNotRunning, got %v", err)
	}
}

func TestStatusToleratesUnavailableNumericProperties(t *testing.T) {
	socketPath := filepath.Join(os.TempDir(), "txtamp-test-unavailable.sock")
	os.Remove(socketPath)
	t.Cleanup(func() {
		os.Remove(socketPath)
	})

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("expected listener, got %v", err)
	}
	defer listener.Close()

	go serveUnavailableDuration(t, listener)

	p := &Player{socketPath: socketPath}
	startFakeProcess(t, p)

	status, err := p.Status()
	if err != nil {
		t.Fatalf("expected status to tolerate unavailable property, got %v", err)
	}

	if status.Elapsed != 12 {
		t.Fatalf("expected elapsed 12, got %d", status.Elapsed)
	}
	if status.Duration != 0 {
		t.Fatalf("expected unknown duration 0, got %d", status.Duration)
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

func serveUnavailableDuration(t *testing.T, listener net.Listener) {
	t.Helper()

	responses := map[string]map[string]any{
		"time-pos": {"data": 12.0, "error": "success"},
		"duration": {"error": "property unavailable"},
		"pause":    {"data": false, "error": "success"},
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
		if err := json.NewEncoder(conn).Encode(responses[property]); err != nil {
			conn.Close()
			return
		}

		conn.Close()
	}
}

func startFakeProcess(t *testing.T, p *Player) {
	t.Helper()

	cmd := exec.Command("sleep", "5")
	if err := cmd.Start(); err != nil {
		t.Fatalf("expected fake process to start, got %v", err)
	}
	t.Cleanup(func() {
		cmd.Process.Kill()
		cmd.Wait()
	})

	p.cmd = cmd
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
