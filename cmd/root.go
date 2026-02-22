package cmd

import (
	"fmt"
	"os"

	"github.com/GianSantos/githabit-cli/internal/habit"
	"github.com/GianSantos/githabit-cli/internal/ui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "githabit",
	Short: "A habit-building CLI that tracks GitHub contributions",
	Long:  "GitHabit gamifies your coding practice with a custom Habit Score and helps you maintain streaks.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = "0.1.0"

	// PersistentPreRun: nudge when today's score is 0 (except for init, help, version, notify, check-reminder)
	rootCmd.PersistentPreRun = func(c *cobra.Command, args []string) {
		switch c.Name() {
		case "init", "help", "check-reminder":
			return
		}
		if c.Parent() != nil && c.Parent().Name() == "notify" && (c.Name() == "start" || c.Name() == "stop") {
			return
		}

		state, err := habit.LoadState()
		if err != nil {
			return
		}
		if state.TodayScore == 0 {
			fmt.Println(ui.NudgeStyle.Render("⚠ Today's Habit Score is 0. Don't break your streak!"))
		}
	}
}
