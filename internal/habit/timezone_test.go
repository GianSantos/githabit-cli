package habit

import (
	"testing"
	"time"
)

func TestLocalMidnight(t *testing.T) {
	// 2025-02-22 14:30:00 in local time
	tm := time.Date(2025, 2, 22, 14, 30, 0, 0, time.Local)
	got := LocalMidnight(tm)
	want := time.Date(2025, 2, 22, 0, 0, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Errorf("LocalMidnight() = %v, want %v", got, want)
	}
}

func TestLocalDateKey(t *testing.T) {
	tm := time.Date(2025, 2, 22, 14, 30, 0, 0, time.Local)
	got := LocalDateKey(tm)
	want := "2025-02-22"
	if got != want {
		t.Errorf("LocalDateKey() = %q, want %q", got, want)
	}
}

func TestUTCBoundsForLocalDay(t *testing.T) {
	// Use a fixed local time
	tm := time.Date(2025, 2, 22, 12, 0, 0, 0, time.Local)
	from, to := UTCBoundsForLocalDay(tm)

	// from should be midnight local in UTC
	midnight := time.Date(2025, 2, 22, 0, 0, 0, 0, time.Local)
	wantFrom := midnight.UTC()
	if !from.Equal(wantFrom) {
		t.Errorf("from = %v, want %v", from, wantFrom)
	}

	// to should be 23:59:59.999999999 local in UTC
	endOfDay := time.Date(2025, 2, 22, 23, 59, 59, 999999999, time.Local)
	wantTo := endOfDay.UTC()
	if !to.Equal(wantTo) {
		t.Errorf("to = %v, want %v", to, wantTo)
	}

	// Span should be ~24 hours
	dur := to.Sub(from)
	if dur < 23*time.Hour || dur > 25*time.Hour {
		t.Errorf("UTCBoundsForLocalDay span = %v, expected ~24h", dur)
	}
}

func TestHoursAgo(t *testing.T) {
	before := time.Now()
	got := HoursAgo(2)
	after := time.Now()
	if got.After(before) || got.Before(after.Add(-3*time.Hour)) {
		t.Errorf("HoursAgo(2) = %v, expected ~2h ago", got)
	}
}
