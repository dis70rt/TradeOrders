FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/api ./cmd/api
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/database ./cmd/database
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/matching ./cmd/match
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/stream ./cmd/stream

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ .