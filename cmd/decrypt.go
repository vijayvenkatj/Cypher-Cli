package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path/filepath"
)

var Decrypt = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a password from the vault.",
	Run: func(cmd *cobra.Command, args []string) {
		backend_url := viperEnvVariable("BACKEND_URL")
		encyptionPassword, _ := cmd.Flags().GetString("encryption-password")
		name, _ := cmd.Flags().GetString("name")

		if name == "" || encyptionPassword == "" {
			fmt.Println("Error: Missing required fields. (encryption_password/name)")
			return
		}

		// Get the user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		configDir := filepath.Join(homeDir, ".cypher-cli")

		// Define new config file path
		tokenPath := filepath.Join(configDir, "token.txt")

		// Read the token file
		tokenBytes, err := os.ReadFile(tokenPath)
		if err != nil {
			fmt.Println("Error reading token file:", err)
			return
		}

		// Prepare the payload
		verificationPayload := map[string]interface{}{
			"encryption_password": createHash(encyptionPassword),
			"token":               string(tokenBytes),
		}
		verificationPayloadJSON, err := json.Marshal(verificationPayload)
		if err != nil {
			fmt.Println("Error marshalling verification payload:", err)
			return
		}

		res, err := http.Post(fmt.Sprintf("%s/verifyEncryptionPassword", backend_url), "application/json", bytes.NewBuffer(verificationPayloadJSON))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			fmt.Println("Error:", res.Status)
			return
		}

		vault, err := GetVault()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, pass := range vault.Passwords {
			if pass.Name == name {
				decryptedPassword, err := DecryptAES(pass.Password, encyptionPassword, pass.Salt)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				fmt.Println("Password:", decryptedPassword.InputJSON)
				return
			}
		}

		fmt.Println("Error: Password not found.")
	},
}

func init() {
	Decrypt.Flags().StringP("name", "n", "", "Name of the password to decrypt.")
}
