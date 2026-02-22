package habit

import (
	"time"
)

// LocalMidnight returns the start of the given date in the local timezone.
func LocalMidnight(t time.Time) time.Time {
	y, m, d := t.Local().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// LocalDateKey returns "YYYY-MM-DD" for the given time in local timezone.
func LocalDateKey(t time.Time) string {
	return t.Local().Format("2006-01-02")
}

// UTCBoundsForLocalDay returns (from, to) in UTC for the entire local day.
// The day is defined by the given time's local date.
func UTCBoundsForLocalDay(t time.Time) (from, to time.Time) {
	midnight := LocalMidnight(t)
	from = midnight.UTC()
	to = midnight.Add(24*time.Hour - time.Nanosecond).UTC()
	return from, to
}

// HoursAgo returns now minus the given hours.
func HoursAgo(hours int) time.Time {
	return time.Now().Add(-time.Duration(hours) * time.Hour)
}
