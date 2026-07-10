package mainview

import (
	"fmt"
	"txtamp/ui"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	width       int
	height      int
	connectedTo string
}

func New(connectedTo string) Model {
	return Model{connectedTo: connectedTo}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	content := fmt.Sprintf(
		`%s

You are connected to: %s

Press ctrl+c to quit.`,
		ui.Title.Render("TxtAmp"),
		ui.Success.Render(m.connectedTo),
	)

	if m.width > 0 && m.height > 0 {
		content = ui.Center(m.width, m.height, content)
	}

	view := tea.NewView(content)
	view.AltScreen = true

	return view
}
