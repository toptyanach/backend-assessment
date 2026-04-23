# Шаг сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник из нашего main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server/main.go

# Финальный легковесный образ
FROM alpine:latest

WORKDIR /app
COPY --from=builder /bin/server /app/server

EXPOSE 8080

CMD ["/app/server"]