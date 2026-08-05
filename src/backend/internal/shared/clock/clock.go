package clock

import "time"

type Clock struct{}

func New() Clock {
	return Clock{}
}

func (Clock) Now() time.Time {
	return time.Now().UTC()
}

type Fixed struct {
	value time.Time
}

func NewFixed(value time.Time) Fixed {
	return Fixed{value: value}
}

func (f Fixed) Now() time.Time {
	return f.value
}
