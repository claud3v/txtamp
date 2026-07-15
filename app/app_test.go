package app

import (
	"testing"
	"txtamp/ui"
)

func TestApplyThemeFromEnv(t *testing.T) {
	t.Setenv("TXTAMP_THEME", "neon-grid")
	defer ui.ApplyTheme(ui.DefaultTheme)

	applyThemeFromEnv()

	if ui.CurrentTheme.Accent != ui.NeonGridTheme.Accent {
		t.Fatalf("expected neon grid theme")
	}
	if ui.CurrentThemeName != "neon-grid" {
		t.Fatalf("expected canonical theme name, got %q", ui.CurrentThemeName)
	}
}

func TestApplyThemeFromEnvIgnoresUnknownTheme(t *testing.T) {
	ui.ApplyTheme(ui.DefaultTheme)
	t.Setenv("TXTAMP_THEME", "not-a-theme")

	applyThemeFromEnv()

	if ui.CurrentTheme.Accent != ui.DefaultTheme.Accent {
		t.Fatalf("expected unknown theme to keep current theme")
	}
}
