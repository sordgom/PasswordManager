/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"net/http"

	"github.com/spf13/cobra"
)

type AppContext struct {
	Client http.Client
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
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	ctx := context.Background()
	appContext = &AppContext{
		Client: *http.DefaultClient,
	}
	ctx = context.WithValue(ctx, "client", appContext.Client)
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
