package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "os"
    "path/filepath"
)

func viperEnvVariable(key string) string {
	viper.SetConfigFile(filepath.Join(os.Getenv("HOME"), ".cypher-cli", ".env"))
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("Error while reading config file:", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		fmt.Println("Invalid type assertion")
	}

	return value
}


var rootCmd = &cobra.Command{
    Use:   "Cypher",
    Short: "Cypher is a CLI tool for performing basic operations",
    Long:  "Cypher is a CLI tool for performing basic operations - Normal Greeting etc.",
}

var usernameFlag string
var masterPasswordFlag string
var encryptionPasswordFlag string
var emailFlag string

func Execute() {
    rootCmd.AddCommand(Register)
    rootCmd.AddCommand(Login)
    rootCmd.AddCommand(Delete)
    rootCmd.AddCommand(Add)
    rootCmd.AddCommand(Show)
    rootCmd.AddCommand(Decrypt)

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Oops. An error occurred while executing Cypher: '%s'\n", err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().StringVarP(&usernameFlag, "username", "u", "", "Specify the Username to use for Master login.")
    rootCmd.PersistentFlags().StringVarP(&masterPasswordFlag, "master-password", "m", "", "Specify the Master Password to use for Master login. (DO NOT FORGET THIS PASSWORD)")
    rootCmd.PersistentFlags().StringVarP(&encryptionPasswordFlag, "encryption-password", "e", "", "Specify the Encryption Password.")
    rootCmd.PersistentFlags().StringVarP(&emailFlag, "email", "l", "", "Specify the Email.")
}
