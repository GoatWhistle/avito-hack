# API «Авито Тамагочи»

Подробная спецификация HTTP- и WebSocket-API с примерами `curl`. Обзорная таблица эндпоинтов —
в [README](../README.md), продуктовая логика — в [CASE.md](CASE.md).

- Базовый URL локально: `http://localhost:8080`
- Префикс всех бизнес-маршрутов: `/api/v1`
- Формат: JSON, `Content-Type: application/json; charset=utf-8` (исключение — загрузка фото,
  там `multipart/form-data`)
- Авторизация: `Authorization: Bearer <token>`
- Время: RFC 3339 в UTC
- Деньги: целое число **копеек**, никаких чисел с плавающей точкой

Дальше в примерах используется переменная `$T` с access-токеном:

```bash
export T=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"anna@demo.avito","password":"demo1234"}' | jq -r .token)
```

---

## Единый конверт ошибок

Все эндпоинты отдают ошибки одинаково (`internal/shared/apierr`). Тела успешных ответов
различаются, тело ошибки — никогда.

```json
{
  "error": {
    "code": "validation_error",
    "message": "field is required",
    "field": "email",
    "request_id": "b9f2c1a4-3d5e-4c81-9f0a-7e1d2c3b4a5f"
  }
}
```

| Поле | Всегда | Описание |
| --- | --- | --- |
| `code` | да | Машинный код, стабильный контракт. Фронтенд переводит его в текст под локаль |
| `message` | да | Техническое описание на английском. **Не предназначено для показа пользователю** |
| `field` | нет | Присутствует только у `validation_error` — имя поля, которое не прошло проверку |
| `request_id` | да | Идентификатор запроса, тот же, что в логах сервера |

### Коды ошибок

| `code` | HTTP | Когда возникает | Что делать клиенту |
| --- | --- | --- | --- |
| `validation_error` | 400 | Поле не прошло валидацию или нарушен доменный инвариант | Показать ошибку у поля из `field` |
| `bad_request` | 400 | Тело не разобралось: битый JSON, превышен `MAX_BODY_BYTES`, неизвестное поле | Исправить запрос |
| `unauthorized` | 401 | Токена нет, он истёк или подделан | Отправить на повторный вход |
| `forbidden` | 403 | Роль не позволяет действие либо пользователь не владелец ресурса | Скрыть действие в UI |
| `not_found` | 404 | Ресурса нет или он удалён | Показать пустое состояние |
| `conflict` | 409 | Дубль по уникальному ключу, недопустимый переход статуса, повторная активация награды | Обновить состояние и показать причину |
| `internal_error` | 500 | Непредвиденная ошибка | Показать общее сообщение; детали ищутся в логах по `request_id` |

Единственное место классификации — `apierr.classify`. Доменные ошибки (`domainerr.InvalidError`,
`ConflictError`, `ErrNotFound`, `ErrForbidden`, `ErrUnauthorized`, `ErrConflict`) разворачиваются
через `errors.Is/As`, поэтому обёртка `fmt.Errorf("...: %w", err)` не превращает 400 в 500.

---

## Пагинация

Списочные эндпоинты используют **keyset-курсор**, а не `OFFSET`. Причина: `OFFSET` на глубоких
страницах заставляет БД пролистывать всё до нужного места, а при вставках между запросами
страницы «прыгают» — элементы дублируются или пропадают.

Курсор — `base64url` от пары `(created_at, id)`, где `id` разрывает ничьи по одинаковому времени.

Параметры запроса:

| Параметр | По умолчанию | Ограничение |
| --- | --- | --- |
| `limit` | 20 | Максимум 100, значения выше молча обрезаются |
| `cursor` | пусто | Значение `next_cursor` из предыдущего ответа |

Формат ответа:

```json
{
  "items": [ { "...": "..." } ],
  "next_cursor": "MjAyNi0wOC0wNVQxMDoxMjozMFoxfDAxOTQuLi4"
}
```

`items` всегда массив — при пустой выборке это `[]`, а не `null`.
Пустой или отсутствующий `next_cursor` означает, что данные закончились.

```bash
curl -s "http://localhost:8080/api/v1/items?limit=2" -H "Authorization: Bearer $T"
curl -s "http://localhost:8080/api/v1/items?limit=2&cursor=MjAyNi0wOC0wNVQx..." -H "Authorization: Bearer $T"
```

Битый курсор — `400 validation_error` с `"field": "cursor"`.

---

## Аутентификация

