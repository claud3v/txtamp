package mainview

import (
	"fmt"
	"strings"
	"txtamp/config"
	"txtamp/navidrome"
	"txtamp/ui"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) toggleQueue() {
	if m.contentMode == queueContent {
		m.contentMode = libraryContent
		return
	}

	m.contentMode = queueContent
	m.focused = songsPane
	m.selectedQueue = clamp(m.selectedQueue, 0, max(len(m.queue)-1, 0))
}

func (m *Model) addSelectedSongToQueue() bool {
	if m.contentMode == queueContent {
		return false
	}

	song, ok := m.selectedPlayableSong()
	if !ok {
		return false
	}

	m.queue = append(m.queue, song)
	m.queueDirty = true
	m.showToast("Added to queue: " + song.Title)
	return true
}

func (m Model) selectedPlayableSong() (navidrome.Song, bool) {
	if m.contentMode == globalSearchContent {
		rows := m.globalSearchRows()
		if len(rows) == 0 {
			return navidrome.Song{}, false
		}

		row := rows[clamp(m.selectedSearchResult, 0, len(rows)-1)]
		if row.kind != searchSongResult {
			return navidrome.Song{}, false
		}

		return row.song, true
	}

	if len(m.songs) == 0 {
		return navidrome.Song{}, false
	}

	index := clamp(m.selectedSong, 0, len(m.songs)-1)
	return m.songs[index], true
}

func (m *Model) moveQueueSelection(delta int) {
	if len(m.queue) == 0 {
		return
	}

	m.selectedQueue = clamp(m.selectedQueue+delta, 0, len(m.queue)-1)
}

func (m *Model) removeSelectedQueueSong() bool {
	if m.contentMode != queueContent || len(m.queue) == 0 {
		return false
	}

	removed := m.queue[clamp(m.selectedQueue, 0, len(m.queue)-1)]
	m.removeQueueSongAt(m.selectedQueue)
	m.showToast("Removed from queue: " + removed.Title)
	return true
}

func (m *Model) clearQueue() bool {
	if len(m.queue) == 0 {
		return false
	}

	m.queue = nil
	m.selectedQueue = 0
	m.queueDirty = true
	m.showToast("Queue cleared")
	return true
}

func (m *Model) playQueueFromTop() tea.Cmd {
	if len(m.queue) == 0 {
		return nil
	}

	m.selectedQueue = 0
	m.showToast("Playing queue")
	return tea.Batch(m.consumeQueuedSongAt(0), clearToast(m.toastID))
}

func (m *Model) moveQueuedSong(delta int) bool {
	if m.contentMode != queueContent || len(m.queue) == 0 {
		return false
	}

	next := clamp(m.selectedQueue+delta, 0, len(m.queue)-1)
	if next == m.selectedQueue {
		return false
	}

	m.queue[m.selectedQueue], m.queue[next] = m.queue[next], m.queue[m.selectedQueue]
	m.selectedQueue = next
	m.queueDirty = true
	return true
}

func (m *Model) playSelectedQueueSong() tea.Cmd {
	if m.contentMode != queueContent || len(m.queue) == 0 {
		return nil
	}

	index := clamp(m.selectedQueue, 0, len(m.queue)-1)
	return m.playSongAtIndex(m.queue[index], -1)
}

func (m *Model) consumeQueuedSongAt(index int) tea.Cmd {
	if len(m.queue) == 0 {
		return nil
	}

	index = clamp(index, 0, len(m.queue)-1)
	song := m.queue[index]
	m.removeQueueSongAt(index)

	return tea.Batch(m.playSongAtIndex(song, -1), m.saveQueue())
}

func (m *Model) removeQueueSongAt(index int) {
	if len(m.queue) == 0 {
		return
	}

	index = clamp(index, 0, len(m.queue)-1)
	m.queue = append(m.queue[:index], m.queue[index+1:]...)
	m.selectedQueue = clamp(m.selectedQueue, 0, max(len(m.queue)-1, 0))
	m.queueDirty = true
}

func (m Model) saveQueue() tea.Cmd {
	songs := append([]navidrome.Song(nil), m.queue...)
	return func() tea.Msg {
		return queueSavedMsg{err: config.SaveQueue(songs)}
	}
}

func (m Model) renderQueue(width, height int) string {
	lines := []string{
		paneTitle("Queue", m.focused == songsPane),
		ui.Subtitle.Render(fmt.Sprintf("%d songs", len(m.queue))),
		"",
	}

	if len(m.queue) == 0 {
		lines = append(lines, ui.Subtitle.Render("Queue is empty"))
	} else {
		start, end := visibleRange(m.selectedQueue, len(m.queue), m.visibleQueueRows(height))
		for i := start; i < end; i++ {
			song := m.queue[i]
			titleWidth := max(width-18, 10)
			title := ui.Truncate(formatNowPlaying(song), titleWidth)
			line := fmt.Sprintf("%-*"+"s %5s", titleWidth, title, formatDuration(song.Duration))
			lines = append(lines, queueSongLine(line, i, m, width-4))
		}
	}

	return ui.MainPane.
		Width(width).
		Height(height).
		MaxHeight(height).
		Render(strings.Join(lines, "\n"))
}

func (m Model) visibleQueueRows(height int) int {
	return max(height-5, 1)
}
