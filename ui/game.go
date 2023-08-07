package ui

import (
	"ai"
	"fmt"

	"github.com/notnil/chess"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	board       *BoardUI
	ai          *ai.ChessAI
	nextMove    *chess.Move
	startSquare chess.Square
}

func NewGame() *Game {
	g := &Game{}
	g.ai = ai.NewChessAi()

	g.board = NewBoardUI()
	g.board.Update(g.ai.Game.Position().Board().SquareMap())
	g.nextMove = nil
	g.startSquare = chess.NoSquare

	return g
}

func (g *Game) Update() bool {
	if g.ai.Game.Position().Turn() == chess.White {
		g.nextMove = g.GetSelectedMove()
	} else {
		g.nextMove = g.ai.Search()
	}

	if g.nextMove != nil {
		g.ai.Move(g.nextMove)
		g.board.Update(g.ai.Game.Position().Board().SquareMap())
		g.nextMove = nil
	}

	if g.ai.Game.Outcome() != chess.NoOutcome {
		fmt.Println(g.ai.Game.Outcome(), g.ai.Game.Method())
		return false
	}

	return true
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

func (g *Game) GetSelectedMove() *chess.Move {
	mouseX, mouseY, mousePressed := sdl.GetMouseState()
	var nextMove *chess.Move = nil
	squareMap := g.ai.Game.Position().Board().SquareMap()
	validMoves := g.ai.Game.ValidMoves()

	if mousePressed == 1 {
		endSquare := g.board.GetSquare(mouseX, mouseY)
		hints := GetHints(validMoves, endSquare)

		if len(hints) > 0 {
			g.board.BlitHints(squareMap, hints, endSquare)
			g.startSquare = endSquare
			nextMove = nil
		} else if g.startSquare != chess.NoSquare {
			nextMove = GetPieceMove(validMoves, g.startSquare, endSquare)

			g.board.ClearHints()
			g.startSquare = chess.NoSquare
		}
	}

	return nextMove
}
