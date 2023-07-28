package main

import (
	"github.com/notnil/chess"
)

func GetHints(moves []*chess.Move, startSquare chess.Square) []*chess.Move {
	var hints []*chess.Move

	for _, move := range moves {
		if startSquare == move.S1() {
			hints = append(hints, move)
		}
	}

	return hints
}

func GetPieceMove(moves []*chess.Move, startSquare chess.Square, endSquare chess.Square,
	squareMap map[chess.Square]chess.Piece) *chess.Move {
	for _, move := range moves {
		if endSquare == move.S2() && startSquare == move.S1() {
			return move
		}
	}

	return nil
}
