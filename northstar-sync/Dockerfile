# --- build stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /sync-service ./cmd/server

# --- final stage ---
FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /sync-service /usr/local/bin/sync-service

USER app
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/sync-service"]