# Shared routes for both the HTTP and the HTTPS server block.
# Not a *.conf file on purpose: nginx.conf only includes conf.d/*.conf,
# so this file is pulled in explicitly and never becomes a server of its own.

add_header X-Content-Type-Options    "nosniff"        always;
add_header X-Frame-Options           "SAMEORIGIN"     always;
add_header Referrer-Policy           "strict-origin-when-cross-origin" always;
add_header Permissions-Policy        "camera=(), microphone=(), geolocation=()" always;

proxy_http_version 1.1;
proxy_set_header Host              $host;
proxy_set_header X-Real-IP         $remote_addr;
proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
proxy_set_header X-Forwarded-Host  $host;

# WebSocket must come before /api so the longer prefix wins.
# Without Upgrade/Connection the handshake degrades to a plain request,
# and without the long read timeout nginx kills an idle socket after 60s.
location /api/v1/ws {
    # ОБЯЗАТЕЛЬНО: клиент передаёт JWT query-параметром (браузерный WebSocket API
    # не умеет кастомные заголовки). Строка запроса попадает в access_log вместе
    # с токеном, поэтому логирование этой локации выключено полностью.
    access_log off;

    proxy_pass http://backend_upstream;

    proxy_set_header Upgrade    $http_upgrade;
    proxy_set_header Connection $connection_upgrade;

    proxy_read_timeout  3600s;
    proxy_send_timeout  3600s;
    proxy_connect_timeout 10s;

    proxy_buffering off;
    proxy_cache     off;
}

location ^~ /api/v1/auth/ {
    limit_req zone=auth_limit burst=10 nodelay;
    proxy_pass http://backend_upstream;
    proxy_read_timeout 60s;
}

location /api/ {
    limit_req zone=api_limit burst=40 nodelay;
    proxy_pass http://backend_upstream;
    proxy_read_timeout 60s;

    # Multipart photo upload, backend enforces MAX_PHOTO_BYTES itself.
    client_max_body_size 8m;
}

location = /healthz {
    proxy_pass http://backend_upstream;
    access_log off;
}

location = /readyz {
    proxy_pass http://backend_upstream;
    access_log off;
}

# Metrics stay internal: exposing them publicly leaks route and timing data.
location = /metrics {
    deny all;
    return 404;
}

# Item photos MUST go through the backend, not straight off the volume.
# internal/server/router.go wraps /uploads/ in uploadGuard: it resolves the item
# id from the path and returns 404 unless the item is published/sold and not
# soft-deleted. Serving the directory with `alias` here would bypass that check
# and expose photos of drafts and deleted listings to anyone guessing a URL.
# Caching is left to the backend's own headers.
location /uploads/ {
    proxy_pass http://backend_upstream;
    proxy_read_timeout 60s;
}

location / {
    proxy_pass http://frontend_upstream;
    proxy_read_timeout 30s;
}
