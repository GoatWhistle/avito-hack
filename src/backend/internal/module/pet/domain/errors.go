package domain

import "errors"

var (
	ErrDuplicateAction = errors.New("action has already been rewarded")
	ErrLimitReached    = errors.New("action reward limit reached")
	ErrConditionNotMet = errors.New("reward condition is not met")
	ErrInvalidAction   = errors.New("invalid action data")
)
