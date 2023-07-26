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

	window, _ := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		ScreenWidth, ScreenHeight, sdl.WINDOW_SHOWN)
	defer window.Destroy()

	surface, _ := window.GetSurface()
	surface.FillRect(nil, 0)

	pieces, _ := img.Load("pieces.png")
	board := NewBoardUI(pieces)

	running := true

	game := chess.NewGame()

	squareMap := game.Position().Board().SquareMap()
	moves := game.ValidMoves()
	board.Update(squareMap)

	selected := chess.NoPiece
	var nextMove *chess.Move = nil

	for running {
		window.UpdateSurface()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			}
		}

		mouseX, mouseY, mousePressed := sdl.GetMouseState()

		if mousePressed == 1 {
			square := board.GetSquare(mouseX, mouseY)
			hints := GetHints(moves, square)

			if len(hints) > 0 {
				board.BlitHints(squareMap, hints, square)
				selected = squareMap[square]
				nextMove = nil
			} else if selected != chess.NoPiece {
				nextMove = GetSelectedMove(moves, selected, square, squareMap)

				game.Move(nextMove)

				squareMap = game.Position().Board().SquareMap()
				moves = game.ValidMoves()

				board.Update(squareMap)
				board.ClearHints()

				selected = chess.NoPiece
			}
		}

		board.Draw(surface)
	}
}
