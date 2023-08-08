package ui

import (
	"github.com/notnil/chess"
	"github.com/veandco/go-sdl2/sdl"
)

// GameUI implementa a interface gráfica do jogo
type GameUI struct {
    // BoardUI é responsável pela 
    // interface do tabuleiro e coordenadas de casas
	board       *BoardUI

    // startSquare armazena a casa do tabuleiro 
    // selecionada, para determinar o próximo movimento
	startSquare chess.Square
}

// NewGameUI retorna uma nova estrutura
// do tipo GameUI
func NewGameUI() *GameUI {
	g := &GameUI{}

	g.board = NewBoardUI()
	g.startSquare = chess.NoSquare

	return g
}

// Update atualiza o tabuleiro com o estado 
// atual do jogo
func (g *GameUI) Update(game *chess.Game) {
    g.board.Update(game.Position().Board().SquareMap())	
}

// Desenha o tabuleiro na superfície surface
func (g *GameUI) Draw(surface *sdl.Surface) {
	g.board.Draw(surface)
}

// GetSelectedMove obtem o atual movimento selecionado.
// Se o usuário selecionou uma peça, as casas para as quais
// essa peça pode se mover são destacadas. Se o usuário 
// selecionou uma das casas, o movimento correspondente é 
// retornado.
//
// As variáveis startSquare e endSquare contém a casa 
// de origem do movimento e a casa de destino, respectivamente.
func (g *GameUI) GetSelectedMove(game *chess.Game) *chess.Move {
	mouseX, mouseY, mousePressed := sdl.GetMouseState()
	var nextMove *chess.Move = nil
	squareMap := game.Position().Board().SquareMap()
	validMoves := game.ValidMoves()

	if mousePressed == 1 {
		endSquare := g.board.GetSquare(mouseX, mouseY)
		hints := GetMovesToSquare(validMoves, endSquare)

		if len(hints) > 0 { 
			// Existem movimentos

			g.board.BlitHints(squareMap, hints, endSquare)
			g.startSquare = endSquare
			nextMove = nil
		} else if g.startSquare != chess.NoSquare {
			// Não existem movimentos, e uma casa estava selecionada

			nextMove = GetValidMove(validMoves, g.startSquare, endSquare)

			g.board.ClearHints()
			g.startSquare = chess.NoSquare
		}
	}

	return nextMove
}

// GetMovesToSquare recebe uma lista de movimentos e uma 
// casa do tabuleiro e retorna uma lista de movimentos que
// tem origem nessa casa
func GetMovesToSquare(moves []*chess.Move, startSquare chess.Square) []*chess.Move {
	var hints []*chess.Move

	for _, move := range moves {
		if startSquare == move.S1() {
			hints = append(hints, move)
		}
	}

	return hints
}

// GetValidMove recebe uma lista de movientos, uma casa de origem
// e uma casa de destino e retorna o movimento correspondente da
// lista. Retorna nil caso não exista um movimento válido.
func GetValidMove(moves []*chess.Move, startSquare chess.Square, endSquare chess.Square) *chess.Move {
	for _, move := range moves {
		if endSquare == move.S2() && startSquare == move.S1() {
			return move
		}
	}

	return nil
}


