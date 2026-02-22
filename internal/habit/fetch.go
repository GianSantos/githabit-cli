package habit

import (
	"context"
	"time"

	"github.com/GianSantos/githabit-cli/internal/api"
)

// FetchTodayScore fetches today's habit score from the API.
// Uses UTC bounds for the user's local "today".
func FetchTodayScore(ctx context.Context, token, login string) (*api.DayContributions, int, error) {
	from, to := UTCBoundsForLocalDay(time.Now())
	dc, err := api.ContributionsQuery(ctx, token, login, from, to)
	if err != nil {
		return nil, 0, err
	}
	score := ScoreFromContributions(dc)
	return dc, score, nil
}

// FetchStreakData fetches contribution calendar for the last 30 days.
func FetchStreakData(ctx context.Context, token, login string) ([]api.ContributionCalendarWeek, error) {
	now := time.Now()
	const thirtyDaysAgo = -30 * 24 * time.Hour
	from := LocalMidnight(now).Add(thirtyDaysAgo)
	to := now
	return api.StreakQuery(ctx, token, login, from, to)
}

// ComputeStreakFromCalendar returns the current streak (consecutive days with contributions ending today).
func ComputeStreakFromCalendar(weeks []api.ContributionCalendarWeek) int {
	dayMap := make(map[string]int)
	for _, w := range weeks {
		for _, d := range w.ContributionDays {
			if d.Date != "" {
				dayMap[d.Date] = d.ContributionCount
			}
		}
	}
	now := time.Now().Local()
	for i := 0; i < 365; i++ {
		d := now.AddDate(0, 0, -i)
		key := d.Format("2006-01-02")
		if dayMap[key] == 0 {
			return i
		}
	}
	return 365
}
