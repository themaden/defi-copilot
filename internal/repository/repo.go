package repository

import (
	"fmt"
	"time"
)

// User Model
type User struct {
	ID         int       `db:"id"`
	TelegramID int64     `db:"telegram_id"`
	Username   string    `db:"username"`
	CreatedAt  time.Time `db:"created_at"`
}

// SaveUser: Saves a new user (Returns ID if exists)
func (db *DB) SaveUser(telegramID int64, username string) (int, error) {
	var id int
	// Using "ON CONFLICT" to just get ID if user exists
	query := `
		INSERT INTO users (telegram_id, username) 
		VALUES ($1, $2) 
		ON CONFLICT (telegram_id) DO UPDATE SET username = EXCLUDED.username
		RETURNING id`

	err := db.QueryRow(query, telegramID, username).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save user: %w", err)
	}
	return id, nil
}

// SaveWallet: Connects wallet to user and saves it
func (db *DB) SaveWallet(userID int, address, encryptedPK string) error {
	query := `
		INSERT INTO wallets (user_id, address, encrypted_pk) 
		VALUES ($1, $2, $3)`

	_, err := db.Exec(query, userID, address, encryptedPK)
	if err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}
	return nil
}

// GetWallet: Fetches user's wallet info
func (db *DB) GetWallet(userID int) (string, string, error) {
	var address, encryptedPK string
	query := `SELECT address, encrypted_pk FROM wallets WHERE user_id = $1 LIMIT 1`

	err := db.QueryRow(query, userID).Scan(&address, &encryptedPK)
	if err != nil {
		return "", "", fmt.Errorf("wallet not found: %w", err)
	}
	return address, encryptedPK, nil
}
