package player

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	ipcTimeout     = 500 * time.Millisecond
	ipcReadyWait   = 2 * time.Second
	ipcRetryPeriod = 50 * time.Millisecond
)

type Player struct {
	cmd        *exec.Cmd
	socketPath string
}

func New() *Player {
	return &Player{
		socketPath: filepath.Join(os.TempDir(), fmt.Sprintf("txtamp-mpv-%d.sock", os.Getpid())),
	}
}

func (p *Player) Play(url string) error {
	if err := p.Stop(); err != nil {
		return err
	}

	if err := os.Remove(p.socketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing mpv socket: %w", err)
	}

	cmd := exec.Command(
		"mpv",
		"--no-video",
		"--idle=no",
		"--force-window=no",
		"--input-ipc-server="+p.socketPath,
		url,
	)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting mpv: %w", err)
	}

	p.cmd = cmd

	go func() {
		cmd.Wait()
		os.Remove(p.socketPath)
	}()

	return nil
}

func (p *Player) TogglePause() error {
	return p.send("cycle", "pause")
}

func (p *Player) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	if err := p.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return fmt.Errorf("stopping mpv: %w", err)
	}

	p.cmd = nil

	return nil
}

func (p *Player) send(command ...any) error {
	conn, err := p.dial()
	if err != nil {
		return fmt.Errorf("connecting to mpv: %w", err)
	}
	defer conn.Close()

	message := struct {
		Command []any `json:"command"`
	}{
		Command: command,
	}

	if err := json.NewEncoder(conn).Encode(message); err != nil {
		return fmt.Errorf("sending mpv command: %w", err)
	}

	return nil
}

func (p *Player) dial() (net.Conn, error) {
	deadline := time.Now().Add(ipcReadyWait)

	for {
		conn, err := net.DialTimeout("unix", p.socketPath, ipcTimeout)
		if err == nil {
			return conn, nil
		}

		if time.Now().After(deadline) {
			return nil, err
		}

		time.Sleep(ipcRetryPeriod)
	}
}
