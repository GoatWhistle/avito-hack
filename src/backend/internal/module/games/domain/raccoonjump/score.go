package raccoonjump

import "time"

const (
	gravity        = 0.45
	springVelocity = 21.0
	framesPerSec   = 60.0
	pixelsPerScore = 10.0
)

const safetyFactor = 3.0

const graceSeconds = 2.0

const maxScorePerSecond = int(
	(springVelocity / 2.0) * framesPerSec / pixelsPerScore * safetyFactor,
)

const (
	MaxScore       = 100000
	MinStreakScore = 25
)

func MaxPlausibleScore(elapsed time.Duration) int {
	seconds := elapsed.Seconds()
	if seconds < 0 {
		seconds = 0
	}

	limit := (seconds + graceSeconds) * float64(maxScorePerSecond)

	if limit >= float64(MaxScore) {
		return MaxScore
	}

	return int(limit)
}

func CountsTowardStreak(score int) bool {
	return score >= MinStreakScore
}
