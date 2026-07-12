package mainview

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestStatusBarShowsBasicCommands(t *testing.T) {
	m := loadedModel()

	content := m.renderStatusBar(100)
	for _, expected := range []string{"Arrows Navigate", "Space Play/Pause", "? Help"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("expected %q in status bar, got:\n%s", expected, content)
		}
	}
}

func TestHelpDialogShowsShortcuts(t *testing.T) {
	m := loadedModel()

	content := m.renderHelpDialog()
	for _, expected := range []string{"Shortcuts", "space", "Play or pause", "?", "Show shortcuts"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("expected %q in help dialog, got:\n%s", expected, content)
		}
	}
}

func TestQuestionMarkTogglesHelp(t *testing.T) {
	m := loadedModel()

	updated, cmd := m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected no command")
	}
	if !m.helpOpen {
		t.Fatal("expected help to open")
	}

	updated, _ = m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	m = updated.(Model)
	if m.helpOpen {
		t.Fatal("expected help to close")
	}
}

func TestEscClosesHelp(t *testing.T) {
	m := loadedModel()
	m.helpOpen = true

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	m = updated.(Model)

	if m.helpOpen {
		t.Fatal("expected help to close")
	}
}

func TestQuestionMarkOpensHelpWhileFiltering(t *testing.T) {
	m := loadedModel()
	m.searching = true
	m.searchPane = songsPane
	m.searchQuery = "iron"

	updated, _ := m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	m = updated.(Model)

	if !m.helpOpen {
		t.Fatal("expected help to open")
	}
	if m.searchQuery != "iron" {
		t.Fatalf("expected search query to stay unchanged, got %q", m.searchQuery)
	}
}

func TestViewOverlaysHelp(t *testing.T) {
	m := loadedModel()
	m.width = 120
	m.height = 30
	m.helpOpen = true

	view := m.View()
	if !strings.Contains(view.Content, "Shortcuts") {
		t.Fatalf("expected help overlay, got:\n%s", view.Content)
	}
	if !strings.Contains(view.Content, "? Help") {
		t.Fatalf("expected status bar behind help, got:\n%s", view.Content)
	}
}
