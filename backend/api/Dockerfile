# Multi-stage build for Geevheit Intelligence 360 API
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/migrate ./cmd/migrate

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/provision ./cmd/provision

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/bin/migrate /app/migrate
COPY --from=builder /app/bin/provision /app/provision

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/api"]

