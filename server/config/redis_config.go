package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sordgom/PasswordManager/server/model"

	"github.com/redis/go-redis/v9"
)

type VaultService interface {
	SaveVaultToRedis(vault *model.Vault) error
	LoadVaultFromRedis(vaultName string) (*model.Vault, error)
}

// This service will handle all redis operations
type RedisVaultService struct {
	client *redis.Client
	Vault  *model.Vault
}

func NewRedisVaultService(client *redis.Client) *RedisVaultService {
	return &RedisVaultService{client: client}
}

func (r *RedisVaultService) SaveVaultToRedis(vault *model.Vault) error {
	ctx := context.Background()

	serializedVault, err := json.Marshal(vault)
	fmt.Printf("SV %s", serializedVault)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, vault.Name, serializedVault, 0).Err()
}

func (r *RedisVaultService) LoadVaultFromRedis(vaultName string) (*model.Vault, error) {
	ctx := context.Background()

	serializedVault, err := r.client.Get(ctx, vaultName).Bytes()
	if err != nil {
		return nil, err
	}

	var vault model.Vault
	err = json.Unmarshal(serializedVault, &vault)
	if err != nil {
		return nil, err
	}

	return &vault, nil
}
