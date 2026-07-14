package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Background lipgloss.Color
	Text       lipgloss.Color
	Muted      lipgloss.Color
	Border     lipgloss.Color
	Accent     lipgloss.Color
	Playing    lipgloss.Color
	Error      lipgloss.Color
	SelectedFg lipgloss.Color
	SelectedBg lipgloss.Color
}

var DefaultTheme = Theme{
	Background: lipgloss.Color("#282A36"),
	Text:       lipgloss.Color("#F8F8F2"),
	Muted:      lipgloss.Color("#7F8C8D"),
	Border:     lipgloss.Color("#3A3F4B"),
	Accent:     lipgloss.Color("#8BE9FD"),
	Playing:    lipgloss.Color("#50FA7B"),
	Error:      lipgloss.Color("#FF5555"),
	SelectedFg: lipgloss.Color("#282A36"),
	SelectedBg: lipgloss.Color("#8BE9FD"),
}

var MonoTheme = Theme{
	Background: lipgloss.Color("#111111"),
	Text:       lipgloss.Color("#EEEEEE"),
	Muted:      lipgloss.Color("#888888"),
	Border:     lipgloss.Color("#444444"),
	Accent:     lipgloss.Color("#FFFFFF"),
	Playing:    lipgloss.Color("#CFCFCF"),
	Error:      lipgloss.Color("#FFFFFF"),
	SelectedFg: lipgloss.Color("#111111"),
	SelectedBg: lipgloss.Color("#EEEEEE"),
}

var AmberTheme = Theme{
	Background: lipgloss.Color("#1A1200"),
	Text:       lipgloss.Color("#FFDFA3"),
	Muted:      lipgloss.Color("#9A7A45"),
	Border:     lipgloss.Color("#5C421C"),
	Accent:     lipgloss.Color("#FFB000"),
	Playing:    lipgloss.Color("#FFD166"),
	Error:      lipgloss.Color("#FF6B35"),
	SelectedFg: lipgloss.Color("#1A1200"),
	SelectedBg: lipgloss.Color("#FFB000"),
}

var BuiltInThemes = map[string]Theme{
	"default": DefaultTheme,
	"mono":    MonoTheme,
	"amber":   AmberTheme,
}

func ThemeByName(name string) (Theme, bool) {
	theme, ok := BuiltInThemes[name]
	return theme, ok
}
