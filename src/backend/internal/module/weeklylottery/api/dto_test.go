package api

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/weeklylottery/app"
	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
)

func TestRunResponseDoesNotSerializeClosedSymbols(t *testing.T) {
	t.Parallel()

	response := toRunResponse(app.RunView{
		ID: "abc", State: domain.StateActive, CreatedAt: time.Now(),
		Slots: []app.SlotView{
			{Index: 0, Opened: true, Symbol: domain.SymbolBicycle},
			{Index: 1, Opened: false, Symbol: ""},
		},
	})

	raw, err := json.Marshal(response)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"id":"abc",
		"state":"active",
		"slots":[
			{"index":0,"opened":true,"symbol":"bicycle"},
			{"index":1,"opened":false}
		],
		"created_at":"`+response.CreatedAt.Format(time.RFC3339Nano)+`"
	}`, string(raw))
}

func TestStateResponseMatchesFrontendContract(t *testing.T) {
	t.Parallel()

	weekStart := time.Date(2026, time.August, 10, 0, 0, 0, 0, time.FixedZone("MSK", 3*60*60))
	slots := make([]app.SlotView, domain.BoardSize)
	for index := range domain.BoardSize {
		slots[index] = app.SlotView{Index: index}
	}
	response := toStateResponse(app.StateView{
		Available: true, WeekStart: weekStart, NextAvailable: weekStart.AddDate(0, 0, 7),
		Run: &app.RunView{
			ID: "run123456789", State: domain.StateActive, CreatedAt: weekStart,
			Slots: slots,
		},
	})

	raw, err := json.Marshal(response)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"available": true,
		"week_start": "2026-08-10T00:00:00+03:00",
		"next_available_at": "2026-08-17T00:00:00+03:00",
		"run": {
			"id": "run123456789",
			"state": "active",
			"slots": [
				{"index": 0, "opened": false},
				{"index": 1, "opened": false},
				{"index": 2, "opened": false},
				{"index": 3, "opened": false},
				{"index": 4, "opened": false},
				{"index": 5, "opened": false},
				{"index": 6, "opened": false},
				{"index": 7, "opened": false},
				{"index": 8, "opened": false}
			],
			"created_at": "2026-08-10T00:00:00+03:00"
		}
	}`, string(raw))
}

func TestPrizeCatalogMatchesFrontendContract(t *testing.T) {
	t.Parallel()

	response := toPrizeCatalogResponse([]app.PrizeCatalogView{{
		ID: "weekly_bicycle_5", Symbol: domain.SymbolBicycle,
		Title: "Скидка 5% на спорт и отдых", Description: "Описание",
		BenefitType: "percent_discount", BenefitValue: 5,
		ScopeType: domain.ScopeTypeCategory, ScopeValue: domain.CategorySport,
	}})

	raw, err := json.Marshal(response)
	require.NoError(t, err)
	assert.JSONEq(t, `[{
		"id": "weekly_bicycle_5",
		"symbol": "bicycle",
		"title": "Скидка 5% на спорт и отдых",
		"description": "Описание",
		"benefit_type": "percent_discount",
		"benefit_value": 5,
		"scope_type": "category",
		"scope_value": "sport"
	}]`, string(raw))
}
