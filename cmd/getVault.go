package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type VaultResponce struct {
	Message   string
	Username  string
	AesString string
	Salt      string
}

type InputJSON struct {
	Passwords []Password `json:"passwords"`
	Username  string     `json:"username"`
}

type Password struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type VaultData struct {
	InputJSON InputJSON `json:"inputJSON"`
}

func GetVault() (InputJSON, error) {
	backend_url := viperEnvVariable("BACKEND_URL")

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return InputJSON{}, err
	}
	configDir := filepath.Join(homeDir, ".cypher-cli")

	// Define new config file paths
	tokenPath := filepath.Join(configDir, "token.txt")
	masterPasswordPath := filepath.Join(configDir, "master_password.txt")

	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		fmt.Println("Please Login to continue.")
		return InputJSON{}, err
	}

	masterBytes, err := os.ReadFile(masterPasswordPath)
	if err != nil {
		fmt.Println("Please Login to continue.")
		return InputJSON{}, err
	}

	payload, err := json.Marshal(map[string]string{
		"token": string(tokenBytes),
	})
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return InputJSON{}, err
	}

	res, err := http.Post(fmt.Sprintf("%s/getVault", backend_url), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error:", err)
		return InputJSON{}, err
	}
	defer res.Body.Close()

	var vaultData VaultResponce
	err = json.NewDecoder(res.Body).Decode(&vaultData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return InputJSON{}, err
	}

	var decryptedDataMap DecryptOutput
	decryptedDataMap, err = DecryptAES(vaultData.AesString, string(masterBytes), vaultData.Salt)
	if err != nil {
		if err.Error() == "Decryption failed." {
			fmt.Println("Error: Decryption failed.")
		}
		return InputJSON{}, err
	}

	var decryptedData VaultData
	decryptedDataBytes, err := json.Marshal(decryptedDataMap)
	if err != nil {
		fmt.Println("Error marshalling decrypted data:", err)
		return InputJSON{}, err
	}
	err = json.Unmarshal(decryptedDataBytes, &decryptedData)
	if err != nil {
		fmt.Println("Error unmarshalling decrypted data:", err)
		return InputJSON{}, err
	}

	return decryptedData.InputJSON, nil
}
