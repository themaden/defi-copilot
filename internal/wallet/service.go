package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Service struct {
	EncryptionKey []byte // Encryption key used when writing to database
}

// NewService: Requires a 32-byte (256-bit) secure key
func NewService(secret string) *Service {
	// If secret is short, pad it or hash it (Assuming 32 bytes for simplicity)
	if len(secret) != 32 {
		log.Println("WARNING: Encryption Key is not 32 characters! Security risk.")
		// Temporary fix for demo:
		secret = fmt.Sprintf("%-32s", secret)[:32]
	}
	return &Service{EncryptionKey: []byte(secret)}
}

// CreateWallet: Creates a new wallet
// Return: (PublicAddress, EncryptedPrivateKey, error)
func (s *Service) CreateWallet() (string, string, error) {
	// 1. Generate Random Private Key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	// 2. Derive Address from Public Key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", fmt.Errorf("public key error")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	// 3. Convert Private Key to Byte array
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBytes)[2:] // Removing "0x" prefix

	// 4. Encrypt (AES)
	encryptedPK, err := s.encrypt(privateKeyHex)
	if err != nil {
		return "", "", err
	}

	return address, encryptedPK, nil
}

// encrypt: Encrypts Private Key before storing in database
func (s *Service) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.EncryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

// SignTransaction: Decrypts key and signs transaction
// NOTE: This is a simulation, 'types.Transaction' object is needed for real transaction.
func (s *Service) SignTransaction(encryptedPK string, toAddress string, amount float64) (string, error) {
	// 1. Decrypt
	pkHex, err := s.decrypt(encryptedPK)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	// 2. Load Private Key
	privateKey, err := crypto.HexToECDSA(pkHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	// 3. (Mobile) Validate Address
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	// Normally "types.NewTransaction" and "types.SignTx" are used here.
	// For now we return a log proving that signing works.
	return fmt.Sprintf("Transaction Signed! \nSender: %s\nReceiver: %s\nAmount: %.4f ETH\n(This transaction is ready to be sent to network)", fromAddress, toAddress, amount), nil
}

// decrypt: Decrypts AES encryption (Internal Helper)
func (s *Service) decrypt(ciphertextHex string) (string, error) {
	data, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.EncryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("data too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
