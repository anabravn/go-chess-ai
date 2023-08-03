package main

import (
	"testing"

	"github.com/notnil/chess"
)

func BenchmarkSearch(b *testing.B) {
	game := chess.NewGame()

	b.Run("Depth = 1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Search(game, 1)
		}
	})

	b.Run("Depth = 2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Search(game, 2)
		}
	})

	b.Run("Depth = 3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Search(game, 3)
		}
	})
}
