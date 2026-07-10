package mainview

import (
	"context"
	"errors"
	"time"
	"txtamp/navidrome"
	"txtamp/player"
	"txtamp/ui"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

const (
	loadTimeout  = 10 * time.Second
	statusPeriod = 1 * time.Second
	restartLimit = 3
)

type focusPane int

const (
	playlistsPane focusPane = iota
	songsPane
)

type Model struct {
	width       int
	height      int
	connectedTo string
	client      navidrome.Client
	player      *player.Player

	focused          focusPane
	selectedPlaylist int
	selectedSong     int
	currentSong      *navidrome.Song
	paused           bool
	elapsed          int
	duration         int
	currentSongIndex int
	playbackID       int
	loading          bool
	err              error

	playlists []navidrome.Playlist
	songs     []navidrome.Song
}

type playlistsLoadedMsg struct {
	playlists []navidrome.Playlist
	err       error
}

type songsLoadedMsg struct {
	songs []navidrome.Song
	err   error
}

type playbackMsg struct {
	song       *navidrome.Song
	songIndex  int
	playbackID int
	paused     bool
	err        error
}

type seekMsg struct {
	err error
}

type playerStatusMsg struct {
	playbackID int
	status     player.Status
	err        error
}

type playerTickMsg struct {
	playbackID int
}

func New(connectedTo string, client navidrome.Client) Model {
	return Model{
		connectedTo: connectedTo,
		client:      client,
		player:      player.New(),
		focused:     playlistsPane,
		loading:     true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadPlaylists()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case playlistsLoadedMsg:
		m.loading = false
		m.err = msg.err
		m.playlists = msg.playlists
		m.songs = nil
		m.selectedPlaylist = 0
		m.selectedSong = 0

		if msg.err != nil || len(m.playlists) == 0 {
			return m, nil
		}

		cmd := m.loadSelectedPlaylist()
		return m, cmd
	case songsLoadedMsg:
		m.loading = false
		m.err = msg.err
		m.songs = msg.songs
		m.selectedSong = 0
	case playbackMsg:
		m.err = msg.err
		if msg.err != nil {
			return m, nil
		}
		if msg.song != nil {
			m.currentSong = msg.song
			m.currentSongIndex = msg.songIndex
			m.playbackID = msg.playbackID
			m.elapsed = 0
			m.duration = msg.song.Duration
		}
		m.paused = msg.paused
		return m, tea.Batch(m.pollPlayerStatus(), tickPlayerStatus(m.playbackID))
	case seekMsg:
		m.err = msg.err
		if msg.err == nil {
			m.elapsed = 0
		}
	case playerTickMsg:
		if m.currentSong == nil || msg.playbackID != m.playbackID {
			return m, nil
		}

		return m, m.pollPlayerStatus()
	case playerStatusMsg:
		if msg.playbackID != m.playbackID {
			return m, nil
		}

		if msg.err != nil {
			if errors.Is(msg.err, player.ErrNotRunning) {
				cmd := m.playNextSong()
				return m, cmd
			}

			m.err = msg.err
			return m, tickPlayerStatus(m.playbackID)
		}

		m.elapsed = msg.status.Elapsed
		m.duration = msg.status.Duration
		m.paused = msg.status.Paused
		if m.playbackFinished() {
			cmd := m.playNextSong()
			return m, cmd
		}
		if m.currentSong != nil {
			return m, tickPlayerStatus(m.playbackID)
		}
	case tea.KeyMsg:
		action, ok := actionForKey(msg.String())
		if !ok {
			return m, nil
		}

		switch action {
		case actionQuit:
			m.player.Stop()
			return m, tea.Quit
		case actionFocusSidebar:
			m.focused = playlistsPane
		case actionFocusMainArea:
			m.focused = songsPane
		case actionMoveUp:
			cmd := m.moveSelection(-1)
			return m, cmd
		case actionMoveDown:
			cmd := m.moveSelection(1)
			return m, cmd
		case actionActivate:
			cmd := m.activateSelection()
			return m, cmd
		case actionPlayPause:
			cmd := m.togglePlayPause()
			return m, cmd
		case actionNextSong:
			cmd := m.playNextSong()
			return m, cmd
		case actionPreviousSong:
			cmd := m.playPreviousSong()
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	layout := ui.NewShellLayout(m.width, m.height)

	sidebar := m.renderPlaylists(layout.SidebarWidth, layout.BodyHeight)
	songs := m.renderSongs(layout.MainWidth, layout.BodyHeight)
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, songs)
	player := m.renderPlayer(layout.Width)

	content := lipgloss.JoinVertical(lipgloss.Left, body, player)

	view := tea.NewView(content)
	view.AltScreen = true

	return view
}

func (m *Model) moveSelection(delta int) tea.Cmd {
	switch m.focused {
	case playlistsPane:
		if len(m.playlists) == 0 {
			return nil
		}

		previous := m.selectedPlaylist
		m.selectedPlaylist = clamp(m.selectedPlaylist+delta, 0, len(m.playlists)-1)
		m.selectedSong = 0
		if m.selectedPlaylist != previous {
			return m.loadSelectedPlaylist()
		}
	case songsPane:
		if len(m.songs) == 0 {
			return nil
		}

		m.selectedSong = clamp(m.selectedSong+delta, 0, len(m.songs)-1)
	}

	return nil
}

func (m *Model) activateSelection() tea.Cmd {
	if m.focused == playlistsPane {
		m.focused = songsPane
		m.selectedSong = 0
		return m.loadSelectedPlaylist()
	}

	if len(m.songs) == 0 {
		return nil
	}

	return m.playSongAt(m.selectedSong)
}

func (m *Model) togglePlayPause() tea.Cmd {
	if m.currentSong == nil {
		return m.activateSelection()
	}

	paused := !m.paused
	player := m.player

	return func() tea.Msg {
		return playbackMsg{
			paused: paused,
			err:    player.TogglePause(),
		}
	}
}

func (m *Model) loadPlaylists() tea.Cmd {
	m.loading = true

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), loadTimeout)
		defer cancel()

		playlists, err := m.client.ListPlaylists(ctx)
		return playlistsLoadedMsg{playlists: playlists, err: err}
	}
}

