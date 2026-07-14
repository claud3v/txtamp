package mainview

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"txtamp/config"
	"txtamp/navidrome"
	"txtamp/player"
	"txtamp/ui"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

const (
	loadTimeout       = 10 * time.Second
	statusPeriod      = 1 * time.Second
	toastPeriod       = 2 * time.Second
	marqueePeriod     = 250 * time.Millisecond
	marqueePauseTicks = 3
	restartLimit      = 3
	seekStep          = 10
	volumeStep        = 5
)

type focusPane int

const (
	modeSelectorPane focusPane = iota
	playlistsPane
	songsPane
)

type sidebarMode int

const (
	artistsMode sidebarMode = iota
	playlistsMode
)

type mainContentMode int

const (
	libraryContent mainContentMode = iota
	globalSearchContent
	queueContent
)

type albumGroup struct {
	album navidrome.Album
	songs []navidrome.Song
}

type Model struct {
	width       int
	height      int
	connectedTo string
	client      navidrome.Client
	player      *player.Player

	focused                    focusPane
	mode                       sidebarMode
	contentMode                mainContentMode
	modeDialogOpen             bool
	helpOpen                   bool
	selectedMode               sidebarMode
	searching                  bool
	searchPane                 focusPane
	searchQuery                string
	goToArtistPending          bool
	globalSearching            bool
	globalSearchQuery          string
	globalSearchSubmittedQuery string
	globalSearchLoading        bool
	globalSearchErr            error
	globalSearchResult         navidrome.SearchResult
	selectedSearchResult       int
	selectedQueue              int
	loadedPlaylistID           string
	loadedArtistID             string
	selectedPlaylist           int
	selectedArtist             int
	selectedArtistRow          int
	selectedSong               int
	collapsedAlbums            map[int]bool
	currentSong                *navidrome.Song
	paused                     bool
	elapsed                    int
	duration                   int
	currentSongIndex           int
	playbackSongs              []navidrome.Song
	playbackSource             string
	playbackID                 int
	loading                    bool
	err                        error
	toast                      string
	toastID                    int
	sidebarMarqueeOffset       int
	queueDirty                 bool

	playlists []navidrome.Playlist
	artists   []navidrome.Artist
	albums    []albumGroup
	songs     []navidrome.Song
	queue     []navidrome.Song
}

type playlistsLoadedMsg struct {
	playlists []navidrome.Playlist
	err       error
}

type songsLoadedMsg struct {
	playlistID string
	songs      []navidrome.Song
	err        error
}

type artistsLoadedMsg struct {
	artists []navidrome.Artist
	err     error
}

type artistLoadedMsg struct {
	artistID string
	albums   []albumGroup
	songs    []navidrome.Song
	err      error
}

type globalSearchLoadedMsg struct {
	query  string
	result navidrome.SearchResult
	err    error
}

type playbackMsg struct {
	song       *navidrome.Song
	songIndex  int
	playbackID int
	paused     bool
	err        error
}

type seekMsg struct {
	elapsed int
	err     error
}

type stopMsg struct {
	err error
}

