package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/GianSantos/githabit-cli/internal/api"
	"github.com/GianSantos/githabit-cli/internal/auth"
	"github.com/GianSantos/githabit-cli/internal/habit"
	"github.com/GianSantos/githabit-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show today's habit score and current streak",
	Long:  "Displays a visual breakdown of today's points and current streak.",
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	token, err := auth.GetToken()
	if err != nil {
		return fmt.Errorf("not initialized: run 'githabit init' first: %w", err)
	}

	state, _ := habit.LoadState()

	// Cache: skip API if last_checked < 1 hour ago
	if habit.IsCacheValid(state) {
		fmt.Println(ui.FormatScore(state.TodayScore,
			state.TodayBreakdown.Commits,
			state.TodayBreakdown.PRs,
			state.TodayBreakdown.Reviews,
			state.TodayBreakdown.Issues))
		fmt.Println(ui.FormatStreak(state.CurrentStreak))
		return nil
	}

	ctx := context.Background()
	login, err := api.GetCurrentUser(ctx, token)
	if err != nil {
		return err
	}

	dc, score, err := habit.FetchTodayScore(ctx, token, login)
	if err != nil {
		return err
	}

	weeks, err := habit.FetchStreakData(ctx, token, login)
	if err != nil {
		return err
	}
	streak := habit.ComputeStreakFromCalendar(weeks)

	// Save to state for cache/nudge
	state.LastChecked = time.Now()
	state.TodayScore = score
	state.CurrentStreak = streak
	state.TodayBreakdown.Commits = dc.Commits
	state.TodayBreakdown.PRs = dc.PRs
	state.TodayBreakdown.Reviews = dc.Reviews
	state.TodayBreakdown.Issues = dc.Issues
	_ = habit.SaveState(state)

	fmt.Println(ui.FormatScore(score, dc.Commits, dc.PRs, dc.Reviews, dc.Issues))
	fmt.Println(ui.FormatStreak(streak))
	return nil
}
