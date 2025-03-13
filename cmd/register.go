package cmd

import (
	"crypto/sha256"
	"fmt"
	"github.com/spf13/cobra"
	"encoding/hex"
	"encoding/json"
	"bytes"
	"net/http"
)

type Vault struct {
	passwords []string
	username string
}

type RegisterPayload struct {
	Username          string `json:"username"`
	Email             string `json:"email"`
	MasterPassword    string `json:"master_password"`
	EncryptionPassword string `json:"encryption_password"`
	Vault             map[string]string `json:"vault"`
}

func createHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
var Register = &cobra.Command{
	Use:   "register",
	Short: "Register using Master credentials.",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		email, _ := cmd.Flags().GetString("email")
		master_password, _ := cmd.Flags().GetString("master-password")
		encryption_password, _ := cmd.Flags().GetString("encryption-password")
		backend_url := viperEnvVariable("BACKEND_URL")

		if(username == "" || email == "" || master_password == "" || encryption_password == ""){
			fmt.Println("Error: Missing required fields.")
			return
		}

		master_hash := createHash(master_password)
		hash := createHash(encryption_password)

		vault := Vault{
			passwords: []string{},
			username:  username,
		}

		inputJSON := map[string]interface{}{
			"passwords": vault.passwords,
			"username":  vault.username,
		}

		encryptedVault, err := convertToAES(map[string]interface{}{
			"inputJSON": inputJSON,
		}, master_password)

		if err != nil {
			fmt.Println("Error encrypting vault:", err)
			return
		}

		payload := RegisterPayload{
			Username:          username,
			Email:             email,
			MasterPassword:    master_hash,
			EncryptionPassword: hash,
			Vault:             encryptedVault,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshalling payload:", err)
			return
		}

		res, err := http.Post(fmt.Sprintf("%s/register", backend_url), "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error sending request")
			return
		}
		defer res.Body.Close()
		
		if(res.StatusCode != 201){
			if(res.StatusCode == 409){
				fmt.Println("Error: User already exists.")
				return
			} else {
				fmt.Println("Error: ",res.Status)
				return
			}
		}
		

		fmt.Println("Successfully registered.")
	},
}
