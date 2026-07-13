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
		{key: ",", expectedAction: actionSeekBackward, found: true},
		{key: ".", expectedAction: actionSeekForward, found: true},
		{key: "1", expectedAction: actionShowArtists, found: true},
		{key: "2", expectedAction: actionShowPlaylists, found: true},
		{key: "/", expectedAction: actionStartSearch, found: true},
		{key: "s", expectedAction: actionGlobalSearch, found: true},
		{key: "a", expectedAction: actionAddToQueue, found: true},
		{key: "q", expectedAction: actionToggleQueue, found: true},
		{key: "d", expectedAction: actionRemoveQueue, found: true},
		{key: "c", expectedAction: actionClearQueue, found: true},
		{key: "P", expectedAction: actionPlayQueue, found: true},
		{key: "J", expectedAction: actionQueueDown, found: true},
		{key: "K", expectedAction: actionQueueUp, found: true},
		{key: "?", expectedAction: actionToggleHelp, found: true},
		{key: "esc", expectedAction: actionCloseDialog, found: true},
		{key: "space", expectedAction: actionPlayPause, found: true},
		{key: "x", expectedAction: actionStopPlayback, found: true},
		{key: "+", expectedAction: actionVolumeUp, found: true},
		{key: "=", expectedAction: actionVolumeUp, found: true},
		{key: "-", expectedAction: actionVolumeDown, found: true},
		{key: "z", found: false},
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
