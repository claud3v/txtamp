package mainview

import (
	"strings"
	"txtamp/navidrome"
	"txtamp/ui"

	"charm.land/bubbles/v2/progress"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderPlayer(width int) string {
	status := "Stopped"
	title := "No song selected"
	metadata := ""
	elapsed := "00:00"
	duration := "00:00"

	if m.currentSong != nil {
		if m.paused {
			status = "Paused"
		} else {
			status = "Playing"
		}

		title = formatSongTitle(*m.currentSong)
		metadata = formatSongMetadata(*m.currentSong)
		elapsed = formatDuration(m.elapsed)
		duration = formatDuration(m.currentDuration())
	}

	innerWidth := max(width-6, 20)
	timeText := elapsed + " / " + duration
	barWidth := max(innerWidth, 8)

	titleLine := joinLeftRight(status+"  "+title, timeText, innerWidth)
	metadataLine := joinLeftRight(metadata, m.visiblePlaybackSource(), innerWidth)
	progressLine := progressBar(m.elapsed, m.currentDuration(), barWidth)
	upNextLine := "Up next: " + ui.Truncate(m.upNextText(), max(innerWidth-9, 8))

	return ui.PlayerBar.
		Width(width - 2).
		Render(titleLine + "\n" + metadataLine + "\n" + progressLine + "\n" + upNextLine)
}

func (m Model) upNextText() string {
	if len(m.queue) > 0 {
		return formatSongTitle(m.queue[0])
	}

	playbackSongs := m.activePlaybackSongs()
	nextIndex := m.currentSongIndex + 1
	if m.currentSong != nil && m.currentSongIndex >= 0 && nextIndex < len(playbackSongs) {
		return formatSongTitle(playbackSongs[nextIndex])
	}

	return "-"
}

func (m Model) visiblePlaybackSource() string {
	if m.playbackSource == "" || m.currentSong == nil {
		return ""
	}
	if m.playbackSource == "Artist: "+m.currentSong.Artist {
		return ""
	}

	return "Source: " + m.playbackSource
}

func formatNowPlaying(song navidrome.Song) string {
	parts := []string{}
	if song.Artist != "" {
		parts = append(parts, song.Artist)
	}
	if song.Album != "" {
		parts = append(parts, song.Album)
	}
	if song.Title != "" {
		parts = append(parts, song.Title)
	}
	if len(parts) == 0 {
		return "Unknown song"
	}

	return strings.Join(parts, " - ")
}

func formatSongTitle(song navidrome.Song) string {
	if song.Title != "" {
		return song.Title
	}

	return "Unknown song"
}

func formatSongMetadata(song navidrome.Song) string {
	parts := []string{}
	if song.Artist != "" {
		parts = append(parts, song.Artist)
	}
	if song.Album != "" {
		parts = append(parts, song.Album)
	}

	return strings.Join(parts, " - ")
}

func joinLeftRight(left, right string, width int) string {
	if right == "" {
		return ui.Truncate(left, width)
	}

	rightWidth := lipgloss.Width(right)
	leftWidth := max(width-rightWidth-2, 1)
	left = ui.Truncate(left, leftWidth)
	padding := max(width-lipgloss.Width(left)-rightWidth, 1)

	return left + strings.Repeat(" ", padding) + right
}

func progressBar(elapsed, duration, width int) string {
	percent := 0.0
	if duration > 0 && elapsed > 0 {
		percent = float64(elapsed) / float64(duration)
		if percent > 1 {
			percent = 1
		}
	}

	bar := progress.New(
		progress.WithWidth(width),
		progress.WithFillCharacters('━', '─'),
		progress.WithoutPercentage(),
	)

	return bar.ViewAs(percent)
}

func (m Model) currentDuration() int {
	if m.duration > 0 {
		return m.duration
	}
	if m.currentSong == nil {
		return 0
	}

	return m.currentSong.Duration
}
