-- migrations/001_initial_schema.sql

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    address VARCHAR(42) NOT NULL, -- Ethereum adresi (0x...)
    encrypted_pk TEXT NOT NULL,   -- Şifrelenmiş Private Key
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- İndexler (Sorgu performansı için)
CREATE INDEX idx_users_telegram_id ON users(telegram_id);
CREATE INDEX idx_wallets_address ON wallets(address);