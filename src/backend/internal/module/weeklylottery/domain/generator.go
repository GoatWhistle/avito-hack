package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	BoardSize      = 9
	lossWeight     = 50
	maxFillerCount = 2
)

type Random interface {
	Intn(limit int) (int, error)
}

type CryptoRandom struct{}

func (CryptoRandom) Intn(limit int) (int, error) {
	if limit <= 0 {
		return 0, fmt.Errorf("random limit must be positive")
	}

	value, err := rand.Int(rand.Reader, big.NewInt(int64(limit)))
	if err != nil {
		return 0, fmt.Errorf("draw random number: %w", err)
	}

	return int(value.Int64()), nil
}

type Generator struct {
	random Random
}

func NewGenerator(random Random) *Generator {
	return &Generator{random: random}
}

func (g *Generator) Generate() ([BoardSize]Symbol, *Prize, error) {
	var board [BoardSize]Symbol
	prize, err := g.drawPrize()
	if err != nil {
		return board, nil, err
	}

	filled := 0
	limits := make(map[Symbol]int, len(prizes))
	for _, symbol := range allSymbols() {
		limits[symbol] = maxFillerCount
	}
	if prize != nil {
		board[0], board[1], board[2] = prize.Symbol, prize.Symbol, prize.Symbol
		filled = 3
		delete(limits, prize.Symbol)
	}

	for index := filled; index < BoardSize; index++ {
		symbol, err := g.drawFiller(limits)
		if err != nil {
			return board, nil, err
		}
		board[index] = symbol
		limits[symbol]--
	}

	if err := g.shuffle(&board); err != nil {
		return board, nil, err
	}
	prizeID := ""
	if prize != nil {
		prizeID = prize.ID
	}
	if err := ValidateBoard(board, prizeID); err != nil {
		return board, nil, fmt.Errorf("validate generated weekly lottery board: %w", err)
	}

	return board, prize, nil
}

func (g *Generator) drawPrize() (*Prize, error) {
	total := lossWeight
	for _, prize := range prizes {
		total += prize.Weight
	}

	draw, err := g.random.Intn(total)
	if err != nil {
		return nil, err
	}
	if draw < lossWeight {
		return nil, nil
	}

	draw -= lossWeight
	for _, prize := range prizes {
		if draw < prize.Weight {
			selected := prize
			return &selected, nil
		}
		draw -= prize.Weight
	}

	return nil, fmt.Errorf("prize draw is outside configured weights")
}

func (g *Generator) shuffle(board *[BoardSize]Symbol) error {
	for index := BoardSize - 1; index > 0; index-- {
		other, err := g.random.Intn(index + 1)
		if err != nil {
			return err
		}
		board[index], board[other] = board[other], board[index]
	}

	return nil
}

func allSymbols() []Symbol {
	return []Symbol{SymbolBicycle, SymbolSmartphone, SymbolSofa, SymbolSneakers, SymbolDelivery, SymbolPromotion}
}

func (g *Generator) drawFiller(limits map[Symbol]int) (Symbol, error) {
	candidates := make([]Symbol, 0, len(limits))
	for _, symbol := range allSymbols() {
		if limits[symbol] > 0 {
			candidates = append(candidates, symbol)
		}
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("weekly lottery board has no filler symbol left")
	}

	index, err := g.random.Intn(len(candidates))
	if err != nil {
		return "", err
	}

	return candidates[index], nil
}
