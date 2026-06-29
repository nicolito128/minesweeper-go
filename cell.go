package minesweeper

type Cell struct {
	X, Y          int
	AdjacentMines uint8
	IsMine        bool
	IsRevealed    bool
	IsFlagged     bool
}

func NewCell(x, y int) *Cell {
	c := new(Cell)
	c.X = x
	c.Y = y
	return c
}

func (c *Cell) Reveal() {
	c.IsRevealed = true
}

func (c *Cell) ToggleFlag() {
	c.IsFlagged = !c.IsFlagged
}

func (c *Cell) IsEmpty() bool {
	return !c.IsMine && c.AdjacentMines == 0
}
