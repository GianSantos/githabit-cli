package ui

import (
	"fmt"
	"time"

	"github.com/GianSantos/githabit-cli/internal/api"
)

// RenderStreakGrid builds a 7x5 ANSI grid for the last 30 days.
// Columns = days of week (Sun-Sat), rows = weeks (5).
func RenderStreakGrid(weeks []api.ContributionCalendarWeek) string {
	// Flatten contribution days and index by date
	dayMap := make(map[string]int)
	for _, w := range weeks {
		for _, d := range w.ContributionDays {
			if d.Date != "" {
				dayMap[d.Date] = d.ContributionCount
			}
		}
	}

	// Determine date range: last 30 days
	now := time.Now().Local()
	start := now.AddDate(0, 0, -29)

	var buf string
	buf += TitleStyle.Render("Last 30 days") + "\n"
	buf += "Sun Mon Tue Wed Thu Fri Sat\n"

	// Build grid: 7 columns (Sun-Sat), 5 rows
	// Each row = one week; columns align with weekday
	grid := make([][]int, 5)
	for i := range grid {
		grid[i] = make([]int, 7)
	}

	for offset := 0; offset < 30; offset++ {
		d := start.AddDate(0, 0, offset)
		key := d.Format("2006-01-02")
		count := dayMap[key]
		level := contributionLevel(count)
		startWeekday := int(start.Weekday()) // 0=Sun, 6=Sat
		cellIdx := offset + startWeekday
		row := cellIdx / 7
		col := cellIdx % 7
		if row < 5 {
			grid[row][col] = level
		}
	}

	for _, r := range grid {
		for _, level := range r {
			buf += StreakCell(level) + "   "
		}
		buf += "\n"
	}
	return buf
}

func contributionLevel(count int) int {
	switch {
	case count == 0:
		return 0
	case count <= 3:
		return 1
	case count <= 6:
		return 2
	case count <= 9:
		return 3
	default:
		return 4
	}
}

func FormatScore(score int, commits, prs, reviews, issues int) string {
	s := fmt.Sprintf("Habit Score: %d pts", score)
	if commits+prs+reviews+issues > 0 {
		s += fmt.Sprintf(" (Commits: %d, PRs: %d, Reviews: %d, Issues: %d)",
			commits, prs, reviews, issues)
	}
	return s
}

func FormatStreak(days int) string {
	if days == 0 {
		return "Current streak: 0 days"
	}
	return fmt.Sprintf("Current streak: %d day%s", days, plural(days))
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
