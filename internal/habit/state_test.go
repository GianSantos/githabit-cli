package habit

import (
	"path/filepath"
	"testing"
	"time"
)

func TestLoadStateFromPath_MissingFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "githabit", "state.json")
	// Don't create the file - simulate first run

	got, err := LoadStateFromPath(path)
	if err != nil {
		t.Fatalf("LoadStateFromPath() error = %v", err)
	}
	if got == nil {
		t.Fatal("LoadStateFromPath() returned nil state")
	}
	if got.LastChecked != (time.Time{}) || got.TodayScore != 0 {
		t.Errorf("LoadStateFromPath() with missing file should return zero state, got %+v", got)
	}
}

func TestSaveStateToPath_LoadStateFromPath_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "githabit", "state.json")

	state := &State{
		LastChecked:   time.Date(2025, 2, 22, 14, 0, 0, 0, time.UTC),
		TodayScore:    42,
		CurrentStreak: 7,
	}
	state.TodayBreakdown.Commits = 2
	state.TodayBreakdown.PRs = 1
	state.TodayBreakdown.Reviews = 0
	state.TodayBreakdown.Issues = 1

	if err := SaveStateToPath(state, path); err != nil {
		t.Fatalf("SaveStateToPath() error = %v", err)
	}

	loaded, err := LoadStateFromPath(path)
	if err != nil {
		t.Fatalf("LoadStateFromPath() after save error = %v", err)
	}
	if loaded.TodayScore != state.TodayScore {
		t.Errorf("TodayScore = %d, want %d", loaded.TodayScore, state.TodayScore)
	}
	if loaded.CurrentStreak != state.CurrentStreak {
		t.Errorf("CurrentStreak = %d, want %d", loaded.CurrentStreak, state.CurrentStreak)
	}
	if loaded.TodayBreakdown.Commits != state.TodayBreakdown.Commits {
		t.Errorf("TodayBreakdown.Commits = %d, want %d", loaded.TodayBreakdown.Commits, state.TodayBreakdown.Commits)
	}
}

func TestIsCacheValid(t *testing.T) {
	tests := []struct {
		name   string
		state  *State
		want   bool
	}{
		{
			name:  "zero time is invalid",
			state: &State{LastChecked: time.Time{}},
			want:  false,
		},
		{
			name:  "recent is valid",
			state: &State{LastChecked: time.Now().Add(-10 * time.Minute)},
			want:  true,
		},
		{
			name:  "old is invalid",
			state: &State{LastChecked: time.Now().Add(-2 * time.Hour)},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCacheValid(tt.state); got != tt.want {
				t.Errorf("IsCacheValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
