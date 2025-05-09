package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	saltSize = 16
)

type StringCipher interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type aesGCMCipher struct{ secret string }

func NewAESGCM(secret string) StringCipher {
	return &aesGCMCipher{secret: secret}
}
func (c *aesGCMCipher) Encrypt(plaintext string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate cryptographic salt: %w", err)
	}

	key := deriveKey(c.secret, salt, 32)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to initialize password credentials: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to initialize GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	fullData := append(salt, append(nonce, ciphertext...)...)
	return base64.URLEncoding.EncodeToString(fullData), nil
}
func (c *aesGCMCipher) Decrypt(ciphertext string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 payload: %w", err)
	}
	if len(data) < saltSize {
		return "", errors.New("invalid ciphertext format")
	}
	salt := data[:saltSize]
	remaining := data[saltSize:]
	key := deriveKey(c.secret, salt, 32)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to initialize password credentials: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to initialize GCM: %w", err)
	}

	if len(remaining) < gcm.NonceSize() {
		return "", errors.New("insufficient ciphertext length")
	}
	nonce := remaining[:gcm.NonceSize()]
	cipherBytes := remaining[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to open GCM: %w", err)
	}
	return string(plaintext), nil
}

type xChaCha20Cipher struct {
	secret string
}

func NewXChaCha20(secret string) StringCipher {
	return &xChaCha20Cipher{secret: secret}
}
func (c *xChaCha20Cipher) Encrypt(plaintext string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate cryptographic salt: %w", err)
	}

	key := deriveKey(c.secret, salt, chacha20poly1305.KeySize)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", fmt.Errorf("failed to initialize XChaCha: %w", err)
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, []byte(plaintext), nil)
	fullData := append(salt, append(nonce, ciphertext...)...)
	return base64.URLEncoding.EncodeToString(fullData), nil
}
func (c *xChaCha20Cipher) Decrypt(ciphertext string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 payload: %w", err)
	}
	if len(data) < saltSize {
		return "", errors.New("invalid ciphertext format")
	}
	salt := data[:saltSize]
	remaining := data[saltSize:]
	key := deriveKey(c.secret, salt, chacha20poly1305.KeySize)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", fmt.Errorf("failed to initialize XChaCha: %w", err)
	}
	if len(remaining) < aead.NonceSize() {
		return "", errors.New("insufficient ciphertext length")
	}
	nonce := remaining[:aead.NonceSize()]
	cipherBytes := remaining[aead.NonceSize():]
	plaintext, err := aead.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to open XChaCha: %w", err)
	}
	return string(plaintext), nil
}

func deriveKey(password string, salt []byte, keyLen uint32) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		3,       // Key derivation rounds
		64*1024, // Mem limit
		4,       // Degree of concurrency
		keyLen,
	)
}
