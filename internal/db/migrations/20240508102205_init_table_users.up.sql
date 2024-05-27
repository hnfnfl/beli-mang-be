CREATE TYPE IF NOT EXISTS role_user AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(30) PRIMARY KEY NOT NULL,
<<<<<<< Updated upstream
    email VARCHAR(30) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role role_user,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_username ON users(username);

CREATE INDEX idx_email ON users(email);

CREATE INDEX idx_role ON users(role);
=======
    email VARCHAR(30) NOT NULL,
    password_hash VARCHAR(255) NOT null,
    role role_user,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    email_role VARCHAR(35) UNIQUE NOT NULL
);

CREATE INDEX idx_username on users(username);
CREATE INDEX idx_email_role on users(email_role);
CREATE INDEX idx_role on users(role);
>>>>>>> Stashed changes
