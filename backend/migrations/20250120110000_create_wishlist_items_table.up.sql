-- migrations/20250120110000_create_wishlist_items_table.up.sql
-- Создание таблицы позиций вишлистов

CREATE TABLE IF NOT EXISTS wishlist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wishlist_id UUID NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    link VARCHAR(2048),
    priority INTEGER DEFAULT 0 CHECK (priority >= 0 AND priority <= 2),
    is_reserved BOOLEAN DEFAULT FALSE,
    reserved_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT item_name_not_empty CHECK (name <> '')
);

CREATE INDEX idx_wishlist_items_wishlist_id ON wishlist_items(wishlist_id);
CREATE INDEX idx_wishlist_items_priority ON wishlist_items(priority);
CREATE INDEX idx_wishlist_items_is_reserved ON wishlist_items(is_reserved);

COMMENT ON TABLE wishlist_items IS 'Таблица позиций (подарков) в вишлистах';
COMMENT ON COLUMN wishlist_items.id IS 'Уникальный идентификатор позиции (UUID)';
COMMENT ON COLUMN wishlist_items.wishlist_id IS 'ID вишлиста';
COMMENT ON COLUMN wishlist_items.name IS 'Название подарка';
COMMENT ON COLUMN wishlist_items.description IS 'Описание подарка';
COMMENT ON COLUMN wishlist_items.link IS 'Ссылка на товар';
COMMENT ON COLUMN wishlist_items.priority IS 'Приоритет (0-low, 1-medium, 2-high)';
COMMENT ON COLUMN wishlist_items.is_reserved IS 'Флаг бронирования';
COMMENT ON COLUMN wishlist_items.reserved_by IS 'Кем забронирован';

CREATE OR REPLACE FUNCTION update_item_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_wishlist_items_updated_at
    BEFORE UPDATE ON wishlist_items
    FOR EACH ROW
    EXECUTE FUNCTION update_item_updated_at_column();