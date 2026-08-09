package domain

import "slices"

type Status string

const (
	StatusDraft      Status = "draft"
	StatusModeration Status = "moderation"
	StatusPublished  Status = "published"
	StatusSold       Status = "sold"
	StatusArchived   Status = "archived"
)

var PublicStatuses = []Status{StatusPublished, StatusSold}

var transitions = map[Status][]Status{
	StatusDraft:      {StatusPublished, StatusModeration, StatusArchived},
	StatusModeration: {StatusPublished, StatusDraft, StatusArchived},
	StatusPublished:  {StatusSold, StatusArchived},
	StatusSold:       {},
	StatusArchived:   {StatusDraft},
}

func (s Status) Valid() bool {
	_, ok := transitions[s]

	return ok
}

func (s Status) CanTransitionTo(target Status) bool {
	return slices.Contains(transitions[s], target)
}

func (s Status) IsTerminal() bool {
	return len(transitions[s]) == 0
}

func (s Status) IsPublic() bool {
	return slices.Contains(PublicStatuses, s)
}

func (s Status) String() string {
	return string(s)
}
