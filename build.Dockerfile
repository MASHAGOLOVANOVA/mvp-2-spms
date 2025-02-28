FROM golang
# Устанавливаем переменные окружения
WORKDIR /app
ENV CGO_ENABLED=0
ENV GOOS=linux

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальное приложение
COPY . .
# Копируем файл server_config.json
COPY web_server/cmd/web_app/server_config.json ./web_server/cmd/web_app/
# Установка make
RUN apt-get update && apt-get install -y make

# Выполнение команды сборки из Makefile
RUN make run-commands-to-build-go

CMD ["./web_server/cmd/web_app/web_app"]

