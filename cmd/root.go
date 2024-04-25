/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"password-manager/pkg"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

type AppContext struct {
	Client *redis.Client
	Vault  *pkg.Vault
}

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
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		vault := ctx.Value("vault").(*pkg.Vault)
		fmt.Println("Vault accessed within Run:", vault)
	},
}

var appContext *AppContext

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	ctx := context.Background()
	client := initRedisClient()

	vault, err := pkg.LoadVaultFromRedis(client)
	if err != nil {
		return fmt.Errorf("failed to load vault: %v", err)
	}
	appContext = &AppContext{
		Client: client,
		Vault:  vault,
	}
	ctx = context.WithValue(ctx, "vault", vault) //I dont think this is necessary
	ctx = context.WithValue(ctx, "client", client)
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.password-manager.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("vault", "v", "", "Vault")
}

func initRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // default db
	})
	return client
}
