CREATE TYPE product_category AS ENUM ('Beverage', 'Food', 'Snack', 'Condiments', 'Additions');

CREATE TABLE IF NOT EXISTS merchant_items (
    item_id VARCHAR(30) PRIMARY KEY NOT NULL,
    merchant_id VARCHAR(30) FOREIGN KEY KEY NOT NULL,
    name VARCHAR(30) NOT NULL,
    product_categories product_category
    price int64 not null,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

-- ALTER TABLE merchant_items
--     ADD CONSTRAINT fk_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(merchant_id) ON DELETE CASCADE;

CREATE INDEX idx_item_id_item on merchant_items(item_id);
CREATE INDEX idx_merchant_id_item on merchant_items(merchant_id);
CREATE INDEX idx_product_category_item on merchant_items(product_categories);
CREATE INDEX idx_name_item on merchant_items(name);
CREATE INDEX IF NOT EXISTS merchant_items_created_at_desc
    ON merchant_items(created_at desc);
CREATE INDEX IF NOT EXISTS merchant_items_created_at_asc
    ON merchant_items(created_at asc);