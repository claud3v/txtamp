package auth

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestTabAndEnterMoveFocus(t *testing.T) {
	m := New()

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(Model)

	if m.focused != urlFocused {
		t.Fatalf("expected URL to be focused, got %v", m.focused)
	}

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if m.focused != usernameFocused {
		t.Fatalf("expected username to be focused, got %v", m.focused)
	}
}

func TestShiftTabMovesFocusBackward(t *testing.T) {
	m := New()

	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	m = updated.(Model)

	if m.focused != connectBtnFocused {
		t.Fatalf("expected connect button to be focused, got %v", m.focused)
	}
}

func TestConnectValidationFocusesInvalidField(t *testing.T) {
	m := NewWithValues("home", "not a url", "john", "secret", nil)
	m.focused = connectBtnFocused

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if cmd != nil {
		t.Fatal("expected validation failure to return no command")
	}
	if m.err == nil || m.err.Error() != "URL is not valid" {
		t.Fatalf("expected URL validation error, got %v", m.err)
	}
	if m.focused != urlFocused {
		t.Fatalf("expected URL to be focused, got %v", m.focused)
	}
}

func TestConnectResultDisplaysError(t *testing.T) {
	expectedErr := errors.New("connection failed")
	m := New()

	updated, _ := m.Update(ConnectResultMsg{Err: expectedErr})
	m = updated.(Model)

	if !errors.Is(m.err, expectedErr) {
		t.Fatalf("expected error to be stored on model, got %v", m.err)
	}
}

func TestValidConnectReturnsCommand(t *testing.T) {
	m := NewWithValues("home", "https://music.example.com", "john", "secret", nil)
	m.focused = connectBtnFocused

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(Model)

	if m.err != nil {
		t.Fatalf("expected no validation error, got %v", m.err)
	}
	if cmd == nil {
		t.Fatal("expected connect command")
	}
}