Access-токен: JWT, HS256, время жизни `JWT_TTL` (по умолчанию 15 минут). Срок хранится в клейме
`exp` внутри токена, отдельного поля в ответе нет. Refresh-токенов в MVP нет — по истечении
требуется повторный вход (см. «[Ограничения MVP](../README.md#ограничения-mvp)» в README).

### POST /api/v1/auth/register

Регистрация. Публичный эндпоинт. Побочный эффект: публикуется событие `user.registered`, по
которому pet-модуль заводит пользователю яйцо.

```bash
curl -i -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "new@example.com",
    "password": "supersecret",
    "full_name": "Новый Пользователь"
  }'
```

| Поле | Тип | Правила |
| --- | --- | --- |
| `email` | string | Обязательное, формат почты, уникальное |
| `password` | string | Обязательное, от 8 до 72 символов (верхний предел — ограничение bcrypt) |
| `full_name` | string | Обязательное, непустое |

Нарушение длины пароля даёт `400` с кодом `validation_error`, полем `password` и текстом,
указывающим границу. Требований к составу символов намеренно нет — регистрация должна
оставаться быстрой.

`201 Created` — тот же конверт, что у входа, чтобы после регистрации не делать второй запрос:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "018f2c3d-4e5f-7a8b-9c0d-1e2f3a4b5c6d",
    "email": "new@example.com",
    "full_name": "Новый Пользователь",
    "role": "user",
    "created_at": "2026-08-05T09:41:12Z"
  }
}
```

Пароль не возвращается и не логируется никогда. Хранится как bcrypt cost 12.

Ошибки: `400 validation_error` (невалидное поле), `409 conflict` (почта занята).

### POST /api/v1/auth/login

```bash
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"anna@demo.avito","password":"demo1234"}'
```

`200 OK`:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "d0000000-0000-4000-a000-000000000001",
    "email": "anna@demo.avito",
    "full_name": "Анна Ковалёва",
    "role": "user",
    "created_at": "2026-06-06T09:41:12Z"
  }
}
```

Ошибки: `401 unauthorized` — одинаковый ответ и на несуществующую почту, и на неверный пароль,
чтобы нельзя было перебором узнать, какие адреса зарегистрированы.

---

## Пользователи

### GET /api/v1/users/me

```bash
curl -s http://localhost:8080/api/v1/users/me -H "Authorization: Bearer $T"
```

`200 OK` — объект `user` того же вида, что в ответе логина.
Ошибки: `401 unauthorized`.

### PATCH /api/v1/users/me

```bash
curl -s -X PATCH http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $T" \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"Анна К."}'
```

`200 OK` — обновлённый объект пользователя.
Ошибки: `400 validation_error`, `401 unauthorized`.

---

## Объявления

Статус-машина (`item/domain/status.go`):

```text
draft ──submit──▶ moderation ──publish──▶ published ──sell──▶ sold
  │                                            │
  └──────────────publish────────────────────────┘
                                             archive
                                                │
                                                ▼
                                            archived ──restore──▶ draft
```

Переход `publish` доступен владельцу напрямую из `draft` — иначе основной сценарий кейса
заблокирован. Путь через `moderation` сохранён как опциональный. Переход `sell` терминальный.
Недопустимый переход — `409 conflict`.

Публикация и продажа порождают события (`item.published`, `item.sold`), на которые реагирует
питомец. Событие публикуется **после коммита** транзакции.

### GET /api/v1/items

Список. Авторизация опциональна: анонимно видны только опубликованные.

| Параметр | Описание |
| --- | --- |
| `status` | Фильтр: `draft`, `moderation`, `published`, `archived` |
| `owner_id` | UUID владельца |
| `search` | Поиск по заголовку (триграммный индекс) |
| `limit`, `cursor` | Пагинация |

```bash
curl -s "http://localhost:8080/api/v1/items?status=published&limit=2"
```

`200 OK`:

```json
{
  "items": [
    {
      "id": "a0000000-0000-4000-a000-000000000001",
      "owner_id": "d0000000-0000-4000-a000-000000000001",
      "owner_name": "Анна Ковалёва",
      "title": "iPhone 13 128GB, идеальное состояние",
      "price": 4990000,
      "status": "published",
      "created_at": "2026-07-18T09:41:12Z"
    }
  ],
  "next_cursor": "MjAyNi0wNy0xOFQwOTo0MToxMlow..."
}
```

`price` — копейки: `4990000` это 49 900 ₽.

Список намеренно легче карточки: без описания и атрибутов. Ошибки: `400 validation_error`
(неизвестный `status` или битый `cursor`).

### GET /api/v1/items/{id}

```bash
curl -s http://localhost:8080/api/v1/items/a0000000-0000-4000-a000-000000000001
```

`200 OK`:

```json
{
  "id": "a0000000-0000-4000-a000-000000000001",
  "owner_id": "d0000000-0000-4000-a000-000000000001",
  "title": "iPhone 13 128GB, идеальное состояние",
  "description": "Продаю iPhone 13 на 128 ГБ в отличном состоянии...",
  "price": 4990000,
  "status": "published",
  "attributes": { "category": "electronics", "condition": "used" },
  "created_at": "2026-07-18T09:41:12Z",
  "updated_at": "2026-07-18T09:41:12Z"
}
```

