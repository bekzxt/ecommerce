FROM golang:1.24

# Сначала копируем весь монорепозиторий
WORKDIR /app
COPY .. .

# Переходим в директорию api-gateway
WORKDIR /app/api-gateway

# Устанавливаем зависимости
RUN go mod tidy

# Собираем бинарник
RUN go build -o app ./cmd/main.go

EXPOSE 8080

CMD ["./app"]
