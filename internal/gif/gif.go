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

// Render displays the embedded GIF animated in the terminal.
// Protocol priority:
//  1. Kitty graphics protocol  (Kitty terminal)
//  2. iTerm2 inline image      (iTerm2 / WezTerm)
//  3. Animated ANSI half-block art (▀ + 24-bit color — works everywhere)
func Render() {
	g, err := gif.DecodeAll(bytes.NewReader(gifData))
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

	renderAnsiAnimated(g)
}

// renderAnsiAnimated plays all GIF frames in-place using ANSI half-block art,
// looping the animation to stay in sync with the audio clip.
func renderAnsiAnimated(g *gif.GIF) {
	bounds := g.Image[0].Bounds()
	srcW := g.Config.Width
	srcH := g.Config.Height
	if srcW == 0 {
		srcW = bounds.Max.X
	}
	if srcH == 0 {
		srcH = bounds.Max.Y
	}

	// Number of terminal rows the image occupies (half-blocks = 2 pixels per row)
	termRows := (srcH * cols / srcW)
	if termRows%2 != 0 {
		termRows++
	}
	termRows /= 2

	// Determine background color for canvas resets between loops
	var bgColor color.Color = color.Transparent
	if g.BackgroundIndex < uint8(len(g.Image[0].Palette)) {
		bgColor = g.Image[0].Palette[g.BackgroundIndex]
	}

	// Loop the animation so it runs alongside the audio clip (~20s).
	// 5 passes of a ~3s GIF ≈ 15s of visible animation before kubectl output appears.
	const loops = 2

	firstFrame := true
	for loop := 0; loop < loops; loop++ {
		// Reset canvas at the start of each loop
		canvas := image.NewRGBA(image.Rect(0, 0, srcW, srcH))
		draw.Draw(canvas, canvas.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

		for i, frame := range g.Image {
			// Composite this frame onto the canvas
			draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)

			rendered := renderFrame(canvas, srcW, srcH)

			if !firstFrame {
				// Move cursor back up to overwrite previous frame
				fmt.Fprintf(os.Stdout, "\x1b[%dA", termRows+1)
			}

			fmt.Fprint(os.Stdout, rendered)
			firstFrame = false

			// Respect frame delay (in 100ths of a second; minimum 60ms)
			delay := g.Delay[i]
			if delay <= 0 {
				delay = 6 // 60ms default
			}
			time.Sleep(time.Duration(delay) * 10 * time.Millisecond)

			// Handle disposal
			switch g.Disposal[i] {
			case gif.DisposalBackground:
				draw.Draw(canvas, frame.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
			case gif.DisposalPrevious:
				draw.Draw(canvas, frame.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
			}
		}
	}

	fmt.Fprint(os.Stdout, "\x1b[0m\n")
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
