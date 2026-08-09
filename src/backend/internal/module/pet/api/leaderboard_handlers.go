package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/httpx"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

type leaderboardEntryResponse struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
	XP         int    `json:"xp"`
	StreakDays int    `json:"streak_days"`
	Rank       int    `json:"rank"`
}

type leaderboardResponse struct {
	Items      []leaderboardEntryResponse `json:"items"`
	MyRank     *int                       `json:"my_rank"`
	NextCursor string                     `json:"next_cursor,omitempty"`
}

type LeaderboardDeps struct {
	Leaderboard  *app.LeaderboardHandler
	Authenticate func(http.Handler) http.Handler
}

type LeaderboardHandlers struct {
	deps LeaderboardDeps
}

func NewLeaderboardHandlers(deps LeaderboardDeps) *LeaderboardHandlers {
	return &LeaderboardHandlers{deps: deps}
}

// @Id getLeaderboard
// @Summary Рейтинг пользователей
// @Description Рейтинг по уровню, затем по опыту, затем по идентификатору (для
// @Description устойчивости порядка). Поле `my_rank` — позиция владельца токена;
// @Description `null`, если она неизвестна или не запрашивалась. Подсчёт ранга —
// @Description самая дорогая часть запроса, поэтому он выполняется только при
// @Description `with_my_rank=true` или `around=me`.
// @Description
// @Description Пагинация использует ОТДЕЛЬНЫЙ формат курсора (`level|xp|user_id`),
// @Description несовместимый с временным курсором остальных списков. Передача
// @Description временного курсора сюда даёт 400 с `field: cursor`.
// @Description
// @Description Параметр `around=me` переключает выдачу на окно вокруг позиции
// @Description текущего пользователя вместо начала таблицы. Любое иное значение
// @Description параметра игнорируется — сравнение строгое, только строка `me`.
// @Tags Leaderboard
// @Produce json
// @Param limit query int false "Размер страницы. По умолчанию 20, максимум 100." default(20) maximum(100)
// @Param cursor query string false "Лидербордный курсор из `next_cursor` предыдущей страницы."
// @Param around query string false "Значение `me` — показать окно вокруг позиции текущего пользователя." Enums(me)
// @Param with_my_rank query bool false "Значение `true` — посчитать и вернуть `my_rank`." default(false)
// @Success 200 {object} leaderboardResponse "Страница рейтинга"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/leaderboard [get]
func (h *LeaderboardHandlers) List(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	query, err := leaderboardQueryFromRequest(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	query.UserID = actor.ID

	result, err := h.deps.Leaderboard.Handle(r.Context(), query)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toLeaderboardResponse(result))
}

func leaderboardQueryFromRequest(r *http.Request) (app.LeaderboardQuery, error) {
	limit, err := httpx.IntQuery(r, "limit", pagination.DefaultLimit)
	if err != nil {
		return app.LeaderboardQuery{}, err
	}

	cursor, err := pagination.DecodeLeaderboardCursor(r.URL.Query().Get("cursor"))
	if err != nil {
		return app.LeaderboardQuery{}, domainerr.NewInvalid("cursor", "invalid cursor")
	}

	return app.LeaderboardQuery{
		Cursor:     cursor,
		Limit:      pagination.NormalizeLimit(limit),
		Around:     r.URL.Query().Get("around") == "me",
		WithMyRank: r.URL.Query().Get("with_my_rank") == "true",
	}, nil
}

func toLeaderboardResponse(result app.LeaderboardResult) leaderboardResponse {
	items := make([]leaderboardEntryResponse, 0, len(result.Items))

	for _, entry := range result.Items {
		items = append(items, leaderboardEntryResponse{
			UserID:     entry.UserID.String(),
			Name:       entry.Name,
			Level:      entry.Level,
			XP:         entry.XP,
			StreakDays: entry.StreakDays,
			Rank:       entry.Rank,
		})
	}

	return leaderboardResponse{Items: items, MyRank: result.MyRank, NextCursor: result.NextCursor}
}
