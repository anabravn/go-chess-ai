package main

import (
	"ai"
	"ui"

	"fmt"
	"os"
	"strconv"

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
	var d int
	var color chess.Color

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Usagem: ./chess-ai [profundidade] [cor]")
		os.Exit(1)
	} else {
		depth, err := strconv.Atoi(args[0])
		d = depth * 2

		if err != nil {
			fmt.Printf("\"%s\" não é uma profundidade válida\n", args[0])
			os.Exit(1)
		}

		switch args[1] {
		case "b", "black":
			color = chess.Black
		default:
			color = chess.White
		}
	}

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()
	img.Init(img.INIT_PNG)

	window := CreateWindow()
	defer window.Destroy()
	surface, _ := window.GetSurface()

	gui := ui.NewGameUI()
	game := chess.NewGame()
	gui.Update(game)

	var nextMove *chess.Move = nil

	over := false
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

		if game.Position().Turn() == color {
			nextMove = gui.GetSelectedMove(game)
		} else {
			nextMove = ai.Search(game, d)
		}

		if nextMove != nil {
			game.Move(nextMove)
			gui.Update(game)
			nextMove = nil
		}

		if !over && game.Outcome() != chess.NoOutcome {
			fmt.Println(game.Outcome(), game.Method())
			over = true
		}

		gui.Draw(surface)
		window.UpdateSurface()
	}
}
