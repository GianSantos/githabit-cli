package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/GianSantos/githabit-cli/internal/api"
	"github.com/GianSantos/githabit-cli/internal/auth"

	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(feedCmd)
}

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "List followed users' activity and estimated daily habit score",
	Long:  "Shows recent activity from users you follow, with estimated habit scores for competition.",
	RunE:  runFeed,
}

func runFeed(cmd *cobra.Command, args []string) error {
	token, err := auth.GetToken()
	if err != nil {
		return fmt.Errorf("not initialized: run 'githabit init' first: %w", err)
	}

	ctx := context.Background()
	login, err := api.GetCurrentUser(ctx, token)
	if err != nil {
		return err
	}

	following, err := api.GetFollowing(ctx, token, login)
	if err != nil {
		return err
	}

	// Limit to first 10 for display
	maxUsers := 10
	if len(following) < maxUsers {
		maxUsers = len(following)
	}

	fmt.Println("Followed users' recent activity:")
	fmt.Println()

	for i := 0; i < maxUsers; i++ {
		u := following[i]
		events, err := api.GetUserEvents(ctx, token, u.Login, 10)
		if err != nil {
			fmt.Printf("  %s: (could not fetch events)\n", u.Login)
			continue
		}
		score := estimateHabitScoreFromEvents(events)
		fmt.Printf("  %s: ~%d pts (from %d recent events)\n", u.Login, score, len(events))
		for j, ev := range events {
			if j >= 3 {
				break
			}
			fmt.Printf("    - %s: %s\n", ev.GetType(), ev.GetRepo().GetName())
		}
		fmt.Println()
	}

	return nil
}

// estimateHabitScoreFromEvents approximates habit score from event types.
func estimateHabitScoreFromEvents(events []*github.Event) int {
	now := time.Now().Local()
	today := now.Format("2006-01-02")
	var score int
	for _, ev := range events {
		created := ev.GetCreatedAt().Time
		if created.Local().Format("2006-01-02") != today {
			continue
		}
		switch ev.GetType() {
		case "PushEvent":
			score += 10
		case "PullRequestEvent":
			score += 15
		case "PullRequestReviewEvent":
			score += 12
		case "IssuesEvent":
			score += 5
		default:
			// Count as minimal activity
			score += 5
		}
	}
	return score
}
