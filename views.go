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
	width := g.board.Width()
	height := g.board.Height()

	viewBoard := make([][]CellView, width)
	for x := range viewBoard {
		viewBoard[x] = make([]CellView, height)
	}

	for cell := range g.board.Cells() {
		viewCell := CellView{
			X:        cell.X,
			Y:        cell.Y,
			Revealed: cell.IsRevealed,
			Flagged:  cell.IsFlagged,
		}

		switch {
		case cell.IsRevealed && cell.IsMine:
			viewCell.Value = -1
		case cell.IsRevealed:
			viewCell.Value = int(cell.AdjacentMines)
		case g.status == StatusLost && cell.IsMine:
			viewCell.Value = -1
		default:
			viewCell.Value = 0
		}

		viewBoard[cell.X][cell.Y] = viewCell
	}

	return GameView{
		Width:     width,
		Height:    height,
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

	width := g.board.Width()
	height := g.board.Height()

	out.Grow(width*height + width)

	for cell := range g.board.Cells() {
		if !cell.IsRevealed {
			if cell.IsFlagged {
				out.WriteByte(SymbolFlag)
			} else {
				out.WriteByte(SymbolUnrevealed)
			}
		} else {
			switch {
			case cell.IsMine:
				out.WriteByte(SymbolMine)
			case cell.AdjacentMines > 0:
				out.Write(strconv.AppendInt(out.AvailableBuffer(), int64(cell.AdjacentMines), 10))
			default:
				out.WriteByte(SymbolEmpty)
			}
		}

		if cell.Y == height-1 {
			out.WriteByte(SymbolBreakln)
		}
	}

	return out.Bytes()
}

func (g *Game) Render(w io.Writer) (err error) {
	_, err = w.Write(g.Bytes())
	return
}
