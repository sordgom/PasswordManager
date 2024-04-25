package pkg

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
)

type Password struct {
	Id   uuid.UUID
	Name string
	Url  string

	username string
	hash     string

	Hint string
}

func (p *Password) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (v *Vault) NewPassword(name, url, username, password, hint string) Password {
	id := uuid.New()

	salt, err := RandomSecret(16)
	if err != nil {
		log.Fatal(err)
	}
	driveKey := DeriveKey(v.masterPassword, salt)
	hash, err := v.EncryptData(password, driveKey)
	if err != nil {
		log.Fatal(err)
	}

	return Password{
		Id:       id,
		Name:     name,
		Url:      url,
		username: username,
		hash:     hash,
		Hint:     hint,
	}
}

func (v *Vault) GetPassword(id uuid.UUID) (Password, error) {
	for _, password := range v.Passwords {
		if password.Id == id {
			return password, nil
		}
	}
	return Password{}, errors.New("Password not found")
}

func (v *Vault) GetPasswords() [][]string {
	// Return a list of password names and hints
	var result [][]string
	for _, password := range v.Passwords {
		result = append(result, []string{password.Name, password.Hint})
	}
	return result
}

func (v *Vault) UpdatePassword(id uuid.UUID, name, url, username, password, hint, masterPassword string) error {
	for i, passwordVar := range v.Passwords {
		if passwordVar.Id == id {
			salt, err := RandomSecret(16)
			if err != nil {
				log.Fatal(err)
			}
			driveKey := DeriveKey(masterPassword, salt)
			hash, err := v.EncryptData(password, driveKey)
			if err != nil {
				log.Fatal(err)
			}

			v.Passwords[i] = Password{
				Id:       id,
				Name:     name,
				Url:      url,
				username: username,
				hash:     hash,
				Hint:     hint,
			}
			return nil
		}
	}
	return errors.New("Password not found")
}

func (v *Vault) DeletePassword(id uuid.UUID) error {
	for i, password := range v.Passwords {
		if password.Id == id {
			v.Passwords = append(v.Passwords[:i], v.Passwords[i+1:]...)
			return nil
		}
	}
	return errors.New("Password not found")
}

func (vault *Vault) AppendPassword(password Password) {
	vault.Passwords = append(vault.Passwords, password)
}

func (p *Password) ReadPassword(masterPassword string) string {
	password, err := DecryptPassword(p.hash, masterPassword)
	if err != nil {
		log.Fatal(err)
	}
	return password
}
