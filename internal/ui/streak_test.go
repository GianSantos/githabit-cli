package ui

import (
	"strings"
	"testing"

	"github.com/GianSantos/githabit-cli/internal/api"
)

func TestFormatScore(t *testing.T) {
	tests := []struct {
		name    string
		score   int
		commits int
		prs     int
		reviews int
		issues  int
		want    string
	}{
		{
			name:    "zero with no breakdown",
			score:   0,
			commits: 0, prs: 0, reviews: 0, issues: 0,
			want: "Habit Score: 0 pts",
		},
		{
			name:    "with breakdown",
			score:   42,
			commits: 2, prs: 1, reviews: 0, issues: 1,
			want: "Habit Score: 42 pts (Commits: 2, PRs: 1, Reviews: 0, Issues: 1)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatScore(tt.score, tt.commits, tt.prs, tt.reviews, tt.issues)
			if got != tt.want {
				t.Errorf("FormatScore() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatStreak(t *testing.T) {
	tests := []struct {
		name  string
		days  int
		want  string
	}{
		{"zero", 0, "Current streak: 0 days"},
		{"one", 1, "Current streak: 1 day"},
		{"plural", 7, "Current streak: 7 days"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatStreak(tt.days)
			if got != tt.want {
				t.Errorf("FormatStreak(%d) = %q, want %q", tt.days, got, tt.want)
			}
		})
	}
}

func TestRenderStreakGrid(t *testing.T) {
	weeks := []api.ContributionCalendarWeek{
		{
			ContributionDays: []struct {
				Date             string
				ContributionCount int
			}{
				{Date: "2025-02-20", ContributionCount: 3},
				{Date: "2025-02-21", ContributionCount: 0},
				{Date: "2025-02-22", ContributionCount: 5},
			},
		},
	}

	got := RenderStreakGrid(weeks)

	if !strings.Contains(got, "Last 30 days") {
		t.Error("RenderStreakGrid should contain 'Last 30 days'")
	}
	if !strings.Contains(got, "Sun Mon Tue Wed Thu Fri Sat") {
		t.Error("RenderStreakGrid should contain weekday headers")
	}
	// Should have 5 rows of data
	lines := strings.Split(got, "\n")
	if len(lines) < 6 {
		t.Errorf("RenderStreakGrid should have header + 5 rows, got %d lines", len(lines))
	}
}
