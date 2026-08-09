# Деплой на VPS

Инструкция для чистого сервера: от `ssh root@...` до работающего домена с HTTPS и WebSocket.

Требования к серверу: 2 vCPU, 4 ГБ RAM, 20 ГБ диска, Ubuntu 22.04/24.04. Меньше тоже поедет
(1 vCPU / 2 ГБ), но сборка фронтенда на месте будет долгой.

---

## 1. Подготовка сервера

```bash
ssh root@SERVER_IP

apt update && apt upgrade -y
apt install -y ca-certificates curl git ufw

install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] \
https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" \
  > /etc/apt/sources.list.d/docker.list

apt update
apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

docker --version && docker compose version
```

Отдельный пользователь вместо работы из-под root:

```bash
adduser --disabled-password --gecos "" deploy
usermod -aG docker deploy
mkdir -p /home/deploy/.ssh
cp /root/.ssh/authorized_keys /home/deploy/.ssh/
chown -R deploy:deploy /home/deploy/.ssh
chmod 700 /home/deploy/.ssh && chmod 600 /home/deploy/.ssh/authorized_keys
```

Файрвол — наружу торчат только SSH и HTTP/HTTPS. Порты 8080, 3000, 5432 и 6379 закрыты, и
`docker-compose.prod.yml` их и не публикует:

```bash
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable
ufw status
```

---

## 2. Клонирование и конфигурация

```bash
su - deploy
git clone https://github.com/GoatWhistle/avito-hack.git
cd avito-hack
cp .env.production.example .env
```

`.env.production.example` — это прод-шаблон. `.env.example` — для локальной разработки, на
сервере он не нужен. Сам `.env` в `.gitignore` и коммититься не должен.

