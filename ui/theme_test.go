package ui

import "testing"

func TestDefaultThemeFeedsSemanticColors(t *testing.T) {
	ApplyTheme(DefaultTheme)
	if ColorAccent != DefaultTheme.Accent {
		t.Fatalf("expected accent color to come from default theme")
	}
	if ColorSelectedFg != DefaultTheme.SelectedFg {
		t.Fatalf("expected selected foreground to come from default theme")
	}
	if ColorSelectedBg != DefaultTheme.SelectedBg {
		t.Fatalf("expected selected background to come from default theme")
	}
}

func TestBuiltInThemes(t *testing.T) {
	for _, name := range []string{
		"default",
		"txtamp-classic",
		"mono",
		"monolith",
		"amber",
		"phosphor-amber",
		"retro",
		"sci-fi",
		"violet-terminal",
		"futuristic",
		"neon-grid",
		"light",
		"paperwhite",
		"dark",
		"deep-space",
	} {
		if _, ok := ThemeByName(name); !ok {
			t.Fatalf("expected built-in theme %q", name)
		}
	}

	if _, ok := ThemeByName("missing"); ok {
		t.Fatal("expected missing theme to be unknown")
	}
}

func TestThemeByNameNormalizesName(t *testing.T) {
	theme, ok := ThemeByName("  VIOLET-TERMINAL  ")
	if !ok {
		t.Fatal("expected theme")
	}
	if theme.Accent != VioletTerminalTheme.Accent {
		t.Fatalf("expected violet terminal theme")
	}
}

func TestApplyThemeRebuildsSemanticStyles(t *testing.T) {
	defer ApplyThemeByName("txtamp-classic")

	ApplyTheme(AmberTheme)
	if ColorAccent != AmberTheme.Accent {
		t.Fatalf("expected accent color to change")
	}
	if CurrentTheme.Accent != AmberTheme.Accent {
		t.Fatalf("expected current theme to change")
	}
}

func TestApplyThemeByNameStoresCanonicalName(t *testing.T) {
	defer ApplyThemeByName("txtamp-classic")

	if !ApplyThemeByName("retro") {
		t.Fatal("expected retro alias to apply")
	}
	if CurrentThemeName != "violet-terminal" {
		t.Fatalf("expected canonical theme name, got %q", CurrentThemeName)
	}
	if CurrentTheme.Accent != VioletTerminalTheme.Accent {
		t.Fatalf("expected violet terminal theme")
	}
}
