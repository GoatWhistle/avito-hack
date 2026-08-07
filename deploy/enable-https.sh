#!/usr/bin/env bash
#
# Включает HTTPS после того, как certbot выпустил сертификат.
#
#   ./deploy/enable-https.sh ваш-домен
#
# Что делает:
#   1. проверяет, что сертификат для домена действительно существует в томе;
#   2. генерирует conf.d/ssl.conf из шаблона с подстановкой домена;
#   3. включает редирект с HTTP на HTTPS (кроме acme-challenge);
#   4. проверяет конфиг через nginx -t и перечитывает его без даунтайма;
#   5. при любой ошибке откатывает изменения и оставляет рабочий HTTP.
#
# Скрипт идемпотентен: повторный запуск с тем же доменом ничего не ломает.

set -euo pipefail

DOMAIN="${1:-}"
if [[ -z "$DOMAIN" ]]; then
    echo "usage: $0 <domain>" >&2
    exit 1
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

COMPOSE="docker compose -f docker-compose.yml -f docker-compose.prod.yml"
NGINX_DIR="deploy/nginx"
SSL_CONF="$NGINX_DIR/conf.d/ssl.conf"
APP_CONF="$NGINX_DIR/conf.d/app.conf"

echo "==> Проверяю сертификат для $DOMAIN"
if ! $COMPOSE exec -T nginx test -f "/etc/letsencrypt/live/$DOMAIN/fullchain.pem"; then
    echo "ОШИБКА: сертификат /etc/letsencrypt/live/$DOMAIN/fullchain.pem не найден." >&2
    echo "Сначала выпустите его:" >&2
    echo "  DOMAIN=$DOMAIN LETSENCRYPT_EMAIL=you@example.com \\" >&2
    echo "    $COMPOSE --profile certbot run --rm certbot" >&2
    exit 1
fi

BACKUP="$(mktemp -d)"
cp "$APP_CONF" "$BACKUP/app.conf"
[[ -f "$SSL_CONF" ]] && cp "$SSL_CONF" "$BACKUP/ssl.conf"

rollback() {
    echo "==> Откат: возвращаю HTTP-конфигурацию"
    cp "$BACKUP/app.conf" "$APP_CONF"
    if [[ -f "$BACKUP/ssl.conf" ]]; then
        cp "$BACKUP/ssl.conf" "$SSL_CONF"
    else
        rm -f "$SSL_CONF"
    fi
    $COMPOSE exec -T nginx nginx -s reload 2>/dev/null || $COMPOSE restart nginx
    echo "Откат завершён, сайт работает по HTTP."
}
trap 'rollback; rm -rf "$BACKUP"' ERR

echo "==> Генерирую $SSL_CONF"
sed "s/__DOMAIN__/$DOMAIN/g" "$NGINX_DIR/ssl.conf.template" > "$SSL_CONF"

echo "==> Включаю редирект HTTP -> HTTPS"
# Заменяем include общих локаций в 80-м блоке на редирект.
# acme-challenge объявлен выше по файлу с префиксом ^~ и продолжает работать:
# без него автопродление сертификата сломается.
if grep -q 'return 301 https://\$host\$request_uri;' "$APP_CONF"; then
    echo "редирект уже был включён"
elif grep -q '^    include /etc/nginx/conf.d/locations.inc;' "$APP_CONF"; then
    sed -i 's|^    include /etc/nginx/conf.d/locations.inc;|    return 301 https://$host$request_uri;|' "$APP_CONF"
    echo "редирект включён"
else
    echo "ОШИБКА: не нашёл, что заменить в $APP_CONF — проверьте файл вручную" >&2
    exit 1
fi

echo "==> Проверяю конфигурацию nginx"
$COMPOSE exec -T nginx nginx -t

echo "==> Перечитываю конфигурацию"
$COMPOSE exec -T nginx nginx -s reload

trap - ERR
rm -rf "$BACKUP"

echo
echo "HTTPS включён для $DOMAIN."
echo
echo "ОСТАЛОСЬ СДЕЛАТЬ ВРУЧНУЮ — иначе фронтенд не заработает:"
echo "  1. в .env: VITE_API_URL=https://$DOMAIN и ALLOWED_ORIGINS=https://$DOMAIN"
echo "  2. пересобрать фронтенд (VITE_* вшиваются в бандл на этапе сборки):"
echo "     $COMPOSE up -d --build frontend backend"
echo
echo "Проверка:"
echo "  curl -sI https://$DOMAIN | head -1                # 200"
echo "  curl -sI https://$DOMAIN/api/v1/ws | head -1      # 401/426, не 502"
