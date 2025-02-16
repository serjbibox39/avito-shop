-- Создание таблицы пользователей
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	role TEXT NOT NULL DEFAULT 'user',
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

-- Создание таблицы товаров
CREATE TABLE merch (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    price INT NOT NULL
);

-- Создание таблицы инвентаря
CREATE TABLE inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid UUID NOT NULL,
	merch TEXT[],
    coins INT CHECK(coins >= 0) DEFAULT 1000
);

-- Создание таблицы транзакций
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_user UUID REFERENCES users(id),
    to_user UUID REFERENCES users(id),
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

WITH i AS (INSERT INTO users(username, role, password) VALUES ('store', 'store', '$2a$10$0uQ73ZFh7feFGwnz9UTWt.GVbxw71oRZ6dr8GBFSHKNCNPvca3Rzi') RETURNING ID)
INSERT INTO inventory (userid) SELECT i.id FROM i;

-- Вставка данных о товарах
INSERT INTO merch (name, price) VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500);