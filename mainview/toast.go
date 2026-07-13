package mainview

import "txtamp/ui"

func (m *Model) showToast(message string) {
	m.toastID++
	m.toast = message
}

func (m Model) renderToast() string {
	return ui.Toast.Render(ui.Truncate(m.toast, 64))
}
