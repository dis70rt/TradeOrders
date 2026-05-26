FROM golang:1.22-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/api ./cmd/api
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/database ./cmd/database
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/matching ./cmd/match
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/stream ./cmd/stream

FROM alpine:3.20
RUN apk --no-cache add ca-certificates \
    && addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/ .

USER appuser