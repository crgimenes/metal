package main

import (
	"image"
	"log"
	"strings"

	"github.com/crgimenes/metal/cmd"
	"github.com/crgimenes/metal/fonts"
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
				drawPix(int(a)+x, int(b)+y, bgColor)
			}
		}
	}
}

var cursorBlinkTimer int
var cursorSetBlink bool = true

func drawCursor(index, fgColor, bgColor byte, x, y int) {
	if cursorSetBlink {
		if cursorBlinkTimer < 15 {
			drawChar(index, fgColor, bgColor, x, y)
		} else {
			drawChar(index, bgColor, fgColor, x, y)
		}
		cursorBlinkTimer++
		if cursorBlinkTimer > 30 {
			cursorBlinkTimer = 0
		}
	} else {
		drawChar(index, bgColor, fgColor, x, y)
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
			if i-1 == cursor {
				drawCursor(videoTextMemory[i], f, b, c*8, r*8)
			} else {
				drawChar(videoTextMemory[i], f, b, c*8, r*8)
			}
			i++
		}
	}
}

func clearVideoTextMode() {
	copy(videoTextMemory[:], make([]byte, len(videoTextMemory)))
	for i := 0; i < len(videoTextMemory); i += 2 {
		videoTextMemory[i] = currentColor
	}
}

func moveLineUp() {
	copy(videoTextMemory[0:], videoTextMemory[columns*2:])
	copy(videoTextMemory[len(videoTextMemory)-columns*2:], make([]byte, columns*2))
	for i := len(videoTextMemory) - columns*2; i < len(videoTextMemory); i += 2 {
		videoTextMemory[i] = currentColor
	}

}

func correctVideoCursor() {
	if cursor < 0 {
		cursor = 0
	}
	for cursor >= rows*columns*2 {
		cursor -= columns * 2
		moveLineUp()
	}
}

func putChar(c byte) {
	correctVideoCursor()
	videoTextMemory[cursor] = currentColor
	cursor++
	correctVideoCursor()
	videoTextMemory[cursor] = c
	cursor++
	correctVideoCursor()
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

var lastKey = struct {
	Time uint64
	Char byte
}{
	0,
	0,
}

var uTime uint64
var c byte

var machine int

//var countaux int
var noKey bool

func keyTreatment(c byte, f func(c byte)) {
	if noKey || lastKey.Char != c || lastKey.Time+20 < uTime {
		f(c)
		noKey = false
		lastKey.Char = c
		lastKey.Time = uTime
	}
}

func getLine() string {
	aux := cursor / (columns * 2)
	var ret string
	for i := aux*(columns*2) + 1; i < aux*(columns*2)+columns*2; i += 2 {
		c := videoTextMemory[i]
		if c == 0 {
			break
		}
		ret += string(videoTextMemory[i])
	}

	ret = strings.TrimSpace(ret)
	return ret
}

func keyboard() {
	for c := 'A'; c <= 'Z'; c++ {
		if ebiten.IsKeyPressed(ebiten.Key(c) - 'A' + ebiten.KeyA) {
			keyTreatment(byte(c), func(c byte) {
				putChar(c)
			})
			return
		}
	}

	for c := '0'; c <= '9'; c++ {
		if ebiten.IsKeyPressed(ebiten.Key(c) - '0' + ebiten.Key0) {
			keyTreatment(byte(c), func(c byte) {
				putChar(c)
			})
			return
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		keyTreatment(byte(' '), func(c byte) {
			putChar(c)
		})
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyComma) {
		keyTreatment(byte(','), func(c byte) {
			putChar(c)
		})
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		keyTreatment(0, func(c byte) {
			cmd.Eval(getLine())
			cursor += columns * 2
			aux := cursor / (columns * 2)
			aux = aux * (columns * 2)
			cursor = aux
			correctVideoCursor()
		})
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		keyTreatment(0, func(c byte) {
			cursor -= 2
			correctVideoCursor()
		})
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		keyTreatment(0, func(c byte) {
			cursor -= columns * 2
			correctVideoCursor()
		})
		return
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		keyTreatment(0, func(c byte) {
			cursor += columns * 2
			correctVideoCursor()
		})
		return
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		keyTreatment(0, func(c byte) {
			cursor -= 2
			correctVideoCursor()
		})
		return
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		keyTreatment(0, func(c byte) {
			cursor += 2
			correctVideoCursor()
		})
		return
	}

	// When the "left mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		//ebitenutil.DebugPrint(screen, "You're pressing the 'LEFT' mouse button.")
	}
	// When the "right mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		//ebitenutil.DebugPrint(screen, "\nYou're pressing the 'RIGHT' mouse button.")
	}
	// When the "middle mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		//ebitenutil.DebugPrint(screen, "\n\nYou're pressing the 'MIDDLE' mouse button.")
	}

	//x, y := ebiten.CursorPosition()
	//fmt.Printf("X: %d, Y: %d\n", x, y)

	// Display the information with "X: xx, Y: xx" format
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

	noKey = true

}

func update(screen *ebiten.Image) error {

	uTime++
	//putChar(2)
	//cursor -= 2

	if machine == 0 {
		bPrintln("METAL BASIC 0.01")
		bPrintln("http://crg.eti.br")
		machine++
	}

	/*
		if countaux > 10 {
			countaux = 0
			putChar(dt)
			dt++
			currentColor = mergeColorCode(0x0, c)
			c++
			if c > 15 {
				c = 0
			}
		}
		countaux++
	*/
	drawVideoTextMode()
	screen.ReplacePixels(img.Pix)
	keyboard()
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
