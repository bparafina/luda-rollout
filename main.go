package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bparafina/krew-rollout/internal/audio"
	"github.com/bparafina/krew-rollout/internal/gif"
	"github.com/bparafina/krew-rollout/internal/passthrough"
	"github.com/bparafina/krew-rollout/internal/playlist"
)

func main() {
	args, skipMedia := parseArgs(os.Args[1:])

	if !skipMedia {
		entry := playlist.Pick()
		dir := playlist.AssetDir()
		audio.Play(filepath.Join(dir, entry.Audio))
		gif.Render(filepath.Join(dir, entry.GIF))
	}

	if err := passthrough.Run(args); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

func parseArgs(args []string) ([]string, bool) {
	cleaned := make([]string, 0, len(args))
	skip := false
	for _, a := range args {
		if a == "--no-rollout" {
			skip = true
		} else {
			cleaned = append(cleaned, a)
		}
	}
	return cleaned, skip
}
