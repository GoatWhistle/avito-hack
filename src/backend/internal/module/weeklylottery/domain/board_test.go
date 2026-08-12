package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
)

func TestValidateBoardRejectsExtraOrMismatchedTriples(t *testing.T) {
	t.Parallel()

	valid := [domain.BoardSize]domain.Symbol{
		domain.SymbolBicycle, domain.SymbolBicycle, domain.SymbolBicycle,
		domain.SymbolSofa, domain.SymbolSofa,
		domain.SymbolSneakers, domain.SymbolSneakers,
		domain.SymbolDelivery, domain.SymbolPromotion,
	}
	require.NoError(t, domain.ValidateBoard(valid, "weekly_bicycle_5"))
	require.Error(t, domain.ValidateBoard(valid, ""))
	require.Error(t, domain.ValidateBoard(valid, "weekly_smartphone_3"))

	extraTriple := valid
	extraTriple[7] = domain.SymbolSofa
	require.Error(t, domain.ValidateBoard(extraTriple, "weekly_bicycle_5"))
}

func TestValidateLosingBoardRejectsAnyTriple(t *testing.T) {
	t.Parallel()

	losing := [domain.BoardSize]domain.Symbol{
		domain.SymbolBicycle, domain.SymbolBicycle,
		domain.SymbolSmartphone, domain.SymbolSmartphone,
		domain.SymbolSofa, domain.SymbolSofa,
		domain.SymbolSneakers, domain.SymbolSneakers,
		domain.SymbolDelivery,
	}
	require.NoError(t, domain.ValidateBoard(losing, ""))

	losing[8] = domain.SymbolBicycle
	require.Error(t, domain.ValidateBoard(losing, ""))
}
