# krew-rollout

A kubectl krew plugin written in Go that intercepts `kubectl rollout` commands and celebrates each deployment with the Ludacris "Rollout (My Business)" audio clip and an animated GIF.

## Project Purpose

This plugin wraps `kubectl rollout` so that every time a user triggers a rollout, the terminal plays audio and displays an animated GIF of the Rollout music video before passing through to the real kubectl rollout command. It is a fun, opinionated developer experience tool.

## Architecture

- **Language**: Go
- **Distribution**: kubectl krew plugin (`kubectl-rollout`)
- **Entry point**: `main.go` — parses args, triggers media, then delegates to `kubectl rollout`
- **Plugin name**: `rollout` (invoked as `kubectl rollout` when installed via krew, shadowing the built-in)
- **Media playback**: embedded or bundled audio + GIF assets
- **Passthrough**: after playing media, executes the real `kubectl rollout` with all original arguments

## Key Behaviors

1. User runs any `kubectl rollout` subcommand (e.g. `kubectl rollout restart`, `kubectl rollout status`)
2. Plugin plays the Ludacris "Rollout" audio clip (non-blocking or brief blocking)
3. Plugin renders the GIF in the terminal (using a terminal GIF renderer or kitty/iTerm protocol)
4. Plugin passes all arguments through to the underlying `kubectl rollout` and streams its output

## File Structure

```
krew-rollout/
├── CLAUDE.md
├── main.go               # Entry point: media trigger + kubectl passthrough
├── media/
│   ├── rollout.mp3       # Ludacris audio clip
│   └── rollout.gif       # Animated GIF
├── internal/
│   ├── audio/            # Audio playback (beep, oto, or exec mpv/afplay)
│   └── gif/              # Terminal GIF rendering
├── go.mod
├── go.sum
└── plugin.yaml           # Krew plugin manifest
```

## Implementation Notes

### Audio Playback
- Use `github.com/faiface/beep` or shell out to `afplay` (macOS), `paplay` (Linux), `ffplay` as fallback
- Audio should play asynchronously so it doesn't block the rollout output
- Embed the audio file using Go `//go:embed`

### GIF Rendering
- Use `github.com/charmbracelet/vhs` or `github.com/nicowillis/termgif` for terminal rendering
- Alternatively, use the iTerm2 inline image protocol or Kitty terminal graphics protocol
- Fall back gracefully if the terminal does not support inline images (print ASCII art or skip)
- Embed the GIF using Go `//go:embed`

### Krew Plugin Manifest (`plugin.yaml`)
- Platform binaries for darwin/amd64, darwin/arm64, linux/amd64, linux/arm64
- Short description: "Plays Rollout by Ludacris on every kubectl rollout"

### kubectl Passthrough
- After media plays, `exec` (or `os/exec`) the real `kubectl rollout` with all `os.Args[1:]`
- Stream stdout/stderr directly so rollout output is not buffered
- Preserve exit codes

## Dependencies

- `github.com/faiface/beep` — cross-platform audio
- `github.com/charmbracelet/bubbletea` or direct ANSI escape sequences — terminal GIF display
- Standard library `os/exec` — kubectl passthrough

## Build & Distribution

```bash
# Build
go build -o kubectl-rollout .

# Install locally for testing
mv kubectl-rollout ~/.krew/bin/

# Package for krew
# Update plugin.yaml with sha256 checksums of release tarballs
```

## Development Conventions

- Keep `main.go` minimal — media and passthrough logic lives in `internal/`
- Embedded assets via `//go:embed media/*`
- No configuration file needed — the plugin always plays media
- Errors in media playback should be swallowed silently; never break the actual rollout
- Support `--no-rollout` flag to skip media for CI/scripting contexts
