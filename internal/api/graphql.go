package api

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// DayContributions holds the contribution counts for a single day.
type DayContributions struct {
	Date     time.Time
	Commits  int
	PRs      int
	Reviews  int
	Issues   int
}

// ContributionsQuery fetches contribution totals for a date range.
// GitHub's contributionsCollection uses UTC, so from/to must be UTC boundaries.
func ContributionsQuery(ctx context.Context, token, login string, from, to time.Time) (*DayContributions, error) {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(httpClient)

	var query struct {
		User struct {
			ContributionsCollection struct {
				TotalCommitContributions            int
				TotalPullRequestContributions       int
				TotalPullRequestReviewContributions int
				TotalIssueContributions             int
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(login),
		"from":  githubv4.DateTime{Time: from},
		"to":    githubv4.DateTime{Time: to},
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("graphql query: %w", err)
	}

	return &DayContributions{
		Date:     from.Truncate(24 * time.Hour),
		Commits:  query.User.ContributionsCollection.TotalCommitContributions,
		PRs:      query.User.ContributionsCollection.TotalPullRequestContributions,
		Reviews:  query.User.ContributionsCollection.TotalPullRequestReviewContributions,
		Issues:   query.User.ContributionsCollection.TotalIssueContributions,
	}, nil
}

// ContributionCalendarWeek represents a week in the contribution calendar.
type ContributionCalendarWeek struct {
	ContributionDays []struct {
		Date             string
		ContributionCount int
	}
}

// StreakQuery fetches the contribution calendar for the last 30 days.
func StreakQuery(ctx context.Context, token, login string, from, to time.Time) ([]ContributionCalendarWeek, error) {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(httpClient)

	var query struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					Weeks []ContributionCalendarWeek
				}
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(login),
		"from":  githubv4.DateTime{Time: from},
		"to":    githubv4.DateTime{Time: to},
	}

	if err := client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("graphql query: %w", err)
	}

	return query.User.ContributionsCollection.ContributionCalendar.Weeks, nil
}
