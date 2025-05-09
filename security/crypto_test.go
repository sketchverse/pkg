package security

import (
	"testing"
)

var sk = "BYq04YTu7DzpjEOEn5CET4ACtIpAIj3dF"

func TestEncryptAESGCM(t *testing.T) {
	cipher := NewAESGCM(sk)
	plaintext := "Hello, world!"
	ciphertext, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Error(err)
	}
	t.Log("ciphertext:", ciphertext)
}

func TestDecryptAESGCM(t *testing.T) {
	cipher := NewAESGCM(sk)
	ciphertext := "sywommIgrgQnjjwJRzfPL74B2GjqkKRYMthr-7EyeOBJSWRqGl-mA2C1VVUq7h1B7Y21z1FDXhkJ"
	plaintext, err := cipher.Decrypt(ciphertext)
	if err != nil {
		t.Error(err)
	}
	t.Log("plaintext:", plaintext)
}

func TestEncryptXChaCha20(t *testing.T) {
	cipher := NewXChaCha20(sk)
	plaintext := "Hello, world!"
	ciphertext, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Error(err)
	}
	t.Log("ciphertext:", ciphertext)
}

func TestDecryptXChaCha20(t *testing.T) {
	cipher := NewXChaCha20(sk)
	ciphertext := "h7qgqIep5emziYUG-DvLhuVcWbB-Hjc7TRQswBHfScbULSGNbRSdTWJW-sFC9QnHyQG0-qeIaDdLRy6c1X5JHWa5Kufo"
	plaintext, err := cipher.Decrypt(ciphertext)
	if err != nil {
		t.Error(err)
	}
	t.Log("plaintext:", plaintext)
}
