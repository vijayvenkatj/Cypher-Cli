package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/spf13/cobra"
)

var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a password from the vault.",
	Run: func(cmd *cobra.Command, args []string) {
		
		master_password,token := Read()
		backend_url := viperEnvVariable("BACKEND_URL")
		username := GetUsername()

		encryptionPassword := PromptForPass("Enter your Encryption Password")

		name := PromptWithMultipleLabels()

		if username == "" || master_password == "" || name == "" || encryptionPassword == "" {
			fmt.Println("Error: Missing required fields (master_password / encryption_password / name).")
			return
		}

		// Verify encryption password
		verificationPayload := map[string]interface{}{
			"encryption_password": createHash(encryptionPassword),
			"token":               token,
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
			"token":     token,
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