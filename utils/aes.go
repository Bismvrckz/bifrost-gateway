package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

func Encrypt(plaintext string) (string, error) {
	trimmedKey := strings.TrimSpace(string(AesEncKey))

	// decodedKey, err := base64.StdEncoding.DecodeString(trimmedKey)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to base64 decode key: %w", err)
	// }

	keyBytes := []byte(trimmedKey)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Seal appends the tag to the ciphertext
	// We also append the result to the nonce to create the final packet:
	// Nonce + Ciphertext + Tag
	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string) (string, error) {
	trimmedKey := strings.TrimSpace(string(AesEncKey))

	// decodedKey, err := base64.StdEncoding.DecodeString("")
	// if err != nil {
	// 	return "", fmt.Errorf("failed to base64 decode key: %w", err)
	// }

	keyBytes := []byte(trimmedKey)

	data, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
