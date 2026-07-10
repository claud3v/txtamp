package auth

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"errors"
	"fmt"
	"txtamp/ui"
)

type focusedTarget int

const (
	aliasFocused focusedTarget = iota
	urlFocused
	usernameFocused
	passwordFocused
	connectBtnFocused
	focusTargetCount
)

type Model struct {
	width  int
	height int

	alias    textinput.Model
	url      textinput.Model
	username textinput.Model
	password textinput.Model

	focused focusedTarget
	err     error
}

func New() Model {
	alias := textinput.New()
	alias.Placeholder = "My Music Server"

	url := textinput.New()
	url.Placeholder = "https://music.example.com"

	username := textinput.New()
	username.Placeholder = "username"

	password := textinput.New()
	password.Placeholder = "password"
	password.EchoMode = textinput.EchoPassword

	m := Model{
		alias:    alias,
		url:      url,
		username: username,
		password: password,
		focused:  aliasFocused,
	}

	m.setFocus()

	return m
}

func NewWithValues(alias, host, username, password string, err error) Model {
	m := New()
	m.alias.SetValue(alias)
	m.url.SetValue(host)
	m.username.SetValue(username)
	m.password.SetValue(password)
	m.err = err
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ConnectResultMsg:
		m.err = msg.Err

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		switch msg.String() {
		case "tab", "enter":
			if m.focused == connectBtnFocused {
				return m, m.handleConnectPress()
			}

			m.focused = (m.focused + 1) % focusTargetCount
			m.setFocus()
		case " ", "space":
			if m.focused == connectBtnFocused {
				return m, m.handleConnectPress()
			}
		case "shift+tab":
			m.focused = (m.focused + focusTargetCount - 1) % focusTargetCount
			m.setFocus()
		}
	}

	var cmd tea.Cmd

	switch m.focused {
	case aliasFocused:
		m.alias, cmd = m.alias.Update(msg)
	case urlFocused:
		m.url, cmd = m.url.Update(msg)
	case usernameFocused:
		m.username, cmd = m.username.Update(msg)
	case passwordFocused:
		m.password, cmd = m.password.Update(msg)
	}

	return m, cmd
}

func (m Model) View() tea.View {
	title := ui.Title.Render("TxtAmp")
	subtitle := ui.Subtitle.Render("It really whips the text llama's ass.")

	connectBtn := ui.Button.Render("Connect")
	if m.focused == connectBtnFocused {
		connectBtn = ui.ButtonFocused.Render("Connect")
	}

	error := ""
	if m.err != nil {
		error = ui.Error.Render(m.err.Error())
	}

	content := fmt.Sprintf(
		`%s
%s


Alias
%s

Host
%s

Username
%s

Password
%s

%s

%s

Press ctrl+c to quit.`,
		title,
		subtitle,
		m.alias.View(),
		m.url.View(),
		m.username.View(),
		m.password.View(),
		error,
		connectBtn,
	)

	if m.width > 0 && m.height > 0 {
		content = ui.Center(m.width, m.height, content)
	}

	view := tea.NewView(content)

	view.AltScreen = true

	return view
}

func (m *Model) setFocus() {
	m.alias.Blur()
	m.url.Blur()
	m.username.Blur()
	m.password.Blur()

	switch m.focused {
	case aliasFocused:
		m.alias.Focus()
	case urlFocused:
		m.url.Focus()
	case usernameFocused:
		m.username.Focus()
	case passwordFocused:
		m.password.Focus()
	}
}

func (m *Model) handleConnectPress() tea.Cmd {
	m.err = nil

	if !IsServerUrlValid(m.url.Value()) {
		m.err = errors.New("URL is not valid")
		m.focused = urlFocused
		m.setFocus()
		return nil
	}

	if !IsNotBlank(m.username.Value()) {
		m.err = errors.New("Username cannot be blank")
		m.focused = usernameFocused
		m.setFocus()
		return nil
	}

	if !IsNotBlank(m.password.Value()) {
		m.err = errors.New("Password cannot be blank")
		m.focused = passwordFocused
		m.setFocus()
		return nil
	}

	return connectServer(m.alias.Value(), m.url.Value(), m.username.Value(), m.password.Value())
}
