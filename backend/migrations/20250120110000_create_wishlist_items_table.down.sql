-- migrations/20250120110000_create_wishlist_items_table.down.sql
-- Удаление таблицы позиций вишлистов

DROP TRIGGER IF EXISTS update_wishlist_items_updated_at ON wishlist_items;
DROP FUNCTION IF EXISTS update_item_updated_at_column();
DROP TABLE IF EXISTS wishlist_items;