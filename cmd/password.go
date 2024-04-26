/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"password-manager/pkg"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// passwordCmd represents the password command
var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Manage your passwords",
	Long: `Usage: Konache password [command]
	command: [add, put, list, get, del]`,
	Run: Run,
}

func Run(cmd *cobra.Command, args []string) {
	appContext, _ := load()

	// Local flags
	addFlag, _ := cmd.Flags().GetBool("add")
	modifyFlag, _ := cmd.Flags().GetBool("put")
	listFlag, _ := cmd.Flags().GetBool("list")
	getFlag, _ := cmd.Flags().GetBool("get")
	delFlag, _ := cmd.Flags().GetBool("del")

	// Table setup
	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()

	if addFlag {
		fmt.Println("Adding password")
		if len(args) != 5 {
			log.Fatal("Please provide the name, url, username, password and hint for the password")
			return
		}

		password := appContext.Vault.NewPassword(args[0], args[1], args[2], args[3], args[4])
		appContext.Vault.AppendPassword(password)

		err := pkg.SaveVaultToRedis(appContext.Client, appContext.Vault)
		if err != nil {
			log.Fatalf("Failed to save vault: %v", err)
		}

		fmt.Printf("\nPassword was added successfully")
	}
	if listFlag {
		fmt.Println("Listing all passwords from Vault", appContext.Vault.Name)

		tbl := table.New("Name", "URL", "Hint")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		passwords := appContext.Vault.GetPasswords()
		for _, password := range passwords {
			tbl.AddRow(password[0], password[1], password[2])
		}

		tbl.Print()
	}
	if modifyFlag {
		fmt.Println("Updating password")
		if len(args) != 3 {
			log.Fatal("Please provide the password name, master password and new password")
			return
		}

		err := appContext.Vault.UpdatePassword(args[0], args[1], args[2])
		if err != nil {
			log.Fatalf("\nFailed to update password: %s", args[0])
		}

		err = pkg.SaveVaultToRedis(appContext.Client, appContext.Vault)
		if err != nil {
			log.Fatalf("Failed to save vault: %v", err)
		}

		fmt.Printf("\nPassword was updated successfully")
	}
	if getFlag {
		fmt.Println("Listing The password value from Vault", appContext.Vault.Name)

		if len(args) != 1 {
			log.Fatal("Please provide the password name")
			return
		}
		password, err := appContext.Vault.GetPassword(args[0])
		if err != nil {
			log.Fatalf("\nFailed to get password: %s", args[0])
		}

		// Ask user to input master password
		fmt.Print("Enter master password: ")
		MasterPassword, err := readPasswordFromStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read password: %v\n", err)
			os.Exit(1)
		}
		if !appContext.Vault.VerifyMasterPassword(MasterPassword) {
			fmt.Println("Master password is incorrect", MasterPassword)
			return
		}

		fmt.Println("Password:", password.ReadPassword(MasterPassword))
	}
	if delFlag {
		fmt.Println("Deleting password")
	}
}

func init() {
	rootCmd.AddCommand(passwordCmd)

	passwordCmd.Flags().BoolP("add", "a", false, "Add a new password to the vault")
	passwordCmd.Flags().BoolP("list", "l", false, "List all passwords in the vault")
	passwordCmd.Flags().BoolP("put", "p", false, "Modify an existing password in the vault")
	passwordCmd.Flags().BoolP("get", "g", false, "Get a specific password from the vault")
	passwordCmd.Flags().BoolP("del", "d", false, "Delete a password from the vault")
	passwordCmd.Flags().StringP("vault", "v", "", "Vault name")
}

func readPasswordFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return password[:len(password)-1], nil
}
