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
    Short: "Cypher - Password Manager. ",
    Long:  "Cypher is an Open Source all on client cloud Password Manager. ",
}

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