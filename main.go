/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"password-manager/cmd"
)

func main() {
	// vault := pkg.Vault{}
	// vaultData, err := json.Marshal(vault)
	// if err != nil {
	// 	log.Fatalf("Error marshaling user: %v", err)
	// }

	// err = client.Set(ctx, vault.Name, vaultData, 0).Err()
	// if err != nil {
	// 	log.Fatalf("Error storing user data in Redis: %v", err)
	// }

	// val, err := client.Get(ctx, vault.Name).Bytes()
	// if err != nil {
	// 	panic(err)
	// }

	// var storedVault pkg.Vault
	// if err := json.Unmarshal(val, &storedVault); err != nil {
	// 	log.Fatalf("Error unmarshaling user: %v", err)
	// }
	// fmt.Println("Vault", storedVault.Name)
	cmd.Execute()
}
