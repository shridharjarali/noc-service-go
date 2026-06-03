# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /noc-services ./cmd/noc-services

# ── Stage 2: Final ────────────────────────────────────────────────────────────
FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /noc-services /noc-services
COPY --from=builder /app/configs /configs

EXPOSE 8099

ENTRYPOINT ["/noc-services"]
