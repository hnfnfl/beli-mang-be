CREATE TABLE IF NOT EXISTS order_product (
    order_id VARCHAR(255) PRIMARY KEY NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    merchant_id VARCHAR(255) NOT NULL,
    item_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE IF NOT EXISTS order_product (
--     order_id VARCHAR(50) PRIMARY KEY NOT NULL,
--     user_id VARCHAR(50) FOREIGN KEY NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
-- CREATE TABLE IF NOT EXISTS order_product_detail (
--     order_id VARCHAR(50) PRIMARY KEY NOT NULL,
--     merchant_id VARCHAR(255) NOT NULL,
--     item_id VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
-- ALTER TABLE merchant_items
--     ADD CONSTRAINT fk_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(merchant_id) ON DELETE CASCADE;
CREATE INDEX idx_id_order ON order_product(order_id);

CREATE INDEX idx_user_id_order ON order_product(user_id);

CREATE INDEX idx_merchant_id_order ON order_product(merchant_id);

CREATE INDEX idx_item_id_order ON order_product(item_id);