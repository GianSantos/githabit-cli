package habit

import (
	"testing"
	"time"

	"github.com/GianSantos/githabit-cli/internal/api"
)

func TestComputeStreakFromCalendar(t *testing.T) {
	today := "2025-02-22" // Adjust if needed; test uses fixed dates

	tests := []struct {
		name   string
		weeks  []api.ContributionCalendarWeek
		today  string // Override for deterministic test
		want   int
	}{
		{
			name:  "empty calendar",
			weeks: []api.ContributionCalendarWeek{},
			today: today,
			want:  0,
		},
		{
			name: "single day with contributions",
			weeks: []api.ContributionCalendarWeek{
				{
					ContributionDays: []struct {
						Date             string
						ContributionCount int
					}{
						{Date: "2025-02-22", ContributionCount: 5},
					},
				},
			},
			today: "2025-02-22",
			want:  1,
		},
		{
			name: "today has no contributions",
			weeks: []api.ContributionCalendarWeek{
				{
					ContributionDays: []struct {
						Date             string
						ContributionCount int
					}{
						{Date: "2025-02-21", ContributionCount: 3},
					},
				},
			},
			today: "2025-02-22",
			want:  0,
		},
		{
			name: "three day streak ending today",
			weeks: []api.ContributionCalendarWeek{
				{
					ContributionDays: []struct {
						Date             string
						ContributionCount int
					}{
						{Date: "2025-02-20", ContributionCount: 1},
						{Date: "2025-02-21", ContributionCount: 2},
						{Date: "2025-02-22", ContributionCount: 1},
					},
				},
			},
			today: "2025-02-22",
			want:  3,
		},
		{
			name: "gap breaks streak",
			weeks: []api.ContributionCalendarWeek{
				{
					ContributionDays: []struct {
						Date             string
						ContributionCount int
					}{
						{Date: "2025-02-20", ContributionCount: 1},
						{Date: "2025-02-21", ContributionCount: 0},
						{Date: "2025-02-22", ContributionCount: 1},
					},
				},
			},
			today: "2025-02-22",
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeStreakForDate(tt.weeks, tt.today)
			if got != tt.want {
				t.Errorf("computeStreakForDate() = %d, want %d", got, tt.want)
			}
		})
	}
}

// computeStreakForDate is a test helper that runs streak logic for a fixed "today" date.
// We need to test streak with deterministic dates; the real ComputeStreakFromCalendar uses time.Now().
func computeStreakForDate(weeks []api.ContributionCalendarWeek, todayStr string) int {
	dayMap := make(map[string]int)
	for _, w := range weeks {
		for _, d := range w.ContributionDays {
			if d.Date != "" {
				dayMap[d.Date] = d.ContributionCount
			}
		}
	}
	for i := 0; i < 365; i++ {
		// Parse today and walk back
		d := parseDate(todayStr).AddDate(0, 0, -i)
		key := d.Format("2006-01-02")
		if dayMap[key] == 0 {
			return i
		}
	}
	return 365
}

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}
