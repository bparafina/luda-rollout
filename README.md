# kubectl-rollout

A kubectl krew plugin that plays Ludacris — *Rollout (My Business)* and renders an animated GIF every time you run `kubectl rollout`.

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

## Media Assets

The audio and GIF assets are **not included** in this repository due to copyright.

You must provide your own copies before building:

```bash
# Trim your audio to ~20 seconds starting at 10s in
ffmpeg -ss 10 -i source.mp3 -t 20 -c copy internal/audio/rollout.mp3

# Generate a GIF from video (3 seconds, 240px wide, 10fps)
ffmpeg -ss 12 -i source.mp4 -t 3 \
  -vf "fps=10,scale=240:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=64[p];[s1][p]paletteuse" \
  -loop 0 internal/gif/rollout.gif
```

Placeholder silent/blank assets work too — the build just needs the files to exist.

## Dependencies

`chafa` is required for GIF rendering in most terminals (tmux, terminal.app, etc).

```bash
# macOS
brew install chafa

# Linux
apt install chafa   # Debian/Ubuntu
dnf install chafa   # Fedora
```

iTerm2 and Kitty users get native inline rendering automatically — no `chafa` needed.

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
