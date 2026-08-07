package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type Deps struct {
	Register      *app.RegisterUserHandler
	Login         *app.LoginUserHandler
	GetProfile    *app.GetProfileHandler
	UpdateProfile *app.UpdateProfileHandler
	Validator     httpx.Validator
	Authenticate  func(http.Handler) http.Handler
	MaxBodyBytes  int64
}

type Handlers struct {
	deps    Deps
	decoder httpx.Decoder
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{deps: deps, decoder: httpx.NewDecoder(deps.Validator, deps.MaxBodyBytes)}
}

// @Id register
// @Summary Регистрация пользователя
// @Description Создаёт пользователя с ролью `user` и заводит для него питомца на стадии `egg`.
// @Description
// @Description Не идемпотентна: повторная регистрация того же e-mail даёт 409.
// @Description Ответ совпадает по форме с `POST /api/v1/auth/login` — сразу приходит
// @Description токен, отдельный вход после регистрации не нужен. Поля `expires_at` нет:
// @Description срок жизни лежит в claim `exp` внутри самого JWT.
// @Description
// @Description Пароль хранится как bcrypt-хеш. Ограничение длины пароля задаётся
// @Description доменом (`user/domain`), а не тегом `validate`, поэтому слишком короткий
// @Description пароль отклоняется с кодом `validation_error`, а не схемной ошибкой.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Учётные данные"
// @Success 201 {object} sessionResponse "Пользователь создан, сессия открыта"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope "E-mail уже зарегистрирован"
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security []
// @Router /api/v1/auth/register [post]
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if !h.decode(w, r, &req) {
		return
	}

	result, err := h.deps.Register.Handle(r.Context(), app.RegisterUserCommand{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.Created(w, sessionResponse{
		Token: result.Token,
		User:  toUserResponse(result.User),
	})
}

// @Id login
// @Summary Вход
// @Description Проверяет пару e-mail/пароль и выдаёт JWT.
// @Description
// @Description Идемпотентна по побочным эффектам, но каждый вызов выпускает новый токен
// @Description с новым `jti`. Срок действия — в клейме `exp`; отдельного поля
// @Description `expires_at` в ответе нет. Неверные учётные данные дают 401 без
// @Description уточнения, что именно неверно.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Учётные данные"
// @Success 200 {object} sessionResponse "Токен выдан"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope "Неверный e-mail или пароль"
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security []
// @Router /api/v1/auth/login [post]
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !h.decode(w, r, &req) {
		return
	}

	result, err := h.deps.Login.Handle(r.Context(), app.LoginUserCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, sessionResponse{
		Token: result.Token,
		User:  toUserResponse(result.User),
	})
}

// @Id getMe
// @Summary Профиль текущего пользователя
// @Description Возвращает профиль владельца токена. Идентификатор берётся из клейма `sub`.
// @Tags Users
// @Produce json
// @Success 200 {object} userResponse "Профиль"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/users/me [get]
func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	user, err := h.deps.GetProfile.Handle(r.Context(), app.GetProfileQuery{UserID: actor.ID})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toUserResponse(user))
}

// @Id updateMe
// @Summary Обновление профиля
// @Description Меняет отображаемое имя. Несмотря на метод PATCH, поле `full_name`
// @Description обязательно (`validate:"required"`) — частичное обновление без него
// @Description отклоняется. Идемпотентна: повтор с тем же значением даёт тот же результат.
// @Description
// @Description Пустая строка и строка из одних пробелов запрещены
// @Description CHECK-ограничением `users_full_name_not_blank`.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body updateProfileRequest true "Новое имя"
// @Success 200 {object} userResponse "Профиль обновлён"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/users/me [patch]
func (h *Handlers) UpdateMe(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	var req updateProfileRequest
	if !h.decode(w, r, &req) {
		return
	}

	user, err := h.deps.UpdateProfile.Handle(r.Context(), app.UpdateProfileCommand{
		UserID:   actor.ID,
		FullName: req.FullName,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toUserResponse(user))
}

func (h *Handlers) decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	return h.decoder.Decode(w, r, dst)
}
