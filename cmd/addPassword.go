package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/spf13/cobra"
)

var URLFlag string

var Add = &cobra.Command{
	Use:   "add",
	Short: "Add a password to the vault.",
	Run: func(cmd *cobra.Command, args []string) {
		backend_url := viperEnvVariable("BACKEND_URL")

		encryptionPassword := PromptForPass("Enter the Encryption password")
		username := PromptWithUI("Enter the Application username")
		password := PromptWithUI("Enter the Application password")
		URL := PromptWithUI("Enter the Application URL or Name")
		masterPassword,token := Read()

		if username == "" || URL == "" || encryptionPassword == "" || masterPassword == "" {
			fmt.Println("Error: Missing required fields. (EncryptionPassword / MasterPassword)")
			return
		}
		if password == "" {
			password = RandomPassword()
			fmt.Println("Password not mentioned: Using", password)
		}

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

		var newPasswords []Password
		newPasswords = vault.Passwords

		encryptedAESPassword, err := convertToAES(map[string]interface{}{
			"inputJSON": password,
		}, encryptionPassword)
		if err != nil {
			fmt.Println("Error encrypting password:", err)
			return
		}

		passwordPayload := Password{
			Id:       len(vault.Passwords) + 1,
			Name:     URL,
			Username: username,
			Password: encryptedAESPassword["aesString"],
			Salt:     encryptedAESPassword["salt"],
		}
		newPasswords = append(newPasswords, passwordPayload)

		newPasswordsJSON, err := json.Marshal(newPasswords)
		if err != nil {
			fmt.Println("Error marshalling new passwords:", err)
			return
		}

		inputJSON := map[string]interface{}{
			"passwords": json.RawMessage(newPasswordsJSON),
			"username":  username,
		}

		encryptedVault, err := convertToAES(map[string]interface{}{
			"inputJSON": inputJSON,
		}, masterPassword)
		if err != nil {
			fmt.Println("Error encrypting vault:", err)
			return
		}

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
			fmt.Println("Error:", err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			fmt.Println("Error:", res.Status)
			return
		}

		fmt.Println("Password added successfully.")
	},
}