Ошибки: `404 not_found` (нет объявления либо оно неопубликовано и вы не владелец).

### GET /api/v1/items/mine

Мои объявления в любом статусе. Требует авторизации, поддерживает `status`, `limit`, `cursor`.

```bash
curl -s "http://localhost:8080/api/v1/items/mine?limit=5" -H "Authorization: Bearer $T"
```

### POST /api/v1/items

Создаёт объявление в статусе `draft`. XP за это не начисляется — питомец реагирует на публикацию,
а не на черновик.

```bash
curl -s -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $T" \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Механическая клавиатура Keychron K2",
    "description": "Беспроводная механическая клавиатура Keychron K2 на коричневых свитчах. Пользовался год дома, все клавиши работают, подсветка белая, регулируется по яркости. Кейкапы не залипают, следов от жидкости нет. В комплекте кабель USB-C, съёмник кейкапов и запасные свитчи. Подключается по Bluetooth сразу к трём устройствам и по проводу. Продаю, потому что перешёл на клавиатуру меньшего формата.",
    "price": 890000,
    "attributes": { "category": "electronics", "condition": "used" }
  }'
```

| Поле | Правила |
| --- | --- |
| `title` | Обязательное, 3–200 символов |
| `description` | До 5000 символов. Для XP за «качественное объявление» нужно **строго больше 200 символов** |
| `price` | Копейки, ≥ 0 |
| `attributes` | Словарь строк, до 50 ключей |

`201 Created` — объект объявления. Ошибки: `400 validation_error`, `401 unauthorized`.

### PATCH /api/v1/items/{id}

Частичное обновление, только владельцем. Передаются лишь изменяемые поля.

```bash
curl -s -X PATCH http://localhost:8080/api/v1/items/a0000000-0000-4000-a000-000000000001 \
  -H "Authorization: Bearer $T" \
  -H 'Content-Type: application/json' \
  -d '{"price": 4750000}'
```

`200 OK`. Ошибки: `400`, `401`, `403 forbidden` (не владелец), `404 not_found`.

### POST /api/v1/items/{id}/status

Смена статуса. **Это ключевой эндпоинт демо**: `publish` запускает цепочку событий, доходящую до
питомца и до WebSocket.

```bash
curl -s -X POST http://localhost:8080/api/v1/items/a0000000-0000-4000-a000-000000000001/status \
  -H "Authorization: Bearer $T" \
  -H 'Content-Type: application/json' \
  -d '{"action":"publish"}'
```

| `action` | Переход | Побочный эффект |
| --- | --- | --- |
| `submit` | `draft → moderation` | — |
| `publish` | `draft`/`moderation → published` | Событие `item.published`, XP питомцу, WS-push |
| `sell` | `published → sold` | Событие `item.sold`, крупное начисление XP, WS-push |
| `archive` | `published → archived` | — |
| `restore` | `archived → draft` | — |

`200 OK` — объявление с новым статусом. Сразу после этого в открытом WebSocket прилетят
`pet.updated` и `xp.gained`.

Ошибки: `400 validation_error` (неизвестный `action`), `403 forbidden` (не владелец),
`409 conflict` (недопустимый переход из текущего статуса).

Повторная публикация того же объявления **не начислит XP второй раз** — мешает
`UNIQUE (user_id, action, subject_id)` на `xp_events`. Само объявление при этом опубликуется
нормально: провал бизнес-правила питомца не отменяет действие пользователя.

### Фотографии

Фото важны не только визуально: наличие фото — одно из условий «качественного объявления»,
за которое начисляется повышенный XP, и главный триггер советов питомца.

**POST /api/v1/items/{id}/photos** — `multipart/form-data`, поле `photo`.
Допустимые типы: `image/jpeg`, `image/png`, `image/webp` (определяются по содержимому файла, а не
по расширению и не по заголовку клиента). Лимит — `MAX_PHOTO_BYTES`, по умолчанию 5 МБ.
До 10 фото на объявление (позиции 0–9).

```bash
curl -s -X POST http://localhost:8080/api/v1/items/a0000000-0000-4000-a000-000000000001/photos \
  -H "Authorization: Bearer $T" \
  -F "photo=@/path/to/photo.jpg"
```

`201 Created`:

```json
{
  "id": "0192a3b4-c5d6-7e8f-9a0b-1c2d3e4f5a6b",
  "item_id": "a0000000-0000-4000-a000-000000000001",
  "url": "/uploads/0192a3b4-c5d6-7e8f-9a0b-1c2d3e4f5a6b.jpg",
  "position": 0,
  "created_at": "2026-08-05T09:41:12Z"
}
```

