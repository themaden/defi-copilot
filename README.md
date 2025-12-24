# DeFi Copilot

DeFi Copilot is a Telegram bot that acts as your personal decentralized finance assistant. It leverages AI to understand natural language commands and helps you manage your crypto assets, check balances, and execute swaps.

## Features

- **AI-Powered Intent Parsing**: Understands natural language commands like "Buy ETH" or "What is my balance?".
- **Wallet Management**: Automatically creates and manages a secure Ethereum wallet for each user.
- **Blockchain Integration**: Connects to the Ethereum network via RPC to fetch balances and potentially execute transactions.
- **Secure Storage**: Encrypts private keys before storing them in the database.

## Architecture

- **cmd/bot**: Entry point of the application.
- **internal/ai**: Handles interaction with OpenAI (or Mock AI) to parse user intents.
- **internal/blockchain**: Manages Ethereum network connections and queries.
- **internal/bot**: Contains the Telegram bot logic and request handlers.
- **internal/config**: Loads configuration from `.env` files.
- **internal/repository**: Manages database connections and user/wallet persistence.
- **internal/wallet**: Handles cryptographic operations like creating wallets and signing transactions.

## Prerequisites

- Go 1.21+
- PostgreSQL
- OpenAI API Key (Optional, defaults to Mock mode)
- Ethereum RPC URL (e.g., Infura, Alchemy)
- Telegram Bot Token

## Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/kullaniciadin/defi-copilot.git
    cd defi-copilot
    ```

2.  Create a `.env` file in the root directory:
    ```env
    DATABASE_URL=postgres://user:password@localhost:5432/defi_copilot?sslmode=disable
    ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/YOUR_INFURA_KEY
    TELEGRAM_BOT_TOKEN=YOUR_TELEGRAM_BOT_TOKEN
    OPENAI_API_KEY=sk-YOUR_OPENAI_KEY
    WALLET_SECRET_KEY=12345678901234567890123456789012 # Must be 32 bytes
    ```

3.  Build and Run:
    ```bash
    go mod tidy
    go run cmd/bot/main.go
    ```

## Usage

Start a chat with your bot on Telegram and send commands like:

-   "Create wallet" (Automatic on /start)
-   "What is my ETH balance?"
-   "Swap 0.5 ETH to USDT"

The bot will interpret your request and respond accordingly.

## License

MIT
