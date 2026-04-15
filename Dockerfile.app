FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Установка migrate CLI
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/app \
    ./cmd/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl postgresql-client bash

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /go/bin/migrate .
COPY --from=builder /app/backend/migrations ./migrations

ENV DATABASE_URL=postgres://postgres:password@localhost:5432/mydb?sslmode=disable \
    JWT_SECRET=your-secret-key-change-in-production \
    APP_PORT=8080 \
    TZ=Europe/Moscow

EXPOSE 8080

# Скрипт запуска с миграциями
RUN echo '#!/bin/sh\n\
echo "Running migrations..."\n\
/app/migrate -path /app/migrations -database "$DATABASE_URL" up || true\n\
echo "Starting application..."\n\
exec /app/app\n\
' > /start.sh && chmod +x /start.sh

CMD ["/start.sh"]