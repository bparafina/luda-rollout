package gif

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
)

//go:embed rollout.gif
var gifData []byte

// Render displays the embedded animated GIF in the terminal.
// Protocol priority:
//  1. iTerm2 inline image protocol
//  2. Kitty graphics protocol
//  3. chafa (Unicode/ANSI block art — works in tmux and most terminals)
//  4. Text banner fallback
//
// Errors are swallowed silently so they never break the actual rollout.
func Render() {
	switch {
	case os.Getenv("TERM_PROGRAM") == "iTerm.app":
		renderITerm2()
	case os.Getenv("KITTY_WINDOW_ID") != "":
		renderKitty()
	default:
		if !renderChafa() {
			renderFallback()
		}
	}
}

func renderITerm2() {
	encoded := base64.StdEncoding.EncodeToString(gifData)
	fmt.Printf("\033]1337;File=inline=1:%s\a", encoded)
}

func renderKitty() {
	encoded := base64.StdEncoding.EncodeToString(gifData)
	fmt.Printf("\033_Ga=T,f=100,m=0;%s\033\\", encoded)
}

// renderChafa shells out to chafa to render the GIF as Unicode block art.
// Returns true if chafa was found and ran successfully.
func renderChafa() bool {
	chafa, err := exec.LookPath("chafa")
	if err != nil {
		return false
	}

	f, err := os.CreateTemp("", "rollout-*.gif")
	if err != nil {
		return false
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(gifData); err != nil || f.Close() != nil {
		return false
	}

	cmd := exec.Command(chafa, "--size", "60x20", "--duration", "3", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

func renderFallback() {
	fmt.Fprintln(os.Stderr, "🎵 Rollout (My Business) - Ludacris 🎵")
}