type volumeMsg struct {
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

type clearToastMsg struct {
	toastID int
}

type queueLoadedMsg struct {
	songs []navidrome.Song
	found bool
	err   error
}

type queueSavedMsg struct {
	err error
}

type sidebarMarqueeTickMsg struct{}

func New(connectedTo string, client navidrome.Client) Model {
	return Model{
		connectedTo: connectedTo,
		client:      client,
		player:      player.New(),
		focused:     playlistsPane,
		mode:        playlistsMode,
		loading:     true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadPlaylists(), loadSavedQueue(), tickSidebarMarquee())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case playlistsLoadedMsg:
		m.playlists = msg.playlists
		if m.mode != playlistsMode {
			return m, nil
		}

		m.loading = false
		m.err = msg.err
		m.songs = nil
		m.selectedPlaylist = 0
		m.selectedSong = 0

		if msg.err != nil || len(m.playlists) == 0 {
			return m, nil
		}

		cmd := m.loadSelectedPlaylist()
		return m, cmd
	case songsLoadedMsg:
		if m.mode != playlistsMode {
			return m, nil
		}
		if len(m.playlists) > 0 && msg.playlistID != m.playlists[m.selectedPlaylist].ID {
			return m, nil
		}

		m.loading = false
		m.err = msg.err
		m.loadedPlaylistID = msg.playlistID
		m.albums = nil
		m.songs = msg.songs
		m.selectedSong = 0
	case artistsLoadedMsg:
		m.artists = msg.artists
		if m.mode != artistsMode {
			return m, nil
		}

		m.loading = false
		m.err = msg.err
		m.albums = nil
		m.songs = nil
		m.selectedArtist = 0
		m.selectedArtistRow = 0
		m.collapsedAlbums = nil
		m.selectedSong = 0

		if msg.err != nil || len(m.artists) == 0 {
			return m, nil
		}

		cmd := m.loadSelectedArtist()
		return m, cmd
	case artistLoadedMsg:
		if m.mode != artistsMode {
			return m, nil
		}
		if len(m.artists) > 0 && msg.artistID != m.artists[m.selectedArtist].ID {
			return m, nil
		}

		m.loading = false
		m.err = msg.err
		m.loadedArtistID = msg.artistID
		m.albums = msg.albums
		m.songs = msg.songs
		m.selectedArtistRow = 0
		m.collapsedAlbums = nil
		m.selectedSong = 0
	case globalSearchLoadedMsg:
		if msg.query != m.globalSearchSubmittedQuery {
			return m, nil
		}

		m.globalSearchLoading = false
		m.globalSearchErr = msg.err
		m.globalSearchResult = msg.result
		m.selectedSearchResult = 0
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
			m.elapsed = msg.elapsed
		}
	case stopMsg:
		m.err = msg.err
		if msg.err == nil {
			m.currentSong = nil
			m.elapsed = 0
			m.duration = 0
			m.paused = false
		}
	case volumeMsg:
		m.err = msg.err
	case playerTickMsg:
		if m.currentSong == nil || msg.playbackID != m.playbackID {
			return m, nil
		}

		return m, m.pollPlayerStatus()
	case clearToastMsg:
		if msg.toastID == m.toastID {
			m.toast = ""
		}
	case queueLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		if msg.found && !m.queueDirty {
			m.queue = msg.songs
			m.selectedQueue = clamp(m.selectedQueue, 0, max(len(m.queue)-1, 0))
		}
	case queueSavedMsg:
		if msg.err != nil {
			m.err = msg.err
		}
	case sidebarMarqueeTickMsg:
		if m.sidebarMarqueeActive() {
			m.sidebarMarqueeOffset++
			return m, tickSidebarMarquee()
		}
		m.sidebarMarqueeOffset = 0
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
		if action, ok := actionForKey(msg.String()); ok && action == actionToggleHelp {
			m.helpOpen = !m.helpOpen
			return m, nil
		}

		if m.helpOpen {
			if action, ok := actionForKey(msg.String()); ok {
				if action == actionQuit {
					m.player.Stop()
					return m, tea.Quit
				}
				if action == actionCloseDialog {
					m.helpOpen = false
				}
			}
			return m, nil
		}

		if m.goToArtistPending {
			cmd := m.handleGoToArtistKey(msg)
			return m, cmd
		}

		if m.globalSearching {
			cmd := m.handleGlobalSearchKey(msg)
			return m, cmd
		}

		if m.searching {
			cmd := m.handleSearchKey(msg)
			return m, cmd
		}

		action, ok := actionForKey(msg.String())
		if !ok {
			return m, nil
		}

		if action == actionQuit {
			m.player.Stop()
			return m, tea.Quit
		}

		if m.modeDialogOpen {
			cmd := m.handleModeDialogAction(action)
			return m, cmd
		}

		switch action {
		case actionFocusSidebar:
			m.focused = playlistsPane
			m.sidebarMarqueeOffset = 0
			return m, tickSidebarMarquee()
		case actionFocusMainArea:
			m.focused = songsPane
			m.sidebarMarqueeOffset = 0
		case actionCloseDialog:
			if m.contentMode == globalSearchContent {
				m.contentMode = libraryContent
				m.globalSearching = false
			} else if m.contentMode == queueContent {
				m.contentMode = libraryContent
			} else {
				m.clearSearch()
			}
		case actionStartSearch:
			m.startSearch()
		case actionGlobalSearch:
			m.startGlobalSearch()
		case actionGoToArtist:
			m.startGoToArtist()
		case actionToggleQueue:
			m.toggleQueue()
		case actionAddToQueue:
			if m.addSelectedSongToQueue() {
				return m, tea.Batch(clearToast(m.toastID), m.saveQueue())
			}
		case actionRemoveQueue:
			if m.removeSelectedQueueSong() {
				return m, m.saveQueue()
			}
		case actionClearQueue:
			if m.clearQueue() {
				return m, tea.Batch(clearToast(m.toastID), m.saveQueue())
			}
		case actionPlayQueue:
			cmd := m.playQueueFromTop()
			if cmd != nil {
				return m, cmd
			}
		case actionQueueUp:
			if m.moveQueuedSong(-1) {
				return m, m.saveQueue()
			}
		case actionQueueDown:
			if m.moveQueuedSong(1) {
				return m, m.saveQueue()
			}
		case actionExpandAlbums:
			m.expandAllAlbums()
		case actionCollapseAlbums:
			m.collapseAllAlbums()
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
			if cmd := m.playSelectedAlbum(); cmd != nil {
				return m, cmd
			}
			cmd := m.togglePlayPause()
			return m, cmd
		case actionStopPlayback:
			cmd := m.stopPlayback()
			return m, withToastClear(cmd, m.toastID)
		case actionNextSong:
			cmd := m.playNextSong()
			return m, cmd
		case actionPreviousSong:
			cmd := m.playPreviousSong()
			return m, cmd
		case actionSeekBackward:
			cmd := m.seekRelative(-seekStep)
			return m, withToastClear(cmd, m.toastID)
		case actionSeekForward:
			cmd := m.seekRelative(seekStep)
			return m, withToastClear(cmd, m.toastID)
		case actionVolumeDown:
			cmd := m.adjustVolume(-volumeStep)
			return m, withToastClear(cmd, m.toastID)
		case actionVolumeUp:
			cmd := m.adjustVolume(volumeStep)
			return m, withToastClear(cmd, m.toastID)
		case actionShowArtists:
			cmd := m.switchMode(artistsMode)
			return m, cmd
		case actionShowPlaylists:
			cmd := m.switchMode(playlistsMode)
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	layout := ui.NewShellLayout(m.width, m.height)

	sidebar := m.renderSidebar(layout.SidebarWidth, layout.BodyHeight)
	mainArea := m.renderMainArea(layout.MainWidth, layout.BodyHeight)
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainArea)
	player := m.renderPlayer(layout.Width)
	status := m.renderStatusBar(layout.Width)

	content := lipgloss.JoinVertical(lipgloss.Left, body, player, status)
	if m.modeDialogOpen {
		content = overlayCentered(content, m.renderModeDialog(), layout.Width, layout.Height)
	}
	if m.helpOpen {
		content = overlayCentered(content, m.renderHelpDialog(), layout.Width, layout.Height)
	}
	if m.toast != "" {
		content = overlayBottomRight(content, m.renderToast(), layout.Width, layout.Height)
	}

	view := tea.NewView(content)
	view.AltScreen = true

	return view
}

