package main

import (
	"image"
	"math/rand"
	"strings"
	"time"

	"github.com/crgimenes/graphos/coreScreen"
	"github.com/crgimenes/metal/cmd"
	"github.com/crgimenes/metal/fonts"
	"github.com/hajimehoshi/ebiten"
)

const (
	rows     = 30
	columns  = 40
	rgbaSize = 4
)

var (
	videoTextMemory  [rows * columns * 2]byte
	cursor           int
	img              *image.RGBA
	font             fonts.Expert118x8
	currentColor     byte = 0x9f
	updateScreen     bool
	cpx, cpy         int
	cursorBlinkTimer int
	cursorSetBlink   bool = true

	uTime uint64

	machine int

	//var countaux int
	noKey bool
	shift bool

	cs *coreScreen.Instance
)

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

func input() {
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
			line := cursor / (columns * 2)
			lineEnd := line*columns*2 + columns*2
			if cursor < 0 {
				cursor = 0
			}

			copy(videoTextMemory[cursor:lineEnd], videoTextMemory[cursor+2:lineEnd])
			videoTextMemory[lineEnd-2] = currentColor
			videoTextMemory[lineEnd-1] = 0

			correctVideoCursor()
		})
		return
	}

	/*
	   KeyMinus: -
	   KeyEqual: =
	   KeyLeftBracket: [
	   KeyRightBracket: ]
	   KeyBackslash:
	   KeySemicolon: ;
	   KeyApostrophe: '
	   KeySlash: /
	   KeyGraveAccent: `
	*/

	shift = ebiten.IsKeyPressed(ebiten.KeyShift)

	if ebiten.IsKeyPressed(ebiten.KeyEqual) {
		if shift {
			keyTreatment('+', func(c byte) {
				putChar(c)
				println("+")
			})
			return
		} else {
			keyTreatment('=', func(c byte) {
				putChar(c)
				println("=")
			})
			return
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		keyTreatment(0, func(c byte) {
			cursor -= columns * 2
			correctVideoCursor()
		})
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		keyTreatment(0, func(c byte) {
			cursor += columns * 2
			correctVideoCursor()
		})
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		keyTreatment(0, func(c byte) {
			cursor -= 2
			correctVideoCursor()
		})
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
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

	cpx, cpy = ebiten.CursorPosition()
	cpx -= cs.Border
	cpy -= cs.Border
	//fmt.Printf("X: %d, Y: %d\n", x, y)

	// Display the information with "X: xx, Y: xx" format
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

	noKey = true

}

func drawChar(index, fgColor, bgColor byte, x, y int) {
	var a, b uint64
	for a = 0; a < 8; a++ {
		for b = 0; b < 8; b++ {
			if font.Bitmap[index][b]&(0x80>>a) != 0 {
				cs.CurrentColor = fgColor
				cs.DrawPix(int(a)+x, int(b)+y)
			} else {
				cs.CurrentColor = bgColor
				cs.DrawPix(int(a)+x, int(b)+y)
			}
		}
	}
}

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

func update(screen *coreScreen.Instance) error {

	uTime++
	//putChar(2)
	//cursor -= 2

	if machine == 0 {
		bPrintln("METAL BASIC 0.01")
		bPrintln("http://crg.eti.br")
		machine++
	}
	//machine++

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

	//bCircle(100, 100, 10)
	//for i := 0; i < len(a); i++ {
	//	bFilledCircle(a[i].X, a[i].Y, 5)

	//	a[i].X += random(-1, +2)
	//	a[i].Y += random(-1, +2)
	//}

	//bLine(100, 100, cpx, cpy)

	//bLine(50, 50, 50, 100)
	//bLine(50, 100, 100, 100)
	//bLine(50, 50, 100, 50)
	//bLine(100, 50, 100, 100)

	//bLine(50, 50, 100, 100)
	//bLine(100, 50, 94, 44)
	//	bBox(50, 50, 100, 100)

	input()
	return nil
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().Unix())
	font.Load()
	clearVideoTextMode()

	cs = coreScreen.Get()

	cs.Border = 10
	cs.Width = 320 + cs.Border*2  // 40 columns
	cs.Height = 240 + cs.Border*2 // 30 rows
	cs.Update = update
	cs.Title = "Metal BASIC 0.01"

	cs.Run()

}
