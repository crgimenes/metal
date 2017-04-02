package main

import (
	"image"
	"log"

	"github.com/crgimenes/metal/pcloader/fonts"
	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 320 // 40 columns
	screenHeight = 240 // 30 rows

	rows     = 30
	columns  = 40
	rgbaSize = 4
)

var (
	videoTextMemory [rows * columns * 2]byte
	cursor          int
	img             *image.RGBA
	square          *ebiten.Image
	font            fonts.Expert118x8
)

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
			// i == color code
			i++
			drawChar(videoTextMemory[i], c*8, r*8)
			i++
		}
	}
}

func clearVideoTextMode() {
	copy(videoTextMemory[:], make([]byte, len(videoTextMemory)))
}

func moveLineUp() {
	copy(videoTextMemory[0:], videoTextMemory[columns*2:])
	copy(videoTextMemory[len(videoTextMemory)-columns*2:], make([]byte, columns*2))
}

func putChar(c byte) {
	// videoTextMemory[cursor] // color code
	cursor++
	videoTextMemory[cursor] = c
	cursor++
	if cursor >= rows*columns*2 {
		// todo:
		// move chars 1 row up
		// subtract one row fron cursor
		cursor -= columns * 2
		moveLineUp()
	}

}

var dt byte

func update(screen *ebiten.Image) error {

	putChar(dt)
	dt++

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

	font.Load()
	//clearVideoTextMode()

	img = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "METAL BASIC 0.01"); err != nil {
		log.Fatal(err)
	}

}
