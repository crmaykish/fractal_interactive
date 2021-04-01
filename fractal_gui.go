package main

import (
	"fmt"
	"math"
	"os"

	"github.com/crmaykish/fractal_core"
	"github.com/veandco/go-sdl2/sdl"
)

const width = 1000
const height = 1000

var iterations = 2000
var iterationGap = 1000

var lightColor = uint32(0xFFFFFF)
var darkColor = uint32(0x000000)

func newColor(colorA, colorB uint32, hue float64) (uint8, uint8, uint8) {
	var ra = uint8(colorA & (0xFF << 16) >> 16)
	var ga = uint8(colorA & (0xFF << 8) >> 8)
	var ba = uint8(colorA & 0xFF)
	var rb = uint8(colorB & (0xFF << 16) >> 16)
	var gb = uint8(colorB & (0xFF << 8) >> 8)
	var bb = uint8(colorB & 0xFF)

	return uint8(float64(rb-ra)*hue) + ra, uint8(float64(gb-ga)*hue) + ga, uint8(float64(bb-ba)*hue) + ba
}

func drawMandelbrot(renderer *sdl.Renderer, m *fractal_core.Mandelbrot) {
	renderer.Clear()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			renderer.SetDrawColor(0, 0, 0, 0xFF)

			var val = fractal_core.GetBuffer(m)[x][y]
			var hue = fractal_core.GetHue(m)[x][y]

			if val < uint32(iterations) {
				var r, g, b = newColor(darkColor, lightColor, hue)
				renderer.SetDrawColor(r, g, b, 0x00)
			}

			renderer.DrawPoint(int32(x), int32(y))
		}
	}
}

func main() {
	fmt.Println("Starting")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	fmt.Println("Creating window")
	window, err := sdl.CreateWindow("Fractal Interactive", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
	}
	defer renderer.Destroy()

	// Create the Mandelbrot generator
	m := fractal_core.Create(width, height, complex(-0.5, 0))

	fractal_core.SetMaxIterations(m, iterations)
	fractal_core.Generate(m)

	drawMandelbrot(renderer, m)

	// window.UpdateSurface()
	renderer.Present()

	running := true

	var regen = false

	var clickx, clicky int32

	// Main loop
	for running {

		// Event polling
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.MouseButtonEvent:
				if event.GetType() == sdl.MOUSEBUTTONDOWN {
					clickx, clicky, _ = sdl.GetMouseState()
					regen = true
				}
			}
		}

		if regen {
			minx, miny, maxx, maxy := fractal_core.GetBounds(m)

			var realX = fractal_core.MapIntToFloat(int(clickx), 0, width, minx, maxx)
			var realY = fractal_core.MapIntToFloat(int(clicky), 0, height, miny, maxy)
			var op = "+"

			if realY < 0 {
				op = "-"
			}

			fractal_core.SetCenter(m, complex(realX, realY))

			fractal_core.ScaleZoom(m, 5.0)

			iterations += iterationGap

			fractal_core.SetMaxIterations(m, iterations)

			fmt.Println("Regenerating...")
			fractal_core.Generate(m)

			drawMandelbrot(renderer, m)

			window.UpdateSurface()

			fmt.Println("Done!")
			fmt.Printf("(%v %s %vi) at %vx\n", realX, op, math.Abs(realY), fractal_core.GetZoom(m))

			regen = false
		}

		renderer.Present()
		sdl.Delay(16)

	}
}
