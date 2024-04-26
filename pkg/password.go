package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Password struct {
	Id   uuid.UUID
	Name string
	Url  string

	Username string
	Hash     string

	Hint string
}

func (p *Password) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (v *Vault) NewPassword(name, url, username, password, hint string) Password {
	id := uuid.New()

	hash, err := v.EncryptPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	return Password{
		Id:       id,
		Name:     name,
		Url:      url,
		Username: username,
		Hash:     hash,
		Hint:     hint,
	}
}

func (v *Vault) GetPassword(name string) (Password, error) {
	for _, password := range v.Passwords {
		if password.Name == name {
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

func (v *Vault) UpdatePassword(name, MasterPassword, password string) error {
	if v.MasterPassword != MasterPassword {
		fmt.Println("incorrect master password")
		return errors.New("incorrect master password")
	}
	for i, passwordVar := range v.Passwords {
		if passwordVar.Name == name {

			hash, err := v.EncryptPassword(password)
			if err != nil {
				log.Fatal(err)
			}

			v.Passwords[i].Name = name
			v.Passwords[i].Hash = hash

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

func (v *Vault) ReadPassword(p *Password) string {
	password, err := v.DecryptPassword(p.Hash)
	if err != nil {
		log.Fatal(err)
	}
	return password
}