func overlayCentered(content, overlay string, width, height int) string {
	overlayWidth := lipgloss.Width(overlay)
	overlayHeight := lipgloss.Height(overlay)
	left := max((width-overlayWidth)/2, 0)
	top := max((height-overlayHeight)/2, 0)

	return overlayAt(content, overlay, width, height, left, top)
}

func overlayBottomRight(content, overlay string, width, height int) string {
	overlayWidth := lipgloss.Width(overlay)
	overlayHeight := lipgloss.Height(overlay)
	left := max(width-overlayWidth-2, 0)
	top := max(height-overlayHeight-4, 0)

	return overlayAt(content, overlay, width, height, left, top)
}

func overlayAt(content, overlay string, width, height, left, top int) string {
	overlayWidth := lipgloss.Width(overlay)
	contentLines := strings.Split(content, "\n")
	overlayLines := strings.Split(overlay, "\n")
	for len(contentLines) < height {
		contentLines = append(contentLines, "")
	}

	for i, overlayLine := range overlayLines {
		target := top + i
		if target < 0 || target >= len(contentLines) {
			continue
		}

		line := contentLines[target]
		prefix := ansi.Cut(line, 0, left)
		prefix += strings.Repeat(" ", max(left-ansi.StringWidth(prefix), 0))

		rightStart := left + overlayWidth
		suffix := ""
		if ansi.StringWidth(line) > rightStart {
			suffix = ansi.Cut(line, rightStart, width)
		}

		contentLines[target] = prefix + overlayLine + suffix
	}

	return strings.Join(contentLines, "\n")
}

func clearToast(toastID int) tea.Cmd {
	return tea.Tick(toastPeriod, func(t time.Time) tea.Msg {
		return clearToastMsg{toastID: toastID}
	})
}

