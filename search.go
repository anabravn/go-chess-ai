package main

import (
	"math"

	"github.com/notnil/chess"
)

func PieceValues(squareMap map[chess.Square]chess.Piece, color chess.Color) (total int) {
	for _, piece := range squareMap {
		if piece.Color() == color {
			total += pieceValues[piece.Type()]
		}
	}

	return
}

func SquareValues(squareMap map[chess.Square]chess.Piece, color chess.Color) (total int) {
	for square, piece := range squareMap {
		if piece.Color() == color {
			file, rank := square.File(), square.Rank()

			if color == chess.Black {
				rank = 7 - rank
			}

			total += piecesSquares[piece.Type()][int(rank)*8+int(file)]
		}
	}

	return
}

func Heuristic(squareMap map[chess.Square]chess.Piece, color chess.Color) int {
	pieces := PieceValues(squareMap, color)
	position := SquareValues(squareMap, color)

	return pieces + position
}

func Eval(game *chess.Game, color chess.Color) float64 {
	if game.Outcome() == chess.Draw {
		return 0
	} else if game.Method() == chess.Checkmate {
		if game.Position().Turn() == color {
			return math.Inf(-1) // Loss
		} else {
			return math.Inf(1) // Win
		}
	}

	return float64(Heuristic(game.Position().Board().SquareMap(), color))
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
		return Eval(game, player), nil
	}

	bestScore := math.Inf(-1)

	moves := OrderMoves(game.ValidMoves())
	for _, move := range moves {
		score, _ := ValorMin(Result(game, move), alpha, beta, depth, player)

		if score > bestScore {
			bestScore, bestMove = score, move
			alpha = math.Max(alpha, bestScore)
		}

		if score >= beta {
			return score, move
		}
	}

	return bestScore, bestMove
}

func ValorMin(game *chess.Game, alpha, beta float64, depth int, player chess.Color) (float64, *chess.Move) {
	var bestMove *chess.Move = nil

	if game.Outcome() != chess.NoOutcome {
		return Eval(game, player), nil
	}

	bestScore := math.Inf(1)

	moves := OrderMoves(game.ValidMoves())
	for _, move := range moves {
		score, _ := ValorMax(Result(game, move), alpha, beta, depth-1, player)

		if score < bestScore {
			bestScore, bestMove = score, move
			beta = math.Min(beta, bestScore)
		}

		if score <= alpha {
			return score, move
		}
	}

	return bestScore, bestMove
}
