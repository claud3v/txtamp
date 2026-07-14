package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorBackground = lipgloss.Color("#282A36")
	ColorText       = lipgloss.Color("#F8F8F2")
	ColorMuted      = lipgloss.Color("#7F8C8D")
	ColorBorder     = lipgloss.Color("#3A3F4B")
	ColorAccent     = lipgloss.Color("#8BE9FD")
	ColorPlaying    = lipgloss.Color("#50FA7B")
	ColorError      = lipgloss.Color("#FF5555")
	ColorSelectedFg = ColorBackground
	ColorSelectedBg = ColorAccent

	// Application title.
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	// Smaller descriptive text.
	Subtitle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Labels beside form fields.
	Label = lipgloss.NewStyle().
		Bold(true).
		Width(12)

	// Normal body text.
	Text = lipgloss.NewStyle().
		Foreground(ColorText)

	// Error messages.
	Error = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorError)

	// Success messages.
	Success = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPlaying)

	// Primary button.
	Button = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAccent)

	// Button when selected.
	ButtonFocused = Button.Copy().
			Foreground(ColorSelectedFg).
			Background(ColorSelectedBg)

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
		BorderForeground(ColorBorder)

	// Main content area.
	MainPane = lipgloss.NewStyle().
			Padding(1, 2)

	// Bottom player status area.
	PlayerBar = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorBorder)

	// One-line command hint footer.
	StatusBar = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(ColorMuted)

	// Pane title.
	PaneTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	// Sidebar mode selector.
	ModeSelector = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorMuted)

	// Active sidebar mode.
	ModeSelectorActive = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSelectedFg).
				Background(ColorSelectedBg).
				Padding(0, 1)

	// Centered modal dialog.
	Dialog = lipgloss.NewStyle().
		Padding(1, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAccent).
		Background(ColorBackground)

	// Short-lived non-blocking notification.
	Toast = lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPlaying).
		Foreground(ColorPlaying).
		Background(ColorBackground)

	// Focused pane title.
	PaneTitleFocused = PaneTitle.Copy().
				Underline(true)

	// Selected row.
	SelectedRow = lipgloss.NewStyle().
			Foreground(ColorAccent)

	// Selected row in the focused pane.
	SelectedRowFocused = lipgloss.NewStyle().
				Foreground(ColorSelectedFg).
				Background(ColorSelectedBg)

	// Currently playing row.
	PlayingRow = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPlaying)

	// Expanded album header.
	AlbumExpanded = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	// Collapsed album header.
	AlbumCollapsed = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Section dividers inside lists.
	SectionHeader = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Player primary title.
	PlayerTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText)

	// Player state label.
	PlayerStatus = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	// Player secondary metadata.
	PlayerMeta = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Player upcoming item.
	PlayerUpNext = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
