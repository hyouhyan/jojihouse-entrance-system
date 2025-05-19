-- ユーザー
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    barcode VARCHAR(64) UNIQUE NOT NULL,
    contact VARCHAR(255),
    remaining_entries INT DEFAULT 0,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_entries INT DEFAULT 0 CHECK (total_entries >= 0),
    allergy VARCHAR(255),
    number INT UNIQUE
);

-- ロール
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

-- ロールとユーザーの中間テーブル
CREATE TABLE user_roles (
    user_id INT REFERENCES users(id),
    role_id INT REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

-- 今居るユーザーの一覧
CREATE TABLE current_users (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    entered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 認証系の情報
CREATE TABLE auth (
    user_id INT PRIMARY KEY UNIQUE REFERENCES users(id),
    password VARCHAR(255),
);