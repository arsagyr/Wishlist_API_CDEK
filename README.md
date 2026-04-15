# Wishlist API

REST API сервис для управления вишлистами к праздникам и событиям.

## Возможности

- Регистрация и аутентификация пользователей (JWT)
- Создание, просмотр, редактирование и удаление вишлистов
- Управление позициями (подарками) внутри вишлистов
- Публичный доступ к вишлистам по уникальной ссылке
- Бронирование подарков без авторизации

## Технический стек

- Go 1.21+
- PostgreSQL
- Docker Compose
- golang-migrate

## Запуск

```bash
docker-compose up --build
```

## API Endpoints

### Аутентификация
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход

### Вишлисты (требуют авторизации)
- `GET /api/wishlists` - Список вишлистов пользователя
- `POST /api/wishlists` - Создать вишлист
- `GET /api/wishlists/{id}` - Получить вишлист
- `PUT /api/wishlists/{id}` - Обновить вишлист
- `DELETE /api/wishlists/{id}` - Удалить вишлист

### Позиции вишлиста (требуют авторизации)
- `GET /api/wishlists/{wishlistId}/items` - Список позиций
- `POST /api/wishlists/{wishlistId}/items` - Добавить позицию
- `GET /api/wishlists/{wishlistId}/items/{id}` - Получить позицию
- `PUT /api/wishlists/{wishlistId}/items/{id}` - Обновить позицию
- `DELETE /api/wishlists/{wishlistId}/items/{id}` - Удалить позицию

### Публичные endpoints
- `GET /api/public/wishlists/{token}` - Получить вишлист по токену
- `POST /api/public/wishlists/{token}/reserve` - Забронировать подарок

## Примеры запросов

### Регистрация
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "Password123!"}'
```

### Вход
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "Password123!"}'
```

### Создание вишлиста
```bash
curl -X POST http://localhost:8080/api/wishlists \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title": "День рождения", "description": "Хочу подарки", "event_date": "2025-06-15"}'
```

### Публичный доступ
```bash
curl http://localhost:8080/api/public/wishlists/{token}
```

### Бронирование подарка
```bash
curl -X POST http://localhost:8080/api/public/wishlists/{token}/reserve \
  -H "Content-Type: application/json" \
  -d '{"item_id": "<item-id>", "reserved_by": "Иван"}'
```

## Переменные окружения

- `DATABASE_URL` - URL подключения к PostgreSQL
- `JWT_SECRET` - Секретный ключ для JWT токенов
- `APP_PORT` - Порт приложения (по умолчанию 8080)