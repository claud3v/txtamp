package ui

import "github.com/charmbracelet/lipgloss"

var (
	CurrentTheme     = DefaultTheme
	CurrentThemeName = "txtamp-classic"

	ColorBackground lipgloss.Color
	ColorText       lipgloss.Color
	ColorMuted      lipgloss.Color
	ColorBorder     lipgloss.Color
	ColorAccent     lipgloss.Color
	ColorPlaying    lipgloss.Color
	ColorError      lipgloss.Color
	ColorSelectedFg lipgloss.Color
	ColorSelectedBg lipgloss.Color

	Title               lipgloss.Style
	Subtitle            lipgloss.Style
	EmptyState          lipgloss.Style
	Label               lipgloss.Style
	Text                lipgloss.Style
	Error               lipgloss.Style
	Success             lipgloss.Style
	Button              lipgloss.Style
	ButtonFocused       lipgloss.Style
	Card                lipgloss.Style
	Page                lipgloss.Style
	Sidebar             lipgloss.Style
	MainPane            lipgloss.Style
	PlayerBar           lipgloss.Style
	StatusBar           lipgloss.Style
	PaneTitle           lipgloss.Style
	ModeSelector        lipgloss.Style
	ModeSelectorActive  lipgloss.Style
	Dialog              lipgloss.Style
	Toast               lipgloss.Style
	PaneTitleFocused    lipgloss.Style
	SelectedRow         lipgloss.Style
	SelectedRowFocused  lipgloss.Style
	PlayingRow          lipgloss.Style
	AlbumExpanded       lipgloss.Style
	AlbumCollapsed      lipgloss.Style
	SectionHeader       lipgloss.Style
	PlayerTitle         lipgloss.Style
	PlayerStatus        lipgloss.Style
	PlayerStatusPlaying lipgloss.Style
	PlayerStatusStopped lipgloss.Style
	PlayerMeta          lipgloss.Style
	PlayerUpNextLabel   lipgloss.Style
	PlayerUpNextTitle   lipgloss.Style
)

func init() {
	ApplyTheme(DefaultTheme)
}

func ApplyTheme(theme Theme) {
	CurrentTheme = theme
	ColorBackground = theme.Background
	ColorText = theme.Text
	ColorMuted = theme.Muted
	ColorBorder = theme.Border
	ColorAccent = theme.Accent
	ColorPlaying = theme.Playing
	ColorError = theme.Error
	ColorSelectedFg = theme.SelectedFg
	ColorSelectedBg = theme.SelectedBg

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	Subtitle = lipgloss.NewStyle().
		Foreground(ColorMuted)

	EmptyState = lipgloss.NewStyle().
		Foreground(ColorMuted).
		Italic(true)

	Label = lipgloss.NewStyle().
		Bold(true).
		Width(12)

	Text = lipgloss.NewStyle().
		Foreground(ColorText)

	Error = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorError)

	Success = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPlaying)

	Button = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAccent)

	ButtonFocused = Button.Copy().
		Foreground(ColorSelectedFg).
		Background(ColorSelectedBg)

	Card = lipgloss.NewStyle().
		Padding(1, 2)

	Page = lipgloss.NewStyle().
		Padding(1, 2)

	Sidebar = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(ColorBorder)

	MainPane = lipgloss.NewStyle().
		Padding(1, 2)

	PlayerBar = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(ColorBorder)

	StatusBar = lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(ColorMuted)

	PaneTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	ModeSelector = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorMuted)

	ModeSelectorActive = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSelectedFg).
		Background(ColorSelectedBg).
		Padding(0, 1)

	Dialog = lipgloss.NewStyle().
		Padding(1, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAccent).
		Background(ColorBackground)

	Toast = lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPlaying).
		Foreground(ColorPlaying).
		Background(ColorBackground)

	PaneTitleFocused = PaneTitle.Copy().
		Underline(true)

	SelectedRow = lipgloss.NewStyle().
		Foreground(ColorAccent)

	SelectedRowFocused = lipgloss.NewStyle().
		Foreground(ColorSelectedFg).
		Background(ColorSelectedBg)

	PlayingRow = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPlaying)

	AlbumExpanded = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	AlbumCollapsed = lipgloss.NewStyle().
		Foreground(ColorMuted)

	SectionHeader = lipgloss.NewStyle().
		Foreground(ColorMuted)

	PlayerTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorText)

	PlayerStatus = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	PlayerStatusPlaying = PlayerStatus.Copy().
		Foreground(ColorPlaying)

	PlayerStatusStopped = PlayerStatus.Copy().
		Foreground(ColorMuted)

	PlayerMeta = lipgloss.NewStyle().
		Foreground(ColorMuted)

	PlayerUpNextLabel = lipgloss.NewStyle().
		Foreground(ColorMuted)

	PlayerUpNextTitle = lipgloss.NewStyle().
		Foreground(ColorText)
}
