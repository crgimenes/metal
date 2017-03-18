package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var square *ebiten.Image

func block(screen *ebiten.Image) {
	if square == nil {

		// Create an 16x16 image
		square, _ = ebiten.NewImage(16, 16, ebiten.FilterNearest)
	}

	// Fill the square with the white color

	square.Fill(color.White)

	// Create an empty option struct
	opts := &ebiten.DrawImageOptions{}

	// Add the Translate effect to the option struct.
	w, h := screen.Size()
	w = (w / 2) - (16 / 2)
	h = (h / 2) - (16 / 2)
	opts.GeoM.Translate(float64(w), float64(h))

	// Draw the square image to the screen with an empty option
	screen.DrawImage(square, opts)
}

func update(screen *ebiten.Image) error {

	screen.Fill(color.NRGBA{0x00, 0x00, 0xff, 0xff})

	block(screen)

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		ebitenutil.DebugPrint(screen, "You're pressing the 'UP' button.")
	}
	// When the "down arrow key" is pressed..
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		ebitenutil.DebugPrint(screen, "\nYou're pressing the 'DOWN' button.")
	}
	// When the "left arrow key" is pressed..
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		ebitenutil.DebugPrint(screen, "\n\nYou're pressing the 'LEFT' button.")
	}
	// When the "right arrow key" is pressed..
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		ebitenutil.DebugPrint(screen, "\n\n\nYou're pressing the 'RIGHT' button.")
	}

	// When the "left mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		ebitenutil.DebugPrint(screen, "You're pressing the 'LEFT' mouse button.")
	}
	// When the "right mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		ebitenutil.DebugPrint(screen, "\nYou're pressing the 'RIGHT' mouse button.")
	}
	// When the "middle mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		ebitenutil.DebugPrint(screen, "\n\nYou're pressing the 'MIDDLE' mouse button.")
	}

	x, y := ebiten.CursorPosition()

	// Display the information with "X: xx, Y: xx" format
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Metal"); err != nil {
		panic(err)
	}

}
