CREATE TYPE role_user AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(30) PRIMARY KEY NOT NULL,
    email VARCHAR(30) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role role_user,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_username ON users(username);

CREATE INDEX idx_email ON users(email);

CREATE INDEX idx_role ON users(role);