package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"os/exec"
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