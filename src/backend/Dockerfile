FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
      -trimpath \
      -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
      -o /out/api ./cmd/api

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
      -o /out/migrate ./cmd/migrate

RUN mkdir -p /out/uploads

FROM golang:1.25-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata wget \
 && go install github.com/air-verse/air@v1.63.0

EXPOSE 8080

ENTRYPOINT ["air", "-c", ".air.toml"]

FROM gcr.io/distroless/static-debian12:nonroot AS runtime

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /out/api /app/api
COPY --from=builder /out/migrate /app/migrate
COPY --from=builder --chown=nonroot:nonroot /out/uploads /data/uploads

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/app/api"]
