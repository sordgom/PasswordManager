package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sordgom/PasswordManager/server/model"
)

type VaultService interface {
	SaveVaultToRedis(ctx context.Context, vault *model.Vault) error
	LoadVaultFromRedis(ctx context.Context, vaultName string) (*model.Vault, error)
}

type RedisVaultService struct {
	client *redis.Client
}

func NewRedisVaultService(client *redis.Client) VaultService {
	return &RedisVaultService{client: client}
}

func (r *RedisVaultService) SaveVaultToRedis(ctx context.Context, vault *model.Vault) error {
	serializedVault, err := json.Marshal(vault)
	if err != nil {
		return err
	}
	fmt.Printf("Serialized Vault: %s\n", serializedVault)
	return r.client.Set(ctx, vault.Name, serializedVault, 0).Err()
}

func (r *RedisVaultService) LoadVaultFromRedis(ctx context.Context, vaultName string) (*model.Vault, error) {
	serializedVault, err := r.client.Get(ctx, vaultName).Bytes()
	if err != nil {
		return nil, err
	}

	var vault model.Vault
	if err = json.Unmarshal(serializedVault, &vault); err != nil {
		return nil, err
	}

	return &vault, nil
}
