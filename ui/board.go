package ui

import (
	"github.com/notnil/chess"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const ScreenWidth, ScreenHeight = 400, 400
const SquareSize = 50

var (
	lightColor = sdl.Color{R: 200, G: 200, B: 200}
	darkColor  = sdl.Color{R: 50, G: 50, B: 50}
)

type BoardUI struct {
	BoardSurface  *sdl.Surface
	HintSurface   *sdl.Surface
	PiecesSurface *sdl.Surface
	PiecesSprites *sdl.Surface
	Size          int
	SquareSize    int
}

var pieceType = map[chess.PieceType]int{
	chess.King:   0,
	chess.Queen:  1,
	chess.Bishop: 2,
	chess.Knight: 3,
	chess.Rook:   4,
	chess.Pawn:   5,
}

var pieceColor = map[chess.Color]int{
	chess.Black: 1,
	chess.White: 0,
}

func (b *BoardUI) blitPiece(piece chess.Piece, file chess.File, rank chess.Rank) {
	var x, y int
	x = 7 - int(file)
	y = 7 - int(rank)

	color := piece.Color()
	t := piece.Type()

	squareRect := sdl.Rect{X: int32(x * SquareSize), Y: int32(y * SquareSize), W: SquareSize, H: SquareSize}
	pieceRect := sdl.Rect{
		X: int32(pieceType[t] * SquareSize),
		Y: int32(pieceColor[color] * SquareSize),
		W: int32(SquareSize),
		H: int32(SquareSize),
	}

	b.PiecesSprites.Blit(&pieceRect, b.PiecesSurface, &squareRect)
}

func (b *BoardUI) Update(squareMap map[chess.Square]chess.Piece) {
	pixel := sdl.MapRGBA(b.HintSurface.Format, 0, 0, 0, 0)
	b.PiecesSurface.FillRect(nil, pixel)

	for rank := chess.Rank8; rank >= chess.Rank1; rank-- {
		for file := chess.FileA; file <= chess.FileH; file++ {
			square := chess.NewSquare(file, rank)
			piece := squareMap[square]
			if piece.Type() != chess.NoPieceType {
				b.blitPiece(piece, file, rank)
			}
		}
	}
}

func fillRect(x, y, w, h int, color sdl.Color, surface *sdl.Surface) {
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
	pixel := sdl.MapRGBA(surface.Format, color.R, color.G, color.B, color.A)
	surface.FillRect(&rect, pixel)
}

func NewBoardUI() *BoardUI {
	size := 8 * SquareSize
	boardSurface, _ := sdl.CreateRGBSurface(0, int32(size), int32(size), 32, 0, 0, 0, 0)

	hintSurface, _ := sdl.CreateRGBSurface(0, int32(size), int32(size), 32, 0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	hintSurface.SetBlendMode(sdl.BLENDMODE_BLEND)

	piecesSurface, _ := sdl.CreateRGBSurface(0, int32(size), int32(size), 32, 0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	piecesSurface.SetBlendMode(sdl.BLENDMODE_BLEND)

	pieces, _ := img.Load("assets/pieces.png")

	var board = &BoardUI{
		Size:          size,
		SquareSize:    SquareSize,
		BoardSurface:  boardSurface,
		HintSurface:   hintSurface,
		PiecesSurface: piecesSurface,
		PiecesSprites: pieces,
	}

	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			var color sdl.Color

			if (file+rank)%2 == 0 {
				color = darkColor
			} else {
				color = lightColor
			}

			x := board.SquareSize * file
			y := board.SquareSize * rank

			fillRect(x, y, board.SquareSize, board.SquareSize, color, board.BoardSurface)
		}
	}

	return board
}

func (b *BoardUI) Draw(surface *sdl.Surface) {
	x := (ScreenWidth - b.Size) / 2
	y := (ScreenHeight - b.Size) / 2

	boardRect := sdl.Rect{X: 0, Y: 0, W: int32(b.Size), H: int32(b.Size)}
	screenRect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(b.Size), H: int32(b.Size)}

	b.BoardSurface.Blit(&boardRect, surface, &screenRect)
	b.HintSurface.Blit(&boardRect, surface, &screenRect)
	b.PiecesSurface.Blit(&boardRect, surface, &screenRect)
}

func (b *BoardUI) BlitHint(hintSquare chess.Square, color sdl.Color) {
	x := int(7-hintSquare.File()) * b.SquareSize
	y := int(7-hintSquare.Rank()) * b.SquareSize
	var rect = sdl.Rect{X: int32(x), Y: int32(y), W: int32(SquareSize), H: int32(SquareSize)}
	pixel := sdl.MapRGBA(b.HintSurface.Format, color.R, color.G, color.B, color.A)

	b.HintSurface.FillRect(&rect, pixel)
}

func (b *BoardUI) ClearHints() {
	pixel := sdl.MapRGBA(b.HintSurface.Format, 0, 0, 0, 0)
	b.HintSurface.FillRect(nil, pixel)
}

func (b *BoardUI) BlitHints(squareMap map[chess.Square]chess.Piece,
	moves []*chess.Move, startSquare chess.Square) {
	hintColor := sdl.Color{R: 0, G: 255, B: 100, A: 200}
	captureColor := sdl.Color{R: 255, G: 100, B: 0, A: 200}
	startColor := sdl.Color{R: 100, G: 100, B: 100, A: 200}

	b.ClearHints()
	b.BlitHint(startSquare, startColor)

	for _, move := range moves {
		if squareMap[move.S2()].Type() != chess.NoPieceType {
			b.BlitHint(move.S2(), captureColor)
		} else {
			b.BlitHint(move.S2(), hintColor)
		}
	}
}

func (b *BoardUI) GetSquare(x, y int32) chess.Square {
	offsetX := int32(ScreenWidth-b.Size) / 2
	offsetY := int32(ScreenHeight-b.Size) / 2
	file := (7 - ((x - offsetX) / SquareSize))
	rank := (7 - ((y - offsetY) / SquareSize))

	return chess.NewSquare(chess.File(file), chess.Rank(rank))
}
