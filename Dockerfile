FROM golang:1.25.3 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o sso ./cmd/app

FROM debian:bookworm-slim
WORKDIR /app

# Установим migrate для запуска миграций
RUN apt-get update && apt-get install -y curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz \
    | tar -xz && mv migrate /usr/local/bin/migrate && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/sso .
COPY migrations ./migrations

EXPOSE 50051
CMD ["./sso"]
