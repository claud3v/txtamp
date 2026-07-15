package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

type ThemeOption struct {
	Name  string
	Label string
	Theme Theme
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

var VioletTerminalTheme = Theme{
	Background: lipgloss.Color("#16091F"),
	Text:       lipgloss.Color("#F3E8FF"),
	Muted:      lipgloss.Color("#8F77A8"),
	Border:     lipgloss.Color("#3B2454"),
	Accent:     lipgloss.Color("#D16BFF"),
	Playing:    lipgloss.Color("#00F5D4"),
	Error:      lipgloss.Color("#FF4FD8"),
	SelectedFg: lipgloss.Color("#16091F"),
	SelectedBg: lipgloss.Color("#D16BFF"),
}

var NeonGridTheme = Theme{
	Background: lipgloss.Color("#061A24"),
	Text:       lipgloss.Color("#D7FBFF"),
	Muted:      lipgloss.Color("#5E8893"),
	Border:     lipgloss.Color("#124452"),
	Accent:     lipgloss.Color("#00E5FF"),
	Playing:    lipgloss.Color("#7CFF6B"),
	Error:      lipgloss.Color("#FF3D81"),
	SelectedFg: lipgloss.Color("#061A24"),
	SelectedBg: lipgloss.Color("#00E5FF"),
}

var PaperwhiteTheme = Theme{
	Background: lipgloss.Color("#F7F3E8"),
	Text:       lipgloss.Color("#1F2933"),
	Muted:      lipgloss.Color("#7B8794"),
	Border:     lipgloss.Color("#D7D0C2"),
	Accent:     lipgloss.Color("#2563EB"),
	Playing:    lipgloss.Color("#047857"),
	Error:      lipgloss.Color("#B91C1C"),
	SelectedFg: lipgloss.Color("#F7F3E8"),
	SelectedBg: lipgloss.Color("#2563EB"),
}

var DeepSpaceTheme = Theme{
	Background: lipgloss.Color("#05070D"),
	Text:       lipgloss.Color("#DCE7F7"),
	Muted:      lipgloss.Color("#65758B"),
	Border:     lipgloss.Color("#1F2A3D"),
	Accent:     lipgloss.Color("#7DD3FC"),
	Playing:    lipgloss.Color("#A7F3D0"),
	Error:      lipgloss.Color("#F87171"),
	SelectedFg: lipgloss.Color("#05070D"),
	SelectedBg: lipgloss.Color("#7DD3FC"),
}

var BuiltInThemeOptions = []ThemeOption{
	{Name: "txtamp-classic", Label: "TxtAmp Classic", Theme: DefaultTheme},
	{Name: "violet-terminal", Label: "Violet Terminal", Theme: VioletTerminalTheme},
	{Name: "neon-grid", Label: "Neon Grid", Theme: NeonGridTheme},
	{Name: "phosphor-amber", Label: "Phosphor Amber", Theme: AmberTheme},
	{Name: "deep-space", Label: "Deep Space", Theme: DeepSpaceTheme},
	{Name: "paperwhite", Label: "Paperwhite", Theme: PaperwhiteTheme},
	{Name: "monolith", Label: "Monolith", Theme: MonoTheme},
}

var BuiltInThemes = map[string]Theme{
	"default":    BuiltInThemeOptions[0].Theme,
	"mono":       MonoTheme,
	"amber":      AmberTheme,
	"retro":      VioletTerminalTheme,
	"sci-fi":     VioletTerminalTheme,
	"futuristic": NeonGridTheme,
	"light":      PaperwhiteTheme,
	"dark":       DeepSpaceTheme,
}

func ThemeByName(name string) (Theme, bool) {
	name = canonicalThemeName(name)
	for _, option := range BuiltInThemeOptions {
		if option.Name == name {
			return option.Theme, true
		}
	}

	theme, ok := BuiltInThemes[name]
	return theme, ok
}

func ApplyThemeByName(name string) bool {
	theme, ok := ThemeByName(name)
	if !ok {
		return false
	}

	CurrentThemeName = canonicalThemeName(name)
	ApplyTheme(theme)
	return true
}

func normalizeThemeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func canonicalThemeName(name string) string {
	switch normalizeThemeName(name) {
	case "", "default":
		return "txtamp-classic"
	case "mono":
		return "monolith"
	case "amber":
		return "phosphor-amber"
	case "retro", "sci-fi":
		return "violet-terminal"
	case "futuristic":
		return "neon-grid"
	case "light":
		return "paperwhite"
	case "dark":
		return "deep-space"
	default:
		return normalizeThemeName(name)
	}
}
