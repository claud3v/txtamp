package mainview

import (
	"fmt"
	"txtamp/ui"
)

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

		nowPlaying = fmt.Sprintf("%s - %s", m.currentSong.Artist, m.currentSong.Title)
		progress = "00:00 / " + formatDuration(m.currentSong.Duration)
	}

	line := fmt.Sprintf("%s  %s  %s", status, bars, ui.Truncate(nowPlaying, max(width-34, 10)))
	if width > 30 {
		line = fmt.Sprintf("%-*s %s", max(width-14, 10), line, progress)
	}

	return ui.PlayerBar.
		Width(width - 2).
		Render(line)
}
