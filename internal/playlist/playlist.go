package playlist

import (
	"math/rand"
	"os"
	"path/filepath"
)

// Entry represents one track in the rollout playlist.
type Entry struct {
	Name   string
	Artist string
	Audio  string // filename within AssetDir, e.g. "break-stuff.mp3"
	GIF    string // filename within AssetDir, e.g. "break-stuff.gif"
}

// Tracks is the full playlist. Assets live in AssetDir().
var Tracks = []Entry{
	{Name: "Rollout (My Business)", Artist: "Ludacris", Audio: "rollout.mp3", GIF: "rollout.gif"},
	{Name: "Break Stuff", Artist: "Limp Bizkit", Audio: "break-stuff.mp3", GIF: "break-stuff.gif"},
	{Name: "Rollin'", Artist: "Limp Bizkit", Audio: "rollin.mp3", GIF: "rollin.gif"},
	{Name: "Ridin'", Artist: "Chamillionaire", Audio: "ridin.mp3", GIF: "ridin.gif"},
	{Name: "Proud Mary", Artist: "CCR", Audio: "proud-mary.mp3", GIF: "proud-mary.gif"},
	{Name: "Roll with the Changes", Artist: "REO Speedwagon", Audio: "roll-with-the-changes.mp3", GIF: "roll-with-the-changes.gif"},
}

// AssetDir returns the directory where playlist assets are stored.
// Override with the KUBECTL_ROLLOUT_ASSETS environment variable.
func AssetDir() string {
	if dir := os.Getenv("KUBECTL_ROLLOUT_ASSETS"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kubectl-rollout")
}

// Pick returns a random Entry whose audio and gif files both exist in AssetDir.
// Falls back to Tracks[0] (Ludacris) when no external assets are available,
// signalling that the caller should use the embedded fallback assets.
func Pick() Entry {
	dir := AssetDir()
	var available []Entry
	for _, t := range Tracks {
		if exists(filepath.Join(dir, t.Audio)) && exists(filepath.Join(dir, t.GIF)) {
			available = append(available, t)
		}
	}
	if len(available) == 0 {
		return Tracks[0]
	}
	return available[rand.Intn(len(available))]
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
