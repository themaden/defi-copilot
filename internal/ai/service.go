package ai

import (
	"context"
	"encoding/json"

	//"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Service struct {
	client *openai.Client
	isMock bool // If true, mock AI runs instead of real AI
}

func NewService(apiKey string) *Service {
	// If key is empty or doesn't start with "sk-", enable Mock mode
	if apiKey == "" || !strings.HasPrefix(apiKey, "sk-") {
		log.Println("⚠️ OpenAI Key missing/invalid. 'MOCK AI' (Simulation Mode) activated.")
		return &Service{isMock: true}
	}

	client := openai.NewClient(apiKey)
	return &Service{client: client, isMock: false}
}

type UserIntent struct {
	Intent string  `json:"intent"`
	Asset  string  `json:"asset"`
	Amount float64 `json:"amount"`
}

// ParseIntent: Selects real AI or Mock logic
func (s *Service) ParseIntent(userMessage string) (*UserIntent, error) {
	if s.isMock {
		return s.mockParse(userMessage)
	}
	return s.openaiParse(userMessage)
}

// --- OPTION A: REAL AI ---
func (s *Service) openaiParse(userMessage string) (*UserIntent, error) {
	systemPrompt := `
    You are a DeFi assistant. Respond in JSON:
    {"intent": "swap/balance/transfer", "asset": "ETH/USDT", "amount": 0}
    `
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: userMessage},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	var intent UserIntent
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &intent); err != nil {
		return nil, err
	}
	return &intent, nil
}

// --- OPTION B: FREE MOCK (SIMULATION) ---
func (s *Service) mockParse(text string) (*UserIntent, error) {
	text = strings.ToLower(text)
	intent := &UserIntent{}

	// 1. Determine Intent (Simple keyword check)
	if strings.Contains(text, "buy") || strings.Contains(text, "swap") || strings.Contains(text, "trade") {
		intent.Intent = "swap"
	} else if strings.Contains(text, "balance") || strings.Contains(text, "how much") {
		intent.Intent = "balance"
	} else {
		intent.Intent = "unknown"
	}

	// 2. Determine Asset
	if strings.Contains(text, "eth") {
		intent.Asset = "ETH"
	} else if strings.Contains(text, "usdt") {
		intent.Asset = "USDT"
	}

	// 3. Determine Amount (Regex number capture)
	re := regexp.MustCompile(`[0-9]+(\.[0-9]+)?`)
	match := re.FindString(text)
	if match != "" {
		amount, _ := strconv.ParseFloat(match, 64)
		intent.Amount = amount
	}

	// Simulate AI delay (for realism)
	// time.Sleep(500 * time.Millisecond)

	return intent, nil
}
