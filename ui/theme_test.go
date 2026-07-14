package ui

import "testing"

func TestDefaultThemeFeedsSemanticColors(t *testing.T) {
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
	for _, name := range []string{"default", "mono", "amber"} {
		if _, ok := ThemeByName(name); !ok {
			t.Fatalf("expected built-in theme %q", name)
		}
	}

	if _, ok := ThemeByName("missing"); ok {
		t.Fatal("expected missing theme to be unknown")
	}
}
