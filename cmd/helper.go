package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"os/exec"
	"os"
	"path/filepath"
	"github.com/manifoldco/promptui"
)


func convertToAES(data map[string]interface{}, key string) (map[string]string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("node", "./cmd/encrypt.js", string(jsonData), key)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("JavaScript Error:", stderr.String())
		return nil, err
	}

	var result map[string]string
	err = json.Unmarshal(out.Bytes(), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type DecryptInput struct {
	AESString string `json:"aesString"`
	Key       string `json:"key"`
	Salt      string `json:"salt"`
}

type DecryptOutput struct {
	InputJSON interface{} `json:"inputJSON"`
}



func DecryptAES(aesString string, key string, salt string) (DecryptOutput, error) {
	cmd := exec.Command("node", "cmd/decrypt.js", aesString, key, salt)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return DecryptOutput{}, fmt.Errorf("%s", stderr.String())
	}

	cleanOutput := strings.TrimSpace(out.String())

	var result DecryptOutput
	err = json.Unmarshal([]byte(cleanOutput), &result)
	if err != nil {
		return DecryptOutput{}, fmt.Errorf("JSON unmarshal error: %v", err)
	}

	return result, nil
}



func Read() (string,string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return "",""
	}
	configDir := filepath.Join(homeDir, ".cypher-cli")
	masterPath := filepath.Join(configDir, "master_password.txt")
	tokenPath := filepath.Join(configDir, "token.txt")

	// Read master password
	masterBytes, err := os.ReadFile(masterPath)
	if err != nil {
		fmt.Println("Please Login to continue.")
		return "",""
	}
	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		fmt.Println("Error reading token file:", err)
		return "",""
	}
	masterPassword := string(masterBytes)
	token := string(tokenBytes)

	return masterPassword,token
}

func GetUsername() (string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return ""
	}
	configDir := filepath.Join(homeDir, ".cypher-cli")
	usernamePath:= filepath.Join(configDir, "username.txt")

	usernameBytes,err := os.ReadFile(usernamePath)
	if err != nil {
		fmt.Println("Please login to continue.")
		return ""
	}

	return string(usernameBytes)
}


func PromptWithUI(label string) (string) {
	prompt := &promptui.Prompt{
		Label: label,
	}
	data,err := prompt.Run()
	if err != nil {
		fmt.Println("Error getting user input.")
		return ""
	}
	return data
}

func PromptForPass(label string) (string) {
	prompt := &promptui.Prompt{
		Label: label,
		Mask: '*',
	}
	data,err := prompt.Run()
	if err != nil {
		fmt.Println("Error getting password.")
		return ""
	}
	return data
}

func PromptWithMultipleLabels() (string) {

	// Getting the Names for the passwords
	var vault InputJSON
	vault,err := GetVault()
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	var items []string

	for _,pass := range vault.Passwords {
		items = append(items,pass.Name)
	}
	prompt := &promptui.Select{
		Label: "Select your Password",
		Items: items,
	}
	_,data,err := prompt.Run()
	if err != nil {
		fmt.Println("Error getting user input.")
		return ""
	}
	return data
}