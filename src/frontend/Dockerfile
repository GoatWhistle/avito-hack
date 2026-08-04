# syntax=docker/dockerfile:1.7

FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

FROM node:22-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG VITE_API_URL=http://localhost:8080
ARG VITE_DEFAULT_LOCALE=ru
ENV VITE_API_URL=${VITE_API_URL} VITE_DEFAULT_LOCALE=${VITE_DEFAULT_LOCALE}

RUN npm run build

FROM node:22-alpine AS dev
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
EXPOSE 8080
ENTRYPOINT ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "8080"]

FROM nginx:1.27-alpine AS runtime

RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/app.conf
RUN esc=$(printf '\033') && sed -i "s/\\\\x1b/${esc}/g" /etc/nginx/conf.d/app.conf
COPY --from=builder /app/dist /usr/share/nginx/html

RUN addgroup -S app && adduser -S app -G app \
 && chown -R app:app /usr/share/nginx/html /var/cache/nginx /var/log/nginx \
 && touch /var/run/nginx.pid && chown app:app /var/run/nginx.pid

USER app
EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
  CMD wget -qO- http://127.0.0.1:8080/health || exit 1
