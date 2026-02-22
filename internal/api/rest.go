package api

import (
	"context"
	"fmt"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

// FollowingUser holds basic info about a user we follow.
type FollowingUser struct {
	Login string
}

// FeedEvent represents an activity event for the feed.
type FeedEvent struct {
	Type      string
	Actor     string
	Repo      string
	CreatedAt string
}

// GetCurrentUser returns the authenticated user's login.
func GetCurrentUser(ctx context.Context, token string) (string, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, ts)
	client := github.NewClient(httpClient)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return "", fmt.Errorf("get current user: %w", err)
	}
	return user.GetLogin(), nil
}

// GetFollowing returns the list of users the current user is following.
func GetFollowing(ctx context.Context, token, login string) ([]FollowingUser, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, ts)
	client := github.NewClient(httpClient)

	opt := &github.ListOptions{PerPage: 100}
	var all []FollowingUser
	for {
		users, resp, err := client.Users.ListFollowing(ctx, login, opt)
		if err != nil {
			return nil, fmt.Errorf("list following: %w", err)
		}
		for _, u := range users {
			all = append(all, FollowingUser{Login: u.GetLogin()})
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return all, nil
}

// GetUserEvents returns public events for a user (used for feed/competition).
func GetUserEvents(ctx context.Context, token, login string, count int) ([]*github.Event, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, ts)
	client := github.NewClient(httpClient)

	events, _, err := client.Activity.ListEventsPerformedByUser(ctx, login, false, &github.ListOptions{
		PerPage: count,
	})
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	return events, nil
}
