package test

import (
	"ai"
	"testing"

	"github.com/notnil/chess"
)

func BenchmarkSearchSelf(b *testing.B) {
	// Jogar contra ela mesma
	b.Run("Depth = 2", func(b *testing.B) {
		game := chess.NewGame()
		turns := 0

		for i := 0; i < 10; i++ {
			for game.Outcome() == chess.NoOutcome {
				move := ai.Search(game, 2)
				game.Move(move)

				if game.Position().Turn() == chess.White {
					turns++
				}
			}
			b.Logf("Result: %s, Turns: %d\n", game.Outcome(), turns)
	}})

	b.Run("Depth = 4", func(b *testing.B) {
		game := chess.NewGame()
		turns := 0

		for i := 0; i < 10; i++ {
			for game.Outcome() == chess.NoOutcome {
				move := ai.Search(game, 4)
				game.Move(move)

				if game.Position().Turn() == chess.White {
					turns++
				}
			}
			b.Logf("Result: %s, Turns: %d\n", game.Outcome(), turns)
	}})
}

