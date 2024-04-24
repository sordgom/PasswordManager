package pkg

import (
	"github.com/google/uuid"
)

type Vault struct {
	Name           string
	masterPassword string

	passwords []Password
}

type Password struct {
	id   uuid.UUID
	name string
	url  string

	hashedPassword string
	hint           string
}

func CreateVault(name string, masterPassword string) Vault {
	vault := Vault{Name: name, masterPassword: masterPassword}
	return vault
}

func (vault *Vault) AppendPassword(password Password) {
	vault.passwords = append(vault.passwords, password)
}

func CreatePassword(name, url, password, hint string) Password {
	id := uuid.New()
	hashedPassword := password // for now
	return Password{
		id:             id,
		name:           name,
		url:            url,
		hashedPassword: hashedPassword,
		hint:           hint,
	}
}

func GenerateHash() string {
	return "hashed_password"
}
