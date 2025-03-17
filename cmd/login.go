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

type LoginResponceType struct {
	Message string
	Token   string
}

var Login = &cobra.Command{
	Use:   "login",
	Short: "Login using Master credentials.",
	Run: func(cmd *cobra.Command, args []string) {
		username := PromptWithUI("Enter the Username")
		master_password := PromptForPass("Enter the Master Password")
		if username == "" || master_password == "" {
			fmt.Println("Error: Missing required fields.")
			return
		}

		backend_url := viperEnvVariable("BACKEND_URL")
		master_hash := createHash(master_password)

		// Get home directory for storing config files
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		configDir := filepath.Join(homeDir, ".cypher-cli")

		// Ensure config directory exists
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			fmt.Println("Error creating config directory:", err)
			return
		}

		// Define file paths
		usernamePath := filepath.Join(configDir, "username.txt")
		tokenPath := filepath.Join(configDir, "token.txt")
		masterPasswordPath := filepath.Join(configDir, "master_password.txt")

		// Write username to file
		err = os.WriteFile(usernamePath, []byte(username), 0600)
		if err != nil {
			fmt.Println("Error writing to username file:", err)
			return
		}

		payload, err := json.Marshal(map[string]string{
			"username":        username,
			"master_password": master_hash,
		})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		res, err := http.Post(fmt.Sprintf("%s/login", backend_url), "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			fmt.Println("Error:", res.Status)
			return
		}

		var LoginResponce LoginResponceType
		err = json.NewDecoder(res.Body).Decode(&LoginResponce)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// Write token and master password to files
		err = os.WriteFile(tokenPath, []byte(LoginResponce.Token), 0600)
		if err != nil {
			fmt.Println("Error writing to token file:", err)
			return
		}

		err = os.WriteFile(masterPasswordPath, []byte(master_password), 0600)
		if err != nil {
			fmt.Println("Error writing to master password file:", err)
			return
		}
		
		fmt.Println(LoginResponce.Message)
	},
}
