package main

import (
	"math"
	"strings"
	"testing"

	"github.com/notnil/chess"
)

func TestPieceValues(t *testing.T) {
	squareMap := chess.StartingPosition().Board().SquareMap()
	want := 390

	t.Run("White", func(t *testing.T) {
		got := PieceValues(squareMap, chess.White)
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("Black", func(t *testing.T) {
		got := PieceValues(squareMap, chess.Black)
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}

func TestSquareValues(t *testing.T) {
	squareMap := chess.StartingPosition().Board().SquareMap()
	want := 295

	t.Run("White", func(t *testing.T) {
		got := SquareValues(squareMap, chess.White)

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("Black", func(t *testing.T) {
		got := SquareValues(squareMap, chess.Black)

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}

func TestEval(t *testing.T) {
	t.Run("No Outcome", func(t *testing.T) {
		game := chess.NewGame()
		want := 295.0 + 390.0
		got := Eval(game, chess.White)

		if got != want {
			t.Errorf("got %f want %f", got, want)
		}
	})

	t.Run("Win", func(t *testing.T) {
		gameStr := "1.f3 e6 2.g4 Qh4#"
		pgnReader := strings.NewReader(gameStr)
		pgn, _ := chess.PGN(pgnReader)

		game := chess.NewGame(pgn)

		want := math.Inf(1)
		got := Eval(game, chess.Black)

		if got != want {
			t.Errorf("got %f want %f", got, want)
		}

	})
}
