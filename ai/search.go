// Package ai implementa o o algoritmo 
// de busca para xadrez
package ai

import (
	"math"

	"github.com/notnil/chess"
)

// PieceValues recebe um mapa de casas para peças e uma cor,
// e retorna a soma dos valores materiais das peças existentes dessa cor.
func PieceValues(squareMap map[chess.Square]chess.Piece, color chess.Color) (total int) {
	for _, piece := range squareMap {
		if piece.Color() == color {
			total += pieceValues[piece.Type()]
		}
	}

	return
}

// SquareValues recebe um mapa de casas para peças e uma cor,
// e retorna a soma dos valores bônus de posição de cada peça
// existente dessa cor
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

// Heuristic recebe um mapa de casas para peças e uma cor, e 
// retorna o valor em ponto flutuante do material e bonus de posições
// para peças dessa cor.
func Heuristic(squareMap map[chess.Square]chess.Piece, color chess.Color) float64 {
	pieces := PieceValues(squareMap, color)
	position := SquareValues(squareMap, color)

	return float64(pieces) + float64(position)
}

// Eval recebe uma posição e uma cor, e retorna a avaliação dessa 
// posição para o jogador dessa cor. Para estados não terminais, a
// pontuação é calculada com base na função heurística.
//
// Casos especiais:
//   Vitória: +inf
//   Derrota: -inf
//   Empate: 0
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

// OrderMoves implementa a ordenação de movimentos. Recebe
// uma lista de movimentos e retorna uma lista ordenada,
// priorizando cheques e capturas.
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

// Search implementa o algoritmo de busca. Recebe uma instancia
// de jogo e uma profundidade, e retorna o próximo movimento.
// Chama recursivamente a função NegaMax.
//
// A profundidade é entendida em termo de movimentos, portanto
// deve ser um número múltiplo de dois para avaliar cada turno
// por completo.
func Search(game *chess.Game, depth int) *chess.Move {
	_, move := NegaMax(game.Position(), math.Inf(-1), math.Inf(1), depth)
	return move
}

// NegaMax implementa o algoritmo de busca minimax com poda alfa-beta
// heurística. Cada jogador age como o jogador MAX, utilizando os 
// parametros inversos da execução anterior.
func NegaMax(position *chess.Position, alpha, beta float64, depth int) (float64, *chess.Move) {
	var bestMove *chess.Move = nil
	player := position.Turn()

	moves := position.ValidMoves()

	if len(moves) == 0 || depth == 0 { // Teste de corte
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
