package main

import (
	"image"
	"image/color"
	"log"

	"github.com/crgimenes/metal/pcloader/fonts"
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

const screenWidth = 320  // 40 columns
const screenHeight = 240 // 30 rows

const rows = 30
const columns = 40
const rgbaSize = 4

var img *image.RGBA

var videoTextMemory [rows * columns]byte

var font fonts.Font8x8
var cursor int

func drawPix(x, y int) {
	pos := 4*y*screenWidth + 4*x
	img.Pix[pos] = 0xff
	img.Pix[pos+1] = 0xff
	img.Pix[pos+2] = 0xff
	img.Pix[pos+3] = 0xff
}

func drawOffPix(x, y int) {
	pos := 4*y*screenWidth + 4*x
	img.Pix[pos] = 0x00
	img.Pix[pos+1] = 0x00
	img.Pix[pos+2] = 0x00
	img.Pix[pos+3] = 0x00
}

func getBit(n int, pos uint64) bool {
	// from right to left
	val := n & (1 << pos)
	return (val > 0)
}

func drawChar(index byte, x, y int) {
	var a, b uint64
	for a = 0; a < 8; a++ {
		for b = 0; b < 8; b++ {
			if font.Bitmap[index][b]&(0x80>>a) != 0 {
				drawPix(int(a)+x, int(b)+y)
			} else {
				drawOffPix(int(a)+x, int(b)+y)
			}
		}
	}
}

func drawVideoTextMode() {
	i := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			drawChar(videoTextMemory[i], c*8, r*8)
			i++
		}
	}
}

func clearVideoTextMode() {
	for i := 0; i < rows*columns; i++ {
		videoTextMemory[i] = 0
	}
}

func putChar(c byte) {
	videoTextMemory[cursor] = c
	cursor++
	if cursor >= rows*columns {
		// todo:
		// move chars 1 row up
		// subtract one row fron cursor
		cursor = 0
	}
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
	//drawChar(0, 100, 100)
	//drawChar(1, 100+8, 100)
	//drawChar(2, 100+8+8, 100)

	putChar(1)
	putChar(0)

	drawVideoTextMode()
	screen.ReplacePixels(img.Pix)

	//block(screen)
	/*
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
	*/
	return nil
}

func main() {

	i := 0
	c := byte(0)
	tot := rows * columns
	for {
		videoTextMemory[i] = c
		c++
		if c > 3 {
			c = 0
		}
		i++
		if i >= tot {
			break
		}
	}

	font.Load()

	img = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "METAL BASIC 0.01"); err != nil {
		log.Fatal(err)
	}

	if err := ebiten.Run(update, 320, 240, 2, "Metal"); err != nil {
		panic(err)
	}

}
