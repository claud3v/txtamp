package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Application title.
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#8BE9FD"))

	// Smaller descriptive text.
	Subtitle = lipgloss.NewStyle().
			Faint(true).
			Foreground(lipgloss.Color("#7F8C8D"))

	// Labels beside form fields.
	Label = lipgloss.NewStyle().
		Bold(true).
		Width(12)

	// Normal body text.
	Text = lipgloss.NewStyle()

	// Error messages.
	Error = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5555"))

	// Success messages.
	Success = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#50FA7B"))

	// Primary button.
	Button = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8BE9FD"))

	// Button when selected.
	ButtonFocused = Button.Copy().
			Foreground(lipgloss.Color("#282A36")).
			Background(lipgloss.Color("#8BE9FD"))

	// Surrounds an entire form.
	Card = lipgloss.NewStyle().
		Padding(1, 2)

	// Page margin.
	Page = lipgloss.NewStyle().
		Padding(1, 2)

	// Main application sidebar.
	Sidebar = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(lipgloss.Color("#3A3F4B"))

	// Main content area.
	MainPane = lipgloss.NewStyle().
			Padding(1, 2)

	// Bottom player status area.
	PlayerBar = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("#3A3F4B"))

	// Pane title.
	PaneTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#8BE9FD"))

	// Sidebar mode selector.
	ModeSelector = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7F8C8D"))

	// Active sidebar mode.
	ModeSelectorActive = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#282A36")).
				Background(lipgloss.Color("#8BE9FD")).
				Padding(0, 1)

	// Centered modal dialog.
	Dialog = lipgloss.NewStyle().
		Padding(1, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8BE9FD")).
		Background(lipgloss.Color("#282A36"))

	// Focused pane title.
	PaneTitleFocused = PaneTitle.Copy().
				Underline(true)

	// Selected row.
	SelectedRow = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD"))

	// Selected row in the focused pane.
	SelectedRowFocused = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282A36")).
				Background(lipgloss.Color("#8BE9FD"))

	// Currently playing row.
	PlayingRow = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#50FA7B"))
)
