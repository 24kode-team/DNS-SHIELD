FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /dns-shield ./cmd/shield

# ── Runtime image (minimal attack surface) ────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S shield && adduser -S shield -G shield

WORKDIR /app
COPY --from=builder /dns-shield /usr/local/bin/dns-shield
COPY configs/ /app/configs/

RUN mkdir -p /etc/dns-shield/tls /app/blocklists && \
    chown -R shield:shield /app /etc/dns-shield

USER shield

EXPOSE 53/udp 53/tcp 853/tcp 8080/tcp

ENTRYPOINT ["dns-shield"]
