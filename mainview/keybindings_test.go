package mainview

import "testing"

func TestActionForKey(t *testing.T) {
	tests := []struct {
		key            string
		expectedAction action
		found          bool
	}{
		{key: "n", expectedAction: actionNextSong, found: true},
		{key: "]", expectedAction: actionNextSong, found: true},
		{key: "p", expectedAction: actionPreviousSong, found: true},
		{key: "[", expectedAction: actionPreviousSong, found: true},
		{key: "1", expectedAction: actionShowArtists, found: true},
		{key: "2", expectedAction: actionShowPlaylists, found: true},
		{key: "/", expectedAction: actionStartSearch, found: true},
		{key: "esc", expectedAction: actionCloseDialog, found: true},
		{key: "space", expectedAction: actionPlayPause, found: true},
		{key: "x", found: false},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			got, found := actionForKey(test.key)
			if found != test.found {
				t.Fatalf("expected found %v, got %v", test.found, found)
			}
			if got != test.expectedAction {
				t.Fatalf("expected action %q, got %q", test.expectedAction, got)
			}
		})
	}
}
