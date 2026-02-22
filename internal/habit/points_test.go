package habit

import (
	"testing"

	"github.com/GianSantos/githabit-cli/internal/api"
)

func TestScoreFromContributions(t *testing.T) {
	tests := []struct {
		name   string
		input  *api.DayContributions
		want   int
	}{
		{
			name: "empty",
			input: &api.DayContributions{Commits: 0, PRs: 0, Reviews: 0, Issues: 0},
			want:  0,
		},
		{
			name:  "single commit",
			input: &api.DayContributions{Commits: 1, PRs: 0, Reviews: 0, Issues: 0},
			want:  10,
		},
		{
			name:  "single PR",
			input: &api.DayContributions{Commits: 0, PRs: 1, Reviews: 0, Issues: 0},
			want:  15,
		},
		{
			name:  "single review",
			input: &api.DayContributions{Commits: 0, PRs: 0, Reviews: 1, Issues: 0},
			want:  12,
		},
		{
			name:  "single issue",
			input: &api.DayContributions{Commits: 0, PRs: 0, Reviews: 0, Issues: 1},
			want:  5,
		},
		{
			name: "mixed contributions",
			input: &api.DayContributions{Commits: 2, PRs: 1, Reviews: 1, Issues: 2},
			want:  2*10 + 15 + 12 + 2*5, // 20 + 15 + 12 + 10 = 57
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScoreFromContributions(tt.input); got != tt.want {
				t.Errorf("ScoreFromContributions() = %d, want %d", got, tt.want)
			}
		})
	}
}
