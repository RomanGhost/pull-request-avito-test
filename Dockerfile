# Базовый образ
FROM golang:1.24 as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY go.mod go.sum ./
#ENV GOPROXY=direct
RUN go mod download

# Копируем весь код
COPY . .

# Сборка приложения
RUN go build -o app ./

# Финальный образ
FROM debian:bookworm-slim

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированное приложение из предыдущего контейнера
COPY --from=builder /app/app .
# Запускаем приложение
CMD ["./app"]