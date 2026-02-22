package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/GianSantos/githabit-cli/internal/api"
	"github.com/GianSantos/githabit-cli/internal/auth"
	"github.com/GianSantos/githabit-cli/internal/habit"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(notifyCmd)
}

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Manage background reminder (crontab on Unix, schtasks on Windows)",
	Long:  "Subcommands: start (schedule 8 PM reminder), stop (remove schedule).",
}

var notifyStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Schedule daily 8 PM reminder",
	RunE:  runNotifyStart,
}

var notifyStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Remove scheduled reminder",
	RunE:  runNotifyStop,
}

var checkReminderSilent bool

var checkReminderCmd = &cobra.Command{
	Use:   "check-reminder",
	Short: "Internal: Check if reminder should fire (used by scheduler)",
	RunE:  runCheckReminder,
	Hidden: true,
}

func init() {
	notifyCmd.AddCommand(notifyStartCmd, notifyStopCmd)
	rootCmd.AddCommand(checkReminderCmd)
	checkReminderCmd.Flags().BoolVar(&checkReminderSilent, "silent", false, "Suppress errors (for cron)")
}

func runNotifyStart(cmd *cobra.Command, args []string) error {
	path, err := exec.LookPath("githabit")
	if err != nil {
		path, _ = os.Executable()
	}
	if path == "" {
		return fmt.Errorf("could not find githabit executable")
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		return addCrontabEntry(path)
	case "windows":
		return addSchtaskEntry(path)
	default:
		return fmt.Errorf("notify not supported on %s", runtime.GOOS)
	}
}

func runNotifyStop(cmd *cobra.Command, args []string) error {
	switch runtime.GOOS {
	case "linux", "darwin":
		return removeCrontabEntry()
	case "windows":
		return removeSchtaskEntry()
	default:
		return fmt.Errorf("notify not supported on %s", runtime.GOOS)
	}
}

func addCrontabEntry(githabitPath string) error {
	// 0 20 * * * = 8:00 PM daily
	entry := fmt.Sprintf("0 20 * * * %s check-reminder --silent\n", githabitPath)
	tmp, err := os.CreateTemp("", "githabit-cron-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	// Read existing crontab
	out, err := exec.Command("crontab", "-l").CombinedOutput()
	existing := string(out)
	if err != nil {
		// crontab -l returns 1 when no crontab
		existing = ""
	}
	// Avoid duplicate
	if !strings.Contains(existing, "githabit check-reminder") {
		_, _ = tmp.WriteString(existing)
		_, _ = tmp.WriteString(entry)
	}
	tmp.Close()
	return exec.Command("crontab", tmp.Name()).Run()
}

func removeCrontabEntry() error {
	out, err := exec.Command("crontab", "-l").CombinedOutput()
	if err != nil {
		return nil // no crontab
	}
	lines := strings.Split(string(out), "\n")
	var kept []string
	for _, l := range lines {
		if !strings.Contains(l, "githabit check-reminder") {
			kept = append(kept, l)
		}
	}
	content := strings.Join(kept, "\n")
	if strings.TrimSpace(content) == "" {
		return exec.Command("crontab", "-r").Run()
	}
	tmp, _ := os.CreateTemp("", "githabit-cron-*")
	defer os.Remove(tmp.Name())
	tmp.WriteString(content)
	if !strings.HasSuffix(content, "\n") {
		tmp.WriteString("\n")
	}
	tmp.Close()
	return exec.Command("crontab", tmp.Name()).Run()
}

func addSchtaskEntry(githabitPath string) error {
	// /sc minute:hour - 20:00 = 8 PM
	cmd := exec.Command("schtasks", "/create", "/tn", "GitHabitReminder",
		"/tr", githabitPath+" check-reminder --silent",
		"/sc", "daily",
		"/st", "20:00",
		"/f")
	return cmd.Run()
}

func removeSchtaskEntry() error {
	return exec.Command("schtasks", "/delete", "/tn", "GitHabitReminder", "/f").Run()
}

func runCheckReminder(cmd *cobra.Command, args []string) error {
	silent := checkReminderSilent

	state, err := habit.LoadState()
	if err != nil {
		return err
	}

	// If cached score is 0, fetch fresh to avoid false positives
	if state.TodayScore == 0 {
		token, err := auth.GetToken()
		if err == nil {
			login, err := api.GetCurrentUser(context.Background(), token)
			if err == nil {
				dc, score, err := habit.FetchTodayScore(context.Background(), token, login)
				if err == nil {
					state.TodayScore = score
					state.TodayBreakdown.Commits = dc.Commits
					state.TodayBreakdown.PRs = dc.PRs
					state.TodayBreakdown.Reviews = dc.Reviews
					state.TodayBreakdown.Issues = dc.Issues
				}
			}
		}
	}

	if state.TodayScore == 0 {
		quote := "The secret of getting ahead is getting started."
		if err := beeep.Alert("GitHabit Reminder", quote, ""); err != nil && !silent {
			return err
		}
	}
	return nil
}
