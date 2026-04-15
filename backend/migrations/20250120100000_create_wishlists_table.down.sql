-- migrations/20250120100000_create_wishlists_table.down.sql
-- Удаление таблицы вишлистов

DROP TRIGGER IF EXISTS update_wishlists_updated_at ON wishlists;
DROP FUNCTION IF EXISTS update_wishlist_updated_at_column();
DROP TABLE IF EXISTS wishlists;