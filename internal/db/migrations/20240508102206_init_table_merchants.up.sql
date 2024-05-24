CREATE TYPE merchant_category AS ENUM ('SmallRestaurant', 'MediumRestaurant', 'LargeRestaurant', 'MerchandiseRestaurant', 'BoothKiosk', 'ConvenienceStore');

CREATE TABLE IF NOT EXISTS merchants (
    merchant_id VARCHAR(30) PRIMARY KEY NOT NULL,
    name VARCHAR(30) NOT NULL,
    merchant_categories merchant_category
    long float64 not null,
    lat float64 not null,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

CREATE INDEX idx_merchant_id on merchants(merchant_id);
CREATE INDEX idx_name_merchant on merchants(name);
CREATE INDEX idx_merchant_category on merchants(merchant_categories);
CREATE INDEX IF NOT EXISTS merchants_created_at_desc
    ON merchants(created_at desc);
CREATE INDEX IF NOT EXISTS merchants_created_at_asc
    ON merchants(created_at asc);