`url` относительный — файл раздаётся статикой (nginx на проде), а не через Go.

**GET /api/v1/items/{id}/photos** — массив фото, отсортированный по `position`.

**DELETE /api/v1/items/{id}/photos/{photoID}** — `204 No Content`, только владельцем.

Ошибки загрузки: `400 validation_error` с `"field": "photo"` — файл отсутствует, слишком большой
или неподдерживаемого типа.

---

## Избранное

Добавление в избранное даёт питомцу +1 XP с лимитом 5 в день. Лимит считается по журналу
`xp_events` на сервере, а не приходит от клиента.

### POST /api/v1/items/{id}/favorite

```bash
curl -i -X POST http://localhost:8080/api/v1/items/a0000000-0000-4000-a000-000000000005/favorite \
  -H "Authorization: Bearer $T"
```

`204 No Content`. Операция идемпотентна на уровне схемы (`PRIMARY KEY (user_id, item_id)`):
повторное добавление не создаёт дубль и не начисляет XP снова.

Ошибки: `401 unauthorized`, `404 not_found`.

### DELETE /api/v1/items/{id}/favorite

`204 No Content`. Удаление из избранного **не отнимает уже начисленный XP** — журнал append-only,
списаний в экономике нет вообще. Это осознанное решение: иначе добавление и удаление по кругу
превращалось бы в способ манипулировать балансом.

### GET /api/v1/favorites

```bash
curl -s "http://localhost:8080/api/v1/favorites?limit=20" -H "Authorization: Bearer $T"
```

`200 OK`:

```json
{
  "items": [
    {
      "item_id": "a0000000-0000-4000-a000-000000000005",
      "owner_id": "d0000000-0000-4000-a000-000000000005",
      "title": "Ноутбук ASUS Vivobook 15, Ryzen 5",
      "price": 4210000,
      "status": "published",
      "photo_url": "/uploads/....jpg",
      "added_at": "2026-08-01T09:41:12Z"
    }
  ],
  "next_cursor": ""
}
```

`photo_url` — первое фото объявления, может отсутствовать.

---

## Питомец, бейджи и награды

Наружу сущность называется «енот»/`raccoon`, внутри кода это модуль `pet`. Публичные пути
сохранены ради совместимости с фронтендом.

### GET /api/v1/raccoon/profile

Состояние питомца с **ленивым пересчётом параметров**: значения деградации применяются в момент
чтения, поэтому питомец, которого не навещали, вернётся голодным без всяких фоновых задач.

```bash
curl -s http://localhost:8080/api/v1/raccoon/profile -H "Authorization: Bearer $T"
```

`200 OK`:

```json
{
  "id": "e0000000-0000-4000-a000-000000000001",
  "user_id": "d0000000-0000-4000-a000-000000000001",
  "name": "Ноти",
  "level": 15,
  "xp": 581,
  "xp_to_next_level": 0,
  "current_streak": 41,
  "stage": "legend",
  "state": "happy",
  "badges": [
    {
      "id": "raccoon_friend",
      "name": "Друг Енота",
      "description": "Достичь 10 уровня",
      "icon_url": "/assets/badges/friend.svg",
      "earned_at": "2026-07-06T09:41:12Z"
    }
  ]
}
```

| Поле | Описание |
| --- | --- |
| `level` | 1–15, выводится из суммарного XP по порогам |
| `xp` | Суммарный накопленный XP |
| `xp_to_next_level` | Порог следующего уровня, не остаток до него; `0` означает достигнутый максимум. Остаток отдаёт `GET /progress` в поле с тем же именем |
| `current_streak` | Дней подряд с заходом, «день» считается по МСК |
| `stage` | Стадия роста: `egg`, `baby`, `teen`, `adult`, `legend` |
| `state` | Текущее настроение: `happy`, `neutral`, `sad`, `sleeping` |
| `badges` | Полученные пользователем бейджи |

Пороги уровней 1–15: 0, 5, 12, 22, 35, 52, 72, 95, 122, 155, 195, 240, 290, 350, 420.

Ошибки: `401 unauthorized`.

### GET /api/v1/badges

Бейджи, **уже полученные** текущим пользователем.

```bash
curl -s http://localhost:8080/api/v1/badges -H "Authorization: Bearer $T"
```

`200 OK` — массив объектов бейджа того же вида, что в профиле; у каждого заполнен `earned_at`.
У нового пользователя ответ — пустой массив `[]`, а не весь каталог со статусами.

Сам каталог из 11 бейджей сидится миграцией `00004`: «Исследователь», «Быстрый палец», «Охотник за
скидками», «Режиссёр Авито», «Честный критик», «Надёжный бро», «Ничего не скроешь», «Молния»,
«Зелёная планета», «Заряженный стрик», «Друг Енота».

