package main

import (
	"fmt"

	"github.com/notnil/chess"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	board      *BoardUI
	game       *chess.Game
	nextMove   *chess.Move
	squareMap  map[chess.Square]chess.Piece
	validMoves []*chess.Move
	searching  bool
}

func NewGame() *Game {
	g := &Game{}

	g.board = NewBoardUI()
	g.game = chess.NewGame()
	g.nextMove = nil
	g.squareMap = g.game.Position().Board().SquareMap()
	g.validMoves = g.game.ValidMoves()
	g.board.Update(g.squareMap)

	return g
}

func (g *Game) Update() {
	if g.game.Position().Turn() == chess.White {
		g.nextMove = g.GetSelectedMove()
	} else if !g.searching {

		fmt.Println("White: ", Eval(g.game.Position(), chess.White))
		fmt.Println("Black: ", Eval(g.game.Position(), chess.Black))

		g.nextMove = Search(g.game, 3)
	}

	if g.nextMove != nil {
		g.game.Move(g.nextMove)
		g.squareMap = g.game.Position().Board().SquareMap()
		g.validMoves = g.game.ValidMoves()
		g.board.Update(g.squareMap)
	}
}

func (g *Game) Draw(surface *sdl.Surface) {
	g.board.Draw(surface)
}

func GetHints(moves []*chess.Move, startSquare chess.Square) []*chess.Move {
	var hints []*chess.Move

	for _, move := range moves {
		if startSquare == move.S1() {
			hints = append(hints, move)
		}
	}

	return hints
}

func GetPieceMove(moves []*chess.Move, startSquare chess.Square, endSquare chess.Square) *chess.Move {
	for _, move := range moves {
		if endSquare == move.S2() && startSquare == move.S1() {
			return move
		}
	}

	return nil
}

var startSquare chess.Square

func (g *Game) GetSelectedMove() *chess.Move {
	mouseX, mouseY, mousePressed := sdl.GetMouseState()
	var nextMove *chess.Move = nil

	if mousePressed == 1 {
		endSquare := g.board.GetSquare(mouseX, mouseY)
		hints := GetHints(g.validMoves, endSquare)

		if len(hints) > 0 {
			g.board.BlitHints(g.squareMap, hints, endSquare)
			startSquare = endSquare
			nextMove = nil
		} else if startSquare != chess.NoSquare {
			nextMove = GetPieceMove(g.validMoves, startSquare, endSquare)

			g.board.ClearHints()
			startSquare = chess.NoSquare
		}
	}

	return nextMove
}
