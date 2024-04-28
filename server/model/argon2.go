package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/argon2"
)

func Compare(password, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatal(err)
	}

	return match
}
func GenerateHash(password string) string {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	return hash
}

func GenerateRandomHash() string {
	secret, err := RandomSecret(32)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := argon2id.CreateHash(string(secret), argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	return hash
}

func RandomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// Derive a key from the master password
func DeriveKey(MasterPassword string, salt []byte) []byte {
	return argon2.IDKey([]byte(MasterPassword), salt, 1, 64*1024, 4, 32) // Returns a 32-byte key
}

func (v *Vault) EncryptPassword(data string) (string, error) {
	salt, err := RandomSecret(16)
	if err != nil {
		log.Fatal(err)
	}
	driveKey := DeriveKey(v.MasterPassword, salt)

	// Encrypt the data
	block, err := aes.NewCipher(driveKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(data), nil)

	// Combine salt, nonce, and ciphertext for the final encrypted data
	finalCiphertext := append(salt, nonce...)
	finalCiphertext = append(finalCiphertext, ciphertext...)

	// Encode the combined slice to base64
	return base64.URLEncoding.EncodeToString(finalCiphertext), nil
}

func (v *Vault) DecryptPassword(encryptedData string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	// Validate minimum size
	saltSize := 16
	if len(data) < saltSize {
		return "", errors.New("encrypted data too short to contain salt")
	}

	salt := data[:saltSize]
	data = data[saltSize:]

	// Derive the key
	key := DeriveKey(v.MasterPassword, salt)

	// Decrypt the data
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("encrypted data too short to contain nonce")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}
	return string(plaintext), nil
}
