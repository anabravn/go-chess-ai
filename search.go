package main

import (
	"math"

	"github.com/notnil/chess"
)

var values = map[chess.PieceType]float64{
	chess.Pawn:   10,
	chess.Knight: 30,
	chess.Bishop: 30,
	chess.Rook:   50,
	chess.Queen:  90,
}

func PieceValues(squareMap map[chess.Square]chess.Piece, color chess.Color) (total float64) {
	for _, piece := range squareMap {
		if piece.Color() == color {
			total += values[piece.Type()]
		}
	}

	return
}

func PieceSquareValues(squareMap map[chess.Square]chess.Piece, color chess.Color) (total float64) {
	for square, piece := range squareMap {
		if piece.Color() == color {
			file, rank := square.File(), square.Rank()

			if color == chess.Black {
				rank = 7 - rank
			}

			total += piecesSquares[piece.Type()][int(file)*8+int(rank)]
		}
	}

	return
}

func Heuristic(squareMap map[chess.Square]chess.Piece, color chess.Color) float64 {
	pieces := PieceValues(squareMap, color)
	position := PieceSquareValues(squareMap, color)
	other := PieceValues(squareMap, color.Other())

	return pieces + position - other
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
	player := game.Position().Turn()
	_, move := ValorMax(game, math.Inf(-1), math.Inf(1), depth, player)
	return move
}

func ValorMax(game *chess.Game, alpha, beta float64, depth int, player chess.Color) (float64, *chess.Move) {
	var bestMove *chess.Move = nil

	if game.Outcome() != chess.NoOutcome || depth == 0 {
		return Utility(game, player), nil
	}

	v := math.Inf(-1)

	moves := OrderMoves(game.ValidMoves())
	for _, move := range moves {
		v2, _ := ValorMin(Result(game, move), alpha, beta, depth, player)

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

func ValorMin(game *chess.Game, alpha, beta float64, depth int, player chess.Color) (float64, *chess.Move) {
	var bestMove *chess.Move = nil

	if game.Outcome() != chess.NoOutcome {
		return Utility(game, player), nil
	}

	v := math.Inf(1)

	moves := OrderMoves(game.ValidMoves())
	for _, move := range moves {
		v2, _ := ValorMax(Result(game, move), alpha, beta, depth-1, player)

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
