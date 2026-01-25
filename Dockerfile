FROM golang:1.21 AS builder

WORKDIR /app
COPY . .
RUN go build -o dragonki ./cmd/bot

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/dragonki .
CMD ["./dragonki"]
