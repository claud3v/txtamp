package mainview

type action string

const (
	actionQuit          action = "quit"
	actionFocusSidebar  action = "focus_sidebar"
	actionFocusMainArea action = "focus_main_area"
	actionMoveUp        action = "move_up"
	actionMoveDown      action = "move_down"
	actionActivate      action = "activate"
	actionPlayPause     action = "play_pause"
	actionNextSong      action = "next_song"
	actionPreviousSong  action = "previous_song"
	actionShowArtists   action = "show_artists"
	actionShowPlaylists action = "show_playlists"
	actionCloseDialog   action = "close_dialog"
	actionStartSearch   action = "start_search"
	actionGlobalSearch  action = "global_search"
	actionToggleHelp    action = "toggle_help"
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
	{Key: "1", Action: actionShowArtists, Description: "Show artists"},
	{Key: "2", Action: actionShowPlaylists, Description: "Show playlists"},
	{Key: "up", Action: actionMoveUp, Description: "Move up"},
	{Key: "down", Action: actionMoveDown, Description: "Move down"},
	{Key: "enter", Action: actionActivate, Description: "Open or play selected item"},
	{Key: " ", Action: actionPlayPause, Description: "Play or pause"},
	{Key: "space", Action: actionPlayPause, Description: "Play or pause"},
	{Key: "n", Action: actionNextSong, Description: "Next song"},
	{Key: "]", Action: actionNextSong, Description: "Next song"},
	{Key: "p", Action: actionPreviousSong, Description: "Previous or restart song"},
	{Key: "[", Action: actionPreviousSong, Description: "Previous or restart song"},
}

func actionForKey(key string) (action, bool) {
	for _, binding := range defaultKeyBindings {
		if binding.Key == key {
			return binding.Action, true
		}
	}

	return "", false
}