func (m *Model) loadSelectedPlaylist() tea.Cmd {
	if len(m.playlists) == 0 {
		return nil
	}

	playlistID := m.playlists[m.selectedPlaylist].ID
	m.loading = true
	m.songs = nil

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), loadTimeout)
		defer cancel()

		songs, err := m.client.GetPlaylist(ctx, playlistID)
		return songsLoadedMsg{songs: songs, err: err}
	}
}

func (m *Model) playSong(song navidrome.Song) tea.Cmd {
	return m.playSongAtIndex(song, m.selectedSong)
}

func (m *Model) playSongAt(index int) tea.Cmd {
	if len(m.songs) == 0 {
		return nil
	}

	index = clamp(index, 0, len(m.songs)-1)
	m.selectedSong = index

	return m.playSongAtIndex(m.songs[index], index)
}

func (m *Model) playSongAtIndex(song navidrome.Song, index int) tea.Cmd {
	client := m.client
	player := m.player
	playbackID := m.playbackID + 1
	m.playbackID = playbackID

	return func() tea.Msg {
		streamURL, err := client.StreamURL(song.ID)
		if err != nil {
			return playbackMsg{err: err}
		}

		if err := player.Play(streamURL); err != nil {
			return playbackMsg{err: err}
		}

		return playbackMsg{
			song:       &song,
			songIndex:  index,
			playbackID: playbackID,
			paused:     false,
		}
	}
}

func (m *Model) playNextSong() tea.Cmd {
	if len(m.songs) == 0 {
		return nil
	}

	nextIndex := m.currentSongIndex + 1
	if nextIndex >= len(m.songs) {
		m.currentSong = nil
		m.elapsed = 0
		m.duration = 0
		m.paused = false
		m.player.Stop()
		return nil
	}

	return m.playSongAt(nextIndex)
}

func (m *Model) playPreviousSong() tea.Cmd {
	if m.currentSong == nil {
		return nil
	}

	if m.elapsed > restartLimit || m.currentSongIndex == 0 {
		return m.seekStart()
	}

	return m.playSongAt(m.currentSongIndex - 1)
}

func (m *Model) seekStart() tea.Cmd {
	player := m.player

	return func() tea.Msg {
		return seekMsg{err: player.SeekStart()}
	}
}

func (m Model) playbackFinished() bool {
	return m.duration > 0 && m.elapsed >= m.duration-1 && !m.paused
}

func (m *Model) pollPlayerStatus() tea.Cmd {
	player := m.player
	playbackID := m.playbackID

	return func() tea.Msg {
		status, err := player.Status()
		return playerStatusMsg{playbackID: playbackID, status: status, err: err}
	}
}

func tickPlayerStatus(playbackID int) tea.Cmd {
	return tea.Tick(statusPeriod, func(t time.Time) tea.Msg {
		return playerTickMsg{playbackID: playbackID}
	})
}

func clamp(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}

	return value
}
