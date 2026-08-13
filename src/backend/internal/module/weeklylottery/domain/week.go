package domain

import "time"

var moscow = time.FixedZone("Europe/Moscow", 3*60*60)

type Week struct {
	Start time.Time
	End   time.Time
	Key   time.Time
}

func WeekAt(now time.Time) Week {
	local := now.In(moscow)
	daysFromMonday := (int(local.Weekday()) + 6) % 7
	start := time.Date(local.Year(), local.Month(), local.Day()-daysFromMonday, 0, 0, 0, 0, moscow)

	key := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	return Week{Start: start, End: start.AddDate(0, 0, 7), Key: key}
}
