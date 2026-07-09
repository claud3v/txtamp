# AGENTS.md

## Project Overview

TXTAMP is a lightweight terminal client for Navidrome.

The project values:

- Simplicity over abstraction
- Fast startup
- Keyboard-first UX
- Small dependencies
- Idiomatic Go
- Readable code

## Guidelines

- Prefer standard library when practical.
- Avoid premature abstractions.
- Keep packages small and cohesive.
- Do not introduce interfaces unless multiple implementations exist.
- Keep dependencies to a minimum.
- Prefer composition over inheritance-like patterns.
- Follow idiomatic Go formatting and naming.
- Write clear error messages.
- Keep functions focused and reasonably small.

## Architecture

The application should remain simple.

Suggested package structure:

```
cmd/
internal/
    api/
    player/
    ui/
    config/
    model/
```

Avoid over-engineering.

Build the smallest thing that works before generalizing.

## UI

- Keyboard-first.
- Responsive.
- Minimal visual clutter.
- Use Bubble Tea idioms.
- Avoid unnecessary animations.

## Playback

Use mpv as the playback engine.

Do not implement custom audio decoding.

## Goal

TXTAMP should feel like a native terminal application, not a GUI squeezed into a terminal.
