package minesweeper

type CellSymbol = byte

const (
	SymbolUnrevealed CellSymbol = 'o'
	SymbolEmpty      CellSymbol = '_'
	SymbolFlag       CellSymbol = '!'
	SymbolMine       CellSymbol = 'x'
	SymbolBreakln    CellSymbol = '\n'
)
