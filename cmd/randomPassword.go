package cmd

import (
	"math/rand"
	"time"
)
func PasswordGenerator(length int) string {
	const (
		lowercase  = "abcdefghijklmnopqrstuvwxyz"
		uppercase  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numeric    = "0123456789"
		punctuation = "!@#$%^&*()_+~`|}{[]:;?><,./-="
	)

	all := lowercase + uppercase + numeric + punctuation
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := make([]byte, length)

	for i := range password {
		password[i] = all[randSource.Intn(len(all))]
	}

	return string(password)
}

func RandomPassword()(string) {
	length := 10
	password := PasswordGenerator(length)
	return password
}
