CREATE TYPE merchant_category AS ENUM (
    'SmallRestaurant',
    'MediumRestaurant',
    'LargeRestaurant',
    'MerchandiseRestaurant',
    'BoothKiosk',
    'ConvenienceStore'
);

CREATE TABLE IF NOT EXISTS merchants (
    merchant_id VARCHAR(30) PRIMARY KEY NOT NULL,
    name VARCHAR(30) NOT NULL,
    merchant_categories merchant_category,
    long double precision NOT NULL,
    lat double precision NOT NULL,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_merchant_id ON merchants(merchant_id);

CREATE INDEX idx_name_merchant ON merchants(name);

CREATE INDEX idx_merchant_category ON merchants(merchant_categories);

CREATE INDEX IF NOT EXISTS merchants_created_at_desc ON merchants(created_at DESC);

CREATE INDEX IF NOT EXISTS merchants_created_at_asc ON merchants(created_at ASC);