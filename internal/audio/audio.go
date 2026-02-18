package audio

import (
	_ "embed"
	"os"
	"os/exec"
	"runtime"
)

//go:embed rollout.mp3
var mp3Data []byte

// Play writes the embedded MP3 to a temp file and starts the OS audio player
// as a detached subprocess. The player outlives the plugin process so the
// full clip plays while kubectl rollout output streams. Errors are silently
// swallowed so audio never breaks the actual rollout.
func Play() {
	f, err := os.CreateTemp("", "rollout-*.mp3")
	if err != nil {
		return
	}
	tmpPath := f.Name()
	// Do not defer os.Remove — the detached player still needs the file.
	// /tmp is cleaned up by the OS.

	if _, err := f.Write(mp3Data); err != nil {
		f.Close()
		return
	}
	if err := f.Close(); err != nil {
		return
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", tmpPath)
	case "linux":
		if path, err := exec.LookPath("paplay"); err == nil {
			cmd = exec.Command(path, tmpPath)
		} else if path, err := exec.LookPath("aplay"); err == nil {
			cmd = exec.Command(path, tmpPath)
		} else if path, err := exec.LookPath("ffplay"); err == nil {
			cmd = exec.Command(path, "-nodisp", "-autoexit", tmpPath)
		}
	case "windows":
		cmd = exec.Command("powershell", "-c",
			`(New-Object Media.SoundPlayer '`+tmpPath+`').PlaySync()`)
	}

	if cmd != nil {
		_ = cmd.Start() // detach — parent exit won't kill the player
	}
}
