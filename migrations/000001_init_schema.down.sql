-- Удаление индексов
DROP INDEX IF EXISTS idx_items_nm_id;
DROP INDEX IF EXISTS idx_items_order_uid;
DROP INDEX IF EXISTS idx_orders_date_created;

-- Удаление таблиц в обратном порядке (с учетом foreign keys)
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS orders;