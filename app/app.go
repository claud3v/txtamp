package app

import (
	"fmt"
	"txtamp/auth"
	"txtamp/mainview"

	tea "charm.land/bubbletea/v2"
)

type screen int

const (
	authScreen screen = iota
	mainScreen
)

type Model struct {
	screen screen
	auth   auth.Model
	main   mainview.Model
}

func New() Model {
	connection, found, err := auth.LoadSavedConnection()
	if err != nil {
		if found {
			return Model{
				screen: authScreen,
				auth: auth.NewWithValues(
					connection.Alias,
					connection.Host,
					connection.Username,
					connection.Password,
					err,
				),
			}
		}

		return Model{
			screen: authScreen,
			auth:   auth.NewWithValues("", "", "", "", err),
		}
	}

	if found {
		return Model{
			screen: mainScreen,
			main:   mainview.New(connection.ConnectedTo),
		}
	}

	return Model{
		screen: authScreen,
		auth:   auth.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(auth.ConnectResultMsg); ok && msg.Err == nil {
		m.screen = mainScreen
		m.main = mainview.New(msg.ConnectedTo)
		return m, nil
	}

	switch m.screen {
	case authScreen:
		authModel, cmd := m.auth.Update(msg)
		m.auth = authModel.(auth.Model)
		return m, cmd
	case mainScreen:
		mainModel, cmd := m.main.Update(msg)
		m.main = mainModel.(mainview.Model)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() tea.View {
	switch m.screen {
	case authScreen:
		return m.auth.View()
	case mainScreen:
		return m.main.View()
	default:
		return tea.NewView("")
	}
}

func Run() error {
	p := tea.NewProgram(New())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("running UI: %w", err)
	}

	return nil
}
