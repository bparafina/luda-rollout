package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/bparafina/krew-rollout/internal/audio"
	"github.com/bparafina/krew-rollout/internal/gif"
	"github.com/bparafina/krew-rollout/internal/passthrough"
)

func main() {
	args, skipMedia := parseArgs(os.Args[1:])

	if !skipMedia {
		audio.Play()
		gif.Render()
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
