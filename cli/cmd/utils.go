package cmd

import (
	"crypto/rand"
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
