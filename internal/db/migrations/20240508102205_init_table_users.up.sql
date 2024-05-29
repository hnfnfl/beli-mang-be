CREATE TYPE role_user AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(30) PRIMARY KEY NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT null,
    role role_user,
    email_role VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_username on users(username);
CREATE INDEX idx_email_role on users(email_role);
CREATE INDEX idx_role on users(role);
