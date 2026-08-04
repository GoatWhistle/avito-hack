# Руководство по архитектуре: Как устроен проект (INTRO.md)

Если вы привыкли к **Clean Architecture** (Handler -> UseCase -> Domain -> Repository), этот проект построен на тех же принципах, но адаптирован под командный хакатон и требования высокой масштабируемости.

Ниже простое объяснение того, что происходит, почему файлов так много и как в этом не запутаться.

---

## 1. Карта соответствия (Clean Arch ↔ Наш проект)

В каждом модуле (например, `src/backend/internal/module/raccoon/`) вы увидите 4 папки:

| Папка в проекте | Ваш привычный аналог | За что отвечает |
| :--- | :--- | :--- |
| **`api/`** | **Handlers / Controllers** | Принимает HTTP-запросы, валидирует JSON (DTO), вызывает UseCase, возвращает HTTP 200/400/500. |
| **`app/`** | **Use Cases / Application** | Бизнес-сценарии приложения. Содержит интерфейсы (**Ports**) и оркерстрацию задач. |
| **`domain/`** | **Domain / Entities** | Чистая бизнес-модель и правила (расчет XP, валидация). Без импортов SQL, HTTP или JSON. |
| **`infra/`** | **Repositories / Adapters** | Работа с базой данных (Postgres, pgx) и внешними сервисами. Реализует интерфейсы из `app/`. |

---

## 2. Почему так много файлов в одной папке?

В небольших проектах часто делают один файл `repository.go` на 1000 строк и один `handler.go` на 1500 строк. В нашем проекте действуют 3 простых правила:

### 1) Один UseCase — один файл
Вместо одного гигантского файла `services.go`, каждый сценарий вынесен в отдельный файл в папке `app/`:
- `create_raccoon.go` — создание
- `update_raccoon.go` — редактирование
- `change_status.go` — смена статуса
- `get_profile.go` — получение профиля

*Плюс:* Удобно читать, легко писать unit-тесты и коммитить без конфликтов в Git.

### 2) Лимит 250 строк на файл
Правило проекта запрещает файлы длиннее 250 строк. Поэтому код разбит на компактные фокусные модули.

### 3) Разделение на Чтение (Queries) и Запись (Commands)
В папке `infra/`:
- `pg_repository.go` — отвечает за **Запись/Изменение** (сохранить агрегат с транзакцией).
- `pg_read_model.go` — отвечает за **Быстрое чтение** (SELECT с пагинацией и фильтрами для списка).

---

## 3. Как устроена работа в команде (Backend #2 vs Backend #3)

Вам не нужно ждать другого разработчика, чтобы писать свой код. Для этого используются **Ports (Интерфейсы)**:

```text
[HTTP Handler / Kafka Consumer] (Backend #3)
            │
            ▼
   [app/ports.go (Интерфейсы)]
            │
    ┌───────┴────────┐
    ▼                ▼
[Stub (Мок)]    [Real Core] (Backend #2)
```

1. **`app/ports.go`** — договор между вами и Backend #2. Здесь описано, *какие методы* должны быть у ядра.
2. **`app/stub_core.go`** — ваша «заглушка». Она отдаёт фейковые данные, чтобы вы могли прямо сейчас писать HTTP-хэндлеры и фронтенд-эндпоинты.
3. Когда Backend #2 допишет реальную бизнес-логику в `domain/` и `app/`, вы просто замените `stub` на реальный класс в `main.go`.

---

## 4. Чек-лист: Что писать вам (Backend #3)

Ваша работа сосредоточена вокруг входов (HTTP/Kafka) и выдачи наград:

1. **В `api/`**: Написать HTTP-контроллеры для `/api/v1/raccoon/profile`, `/api/v1/badges`, `/api/v1/rewards/claim`.
2. **В `app/`**: Написать сценарии связки `GetRaccoonProfileUseCase` и `ClaimRewardUseCase`.
3. **В `infra/`**: Написать адаптер генерации промокодов `PromocodeGeneratorAdapter`.

Вам **не нужно** менять математику уровней и таблицы в `domain/` — за это отвечает Backend #2!

---

## 5. Простой наглядный пример: Прослеживаем один HTTP-эндпоинт

Давайте разберем на пальцах, как работает эндпоинт **`GET /api/v1/raccoon/profile`** от HTTP до возврата данных:

### Шаг 1: DTO в `api/dto.go` (Формат JSON для фронтенда)
```go
type RaccoonProfileResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Level int    `json:"level"`
    XP    int    `json:"xp"`
}
```

### Шаг 2: Handler в `api/read_handlers.go` (HTTP Контроллер)
```go
func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
    // 1. Достаем ID пользователя из JWT-токена
    userID := httpx.MustActorID(r.Context())

    // 2. Вызываем UseCase
    profile, err := h.deps.GetProfileUseCase.Execute(r.Context(), userID)
    if err != nil {
        httpx.WriteError(w, err)
        return
    }

    // 3. Отдаем JSON 200 OK
    httpx.WriteJSON(w, http.StatusOK, toProfileResponse(profile))
}
```

### Шаг 3: UseCase в `app/get_profile.go` (Сценарий)
```go
type GetRaccoonProfileUseCase struct {
    reader RaccoonStateReader // Зависим от ИНТЕРФЕЙСА, а не базы!
}

func (uc *GetRaccoonProfileUseCase) Execute(ctx context.Context, userID uuid.UUID) (*RaccoonProfileView, error) {
    // Вызываем метод интерфейса
    return uc.reader.GetRaccoonProfile(ctx, userID)
}
```

### Шаг 4: Port & Stub в `app/ports.go` и `app/stub_core.go` (Заглушка)
```go
// В app/ports.go (Интерфейс)
type RaccoonStateReader interface {
    GetRaccoonProfile(ctx context.Context, userID uuid.UUID) (*RaccoonProfileView, error)
}

// В app/stub_core.go (Фейковые данные для теста)
func (s *StubRaccoonStateReader) GetRaccoonProfile(ctx context.Context, userID uuid.UUID) (*RaccoonProfileView, error) {
    return &RaccoonProfileView{
        ID:    userID,
        Name:  "Rocky",
        Level: 3,
        XP:    350,
    }, nil
}
```

### Итог:
Когда фронтенд делает `GET /api/v1/raccoon/profile`:
1. `Handler` берёт токен.
2. `UseCase` вызывает `RaccoonStateReader`.
3. `Stub` сразу возвращает Енота уровня 3 с 350 XP.
4. `Handler` возвращает `{ "name": "Rocky", "level": 3, "xp": 350 }`.

Всё работает прямо сейчас! А когда **Backend #2** напишет Postgres или Redis, вы просто передадите в `GetRaccoonProfileUseCase` вместо `Stub` класс `PgReadModel` от Backend #2!

