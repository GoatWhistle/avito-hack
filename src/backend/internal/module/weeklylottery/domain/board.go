package domain

import "fmt"

func ValidateBoard(board [BoardSize]Symbol, prizeID string) error {
	counts := make(map[Symbol]int, len(prizes))
	valid := make(map[Symbol]struct{}, len(prizes))
	for _, prize := range prizes {
		valid[prize.Symbol] = struct{}{}
	}

	for _, symbol := range board {
		if _, ok := valid[symbol]; !ok {
			return fmt.Errorf("weekly lottery board contains unknown symbol %q", symbol)
		}
		counts[symbol]++
	}

	if prizeID == "" {
		for _, count := range counts {
			if count >= 3 {
				return fmt.Errorf("losing weekly lottery board contains a triple")
			}
		}

		return nil
	}

	prize, ok := PrizeByID(prizeID)
	if !ok {
		return fmt.Errorf("weekly lottery board references unknown prize %q", prizeID)
	}
	if counts[prize.Symbol] != 3 {
		return fmt.Errorf("weekly lottery prize symbol must occur exactly three times")
	}
	for symbol, count := range counts {
		if symbol != prize.Symbol && count >= 3 {
			return fmt.Errorf("weekly lottery board contains an extra triple")
		}
	}

	return nil
}
