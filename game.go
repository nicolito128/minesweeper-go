package minesweeper

import (
	"iter"
	"math/rand"

	"github.com/nicolito128/deterministic-minesweeper-go/prng"
)

type GameStatus uint8

const (
	StatusUnknown GameStatus = iota
	StatusPlaying
	StatusLost
	StatusWon
)

func (gs GameStatus) String() string {
	switch gs {
	case StatusUnknown:
		return "unknown"
	case StatusPlaying:
		return "playing"
	case StatusLost:
		return "lost"
	case StatusWon:
		return "won"
	default:
		return ""
	}
}

type Game struct {
	// total amount of mines
	minesTotal  int
	flagsPlaced int

	revealedCells int

	// board size
	width, height int // TODO
	board         [][]Cell

	// deterministic randomness by seed
	randSrc *prng.SeededSource64
	random  *rand.Rand

	// current status of the game (GameState)
	status GameStatus
}

func NewGame(width, height, minesTotal int) (*Game, error) {
	return NewGameSeeded(rand.Uint64(), width, height, minesTotal)
}

func NewGameSeeded(seed uint64, width, height, minesTotal int) (*Game, error) {
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidSize
	}
	if minesTotal <= 0 {
		return nil, ErrInvalidTotalMines
	}
	g := new(Game)
	g.width = width
	g.height = height
	g.minesTotal = minesTotal
	g.randSrc = prng.NewSeededSource64(seed)
	g.random = rand.New(g.randSrc)

	g.board = make([][]Cell, width)
	for i := range width {
		row := make([]Cell, height)
		for j := range height {
			row[j] = NewCell(false)
		}
		g.board[i] = row
	}

	return g, nil
}

func (g *Game) Status() GameStatus {
	return g.status
}

func (g *Game) Seed() uint64 {
	return g.randSrc.State()
}

func (g *Game) Rand() *rand.Rand {
	return g.random
}

func (g *Game) Mines() int {
	return g.minesTotal
}

func (g *Game) FlagsPlaced() int {
	return g.flagsPlaced
}

func (g *Game) Handle(act *Action) error {
	if act == nil {
		return ErrInvalidAction
	}

	x, y := act.X, act.Y
	if g.status == StatusLost || g.status == StatusWon {
		return nil
	}

	if !g.inBounds(x, y) {
		return ErrOutOfBounds
	}

	switch act.Kind {
	case ActionRevealCell:
		if g.status == StatusUnknown {
			return g.Start(x, y)
		}
		return g.RevealCell(x, y)
	case ActionToggleFlag:
		if g.status == StatusUnknown {
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
	if g.status != StatusUnknown {
		return nil
	}

	g.generateInitialBoard(x, y, g.minesTotal)
	g.status = StatusPlaying
	g.RevealCell(x, y)
	return nil
}

func (g *Game) ToggleFlag(x, y int) error {
	if !g.inBounds(x, y) {
		return ErrOutOfBounds
	}
	if g.board[x][y].revealed {
		return ErrInvalidAction
	}

	g.board[x][y].flagged = !g.board[x][y].flagged
	if g.board[x][y].flagged {
		g.flagsPlaced++
	} else {
		g.flagsPlaced--
	}

	return nil
}

func (g *Game) RevealCell(startX, startY int) error {
	if !g.inBounds(startX, startY) {
		return ErrOutOfBounds
	}
	if g.board[startX][startY].revealed || g.board[startX][startX].flagged {
		return ErrInvalidAction
	}

	if g.board[startX][startY].kind == CellMine {
		g.board[startX][startY].revealed = true
		g.onRevealedMine()
		return nil
	}

	if g.board[startX][startY].kind == CellCount {
		g.board[startX][startY].revealed = true
		g.revealedCells++
		return nil
	}

	queue := make([][2]int, 0)

	g.board[startX][startY].revealed = true
	g.revealedCells++
	queue = append(queue, [2]int{startX, startY})

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for nx, ny := range g.adjacentCells(curr[0], curr[1]) {
			if nx == 0 && ny == 0 {
				continue
			}

			if !g.board[nx][ny].revealed && !g.board[nx][ny].flagged {
				g.board[nx][ny].revealed = true
				g.revealedCells++

				if g.board[nx][ny].kind == CellEmpty {
					queue = append(queue, [2]int{nx, ny})
				}
			}
		}
	}

	g.onWinCase()

	return nil
}

func (g *Game) generateInitialBoard(initX, initY int, totalminesTotal int) {
	minesTotalPlaced := 0

	for minesTotalPlaced < totalminesTotal {
		randx := g.random.Intn(g.width)
		randy := g.random.Intn(g.height)

		if abs(randx-initX) <= 1 && abs(randy-initY) <= 1 {
			continue
		}

		if g.board[randx][randy].kind == CellMine {
			continue
		}

		g.board[randx][randy].kind = CellMine
		minesTotalPlaced++
	}

	g.calculateNeighborCounters()
}

func (g *Game) calculateNeighborCounters() {
	for x, y := range g.cells() {
		if g.board[x][y].kind == CellMine {
			continue
		}

		minesCount := 0
		for nx, ny := range g.adjacentCells(x, y) {
			if g.board[nx][ny].kind == CellMine {
				minesCount++
			}
		}

		if minesCount > 0 {
			g.board[x][y].kind = CellCount
			g.board[x][y].value = minesCount
		}
	}
}

func (g *Game) onRevealedMine() {
	g.status = StatusLost

	for x, y := range g.cells() {
		g.board[x][y].revealed = true
		g.revealedCells++
	}
}

func (g *Game) onWinCase() {
	if g.revealedCells == (g.width*g.height - g.minesTotal) {
		g.status = StatusWon
		for x, y := range g.cells() {
			g.board[x][y].flagged = true
		}
	}
}

func (g *Game) inBounds(x, y int) bool {
	return x >= 0 && x < g.width && y >= 0 && y < g.height
}

func (g *Game) cells() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for x := range g.width {
			for y := range g.height {
				if !yield(x, y) {
					return
				}
			}
		}
	}
}

func (g *Game) adjacentCells(x, y int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		if g.inBounds(x, y) {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					nx, ny := x+dx, y+dy

					if g.inBounds(nx, ny) {
						if !yield(nx, ny) {
							return
						}
					}
				}
			}
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
