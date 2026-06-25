package minesweeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type CellView struct {
	X        int  `json:"x"`
	Y        int  `json:"y"`
	Revealed bool `json:"revealed"`
	Flagged  bool `json:"flagged"`

	// -1 (Game Over)
	// 0-8 revealed adjacents value
	// 0 unrevealed cell
	Value int `json:"value"`
}

func (c CellView) MarshalJSON() ([]byte, error) {
	type Alias CellView
	return json.Marshal(Alias(c))
}

func (c *CellView) UnmarshalJSON(data []byte) error {
	type Alias CellView
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*c = CellView(aux)
	return nil
}

type GameView struct {
	Width     int          `json:"width"`
	Height    int          `json:"height"`
	MinesLeft int          `json:"mines_left"`
	Status    string       `json:"status"` // "playing", "won", "lost"
	Board     [][]CellView `json:"board"`
}

func (g GameView) MarshalJSON() ([]byte, error) {
	type Alias GameView

	if g.Board == nil {
		g.Board = [][]CellView{}
	}

	return json.Marshal(Alias(g))
}

func (g *GameView) UnmarshalJSON(data []byte) error {
	type Alias GameView
	var aux Alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	g.Width = aux.Width
	g.Height = aux.Height
	g.MinesLeft = aux.MinesLeft
	g.Status = aux.Status
	g.Board = aux.Board

	if g.Width < 0 || g.Height < 0 {
		return fmt.Errorf("invalid board width or height: %dx%d", g.Width, g.Height)
	}

	return nil
}

func (g *Game) View() GameView {
	viewBoard := make([][]CellView, g.width)
	for x := range viewBoard {
		viewBoard[x] = make([]CellView, g.height)
	}

	for x, y := range g.cells() {
		internalCell := g.board[x][y]

		viewCell := CellView{
			X:        x,
			Y:        y,
			Revealed: internalCell.revealed,
			Flagged:  internalCell.flagged,
		}

		if internalCell.revealed {
			if internalCell.kind == CellMine {
				viewCell.Value = -1
			} else {
				viewCell.Value = internalCell.value // Número 0-8
			}
		} else if g.status == StatusLost && internalCell.kind == CellMine {
			viewCell.Value = -1
		} else {
			viewCell.Value = 0
		}

		viewBoard[x][y] = viewCell
	}

	return GameView{
		Width:     g.width,
		Height:    g.height,
		MinesLeft: g.minesTotal - g.flagsPlaced,
		Status:    g.status.String(),
		Board:     viewBoard,
	}
}

func (g *Game) ToJSON() ([]byte, error) {
	currentView := g.View()
	return json.Marshal(currentView)
}

func (g *Game) String() string {
	return string(g.Bytes())
}

func (g *Game) Bytes() []byte {
	var out bytes.Buffer

	out.Grow(g.width*g.height + g.width)

	for x, y := range g.cells() {
		c := g.board[x][y]

		if !c.revealed {
			if c.flagged {
				out.WriteByte(SymbolFlag)
			} else {
				out.WriteByte(SymbolUnrevealed)
			}
		} else {
			switch c.kind {
			case CellEmpty:
				out.WriteByte(SymbolEmpty)
			case CellMine:
				out.WriteByte(SymbolMine)
			case CellCount:
				out.Write(strconv.AppendInt(out.AvailableBuffer(), int64(c.value), 10))
			}
		}

		if y == g.height-1 {
			out.WriteByte(SymbolBreakln)
		}
	}

	return out.Bytes()
}

func (g *Game) Render(w io.Writer) (err error) {
	_, err = w.Write(g.Bytes())
	return
}
