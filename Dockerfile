# syntax=docker/dockerfile:1

# El builder corre siempre en la plataforma nativa del runner (amd64),
# y Go cross-compila al target: evita emular toda la compilación con QEMU.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /app

# Args inyectados automáticamente por buildx (linux/arm64, linux/amd64)
ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o main ./cmd/alum-bot

FROM alpine:3.21 AS app

RUN apk add --no-cache tzdata
WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/messages .

ENTRYPOINT ["/app/main"]
