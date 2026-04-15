-- migrations/20250115130000_create_users_table.down.sql
-- Откат миграции

-- Удаляем триггеры
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаляем таблицы (в правильном порядке из-за внешних ключей)
DROP TABLE IF EXISTS login_attempts;
DROP TABLE IF EXISTS users;

-- Удаляем расширение (если оно больше не нужно)
-- DROP EXTENSION IF EXISTS "pgcrypto"; -- Осторожно: может использоваться другими таблицами