package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
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
	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Metal"); err != nil {
		panic(err)
	}

}
