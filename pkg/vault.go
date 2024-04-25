package pkg

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Vault struct {
	Name           string
	masterPassword string

	Passwords []Password
}

func (v *Vault) New(name, masterPassword string) {
	v.Name = name
	v.masterPassword = masterPassword
}

func SaveVaultToRedis(client *redis.Client, vault *Vault) error {
	ctx := context.Background()

	serializedVault, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	return client.Set(ctx, vault.Name, serializedVault, 0).Err() //Will make this better I swear
}

func LoadVaultFromRedis(client *redis.Client, vaultName string) (*Vault, error) {
	ctx := context.Background()
	serializedVault, err := client.Get(ctx, vaultName).Bytes()
	if err != nil {
		return nil, err
	}

	var vault Vault
	err = json.Unmarshal(serializedVault, &vault)
	if err != nil {
		return nil, err
	}

	return &vault, nil
}
