FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
RUN go build -o dragonki ./cmd/bot

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 🔥 COPIA O BINÁRIO
COPY --from=builder /app/dragonki .

# 🔥 COPIA OS ASSETS
COPY --from=builder /app/assets ./assets

CMD ["./dragonki"]
