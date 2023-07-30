package main

import (
	"github.com/notnil/chess"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()
	img.Init(img.INIT_PNG)

	window, _ := sdl.CreateWindow("Go Chess AI",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		ScreenWidth, ScreenHeight, sdl.WINDOW_SHOWN)
	defer window.Destroy()
	surface, _ := window.GetSurface()
	surface.FillRect(nil, 0)

	board := NewBoardUI()
	game := chess.NewGame()
	var nextMove *chess.Move

	squareMap := game.Position().Board().SquareMap()
	moves := game.ValidMoves()
	board.Update(squareMap)

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

		if game.Position().Turn() == chess.White {
			nextMove = board.GetSelectedMove(squareMap, moves)
		} else {
			nextMove = Search(game, 2)
		}

		if nextMove != nil {
			game.Move(nextMove)
			squareMap = game.Position().Board().SquareMap()
			moves = game.ValidMoves()
			board.Update(squareMap)
		}

		board.Draw(surface)
	}
}