### POST /api/v1/rewards/claim

Забрать награду, доступную по уровню, и получить персональный промокод.

```bash
curl -s -X POST http://localhost:8080/api/v1/rewards/claim \
  -H "Authorization: Bearer $T" \
  -H 'Content-Type: application/json' \
  -d '{"reward_id":"free_delivery_500"}'
```

`200 OK`:

```json
{
  "reward_id": "free_delivery_500",
  "promocode": "MZXW6YTBOI2A-K4RTQZ7NB3VXY2LM"
}
```

Структура кода: `base32(nonce)` + `-` + усечённая `HMAC-SHA256(user_id ‖ reward_id ‖ nonce)`,
подписанная серверным секретом `REWARD_HMAC_SECRET`. Подробности защиты — в разделе
«[Безопасность наград](../README.md#безопасность-наград-подпись-промокодов-hmac)» README.

Что это значит на практике:

- Код бесполезен без аккаунта, которому он выдан: подпись включает `user_id`, и активация
  проверяет её для пользователя из токена.
- Подобрать или сгенерировать код без серверного секрета нельзя.
- Повторная выдача той же награды невозможна: `UNIQUE (user_id, reward_id)`.
- Повторная активация невозможна: переход выполняется как `UPDATE ... WHERE status = 'granted'`,
  и из двух параллельных запросов строку изменит ровно один.

Ошибки: `400 validation_error` (пустой `reward_id`), `401 unauthorized`,
`403 forbidden` (условие награды не выполнено), `409 conflict` (награда уже выдана).

### Каталог наград

Сидится миграцией `00004`, по одной награде на уровень.

| Уровень | `reward_id` | Название | Тип |
| --- | --- | --- | --- |
| 1 | `starter_status` | Статус «Новичок» | cosmetic |
| 2 | `raccoon_accessory` | Кепка и очки для Енота | cosmetic |
| 3 | `services_discount_10` | Скидка 10% на «Выделение» | promo |
| 4 | `attentive_badge` | Бейдж «Внимательный» | cosmetic |
| 5 | `listing_badge_or_delivery_discount` | Значок на объявление или −20% на доставку | promo |
| 6 | `raccoon_avatar_frame` | Цветная рамка аватарки | cosmetic |
| 7 | `favorites_personal_offer` | Персональный промокод по избранному | promo |
| 8 | `dialogue_master_badge` | Бейдж «Мастер диалога» | cosmetic |
| 9 | `autoteka_discount_30` | Скидка 30% на Автотеку | promo |
| 10 | `free_delivery_500` | Бесплатная доставка до 500 ₽ | promo |
| 11 | `xl_listing_discount_50` | Скидка 50% на XL-объявление | promo |
| 12 | `treasure_hunter_badge` | Бейдж «Охотник за сокровищами» | cosmetic |
| 13 | `views_boost_24h` | Увеличение просмотров на 24 часа | promo |
| 14 | `raccoon_sticker_pack` | Стикер-пачка с Енотом | cosmetic |
| 15 | `free_delivery_and_avito_guru_badge` | Купон и бейдж «Гуру Авито» | promo |

Статусы выданной награды: `granted` → `activated` → `expired`.

### Начисления XP

Все значения — из `pet/domain/economy.go`. Клиент не участвует в начислении вообще: XP появляется
только как реакция бэкенда на доменное событие.

| Действие | XP | Лимит | Условие |
| --- | --- | --- | --- |
| Публикация объявления | 50 | 3/день | Уникальное объявление, `pet/domain/listing.go` |
| Продажа объявления | 100 | без дневного лимита | Уникальное объявление |
| Ежедневный заход | 1 | 1/день | ×1.5 при серии от 7 дней |
| Подписка на поиск | 3 | 3/неделю | — |
| Добавление в избранное | 1 | 5/день | Уникальное объявление |
| Новый диалог | 2 | — | Уникальный собеседник |
| Быстрый ответ | 3 | — | В течение часа и в дневное время (8:00–23:00 МСК) |
| Качественное объявление | 2 | 1/объявление | Фото + цена + описание **строго больше 200 символов** |
| То же с видео | 4 | 1/объявление | Дополнительно видео |
| Текстовый отзыв | 1 | — | Подтверждённая сделка |
| Отзыв с фото | 2 | — | Подтверждённая сделка + вложение |
| Отзыв с видео | 4 | — | Подтверждённая сделка + видео длиннее 30 секунд |
| Бонус за серию | 30 / 70 / 150 / 300 | — | Серия 3 / 7 / 14 / 30 дней |

Итоговое начисление домножается на состояние питомца: сытость ниже 30 — ×0.5, настроение от 70 —
×1.25. Положительное начисление никогда не округляется до нуля — минимум 1 XP.

Действия, для которых в MVP ещё нет источника событий (диалоги, отзывы, подписки на поиск),
реализованы и покрыты тестами в домене, но пока не вызываются: соответствующих модулей нет.
Экономика работает на публикациях, продажах, избранном и ежедневных заходах.

---

## WebSocket

Реалтайм-состояние питомца. Реализация — `coder/websocket`, hub в памяти процесса
(`internal/shared/ws`).

### Подключение

```text
ws://localhost:8080/api/v1/ws?token=<access_token>
wss://ваш-домен/api/v1/ws?token=<access_token>
```

Токен передаётся **query-параметром**: браузерный `WebSocket` не умеет задавать заголовки при
рукопожатии. Валидируется тем же кодом, что и HTTP-авторизация.

Особенности реализации:

- Маршрут зарегистрирован **вне** `chimw.Timeout`: общий таймаут запроса (30 с) убивал бы
  долгоживущее соединение.
- Origin проверяется по `ALLOWED_ORIGINS`. Неуказанный домен — рукопожатие отклоняется.
- Лимит одного сообщения — 4096 байт.
- Один пользователь может держать несколько вкладок: hub хранит `map[userID][]conn` и рассылает
  событие во все его соединения.

Ошибки рукопожатия отдаются обычным HTTP-конвертом до апгрейда: `401 unauthorized` при
отсутствующем или невалидном токене.

Проверить руками:

```bash
# ожидается 401 — апгрейд доходит до бэкенда, но токена нет
curl -sI "http://localhost:8080/api/v1/ws" | head -1

# с токеном (нужен websocat или wscat)
websocat "ws://localhost:8080/api/v1/ws?token=$T"
```

### Формат сообщений

Один конверт в обе стороны:

```json
{ "type": "...", "request_id": "...", "payload": { } }
```

`request_id` опционален, задаётся клиентом и возвращается в ответе — по нему клиент сопоставляет
ответ со своим запросом. У push-событий, инициированных сервером, его нет.

### Клиент → сервер

| `type` | Ответ | Описание |
| --- | --- | --- |
| `ping` | `pong` | Heartbeat |
| `pet.get` | `pet.state` | Запросить состояние с применённой деградацией |
| `pet.pet` | `pet.state` | Погладить питомца: безлимитно, поднимает настроение |

```json
{ "type": "ping", "request_id": "1" }
{ "type": "pet.get", "request_id": "2" }
{ "type": "pet.pet", "request_id": "3" }
```

Неизвестный `type` возвращает `error` с кодом `unknown_message`; соединение остаётся живым.

Все остальные действия выполняются обычным REST. WebSocket сознательно оставлен каналом
уведомлений, а не вторым API — так проще тестировать, отлаживать и восстанавливаться после
разрыва.

### Сервер → клиент: ответы

```json
{ "type": "pong", "request_id": "1" }
```

```json
{
  "type": "pet.state",
  "request_id": "2",
  "payload": {
    "id": "e0000000-0000-4000-a000-000000000001",
    "user_id": "d0000000-0000-4000-a000-000000000001",
    "name": "Ноти",
    "stage": "legend",
    "state": "happy",
    "level": 15,
    "xp": 581,
    "next_level_xp": 0,
    "satiety": 92,
    "happiness": 88,
    "energy": 100,
    "streak_days": 41,
    "freezes": 3,
    "is_hatched": true,
    "hatched_at": "2026-06-10T09:41:12Z",
    "last_checkin_date": "2026-08-05T00:00:00Z",
    "last_decay_time": "2026-08-05T08:50:56Z",
    "updated_at": "2026-08-05T09:50:56Z"
  }
}
```

```json
{
  "type": "error",
  "request_id": "9",
  "payload": { "code": "unknown_message", "message": "unsupported message type" }
}
```

Коды ошибок WebSocket: `unknown_message` (неподдерживаемый тип), `internal_error` (не удалось
обновить питомца). Они не пересекаются с HTTP-кодами `apierr` — это отдельный, более узкий набор.

Поля `payload` питомца:

| Поле | Тип | Описание |
| --- | --- | --- |
| `stage` | string | `egg`, `baby`, `teen`, `adult`, `legend` |
| `state` | string | `happy`, `neutral`, `sad`, `sleeping` — выводится из параметров |
| `level` | int | 1–15 |
| `xp` | int | Суммарный XP |
| `next_level_xp` | int | Порог следующего уровня, `0` на максимуме |
| `satiety` | int | Сытость 0–100, ниже 30 — множитель XP ×0.5 |
| `happiness` | int | Настроение 0–100, от 70 — множитель ×1.25 |
| `energy` | int | Энергия 0–100 |
| `streak_days` | int | Серия заходов подряд |
| `freezes` | int | Запас заморозок серии |
| `is_hatched` | bool | Вылупился ли питомец |
| `hatched_at` | timestamp | Момент вылупления, отсутствует у яйца |
| `last_checkin_date` | date | Последний засчитанный чек-ин, отсутствует до первого |
| `last_decay_time` | timestamp | Точка отсчёта ленивого пересчёта параметров |

### Сервер → клиент: push-события

Приходят без запроса, порождаются действиями пользователя на доске объявлений.

| `type` | `payload` | Когда |
| --- | --- | --- |
| `pet.updated` | Полное состояние питомца (как в `pet.state`) | Любое изменение параметров |
| `xp.gained` | `{"amount": 2, "reason": "item.published", "total": 37}` | Начислен XP |
| `level.up` | `{"level": 4}` | Достигнут новый уровень |
| `reward.granted` | `{"reward_id": "free_delivery_500", "title": "Бесплатная доставка"}` | Выдана награда |
| `streak.updated` | `{"days": 7, "milestone": true}` | Изменилась серия; `milestone` — достигнут рубеж 3/7/14/30 |

Значения `reason` у `xp.gained` совпадают с типами доменных событий: `item.published`,
`item.sold`, `favorite.added`, `user.registered`.

Порядок при одном действии: сначала `pet.updated`, затем `xp.gained` (если XP начислен), затем
`level.up` (если уровень вырос).

### Поведение клиента

- **Реконнект** с экспоненциальной задержкой при разрыве.
- **После реконнекта — рефетч через REST.** Источник правды — REST, WebSocket только оптимизация.
  Пропущенный за время разрыва push не приводит к рассинхрону.
- **Heartbeat `ping`** раз в 30 секунд: удерживает соединение через прокси и позволяет заметить
  мёртвый сокет раньше TCP-таймаута.
- **Токен живёт 15 минут.** После истечения соединение нужно переустановить с новым токеном.

Пример на TypeScript:

```typescript
const socket = new WebSocket(`${wsUrl}/api/v1/ws?token=${token}`);

socket.onmessage = (event) => {
  const message = JSON.parse(event.data) as ServerMessage;

  switch (message.type) {
    case 'pet.updated':
    case 'pet.state':
      applyPetState(message.payload);
      break;
    case 'xp.gained':
      showXpToast(message.payload.amount);
      break;
    case 'level.up':
      playLevelUpAnimation(message.payload.level);
      break;
    case 'reward.granted':
      invalidateRewards();
      break;
  }
};
```

---

## Служебные эндпоинты

Вне префикса `/api/v1`.

| Метод | Путь | Ответ | Назначение |
| --- | --- | --- | --- |
| GET | `/healthz` | `{"status":"ok"}` | Живость процесса. Не трогает БД — отвечает, пока процесс жив |
| GET | `/readyz` | `{"status":"ready"}` либо `503 {"status":"unavailable"}` | Готовность принимать трафик, пингует БД с таймаутом 2 с |
| GET | `/metrics` | Формат Prometheus | RED-метрики по маршрутам. На проде закрыт от внешнего доступа |

Разделение `healthz` и `readyz` намеренное: при недоступной БД процесс жив (перезапускать
бессмысленно), но трафик принимать не готов — балансировщик должен убрать его из ротации,
а не рестартовать.

---

## Питомец, прогресс и чек-ин по REST

Всё, что ниже, требует заголовка `Authorization: Bearer $T`.

### GET /api/v1/pet

Состояние питомца с параметрами, пересчитанными на момент запроса: домен применяет деградацию
от `last_decay_time` лениво, при чтении.

```bash
curl -s http://localhost:8080/api/v1/pet -H "Authorization: Bearer $T"
```

`200 OK` — объект с полями `id`, `user_id`, `name`, `stage`, `state`, `level`, `xp`,
`next_level_xp`, `satiety`, `happiness`, `energy`, `streak_days`, `freezes`, `is_hatched`,
`hatched_at`, `last_checkin_date`, `last_decay_time`, `updated_at`.

### POST /api/v1/pet/actions/stroke

Погладить питомца — поднимает настроение. То же действие доступно и через WebSocket (`pet.pet`).

### POST /api/v1/checkin

Ежедневный чек-ин: продлевает серию, начисляет опыт и повышает сытость. Повторный вызов в тот же
день по московскому времени возвращает ошибку дублирования.

`200 OK` — объект с `pet`, `xp_granted`, `level`, `previous_level`, `next_level_xp`,
`unlocked_rewards` и вложенным `streak` (`days`, `continued`, `freeze_used`, `reset`,
`milestone_bonus`, `milestone_reached`, `freezes_left`).

### GET /api/v1/progress

Короткая сводка прогресса: `level`, `xp`, `next_level_xp`, `xp_to_next_level`, `is_max_level`,
`stage`, `streak_days`, `freezes`.

---

## Награды по уровням

### GET /api/v1/rewards

Каталог наград со статусом для текущего пользователя. Ответ — конверт `{"items": [...]}`;
у элемента поля `id`, `title`, `description`, `kind`, `condition_type`, `condition_value`,
`unlocked`, `claimed`, `progress_current`, `progress_target`, а у уже выданных — ещё и `status`.

`progress_current` и `progress_target` дают прогресс-бар к условию: например «уровень 4 из 5».

### GET /api/v1/rewards/my

Выданные пользователю награды, тоже конвертом `{"items": [...]}`: `reward_id`, `title`,
`description`, `kind`, `status`, `granted_at`. Поля `code`, `activated_at` и `expires_at`
появляются после активации.

### POST /api/v1/rewards/{id}/activate

Активация выданной награды. Возвращает подписанный промокод. Повторная активация невозможна:
переход выполняется как `UPDATE ... WHERE status = 'granted'`, и из двух параллельных запросов
строку изменит ровно один.

---

## Ежедневная сводка

### GET /api/v1/summary/today

Сводка за прошедший день от лица питомца. Генерируется лениво при первом запросе и сохраняется в
`daily_summaries` — повторные обращения читают уже готовую запись.

```bash
curl -s http://localhost:8080/api/v1/summary/today -H "Authorization: Bearer $T"
```

`200 OK`:

```json
{
  "id": "9e9b59d1-6072-46b8-84bf-38289aa15ea7",
  "date": "2026-08-06",
  "message": "Ты набрал 2 XP, я подрос до 15 уровня. Серия держится уже 41 дней подряд.",
  "advice": {
    "text": "У объявления «iPhone 13 128GB» нет ни одной фотографии — с фото такие продаются заметно быстрее.",
    "item_id": "a0000000-0000-4000-a000-000000000001",
    "action": "add_photo"
  },
  "generated_by": "template",
  "facts": {
    "total_xp": 2,
    "actions": [{ "action": "daily_login", "count": 1, "amount": 2 }],
    "level": 15,
    "previous_level": 2,
    "leveled_up": true,
    "streak_days": 41,
    "leaderboard_rank": 1,
    "issues_count": 3
  },
  "created_at": "2026-08-07T15:51:58Z"
}
```

| Поле | Описание |
| --- | --- |
| `message` | Текст от первого лица, собранный поверх `facts` |
| `advice` | Ровно один совет с `item_id` и кодом действия — фронтенд делает из него ссылку на объявление. Отсутствует, если придраться не к чему |
| `facts` | Детерминированные факты дня: агрегаты XP, уровень, параметры, серия, позиция в лидерборде |
| `generated_by` | `template` либо `llm` — чем собран текст. На прогрессию не влияет |

Ошибки: `401 unauthorized`.

### GET /api/v1/summary/history

История сводок с keyset-пагинацией (`limit`, `cursor`), новые сверху.

```bash
curl -s "http://localhost:8080/api/v1/summary/history?limit=5" -H "Authorization: Bearer $T"
```

`200 OK` — конверт `{"items": [...], "next_cursor": "..."}`, элементы того же вида,
что в `/summary/today`.

---

## Лидерборд

### GET /api/v1/leaderboard

Рейтинг по уровню и опыту с keyset-курсором.

```bash
curl -s "http://localhost:8080/api/v1/leaderboard?around=me" -H "Authorization: Bearer $T"
```

| Параметр | Назначение |
| --- | --- |
| `cursor` | Курсор следующей страницы из поля `next_cursor` |
| `around` | Значение `me` — вернуть окрестность позиции текущего пользователя |

Поле `name` — имя пользователя, а не кличка питомца: email в рейтинге не раскрывается.

`200 OK`:

```json
{
  "items": [
    {
      "user_id": "d0000000-0000-4000-a000-000000000001",
      "name": "Анна Ковалёва",
      "level": 15,
      "xp": 581,
      "streak_days": 41,
      "rank": 1
    }
  ],
  "my_rank": 7,
  "next_cursor": "..."
}
```

---

## Что ещё не реализовано

Помечено честно, чтобы никто не искал этого в коде:

| Что | Статус |
| --- | --- |
| Задания дня `GET /tasks/today` | Не реализовано, отвечает `404` |
| Действия `POST /pet/actions/{play,feed}` | Не реализованы, доступно только «погладить» |
| Генерация текста сводки через LLM | Клиент к Anthropic API и промпт написаны (`internal/shared/llm`), но в сборку приложения не подключены: текст собирается шаблонным генератором, в ответе `generated_by: "template"` |
| Диалоги, отзывы, подписки на поиск как источники событий | Правила экономики написаны и покрыты тестами в домене, но вызывать их пока некому: соответствующих модулей нет |

Экономика работает на публикациях, продажах, избранном и ежедневных заходах.
