-- migrations/20250120100000_create_wishlists_table.up.sql
-- Создание таблицы вишлистов

CREATE TABLE IF NOT EXISTS wishlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date DATE,
    public_token VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT title_not_empty CHECK (title <> '')
);

CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);
CREATE INDEX idx_wishlists_public_token ON wishlists(public_token);
CREATE INDEX idx_wishlists_created_at ON wishlists(created_at);

COMMENT ON TABLE wishlists IS 'Таблица вишлистов пользователей';
COMMENT ON COLUMN wishlists.id IS 'Уникальный идентификатор вишлиста (UUID)';
COMMENT ON COLUMN wishlists.user_id IS 'ID владельца вишлиста';
COMMENT ON COLUMN wishlists.title IS 'Название события/праздника';
COMMENT ON COLUMN wishlists.description IS 'Описание вишлиста';
COMMENT ON COLUMN wishlists.event_date IS 'Дата события';
COMMENT ON COLUMN wishlists.public_token IS 'Публичный токен для доступа по ссылке';

CREATE OR REPLACE FUNCTION update_wishlist_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_wishlists_updated_at
    BEFORE UPDATE ON wishlists
    FOR EACH ROW
    EXECUTE FUNCTION update_wishlist_updated_at_column();