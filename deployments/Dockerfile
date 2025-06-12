# Стадия сборки
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN go build -o backend ./cmd/backend/main.go

# Финальный минимальный образ
FROM debian:bookworm-slim

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/backend .

# Копируем нужные дополнительные файлы
COPY --from=builder /app/migrations ./migrations

# Пример: копируем configs, static, templates и т.д. — если нужно
# COPY --from=builder /app/configs ./configs
# COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./backend"]