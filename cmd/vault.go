/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
		shouldGenerate, _ := cmd.Flags().GetBool("generate")

		if len(args) > 2 || (len(args) == 1 && !shouldGenerate) || len(args) == 0 {
			fmt.Println("Please provide a name and a master password for the vault")
			return
		}
		var hashedPassword string
		if shouldGenerate {
			hashedPassword = pkg.GenerateHash()
		} else {
			hashedPassword = args[1]
		}
		vault := pkg.CreateVault(args[0], hashedPassword)
		fmt.Printf("vault %s created", vault.Name)
	},
}

func init() {
	vaultCmd.Flags().BoolP("generate", "g", false, "generate a master password")
	rootCmd.AddCommand(vaultCmd)
}
