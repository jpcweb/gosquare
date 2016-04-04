package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math"
	"runtime"
	"time"
)

func init() {
	runtime.LockOSThread()
}

type interactCoord struct {
	XClic int32
	YClic int32
	XCoef float64
	YCoef float64
}

type Coord struct {
	X []int32
	Y []int32
}

func drawAndWait(renderer *sdl.Renderer, rect sdl.Rect) {
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.FillRect(&rect)
	renderer.Present()

	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()
	time.Sleep(time.Millisecond * 1)
}

func calc(speed float64, coef float64, ratio float64) int32 {
	return int32((speed * coef) * math.Sin(ratio))
}

func defineCoef(diff int32) float64 {
	var coef float64
	if diff > 0 {
		coef = 1
	} else {
		coef = -1
	}
	return coef
}

/*
	pos > rect position (x or y)
	clicPos > mouse position (x or y)
*/
func getPosition(pos int32, clicPos int32, speed float64, coef float64, ratio float64) int32 {
	if (coef == 1 && pos > clicPos) || (coef == -1 && pos < clicPos) {
		pos = clicPos
	} else {
		pos += calc(speed, coef, ratio)
	}
	return pos
}

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	var renderer *sdl.Renderer
	renderer, err = sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		log.Panic(err)
	}
	defer renderer.Destroy()

	rect := sdl.Rect{0, 0, 200, 200}

	var ev sdl.Event
	var ic interactCoord
	running := true
	var x int32
	var y int32
	var action bool = false
	var speed float64 = 24
	var ratio float64 = 0.1
	for running {
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&rect)
		renderer.Present()

		ev = sdl.PollEvent()

		if ev != nil {
			switch t := ev.(type) { //Event is an empty interface so we need to use switch to match type
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				if t.Type == sdl.MOUSEBUTTONDOWN {
					action = false
				}
				if t.Type == sdl.MOUSEBUTTONUP {
					var coord Coord
					var breaker string
					/*You can even prepare a move array*/
					action = true
					ic.XClic = t.X - (rect.W / 2)
					ic.YClic = t.Y - (rect.W / 2)

					ic.XCoef = defineCoef(ic.XClic - x)
					ic.YCoef = defineCoef(ic.YClic - y)

					breaker = "y"
					if ((ic.XClic-rect.X)*int32(ic.XCoef))-((ic.YClic-rect.Y)*int32(ic.YCoef)) > 0 {
						breaker = "x"
					}
					for {
						if !action {
							break
						}

						x = getPosition(x, ic.XClic, speed, ic.XCoef, ratio)
						if x == ic.XClic && breaker == "x" {
							break
						}
						coord.X = append(coord.X, x)

						y = getPosition(y, ic.YClic, speed, ic.YCoef, ratio)
						if y == ic.YClic && breaker == "y" {
							break
						}
						coord.Y = append(coord.Y, y)
					}
					/**/
					for k, _ := range coord.X {
						rect.X = coord.X[k]
						if len(coord.Y) > k {
							rect.Y = coord.Y[k]
						}
						drawAndWait(renderer, rect)
					}
				}
			default: //Locked escape error
				continue
			}
		}
		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()
	}

	//sdl.Delay(1000)
	//sdl.Quit()
}
