package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"

	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/argon2"
)

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
func DeriveKey(masterPassword string, salt []byte) []byte {
	return argon2.IDKey([]byte(masterPassword), salt, 1, 64*1024, 4, 32) // Returns a 32-byte key
}

func (v *Vault) EncryptData(data string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
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
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptPassword(encryptedData string, masterPassword string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	saltSize := 16  // For example, 16 bytes salt
	nonceSize := 12 // For example, 12 bytes nonce (typical for GCM)

	// Validate sizes
	if len(data) < saltSize+nonceSize { // Minimum length check
		return "", errors.New("encrypted data is too short")
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	key := DeriveKey(masterPassword, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
