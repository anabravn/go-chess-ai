package main

import (
	"ui"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func CreateWindow() *sdl.Window {
	window, _ := sdl.CreateWindow("Go Chess AI",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		ui.ScreenWidth, ui.ScreenHeight, sdl.WINDOW_SHOWN)
	surface, _ := window.GetSurface()
	surface.FillRect(nil, 0)

	return window
}

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()
	img.Init(img.INIT_PNG)

	window := CreateWindow()
	defer window.Destroy()
	surface, _ := window.GetSurface()

	game := ui.NewGame()

	running := true
	for running {
		window.UpdateSurface()

		event := sdl.WaitEventTimeout(10)
		if event != nil {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			}
		}

		running = game.Update()
		game.Draw(surface)
	}
}
