package main

import (
	"image"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/crgimenes/metal/cmd"
	"github.com/crgimenes/metal/fonts"
	"github.com/hajimehoshi/ebiten"
)

const (
	border       = 10
	screenWidth  = 320 + border*2 // 40 columns
	screenHeight = 240 + border*2 // 30 rows

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
	currentColor    byte = 0x9f
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
	x += border
	y += border
	if x < border || y < border || x >= screenWidth-border || y >= screenHeight-border {
		return
	}
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
var shift bool

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

var cpx, cpy int

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
	cpx -= border
	cpy -= border
	//fmt.Printf("X: %d, Y: %d\n", x, y)

	// Display the information with "X: xx, Y: xx" format
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

	noKey = true

}

func clearVideo() {
	for i := 0; i < screenHeight*screenWidth*4; i += 4 {
		img.Pix[i] = CGAColors[9].R
		img.Pix[i+1] = CGAColors[9].G
		img.Pix[i+2] = CGAColors[9].B
		img.Pix[i+3] = 0xff
	}
}

/*
func bLine(x1, y1, x2, y2 int) {
	dx := x2 - x1
	dy := y2 - y1
	for x := x1; x < x2; x++ {
		y := y1 + dy*(x-x1)/dx
		drawPix(x, y, 0xf)
	}
}
*/

func bLine(x1, y1, x2, y2 int) {
	var x, y, dx, dy, dx1, dy1, px, py, xe, ye, i int
	dx = x2 - x1
	dy = y2 - y1
	if dx < 0 {
		dx1 = -dx
	} else {
		dx1 = dx
	}

	if dy < 0 {
		dy1 = -dy
	} else {
		dy1 = dy
	}
	px = 2*dy1 - dx1
	py = 2*dx1 - dy1
	if dy1 <= dx1 {
		if dx >= 0 {
			x = x1
			y = y1
			xe = x2
		} else {
			x = x2
			y = y2
			xe = x1
		}
		drawPix(x, y, 0xf)
		for i = 0; x < xe; i++ {
			x = x + 1
			if px < 0 {
				px = px + 2*dy1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					y = y + 1
				} else {
					y = y - 1
				}
				px = px + 2*(dy1-dx1)
			}
			drawPix(x, y, 0xf)
		}
	} else {
		if dy >= 0 {
			x = x1
			y = y1
			ye = y2
		} else {
			x = x2
			y = y2
			ye = y1
		}
		drawPix(x, y, 0xf)
		for i = 0; y < ye; i++ {
			y = y + 1
			if py <= 0 {
				py = py + 2*dx1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					x = x + 1
				} else {
					x = x - 1
				}
				py = py + 2*(dx1-dy1)
			}
			drawPix(x, y, 0xf)
		}
	}
}

func bBox(x1, y1, x2, y2 int) {
	for y := y1; y <= y2; y++ {
		drawPix(x1, y, 0xf)
		drawPix(x2, y, 0xf)
	}
	for x := x1; x <= x2; x++ {
		drawPix(x, y1, 0xf)
		drawPix(x, y2, 0xf)
	}
}

func bCircle(x0, y0, radius int) {
	x := radius
	y := 0
	e := 0

	for x >= y {
		drawPix(x0+x, y0+y, 0xf)
		drawPix(x0+y, y0+x, 0xf)
		drawPix(x0-y, y0+x, 0xf)
		drawPix(x0-x, y0+y, 0xf)
		drawPix(x0-x, y0-y, 0xf)
		drawPix(x0-y, y0-x, 0xf)
		drawPix(x0+y, y0-x, 0xf)
		drawPix(x0+x, y0-y, 0xf)

		if e <= 0 {
			y += 1
			e += 2*y + 1
		}
		if e > 0 {
			x -= 1
			e -= 2*x + 1
		}
	}
}

func bFilledCircle(x0, y0, radius int) {
	x := radius
	y := 0
	xChange := 1 - (radius << 1)
	yChange := 0
	radiusError := 0

	for x >= y {
		for i := x0 - x; i <= x0+x; i++ {
			drawPix(i, y0+y, 0xf)
			drawPix(i, y0-y, 0xf)
		}
		for i := x0 - y; i <= x0+y; i++ {
			drawPix(i, y0+x, 0xf)
			drawPix(i, y0-x, 0xf)
		}

		y++
		radiusError += yChange
		yChange += 2
		if ((radiusError << 1) + xChange) > 0 {
			x--
			radiusError += xChange
			xChange += 2
		}
	}
}

var ranking = 1000000

func update(screen *ebiten.Image) error {

	uTime++
	//putChar(2)
	//cursor -= 2

	if machine == 0 {
		bPrintln("METAL BASIC 0.01")
		bPrintln("http://crg.eti.br")
		machine++
	}
	machine++

	if machine > 10 {
		return nil
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

	//bCircle(100, 100, 10)
	for i := 0; i < len(a); i++ {
		bFilledCircle(a[i].X, a[i].Y, 5)

		a[i].X += random(-1, +2)
		a[i].Y += random(-1, +2)
	}

	bLine(100, 100, cpx, cpy)

	//bLine(50, 50, 50, 100)
	//bLine(50, 100, 100, 100)
	//bLine(50, 50, 100, 50)
	//bLine(100, 50, 100, 100)

	//bLine(50, 50, 100, 100)
	//bLine(100, 50, 94, 44)
	//	bBox(50, 50, 100, 100)

	screen.ReplacePixels(img.Pix)
	keyboard()
	return nil
}

func distance(x1, y1, x2, y2 int) int {
	first := math.Pow(float64(x2-x1), 2)
	second := math.Pow(float64(y2-y1), 2)
	return int(math.Sqrt(first + second))
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

type dot struct {
	X int
	Y int
	B bool
}

var a []dot

func main() {
	rand.Seed(time.Now().Unix())
	font.Load()
	clearVideoTextMode()

	for i := 0; i < 10; i++ {
		x, y := random(0, screenWidth), random(0, screenHeight)
		d := dot{X: x, Y: y}
		a = append(a, d)
	}

	img = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	clearVideo()
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "METAL BASIC 0.01"); err != nil {
		log.Fatal(err)
	}

}
