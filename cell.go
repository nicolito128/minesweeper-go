package minesweeper

type CellKind uint8

const (
	CellEmpty CellKind = iota
	CellCount
	CellBomb
)

type Cell struct {
	kind     CellKind
	value    int
	revealed bool
	flagged  bool
}

func NewCell(bomb bool) Cell {
	c := Cell{}
	if bomb {
		c.kind = CellBomb
	}
	return c
}