func withToastClear(cmd tea.Cmd, toastID int) tea.Cmd {
	if cmd == nil {
		return nil
	}

	return tea.Batch(cmd, clearToast(toastID))
}

func loadSavedQueue() tea.Cmd {
	return func() tea.Msg {
		songs, found, err := config.LoadQueue()
		return queueLoadedMsg{songs: songs, found: found, err: err}
	}
}

func (m Model) renderModeDialog() string {
	width := 22
	rows := []string{
		modeDialogRow("Artists", "1", m.selectedMode == artistsMode, width),
		modeDialogRow("Playlists", "2", m.selectedMode == playlistsMode, width),
	}

	return ui.Dialog.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func modeDialogRow(label, key string, selected bool, width int) string {
	text := lipgloss.JoinHorizontal(
		lipgloss.Top,
		label,
		lipgloss.PlaceHorizontal(max(width-lipgloss.Width(label), 1), lipgloss.Right, key),
	)

	if selected {
		return ui.SelectedRowFocused.Width(width).Render(text)
	}

	return lipgloss.NewStyle().Width(width).Render(text)
}

func (m *Model) moveSelection(delta int) tea.Cmd {
	if m.contentMode == globalSearchContent && m.focused == songsPane {
		m.moveGlobalSearchSelection(delta)
		return nil
	}
	if m.contentMode == queueContent && m.focused == songsPane {
		m.moveQueueSelection(delta)
		return nil
	}

	switch m.focused {
	case modeSelectorPane:
		if delta > 0 {
			m.focused = playlistsPane
			m.sidebarMarqueeOffset = 0
			return tickSidebarMarquee()
		}
	case playlistsPane:
		if delta < 0 && m.selectedSidebarPosition() == 0 {
			m.focused = modeSelectorPane
			m.sidebarMarqueeOffset = 0
			return nil
		}

		return m.moveSidebarSelection(delta)
	case songsPane:
		if m.mode == artistsMode {
			m.moveArtistRowSelection(delta)
			return nil
		}

		songs := m.filteredSongs()
		if len(songs) == 0 {
			return nil
		}

		position := m.selectedSongPosition(songs)
		position = clamp(position+delta, 0, len(songs)-1)
		m.selectedSong = songs[position].index
	}

	return nil
}

func (m *Model) activateSelection() tea.Cmd {
	if m.focused == modeSelectorPane {
		m.openModeDialog()
		return nil
	}

	if m.contentMode == globalSearchContent && m.focused == songsPane {
		return m.activateGlobalSearchResult()
	}
	if m.contentMode == queueContent && m.focused == songsPane {
		return m.playSelectedQueueSong()
	}

	if m.focused == playlistsPane {
		m.focused = songsPane
		m.sidebarMarqueeOffset = 0
		m.selectedSong = 0
		if m.selectedSidebarItemLoaded() {
			return nil
		}

		return m.loadSelectedSidebarItem()
	}

	if m.mode == artistsMode {
		return m.activateArtistRow()
	}

	songs := m.filteredSongs()
	if len(songs) == 0 {
		return nil
	}

	return m.playSongAt(m.selectedSong)
}

func (m *Model) openModeDialog() {
	m.modeDialogOpen = true
	m.selectedMode = m.mode
}

func (m *Model) handleModeDialogAction(action action) tea.Cmd {
	switch action {
	case actionCloseDialog:
		m.modeDialogOpen = false
	case actionMoveUp, actionMoveDown:
		m.toggleSelectedMode()
	case actionActivate:
		return m.applySelectedMode()
	case actionShowArtists:
		m.selectedMode = artistsMode
		return m.applySelectedMode()
	case actionShowPlaylists:
		m.selectedMode = playlistsMode
		return m.applySelectedMode()
	}

	return nil
}

func (m *Model) toggleSelectedMode() {
	if m.selectedMode == artistsMode {
		m.selectedMode = playlistsMode
		return
	}

	m.selectedMode = artistsMode
}

func (m *Model) applySelectedMode() tea.Cmd {
	mode := m.selectedMode
	m.modeDialogOpen = false
	return m.switchMode(mode)
}

func (m *Model) moveSidebarSelection(delta int) tea.Cmd {
	switch m.mode {
	case playlistsMode:
		playlists := m.filteredPlaylists()
		if len(playlists) == 0 {
			return nil
		}

		previous := m.selectedPlaylist
		position := m.selectedPlaylistPosition(playlists)
		position = clamp(position+delta, 0, len(playlists)-1)
		m.selectedPlaylist = playlists[position].index
		m.selectedSong = 0
		if m.selectedPlaylist != previous {
			m.sidebarMarqueeOffset = 0
			return tea.Batch(m.loadSelectedPlaylist(), tickSidebarMarquee())
		}
	case artistsMode:
		artists := m.filteredArtists()
		if len(artists) == 0 {
			return nil
		}

		previous := m.selectedArtist
		position := m.selectedArtistPosition(artists)
		position = clamp(position+delta, 0, len(artists)-1)
		m.selectedArtist = artists[position].index
		m.selectedArtistRow = 0
		m.collapsedAlbums = nil
		m.selectedSong = 0
		if m.selectedArtist != previous {
			m.sidebarMarqueeOffset = 0
			return tea.Batch(m.loadSelectedArtist(), tickSidebarMarquee())
		}
	}

	return nil
}

func (m *Model) switchMode(mode sidebarMode) tea.Cmd {
	if m.mode == mode {
		return nil
	}

	m.mode = mode
	m.focused = playlistsPane
	m.sidebarMarqueeOffset = 0
	m.selectedSong = 0
	m.songs = nil
	m.albums = nil
	m.selectedArtistRow = 0
	m.collapsedAlbums = nil

	switch mode {
	case playlistsMode:
		if len(m.playlists) == 0 {
			return tea.Batch(m.loadPlaylists(), tickSidebarMarquee())
		}

		return tea.Batch(m.loadSelectedSidebarItem(), tickSidebarMarquee())
	case artistsMode:
		if len(m.artists) == 0 {
			return tea.Batch(m.loadArtists(), tickSidebarMarquee())
		}

		return tea.Batch(m.loadSelectedSidebarItem(), tickSidebarMarquee())
	default:
		return nil
	}
}

func (m *Model) loadSelectedSidebarItem() tea.Cmd {
	if m.selectedSidebarItemLoaded() {
		return nil
	}

	switch m.mode {
	case playlistsMode:
		return m.loadSelectedPlaylist()
	case artistsMode:
		return m.loadSelectedArtist()
	default:
		return nil
	}
}

func (m Model) selectedSidebarItemLoaded() bool {
	switch m.mode {
	case playlistsMode:
		return len(m.playlists) > 0 && m.loadedPlaylistID == m.playlists[m.selectedPlaylist].ID
	case artistsMode:
		return len(m.artists) > 0 && m.loadedArtistID == m.artists[m.selectedArtist].ID
	default:
		return false
	}
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
		return songsLoadedMsg{playlistID: playlistID, songs: songs, err: err}
	}
}

