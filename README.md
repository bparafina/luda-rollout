# kubectl-rollout

A kubectl krew plugin that plays Ludacris — *Rollout (My Business)* and renders an animated GIF every time you run `kubectl rollout`.

![demo](demo/demo.gif)

## Install

### Option 1 — krew (local manifest)

```bash
kubectl krew install --manifest=plugin.yaml --archive=kubectl-rollout.tar.gz
```

Build and package first:

```bash
make build
tar -czf kubectl-rollout.tar.gz kubectl-rollout
kubectl krew install --manifest=plugin.yaml --archive=kubectl-rollout.tar.gz
```

### Option 2 — manual install (fastest)

```bash
make install
```

This builds the binary and copies it directly to `~/.krew/bin/kubectl-rollout`.

### Option 3 — krew index (once published)

```bash
kubectl krew install rollout
```

> Not yet available in the krew index. Requires cutting a GitHub release and submitting to [krew-index](https://github.com/kubernetes-sigs/krew-index).

## Playlist

On each `kubectl rollout` invocation the plugin picks a random track from its playlist.
Assets live in `~/.kubectl-rollout/` (override with `KUBECTL_ROLLOUT_ASSETS` env var).

Each clip is cut to the moment in the song that references rolling.

| Track | Artist | Rolling cue | Audio file | GIF file |
|---|---|---|---|---|
| Rollout (My Business) | Ludacris | "Rollout!" hook | `rollout.mp3` | `rollout.gif` |
| Break Stuff | Limp Bizkit | Chorus peak | `break-stuff.mp3` | `break-stuff.gif` |
| Rollin' | Limp Bizkit | "Rollin' rollin' rollin'" | `rollin.mp3` | `rollin.gif` |
| Ridin' | Chamillionaire | "They see me rollin'" | `ridin.mp3` | `ridin.gif` |
| Proud Mary | CCR | "Rolling on the river" | `proud-mary.mp3` | `proud-mary.gif` |
| Roll with the Changes | REO Speedwagon | "Roll with the changes" | `roll-with-the-changes.mp3` | `roll-with-the-changes.gif` |
| Rock and Roll | Led Zeppelin | Opening riff | `rock-and-roll.mp3` | `rock-and-roll.gif` |
| Start Me Up | The Rolling Stones | Opening | `start-me-up.mp3` | `start-me-up.gif` |
| Jumpin' Jack Flash | The Rolling Stones | Opening | `jumpin-jack-flash.mp3` | `jumpin-jack-flash.gif` |
| Thunderstruck | AC/DC | "THUNDER!" drop | `thunderstruck.mp3` | `thunderstruck.gif` |
| Roll With It | Steve Winwood | "Just roll with it baby" | `roll-with-it.mp3` | `roll-with-it.gif` |
| Like a Rolling Stone | Bob Dylan | "How does it feel?" chorus | `like-a-rolling-stone.mp3` | `like-a-rolling-stone.gif` |

Entries with both files present are eligible for random selection. If `~/.kubectl-rollout/`
is empty the embedded Ludacris clip is used as the fallback.

## Media Assets

Audio and GIF assets are **not included** in this repository due to copyright.
Provide your own copies in `~/.kubectl-rollout/` using [yt-dlp](https://github.com/yt-dlp/yt-dlp) and [ffmpeg](https://ffmpeg.org/).

### Prerequisites

```bash
brew install yt-dlp ffmpeg   # macOS
# apt install yt-dlp ffmpeg  # Debian/Ubuntu
mkdir -p ~/.kubectl-rollout
```

### Download a track

```bash
# 1. Download audio
TRACK="Limp Bizkit Break Stuff"
SLUG="break-stuff"
START=45   # seconds into the song where the "rolling" reference hits

yt-dlp "ytsearch1:${TRACK}" -x --audio-format mp3 \
  -o "/tmp/${SLUG}-full.%(ext)s" --no-playlist
ffmpeg -ss ${START} -i /tmp/${SLUG}-full.mp3 -t 20 -c copy \
  ~/.kubectl-rollout/${SLUG}.mp3

# 2. Download GIF
yt-dlp "ytsearch1:${TRACK} official video" \
  -f 'bestvideo[ext=mp4]/bestvideo' \
  -o "/tmp/${SLUG}-video.%(ext)s" --no-playlist
ffmpeg -ss ${START} -i /tmp/${SLUG}-video.mp4 -t 3 \
  -vf 'fps=10,scale=240:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=64[p];[s1][p]paletteuse' \
  -loop 0 ~/.kubectl-rollout/${SLUG}.gif
```

### Suggested prompts for Claude / ChatGPT

Use an LLM to generate the full download script for all tracks at once:

> **Prompt:** Generate a bash script using yt-dlp and ffmpeg to download 20-second audio clips and 3-second GIF animations for the following tracks, cut to the timestamp where the song references "rolling". Output each file to `~/.kubectl-rollout/` using the filename listed. Tracks:
>
> - Rollout (My Business) — Ludacris → `rollout.mp3` / `rollout.gif` (start: 10s)
> - Break Stuff — Limp Bizkit → `break-stuff.mp3` / `break-stuff.gif` (start: 45s)
> - Rollin' — Limp Bizkit → `rollin.mp3` / `rollin.gif` (start: 48s)
> - Ridin' — Chamillionaire → `ridin.mp3` / `ridin.gif` (start: 5s)
> - Proud Mary — CCR → `proud-mary.mp3` / `proud-mary.gif` (start: 52s)
> - Roll with the Changes — REO Speedwagon → `roll-with-the-changes.mp3` / `roll-with-the-changes.gif` (start: 120s)
> - Rock and Roll — Led Zeppelin → `rock-and-roll.mp3` / `rock-and-roll.gif` (start: 5s)
> - Start Me Up — The Rolling Stones → `start-me-up.mp3` / `start-me-up.gif` (start: 0s)
> - Jumpin' Jack Flash — The Rolling Stones → `jumpin-jack-flash.mp3` / `jumpin-jack-flash.gif` (start: 0s)
> - Thunderstruck — AC/DC → `thunderstruck.mp3` / `thunderstruck.gif` (start: 13s)
> - Roll With It — Steve Winwood → `roll-with-it.mp3` / `roll-with-it.gif` (start: 47s)
> - Like a Rolling Stone — Bob Dylan → `like-a-rolling-stone.mp3` / `like-a-rolling-stone.gif` (start: 60s)

### Build placeholders (required for compiling)

The build embeds a fallback `rollout.mp3` and `rollout.gif`. Create minimal placeholders if you don't have the real files:

```bash
# Silent 1s MP3
ffmpeg -f lavfi -i anullsrc=r=44100:cl=mono -t 1 internal/audio/rollout.mp3

# 1×1 GIF
printf 'GIF89a\x01\x00\x01\x00\x00\x00\x00!\xf9\x04\x00\x00\x00\x00\x00,\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x00;' \
  > internal/gif/rollout.gif
```

## Usage

```bash
# Works as a drop-in replacement for kubectl rollout
kubectl rollout restart deployment/myapp -n production
kubectl rollout status deployment/myapp -n production
kubectl rollout undo deployment/myapp

# Skip media (for CI / scripting)
kubectl rollout --no-rollout restart deployment/myapp
```

## Build from source

```bash
git clone https://github.com/bparafina/krew-rollout
cd krew-rollout
make build       # build for current platform
make build-all   # cross-compile linux/darwin amd64/arm64
make install     # build + copy to ~/.krew/bin/
```
