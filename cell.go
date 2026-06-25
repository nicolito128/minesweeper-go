package minesweeper

type CellKind uint8

const (
	CellEmpty CellKind = iota
	CellCount
	CellMine
)

type Cell struct {
	kind     CellKind
	value    int
	revealed bool
	flagged  bool
}

func NewCell(mine bool) Cell {
	c := Cell{}
	if mine {
		c.kind = CellMine
	}
	return c
}

type CellSymbol = byte

const (
	SymbolUnrevealed CellSymbol = 'o'
	SymbolEmpty      CellSymbol = '_'
	SymbolFlag       CellSymbol = '!'
	SymbolMine       CellSymbol = 'x'
	SymbolBreakln    CellSymbol = '\n'
)
