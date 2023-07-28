package main

import (
	"math"

	"github.com/notnil/chess"
)

var values = map[chess.PieceType]int{
	chess.NoPieceType: 0,
	chess.King:        0,
	chess.Pawn:        1,
	chess.Knight:      3,
	chess.Bishop:      3,
	chess.Rook:        5,
	chess.Queen:       9,
}

func PieceValues(squareMap map[chess.Square]chess.Piece, color chess.Color) (total int) {
	for _, piece := range squareMap {
		if piece.Color() == color {
			total += values[piece.Type()]
		}
	}

	return
}

func Heuristic(squareMap map[chess.Square]chess.Piece, color chess.Color) float64 {
	pieces := float64(PieceValues(squareMap, color))
	other := float64(PieceValues(squareMap, color.Other()))
	return pieces - other
}

func Utility(game *chess.Game, color chess.Color) float64 {
	if game.Outcome() == chess.Draw {
		return 0
	} else if game.Method() == chess.Checkmate {
		if game.Position().Turn() == color {
			return math.Inf(-1) // Loss
		} else {
			return math.Inf(1) // Win
		}
	}

	return Heuristic(game.Position().Board().SquareMap(), color)
}

func Result(game *chess.Game, move *chess.Move) *chess.Game {
	newGame := game.Clone()
	newGame.Move(move)

	return newGame
}

func OrderMoves(moves []*chess.Move) []*chess.Move {
	var checkMoves, captureMoves,
		orderedMoves, rest []*chess.Move

	for _, move := range moves {
		capture := move.HasTag(chess.Capture)
		check := move.HasTag(chess.Check)
		move.S1()

		if check {
			checkMoves = append(checkMoves, move)
		} else if capture {
			captureMoves = append(captureMoves, move)
		} else {
			rest = append(rest, move)
		}
	}

	orderedMoves = append(checkMoves, captureMoves...)
	orderedMoves = append(orderedMoves, rest...)

	return orderedMoves
}

func Search(game *chess.Game, depth int) *chess.Move {
	_, move := ValorMax(game, math.Inf(-1), math.Inf(1), depth)
	return move
}

func ValorMax(game *chess.Game, alpha, beta float64, depth int) (float64, *chess.Move) {
	turn := game.Position().Turn()
	var bestMove *chess.Move = nil

	if game.Outcome() != chess.NoOutcome || depth == 0 {
		return Utility(game, turn), nil
	}

	v := math.Inf(-1)

	moves := OrderMoves(game.ValidMoves())
	for _, move := range moves {
		v2, _ := ValorMin(Result(game, move), alpha, beta, depth)
		// cada nível é uma jogada, não um turno

		if v2 > v {
			v, bestMove = v2, move
			alpha = math.Max(alpha, v)
		}

		if v >= beta {
			return v, bestMove
		}
	}

	return v, bestMove
}

func ValorMin(game *chess.Game, alpha, beta float64, depth int) (float64, *chess.Move) {
	turn := game.Position().Turn()
	var bestMove *chess.Move = nil

	if game.Outcome() != chess.NoOutcome {
		return Utility(game, turn), nil
	}

	v := math.Inf(1)

	moves := OrderMoves(game.ValidMoves())
	for _, move := range moves {
		v2, _ := ValorMax(Result(game, move), alpha, beta, depth-1)

		if v2 < v {
			v, bestMove = v2, move
			beta = math.Min(beta, v)
		}
		if v <= alpha {
			return v, bestMove
		}
	}

	return v, bestMove
}
