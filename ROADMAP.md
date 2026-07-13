# Roadmap

TXTAMP is currently focused on becoming a dependable daily-use terminal music player before adding richer visual features.

## Done

- Connect to Navidrome.
- Save encrypted credentials.
- Restore saved sessions.
- Browse playlists.
- Browse artists with albums and songs.
- Search locally within the current pane.
- Search globally across artists, albums, and songs.
- Play songs through mpv.
- Play, pause, next, and previous.
- Automatically play the next song.
- Show playback progress.
- Show keyboard shortcuts.
- Add songs to a queue.
- View, remove, and reorder queued songs.
- Persist the queue in `~/.txtamp/queue.json`.

## Next

- Finish queue controls.
  - Clear the queue.
  - Play the queue from the top.
  - Show stronger queue feedback in the player or status bar.
  - Decide whether queue playback should have a distinct mode from normal playlist playback.

- Improve playback controls.
  - Seek forward and backward.
  - Stop playback.
  - Control volume.
  - Make playing, paused, and stopped states more obvious.

- Make albums first-class.
  - Navigate artist -> albums -> songs more explicitly.
  - Add an album-focused view.
  - Support playing or queueing a whole album.

- Polish the main UI.
  - Tighten queue, search, sidebar, and player layout consistency.
  - Improve empty, loading, and error states.
  - Keep shortcut hints accurate as bindings grow.

- Add richer music metadata.
  - Lyrics.
  - Album art, if terminal support is practical.

## Later

- User-configurable keybindings.
- Persist more local state, such as last selected view or playlist.
- Better queue recovery if saved songs are missing or changed on the server.
- Optional shuffle and repeat modes.
- Packaging or install instructions.
