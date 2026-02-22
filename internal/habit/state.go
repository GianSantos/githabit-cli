package habit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
)

const stateFilename = "githabit/state.json"
// 1 hour cache valid duration
const cacheValidDuration = 1 * time.Hour

// State holds cached data and last check time.
type State struct {
	LastChecked     time.Time `json:"last_checked"`
	TodayScore      int      `json:"today_score"`
	CurrentStreak   int      `json:"current_streak"`
	TodayBreakdown  struct {
		Commits int `json:"commits"`
		PRs     int `json:"prs"`
		Reviews int `json:"reviews"`
		Issues  int `json:"issues"`
	} `json:"today_breakdown"`
}

func StatePath() (string, error) {
	dir, err := xdg.StateFile(stateFilename)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func LoadState() (*State, error) {
	path, err := StatePath()
	if err != nil {
		return nil, err
	}
	return LoadStateFromPath(path)
}

// LoadStateFromPath reads state from the given path. Used for testing.
func LoadStateFromPath(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{}, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse state: %w", err)
	}
	return &s, nil
}

func SaveState(s *State) error {
	path, err := StatePath()
	if err != nil {
		return err
	}
	return SaveStateToPath(s, path)
}

// SaveStateToPath writes state to the given path. Used for testing.
func SaveStateToPath(s *State, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func IsCacheValid(s *State) bool {
	if s.LastChecked.IsZero() {
		return false
	}
	return time.Since(s.LastChecked) < cacheValidDuration
}
