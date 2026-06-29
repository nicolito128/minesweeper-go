package minesweeper

import (
	"math/rand"

	"github.com/nicolito128/minesweeper-go/prng"
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

	board *Board

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
	if width < 3 || height < 3 {
		return nil, ErrInvalidBoard
	}
	if minesTotal <= 0 {
		return nil, ErrInvalidTotalMines
	}
	if minesTotal >= (width * height / 2) {
		return nil, ErrInvalidBoard
	}
	g := new(Game)
	g.minesTotal = minesTotal
	g.randSrc = prng.NewSeededSource64(seed)
	g.random = rand.New(g.randSrc)
	g.board = NewBoard(width, height)

	return g, nil
}

func (g *Game) Board() *Board {
	return g.board
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

	if !g.board.InBounds(x, y) {
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
	if !g.board.InBounds(x, y) {
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
	if !g.board.InBounds(x, y) {
		return ErrOutOfBounds
	}
	if g.board.Cell(x, y).IsRevealed {
		return ErrInvalidAction
	}

	g.board.Cell(x, y).IsFlagged = !g.board.Cell(x, y).IsFlagged
	if g.board.Cell(x, y).IsFlagged {
		g.flagsPlaced++
	} else {
		g.flagsPlaced--
	}

	return nil
}

func (g *Game) RevealCell(startX, startY int) error {
	if !g.board.InBounds(startX, startY) {
		return ErrOutOfBounds
	}
	if g.board.Cell(startX, startY).IsRevealed || g.board.Cell(startX, startY).IsFlagged {
		return ErrInvalidAction
	}

	if g.board.Cell(startX, startY).IsMine {
		g.board.Cell(startX, startY).Reveal()
		g.onRevealedMine()
		return nil
	}

	if g.board.Cell(startX, startY).AdjacentMines > 0 {
		g.board.Cell(startX, startY).Reveal()
		g.revealedCells++
		return nil
	}

	queue := make([][2]int, 0)

	g.board.Cell(startX, startY).IsRevealed = true
	g.revealedCells++
	queue = append(queue, [2]int{startX, startY})

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for ac := range g.board.AdjacentCells(curr[0], curr[1]) {
			if !ac.IsRevealed && !ac.IsFlagged {
				ac.IsRevealed = true
				g.revealedCells++

				if ac.IsEmpty() {
					queue = append(queue, [2]int{ac.X, ac.Y})
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
		randx := g.random.Intn(g.board.width)
		randy := g.random.Intn(g.board.height)

		if abs(randx-initX) <= 1 && abs(randy-initY) <= 1 {
			continue
		}

		if g.board.Cell(randx, randy).IsMine {
			continue
		}

		g.board.Cell(randx, randy).IsMine = true
		minesTotalPlaced++
	}

	g.calculateNeighborCounters()
}

func (g *Game) calculateNeighborCounters() {
	for c := range g.board.Cells() {
		if c.IsMine {
			continue
		}

		minesCount := uint8(0)
		for ac := range g.board.AdjacentCells(c.X, c.Y) {
			if ac.IsMine {
				minesCount++
			}
		}

		if minesCount > 0 {
			c.AdjacentMines = minesCount
		}
	}
}

func (g *Game) onRevealedMine() {
	g.status = StatusLost

	for c := range g.board.Cells() {
		if !c.IsRevealed {
			c.IsRevealed = true
			g.revealedCells++
		}
	}
}

func (g *Game) onWinCase() {
	if g.revealedCells == (g.board.width*g.board.height - g.minesTotal) {
		g.status = StatusWon
		for c := range g.board.Cells() {
			if c.IsMine {
				c.IsFlagged = true
				c.IsRevealed = false
			} else {
				c.IsRevealed = true
				c.IsFlagged = false
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
