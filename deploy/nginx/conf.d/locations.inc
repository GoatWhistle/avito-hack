include /etc/nginx/conf.d/security-headers.inc;

proxy_http_version 1.1;
proxy_set_header Host              $host;
proxy_set_header X-Real-IP         $remote_addr;
proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
proxy_set_header X-Forwarded-Host  $host;

location /api/v1/ws {
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

location = /metrics {
    deny all;
    return 404;
}

location /uploads/ {
    proxy_pass http://backend_upstream;
    proxy_read_timeout 60s;
}

location / {
    proxy_pass http://frontend_upstream;
    proxy_read_timeout 30s;

    proxy_hide_header X-Content-Type-Options;
    proxy_hide_header X-Frame-Options;
    proxy_hide_header Referrer-Policy;
    proxy_hide_header Permissions-Policy;
    proxy_hide_header Content-Security-Policy;
}
