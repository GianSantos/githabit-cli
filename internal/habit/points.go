package habit

import "github.com/GianSantos/githabit-cli/internal/api"

// Point weights per the spec.
const (
	PointsCommit   = 10
	PointsPR       = 15
	PointsReview   = 12
	PointsIssue    = 5
)

// ScoreFromContributions computes the Habit Score from daily contributions.
func ScoreFromContributions(d *api.DayContributions) int {
	return d.Commits*PointsCommit +
		d.PRs*PointsPR +
		d.Reviews*PointsReview +
		d.Issues*PointsIssue
}
