package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

var Decrypt = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a password from the vault.",
	Run: func(cmd *cobra.Command, args []string) {
		
		backend_url := viperEnvVariable("BACKEND_URL")

		encryptionPassword := PromptForPass("Enter the Encryption Password")

		// Getting the Names for the passwords
		var vault InputJSON
		vault,err := GetVault()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		name := PromptWithMultipleLabels()
		fmt.Println(name)

		if name == "" || encryptionPassword == "" {
			fmt.Println("Error: Missing required fields. (encryption_password/name)")
			return
		}

		_,token := Read()

		// Prepare the payload
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

		for _, pass := range vault.Passwords {
			if pass.Name == name {
				decryptedPassword, err := DecryptAES(pass.Password, encryptionPassword, pass.Salt)
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