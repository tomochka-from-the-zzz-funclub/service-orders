# # FROM golang:latest

# # COPY go.mod go.sum ./
# # RUN go mod download

# # COPY . ./
# # WORKDIR ./cmd
# # RUN go build -o main .
# # CMD ["./main"]

# # Этап сборки
# FROM golang:latest AS builder

# # Установим рабочую директорию
# WORKDIR /app

# # Скопируем только зависимости
# COPY go.mod go.sum ./

# # Загрузим зависимости
# RUN go mod download

# # Скопируем остальные файлы проекта
# COPY . ./

# # Соберем приложение
# RUN go build -o main ./cmd

# # Этап выполнения
# FROM alpine:latest

# # Установим рабочую директорию
# WORKDIR /app

# # Копируем скомпилированное приложение из этапа сборки
# COPY --from=builder /app/main .

# # Задаем команду, которая будет выполнена при запуске контейнера
# CMD ["./main"]


# FROM golang:1.23 AS builder

# WORKDIR /app

# COPY . .

# RUN go mod tidy
# RUN go build -o api-gateway

# FROM golang:1.23
# COPY --from=builder /app/api-gateway /api-gateway
# COPY wait-for-it.sh /app/wait-for-it.sh
# RUN chmod +x /app/wait-for-it.sh
# CMD ["/app/wait-for-it.sh", "kafka:9092", "--", "/api-gateway"]



FROM golang:1.23 AS builder
# Устанавливаем рабочую директорию на уровень проекта
WORKDIR /app
# Копируем все файлы проекта в рабочую директорию
COPY . .
# Убедитесь, что переходим в директорию, где находится ваш main.go
WORKDIR /app/cmd
# Упорядочиваем зависимости
RUN go mod tidy
# Собираем исполняемый файл
RUN go build -o api-gateway
FROM golang:1.23
# Копируем собранный исполняемый файл из предыдущего образа
COPY --from=builder /app/cmd/api-gateway /api-gateway
# Копируем wait-for-it.sh в корневую директорию
COPY wait-for-it.sh /app/wait-for-it.sh
# Даем права на выполнение для wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

# Устанавливаем команду для запуска
CMD ["/app/wait-for-it.sh", "kafka:9092", "--", "/api-gateway"]