Из шаблона обязательно заполняются `POSTGRES_PASSWORD`, `DATABASE_URL`, `JWT_SECRET`,
`REWARD_HMAC_SECRET`, `ALLOWED_ORIGINS`, `DOMAIN` и `LETSENCRYPT_EMAIL` — с оставленными
плейсхолдерами бэкенд либо не стартует, либо стартует небезопасно. Разбор всех переменных — в
разделе [11](#11-переменные-окружения-подробно).

Сгенерировать секреты (никогда не оставляйте значения из примера):

```bash
echo "POSTGRES_PASSWORD=$(openssl rand -base64 24 | tr -d '/+=')"
echo "JWT_SECRET=$(openssl rand -hex 32)"
echo "REWARD_HMAC_SECRET=$(openssl rand -hex 32)"
```

Заполнить `.env`. Значения, которые обязательно меняются относительно локального запуска:

| Переменная | Значение на проде | Почему |
| --- | --- | --- |
| `POSTGRES_PASSWORD` | сгенерированный пароль | БД внутри сети, но пароль по умолчанию недопустим |
| `DATABASE_URL` | `postgres://avito:ПАРОЛЬ@postgres:5432/avito?sslmode=disable` | пароль должен совпадать с `POSTGRES_PASSWORD` |
| `JWT_SECRET` | 64 hex-символа | подпись access-токенов |
| `REWARD_HMAC_SECRET` | 64 hex-символа | подпись промокодов; смена секрета инвалидирует все выданные коды |
| `ALLOWED_ORIGINS` | `https://ваш-домен` | CORS и проверка Origin при WebSocket-рукопожатии |
| `VITE_API_URL` | оставить **пустым** | фронт ходит относительными путями, nginx проксирует |
| `LOG_FORMAT` | `json` | ставится автоматически в prod-оверрайде |

`REWARD_HMAC_SECRET` и `JWT_SECRET` должны быть не короче 16 символов — иначе `Load()` в
`internal/config` вернёт ошибку и контейнер `backend` уйдёт в рестарт-луп. Рекомендуется 64
hex-символа. Секрет наград задаётся **один раз до демо**: его смена на работающем проде навсегда
инвалидирует все ранее выданные промокоды, подписи перестанут сходиться.

Важно про `VITE_*`: это build-time переменные, они вшиваются в статический бандл на этапе
`npm run build`. Правка `.env` без пересборки фронтенда не изменит вообще ничего:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build frontend
```

---

## 3. Первый запуск

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps
```

Порядок запуска задан через `depends_on` с условиями: postgres → migrate (одноразовый, применяет
миграции и завершается) → backend → frontend → nginx.

Проверка:

```bash
curl -s http://localhost/healthz
curl -s http://localhost/readyz
curl -si http://localhost/ | head -1
```

Сайт уже доступен по `http://SERVER_IP`. Демо-данные (12 пользователей с питомцами разных стадий,
объявления, награды) приезжают миграцией `00005_seed_demo_data.sql` автоматически — отдельного шага
для сида нет, лидерборд у жюри заполнен сразу.

Общий пароль всех демо-аккаунтов — `demo1234`, логины вида `anna@demo.avito`, `boris@demo.avito`.
Быстрая проверка, что данные на месте:

```bash
curl -s -X POST http://localhost/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"anna@demo.avito","password":"demo1234"}'
```

---

## 4. Домен

В DNS-панели регистратора: `A`-запись `@` → IP сервера, и `A`-запись `www` → тот же IP.
Дождаться распространения:

```bash
dig +short ваш-домен
```

Пока запись не резолвится, выпускать сертификат бессмысленно — Let's Encrypt проверяет домен
HTTP-запросом на него же.

---

## 5. HTTPS

Certbot ходит по `http://ваш-домен/.well-known/acme-challenge/`, эта локация в nginx открыта
всегда и никуда не редиректится.

```bash
export DOMAIN=ваш-домен
export LETSENCRYPT_EMAIL=you@example.com

docker compose -f docker-compose.yml -f docker-compose.prod.yml \
  --profile certbot run --rm certbot
```

После успешного выпуска включить TLS одной командой — руками конфиги править не нужно:

```bash
./deploy/enable-https.sh ваш-домен
```

Скрипт проверяет наличие сертификата в томе, рендерит `deploy/nginx/ssl.conf.template` в
`conf.d/ssl.conf` с подстановкой домена, переключает 80-й порт на редирект (оставляя открытым
`acme-challenge` для автопродления), прогоняет `nginx -t` и перечитывает конфиг. Если что-то
не сходится — откатывает всё обратно на рабочий HTTP.

Проверка, включая WebSocket:

```bash
curl -sI https://ваш-домен | head -1
curl -sI https://ваш-домен/api/v1/ws | head -1
```

Второй запрос должен вернуть `401` (нет токена), а не `502` — значит, апгрейд доходит до бэкенда.

**После включения HTTPS обязательно** поправить `.env` (`VITE_API_URL=https://...`,
`ALLOWED_ORIGINS=https://...`) и пересобрать фронтенд, иначе браузер заблокирует запросы с
HTTPS-страницы на HTTP-API (mixed content), а `wss://` не пройдёт проверку Origin:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build frontend backend
```

Автопродление сертификата (раз в сутки, certbot сам решает, пора ли):

```bash
crontab -e
```

```cron
0 4 * * * cd /home/deploy/avito-hack && DOMAIN=ваш-домен LETSENCRYPT_EMAIL=you@example.com docker compose -f docker-compose.yml -f docker-compose.prod.yml --profile certbot run --rm certbot renew --quiet && docker compose -f docker-compose.yml -f docker-compose.prod.yml exec -T nginx nginx -s reload
```

---

## 6. Обновление

```bash
cd /home/deploy/avito-hack
git pull
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
docker image prune -f
```

Миграции применяются сами: контейнер `migrate` отрабатывает до старта бэкенда. Простой —
секунды на пересоздание контейнеров. Нулевого даунтайма в MVP нет и не планировалось.

Откат на предыдущий коммит:

```bash
git log --oneline -5
git checkout <sha>
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

Откат схемы БД, если новый коммит принёс несовместимую миграцию:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml run --rm migrate down
```

### Аварийный откат прямо на защите

Правило: **в день защиты не обновляем прод**. Если всё-таки обновились и сломалось — счёт идёт
на минуты, поэтому готовимся заранее.

Перед демо зафиксировать заведомо рабочий коммит и снять дамп:

```bash
git rev-parse HEAD > ~/GOOD_SHA
docker compose -f docker-compose.yml -f docker-compose.prod.yml exec -T postgres \
  pg_dump -U avito avito | gzip > ~/pre-demo.sql.gz
```

Откат за ~60 секунд (образы уже в кеше, пересборка быстрая):

```bash
cd /home/deploy/avito-hack
git checkout $(cat ~/GOOD_SHA)
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

Если БД пришла в негодность — восстановить дамп:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml stop backend
gunzip -c ~/pre-demo.sql.gz | docker compose -f docker-compose.yml -f docker-compose.prod.yml \
  exec -T postgres psql -U avito -d avito
docker compose -f docker-compose.yml -f docker-compose.prod.yml start backend
```

Самый быстрый вариант «вернуть как было», если сломан только код, а данные целы:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml restart backend frontend
```

---

## 7. Kafka (опционально)

Сервис `kafka` объявлен в `docker-compose.yml` под профилем `kafka` и по умолчанию **не
поднимается**. Без него события питомца идут внутрипроцессной шиной — весь пользовательский
сценарий работает полностью.

Включать Kafka для демо стоит только если её нужно показать жюри: брокер добавляет ~1 ГБ RAM и
20–40 секунд к холодному старту (KRaft-инициализация), то есть это самый медленный сервис в стеке.

```bash
COMPOSE_PROFILES=kafka docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

В `.env` при этом задать `KAFKA_BROKERS=kafka:9092`. Пустое значение = шина в процессе.
Порт брокера наружу не публикуется ни в одном режиме.

---

## 8. Эксплуатация

Логи (в проде — JSON, ротация 20 МБ × 5 файлов):

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f --tail=100 backend
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f --tail=100 nginx
```

Ресурсы и здоровье:

```bash
docker stats --no-stream
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps
```

Бэкап БД:

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml exec -T postgres \
  pg_dump -U avito avito | gzip > backup-$(date +%F).sql.gz
```

Восстановление:

```bash
gunzip -c backup-2026-08-05.sql.gz | \
  docker compose -f docker-compose.yml -f docker-compose.prod.yml exec -T postgres psql -U avito -d avito
```

Бэкап загруженных фото. Файлы лежат в docker-томе `avito-hack_uploads`, а не в репозитории:
дамп БД без них восстановит объявления с битыми картинками, поэтому бэкапить нужно оба.

```bash
docker run --rm -v avito-hack_uploads:/data -v "$PWD":/backup alpine \
  tar czf /backup/uploads-$(date +%F).tar.gz -C /data .
```

Восстановление фото:

```bash
docker run --rm -v avito-hack_uploads:/data -v "$PWD":/backup alpine \
  tar xzf /backup/uploads-2026-08-07.tar.gz -C /data
```

Демо-фото из образа попадают в том только при его первом создании: docker наполняет
том содержимым каталога образа, если том пуст, и не трогает непустой. После обновления
`src/backend/assets/demo-photos/` на стенде, где том уже существует, картинки останутся
старыми — том нужно наполнить явно:

```bash
docker run --rm -v avito-hack_uploads:/data \
  -v "$PWD/src/backend/assets/demo-photos":/src:ro alpine \
  sh -c 'cp -rf /src/. /data/ && chown -R 65532:65532 /data && chmod -R a+r /data'
```

Владелец 65532 — пользователь `nonroot`, под которым работает бэкенд: без `chown`
скопированные файлы останутся за root и отдадут 404 на публикации.

Ежедневный бэкап обоих хранилищ по крону:

```cron
30 3 * * * cd /home/deploy/avito-hack && docker compose -f docker-compose.yml -f docker-compose.prod.yml exec -T postgres pg_dump -U avito avito | gzip > /home/deploy/backups/db-$(date +\%F).sql.gz && docker run --rm -v avito-hack_uploads:/data -v /home/deploy/backups:/backup alpine tar czf /backup/uploads-$(date +\%F).tar.gz -C /data . && find /home/deploy/backups -mtime +7 -delete
```

Подключиться к БД (снаружи порт закрыт, только через контейнер):

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml exec postgres psql -U avito -d avito
```

---

## 9. Как устроены конфиги nginx

Конфигурация прокси намеренно разбита на несколько файлов, и почти каждое решение в ней
продиктовано граблями, на которые уже наступили. Ниже — что именно и почему, чтобы это не
пришлось выяснять заново.

### Заголовки безопасности вынесены в отдельный файл

`add_header` в nginx **не наследуется** в `location`, если внутри этого `location` объявлен
хотя бы один свой `add_header`. Наследование работает по принципу «всё или ничего»: одна
строка `add_header Cache-Control ...` внутри блока молча отменяет весь набор серверных
заголовков, объявленных уровнем выше.

Именно так проект однажды уже отдавал главную страницу без единого заголовка безопасности:
локации `/assets/` и `= /index.html` задавали свой `Cache-Control`, и вместе с ним теряли
X-Frame-Options, CSP и всё остальное. Дыра была не в списке заголовков, а в том, где он
объявлен.

Поэтому набор живёт в отдельном файле и подключается `include`-ом **в каждый `location`,
где есть собственный `add_header`**:

* прод — `deploy/nginx/conf.d/security-headers.inc`, подключается из `locations.inc` на
  уровне `server` и дополнительно внутрь нужных локаций;
* dev/standalone-фронтенд — `src/frontend/security-headers.inc`.

Дублировать строки руками нельзя: две копии рано или поздно разъезжаются, `include` — нет.

Правило для новых локаций: **добавил `add_header` в блок — сразу добавь туда же `include`
security-headers**. Иначе блок начнёт отдавать ответы без защиты, и это не вызовет ни ошибки
конфигурации, ни предупреждения при `nginx -t`.

Оба `.inc`-файла намеренно не имеют расширения `.conf`: `nginx.conf` инклюдит `conf.d/*.conf`,
и `.inc` не должен подключаться сам по себе как отдельный сервер.

### Почему в CSP есть `unsafe-inline` и `wasm-unsafe-eval`

Обе директивы обязательны, убирать их нельзя — приложение сломается.

`script-src 'unsafe-inline'`. React Router в SPA-режиме (`ssr: false`) генерирует `index.html`
на этапе сборки и вшивает в него инлайновые `<script>` с `window.__reactRouterContext` и
streamController. Подставить nonce негде: HTML статический, серверного рендера нет. Убрать
`'unsafe-inline'` — получить белый экран.

`script-src 'wasm-unsafe-eval'`. Плеер dotLottie компилирует `dotlottie-player.wasm` через
`WebAssembly.instantiateStreaming`. Без этой директивы браузер валит инициализацию, и питомец
молча падает в статичный SVG-фолбэк — то есть главная фича демо перестаёт быть живой.
Директива разрешает ровно компиляцию WebAssembly и **не** открывает `eval()`.

Внешние CDN в списках источников не нужны: wasm лежит в `public/lottie/` и отдаётся своим же
nginx (см. `setWasmUrl` в `app/components/pet-avatar`).

`connect-src 'self'` обходится без внешних origin, потому что `VITE_API_URL` пуст: фронт ходит
относительными путями, nginx рядом проксирует `/api/` и `/uploads/` на бэкенд. Запросы
same-origin, включая WebSocket `/api/v1/ws`. Если соберёте бандл с абсолютным `VITE_API_URL` —
этот origin придётся добавить в `connect-src` явно, иначе браузер заблокирует запросы.

`object-src 'none'` и `base-uri 'self'` закрывают устаревшие векторы инъекции. `frame-ancestors`
согласован с X-Frame-Options: `'none'` / `DENY` на dev-фронте, `'self'` / `SAMEORIGIN` на проде.

### Заголовки апстрима режутся на прокси

В `locations.inc` блок `location /` для фронтенда содержит `proxy_hide_header` на все пять
заголовков безопасности. Причина: фронтенд — это свой nginx, который сам выставляет полный
набор. За прокси они дублируются, и клиент получает по два X-Frame-Options и две CSP с
**разными** значениями (`DENY` от фронта и `SAMEORIGIN` отсюда). Браузер применит самую строгую
политику, но конфигурация становится непредсказуемой. Источник истины — прокси, он ближе к
клиенту, поэтому заголовки апстрима срезаются.

### `default_type`, а не `add_header Content-Type`

В health-локациях (`/nginx-health`, `/health`) тип ответа задаётся через `default_type`.
`add_header` именно **добавляет** строку, а не заменяет — в ответе оказывалось сразу два
заголовка `Content-Type`.

### Порядок локаций и лимиты запросов

`location ^~ /api/v1/auth/` объявлен до `location /api/`: более длинный префикс с `^~` должен
выигрывать, иначе логин попадёт в общую зону и получит лимит 20 r/s вместо 5 r/s. Отдельная
зона для авторизации нужна потому, что один логин стоит примерно 280 мс CPU (bcrypt cost=12) —
незалимиченный `/auth/` это одновременно и брутфорс, и дешёвый DoS с одного ноутбука.

`location /api/v1/ws` тоже идёт раньше `/api/`, и по той же причине: без `Upgrade`/`Connection`
рукопожатие деградирует до обычного запроса, а без длинного `proxy_read_timeout` nginx убивает
простаивающий сокет через 60 секунд.

Логирование WebSocket-локации выключено полностью и намеренно: браузерный WebSocket API не умеет
кастомные заголовки, поэтому клиент передаёт JWT query-параметром — строка запроса попала бы в
access_log вместе с токеном.

Те же лимиты продублированы в dev-конфиге `src/frontend/nginx.conf` (зоны `dev_api_limit` и
`dev_auth_limit`): там фронт проксирует на бэкенд напрямую, и без них `/api/` был бы открыт
вообще без ограничений.

### `/uploads/` обязан идти через бэкенд

Фотографии объявлений проксируются на бэкенд, а не раздаются напрямую с тома. `internal/server/router.go`
оборачивает `/uploads/` в `uploadGuard`: он достаёт id объявления из пути и возвращает 404, если
объявление не опубликовано, продано с ограничениями или удалено мягким удалением. Раздача
директории через `alias` обошла бы эту проверку и открыла фото черновиков и удалённых объявлений
любому, кто угадает URL. Кеширование оставлено на собственные заголовки бэкенда.

### `/metrics` закрыт снаружи

Публичные метрики утекают маршрутами и таймингами, поэтому локация отвечает 404.

### HTTP-only по умолчанию и включение HTTPS

`deploy/nginx/conf.d/app.conf` едет только с HTTP-блоком, чтобы
`docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d` поднимался на голом IP
вообще без сертификата: стек обязан встать раньше, чем certbot сможет проверить домен по порту 80.

HTTPS руками не настраивается. `deploy/nginx/ssl.conf.template` лежит уровнем выше `conf.d/` и
имеет расширение `.template` именно поэтому: nginx не стартует, если `ssl_certificate` указывает
на несуществующий файл, — пока сертификат не выпущен, конфига не должно существовать, иначе ляжет
весь стек, включая порт 80, по которому certbot и проходит проверку.

Включает HTTPS скрипт `./deploy/enable-https.sh ваш-домен`. Он идемпотентен и делает пять вещей:
проверяет, что сертификат для домена действительно есть в томе; рендерит `conf.d/ssl.conf` из
шаблона с подстановкой домена; заменяет `include` общих локаций в 80-м блоке на редирект на HTTPS;
валидирует результат через `nginx -t` и перечитывает конфиг без даунтайма; при любой ошибке
откатывает всё назад и оставляет рабочий HTTP.

`location ^~ /.well-known/acme-challenge/` остаётся открытым на 80-м порту навсегда и объявлен
выше по файлу, чем редирект: без него сломается автопродление сертификата.

HSTS выставлен с намеренно коротким `max-age=86400` (сутки) — это хакатон-демо, и если с
сертификатом что-то пойдёт не так, браузеры отпустят домен через день, а не через год. По той же
причине нет `includeSubDomains` и `preload`: откатить их нельзя.

### Прочие мелочи

`text/html` сжимается nginx безусловно и в список `gzip_types` не добавляется.

`MAX_PHOTO_BYTES` (5 МБ) должен оставаться меньше `client_max_body_size` в
`deploy/nginx/nginx.conf` (8 МБ), иначе nginx отрежет запрос раньше бэкенда и пользователь
получит 413 вместо внятной ошибки валидации.

---

## 10. Что ломается чаще всего

| Симптом | Причина | Что делать |
| --- | --- | --- |
| WebSocket не подключается, в консоли 502 или мгновенный close | nginx проксирует `/api/v1/ws` без `Upgrade`/`Connection` | проверить, что подключён `locations.inc`, там блок `location /api/v1/ws` идёт до `location /api/` |
| WebSocket рвётся каждую минуту | дефолтный `proxy_read_timeout 60s` | в `locations.inc` для `/api/v1/ws` стоит `3600s`, убедиться что конфиг перечитан |
| WebSocket закрывается сразу после handshake | домен не указан в `ALLOWED_ORIGINS`, бэкенд отбрасывает по Origin | добавить `https://домен` в `.env`, перезапустить backend |
| Mixed content в консоли браузера | `VITE_API_URL` остался `http://` | поправить `.env` и **пересобрать** фронтенд с `--build` |
| Бэкенд не стартует, в логах про secret | `JWT_SECRET` или `REWARD_HMAC_SECRET` короче 16 символов | сгенерировать нормальные секреты |
| `migrate` падает и перезапускается | конфликт версий или недоступная БД | `logs migrate`, проверить `DATABASE_URL` и что пароль совпадает с `POSTGRES_PASSWORD` |
| 502 на всё подряд | бэкенд не прошёл healthcheck | `logs backend`, проверить `readyz` и доступность postgres |
| Промокоды перестали активироваться | сменили `REWARD_HMAC_SECRET` | старые коды подписаны прежним секретом и невалидны навсегда — это ожидаемое поведение |
| Кончилось место на диске | логи и старые образы | `docker image prune -a -f`, `docker system df` |
| Фото объявлений отдают 404 | так и задумано для черновиков: `/uploads/` проксируется на бэкенд, а тот проверяет, что объявление опубликовано и не удалено | если 404 на опубликованном — `logs backend`, проверить том `uploads` и `UPLOAD_DIR` |
| Фронтенд собирается, но контейнер отдаёт 403/404 на всё | nginx раздаёт не ту директорию: React Router 8 кладёт сборку в `build/client` | проверить `COPY --from=builder /app/build/client` в `src/frontend/Dockerfile` — путь должен совпадать с реальным выходом сборки |
| Лидерборд пустой у жюри | сид-миграция не отработала | `logs migrate`; `exec postgres psql -U avito -d avito -c 'select count(*) from users;'` — на чистой базе должно быть 12 |
| `nginx` не стартует после включения HTTPS | сертификата нет по указанному пути | `./deploy/enable-https.sh` сам откатывается; вручную — удалить `deploy/nginx/conf.d/ssl.conf` и перезапустить nginx |
| Стек долго поднимается | профиль `kafka` включён, брокеру нужно 20–40 с | для демо Kafka не обязательна: снять `COMPOSE_PROFILES=kafka` и очистить `KAFKA_BROKERS` |
| Заголовки безопасности пропали на части страниц | в `location` добавили свой `add_header`, а `include` security-headers — нет | добавить `include` в этот блок, см. раздел 9 |
| Питомец показывается статичной картинкой вместо анимации | из CSP убрали `wasm-unsafe-eval`, плеер dotLottie не инициализируется | вернуть директиву в `script-src`, см. раздел 9 |
| Белый экран, в консоли ошибки CSP про inline script | из CSP убрали `'unsafe-inline'` | вернуть директиву: React Router вшивает инлайновые скрипты в статический `index.html` |

---

## 11. Переменные окружения подробно

Полные списки — в `.env.example` (локальная разработка) и `.env.production.example` (прод).
Ниже то, что не выводится из названия переменной.

### `VITE_API_URL` по умолчанию пустой — и это правильно

Пустое значение означает, что фронт запрашивает относительный `/api/v1` своего же origin, а
nginx рядом проксирует `/api/` и `/uploads/` на `backend:8080`. Что это даёт:

* домен меняется без пересборки бандла — стенд одинаково открывается по `localhost:3000`,
  `127.0.0.1:3000` и по доменному имени;
* mixed content невозможен в принципе: схема запроса всегда совпадает со схемой страницы;
* запросы same-origin, то есть CORS в них вообще не участвует.

Абсолютный URL нужен, только если API вынесен на отдельный origin. В этом случае придётся
добавить origin фронта в `ALLOWED_ORIGINS`, прописать его же в `connect-src` политики CSP и
следить за совпадением схем.

### `ALLOWED_ORIGINS`

Используется и для CORS, и для проверки Origin при WebSocket-рукопожатии, поэтому после
включения HTTPS здесь обязан быть `https://` — иначе сокет закроется сразу после handshake.
Формат: точная схема плюс домен, через запятую без пробелов.

Поскольку фронт ходит на API относительными путями, его собственные запросы под CORS не
попадают. Список нужен для прямых обращений к порту 8080: Swagger UI, curl из devtools,
сторонний фронтенд. В локальном `.env.example` перечислены и `localhost`, и `127.0.0.1` —
для браузера это разные origin, и стенд должен открываться одинаково по обоим.

### Пул подключений к PostgreSQL

`DB_MAX_CONNS` — верхняя граница одновременных соединений от **одного** инстанса бэкенда.
Суммарно по всем инстансам значение не должно превышать `max_connections` постгреса (по
умолчанию 100), иначе новые соединения начнут отвергаться.

`DB_MIN_CONNS` — тёплый резерв, который держится открытым. Он **не должен быть больше**
`DB_MAX_CONNS`: иначе `Load()` вернёт ошибку и контейнер уйдёт в рестарт-луп.

`DB_MAX_CONN_LIFETIME` и `DB_MAX_CONN_IDLE_TIME` заставляют пересоздавать соединения — это
страхует от утечки серверных ресурсов и подхватывает смену параметров постгреса.
`DB_HEALTH_CHECK_PERIOD` — период фоновой проверки живости соединений.

`DB_CONNECT_TIMEOUT` ограничивает не только установку соединения, но и стартовый ping: слишком
малое значение уронит бэкенд, если постгрес поднимается дольше.

Значения в `.env.example` совпадают с дефолтами `postgres.DefaultPoolOptions()` и перечислены
явно только для наглядности. То же касается таймаутов HTTP-сервера — их дублируют из
`internal/config`, чтобы было видно, что поднимать для медленной сети.

### Загрузка файлов

`MAX_PHOTO_BYTES` (5 МБ) должен быть меньше `client_max_body_size` в `deploy/nginx/nginx.conf`
(8 МБ). Иначе nginx отрежет запрос раньше бэкенда, и пользователь получит 413 вместо внятной
ошибки валидации.

### Kafka

Пустой `KAFKA_BROKERS` означает, что события идут внутрипроцессной шиной и брокер не нужен
вообще. Kafka-путь дополнительно требует сборки с тегом: `go build -tags kafka ./...`.
Для демо с брокером: `KAFKA_BROKERS=kafka:9092`, сервис `kafka` есть в `docker-compose.yml`
под профилем (см. раздел 7). Порт наружу не публикуется ни в одном режиме.

Сам брокер работает в режиме KRaft, без ZooKeeper. Лимит памяти контейнера ниже 1 ГБ ставить
нельзя: 768 МБ heap плюс оверхед JVM дают OOMKill прямо на старте.

### Логи

`LOG_FORMAT` и `LOG_COLOR` в `docker-compose.prod.yml` жёстко переопределены на `json`/`never`.
Значения из `.env` влияют только на запуск без prod-оверрайда.

### `ANTHROPIC_API_KEY`

Переменная необязательная: без ключа приложение работает штатно, просто не генерируется
ежедневная сводка от питомца. Ключ берётся в консоли Anthropic и хранится только в `.env` на
сервере — коммитить его нельзя.

### Демо-фото в образе бэкенда

Демо-фотографии объявлений попадают в образ, а не монтируются в том напрямую. Docker при первом
старте заполняет пустой named volume содержимым каталога из образа, поэтому после
`docker compose down -v && up -d` каталог объявлений сразу с картинками. Пути внутри —
`/uploads/{item_id}/{файл}`, ровно как их пишет `LocalPhotoStorage`, а строки `item_photos`
ссылаются на них из миграции `00011`.
