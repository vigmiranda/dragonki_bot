FROM golang:1.21

WORKDIR /app
COPY . .

WORKDIR /app/cmd/bot
RUN go build -o /app/dragonki

CMD ["/app/dragonki"]


