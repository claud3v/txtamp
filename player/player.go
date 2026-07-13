package player

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const (
	ipcTimeout     = 500 * time.Millisecond
	ipcReadyWait   = 2 * time.Second
	ipcRetryPeriod = 50 * time.Millisecond
)

var ErrNotRunning = errors.New("mpv is not running")

type Player struct {
	cmd        *exec.Cmd
	socketPath string
	mu         sync.Mutex
}

type Status struct {
	Elapsed  int
	Duration int
	Paused   bool
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

	p.mu.Lock()
	p.cmd = cmd
	p.mu.Unlock()

	go func() {
		cmd.Wait()
		os.Remove(p.socketPath)
		p.mu.Lock()
		if p.cmd == cmd {
			p.cmd = nil
		}
		p.mu.Unlock()
	}()

	return nil
}

func (p *Player) TogglePause() error {
	response, err := p.send("cycle", "pause")
	if err != nil {
		return err
	}

	return checkResponse(response)
}

func (p *Player) SeekStart() error {
	response, err := p.send("seek", 0, "absolute")
	if err != nil {
		return err
	}

	return checkResponse(response)
}

func (p *Player) Seek(seconds int) error {
	response, err := p.send("seek", seconds, "relative")
	if err != nil {
		return err
	}

	return checkResponse(response)
}

func (p *Player) AdjustVolume(delta int) error {
	response, err := p.send("add", "volume", delta)
	if err != nil {
		return err
	}

	return checkResponse(response)
}

func (p *Player) Status() (Status, error) {
	elapsed, err := p.getNumberProperty("time-pos")
	if err != nil {
		return Status{}, err
	}

	duration, err := p.getNumberProperty("duration")
	if err != nil {
		return Status{}, err
	}

	paused, err := p.getBoolProperty("pause")
	if err != nil {
		return Status{}, err
	}

	return Status{
		Elapsed:  int(elapsed),
		Duration: int(duration),
		Paused:   paused,
	}, nil
}

func (p *Player) Stop() error {
	p.mu.Lock()
	cmd := p.cmd
	p.cmd = nil
	p.mu.Unlock()

	if cmd == nil || cmd.Process == nil {
		return nil
	}

	if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return fmt.Errorf("stopping mpv: %w", err)
	}

	return nil
}

func (p *Player) send(command ...any) (mpvResponse, error) {
	if !p.running() {
		return mpvResponse{}, ErrNotRunning
	}

	conn, err := p.dial()
	if err != nil {
		if !p.running() {
			return mpvResponse{}, ErrNotRunning
		}

		return mpvResponse{}, fmt.Errorf("connecting to mpv: %w", err)
	}
	defer conn.Close()

	message := struct {
		Command []any `json:"command"`
	}{
		Command: command,
	}

	if err := json.NewEncoder(conn).Encode(message); err != nil {
		return mpvResponse{}, fmt.Errorf("sending mpv command: %w", err)
	}

	var response mpvResponse
	if err := json.NewDecoder(bufio.NewReader(conn)).Decode(&response); err != nil {
		return mpvResponse{}, fmt.Errorf("reading mpv response: %w", err)
	}

	return response, nil
}

func (p *Player) running() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.cmd != nil && p.cmd.Process != nil
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

func (p *Player) getNumberProperty(name string) (float64, error) {
	response, err := p.send("get_property", name)
	if err != nil {
		return 0, err
	}
	if response.Error == "property unavailable" {
		return 0, nil
	}
	if err := checkResponse(response); err != nil {
		return 0, err
	}

	value, ok := response.Data.(float64)
	if !ok {
		return 0, fmt.Errorf("mpv property %s was not a number", name)
	}

	return value, nil
}

func (p *Player) getBoolProperty(name string) (bool, error) {
	response, err := p.send("get_property", name)
	if err != nil {
		return false, err
	}
	if response.Error == "property unavailable" {
		return false, nil
	}
	if err := checkResponse(response); err != nil {
		return false, err
	}

	value, ok := response.Data.(bool)
	if !ok {
		return false, fmt.Errorf("mpv property %s was not a boolean", name)
	}

	return value, nil
}

type mpvResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

func checkResponse(response mpvResponse) error {
	if response.Error == "success" {
		return nil
	}

	return fmt.Errorf("mpv command failed: %s", response.Error)
}
