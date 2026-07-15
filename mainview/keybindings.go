package mainview

type action string

const (
	actionQuit           action = "quit"
	actionFocusSidebar   action = "focus_sidebar"
	actionFocusMainArea  action = "focus_main_area"
	actionMoveUp         action = "move_up"
	actionMoveDown       action = "move_down"
	actionActivate       action = "activate"
	actionPlayPause      action = "play_pause"
	actionStopPlayback   action = "stop_playback"
	actionNextSong       action = "next_song"
	actionPreviousSong   action = "previous_song"
	actionToggleRepeat   action = "toggle_repeat"
	actionToggleShuffle  action = "toggle_shuffle"
	actionSeekBackward   action = "seek_backward"
	actionSeekForward    action = "seek_forward"
	actionVolumeUp       action = "volume_up"
	actionVolumeDown     action = "volume_down"
	actionAddToQueue     action = "add_to_queue"
	actionToggleQueue    action = "toggle_queue"
	actionRemoveQueue    action = "remove_queue"
	actionClearQueue     action = "clear_queue"
	actionPlayQueue      action = "play_queue"
	actionQueueUp        action = "queue_up"
	actionQueueDown      action = "queue_down"
	actionExpandAlbums   action = "expand_albums"
	actionCollapseAlbums action = "collapse_albums"
	actionShowArtists    action = "show_artists"
	actionShowAlbums     action = "show_albums"
	actionShowPlaylists  action = "show_playlists"
	actionCloseDialog    action = "close_dialog"
	actionStartSearch    action = "start_search"
	actionGlobalSearch   action = "global_search"
	actionGoToArtist     action = "go_to_artist"
	actionToggleTheme    action = "toggle_theme"
	actionToggleHelp     action = "toggle_help"
)

type keyBinding struct {
	Key         string
	Action      action
	Description string
}

var defaultKeyBindings = []keyBinding{
	{Key: "ctrl+c", Action: actionQuit, Description: "Quit"},
	{Key: "?", Action: actionToggleHelp, Description: "Show shortcuts"},
	{Key: "esc", Action: actionCloseDialog, Description: "Close dialog"},
	{Key: "left", Action: actionFocusSidebar, Description: "Focus sidebar"},
	{Key: "right", Action: actionFocusMainArea, Description: "Focus songs"},
	{Key: "/", Action: actionStartSearch, Description: "Filter focused pane"},
	{Key: "s", Action: actionGlobalSearch, Description: "Global search"},
	{Key: "g", Action: actionGoToArtist, Description: "Go to sidebar letter"},
	{Key: "t", Action: actionToggleTheme, Description: "Switch theme"},
	{Key: "a", Action: actionAddToQueue, Description: "Add selected song or album to queue"},
	{Key: "q", Action: actionToggleQueue, Description: "Show queue"},
	{Key: "d", Action: actionRemoveQueue, Description: "Remove queued song"},
	{Key: "c", Action: actionClearQueue, Description: "Clear queue"},
	{Key: "P", Action: actionPlayQueue, Description: "Play queue from top"},
	{Key: "J", Action: actionQueueDown, Description: "Move queued song down"},
	{Key: "K", Action: actionQueueUp, Description: "Move queued song up"},
	{Key: "E", Action: actionExpandAlbums, Description: "Expand all albums"},
	{Key: "C", Action: actionCollapseAlbums, Description: "Collapse all albums"},
	{Key: "1", Action: actionShowArtists, Description: "Show artists"},
	{Key: "2", Action: actionShowAlbums, Description: "Show albums"},
	{Key: "3", Action: actionShowPlaylists, Description: "Show playlists"},
	{Key: "up", Action: actionMoveUp, Description: "Move up"},
	{Key: "down", Action: actionMoveDown, Description: "Move down"},
	{Key: "enter", Action: actionActivate, Description: "Open, play, or toggle selected item"},
	{Key: " ", Action: actionPlayPause, Description: "Play or pause"},
	{Key: "space", Action: actionPlayPause, Description: "Play or pause"},
	{Key: "x", Action: actionStopPlayback, Description: "Stop playback"},
	{Key: "n", Action: actionNextSong, Description: "Next song"},
	{Key: "]", Action: actionNextSong, Description: "Next song"},
	{Key: "p", Action: actionPreviousSong, Description: "Previous or restart song"},
	{Key: "[", Action: actionPreviousSong, Description: "Previous or restart song"},
	{Key: "r", Action: actionToggleRepeat, Description: "Cycle repeat mode"},
	{Key: "z", Action: actionToggleShuffle, Description: "Toggle shuffle"},
	{Key: ",", Action: actionSeekBackward, Description: "Seek backward 10 seconds"},
	{Key: ".", Action: actionSeekForward, Description: "Seek forward 10 seconds"},
	{Key: "+", Action: actionVolumeUp, Description: "Volume up"},
	{Key: "=", Action: actionVolumeUp, Description: "Volume up"},
	{Key: "-", Action: actionVolumeDown, Description: "Volume down"},
}

func actionForKey(key string) (action, bool) {
	for _, binding := range defaultKeyBindings {
		if binding.Key == key {
			return binding.Action, true
		}
	}

	return "", false
}
