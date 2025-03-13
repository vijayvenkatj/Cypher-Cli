package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
)

var deletepasswordFlag string

var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a password from the vault.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		configDir := filepath.Join(homeDir, ".cypher-cli")

		// Define new config file paths
		usernamePath := filepath.Join(configDir, "username.txt")
		masterPasswordPath := filepath.Join(configDir, "master_password.txt")
		tokenPath := filepath.Join(configDir, "token.txt")

		usernameBytes, err := os.ReadFile(usernamePath)
		if err != nil {
			fmt.Println("Error reading username:", err)
			return
		}
		backend_url := viperEnvVariable("BACKEND_URL")
		username := string(usernameBytes)

		masterBytes, err := os.ReadFile(masterPasswordPath)
		if err != nil {
			fmt.Println("Please login to continue.")
			return
		}
		master_password := string(masterBytes)

		encyptionPassword, _ := cmd.Flags().GetString("encryption-password")
		name, _ := cmd.Flags().GetString("name")

		if username == "" || master_password == "" || name == "" || encyptionPassword == "" {
			fmt.Println("Error: Missing required fields (master_password / encryption_password / name).")
			return
		}

		tokenBytes, err := os.ReadFile(tokenPath)
		if err != nil {
			fmt.Println("Error reading token file:", err)
			return
		}

		// Verify encryption password
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
			fmt.Println("Error verifying encryption password:", err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			fmt.Println("Error:", res.Status)
			return
		}

		// Get current vault
		vault, err := GetVault()
		if err != nil {
			fmt.Println("Error retrieving vault:", err)
			return
		}

		passwords := vault.Passwords
		var newPasswords []Password
		found := false

		for _, pass := range passwords {
			if pass.Name == name {
				found = true
				continue
			}
			newPasswords = append(newPasswords, pass)
		}

		if !found {
			fmt.Println("Error: Password entry not found for", name)
			return
		}

		vault.Passwords = newPasswords
		fmt.Println("Deleted password entry for:", name)

		// Encrypt updated vault
		inputJSON := map[string]interface{}{
			"passwords": newPasswords,
			"username":  username,
		}

		encryptedVault, err := convertToAES(map[string]interface{}{
			"inputJSON": inputJSON,
		}, master_password)
		if err != nil {
			fmt.Println("Error encrypting vault:", err)
			return
		}

		// Prepare and send update request
		payload := map[string]interface{}{
			"aesString": encryptedVault["aesString"],
			"salt":      encryptedVault["salt"],
			"username":  username,
			"token":     string(tokenBytes),
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshalling payload:", err)
			return
		}

		res, err = http.Post(fmt.Sprintf("%s/updateVault", backend_url), "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error updating vault:", err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			fmt.Println("Error updating vault:", res.Status)
			return
		}

		fmt.Println("Password entry deleted successfully.")
	},
}

func init() {
	Delete.Flags().StringVarP(&deletepasswordFlag, "name", "n", "", "URL/Name for the password to be deleted")
}
