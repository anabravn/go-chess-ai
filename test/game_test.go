package test

import (
	"ai"
	"testing"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

func BenchmarkSelf(b *testing.B) {
	// Jogar contra ela mesma
	b.Run("Depth = 2", func(b *testing.B) {
		game := chess.NewGame()

		for i := 0; i < b.N; i++ {
			for game.Outcome() == chess.NoOutcome {
				move := ai.Search(game, 2)
				game.Move(move)
			}
			if game.Outcome() != chess.Draw || len(game.Moves()) != 38 {
				b.Error()
			}
		}
	})

	b.Run("Depth = 4", func(b *testing.B) {
		game := chess.NewGame()

		for i := 0; i < b.N; i++ {
			for game.Outcome() == chess.NoOutcome {
				move := ai.Search(game, 4)
				game.Move(move)
			}
			if game.Outcome() != chess.Draw || len(game.Moves()) != 22 {
				b.Error()
			}
		}
	})
}

func BenchmarkStockfish(b *testing.B) {
	// set up engine to use stockfish exe
	eng, err := uci.New("./stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()

	// initialize uci with new game
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}

	b.Run("White", func(b *testing.B) {
		win, draw, loss := 0, 0, 0

		for i := 0; i < b.N; i++ {
			game := chess.NewGame()

			for game.Outcome() == chess.NoOutcome {
				move := ai.Search(game, 4)
				game.Move(move)

				if game.Outcome() != chess.NoOutcome {
					break
				}

				cmdPos := uci.CmdPosition{Position: game.Position()}
				cmdGo := uci.CmdGo{MoveTime: time.Second / 100}

				if err := eng.Run(cmdPos, cmdGo); err != nil {
					panic(err)
				}

				bestMove := eng.SearchResults().BestMove
				if err := game.Move(bestMove); err != nil {
					panic(err)
				}
			}

			switch game.Outcome() {
			case chess.WhiteWon:
				win++
			case chess.BlackWon:
				loss++
			case chess.Draw:
				draw++
			}
		}

		b.Logf("Wins: %d, Draws: %d, Losses: %d", win, draw, loss)
	})

	b.Run("Black", func(b *testing.B) {
		win, draw, loss := 0, 0, 0

		for i := 0; i < b.N; i++ {
			game := chess.NewGame()

			for game.Outcome() == chess.NoOutcome {
				cmdPos := uci.CmdPosition{Position: game.Position()}
				cmdGo := uci.CmdGo{MoveTime: time.Second / 100}

				if err := eng.Run(cmdPos, cmdGo); err != nil {
					panic(err)
				}

				move := eng.SearchResults().BestMove
				if err := game.Move(move); err != nil {
					panic(err)
				}

				if game.Outcome() != chess.NoOutcome {
					break
				}

				nextMove := ai.Search(game, 4)
				game.Move(nextMove)
			}

			switch game.Outcome() {
			case chess.WhiteWon:
				loss++
			case chess.BlackWon:
				win++
			case chess.Draw:
				draw++
			}
		}

		b.Logf("Wins: %d, Draws: %d, Losses: %d", win, draw, loss)
	})
}
