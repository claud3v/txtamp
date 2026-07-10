package mainview

import (
	"fmt"
	"strings"
	"txtamp/ui"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

type focusPane int

const (
	playlistsPane focusPane = iota
	songsPane
)

type playlist struct {
	name  string
	songs []song
}

type song struct {
	title    string
	artist   string
	duration string
}

type Model struct {
	width       int
	height      int
	connectedTo string

	focused          focusPane
	selectedPlaylist int
	selectedSong     int
	currentSong      *song
	paused           bool

	playlists []playlist
}

func New(connectedTo string) Model {
	return Model{
		connectedTo: connectedTo,
		focused:     playlistsPane,
		playlists:   dummyPlaylists(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "left":
			m.focused = playlistsPane
		case "right":
			m.focused = songsPane
		case "up":
			m.moveSelection(-1)
		case "down":
			m.moveSelection(1)
		case "enter":
			m.activateSelection()
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

func (m Model) renderPlaylists(width, height int) string {
	lines := []string{
		ui.Title.Render("TxtAmp"),
		ui.Subtitle.Render("Connected: " + m.connectedTo),
		"",
		paneTitle("Playlists", m.focused == playlistsPane),
	}

	for i, playlist := range m.playlists {
		line := selectableLine(playlist.name, i == m.selectedPlaylist, m.focused == playlistsPane, width-4)
		lines = append(lines, line)
	}

	return ui.Sidebar.
		Width(width).
		Height(height).
		Render(strings.Join(lines, "\n"))
}

func (m Model) renderSongs(width, height int) string {
	currentPlaylist := m.playlists[m.selectedPlaylist]

	lines := []string{
		paneTitle(currentPlaylist.name, m.focused == songsPane),
		ui.Subtitle.Render(fmt.Sprintf("%d songs", len(currentPlaylist.songs))),
		"",
	}

	for i, song := range currentPlaylist.songs {
		titleWidth := max(width-18, 10)
		title := ui.Truncate(song.title, titleWidth)
		line := fmt.Sprintf("%-*"+"s %5s", titleWidth, title, song.duration)
		line = selectableLine(line, i == m.selectedSong, m.focused == songsPane, width-4)
		lines = append(lines, line)
	}

	return ui.MainPane.
		Width(width).
		Height(height).
		Render(strings.Join(lines, "\n"))
}

func (m Model) renderPlayer(width int) string {
	status := "Stopped"
	nowPlaying := "No song selected"
	progress := "00:00 / 00:00"
	bars := "[      ]"

	if m.currentSong != nil {
		if m.paused {
			status = "Paused"
			bars = "[||    ]"
		} else {
			status = "Playing"
			bars = "[||||  ]"
		}

		nowPlaying = fmt.Sprintf("%s - %s", m.currentSong.artist, m.currentSong.title)
		progress = "00:00 / " + m.currentSong.duration
	}

	line := fmt.Sprintf("%s  %s  %s", status, bars, ui.Truncate(nowPlaying, max(width-34, 10)))
	if width > 30 {
		line = fmt.Sprintf("%-*s %s", max(width-14, 10), line, progress)
	}

	return ui.PlayerBar.
		Width(width - 2).
		Render(line)
}

func (m *Model) moveSelection(delta int) {
	switch m.focused {
	case playlistsPane:
		m.selectedPlaylist = clamp(m.selectedPlaylist+delta, 0, len(m.playlists)-1)
		m.selectedSong = 0
	case songsPane:
		songs := m.playlists[m.selectedPlaylist].songs
		m.selectedSong = clamp(m.selectedSong+delta, 0, len(songs)-1)
	}
}

func (m *Model) activateSelection() {
	if m.focused == playlistsPane {
		m.focused = songsPane
		m.selectedSong = 0
		return
	}

	song := m.playlists[m.selectedPlaylist].songs[m.selectedSong]
	m.currentSong = &song
	m.paused = false
}

func (m *Model) togglePlayPause() {
	if m.currentSong == nil {
		m.activateSelection()
		return
	}

	m.paused = !m.paused
}

func paneTitle(title string, focused bool) string {
	if focused {
		return ui.PaneTitleFocused.Render(title)
	}

	return ui.PaneTitle.Render(title)
}

func selectableLine(text string, selected, focused bool, width int) string {
	prefix := "  "
	if selected {
		prefix = "> "
	}

	line := prefix + ui.Truncate(text, max(width-2, 1))
	if selected && focused {
		return ui.SelectedRowFocused.Width(width).Render(line)
	}
	if selected {
		return ui.SelectedRow.Width(width).Render(line)
	}

	return lipgloss.NewStyle().Width(width).Render(line)
}

func dummyPlaylists() []playlist {
	return []playlist{
		{
			name: "Aerosmith Top Songs",
			songs: []song{
				{title: "Dream On", artist: "Aerosmith", duration: "4:28"},
				{title: "Sweet Emotion", artist: "Aerosmith", duration: "4:34"},
				{title: "Walk This Way", artist: "Aerosmith", duration: "3:40"},
				{title: "Crazy", artist: "Aerosmith", duration: "5:17"},
				{title: "Janie's Got a Gun", artist: "Aerosmith", duration: "5:30"},
			},
		},
		{
			name: "Late Night Drive",
			songs: []song{
				{title: "Midnight City", artist: "M83", duration: "4:04"},
				{title: "Nightcall", artist: "Kavinsky", duration: "4:18"},
				{title: "A Real Hero", artist: "College", duration: "4:27"},
				{title: "Under Your Spell", artist: "Desire", duration: "3:52"},
			},
		},
		{
			name: "Sunday Reset",
			songs: []song{
				{title: "Harvest Moon", artist: "Neil Young", duration: "5:03"},
				{title: "Pink Moon", artist: "Nick Drake", duration: "2:06"},
				{title: "Landslide", artist: "Fleetwood Mac", duration: "3:19"},
				{title: "Into the Mystic", artist: "Van Morrison", duration: "3:25"},
			},
		},
		{
			name: "Tiny Desk Energy",
			songs: []song{
				{title: "Them Changes", artist: "Thundercat", duration: "3:08"},
				{title: "Come Down", artist: "Anderson .Paak", duration: "2:56"},
				{title: "Dang!", artist: "Mac Miller", duration: "5:05"},
				{title: "Am I Wrong", artist: "Anderson .Paak", duration: "4:13"},
			},
		},
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
