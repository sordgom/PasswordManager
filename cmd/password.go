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

// passwordCmd represents the password command
var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Manage your passwords",
	Long: `Usage: Konache password [command]
	command: [add, put, list, get, del]`,
	Run: func(cmd *cobra.Command, args []string) {
		client := appContext.Client
		vault := cmd.Context().Value("vault").(*pkg.Vault)
		if vault == nil {
			fmt.Println("Failed to retrieve vault from context")
			return
		}

		addFlag, _ := cmd.Flags().GetBool("add")
		modifyFlag, _ := cmd.Flags().GetBool("put")
		listFlag, _ := cmd.Flags().GetBool("list")
		getFlag, _ := cmd.Flags().GetBool("get")
		delFlag, _ := cmd.Flags().GetBool("del")

		if addFlag {
			fmt.Println("Adding password")
			if len(args) != 5 {
				log.Fatal("Please provide the name, url, username, password and hint for the password")
				return
			}

			password := vault.NewPassword(args[0], args[1], args[2], args[3], args[4])
			vault.AppendPassword(password)

			err := pkg.SaveVaultToRedis(client, vault)
			if err != nil {
				log.Fatalf("Failed to save vault: %v", err)
			}

			fmt.Printf("\nPassword was added successfully")
		}
		if listFlag {
			fmt.Println("Listing all passwords from Vault", vault.Name)
			// Fix the formatting
			fmt.Println("Name | Url | Hint")
			passwords := vault.Passwords
			for _, password := range passwords {
				fmt.Println(password.Name, " | ", password.Url, " | ", password.Hint, " | ")
			}
		}
		if modifyFlag {
			fmt.Println("Updating password")
			if len(args) != 3 {
				log.Fatal("Please provide the password name, master password and new password")
				return
			}

			err := vault.UpdatePassword(args[0], args[1], args[2])
			if err != nil {
				log.Fatalf("\nFailed to update password: %s", args[0])
			}

			err = pkg.SaveVaultToRedis(client, vault)
			if err != nil {
				log.Fatalf("Failed to save vault: %v", err)
			}

			fmt.Printf("\nPassword was updated successfully")
		}
		if getFlag {
			fmt.Println("Getting password")
		}
		if delFlag {
			fmt.Println("Deleting password")
		}
	},
}

func init() {
	rootCmd.AddCommand(passwordCmd)

	passwordCmd.Flags().BoolP("add", "a", false, "Add a new password to the vault")
	passwordCmd.Flags().BoolP("list", "l", false, "List all passwords in the vault")
	passwordCmd.Flags().BoolP("put", "p", false, "Modify an existing password in the vault")
	passwordCmd.Flags().BoolP("get", "g", false, "Get a specific password from the vault")
	passwordCmd.Flags().BoolP("del", "d", false, "Delete a password from the vault")
}
