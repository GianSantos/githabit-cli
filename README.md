# GitHabit

A habit-building CLI tool that tracks GitHub contributions using a custom "Habit Score" to gamify coding practice.

## Install

```bash
go build -o githabit .
# or: go install
```

## Quick Start

1. **Initialize** (prompts for GitHub PAT, validates scopes, saves to keyring):
   ```bash
   githabit init
   ```

2. **Check status** (today's score and streak):
   ```bash
   githabit status
   ```

3. **View streak grid** (30-day ANSI color grid):
   ```bash
   githabit streak
   ```

4. **See followed users' activity**:
   ```bash
   githabit feed
   ```

5. **Schedule 8 PM daily reminder** (crontab on Unix, schtasks on Windows):
   ```bash
   githabit notify start
   ```

## Habit Score

| Activity | Points |
|----------|--------|
| Commit | 10 |
| Pull Request | 15 |
| PR Review | 12 |
| Issue Created | 5 |

## Requirements

- GitHub Personal Access Token with `repo`, `read:user`, `read:org` scopes
- Go 1.21+
