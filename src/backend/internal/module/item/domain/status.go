package domain

import "slices"

type Status string

const (
	StatusDraft      Status = "draft"
	StatusModeration Status = "moderation"
	StatusPublished  Status = "published"
	StatusArchived   Status = "archived"
)

var transitions = map[Status][]Status{
	StatusDraft:      {StatusModeration, StatusArchived},
	StatusModeration: {StatusPublished, StatusDraft, StatusArchived},
	StatusPublished:  {StatusArchived},
	StatusArchived:   {StatusDraft},
}

func (s Status) Valid() bool {
	_, ok := transitions[s]

	return ok
}

func (s Status) CanTransitionTo(target Status) bool {
	return slices.Contains(transitions[s], target)
}

func (s Status) String() string {
	return string(s)
}
