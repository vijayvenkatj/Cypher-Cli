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

var URLFlag string

var Add = &cobra.Command{
	Use:   "add",
	Short: "Add a password to the vault.",
	Run: func(cmd *cobra.Command, args []string) {
		backend_url := viperEnvVariable("BACKEND_URL")
		encyptionPassword, _ := cmd.Flags().GetString("encryption-password")

		// Get the user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		configDir := filepath.Join(homeDir, ".cypher-cli")

		// Define new config file paths
		masterPath := filepath.Join(configDir, "master_password.txt")
		tokenPath := filepath.Join(configDir, "token.txt")

		// Read master password
		masterBytes, err := os.ReadFile(masterPath)
		if err != nil {
			fmt.Println("Please Login to continue.")
			return
		}
		masterPassword := string(masterBytes)

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		URL, _ := cmd.Flags().GetString("URL")

		if username == "" || URL == "" || encyptionPassword == "" || masterPassword == "" {
			fmt.Println("Error: Missing required fields. (EncryptionPassword / MasterPassword)")
			return
		}
		if password == "" {
			password = RandomPassword()
			fmt.Println("Password not mentioned: Using", password)
		}

		// Read token file
		tokenBytes, err := os.ReadFile(tokenPath)
		if err != nil {
			fmt.Println("Error reading token file:", err)
			return
		}

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

		var newPasswords []Password
		newPasswords = vault.Passwords

		encryptedAESPassword, err := convertToAES(map[string]interface{}{
			"inputJSON": password,
		}, encyptionPassword)
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
			"token":     string(tokenBytes),
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

func init() {
	Add.Flags().StringP("URL", "U", "", "URL to send the password to.")
	Add.MarkFlagRequired("URL")
	Add.Flags().StringP("username", "u", "", "Specify the username for the website.")
	Add.MarkFlagRequired("username")
	Add.Flags().StringP("password", "p", "", "Specify the password for the website. (Ignore if you need a random password assigned)")
}
