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

`.env.production.example` — это прод-шаблон: в нём каждая переменная снабжена пометкой,
что и почему нужно заменить. `.env.example` — для локальной разработки, на сервере он не нужен.

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
| `VITE_API_URL` | `https://ваш-домен` (без `/api`) | вшивается в бандл на этапе сборки |
| `LOG_FORMAT` | `json` | ставится автоматически в prod-оверрайде |

`REWARD_HMAC_SECRET` и `JWT_SECRET` должны быть не короче 16 символов — иначе бэкенд откажется
стартовать (проверка в `internal/config`).

Важно про `VITE_*`: это build-time переменные, они попадают в статический бандл. После их
изменения фронтенд нужно пересобрать (`--build`), рестарта контейнера недостаточно.

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

Ежедневный бэкап обоих хранилищ по крону:

```cron
30 3 * * * cd /home/deploy/avito-hack && docker compose -f docker-compose.yml -f docker-compose.prod.yml exec -T postgres pg_dump -U avito avito | gzip > /home/deploy/backups/db-$(date +\%F).sql.gz && docker run --rm -v avito-hack_uploads:/data -v /home/deploy/backups:/backup alpine tar czf /backup/uploads-$(date +\%F).tar.gz -C /data . && find /home/deploy/backups -mtime +7 -delete
```

Подключиться к БД (снаружи порт закрыт, только через контейнер):

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml exec postgres psql -U avito -d avito
```

---

## 9. Что ломается чаще всего

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
