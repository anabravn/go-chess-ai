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

func GetSelectedMove(moves []*chess.Move, piece chess.Piece, square chess.Square,
	squareMap map[chess.Square]chess.Piece) *chess.Move {
	for _, move := range moves {
		if square == move.S2() {
			s2 := squareMap[move.S1()]
			if s2 == piece {
				return move
			}
		}
	}

	return nil
}
