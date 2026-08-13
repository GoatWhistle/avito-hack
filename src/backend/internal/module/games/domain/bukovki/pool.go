package bukovki

import (
	"context"
	"fmt"
	"math/rand"
)

type MemoryWordPool struct {
	words []string
	set   map[string]bool
}

func NewMemoryWordPool() *MemoryWordPool {
	list := []string{
		"apple", "ghost", "brain", "house", "train",
		"mouse", "plant", "water", "earth", "light",
	}

	set := make(map[string]bool, len(list))
	for _, w := range list {
		set[w] = true
	}

	return &MemoryWordPool{
		words: list,
		set:   set,
	}
}

func (p *MemoryWordPool) RandomWord(ctx context.Context, length int) (string, error) {
	if len(p.words) == 0 {
		return "", fmt.Errorf("no words available in pool")
	}
	idx := rand.Intn(len(p.words))
	return p.words[idx], nil
}

func (p *MemoryWordPool) IsValidWord(ctx context.Context, word string) bool {
	return len(word) == wordLength
}