func (m *Model) loadArtists() tea.Cmd {
	m.loading = true

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), loadTimeout)
		defer cancel()

		artists, err := m.client.ListArtists(ctx)
		return artistsLoadedMsg{artists: artists, err: err}
	}
}

func (m *Model) loadSelectedArtist() tea.Cmd {
	if len(m.artists) == 0 {
		return nil
	}

	artistID := m.artists[m.selectedArtist].ID
	m.loading = true
	m.albums = nil
	m.songs = nil

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), loadTimeout)
		defer cancel()

		albums, err := m.client.GetArtistAlbums(ctx, artistID)
		if err != nil {
			return artistLoadedMsg{artistID: artistID, err: err}
		}

		var groups []albumGroup
		var songs []navidrome.Song
		for _, album := range albums {
			albumSongs, err := m.client.GetAlbumSongs(ctx, album.ID)
			if err != nil {
				return artistLoadedMsg{artistID: artistID, err: err}
			}

			groups = append(groups, albumGroup{album: album, songs: albumSongs})
			songs = append(songs, albumSongs...)
		}

		return artistLoadedMsg{artistID: artistID, albums: groups, songs: songs}
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

	return m.playSongFromList(m.songs[index], index, m.songs, m.currentLibraryPlaybackSource())
}

func (m *Model) playSongAtIndex(song navidrome.Song, index int) tea.Cmd {
	return m.playSongFromList(song, index, nil, "")
}

