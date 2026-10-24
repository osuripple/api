package common

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomString returns a cryptographically secure lowercase hexadecimal string.
func RandomString(n int) (string, error) {
	b := make([]byte, n/2+n%2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b)[:n], nil
}
