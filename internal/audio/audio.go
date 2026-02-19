package audio

import (
	_ "embed"
	"os"
	"os/exec"
	"runtime"
)

//go:embed rollout.mp3
var mp3Data []byte

// Play starts the OS audio player for the given path as a detached subprocess.
// If path is empty or the file does not exist, the embedded Ludacris clip is used.
// Errors are silently swallowed so audio never breaks the actual rollout.
func Play(path string) {
	if path == "" || !fileExists(path) {
		f, err := os.CreateTemp("", "rollout-*.mp3")
		if err != nil {
			return
		}
		// Do not defer os.Remove — the detached player still needs the file.
		if _, err := f.Write(mp3Data); err != nil {
			f.Close()
			return
		}
		if err := f.Close(); err != nil {
			return
		}
		path = f.Name()
	}
	startPlayer(path)
}

func startPlayer(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", path)
	case "linux":
		if p, err := exec.LookPath("paplay"); err == nil {
			cmd = exec.Command(p, path)
		} else if p, err := exec.LookPath("aplay"); err == nil {
			cmd = exec.Command(p, path)
		} else if p, err := exec.LookPath("ffplay"); err == nil {
			cmd = exec.Command(p, "-nodisp", "-autoexit", path)
		}
	case "windows":
		cmd = exec.Command("powershell", "-c",
			`(New-Object Media.SoundPlayer '`+path+`').PlaySync()`)
	}
	if cmd != nil {
		_ = cmd.Start() // detach — parent exit won't kill the player
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
