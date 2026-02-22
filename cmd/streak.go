package cmd

import (
	"context"
	"fmt"

	"github.com/GianSantos/githabit-cli/internal/api"
	"github.com/GianSantos/githabit-cli/internal/auth"
	"github.com/GianSantos/githabit-cli/internal/habit"
	"github.com/GianSantos/githabit-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(streakCmd)
}

var streakCmd = &cobra.Command{
	Use:   "streak",
	Short: "Show 30-day contribution grid",
	Long:  "Displays a 7x5 ANSI color grid of the last 30 days of contributions.",
	RunE:  runStreak,
}

func runStreak(cmd *cobra.Command, args []string) error {
	token, err := auth.GetToken()
	if err != nil {
		return fmt.Errorf("not initialized: run 'githabit init' first: %w", err)
	}

	ctx := context.Background()
	login, err := api.GetCurrentUser(ctx, token)
	if err != nil {
		return err
	}

	weeks, err := habit.FetchStreakData(ctx, token, login)
	if err != nil {
		return err
	}

	fmt.Println(ui.RenderStreakGrid(weeks))
	return nil
}
