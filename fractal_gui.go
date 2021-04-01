package main

import (
	"fmt"
	"math"
	"os"

	"github.com/crmaykish/fractal_core"
	"github.com/veandco/go-sdl2/sdl"
)

const width = 500
const height = 500

var iterations = 1000

func drawMandelbrot(renderer *sdl.Renderer, m *fractal_core.Mandelbrot) {
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var color = uint8(0)

			var cellValue = fractal_core.GetBuffer(m)[x][y]

			if cellValue < uint32(iterations) {
				color = uint8(fractal_core.MapIntToFloat(int(cellValue), 0, iterations, 0, math.Pow(2, 8)-1))
			}

			renderer.SetDrawColor(color, color, color, 0xFF)
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
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	renderer.Clear()

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
					fmt.Printf("click at %d, %d\n", clickx, clicky)
					regen = true
				}
			}
		}

		if regen {
			minx, miny, maxx, maxy := fractal_core.GetBounds(m)

			var realX = fractal_core.MapIntToFloat(int(clickx), 0, width, minx, maxx)
			var realY = fractal_core.MapIntToFloat(int(clicky), 0, height, miny, maxy)
			fmt.Printf("(%f + %fi)\n", realX, realY)

			fractal_core.SetCenter(m, complex(realX, realY))

			fractal_core.ScaleZoom(m, 10.0)

			iterations += 1000

			fractal_core.SetMaxIterations(m, iterations)

			fmt.Println("Regenerating...")
			fractal_core.Generate(m)

			drawMandelbrot(renderer, m)

			window.UpdateSurface()

			fmt.Println("Done!")

			regen = false
		}

		renderer.Present()
		sdl.Delay(16)

	}
}
