CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    barcode CHAR(13) UNIQUE NOT NULL,
    contact VARCHAR(255),
    remaining_entries INT DEFAULT 0 CHECK (remaining_entries >= 0),
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_entries INT DEFAULT 0 CHECK (total_entries >= 0),
    total_stay_time INTERVAL DEFAULT '0 minute'
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- 初期データの注入
INSERT INTO roles (name) 
VALUES 
    ('member'),
    ('student'),
    ('system-admin'),
    ('house-admin'),
    ('guest');

CREATE TABLE user_roles (
    user_id INT REFERENCES users(id),
    role_id INT REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);