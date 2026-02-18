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
