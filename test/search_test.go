package test

import (
	"ai"
	"testing"

	"github.com/notnil/chess"
)

func BenchmarkSearch(b *testing.B) {
	game := chess.NewGame()

	b.Run("Depth = 2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ai.Search(game, 2)
		}
	})

	b.Run("Depth = 4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ai.Search(game, 4)
		}
	})

	b.Run("Depth = 6", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ai.Search(game, 6)
		}
	})

	b.Run("Depth = 8", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ai.Search(game, 8)
		}
	})
}
