package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Vault struct {
	Name           string
	masterPassword string

	passwords []Password
}

func New(name, masterPassword string) Vault {
	return Vault{
		Name:           name,
		masterPassword: masterPassword,
	}
}

func (v *Vault) VerifyMasterPassword(newPassword string) bool {
	fmt.Print("MP", v.masterPassword)
	return v.masterPassword == newPassword
}

func SaveVaultToRedis(client *redis.Client, vault *Vault) error {
	ctx := context.Background()

	serializedVault, err := json.Marshal(vault)
	fmt.Printf("SV %s", serializedVault)
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
