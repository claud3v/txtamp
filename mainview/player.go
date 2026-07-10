package mainview

import (
	"fmt"
	"strings"
	"txtamp/navidrome"
	"txtamp/ui"
)

func (m Model) renderPlayer(width int) string {
	status := "-"
	nowPlaying := "No song selected"
	elapsed := "00:00"
	duration := "00:00"

	if m.currentSong != nil {
		if m.paused {
			status = "||"
		} else {
			status = ">"
		}

		nowPlaying = formatNowPlaying(*m.currentSong)
		duration = formatDuration(m.currentSong.Duration)
	}

	innerWidth := max(width-6, 20)
	timeText := elapsed + " / " + duration
	barWidth := max(innerWidth-len(timeText)-3, 8)

	titleLine := fmt.Sprintf("%s  %s", status, ui.Truncate(nowPlaying, max(innerWidth-3, 8)))
	progressLine := fmt.Sprintf("%s  %s", progressBar(0, m.currentSongDuration(), barWidth), timeText)

	return ui.PlayerBar.
		Width(width - 2).
		Render(titleLine + "\n" + progressLine)
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

func progressBar(elapsed, duration, width int) string {
	if width <= 0 {
		return ""
	}
	if width == 1 {
		return "["
	}
	if width == 2 {
		return "[]"
	}

	innerWidth := width - 2
	filled := 0
	if duration > 0 && elapsed > 0 {
		filled = min(innerWidth, elapsed*innerWidth/duration)
	}

	return "[" + strings.Repeat("=", filled) + strings.Repeat("-", innerWidth-filled) + "]"
}

func (m Model) currentSongDuration() int {
	if m.currentSong == nil {
		return 0
	}

	return m.currentSong.Duration
}
