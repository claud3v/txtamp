package mainview

import (
	"fmt"
	"strings"
	"txtamp/ui"
)

func (m Model) renderStatusBar(width int) string {
	text := fmt.Sprintf("Queue %d  %s  Space Play/Pause  q Queue  ? Help", len(m.queue), m.contextualStatusHint())
	return ui.StatusBar.
		Width(width).
		Render(ui.Truncate(text, max(width-4, 1)))
}

func (m Model) contextualStatusHint() string {
	if m.mode == artistsMode && m.focused == playlistsPane && m.contentMode == libraryContent {
		if m.goToArtistPending {
			return "Letter Go To Artist  Esc Cancel"
		}
		return "Arrows Navigate  Enter Open  g Go"
	}
	if m.selectedArtistAlbumRow() != nil && m.focused == songsPane && m.contentMode == libraryContent {
		return "Arrows Navigate  Space Play Album  Enter Toggle  a Queue Album"
	}

	return "Arrows Navigate  Enter Play  a Add"
}

func (m Model) renderHelpDialog() string {
	width := 42
	lines := []string{
		ui.PaneTitle.Render("Shortcuts"),
		"",
	}

	for _, binding := range helpBindings() {
		lines = append(lines, helpBindingLine(binding, width))
	}

	return ui.Dialog.Render(strings.Join(lines, "\n"))
}

func helpBindings() []keyBinding {
	bindings := make([]keyBinding, 0, len(defaultKeyBindings))
	for _, binding := range defaultKeyBindings {
		if binding.Key == " " {
			continue
		}
		bindings = append(bindings, binding)
	}

	return bindings
}

func helpBindingLine(binding keyBinding, width int) string {
	key := displayKey(binding.Key)
	keyWidth := 12
	descWidth := max(width-keyWidth-1, 8)
	return ui.PaneTitle.Width(keyWidth).Render(key) + " " + ui.Truncate(binding.Description, descWidth)
}

func displayKey(key string) string {
	switch key {
	case "space":
		return "space"
	case "left":
		return "left"
	case "right":
		return "right"
	case "up":
		return "up"
	case "down":
		return "down"
	case "enter":
		return "enter"
	case "esc":
		return "esc"
	default:
		return key
	}
}
