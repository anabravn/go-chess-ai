package ai

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

			if color == chess.White {
				rank = 7 - rank
			}

			total += piecesSquares[piece.Type()][int(rank)*8+int(file)]
		}
	}

	return
}

func Heuristic(squareMap map[chess.Square]chess.Piece, color chess.Color) float64 {
	pieces := PieceValues(squareMap, color)
	position := SquareValues(squareMap, color)

	return float64(pieces) + float64(position)
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

	player := Heuristic(position.Board().SquareMap(), color)
	other := Heuristic(position.Board().SquareMap(), color.Other())

	return player - other
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
	_, move := NegaMax(game.Position(), math.Inf(-1), math.Inf(1), depth)
	return move
}

func NegaMax(position *chess.Position, alpha, beta float64, depth int) (float64, *chess.Move) {
	var bestMove *chess.Move = nil
	player := position.Turn()

	moves := position.ValidMoves()

	if len(moves) == 0 || depth == 0 {
		return Eval(position, player), nil
	}

	bestScore := math.Inf(-1)

	for _, move := range OrderMoves(moves) {
		score, _ := NegaMax(position.Update(move), -beta, -alpha, depth-1)
		score = -score

		if score >= beta {
			return score, move
		}

		if score > bestScore {
			bestScore, bestMove = score, move
			alpha = math.Max(alpha, score)
		}

	}

	return bestScore, bestMove
}
