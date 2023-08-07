package ai

import (
	"github.com/notnil/chess"
)

type ChessAI struct {
	Game       *chess.Game
	validMoves []*chess.Move
	Searching  bool
	done       chan bool
	move       chan *chess.Move
}

func NewChessAi() *ChessAI {
	ai := &ChessAI{}

	ai.Game = chess.NewGame()
	ai.validMoves = ai.Game.ValidMoves()
	ai.Searching = false

	return ai
}

func (ai *ChessAI) Search() *chess.Move {
	ai.Searching = true
	var nextMove *chess.Move = nil

	nextMove = Search(ai.Game, 4)
	return nextMove
}

func (ai *ChessAI) Move(move *chess.Move) {
	ai.Game.Move(move)
}
