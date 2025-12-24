package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	EthereumRPCURL string
}

func LoadConfig() *Config {
	// Try to load .env file, check system envs if not found

	if err := godotenv.Load(".env"); err != nil {
		if err := godotenv.Load("../../.env"); err != nil { // If called from subfolder
			log.Println("⚠️ Warning: .env file could not be loaded. Checking system environment variables.")
		}
	}

	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		log.Fatal("ERROR: ETHEREUM_RPC_URL is not set!")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("ERROR: DATABASE_URL environment variable is not set!")
	}

	return &Config{
		DatabaseURL:    dbURL,
		EthereumRPCURL: rpcURL,
	}
}
