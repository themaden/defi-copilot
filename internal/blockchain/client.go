package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthereumClient wraps the Geth client
type EthereumClient struct {
	Client *ethclient.Client
}

// NewEthereumClient connects to the network and returns the client instance
func NewEthereumClient(rpcURL string) (*EthereumClient, error) {
	// Set 10 seconds timeout to prevent hanging if network is unresponsive
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum network: %w", err)
	}

	return &EthereumClient{Client: client}, nil
}

// GetETHBalance returns the ETH balance (in Wei) of the given address
func (e *EthereumClient) GetETHBalance(address string) (*big.Float, error) {
	// Convert String address to Hex Address format
	account := common.HexToAddress(address)

	// nil = get balance at latest block
	balanceWei, err := e.Client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}

	// Convert Wei to Ether (1 Ether = 10^18 Wei)
	// Float kullanarak insan okunabilir formata çeviriyoruz
	fBalance := new(big.Float)
	fBalance.SetString(balanceWei.String())
	val := new(big.Float).Quo(fBalance, big.NewFloat(1e18))

	return val, nil
}
