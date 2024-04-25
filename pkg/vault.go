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

	return client.Set(ctx, "vaultKey", serializedVault, 0).Err() //Will make this better I swear
}

func LoadVaultFromRedis(client *redis.Client) (*Vault, error) {
	ctx := context.Background()
	serializedVault, err := client.Get(ctx, "vaultKey").Bytes()
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
