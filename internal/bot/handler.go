package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/kullaniciadin/defi-copilot/internal/ai"
	"github.com/kullaniciadin/defi-copilot/internal/blockchain"
	"github.com/kullaniciadin/defi-copilot/internal/repository"
	"github.com/kullaniciadin/defi-copilot/internal/wallet"
	tele "gopkg.in/telebot.v3"
)

// DeFiBot: Struct that holds all services
type DeFiBot struct {
	Bot        *tele.Bot
	DB         *repository.DB
	AI         *ai.Service
	Wallet     *wallet.Service
	Blockchain *blockchain.EthereumClient
}

// NewBot: Initializes and configures the bot
func NewBot(token string, db *repository.DB, ai *ai.Service, w *wallet.Service, bc *blockchain.EthereumClient) *DeFiBot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Bot failed to start: %v", err)
	}

	bot := &DeFiBot{Bot: b, DB: db, AI: ai, Wallet: w, Blockchain: bc}

	// --- COMMANDS ---
	b.Handle("/start", bot.handleStart)
	b.Handle(tele.OnText, bot.handleMessage)

	return bot
}

// Start: Starts listening for updates
func (b *DeFiBot) Start() {
	log.Println("🤖 Telegram Bot Listening...")
	b.Bot.Start()
}

// /start command
func (b *DeFiBot) handleStart(c tele.Context) error {
	user := c.Sender()

	// 1. Register User
	userID, err := b.DB.SaveUser(user.ID, user.Username)
	if err != nil {
		return c.Send("❌ System error: Could not save user.")
	}

	// 2. Check Exising Wallet
	_, _, err = b.DB.GetWallet(userID)
	if err != nil {
		// Create wallet if none exists
		addr, encPK, _ := b.Wallet.CreateWallet()
		b.DB.SaveWallet(userID, addr, encPK)
		return c.Send(fmt.Sprintf("👋 Welcome! I created a secure wallet for you.\n\n📬 Address: `%s`\n\nYou can ask me 'What is my balance?' or 'Buy ETH'.", addr))
	}

	return c.Send("👋 Welcome back! Your wallet is ready. How can I help you?")
}

// Normal messages (AI Integration)
func (b *DeFiBot) handleMessage(c tele.Context) error {
	msg := c.Text()

	// 1. Parse Intent with AI
	intent, err := b.AI.ParseIntent(msg)
	if err != nil {
		return c.Send("🧠 I'm confused, could you repeat that?")
	}

	// 2. Execute Action based on Intent
	switch intent.Intent {
	case "balance":
		// Instead of Vitalik's address, you should fetch the user's address from DB
		// Showing demo/mock response for now.
		// Advanced: db.GetWallet() -> blockchain.GetBalance(addr)
		return c.Send("💰 Balance query requested. This feature is currently in maintenance (Demo).")

	case "swap":
		return c.Send(fmt.Sprintf("🔄 Swap Detected!\nAsset: %s\nAmount: %.2f\n\nPreparing transaction (signature required)...", intent.Asset, intent.Amount))

	default:
		return c.Send(fmt.Sprintf("🤔 I didn't fully understand your intent, but you said: %s", msg))
	}
}
