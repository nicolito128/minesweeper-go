package main

import (
	"flag"
	"fmt"
	"math/rand"

	minesweeper "github.com/nicolito128/deterministic-minesweeper-go"
)

var (
	widthFlag   = flag.Int("w", 9, "Width of the board")
	heightFlag  = flag.Int("h", 9, "Height of the board")
	minesFlag   = flag.Int("mines", 10, "Total amount of mines")
	seedFlag    = flag.Uint64("seed", 0, "Seed for randomness")
	actionsFlag = flag.String("act", "", "Chain of actions")
)

func main() {
	flag.Parse()

	width := *widthFlag
	height := *heightFlag
	mines := *minesFlag
	seed := *seedFlag
	if seed == 0 {
		seed = rand.Uint64()
	}

	g, err := minesweeper.NewGameSeeded(seed, width, height, mines)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[ Seed: %v ]\n", g.Seed())

	actions, err := minesweeper.ParseActions(*actionsFlag)
	if err != nil {
		panic(err)
	}

	printGame(g)
	for _, action := range actions {
		if err := g.Handle(action); err != nil {
			fmt.Println("Error:", err)
		}
		printGame(g)
	}
}

func printGame(g *minesweeper.Game) {
	fmt.Printf("[ Mines: %v ]\n", g.Mines()-g.FlagsPlaced())
	fmt.Printf("[ Status: %s ]\n", g.Status())
	fmt.Println(g)
}
