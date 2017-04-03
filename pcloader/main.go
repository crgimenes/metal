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
	currentColor    byte = 0x0f
)

var CGAColors = []struct {
	R byte
	G byte
	B byte
}{
	{0, 0, 0},
	{0, 0, 170},
	{0, 170, 0},
	{0, 170, 170},
	{170, 0, 0},
	{170, 0, 170},
	{170, 85, 0},
	{170, 170, 170},
	{85, 85, 85},
	{85, 85, 255},
	{85, 255, 85},
	{85, 255, 255},
	{255, 85, 85},
	{255, 85, 255},
	{255, 255, 85},
	{255, 255, 255},
}

func mergeColorCode(b, f byte) byte {
	return (f & 0xff) | (b << 4)
}

func drawPix(x, y int, color byte) {
	pos := 4*y*screenWidth + 4*x
	img.Pix[pos] = CGAColors[color].R
	img.Pix[pos+1] = CGAColors[color].G
	img.Pix[pos+2] = CGAColors[color].B
	img.Pix[pos+3] = 0xff
}

func getBit(n int, pos uint64) bool {
	// from right to left
	val := n & (1 << pos)
	return (val > 0)
}

func drawChar(index, fgColor, bgColor byte, x, y int) {
	var a, b uint64
	for a = 0; a < 8; a++ {
		for b = 0; b < 8; b++ {
			if font.Bitmap[index][b]&(0x80>>a) != 0 {
				drawPix(int(a)+x, int(b)+y, fgColor)
			} else {
				//drawOffPix(int(a)+x, int(b)+y)
				drawPix(int(a)+x, int(b)+y, bgColor)
			}
		}
	}
}

func drawVideoTextMode() {
	i := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			color := videoTextMemory[i]
			f := color & 0x0f
			b := color & 0xf0 >> 4
			i++
			drawChar(videoTextMemory[i], f, b, c*8, r*8)
			i++
		}
	}
}

func clearVideoTextMode() {
	copy(videoTextMemory[:], make([]byte, len(videoTextMemory)))
	for i := 0; i < len(videoTextMemory); i += 2 {
		videoTextMemory[i] = 0x0F
	}
}

func moveLineUp() {
	copy(videoTextMemory[0:], videoTextMemory[columns*2:])
	copy(videoTextMemory[len(videoTextMemory)-columns*2:], make([]byte, columns*2))
	for i := len(videoTextMemory) - columns*2; i < len(videoTextMemory); i += 2 {
		videoTextMemory[i] = 0x0F
	}

}

func putChar(c byte) {
	videoTextMemory[cursor] = currentColor
	cursor++
	videoTextMemory[cursor] = c
	cursor++
	if cursor >= rows*columns*2 {
		cursor -= columns * 2
		moveLineUp()
	}
}

func bPrint(msg string) {
	for i := 0; i < len(msg); i++ {
		c := msg[i]

		switch c {
		case 13:
			cursor += columns * 2
			continue
		case 10:
			aux := cursor / (columns * 2)
			aux = aux * (columns * 2)
			cursor = aux
			continue
		}
		putChar(msg[i])
	}
}

func bPrintln(msg string) {
	msg += "\r\n"
	bPrint(msg)
}

//var dt byte
//var c byte

var machine int

func update(screen *ebiten.Image) error {

	if machine == 0 {
		bPrintln("teste0")
		bPrintln("teste1")
		bPrintln("teste2")
		machine++
	}

	//putChar(dt)
	//dt++
	//currentColor = mergeColorCode(0x0, c)
	//c++
	//if c > 15 {
	//	c = 0
	//}

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
	clearVideoTextMode()

	img = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "METAL BASIC 0.01"); err != nil {
		log.Fatal(err)
	}

}
