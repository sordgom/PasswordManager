/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"password-manager/pkg"

	"github.com/spf13/cobra"
)

// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Create a new vault with a master password",
	Long: `Usage: Konache vault [vault_name] [master_password]
	vault_name: the name of the vault
	master_password: the master password for the vault`,
	Run: func(cmd *cobra.Command, args []string) {
		client := appContext.Client
		vaultFlag, _ := cmd.Flags().GetString("vault")

		shouldGenerate, _ := cmd.Flags().GetBool("generate")
		if len(args) > 1 || (len(args) == 0 && !shouldGenerate) {
			fmt.Println("Please provide a name and a master password for the vault")
			return
		}

		var hashedPassword string //Should be hashed first
		if shouldGenerate {
			hashedPassword = pkg.GenerateRandomHash()
		} else {
			hashedPassword = pkg.GenerateHash(args[0])
		}

		vault := pkg.New(vaultFlag, hashedPassword)

		err := pkg.SaveVaultToRedis(client, &vault)
		if err != nil {
			log.Fatalf("Failed to save vault: %v", err)
		}

		fmt.Printf("vault %s created", vault.Name)
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)

	vaultCmd.Flags().BoolP("generate", "g", false, "generate a master password")
}
