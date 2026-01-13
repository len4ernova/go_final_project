# Используем официальный образ Go в качестве базового.
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера.
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей.
COPY go.mod go.sum ./

# Загружаем зависимости.
RUN go mod download

# Копируем исходный код приложения.
COPY . .

# Собираем приложение.
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Используем минимальный образ Alpine Linux для запуска приложения.
FROM alpine:latest

# Устанавливаем рабочую директорию внутри контейнера.
WORKDIR /app

# Устанавливаем глоб. переменнные
ENV TODO_PORT=7540
# Копируем исполняемый файл из билдера.
COPY --from=builder /app/main .
RUN mkdir web
COPY --from=builder /app/web web/

# Открываем порт.
EXPOSE $TODO_PORT

# Запускаем приложение.
CMD ["./main"]
