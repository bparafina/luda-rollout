package gif

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"strings"
	"time"

	"github.com/BourgeoisBear/rasterm"
)

//go:embed rollout.gif
var gifData []byte

const cols = 80 // output width in terminal columns

// Render displays the GIF at path animated in the terminal.
// nowPlaying is shown as a header above the animation (may be empty).
// If path is empty or the file does not exist, the embedded rollout GIF is used.
// Protocol priority:
//  1. Kitty graphics protocol  (Kitty terminal)
//  2. iTerm2 inline image      (iTerm2 / WezTerm)
//  3. Animated ANSI half-block art in the alternate screen buffer
func Render(path, nowPlaying string) {
	data := gifData
	if path != "" {
		if d, err := os.ReadFile(path); err == nil {
			data = d
		}
	}

	g, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil || len(g.Image) == 0 {
		renderFallback()
		return
	}

	if !rasterm.IsTmuxScreen() {
		if rasterm.IsKittyCapable() {
			if err := rasterm.KittyWriteImage(os.Stdout, g.Image[0], rasterm.KittyImgOpts{}); err == nil {
				return
			}
		}
		if rasterm.IsItermCapable() {
			if err := rasterm.ItermWriteImage(os.Stdout, g.Image[0]); err == nil {
				return
			}
		}
	}

	renderAnsiAnimated(g, nowPlaying)
}

// renderAnsiAnimated plays all GIF frames using the alternate screen buffer so
// the animation never scrolls the terminal. Each frame is drawn by moving the
// cursor to the top-left of the GIF area with \x1b[H, which is instantaneous
// and produces no flicker or scroll. The alternate screen is exited when the
// animation finishes and the terminal returns to exactly where it was.
func renderAnsiAnimated(g *gif.GIF, nowPlaying string) {
	bounds := g.Image[0].Bounds()
	srcW := g.Config.Width
	srcH := g.Config.Height
	if srcW == 0 {
		srcW = bounds.Max.X
	}
	if srcH == 0 {
		srcH = bounds.Max.Y
	}

	// Determine background color for canvas resets between loops
	var bgColor color.Color = color.Transparent
	if g.BackgroundIndex < uint8(len(g.Image[0].Palette)) {
		bgColor = g.Image[0].Palette[g.BackgroundIndex]
	}

	// Enter alternate screen and hide cursor — the terminal saves its current
	// state and presents a clean buffer for the duration of the animation.
	fmt.Fprint(os.Stdout, "\x1b[?1049h\x1b[?25l")
	defer fmt.Fprint(os.Stdout, "\x1b[?25h\x1b[?1049l")

	// Print now-playing header at the very top of the alternate screen.
	if nowPlaying != "" {
		fmt.Fprintf(os.Stdout, "\x1b[1;1H%s\n\n", nowPlaying)
	}

	// Loop the animation so it runs alongside the audio clip (~20s).
	// 2 passes of a ~3s GIF ≈ 6s of visible animation before kubectl output appears.
	const loops = 2

	for loop := 0; loop < loops; loop++ {
		canvas := image.NewRGBA(image.Rect(0, 0, srcW, srcH))
		draw.Draw(canvas, canvas.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

		for i, frame := range g.Image {
			draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)

			rendered := renderFrame(canvas, srcW, srcH)

			// Jump to row 3 (below header + blank line) and redraw the frame.
			// \x1b[3;1H positions the cursor without scrolling anything.
			if nowPlaying != "" {
				fmt.Fprint(os.Stdout, "\x1b[3;1H")
			} else {
				fmt.Fprint(os.Stdout, "\x1b[1;1H")
			}
			fmt.Fprint(os.Stdout, rendered)

			delay := g.Delay[i]
			if delay <= 0 {
				delay = 6 // 60ms default
			}
			time.Sleep(time.Duration(delay) * 10 * time.Millisecond)

			switch g.Disposal[i] {
			case gif.DisposalBackground, gif.DisposalPrevious:
				draw.Draw(canvas, frame.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
			}
		}
	}
}

// renderFrame renders a full canvas as an ANSI half-block string.
func renderFrame(src image.Image, srcW, srcH int) string {
	rows := srcH * cols / srcW
	if rows%2 != 0 {
		rows++
	}

	var sb strings.Builder
	for row := 0; row < rows; row += 2 {
		for col := 0; col < cols; col++ {
			sx := col * srcW / cols
			syTop := row * srcH / rows
			syBot := (row + 1) * srcH / rows

			tr, tg, tb := rgb(src.At(sx, syTop))
			br, bg, bb := rgb(src.At(sx, syBot))

			fmt.Fprintf(&sb, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀",
				tr, tg, tb, br, bg, bb)
		}
		sb.WriteString("\x1b[0m\n")
	}
	return sb.String()
}

// rgb extracts 8-bit R, G, B from any color.
func rgb(c color.Color) (uint8, uint8, uint8) {
	r, g, b, _ := c.RGBA()
	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
}

func renderFallback() {
	fmt.Fprintln(os.Stderr, "🎵 Rollout (My Business) - Ludacris 🎵")
}
