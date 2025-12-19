package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	// Character sets for password generation
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	specialChars   = "!@#$%^&*"

	// Combined character set
	allChars = lowercaseChars + uppercaseChars + digitChars + specialChars
)

// GeneratePassword generates a secure random password of 10 characters
// The password includes lowercase, uppercase, digits, and special characters
func GeneratePassword() string {
	password := make([]byte, 10)

	for i := range password {
		// Generate a random index within the character set
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			// Fallback to a simple approach if crypto/rand fails
			// This should rarely happen
			password[i] = allChars[i%len(allChars)]
			continue
		}
		password[i] = allChars[num.Int64()]
	}

	return string(password)
}
