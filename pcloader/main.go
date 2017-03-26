package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/crgimenes/metal/pcloader/fonts"
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

var img *image.RGBA
var screenWidth = 320
var screenHeight = 240
var font fonts.Font8x8

func drawPix(x, y int) {
	pos := 4*y*screenWidth + 4*x
	img.Pix[pos] = 0xff
	img.Pix[pos+1] = 0xff
	img.Pix[pos+2] = 0xff
	img.Pix[pos+3] = 0xff
}

func update(screen *ebiten.Image) error {

	screen.Fill(color.NRGBA{0x00, 0x00, 0xff, 0xff})
	////
	/*
		drawPix(100, 100)
		drawPix(101, 100)
		drawPix(102, 100)
		drawPix(103, 100)
		drawPix(104, 100)
		drawPix(105, 100)
		drawPix(100, 100)
		drawPix(101, 101)
		drawPix(102, 102)
		drawPix(103, 103)
		drawPix(104, 104)
		drawPix(105, 105)
	*/
	var xa uint64
	var ya uint64

	for xa = 0; xa < 8; xa++ {
		for ya = 0; ya < 8; ya++ {
			if font.Bitmap[1][ya]&(0x80>>xa) != 0 {
				drawPix(int(xa)+100, int(ya)+100)
			}
		}
	}
	screen.ReplacePixels(img.Pix)

	block(screen)

	////
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

	font.Load()

	img = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "METAL BASIC 0.01"); err != nil {
		log.Fatal(err)
	}

	if err := ebiten.Run(update, 320, 240, 2, "Metal"); err != nil {
		panic(err)
	}

}
