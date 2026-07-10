package mainview

import (
	"context"
	"time"
	"txtamp/navidrome"
	"txtamp/ui"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

const loadTimeout = 10 * time.Second

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

	focused          focusPane
	selectedPlaylist int
	selectedSong     int
	currentSong      *navidrome.Song
	paused           bool
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

func New(connectedTo string, client navidrome.Client) Model {
	return Model{
		connectedTo: connectedTo,
		client:      client,
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "left":
			m.focused = playlistsPane
		case "right":
			m.focused = songsPane
		case "up":
			cmd := m.moveSelection(-1)
			return m, cmd
		case "down":
			cmd := m.moveSelection(1)
			return m, cmd
		case "enter":
			cmd := m.activateSelection()
			return m, cmd
		case " ", "space":
			m.togglePlayPause()
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

	song := m.songs[m.selectedSong]
	m.currentSong = &song
	m.paused = false

	return nil
}

func (m *Model) togglePlayPause() {
	if m.currentSong == nil {
		m.activateSelection()
		return
	}

	m.paused = !m.paused
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

func clamp(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}

	return value
}
