package main

import (
	"ui"
	"ai"
	"fmt"

	"github.com/notnil/chess"
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

	gameUi := ui.NewGameUI()
	game := chess.NewGame()
	gameUi.Update(game)

	var nextMove *chess.Move = nil

	running := true
	for running {
		event := sdl.WaitEventTimeout(10)
		if event != nil {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			}
		}

		if game.Position().Turn() == chess.White {
		    nextMove = gameUi.GetSelectedMove(game)
		} else {
		    nextMove = ai.Search(game, 2)
		}

		if nextMove != nil {
		    game.Move(nextMove)
		    gameUi.Update(game)
		    nextMove = nil
		}

		if game.Outcome() != chess.NoOutcome {
		    fmt.Println(game.Outcome(), game.Method())
		    running = false
		}

		gameUi.Draw(surface)
		window.UpdateSurface()
	}
}
