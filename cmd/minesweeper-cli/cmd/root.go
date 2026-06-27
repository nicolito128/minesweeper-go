package cmd

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/nicolito128/minesweeper-go"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "minesweeper-cli",
	Short: "A deterministic minesweeper CLI utility.",
	Long: `You can simulate a game of minesweeper passing a seed or taking a generated seed by the program.
	Example usage:
		minesweeper-cli --seed 123456789 --width 12 --height 24 -a "<actions>"
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		seed, _ := cmd.Flags().GetUint64("seed")
		if seed == 0 {
			seed = rand.Uint64()
		}

		width, err := cmd.Flags().GetInt("width")
		if err != nil {
			return err
		}

		height, err := cmd.Flags().GetInt("height")
		if err != nil {
			return err
		}

		mines, err := cmd.Flags().GetInt("mines")
		if err != nil {
			return err
		}

		actionsStr, err := cmd.Flags().GetString("actions")
		if err != nil {
			return err
		}

		actions, err := minesweeper.ParseActions(actionsStr)
		if err != nil {
			return err
		}

		game, err := minesweeper.NewGameSeeded(seed, width, height, mines)
		if err != nil {
			return err
		}

		for _, action := range actions {
			game.Handle(action)
		}

		b, err := game.View().MarshalJSON()
		if err != nil {
			return err
		}

		fmt.Println(string(b))
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().Uint64P("seed", "s", 0, "Seed to fill the randomness")
	rootCmd.Flags().StringP("actions", "a", "", "List of interactions with the current game")
	rootCmd.Flags().Int("width", 9, "Width of the board")
	rootCmd.Flags().Int("height", 9, "Height of the board")
	rootCmd.Flags().IntP("mines", "m", 10, "Total amount of mines for the game")
}
