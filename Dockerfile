FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./


RUN go mod tidy && go mod download

COPY . .

RUN go build -o main cmd/alum-bot/*

FROM alpine:3.21 AS app

RUN apk add --no-cache tzdata
WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/messages .

ENTRYPOINT ["/app/main"]