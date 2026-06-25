package minesweeper

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/nicolito128/deterministic-minesweeper-go/prng"
)

type Coor struct {
	x, y int
}

type GameState uint8

const (
	StateUnknown GameState = iota
	StatePlaying
	StatePlayerLoses
	StatePlayerWins
)

func (gs GameState) String() string {
	switch gs {
	case StateUnknown:
		return "unknown"
	case StatePlaying:
		return "playing"
	case StatePlayerLoses:
		return "player loses"
	case StatePlayerWins:
		return "player wins"
	default:
		return ""
	}
}

type Game struct {
	// total amount of mines
	totalMines int
	// revealed mines
	countedMines int

	// board size
	size  int
	board [][]Cell

	// deterministic randomness by seed
	randSrc *prng.SeededSource64
	random  *rand.Rand

	// current state of the game (GameState)
	state GameState
}

func NewGame(size, totalMines int) (*Game, error) {
	return NewGameSeeded(rand.Uint64(), size, totalMines)
}

func NewGameSeeded(seed uint64, size, totalMines int) (*Game, error) {
	if size <= 0 {
		return nil, ErrInvalidSize
	}
	g := new(Game)
	g.size = size
	g.totalMines = totalMines
	g.randSrc = prng.NewSeededSource64(seed)
	g.random = rand.New(g.randSrc)

	g.board = make([][]Cell, size)
	for i := range size {
		row := make([]Cell, size)
		for j := range size {
			row[j] = NewCell(false)
		}
		g.board[i] = row
	}

	return g, nil
}

func (g *Game) State() GameState {
	return g.state
}

func (g *Game) Seed() uint64 {
	return g.randSrc.State()
}

func (g *Game) Rand() *rand.Rand {
	return g.random
}

func (g *Game) Mines() int {
	return g.totalMines
}

func (g *Game) CountedMines() int {
	return g.countedMines
}

func (g *Game) String() string {
	var out strings.Builder
	out.Grow(g.size*g.size + g.size)
	for i := range g.size {
		for j := range g.size {
			c := g.board[i][j]
			if !c.revealed {
				if !c.flagged {
					out.WriteByte('o')
					continue
				}
				if c.flagged {
					out.WriteByte('!')
					continue
				}
			}
			if c.revealed {
				switch c.kind {
				case CellEmpty:
					out.WriteByte(' ')
				case CellMine:
					out.WriteByte('x')
				case CellCount:
					fmt.Fprintf(&out, "%d", c.value)
				}
			}
		}
		out.WriteByte('\n')
	}
	return out.String()
}

func (g *Game) Handle(act *Action) error {
	if act == nil {
		return ErrInvalidAction
	}

	x, y := act.X, act.Y
	if g.state == StatePlayerLoses || g.state == StatePlayerWins {
		return nil
	}

	if !g.inBounds(x, y) {
		return ErrOutOfBounds
	}

	switch act.Kind {
	case ActionRevealCell:
		if g.state == StateUnknown {
			g.Start(x, y)
		}
		return g.RevealCell(x, y)
	case ActionToggleFlag:
		if g.state == StateUnknown {
			return nil
		}
		return g.ToggleFlag(x, y)
	}

	return nil
}

func (g *Game) Start(x, y int) error {
	if !g.inBounds(x, y) {
		return ErrOutOfBounds
	}
	if g.state != StateUnknown {
		return nil
	}

	g.generateInitialBoard(x, y, g.totalMines)
	g.state = StatePlaying
	return nil
}

func (g *Game) ToggleFlag(x, y int) error {
	if !g.inBounds(x, y) {
		return ErrOutOfBounds
	}
	if g.board[y][x].revealed {
		return ErrInvalidAction
	}

	g.board[x][y].flagged = !g.board[x][y].flagged
	g.countedMines += 1

	return nil
}

func (g *Game) RevealCell(startX, startY int) error {
	if !g.inBounds(startX, startY) {
		return ErrOutOfBounds
	}
	if g.board[startX][startY].revealed || g.board[startY][startX].flagged {
		return ErrInvalidAction
	}

	if g.board[startX][startY].kind == CellMine {
		g.board[startX][startY].revealed = true
		g.onRevealedBomb()
		return nil
	}

	if g.board[startX][startY].kind == CellCount {
		g.board[startX][startY].revealed = true
		return nil
	}

	queue := make([]Coor, 0)

	g.board[startX][startY].revealed = true
	queue = append(queue, Coor{x: startX, y: startY})

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}

				nx, ny := curr.x+dx, curr.y+dy

				if g.inBounds(nx, ny) && !g.board[nx][ny].revealed && !g.board[nx][ny].flagged {
					g.board[nx][ny].revealed = true

					if g.board[nx][ny].kind == CellEmpty {
						queue = append(queue, Coor{x: nx, y: ny})
					}
				}
			}
		}
	}

	return nil
}

func (g *Game) generateInitialBoard(initX, initY int, totaltotalMines int) {
	totalMinesPlaced := 0

	for totalMinesPlaced < totaltotalMines {
		x := g.random.Intn(g.size)
		y := g.random.Intn(g.size)

		if abs(x-initX) <= 1 && abs(y-initY) <= 1 {
			continue
		}

		if g.board[y][x].kind == CellMine {
			continue
		}

		g.board[y][x].kind = CellMine
		totalMinesPlaced++
	}

	g.calculateNeighborCounters()
}

func (g *Game) calculateNeighborCounters() {
	for i := range g.size {
		for j := range g.size {
			if g.board[i][j].kind == CellMine {
				continue
			}

			bombCount := 0
			for dj := -1; dj <= 1; dj++ {
				for di := -1; di <= 1; di++ {
					ni, nj := i+di, j+dj

					if g.inBounds(ni, nj) {
						if g.board[ni][nj].kind == CellMine {
							bombCount++
						}
					}
				}
			}

			if bombCount > 0 {
				g.board[i][j].kind = CellCount
				g.board[i][j].value = bombCount
			}
		}
	}
}

func (g *Game) onRevealedBomb() {
	g.state = StatePlayerLoses

	for i := range g.size {
		for j := range g.size {
			g.board[i][j].revealed = true
		}
	}
}

func (g *Game) inBounds(x, y int) bool {
	return x >= 0 && x < g.size && y >= 0 && y < g.size
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
