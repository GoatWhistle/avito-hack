package domain

import "maps"

type Attributes map[string]string

func NewAttributes(raw map[string]string) Attributes {
	if len(raw) == 0 {
		return Attributes{}
	}

	return Attributes(maps.Clone(raw))
}

func (a Attributes) Clone() Attributes {
	if a == nil {
		return Attributes{}
	}

	return Attributes(maps.Clone(a))
}

func (a Attributes) Get(key string) (string, bool) {
	value, ok := a[key]

	return value, ok
}

func (a Attributes) Len() int {
	return len(a)
}
