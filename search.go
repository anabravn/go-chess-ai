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

func Eval(position *chess.Position, color chess.Color) float64 {
	switch position.Status() {
	case chess.Stalemate, chess.InsufficientMaterial,
		chess.ThreefoldRepetition, chess.FivefoldRepetition,
		chess.FiftyMoveRule, chess.SeventyFiveMoveRule:
		return 0.0
	case chess.Checkmate:
		if position.Turn() == color {
			return math.Inf(-1) // Loss
		} else {
			return math.Inf(1) // Win
		}
	}

	return float64(Heuristic(position.Board().SquareMap(), color))
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
	_, move := ValorMax(game.Position(), math.Inf(-1), math.Inf(1), depth, player)
	return move
}

func ValorMax(position *chess.Position, alpha, beta float64, depth int, player chess.Color) (float64, *chess.Move) {
	var bestMove *chess.Move = nil
	moves := OrderMoves(position.ValidMoves())

	if len(moves) == 0 || depth == 0 {
		return Eval(position, player), nil
	}

	bestScore := math.Inf(-1)

	for _, move := range moves {
		score, _ := ValorMin(position.Update(move), alpha, beta, depth, player)

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

func ValorMin(position *chess.Position, alpha, beta float64, depth int, player chess.Color) (float64, *chess.Move) {
	var bestMove *chess.Move = nil
	moves := OrderMoves(position.ValidMoves())

	if len(moves) == 0 || depth == 0 {
		return Eval(position, player), nil
	}

	bestScore := math.Inf(1)

	for _, move := range moves {
		score, _ := ValorMax(position.Update(move), alpha, beta, depth-1, player)

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
