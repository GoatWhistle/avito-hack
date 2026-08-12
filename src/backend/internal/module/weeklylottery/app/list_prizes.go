package app

import "github.com/avito-hack/backend/internal/module/weeklylottery/domain"

type ListPrizesHandler struct{}

func NewListPrizesHandler() *ListPrizesHandler {
	return &ListPrizesHandler{}
}

func (h *ListPrizesHandler) Handle() []PrizeCatalogView {
	prizes := domain.Prizes()
	out := make([]PrizeCatalogView, 0, len(prizes))
	for _, prize := range prizes {
		out = append(out, PrizeCatalogView{
			ID: prize.ID, Symbol: prize.Symbol, Title: prize.Title, Description: prize.Description,
			BenefitType: prize.BenefitType, BenefitValue: prize.BenefitValue,
			ScopeType: prize.ScopeType, ScopeValue: prize.ScopeValue,
		})
	}

	return out
}
