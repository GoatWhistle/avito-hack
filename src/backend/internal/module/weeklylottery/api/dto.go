package api

import (
	"time"

	"github.com/avito-hack/backend/internal/module/weeklylottery/app"
)

type SlotResponse struct {
	Index  int    `json:"index"`
	Opened bool   `json:"opened"`
	Symbol string `json:"symbol,omitempty"`
}

type PrizeResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	BenefitType  string    `json:"benefit_type"`
	BenefitValue int       `json:"benefit_value"`
	ScopeType    string    `json:"scope_type"`
	ScopeValue   string    `json:"scope_value"`
	Code         string    `json:"code"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type PrizeCatalogResponse struct {
	ID           string `json:"id"`
	Symbol       string `json:"symbol"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	BenefitType  string `json:"benefit_type"`
	BenefitValue int    `json:"benefit_value"`
	ScopeType    string `json:"scope_type"`
	ScopeValue   string `json:"scope_value"`
}

type RunResponse struct {
	ID        string         `json:"id"`
	State     string         `json:"state"`
	Slots     []SlotResponse `json:"slots"`
	Prize     *PrizeResponse `json:"prize,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

type StateResponse struct {
	Available     bool         `json:"available"`
	WeekStart     time.Time    `json:"week_start"`
	NextAvailable time.Time    `json:"next_available_at"`
	Run           *RunResponse `json:"run,omitempty"`
}

type RevealResponse struct {
	Index  int         `json:"index"`
	Symbol string      `json:"symbol"`
	Run    RunResponse `json:"run"`
}

func toRunResponse(view app.RunView) RunResponse {
	slots := make([]SlotResponse, 0, len(view.Slots))
	for _, slot := range view.Slots {
		slots = append(slots, SlotResponse{
			Index: slot.Index, Opened: slot.Opened, Symbol: string(slot.Symbol),
		})
	}

	response := RunResponse{ID: view.ID, State: string(view.State), Slots: slots, CreatedAt: view.CreatedAt}
	if view.Prize != nil {
		response.Prize = &PrizeResponse{
			ID: view.Prize.ID, Title: view.Prize.Title, Description: view.Prize.Description,
			BenefitType: view.Prize.BenefitType, BenefitValue: view.Prize.BenefitValue,
			ScopeType: view.Prize.ScopeType, ScopeValue: view.Prize.ScopeValue,
			Code: view.Prize.Code, ExpiresAt: view.Prize.ExpiresAt,
		}
	}

	return response
}

func toStateResponse(view app.StateView) StateResponse {
	response := StateResponse{
		Available: view.Available, WeekStart: view.WeekStart, NextAvailable: view.NextAvailable,
	}
	if view.Run != nil {
		run := toRunResponse(*view.Run)
		response.Run = &run
	}

	return response
}

func toPrizeCatalogResponse(views []app.PrizeCatalogView) []PrizeCatalogResponse {
	out := make([]PrizeCatalogResponse, 0, len(views))
	for _, view := range views {
		out = append(out, PrizeCatalogResponse{
			ID: view.ID, Symbol: string(view.Symbol), Title: view.Title, Description: view.Description,
			BenefitType: view.BenefitType, BenefitValue: view.BenefitValue,
			ScopeType: view.ScopeType, ScopeValue: view.ScopeValue,
		})
	}

	return out
}