func (m *Model) playSongFromList(song navidrome.Song, index int, songs []navidrome.Song, source string) tea.Cmd {
	client := m.client
	player := m.player
	playbackID := m.playbackID + 1
	m.playbackID = playbackID
	m.playbackSongs = append([]navidrome.Song(nil), songs...)
	if source != "" {
		m.playbackSource = source
	}

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
	if len(m.queue) > 0 {
		return m.consumeQueuedSongAt(0)
	}

	playbackSongs := m.activePlaybackSongs()
	if len(playbackSongs) == 0 {
		return nil
	}
	if m.currentSongIndex < 0 {
		return nil
	}

	nextIndex := m.currentSongIndex + 1
	if nextIndex >= len(playbackSongs) {
		m.currentSong = nil
		m.elapsed = 0
		m.duration = 0
		m.paused = false
		m.player.Stop()
		return nil
	}

	m.syncVisibleSelectionToPlaybackSong(playbackSongs[nextIndex], nextIndex)
	return m.playSongFromList(playbackSongs[nextIndex], nextIndex, playbackSongs, m.playbackSource)
}

func (m *Model) playPreviousSong() tea.Cmd {
	if m.currentSong == nil {
		return nil
	}
	if m.currentSongIndex < 0 {
		return m.seekStart()
	}
	playbackSongs := m.activePlaybackSongs()
	if len(playbackSongs) == 0 {
		return m.seekStart()
	}

	if m.elapsed > restartLimit || m.currentSongIndex == 0 {
		return m.seekStart()
	}

	previousIndex := m.currentSongIndex - 1
	m.syncVisibleSelectionToPlaybackSong(playbackSongs[previousIndex], previousIndex)
	return m.playSongFromList(playbackSongs[previousIndex], previousIndex, playbackSongs, m.playbackSource)
}

func (m Model) currentLibraryPlaybackSource() string {
	switch m.mode {
	case playlistsMode:
		if len(m.playlists) > 0 {
			return "Playlist: " + m.playlists[clamp(m.selectedPlaylist, 0, len(m.playlists)-1)].Name
		}
	case artistsMode:
		if len(m.artists) > 0 {
			return "Artist: " + m.artists[clamp(m.selectedArtist, 0, len(m.artists)-1)].Name
		}
	}

	return ""
}

func (m Model) activePlaybackSongs() []navidrome.Song {
	if len(m.playbackSongs) > 0 {
		return m.playbackSongs
	}
	if m.currentSong == nil || m.currentSongIndex < 0 || m.currentSongIndex >= len(m.songs) {
		return nil
	}
	if m.songs[m.currentSongIndex].ID != m.currentSong.ID {
		return nil
	}

	return m.songs
}

func (m *Model) syncVisibleSelectionToPlaybackSong(song navidrome.Song, index int) {
	if index < 0 || index >= len(m.songs) {
		return
	}
	if m.songs[index].ID != song.ID {
		return
	}

	m.selectedSong = index
}

func (m *Model) seekStart() tea.Cmd {
	player := m.player

	return func() tea.Msg {
		return seekMsg{elapsed: 0, err: player.SeekStart()}
	}
}

func (m *Model) seekRelative(seconds int) tea.Cmd {
	if m.currentSong == nil {
		return nil
	}

	elapsed := max(m.elapsed+seconds, 0)
	if m.currentDuration() > 0 {
		elapsed = min(elapsed, m.currentDuration())
	}
	m.showToast(formatSeekToast(seconds))
	player := m.player

	return func() tea.Msg {
		return seekMsg{elapsed: elapsed, err: player.Seek(seconds)}
	}
}

func (m *Model) stopPlayback() tea.Cmd {
	if m.currentSong == nil {
		return nil
	}

	player := m.player
	m.playbackID++
	m.showToast("Stopped")

	return func() tea.Msg {
		return stopMsg{err: player.Stop()}
	}
}

func (m *Model) adjustVolume(delta int) tea.Cmd {
	if m.currentSong == nil {
		return nil
	}

	player := m.player
	m.showToast(formatVolumeToast(delta))

	return func() tea.Msg {
		return volumeMsg{err: player.AdjustVolume(delta)}
	}
}

func formatSeekToast(seconds int) string {
	if seconds > 0 {
		return fmt.Sprintf("+%ds", seconds)
	}

	return fmt.Sprintf("%ds", seconds)
}

func formatVolumeToast(delta int) string {
	if delta > 0 {
		return "Volume Up"
	}

	return "Volume Down"
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

func tickSidebarMarquee() tea.Cmd {
	return tea.Tick(marqueePeriod, func(t time.Time) tea.Msg {
		return sidebarMarqueeTickMsg{}
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
