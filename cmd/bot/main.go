package main

import (
	"log"
	"os"

	"github.com/kullaniciadin/defi-copilot/internal/ai"
	"github.com/kullaniciadin/defi-copilot/internal/blockchain"
	"github.com/kullaniciadin/defi-copilot/internal/bot"
	"github.com/kullaniciadin/defi-copilot/internal/config"
	"github.com/kullaniciadin/defi-copilot/internal/repository"
	"github.com/kullaniciadin/defi-copilot/internal/wallet"
)

func main() {
	// 1. Config & Env
	cfg := config.LoadConfig()

	// 2. Database
	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	db.Migrate() // Ensure tables exist

	// 3. Blockchain
	ethClient, err := blockchain.NewEthereumClient(cfg.EthereumRPCURL)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Services (AI & Wallet)
	aiService := ai.NewService(os.Getenv("OPENAI_API_KEY"))
	walletSecret := os.Getenv("WALLET_SECRET_KEY")
	if walletSecret == "" {
		walletSecret = ""
	}
	walletService := wallet.NewService(walletSecret)

	// 5. Start Bot
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Please add TELEGRAM_BOT_TOKEN to .env file!")
	}

	defiBot := bot.NewBot(botToken, db, aiService, walletService, ethClient)
	defiBot.Start()
}
