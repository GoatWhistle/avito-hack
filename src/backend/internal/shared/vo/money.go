package vo

import (
	"errors"
	"fmt"
)

var ErrNegativeMoney = errors.New("money: negative amount")

type Money struct {
	kopeks int64
}

func NewMoney(kopeks int64) (Money, error) {
	if kopeks < 0 {
		return Money{}, ErrNegativeMoney
	}

	return Money{kopeks: kopeks}, nil
}

func MustMoney(kopeks int64) Money {
	m, err := NewMoney(kopeks)
	if err != nil {
		panic(err)
	}

	return m
}

func ZeroMoney() Money {
	return Money{}
}

func (m Money) Kopeks() int64 {
	return m.kopeks
}

func (m Money) IsZero() bool {
	return m.kopeks == 0
}

func (m Money) Add(other Money) Money {
	return Money{kopeks: m.kopeks + other.kopeks}
}

func (m Money) Equal(other Money) bool {
	return m.kopeks == other.kopeks
}

func (m Money) String() string {
	return fmt.Sprintf("%d.%02d", m.kopeks/100, m.kopeks%100)
}
