FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./


RUN go mod tidy && go mod download

COPY . .

RUN go build -o main cmd/alum-bot/*

FROM alpine:3.21 AS app

WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ["/app/main"]