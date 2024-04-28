package cmd

import (
	"errors"
	"log"
)

func load() (AppContext, error) {
	client := appContext.Client
	vault := appContext.Vault
	if vault == nil {
		log.Fatal("Failed to find Vault")
		return AppContext{}, errors.New("invalid vault")
	}
	return AppContext{
		client,
		vault,
	}, nil
}
