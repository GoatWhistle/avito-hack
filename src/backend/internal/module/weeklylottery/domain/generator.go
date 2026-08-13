package domain

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	BoardSize  = 9
	lossWeight = 50
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

	symbols := allSymbols()
	if prize != nil {
		board[0], board[1], board[2] = prize.Symbol, prize.Symbol, prize.Symbol
		symbols = without(symbols, prize.Symbol)
		for index := 3; index < BoardSize; index++ {
			board[index] = symbols[(index-3)/2]
		}
	} else {
		for index := range BoardSize {
			board[index] = symbols[(index/2)%len(symbols)]
		}
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

func without(symbols []Symbol, excluded Symbol) []Symbol {
	out := make([]Symbol, 0, len(symbols)-1)
	for _, symbol := range symbols {
		if symbol != excluded {
			out = append(out, symbol)
		}
	}

	return out
}
