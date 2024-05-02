/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	client := appContext.Client

	// Local flags
	addFlag, _ := cmd.Flags().GetBool("add")
	modifyFlag, _ := cmd.Flags().GetBool("put")
	listFlag, _ := cmd.Flags().GetBool("list")
	getFlag, _ := cmd.Flags().GetBool("get")
	delFlag, _ := cmd.Flags().GetBool("del")

	// Global flags
	vaultName, _ := cmd.Flags().GetString("vault")

	// Table setup
	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()

	if addFlag {
		fmt.Println("Adding password")
		if len(args) != 5 {
			log.Fatal("Please provide the name, url, username, password and hint for the password")
			return
		}

		req := strings.NewReader(fmt.Sprintf(`{"name": "%s", "url": "%s", "username": "%s", "password": "%s", "hint": "%s"}`, args[0], args[1], args[2], args[3], args[4]))
		res, err := client.Post(fmt.Sprintf(`http://localhost:8080/password?vault_name=%s`, vaultName), "application/json", req)
		if err != nil {
			log.Fatalf("Failed to add password: %v", err)
		}
		fmt.Printf("\nPassword was added successfully")
		defer res.Body.Close()
	}
	if listFlag {
		fmt.Println("Listing all passwords from Vault", vaultName)
		time.Sleep(1 * time.Second)

		tbl := table.New("Name", "URL", "Hint")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		res, err := client.Get(fmt.Sprintf(`http://localhost:8080/passwords?vault_name=%s`, vaultName))
		if err != nil || res.StatusCode != 200 {
			log.Fatalf("Failed to list passwords: %v", err)
		}
		defer res.Body.Close()

		var passwords []getPasswordResponse
		if err := json.NewDecoder(res.Body).Decode(&passwords); err != nil {
			log.Fatal("Error decoding JSON response:", err)
		}

		for _, password := range passwords {
			tbl.AddRow(password.Name, password.Password, password.Hint)
		}

		tbl.Print()
	}
	if getFlag {
		fmt.Println("Listing The password value from Vault", vaultName)
		if len(args) != 1 {
			log.Fatal("Please provide the password name")
			return
		}

		// Verify master password
		fmt.Print("Enter master password: ")
		masterPassword, err := readPasswordFromStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read password: %v\n", err)
			os.Exit(1)
		}
		req := strings.NewReader(fmt.Sprintf(`{"vault_name": "%s", "master_password": "%s"}`, vaultName, masterPassword))
		res, err := client.Post("http://localhost:8080/vault/verify", "application/json", req)
		if err != nil {
			log.Fatalf("Failed to Verify master password: %v", err)
		}
		if res.StatusCode != 200 {
			log.Fatalf("Failed to Verify master password: %v", res.Status)
		}

		// Fetch a single password from the vault
		res, err = client.Get(fmt.Sprintf(`http://localhost:8080/password?vault_name=%s&password_name=%s`, vaultName, args[0]))
		if err != nil || res.StatusCode != 200 {
			log.Fatalf("Failed to list passwords: %v", err)
		}
		defer res.Body.Close()
		var password getPasswordResponse
		if err := json.NewDecoder(res.Body).Decode(&password); err != nil {
			log.Fatal("Error decoding JSON response:", err)
		}

		fmt.Println("Password:", password.Password)
	}
	if modifyFlag {
		// fmt.Println("Updating password")
		// if len(args) != 1 {
		// 	log.Fatal("Please provide the password name")
		// 	return
		// }

		// // Ask user to input master password
		// fmt.Print("Enter master password: ")
		// MasterPassword, err := readPasswordFromStdin()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to read password: %v\n", err)
		// 	os.Exit(1)
		// }
		// if !appContext.Vault.VerifyMasterPassword(MasterPassword) {
		// 	fmt.Println("Master password is incorrect", MasterPassword)
		// 	return
		// }

		// // Ask user to input new password and to confirm it
		// fmt.Print("Enter new password: ")
		// NewPassword, err := readPasswordFromStdin()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to read new password: %v\n", err)
		// 	os.Exit(1)
		// }
		// fmt.Print("Confirm new password: ")
		// ConfirmNewPassword, err := readPasswordFromStdin()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to read confirm password: %v\n", err)
		// 	os.Exit(1)
		// }

		// err = appContext.Vault.UpdatePassword(args[0], NewPassword, ConfirmNewPassword)
		// if err != nil {
		// 	log.Fatalf("\nFailed to update password: %s", args[0])
		// }

		// err = pkg.SaveVaultToRedis(appContext.Client, appContext.Vault)
		// if err != nil {
		// 	log.Fatalf("Failed to save vault: %v", err)
		// }

		// fmt.Printf("\nPassword was updated successfully")
	}
	if delFlag {
		// fmt.Println("Deleting password")
		// if len(args) != 1 {
		// 	log.Fatal("Please provide the password name")
		// 	return
		// }

		// // Ask user to input master password
		// fmt.Print("Enter master password: ")
		// MasterPassword, err := readPasswordFromStdin()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to read password: %v\n", err)
		// 	os.Exit(1)
		// }
		// if !appContext.Vault.VerifyMasterPassword(MasterPassword) {
		// 	fmt.Println("Master password is incorrect", MasterPassword)
		// 	return
		// }

		// err = appContext.Vault.DeletePassword(args[0])
		// if err != nil {
		// 	log.Fatalf("\nFailed to delete password: %s", args[0])
		// }
		// fmt.Println("Password was deleted successfully")

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

type getPasswordResponse struct {
	Name     string `json:"name"`
	Hint     string `json:"hint"`
	Password string `json:"password"`
}
