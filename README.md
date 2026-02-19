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

| Track | Artist | Audio file | GIF file |
|---|---|---|---|
| Rollout (My Business) | Ludacris | `rollout.mp3` | `rollout.gif` |
| Break Stuff | Limp Bizkit | `break-stuff.mp3` | `break-stuff.gif` |
| Rollin' | Limp Bizkit | `rollin.mp3` | `rollin.gif` |
| Ridin' | Chamillionaire | `ridin.mp3` | `ridin.gif` |
| Proud Mary | CCR | `proud-mary.mp3` | `proud-mary.gif` |
| Roll with the Changes | REO Speedwagon | `roll-with-the-changes.mp3` | `roll-with-the-changes.gif` |

Entries with both files present are eligible for random selection. If `~/.kubectl-rollout/`
is empty the embedded Ludacris clip is used as the fallback.

## Media Assets

Audio and GIF assets are **not included** in this repository due to copyright.
Provide your own copies in `~/.kubectl-rollout/` using yt-dlp + ffmpeg:

```bash
mkdir -p ~/.kubectl-rollout

# Example: download audio clip
yt-dlp 'ytsearch1:Limp Bizkit Break Stuff' -x --audio-format mp3 \
  -o '/tmp/break-stuff-full.%(ext)s' --no-playlist
ffmpeg -ss 15 -i /tmp/break-stuff-full.mp3 -t 20 -c copy \
  ~/.kubectl-rollout/break-stuff.mp3

# Example: generate GIF from music video
yt-dlp 'ytsearch1:Limp Bizkit Break Stuff official video' \
  -f 'bestvideo[ext=mp4]/bestvideo' -o '/tmp/break-stuff-video.%(ext)s' --no-playlist
ffmpeg -ss 20 -i /tmp/break-stuff-video.mp4 -t 3 \
  -vf 'fps=10,scale=240:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=64[p];[s1][p]paletteuse' \
  -loop 0 ~/.kubectl-rollout/break-stuff.gif
```

The build also requires placeholder assets for the embedded fallback:

```bash
# Minimal placeholder — silent 1s MP3 and 1x1 GIF (or copy a real clip)
ffmpeg -f lavfi -i anullsrc=r=44100:cl=mono -t 1 internal/audio/rollout.mp3
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
