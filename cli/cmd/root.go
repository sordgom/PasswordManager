/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/sordgom/PasswordManager/cli/pkg"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

type AppContext struct {
	Client *redis.Client
	Vault  *pkg.Vault
}

var appContext *AppContext

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Konache",
	Short: "Konache is a simple Password Manager",
	Long: `
Konache is a  is a CLI Password manager that allows users to:
- Create new password vaults with a master password.
- Add new passwords to the vaults.
- List all the passwords in a vault.
- Retrieve a password from a vault.
PS: we never store or know about your passwords.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		vaultFlag, _ := cmd.Flags().GetString("vault")
		if vaultFlag == "" {
			fmt.Println("Please specify a vault with -v or --vault")
			return
		}
		vault, err := pkg.LoadVaultFromRedis(appContext.Client, vaultFlag)
		if err != nil {
			fmt.Errorf("failed to load vault: %v", err)
			return
		}
		appContext.Vault = vault
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	ctx := context.Background()
	client := initRedisClient()
	appContext = &AppContext{
		Client: client,
	}
	ctx = context.WithValue(ctx, "client", client)
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.password-manager.yaml)")
	rootCmd.PersistentFlags().StringP("vault", "v", "", "Vault")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // default db
	})
	return client
}
