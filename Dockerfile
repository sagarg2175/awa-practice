FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o task-api .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/task-api .
COPY --from=builder /app/config.yaml .

EXPOSE 8080

CMD ["./task-api"]