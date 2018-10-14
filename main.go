package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gonutz/prototype/draw"
)

const tileSize = 60

/*
Example field, S are my ships, o means miss, x means hit, left field is mine,
right is hers.

 1234567890 1234567890
A       S  A          A
B  S    S  B          B
C  S    S  C   o  x   C
D  S    S  D      xo  D
E  S       E  o       E
F  S       F          F
G          G   o  o   G
H     SS   H  oxxx    H
I          I          I
J SSS  SSS J          J
 1234567890 1234567890
*/

func main() {
	rand.Seed(time.Now().UnixNano())

	const (
		stateSettingShips = iota
		stateShooting
	)
	state := stateSettingShips
	var ships []ship
	var home, sea seaField

	restart := func() {
		ships = []ship{
			newShip(5),
			newShip(4),
			newShip(3),
			newShip(3),
			newShip(2),
		}
		state = stateSettingShips
		home = seaField{}
		sea = seaField{}
	}
	restart()

	const windowW, windowH = 23 * tileSize, 12 * tileSize
	check(draw.RunWindow("Ships", windowW, windowH, func(window draw.Window) {
		if window.WasKeyPressed(draw.KeyEscape) {
			window.Close()
		}
		if window.WasKeyPressed(draw.KeyF2) {
			restart()
		}

		if state == stateSettingShips {
			setI := -1
			for i := range ships {
				if !ships[i].isSet {
					setI = i
					break
				}
			}
			if setI == -1 {
				state = stateShooting
			} else {
				s := &ships[setI]
				s.moving = true
				for _, c := range window.Clicks() {
					if c.Button == draw.RightButton {
						s.vert = !s.vert
						break
					}
				}
				mx, my := window.MousePosition()
				tileX := (mx - tileSize) / tileSize
				tileY := (my - tileSize) / tileSize
				if 0 <= tileX && tileX < 10 && 0 <= tileY && tileY < 10 {
					if s.vert {
						s.x, s.y = tileX, tileY-s.length/2
						if s.y < 0 {
							s.y = 0
						}
						if s.y+s.length > 10 {
							s.y = 10 - s.length
						}
					} else {
						s.x, s.y = tileX-s.length/2, tileY
						if s.x < 0 {
							s.x = 0
						}
						if s.x+s.length > 10 {
							s.x = 10 - s.length
						}
					}
				}

				for _, c := range window.Clicks() {
					if c.Button == draw.LeftButton {
						s.isSet = true
						s.moving = false
						break
					}
				}
			}
		} else if state == stateShooting {
			mx, my := window.MousePosition()
			tileX := (mx - 12*tileSize) / tileSize
			tileY := (my - tileSize) / tileSize
			if 0 <= tileX && tileX < 10 && 0 <= tileY && tileY < 10 {
				for _, c := range window.Clicks() {
					if c.Button == draw.RightButton {
						if sea[tileX][tileY] == seaHit {
							sea[tileX][tileY] = seaEmpty
						} else {
							sea[tileX][tileY] = seaHit
						}
						break
					} else if c.Button == draw.LeftButton {
						if sea[tileX][tileY] == seaMiss {
							sea[tileX][tileY] = seaEmpty
						} else {
							sea[tileX][tileY] = seaMiss
						}
						break
					}
				}
			} else {
				tileX := (mx - tileSize) / tileSize
				tileY := (my - tileSize) / tileSize
				if 0 <= tileX && tileX < 10 && 0 <= tileY && tileY < 10 {
					if len(window.Clicks()) > 0 {
						if home[tileX][tileY] != seaEmpty {
							home[tileX][tileY] = seaEmpty
						} else {
							home[tileX][tileY] = seaHit
						}
					}
				}
			}
		}

		const textScale = 2
		textColor := draw.White
		charW, charH := window.GetScaledTextSize("A", textScale)
		textDx, textDy := (tileSize-charW)/2, (tileSize-charH)/2
		for i := 1; i <= 10; i++ {
			n := strconv.Itoa(i)
			window.DrawScaledText(n, i*tileSize+textDx, textDy, textScale, textColor)
			window.DrawScaledText(n, (i+11)*tileSize+textDx, textDy, textScale, textColor)
			window.DrawScaledText(n, i*tileSize+textDx, 11*tileSize+textDy, textScale, textColor)
			window.DrawScaledText(n, (i+11)*tileSize+textDx, 11*tileSize+textDy, textScale, textColor)
		}
		for c := 'A'; c <= 'J'; c++ {
			y := int(c-'A'+1)*tileSize + textDy
			window.DrawScaledText(string(c), textDx, y, textScale, textColor)
			window.DrawScaledText(string(c), 11*tileSize+textDx, y, textScale, textColor)
			window.DrawScaledText(string(c), 22*tileSize+textDx, y, textScale, textColor)
		}
		seaColor := draw.Blue
		window.FillRect(tileSize, tileSize, 10*tileSize, 10*tileSize, seaColor)
		window.FillRect(12*tileSize, tileSize, 10*tileSize, 10*tileSize, seaColor)

		margin := tileSize / 4
		{
			col := draw.White
			col.A = 0.1
			for y := 0; y < 10; y++ {
				for x := 0; x < 10; x++ {
					window.FillEllipse(
						(x+1)*tileSize+margin,
						(y+1)*tileSize+margin,
						tileSize-2*margin,
						tileSize-2*margin,
						col,
					)

					seaCol := col
					if sea[x][y] == seaMiss {
						seaCol = draw.White
					}
					if sea[x][y] == seaHit {
						seaCol = draw.Red
					}
					window.FillEllipse(
						(x+12)*tileSize+margin,
						(y+1)*tileSize+margin,
						tileSize-2*margin,
						tileSize-2*margin,
						seaCol,
					)
				}
			}
		}

		for _, s := range ships {
			if s.isSet || s.moving {
				startX, startY := s.x, s.y
				endX, endY := startX, startY
				if s.vert {
					endY += s.length - 1
				} else {
					endX += s.length - 1
				}
				for y := startY; y <= endY; y++ {
					for x := startX; x <= endX; x++ {
						window.FillRect(
							(x+1)*tileSize,
							(y+1)*tileSize,
							tileSize,
							tileSize,
							draw.Gray,
						)
					}
				}
			}
		}

		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				if home[x][y] != seaEmpty {
					window.FillEllipse(
						(x+1)*tileSize+margin,
						(y+1)*tileSize+margin,
						tileSize-2*margin,
						tileSize-2*margin,
						draw.Black,
					)
				}
			}
		}
	}))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func newShip(length int) ship {
	return ship{
		length: length,
		vert:   rand.Intn(2) == 0,
	}
}

type ship struct {
	length int
	vert   bool
	x, y   int
	isSet  bool
	moving bool
}

type seaField [10][10]int

const (
	seaEmpty = iota
	seaMiss
	seaHit
)
