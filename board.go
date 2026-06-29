package minesweeper

import (
	"iter"
)

type Board struct {
	width, height int
	cells         []Cell
}

func NewBoard(width, height int) *Board {
	if width < 0 {
		width = -width
	}
	if height < 0 {
		height = -height
	}
	cells := make([]Cell, width*height)
	board := &Board{
		width:  width,
		height: height,
		cells:  cells,
	}
	board.initCellCoordinates()
	return board
}

func (b *Board) initCellCoordinates() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			cell := &b.cells[(y*b.width)+x]
			cell.X = x
			cell.Y = y
		}
	}
}

func (b *Board) Width() int {
	return b.width
}

func (b *Board) Height() int {
	return b.height
}

func (b *Board) Cell(x, y int) *Cell {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return nil
	}

	index := (y * b.width) + x
	return &b.cells[index]
}

func (b *Board) InBounds(x, y int) bool {
	return x >= 0 && x < b.width && y >= 0 && y < b.height
}

func (b *Board) Cells() iter.Seq[*Cell] {
	return func(yield func(*Cell) bool) {
		for x := range b.width {
			for y := range b.height {
				if !yield(b.Cell(x, y)) {
					return
				}
			}
		}
	}
}

func (b *Board) AdjacentCells(x, y int) iter.Seq[*Cell] {
	return func(yield func(*Cell) bool) {
		if b.InBounds(x, y) {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					nx, ny := x+dx, y+dy
					if nx == 0 && ny == 0 {
						continue
					}

					if b.InBounds(nx, ny) {
						if !yield(b.Cell(nx, ny)) {
							return
						}
					}
				}
			}
		}
	}
}
